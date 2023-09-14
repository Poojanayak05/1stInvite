package Team_5

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TravellingAllowance2 struct {
	ID                string `json:"ta_id"`
	EmpID             string `json:"emp_id"`
	PoaID             string `json:"poa_id"`
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
	Status            string `json:"status"`
	EntryDate         string `json:"entry_date"`
	UpdatedDate       string `json:"update_date"`
	TotalTA           string `json:"total_ta"`
	Purpose           string `json:"purpose"`
	ModeOfTravel      string `json:"mode_of_travel"`
	RatePerKM         string `json:"rate_per_KM"`
	DA                string `json:"da"`
}

func UpdateTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
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
		TAID              string   `json:"ta_id"`
		EmpID             string   `json:"emp_id"`
		PoaID             string   `json:"poa_id"`
		StartOdometer     string   `json:"start_odometer"`
		StartLocationName string   `json:"start_location_name"`
		EndOdometer       string   `json:"end_odometer"`
		EndLocationName   string   `json:"end_location_name"`
		Printing          string   `json:"printing"`
		Telephone         string   `json:"telephone"`
		Stationery        string   `json:"stationery"`
		Others            string   `json:"others"`
		OtherText         string   `json:"other_text"`
		Klmtr             string   `json:"klmtr"`
		ModeOfTravel      string   `json:"mode_of_travel"`
		Purpose           string   `json:"purpose"`
		RatePerKM         string   `json:"rate_per_KM"`
		DA                string   `json:"da"`
		UpdateDate        string   `json:"update_date"`
		Files             []string `json:"files"`
	}

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error decoding JSON", "success": false})
		return
	}

	query := `UPDATE travelling_allowance SET
		emp_id=?, poa_id=?, start_odometer=?, start_location_name=?, end_odometer=?, end_location_name=?,
		printing=?, telephone=?, stationery=?, others=?, other_text=?, klmtr=?, status=0, entry_date=?, update_date=?,
		total_ta=?, approved_by='', approve_date='', purpose=?, mode_of_travel=?, rate_per_KM=?, da=?
		WHERE id=?`
	// Execute the query and fetch the data
	stmt, err := DB.Prepare(query)

	if err != nil {
		log.Println("Error preparing statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Error preparing statement", "success": false, "error": err})
		return
	}
	defer stmt.Close()

	// Get current date if Date is not provided
	date := request.UpdateDate
	if date == "" {
		date = time.Now().Format("06/01/02")
	}

	// Calculate total_ta based on mode_of_travel, klmtr, and other variables
	klmtr, _ := strconv.ParseFloat(request.Klmtr, 64)
	da, _ := strconv.ParseFloat(request.DA, 64)
	printing, _ := strconv.ParseFloat(request.Printing, 64)
	telephone, _ := strconv.ParseFloat(request.Telephone, 64)
	stationery, _ := strconv.ParseFloat(request.Stationery, 64)
	others, _ := strconv.ParseFloat(request.Others, 64)

	var ta float64
	if request.ModeOfTravel == "2" || request.ModeOfTravel == "3" {
		ta = (klmtr * 3.50) + da + printing + telephone + stationery + others
	} else {
		ta = klmtr + da + printing + telephone + stationery + others
	}
	if request.EmpID == "" || request.EmpID == "0" || request.PoaID == "" || request.PoaID == "0" || request.TAID == "" || request.TAID == "0" {
		log.Println("Error preparing statement:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Send the Required Parameters", "success": false})
		return
	}

	_, err = stmt.Exec(request.EmpID, request.PoaID, request.StartOdometer, request.StartLocationName, request.EndOdometer,
		request.EndLocationName, request.Printing, request.Telephone, request.Stationery, request.Others, request.OtherText,
		request.Klmtr, date, date, ta, request.Purpose, request.ModeOfTravel, request.RatePerKM, request.DA, request.TAID)

	if err != nil {
		log.Println("Error updating travelling_allowance:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Error updating travelling_allowance", "success": false})
		return
	}

	for _, file := range request.Files {
		_, err = DB.Exec("UPDATE taAttachments SET ta_id=? WHERE id=?", request.TAID, file)
		if err != nil {
			log.Println("Error updating taAttachments:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Error updating taAttachments", "success": false})
			return
		}
	}
	//Reponse
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully Updated",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
