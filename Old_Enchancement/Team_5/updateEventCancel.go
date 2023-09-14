package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func UpdateEventCancel(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	// Define the request structure
	type Request struct {
		Poa_ID string `json:"poa_id"`
	}

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// Set response status to 400 Bad Request if the request body is invalid
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}

	if request.Poa_ID != "" {

		// Convert string value to an integer
		ID, _ := strconv.Atoi(request.Poa_ID)
		//status 1 - Reschedule, 2 - Cancelled
		updatePoa := fmt.Sprintf("UPDATE tbl_poa SET status = '2' WHERE primary_id = '%d'", ID)
		rows, _ := DB.Exec(updatePoa)
		rowsAffected, _ := rows.RowsAffected()

		if rowsAffected == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Poa Cancelled Unsuccessfully", "success": false, "error": err})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Poa Cancelled Successfully", "success": true})
		return

	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter all Required fields", "success": false, "error": err})
		return
	}

}
