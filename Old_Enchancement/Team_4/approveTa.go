package Team_4

//Done by Sushmitha
import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type TaEmployee struct {
	ID             int       `json:"id"`
	CountryID      int       `json:"countryID"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	DOJ            time.Time `json:"doj"`
	OfficeMailID   string    `json:"officeMailId"`
	PersonalMailID string    `json:"personalMailId"`
	ContactNum     string    `json:"contactNum"`
	Address        string    `json:"address"`
	Address2       string    `json:"address2"`
	Address3       string    `json:"address3"`
	Pincode        string    `json:"pincode"`
	EmpRole        int       `json:"empRole"`
	LicenseNumber  string    `json:"license_number"`
	SupervisorID   int       `json:"supervisorId"`
	ProfilePic     string    `json:"profile_pic"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	CreatedBy      string    `json:"createdBy"`
	LastUpdatedAt  time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy  string    `json:"lastUpdatedBy"`
	WorkNum        string    `json:"workNum"`
}
type RequestBody struct {
	TAID          []string `json:"ta_id"`
	UserID        string   `json:"user_id"`
	ExtraComments string   `json:"extra_comments"`
	Status        int      `json:"status"`
}
type obj struct {
	ID                int
	EmpID             int
	POAID             int
	StartOdometer     string
	StartLocationName string
	EndOdometer       string
	EndLocationName   string
	Printing          int
	Telephone         int
	Stationery        int
	Others            int
	OtherText         string
	Kilometer         int
	Status            int
	EntryDate         string
	ApprovedBy        sql.NullInt64
	ApproveDate       string
	MailSend          sql.NullInt64
	ExtraComments     string
	TotalTA           float64
	Purpose           string
	ModeOfTravel      int
	RatePerKM         float64
	DA                int
	VerifyedBy        sql.NullInt64
	VerifyiedDate     sql.NullString
	UpdateDate        string
	ExtraDate         sql.NullString
	ExtraOMComments   sql.NullString
	TaName            string
	CheckInLocation   string
}

var res1 obj

var data RequestBody
var EmpId int

func ApproveTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"status": "404 Not Found", "Error": "Method not found"})
		return
	}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "message": err.Error()})
		return
	}
	taIDs := data.TAID
	stringOpsIds := make([]string, len(taIDs))

	// for i, id := range taIDs {
	// 	stringOpsIds[i] = id
	// }
	copy(stringOpsIds, taIDs)
	result := strings.Join(stringOpsIds, ",")

	if data.UserID != "" {

		Td := data.TAID
		status := data.Status
		if data.Status == 0 {
			status = 1
		}

		ExtraComments := data.ExtraComments
		currentTime := time.Now()
		date := currentTime.Add(5*time.Hour + 30*time.Minute).Format("2006-01-02")

		var length = len(data.TAID)
		for i := 0; i < length; i++ {
			query := fmt.Sprintf("UPDATE travelling_allowance SET status=%d, approved_by=%s, approve_date='%s', extra_comments='%s' WHERE id=%s",

				status, data.UserID, date, ExtraComments, Td[i])

			_, err := DB.Exec(query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
				return
			}
		}

		updateQuery := fmt.Sprintf("SELECT DISTINCT emp_id  FROM travelling_allowance WHERE id in (%s)", result)

		rows, err := DB.Query(updateQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}

		for rows.Next() {
			rows.Scan(&EmpId)
		}

		rows.Close()
		qury1 := fmt.Sprintf("SELECT DATE_FORMAT(entry_date, '%%M %%Y') as entry_date FROM travelling_allowance WHERE emp_id = %d AND id IN (%s)", EmpId, result)

		rows2, err2 := DB.Query(qury1)
		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		var entryDate string
		if rows2.Next() {
			if err := rows2.Scan(&entryDate); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
				return
			}
		}

		q := fmt.Sprintf("SELECT travelling_allowance.*,tbl_poa.name AS 'Ta_Name',COALESCE(tbl_poa.check_in_location,'') FROM travelling_allowance JOIN tbl_poa ON travelling_allowance.poa_id=tbl_poa.id WHERE travelling_allowance.id in(%s) and  emp_id=%d", result, EmpId)
		rows, err = DB.Query(q)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		defer rows.Close()

		var location string
		_ = location
		var cnt int = 0
		var rowsHTML, modeOfTravel string
		var approve, verify, totalDa, totalStationery, DailyAllowance, totalTelephone, totalPrinting, totalOthers, grandTotal int
		var totalKlmtr, totalTa, ratePErKm, fare_amount float64
		_, _, _, _ = approve, verify, DailyAllowance, totalStationery

		// Process fetched rows

		for rows.Next() {
			err := rows.Scan(
				&res1.ID,
				&res1.EmpID,
				&res1.POAID,
				&res1.StartOdometer,
				&res1.StartLocationName,
				&res1.EndOdometer,
				&res1.EndLocationName,
				&res1.Printing,
				&res1.Telephone,
				&res1.Stationery,
				&res1.Others,
				&res1.OtherText,
				&res1.Kilometer,
				&res1.Status,
				&res1.EntryDate,
				&res1.ApprovedBy,
				&res1.ApproveDate,
				&res1.MailSend,
				&res1.ExtraComments,
				&res1.TotalTA,
				&res1.Purpose,
				&res1.ModeOfTravel,
				&res1.RatePerKM,
				&res1.DA,
				&res1.VerifyedBy,
				&res1.VerifyiedDate,
				&res1.UpdateDate,
				&res1.ExtraDate,
				&res1.ExtraOMComments,
				&fare_amount,
				&res1.TaName,
				&res1.CheckInLocation,
			)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			if res1.ModeOfTravel == 1 {
				modeOfTravel = "BUS"
			} else if res1.ModeOfTravel == 2 {
				modeOfTravel = "CAR"
			} else if res1.ModeOfTravel == 3 {
				modeOfTravel = "BIKE"
			} else if res1.ModeOfTravel == 4 {
				modeOfTravel = "AUTO"
			} else if res1.ModeOfTravel == 5 {
				modeOfTravel = "Train"
			} else {
				modeOfTravel = ""
			}

			if res1.RatePerKM == 1 {
				ratePErKm = 3.5
			} else {
				ratePErKm = 0
			}

			if res1.DA == 1 {
				DailyAllowance = 200
			} else if res1.DA == 2 {
				DailyAllowance = 150
			} else if res1.DA == 3 {
				DailyAllowance = 100
			} else {
				DailyAllowance = 0
			}
			cnt += 1
			totalKlmtr += float64(res1.Kilometer)
			totalTa += (float64(res1.Kilometer) * ratePErKm) + fare_amount
			totalDa += int(DailyAllowance)
			totalTelephone += res1.Telephone
			totalPrinting += res1.Printing

			totalOthers += res1.Others
			grandTotal += int(res1.TotalTA)

			rowsHTML += fmt.Sprintf(`
	<tr>
		<td>%d</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%d</td>
		<td>%s</td>
		<td>%.2f</td>
		<td>%.2f</td>
		<td>%.2f</td>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%s</td>
		<td>%d</td>
		<td>%d</td>
	</tr>`, cnt, res1.EntryDate, res1.StartLocationName, res1.EndLocationName,
				res1.Purpose, res1.StartOdometer, res1.EndOdometer, res1.Kilometer, modeOfTravel, ratePErKm, fare_amount, totalTa,
				DailyAllowance, res1.Telephone, res1.Printing, res1.OtherText,
				res1.Others, grandTotal)

		}

		htmlTemplate := `
						<tr>
							<td></td>
							<td></td>
							<td></td>
							<td></td>
							<td></td>
							<td></td>
							<th>Total</th>
							<td>%.2f</td>
							<td></td>
							<td></td>
							<td>%.2f</td>
							<td>%d</td>
							<td>%d</td>
							<td>%d</td>
							<td></td>
							<td>%d</td>
							<td>%d</td>
						</tr>`

		//   HTML template for TA
		totalRowHTML := fmt.Sprintf(htmlTemplate, totalKlmtr, totalTa, totalDa, totalTelephone, totalPrinting, totalOthers, grandTotal)
		// totalRowHTML += rowsHTML
		rowsHTML += totalRowHTML
		e := DB.QueryRow("SELECT  first_name,last_name,roleName FROM employee e JOIN roles_Master rm ON e.empRole =rm.id WHERE e.status = 1 AND e.id = ?", EmpId)
		var em struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			RoleName  string `json:"roleName"`
		}
		if err := e.Scan(&em.FirstName, &em.LastName, &em.RoleName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		f := fmt.Sprintf("SELECT  e.first_name,e.last_name FROM employee e JOIN travelling_allowance ta ON e.id =ta.approved_by WHERE e.status = 1 AND ta.id in (%s)", result)
		rows3, err := DB.Query(f)
		var app struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		for rows3.Next() {
			rows3.Scan(&app.FirstName, &app.LastName)
		}

		rows3.Close()

		g := fmt.Sprintf("SELECT  e.first_name,e.last_name FROM employee e JOIN travelling_allowance ta ON e.id =ta.verifyed_by WHERE e.status = 1 AND ta.id in (%s)", result)
		rows4, err := DB.Query(g)
		var ver struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		for rows4.Next() {
			rows4.Scan(&ver.FirstName, &ver.LastName)
		}

		rows4.Close()

		// Generate HTML content
		var htmlContent strings.Builder
		htmlContent.WriteString(`
		<html>
		<head>
		<meta charset='UTF-8'>
		<meta http-equiv='X-UA-Compatible' content='IE=edge'>
		<meta name='viewport' content='width=device-width, initial-scale=1.0'>
		<title>Travel Allowance Document</title>
		<style>
			table, td, th {
				border: 1px solid;
			}
			table {
				table-layout: auto;
				width: 100%;
				border-collapse: collapse;
			}
		</style>
		</head>
		<body>
		<table>
		<tr>
			<th style='height: 50px;' colspan='6'>Buzz India Trust</th>
		</tr>
		<tr>
			<th style='height: 50px;' colspan='6'>Staff Travel Settlement Format</th>
		</tr>
		<tr style='height: 50px;'>
				<td>Staff Name:</td>
				<td>` + em.FirstName + ` ` + em.LastName + `</td>
				<td>Designation:</td>
				<td colspan='3'>` + em.RoleName + `</td>
		</tr>
		<tr style='height: 50px;'>
				<td>Approve Date:</td>
				<td>` + date + `</td>
				<td>TA Month:</td>
				<td colspan='5'>` + res1.EntryDate + `</td>
		</tr>
		</table>
		<table>
			<tr>
				<th>Sl</th>
				<th>Date</th>
				<th>From</th>
				<th>To</th>
				<th>Purpose</th>
				<th>Start Odo</th>
				<th>End Odo</th>
				<th>Total KM</th>
				<th>Mode of Travel</th>
				<th>Rate per KM</th>
				<th>Fare Amount</th>
				<th>Total TA</th>
				<th>Food Exp</th>
				<th>Tele</th>
				<th>Print & Stati</th>
				<th colspan='2'>Oth exp</th>
				<th>Total</th>
			</tr>
		` + rowsHTML + `
		</table>
		<table style='margin-top: 20px;'>
			<tr>
				<th>Submitted by</th>
				<td> ` + em.FirstName + ` ` + em.LastName + `</td>
				<th>Verified By Ops /Sr.Ops Manager</th>
				<td>` + ver.FirstName + ` ` + ver.LastName + `</td>
				<th>Approved by Programme Manager</th>
				<td>` + app.FirstName + ` ` + app.LastName + `</td>
			</tr>
			<tr>
				<td colspan='2' style='height: 50px;'></td>
				<td colspan='2' style='height: 50px;'></td>
				<td colspan='2' style='height: 50px;'></td>
			</tr>
		</table>
		</body>
		</html>
	`)

		htmlStr := htmlContent.String()

		if false {
			query := fmt.Sprintf("SELECT * FROM employee e JOIN roles_Master rm ON e.empRole = rm.id WHERE e.status = 1 AND e.id = %d", res1.ApprovedBy.Int64)

			rows, err = DB.Query(query)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
				return
			}
			var employee TaEmployee

			if rows.Next() {
				err = rows.Scan(
					&employee.ID,
					&employee.CountryID,
					&employee.FirstName,
					&employee.LastName,
					&employee.Gender,
					&employee.DOJ,
					&employee.OfficeMailID,
					&employee.PersonalMailID,
					&employee.ContactNum,
					&employee.Address,
					&employee.Address2,
					&employee.Address3,
					&employee.Pincode,
					&employee.EmpRole,
					&employee.LicenseNumber,
					&employee.SupervisorID,
					&employee.ProfilePic,
					&employee.Status,
					&employee.CreatedAt,
					&employee.CreatedBy,
					&employee.LastUpdatedAt,
					&employee.LastUpdatedBy,
					&employee.WorkNum,
				)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
					return
				}
			}
		}
		file, err := os.Create("output.html")
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}
		defer file.Close()

		_, err = file.WriteString(htmlStr)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}

		err = file.Sync()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}

		json.NewEncoder(w).Encode("Email Sent Successfully")
		err = file.Close()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}
	}
	var officeMailID string
	err = DB.QueryRow("SELECT OfficeMailId FROM employee WHERE id=?", EmpId).Scan(&officeMailID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}

	t := template.New("output.html")

	t, err = t.ParseFiles("output.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, nil); err != nil {
		log.Println(err)
	}
	results := tpl.String()

	// Convert HTML to PDF using wkhtmltopdf command
	cmd := exec.Command("wkhtmltopdf", "-O", "landscape", "-", "Travel Allowance.pdf")
	cmd.Stdin = bytes.NewBufferString(results)

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	var month, body, subject string
	err = DB.QueryRow("SELECT CONCAT(MONTHNAME(entry_date),' ', YEAR(entry_date)) FROM travelling_allowance WHERE emp_id = ?", EmpId).Scan(&month)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
		return
	}
	if res1.Status == 1 {
		body = fmt.Sprintf("Your approved TA for the month of  %s", month)
		subject = fmt.Sprintf("Your approved TA for the month of  %s", month)
	} else {
		body = fmt.Sprintf("Your TA for the month of  %s has been rejected.\n Remarks: %s ", month, data.ExtraComments)
		subject = fmt.Sprintf("Your TA for the month of  %s has been rejected", month)
	}
	mail := fmt.Sprintf("UPDATE travelling_allowance SET mail_send=1 WHERE id IN (%s)", result)
	_, err = DB.Query(mail)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
		return
	}

	// Set the email body with the month information

	m := gomail.NewMessage()
	m.SetHeader("From", "travelallowance@buzzindia.co")
	m.SetHeader("To", officeMailID)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Attach the PDF to the email
	m.Attach("Travel Allowance.pdf")

	// Send the email
	d := gomail.NewDialer("smtp.gmail.com", 587, "travelallowance@buzzindia.co", "ybsfwycrqfvemsbs")
	if err := d.DialAndSend(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400", "message": "Invalid Office Email Address", "success": false, "error": err.Error()})
		return
	}

}
