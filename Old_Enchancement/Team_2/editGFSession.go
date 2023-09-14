package Team_2

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Response represents the JSON response structure
// type GFSessionResponse struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// 	Success bool   `json:"success"`
// }

// GFSession represents the structure of GF session data
type GFSession struct {
	ID          string `json:"gf_session_id"`
	TBID        string `json:"tb_id"`
	Name        string `json:"name"`
	TBName      string `json:"tb_name"`
	PlanDate    string `json:"plan_date,omitempty"`
	SessionType string `json:"gf_session_type,omitempty"`
	UserID      int    `json:"user_id"`
	Check_in    string `json:"check_in"`
}

func EditGFSessionMethod(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// var response GFSessionResponse

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request", "success": false})
		return
	}
	// Parse request data
	var request GFSession
	err = json.Unmarshal(body, &request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Data", "success": false})
		return
	}

	var date time.Time
	if request.PlanDate != "" {
		parsedDate, err := time.Parse("02-01-2006 15:04:05", request.PlanDate)
		if err != nil {
			// Handle the error if the date string is in an incorrect
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid date format", "success": false, "error": err.Error()})

		} else {
			date = parsedDate
		}
	}

	// Retrieve existing GF session data
	// var existingSession GFSession
	err = db.QueryRow("SELECT COALESCE(session_type,''), DATE_FORMAT(date, '%Y-%m-%d %h:%i:%s') as date, user_id,COALESCE(check_in,'') FROM tbl_poa WHERE id = ?", request.ID).Scan(
		&request.SessionType, &request.PlanDate, &request.UserID, &request.Check_in)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to retrieve GF session data", "success": false})
		return

	}
	exist_date, _ := time.Parse("02-01-2006 15:04:00", request.PlanDate)

	// Duplicate date & time validation
	if exist_date != date {
		var count int
		err = db.QueryRow("SELECT COUNT(id) FROM tbl_poa WHERE date = ? AND user_id = ? AND status != 2", request.PlanDate, request.UserID).Scan(&count)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to validate duplicate GF session", "success": false})
			return
		}
		// if count > 0 {
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusConflict, "message": "Poa already exist for date & time", "success": false})
		// 	return
		// }

		// Generate GF session name
		// name := GetGFSessionName(request.TBName, request.SessionType)

		// Update GF session data
		if request.Check_in == "" {
			// Check-in is not null, so you can proceed with updating the date

			// Update GF session data
			sqlResult, err := db.Exec("UPDATE tbl_poa SET date = ? WHERE id = ?",
				date, request.ID)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to update GF session data", "success": false})
				return
			}
			rowsAffected, _ := sqlResult.RowsAffected()
			if rowsAffected > 0 {
				// Update successfully
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "GF session updated successfully", "success": true})
			} else {
				// Update did not affect any rows
				json.NewEncoder(w).Encode(map[string]interface{}{"code": 404, "message": "GF session already exists", "success": false})
			}
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 404, "message": "Can't update already checked In session", "success": false})

		}

	}
}

func GetGFSessionName(name string, sessionType int) string {
	nameArr := strings.Split(name, "_")
	var sessionName string

	switch sessionType {
	case 1:
		sessionName = "CM"
	case 3:
		sessionName = "BV"
	default:
		sessionName = "VV"
	}

	return nameArr[0] + "_" + sessionName
}
