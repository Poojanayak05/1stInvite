package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type TravellingAllowance1 struct {
	ID                string  `json:"id"`
	EmpID             string  `json:"emp_id"`
	PoaID             string  `json:"poa_id"`
	StartOdometer     string  `json:"start_odometer"`
	StartLocationName string  `json:"start_location_name"`
	EndOdometer       string  `json:"end_odometer"`
	EndLocationName   string  `json:"end_location_name"`
	Printing          string  `json:"printing"`
	Telephone         string  `json:"telephone"`
	Stationery        string  `json:"stationery"`
	Others            string  `json:"others"`
	OtherText         string  `json:"other_text"`
	Klmtr             string  `json:"klmtr"`
	Status            string  `json:"status"`
	Entry_date        string  `json:"entry_date"`
	ApprovedBy        *string `json:"approved_by"`
	ApproveDate       *string `json:"approve_date"`
	MailSend          *string `json:"mail_send"`
	ExtraComments     *string `json:"extra_comments"`
	TotalTA           string  `json:"total_ta"`
	Purpose           string  `json:"purpose"`
	ModeOfTravel      string  `json:"mode_of_travel"`
	RatePerKM         string  `json:"rate_per_KM"`
	DA                string  `json:"da"`
	VerifyedBy        *string `json:"verifyed_by"`
	VerifyiedDate     *string `json:"verifyied_date"`
	UpdateDate        string  `json:"update_date"`
	ExtraDate         *string `json:"extra_date"`
	ExtraOMComments   *string `json:"extra_om_comments"`
	TaName            string  `json:"Ta_Name"`
	PoaName           string  `json:"poa_name"`
	ModeOfTravelName  string  `json:"mode_of_travel_name"`
	RatePerKMName     string  `json:"rate_per_KM_name"`
	DAName            *string `json:"da_name"`
	FareAmount        string  `json:"fare_amount"`
}

func ListTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	type Request struct {
		EmpID string `json:"emp_id"`
		Date  string `json:"date"`
	}
	var taList []TravellingAllowance1
	//data := make(map[string]interface{})
	var data Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error decoding JSON", "success": false})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dateFilter := fmt.Sprintf("MONTH(entry_date) = MONTH('%s') AND", data.Date)
	query := fmt.Sprintf(`SELECT ta.id,ta.emp_id,ta.poa_id,ta.start_odometer,ta.start_location_name,ta.end_odometer,
    ta.end_location_name,ta.printing,telephone,ta.stationery,ta.others,ta.other_text,ta.klmtr,ta.status,ta.entry_date,
    COALESCE(ta.approved_by,''),COALESCE(ta.approve_date,''),COALESCE(ta.mail_send,''),COALESCE(ta.extra_comments,''),
    ta.total_ta,ta.purpose,ta.mode_of_travel,ta.rate_per_KM,ta.da,COALESCE(ta.fare_amount,''),COALESCE(ta.verifyed_by,''),
    COALESCE(ta.verifyied_date,''),ta.update_date,COALESCE(ta.extra_date,''),COALESCE(ta.extra_om_comments,''),
    CONCAT(tbl_poa.name,'_',ta.entry_date) as 'Ta_Name',tbl_poa.name as 'poa_name' FROM travelling_allowance ta
    JOIN tbl_poa ON ta.poa_id=tbl_poa.id WHERE %s emp_id=%s`, dateFilter, data.EmpID)
	// Execute the query and fetch the data
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": taList, "message": "TA Data", "success": true})
		return
	}
	defer rows.Close()
	// Create a slice to hold the retrieved TravellingAllowance1
	// Iterate over the rows and append data to the slice
	for rows.Next() {
		var ta TravellingAllowance1
		err := rows.Scan(
			&ta.ID, &ta.EmpID, &ta.PoaID, &ta.StartOdometer, &ta.StartLocationName, &ta.EndOdometer, &ta.EndLocationName,
			&ta.Printing, &ta.Telephone, &ta.Stationery, &ta.Others, &ta.OtherText, &ta.Klmtr, &ta.Status, &ta.Entry_date,
			&ta.ApprovedBy, &ta.ApproveDate, &ta.MailSend,
			&ta.ExtraComments, &ta.TotalTA, &ta.Purpose,
			&ta.ModeOfTravel, &ta.RatePerKM, &ta.DA, &ta.FareAmount,
			&ta.VerifyedBy, &ta.VerifyiedDate, &ta.UpdateDate,
			&ta.ExtraDate, &ta.ExtraOMComments, &ta.TaName, &ta.PoaName,
		)
		if err != nil {
			log.Println("Error scanning row:", err)

			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan the", "success": false})
			return

		}
		switch ta.ModeOfTravel {
		case "0":
			ta.ModeOfTravelName = ""
		case "1":
			ta.ModeOfTravelName = "Bus"
		case "2":
			ta.ModeOfTravelName = "Car"
		case "3":
			ta.ModeOfTravelName = "Bike"
		case "4":
			ta.ModeOfTravelName = "Auto"
		case "5":
			ta.ModeOfTravelName = "Train/Metro"
		case "6":
			ta.ModeOfTravelName = "Taxi"
		}
		if ta.RatePerKM == "0" {
			ta.RatePerKMName = ""
		} else {
			ta.RatePerKMName = "3.50"
		}
		switch ta.DA {
		case "0":
			ta.DAName = nil
		case "1":
			daName := "Bangalore - 200"
			ta.DAName = &daName
		case "2":
			daName := "Other District - 150"
			ta.DAName = &daName
		case "3":
			daName := "Others - 100"
			ta.DAName = &daName
		}
		if ta.ApprovedBy != nil && *ta.ApprovedBy == "" {
			ta.ApprovedBy = nil
		}
		if ta.ApproveDate != nil && *ta.ApproveDate == "" {
			ta.ApproveDate = nil
		}
		if ta.MailSend != nil && *ta.MailSend == "" {
			ta.MailSend = nil
		}
		if ta.ExtraComments != nil && *ta.ExtraComments == "" {
			ta.ApproveDate = nil
		}
		if ta.VerifyedBy != nil && *ta.VerifyedBy == "" {
			ta.VerifyedBy = nil
		}
		if ta.VerifyiedDate != nil && *ta.VerifyiedDate == "" {
			ta.VerifyiedDate = nil
		}
		if ta.ExtraOMComments != nil && *ta.ExtraOMComments == "" {
			ta.ExtraOMComments = nil
		}
		if ta.ExtraDate != nil && *ta.ExtraDate == "" {
			ta.ExtraDate = nil
		}
		if ta.DAName != nil && *ta.DAName == "" {
			ta.VerifyiedDate = nil
		}
		taList = append(taList, ta)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"data": taList, "success": true, "message": "TA Data"})

}
