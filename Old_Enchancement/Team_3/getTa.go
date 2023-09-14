package Team_3

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func GetTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	type Travel struct {
		Id                  string `json:"id"`
		Emp_id              string `json:"emp_id"`
		Poa_id              string `json:"poa_id"`
		Start_odometer      string `json:"start_odometer"`
		Start_location_name string `json:"start_location_name"`
		End_odometer        string `json:"end_odometer"`
		End_location_name   string `json:"end_location_name"`
		Printing            string `json:"printing"`
		Telephone           string `json:"telephone"`
		Stationery          string `json:"stationery"`
		Others              string `json:"others"`
		Other_text          string `json:"other_text"`
		Klmtr               string `json:"klmtr"`
		Status              string `json:"status"`
		Entry_date          string `json:"entry_date"`
		Approved_by         string `json:"approved_by"`
		Approve_date        string `json:"approve_date"`
		Mail_send           string `json:"mail_send"`
		Extra_comments      string `json:"extra_comments"`
		Total_ta            string `json:"total_ta"`
		Purpose             string `json:"purpose"`
		Mode_of_travel      string `json:"mode_of_travel"`
		Rate_per_KM         string `json:"rate_per_KM"`
		Da                  string `json:"da"`
		Verifyed_by         string `json:"verifyed_by"`
		Verifyied_date      string `json:"verifyied_date"`
		Update_date         string `json:"update_date"`
		Extra_date          string `json:"extra_date"`
		Extra_om_comments   string `json:"extra_om_comments"`
		Fare_amount         string `json:"fare_amount"`
	}

	type Target struct {
		Id         string `json:"Id"`
		Image_name string `json:"image_name"`
		Ta_id      string `json:"ta_id"`
		Emp_id     string `json:"emp_id"`
	}

	type Response1 struct {
		Data    []Travel `json:"data"`
		Files   []Target `json:"files"`
		Success bool     `json:"success"`
		Message string   `json:"message"`
	}
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	var TargetList []Travel
	var Targets []Target
	if ta_id, ok := request["ta_id"].(string); ok {
		rows, err := DB.Query("select id, emp_id,poa_id,start_odometer,start_location_name,end_odometer,end_location_name,printing,telephone,stationery,others,other_text,klmtr,status,entry_date,COALESCE(approved_by,''),COALESCE(approve_date,''),COALESCE(mail_send,''),COALESCE(extra_comments,''),total_ta,purpose,mode_of_travel,rate_per_KM,da,COALESCE(verifyed_by,''),COALESCE(verifyied_date,''),update_date,COALESCE(extra_date,''),COALESCE(extra_om_comments,''),COALESCE(fare_amount,'') from travelling_allowance where id=?", ta_id)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
		defer rows.Close()

		var tadata Travel
		for rows.Next() {

			err := rows.Scan(&tadata.Id, &tadata.Emp_id, &tadata.Poa_id, &tadata.Start_odometer, &tadata.Start_location_name, &tadata.End_odometer, &tadata.End_location_name, &tadata.Printing, &tadata.Telephone, &tadata.Stationery, &tadata.Others, &tadata.Other_text, &tadata.Klmtr, &tadata.Status, &tadata.Entry_date, &tadata.Approved_by, &tadata.Approve_date, &tadata.Mail_send, &tadata.Extra_comments, &tadata.Total_ta, &tadata.Purpose, &tadata.Mode_of_travel, &tadata.Rate_per_KM, &tadata.Da, &tadata.Verifyed_by, &tadata.Verifyied_date, &tadata.Update_date, &tadata.Extra_date, &tadata.Extra_om_comments, &tadata.Fare_amount)
			if err != nil {
				log.Println(err)
				// json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})

			}
			TargetList = append(TargetList, tadata)
		}

		rows1, err := DB.Query("select * from taAttachments where ta_id=?", ta_id)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
		defer rows1.Close()
		var target1 Target
		for rows1.Next() {

			err := rows1.Scan(&target1.Id, &target1.Image_name, &target1.Ta_id, &target1.Emp_id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})

			}
			Targets = append(Targets, target1)

		}

	}

	Response := Response1{
		Data:    TargetList,
		Files:   Targets,
		Success: true,
		Message: "TA Data",
	}
	json.NewEncoder(w).Encode(Response)
}
