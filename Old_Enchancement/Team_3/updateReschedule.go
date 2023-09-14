package Team_3

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database configuration

// Define a struct to hold the request data
type Request13 struct {
	PoaID    string `json:"poa_id"`
	DateTime string `json:"date_time"`
}

// Define a struct to hold the response data
type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"message"`
	Success bool   `json:"success"`
}

func UpdateReschedule(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Set the response content type to JSONw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	// Read the request body
	var request Request13
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}

	// Check the date time format
	dateTime, err := time.Parse("2006-01-02 15:04", request.DateTime)
	if err != nil {
		response := Response{
			Code:    http.StatusBadRequest,
			Msg:     "Not Proper Date Time Format",
			Success: false,
		}
		jsonResponse(w, response, http.StatusBadRequest)
		return
	}

	currentDateTime := time.Now()
	if dateTime.Before(currentDateTime) {
		response := Response{
			Code:    http.StatusBadRequest,
			Msg:     "Date and time must be greater than current date and time",
			Success: false,
		}
		jsonResponse(w, response, http.StatusBadRequest)
		return
	}

	// Duplicate date & time validation

	query := "SELECT date, user_id FROM tbl_poa WHERE id = ?"
	var existDate string
	var userID int
	err = DB.QueryRow(query, request.PoaID).Scan(&existDate, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ID doesn't exist in tbl_poa
			response := Response{
				Code:    http.StatusNotFound,
				Msg:     "Poa doesn't exist",
				Success: false,
			}
			jsonResponse(w, response, http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to query database1", http.StatusInternalServerError)
		return
	}

	if existDate != dateTime.Format("2006-01-02 15:04") {
		countQuery := "SELECT COUNT(id) FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'"
		var count int
		err = DB.QueryRow(countQuery, dateTime.Format("2006-01-02 15:04"), userID).Scan(&count)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			response := Response{
				Code:    http.StatusConflict,
				Msg:     "Poa already exists for date & time",
				Success: false,
			}
			jsonResponse(w, response, http.StatusConflict)
			return
		}
	}

	// Perform the database transaction
	t, err := DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
		return
	}

	updatePoa := "UPDATE tbl_poa SET date = ?, status = '1' WHERE primary_id = ?"
	_, err = t.Exec(updatePoa, dateTime.Format("2006-01-02 15:04"), request.PoaID)
	if err != nil {
		t.Rollback()
		http.Error(w, "Failed to update tbl_poa", http.StatusInternalServerError)
		return
	}

	insertQuery := "INSERT INTO tbl_poa_history (poa_id, original_date, reschedule_date) VALUES (?, ?, ?)"
	_, err = t.Exec(insertQuery, request.PoaID, existDate, dateTime.Format("2006-01-02 15:04"))
	if err != nil {
		t.Rollback()
		http.Error(w, "Failed to insert into tbl_poa_history", http.StatusInternalServerError)
		return
	}

	err = t.Commit()
	if err != nil {
		http.Error(w, "Failed to commit database transaction", http.StatusInternalServerError)
		return
	}

	response := Response{
		Code:    http.StatusOK,
		Msg:     "Reschedule Changed Successfully",
		Success: true,
	}
	jsonResponse(w, response, http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	responseJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
