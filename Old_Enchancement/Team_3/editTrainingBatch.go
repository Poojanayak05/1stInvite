package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RequestEditTB struct {
	BatchName            string `json:"batch_name"`
	SubVillage           string `json:"sub_village"`
	ContactPerson        string `json:"contact_person"`
	NumberOfParticipants string `json:"number_of_participants"`
	Day2                 string `json:"day2"`
	Day1                 string `json:"day1"`
	TBId                 string `json:"tb_id"`
	LocationID           string `json:"location_id"`
	ContactNumber        string `json:"contact_number"`
}

type Response11 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func UpdateTrainingBatch(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var p RequestEditTB
	err1 := json.NewDecoder(r.Body).Decode(&p)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err1})
		return
	}

	id := p.TBId
	name := p.BatchName
	response := Response11{}

	tx, err := DB.Begin()
	if err != nil {
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Recovered in updateTrainingBatch:", r)
		}
	}()

	stmt, err := tx.Prepare("SELECT GROUP_CONCAT(id) as ids, location_id,user_id, name as batch_name, sub_village as sub_village, GROUP_CONCAT(DATE_FORMAT(date, '%Y-%m-%s %H:%i')) as dates,project_id FROM tbl_poa WHERE tb_id = ? AND type = '1' group by tb_id, user_id")
	if err != nil {
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	var ids, batchName, dbSubVillage, dbDates, projectID string
	var userID int
	err = stmt.QueryRow(id).Scan(&ids, &p.LocationID, &userID, &batchName, &dbSubVillage, &dbDates, &projectID)
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	err = stmt.Close()
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	// Duplicate record for date & time validation
	dates := strings.Split(dbDates, ",")
	existDay1 := dates[0]
	existDay2 := dates[1]

	// Tb can be created only if the project is active. If project endDate is less than tb date then don't allow to create
	getProjectEndDateStmt, err := tx.Prepare("SELECT endDate FROM project WHERE id = ?")
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	var projEndDate string
	err = getProjectEndDateStmt.QueryRow(projectID).Scan(&projEndDate)
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	err = getProjectEndDateStmt.Close()
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	projEndDate = projEndDate + " 23:59:59"
	t, _ := time.Parse("2006-01-02 15:04:05", p.Day2)

	projEndDateParsed, err := time.Parse("2006-01-02 15:04:05", projEndDate)
	if err != nil {
		// Handle the error if the string cannot be parsed as a valid time
		// Optionally, you can provide a default value for projEndDateParsed in case of an error.
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	// Use projEndDateParsed as a time.Time value in your code
	if t.After(projEndDateParsed) {
		// Perform the desired actions if t is after projEndDateParsed
		fmt.Println(projEndDateParsed)
	}

	// Day 1 validation
	if p.Day1 != "1970-01-01 00:00" && existDay1 != p.Day1 {
		stmt, err = tx.Prepare("SELECT count(id) as count FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'")
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		var count int
		err = stmt.QueryRow(p.Day1, userID).Scan(&count)
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		err = stmt.Close()
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		if count > 0 {
			tx.Rollback()
			response.Code = 409
			response.Message = "Day1 already exists for date & time"
			response.Success = false
			sendResponse(w, response)
			return
		}
	}

	// Day 2 validation
	if p.Day2 != "1970-01-01 00:00" && existDay2 != p.Day2 {
		stmt, err = tx.Prepare("SELECT count(id) as count FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'")
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		var count int
		err = stmt.QueryRow(p.Day2, userID).Scan(&count)
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		err = stmt.Close()
		if err != nil {
			tx.Rollback()
			response.Code = 500
			response.Message = "Internal server error"
			response.Success = false
			sendResponse(w, response)
			return
		}

		if count > 0 {
			tx.Rollback()
			response.Code = 409
			response.Message = "Day2 already exists for date & time"
			response.Success = false
			sendResponse(w, response)
			return
		}
	}

	// if val1 != val2 || p.SubVillage != "" {
	// 	name = getName(batchName, p.SubVillage)
	// } else {
	// 	name = getName(batchName, name)
	// }

	idsArr := strings.Split(ids, ",")
	id1 := idsArr[0]
	id2 := idsArr[1]

	updateQuery1 := "UPDATE tbl_poa SET name = ?, sub_village = ?, location_id = ?, participants = ?, contact_person = ?, contact_number = ?, date = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery1, name, p.SubVillage, p.LocationID, p.NumberOfParticipants, p.ContactPerson, p.ContactNumber, p.Day1, id1)
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	updateQuery2 := "UPDATE tbl_poa SET name = ?, sub_village = ?, location_id = ?, participants = ?, contact_person = ?, contact_number = ?, date = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery2, name, p.SubVillage, p.LocationID, p.NumberOfParticipants, p.ContactPerson, p.ContactNumber, p.Day2, id2)
	if err != nil {
		tx.Rollback()
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	err = tx.Commit()
	if err != nil {
		response.Code = 500
		response.Message = "Internal server error"
		response.Success = false
		sendResponse(w, response)
		return
	}

	response.Code = 200
	response.Message = "Training Batch Updated Successfully"
	response.Success = true
	sendResponse(w, response)
}

// func getIntFromRequest(request map[string]interface{}, key string) int {
// 	value, ok := request[key].(float64)
// 	if !ok {
// 		return 0
// 	}
// 	return int(value)
// }

// func getStringFromRequest(request map[string]interface{}, key string) string {
// 	value, ok := request[key].(string)
// 	if !ok {
// 		return ""
// 	}
// 	return value
// }

// func getDateFromRequest(request map[string]interface{}, key string) time.Time {
// 	value, ok := request[key].(string)
// 	if !ok {
// 		return time.Time{}
// 	}
// 	parsedTime, err := time.Parse("2006-01-02 15:04:05", value)
// 	if err != nil {
// 		return time.Time{}
// 	}
// 	return parsedTime
// }

func sendResponse(w http.ResponseWriter, response Response11) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(response.Code)
	w.Write(jsonData)
}
