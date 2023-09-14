package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func UpdateEnrolledGelathi(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	// Define the request structure
	type Request struct {
		Gelathi_ID string `json:"gelathi_id"`
		Circle_ID  string `json:"circle_id"`
		Flag       string `json:"flag"`
	}
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// Set response status to 400 Bad Request if the request body is invalid
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var queryCircle, queryType string
	if request.Gelathi_ID != "" && request.Circle_ID != "" && request.Flag != "" {
		// Convert string values to integers
		GelathiID, _ := strconv.Atoi(request.Gelathi_ID)
		CircleID, _ := strconv.Atoi(request.Circle_ID)
		Flag, _ := strconv.Atoi(request.Flag)

		if Flag == 1 {
			//INSERT query
			queryCircle = fmt.Sprintf("INSERT INTO gelathi_circle (circle_id, gelathi_id) VALUES ('%d', '%d')", CircleID, GelathiID)
			rows, _ := DB.Exec(queryCircle)
			queryType = "Gelathi Added Successfully"
			rowsAffected, _ := rows.RowsAffected()

			if rowsAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Gelathi Added Unsuccessfully", "success": false, "error": err})
				return
			}

		} else {
			//  DELETE query
			queryCircle = fmt.Sprintf("DELETE FROM gelathi_circle WHERE circle_id = '%d' AND gelathi_id = '%d'", CircleID, GelathiID)
			rows, _ := DB.Exec(queryCircle)
			queryType = "Gelathi Removed Successfully"
			rowsAffected, _ := rows.RowsAffected()

			if rowsAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Gelathi Removed Unsuccessfully", "success": false, "error": err})
				return
			}

		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": queryType, "success": true})
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter all Required fields", "success": false, "error": err})
		return
	}

}
