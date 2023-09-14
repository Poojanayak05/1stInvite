package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Response represents the response structure

type RESP struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// POA represents the POA structure
type POA struct {
	PoaID       string `json:"poa_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AllDay      int    `json:"all_day"`
	Date        string `json:"date"`
	Date2       string `json:"date2"`
	UserID      string `json:"user_id"`
}

func UpdateRescheduleEvent(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	// Read the request body

	// Initialize a POA variable to hold the decoded data
	// Decode the JSON request
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print("line 44")

		return
	}

	defer r.Body.Close()

	var request POA
	err = json.Unmarshal(data, &request)
	fmt.Print("data:", request)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed to Unmarshal", "success": false, "error": err.Error()})

	}

	poaid, _ := strconv.Atoi(request.PoaID)

	// Validate the POA ID\
	if poaid == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid POA ID", "success": false, "error": err.Error()})

		return
	} else {

		// Update the POA
		err = updatePOA(w, DB, request)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})

			return
		}

	}
}

func updatePOA(w http.ResponseWriter, db *sql.DB, request POA) error {

	poaid, _ := strconv.Atoi(request.PoaID)
	response := make(map[string]interface{})
	reqDate, _ := parseDateTime(request.Date)
	reqDate2, _ := parseDateTime(request.Date2)

	// Build the SQL query and parameters
	var params []interface{}
	// query := "UPDATE tbl_poa SET"
	fields := []string{}

	if request.Name != "" {
		fields = append(fields, "name = ?")
		params = append(params, request.Name)
	}

	if request.Description != "" {
		fields = append(fields, "description = ?")
		params = append(params, request.Description)
	}

	if request.AllDay > 0 {
		fields = append(fields, "all_day = ?")
		params = append(params, request.AllDay)
	}

	if request.Date != "" {

		fields = append(fields, "date = ?")
		params = append(params, reqDate)
	}

	if request.Date2 != "" {

		fields = append(fields, "date2 = ?")
		params = append(params, reqDate2)
	}

	// Check if any fields are being updated
	if len(fields) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No fields to update", "success": false})
	}

	// Duplicate date & time validation start

	oldDate, uid, err := getPOADate(db, poaid)
	if err != nil {

		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
		return err
	}

	if oldDate == reqDate {
		count, err := countPOAsWithDate(db, reqDate, uid)
		if err != nil {

			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
			return err
		}

		if count > 0 {
			return json.NewEncoder(w).Encode(map[string]interface{}{"code": 409, "message": "POA already exists for date & time", "success": false})

		}
	}

	// Duplicate date & time validation End

	// Prepare the query
	// tblFields := "DATE_FORMAT(date, '%Y-%m-%d %h:%i') as date, DATE_FORMAT(date2, '%Y-%m-%d %h:%i') as date2"
	q := ("SELECT DATE_FORMAT(date, '%Y-%m-%d %H:%i') as date, DATE_FORMAT(date2, '%Y-%m-%d %H:%i') as date2 FROM tbl_poa WHERE id = ?")

	// Execute the query
	row := db.QueryRow(q, poaid)

	var fetchedDate, fetchedDate2 string

	// Retrieve the data
	err = row.Scan(&fetchedDate, &fetchedDate2)
	fetchedD, _ := time.Parse("2006-01-02 15:04", fetchedDate)

	fetchedD2, _ := time.Parse("2006-01-02 15:04", fetchedDate2)

	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})

		}
	}
	// Check if the fetched dates are different from the provided dates

	if fetchedD != reqDate || fetchedD2 != reqDate2 {

		fields = append(fields, "status = '1'")
		// Update the POA
		query := fmt.Sprintf("UPDATE tbl_poa SET %s WHERE primary_id = ?", strings.Join(fields, ","))
		params = append(params, poaid)
		_, err = db.Exec(query, params...)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
		}

		response["code"] = 200
		response["success"] = true
		response["message"] = "Reschedule Changed Successfully"
		// Convert response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false, "error": err.Error()})
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send response
		_, err = w.Write(responseJSON)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to write response", "success": false, "error": err.Error()})
		}
	} else {
		response["code"] = 200
		response["success"] = false
		response["message"] = "Reschedule Changed Unsuccessfully"
		// Convert response to JSON
		responseJSON, err := json.Marshal(response)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false, "error": err.Error()})
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send response
		_, err = w.Write(responseJSON)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to write response", "success": false, "error": err.Error()})
		}
	}

	return nil
}

func parseDateTime(datetime string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", datetime)

	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func getPOADate(db *sql.DB, id int) (time.Time, int, error) {
	query := "SELECT DATE_FORMAT(date, '%Y-%m-%d %H:%i') as date,user_id FROM tbl_poa WHERE id = ?"
	var date string
	var user_id int
	err := db.QueryRow(query, id).Scan(&date, &user_id)
	if err != nil {

		return time.Time{}, 0, err
	}
	// t,_:=parseDateTime(date)
	t, _ := time.Parse("2006-01-02 15:04", date)

	return t, user_id, nil
}

func countPOAsWithDate(db *sql.DB, date time.Time, uid int) (int, error) {
	query := "SELECT COUNT(id) as count FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'"
	var count int
	err := db.QueryRow(query, date, uid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
