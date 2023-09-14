package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PostData struct {
	EmpID             string `json:"emp_id"`
	POAID             string `json:"poa_id"`
	StartOdometer     string `json:"start_odometer"`
	StartLocationName string `json:"start_location_name"`
	EndOdometer       string `json:"end_odometer"`
	EndLocationName   string `json:"end_location_name"`
	Printing          string `json:"printing"`
	Telephone         string `json:"telephone"`
	Stationery        string `json:"stationery"`
	Others            string `json:"others"`
	OtherText         string `json:"other_text"`
	Klmtr             string `json:"klmtr"`
	InsideBangalore   string `json:"insideBangalore"`
	Date              string `json:"date"`
	ModeOfTravel      string `json:"mode_of_travel"`
	RatePerKM         string `json:"rate_per_KM"`
	Files             []int  `json:"files"`
	DA                string `json:"da"`
	FareAmount        string `json:"fare_amount"`
}

type Responseh struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func AddNewTA(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	response := make(map[string]interface{})

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Read Body Error", "success": false, "error": err.Error()})
		return
	}

	defer r.Body.Close()

	var request PostData
	err = json.Unmarshal(data, &request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to unmarshal request", "success": false, "error": err.Error()})
		return
	}
	if request.ModeOfTravel == "2" || request.ModeOfTravel == "3" {
		if request.EmpID != "" && request.POAID != "" && request.StartOdometer != "" {
			empID, _ := strconv.Atoi(request.EmpID)
			poaID, _ := strconv.Atoi(request.POAID)

			startOdometer := request.StartOdometer
			startLocationName := request.StartLocationName
			endOdometer := request.EndOdometer
			endLocationName := request.EndLocationName

			stationery, _ := strconv.ParseFloat(request.Stationery, 64)
			purpose := '-'

			var printing float64
			if request.Printing != "" {
				printing, _ = strconv.ParseFloat(request.Printing, 64)
			} else {
				printing = 0
			}

			var fareAmount float64
			if request.FareAmount != "" {
				fareAmount, _ = strconv.ParseFloat(request.FareAmount, 64)
			} else {
				fareAmount = 0
			}

			var others float64
			if request.Others != "" {
				others, _ = strconv.ParseFloat(request.Others, 64)
			} else {
				others = 0
			}
			otherText := request.OtherText

			var klmtr float64
			if request.Klmtr != "" {
				klmtr, _ = strconv.ParseFloat(request.Klmtr, 64)
			} else {
				klmtr = 0
			}

			// date := request.Date

			var modeOfTravel float64
			if request.ModeOfTravel != "" {
				modeOfTravel, _ = strconv.ParseFloat(request.ModeOfTravel, 64)
			} else {
				modeOfTravel = 0
			}

			var ratePerKM float64
			if request.RatePerKM != "" {
				ratePerKM, _ = strconv.ParseFloat(request.RatePerKM, 64)
			} else {
				ratePerKM = 0
			}

			var da float64
			if request.DA != "" {
				da, _ = strconv.ParseFloat(request.DA, 64)
			} else {
				da = 0
			}

			var telephone float64
			if request.Telephone != "" {
				telephone, _ = strconv.ParseFloat(request.Telephone, 64)
			} else {
				telephone = 0
			}

			rate := 0.0
			if ratePerKM == 1 {
				rate = 3.50
			} else if ratePerKM == 2 {
				rate = 10.0
			}

			das := 0
			switch da {
			case 1:
				das = 200
			case 2:
				das = 150
			case 3:
				das = 100
			default:
				das = 0
			}

			currentMonth := time.Now().Month()
			currentYear := time.Now().Year()

			if telephone > 0 {
				query := "SELECT id FROM travelling_allowance WHERE emp_id = ? AND MONTH(entry_date) = ? AND YEAR(entry_date) = ? AND telephone > 0"
				rows, err := DB.Query(query, empID, currentMonth, currentYear)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
					return
				}
				defer rows.Close()

				if rows.Next() {
					response["success"] = true
					response["message"] = "Telephone bill already added for this month"

					jsonData, err := json.Marshal(response)
					if err != nil {
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to marshal response", "success": false, "error": err.Error()})
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)
					return
				}
			}

			var ta float64

			if modeOfTravel == 2 || modeOfTravel == 3 {
				ta = (klmtr * rate) + float64(das) + printing + telephone + stationery + others
			} else {
				ta = fareAmount + float64(das) + printing + telephone + stationery + others
			}
			date1 := time.Now().Add(5*time.Hour + 30*time.Minute)
			qry := fmt.Sprintf("INSERT INTO travelling_allowance (emp_id, poa_id, start_odometer, start_location_name, end_odometer, end_location_name, printing, telephone, stationery, others, other_text, klmtr, status, entry_date, update_date, total_ta, purpose, mode_of_travel, rate_per_KM, da, fare_amount) VALUES (%d, %d, '%s', '%s', '%s', '%s', %f, %f, %f, %f, '%s', %f, 0, '%s', '%s', %f, '%c', %f, %f, %f, %f)", empID, poaID, startOdometer, startLocationName, endOdometer, endLocationName, printing, telephone, stationery, others, otherText, klmtr, date1, date1, ta, purpose, modeOfTravel, ratePerKM, da, fareAmount)
			_, err = DB.Exec(qry)
			if err != nil {
				response["success"] = false
				response["message"] = "Failed to add record"
				response["error"] = err.Error()

				jsonData, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}

			idQuery := "SELECT LAST_INSERT_ID()"
			var id int
			err = DB.QueryRow(idQuery).Scan(&id)
			if err != nil {
				response["success"] = false
				response["message"] = "Failed to retrieve ID"
				response["error"] = err.Error()

				jsonData, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}

			for _, file := range request.Files {
				updateQuery := fmt.Sprintf("UPDATE taAttachments SET ta_id = %d WHERE id = %d", id, file)
				_, err := DB.Exec(updateQuery)
				if err != nil {
					log.Println("ERROR>>", err)

					response["success"] = false
					response["message"] = "Failed to update taAttachments"
					response["error"] = err.Error()

					jsonData, _ := json.Marshal(response)
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)
					return
				}
			}

			response["success"] = true
			response["message"] = "Added Successfully"

			jsonData, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	} else if request.ModeOfTravel == "1" || request.ModeOfTravel == "4" || request.ModeOfTravel == "5" || request.ModeOfTravel == "6" {
		if request.EmpID != "" && request.POAID != "" {
			empID, _ := strconv.Atoi(request.EmpID)
			poaID, _ := strconv.Atoi(request.POAID)

			startOdometer := request.StartOdometer
			startLocationName := request.StartLocationName
			endOdometer := request.EndOdometer
			endLocationName := request.EndLocationName

			stationery, _ := strconv.ParseFloat(request.Stationery, 64)
			purpose := '-'

			var printing float64
			if request.Printing != "" {
				printing, _ = strconv.ParseFloat(request.Printing, 64)
			} else {
				printing = 0
			}

			var fareAmount float64
			if request.FareAmount != "" {
				fareAmount, _ = strconv.ParseFloat(request.FareAmount, 64)
			} else {
				fareAmount = 0
			}

			var others float64
			if request.Others != "" {
				others, _ = strconv.ParseFloat(request.Others, 64)
			} else {
				others = 0
			}
			otherText := request.OtherText

			var klmtr float64
			if request.Klmtr != "" {
				klmtr, _ = strconv.ParseFloat(request.Klmtr, 64)
			} else {
				klmtr = 0
			}

			// date := request.Date

			var modeOfTravel float64
			if request.ModeOfTravel != "" {
				modeOfTravel, _ = strconv.ParseFloat(request.ModeOfTravel, 64)
			} else {
				modeOfTravel = 0
			}

			var ratePerKM float64
			if request.RatePerKM != "" {
				ratePerKM, _ = strconv.ParseFloat(request.RatePerKM, 64)
			} else {
				ratePerKM = 0
			}

			var da float64
			if request.DA != "" {
				da, _ = strconv.ParseFloat(request.DA, 64)
			} else {
				da = 0
			}

			var telephone float64
			if request.Telephone != "" {
				telephone, _ = strconv.ParseFloat(request.Telephone, 64)
			} else {
				telephone = 0
			}

			rate := 0.0
			if ratePerKM == 1 {
				rate = 3.50
			} else if ratePerKM == 2 {
				rate = 10.0
			}

			das := 0
			switch da {
			case 1:
				das = 200
			case 2:
				das = 150
			case 3:
				das = 100
			default:
				das = 0
			}

			currentMonth := time.Now().Month()
			currentYear := time.Now().Year()

			if telephone > 0 {
				query := "SELECT id FROM travelling_allowance WHERE emp_id = ? AND MONTH(entry_date) = ? AND YEAR(entry_date) = ? AND telephone > 0"
				rows, err := DB.Query(query, empID, currentMonth, currentYear)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
					return
				}
				defer rows.Close()

				if rows.Next() {
					response["success"] = true
					response["message"] = "Telephone bill already added for this month"

					jsonData, err := json.Marshal(response)
					if err != nil {
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to marshal response", "success": false, "error": err.Error()})
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)
					return
				}
			}

			var ta float64

			if modeOfTravel == 2 || modeOfTravel == 3 {
				ta = (klmtr * rate) + float64(das) + printing + telephone + stationery + others
			} else {
				ta = fareAmount + float64(das) + printing + telephone + stationery + others
			}
			date1 := time.Now().Add(5*time.Hour + 30*time.Minute)
			qry := fmt.Sprintf("INSERT INTO travelling_allowance (emp_id, poa_id, start_odometer, start_location_name, end_odometer, end_location_name, printing, telephone, stationery, others, other_text, klmtr, status, entry_date, update_date, total_ta, purpose, mode_of_travel, rate_per_KM, da, fare_amount) VALUES (%d, %d, '%s', '%s', '%s', '%s', %f, %f, %f, %f, '%s', %f, 0, '%s', '%s', %f, '%c', %f, %f, %f, %f)", empID, poaID, startOdometer, startLocationName, endOdometer, endLocationName, printing, telephone, stationery, others, otherText, klmtr, date1, date1, ta, purpose, modeOfTravel, ratePerKM, da, fareAmount)
			_, err = DB.Exec(qry)
			if err != nil {
				response["success"] = false
				response["message"] = "Failed to add record"
				response["error"] = err.Error()

				jsonData, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}

			idQuery := "SELECT LAST_INSERT_ID()"
			var id int
			err = DB.QueryRow(idQuery).Scan(&id)
			if err != nil {
				response["success"] = false
				response["message"] = "Failed to retrieve ID"
				response["error"] = err.Error()

				jsonData, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}

			for _, file := range request.Files {
				updateQuery := fmt.Sprintf("UPDATE taAttachments SET ta_id = %d WHERE id = %d", id, file)
				_, err := DB.Exec(updateQuery)
				if err != nil {
					log.Println("ERROR>>", err)

					response["success"] = false
					response["message"] = "Failed to update taAttachments"
					response["error"] = err.Error()

					jsonData, _ := json.Marshal(response)
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonData)
					return
				}
			}

			response["success"] = true
			response["message"] = "Added Successfully"

			jsonData, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	} else {
		response["success"] = true
		response["message"] = "Send Required Parameters"

		jsonData, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
