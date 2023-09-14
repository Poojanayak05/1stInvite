package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type EmailExits struct {
	Office_Email_ID string `json:"office_email_id"`
}

func GetEmailExits(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Check if the HTTP method is GET
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	var request EmailExits
	var data []EmailExits
	// Decode the JSON request body into the 'request' struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// Set response status to 400 Bad Request if the request body is invalid
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	if request.Office_Email_ID != "" {
		// Construct the SQL query to check if the office email ID already exists
		fields := "emp.officeMailId"
		query := fmt.Sprintf("SELECT '%s'  FROM employee emp where emp.officeMailId = '%s' and emp.status = 1", fields, request.Office_Email_ID)

		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			// Set response status to 400 Bad Request if the database query fails
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var mail string
			// Set response status to 400 Bad Request if there is an error scanning the database result
			err := rows.Scan(&mail)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			data = append(data, request)
		}
	}

	if data != nil {
		// Set response header to indicate JSON content type
		w.Header().Set("Content-Type", "application/json")
		// Return JSON response indicating that the office email ID already exists
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 409, "message": "Office Email Id already exist", "success": false})
		return
	} else {
		// Set response header to indicate JSON content type
		w.Header().Set("Content-Type", "application/json")
		// Return JSON response indicating success
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Successfully", "success": true})
		return
	}
}
