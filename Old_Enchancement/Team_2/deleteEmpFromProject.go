package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func DeleteEmpFromProject(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type RequestBody struct {
		ProjectID string `json:"project_id"`
		EmpID     string `json:"emp_id"`
	}
	decoder := json.NewDecoder(r.Body)

	var data RequestBody

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	if data.EmpID != "" && data.ProjectID != "" {

		empID := data.EmpID
		projectID := data.ProjectID

		deleteQuery := fmt.Sprintf("DELETE FROM project_emps WHERE emp_id = '%s' AND project_id = '%s'", empID, projectID)

		result, err := DB.Exec(deleteQuery)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		if rowsAffected != 0 {
			updateHistory := fmt.Sprintf("UPDATE project_emps_history SET endDate = current_date() WHERE projectId='%s' AND emp_id = '%s'", projectID, empID)

			result, err := DB.Exec(updateHistory)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
				return
			}
			check, err := result.RowsAffected()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			fmt.Println(check)
			json.NewEncoder(w).Encode(map[string]interface{}{"message": "Deleted Successfully", "success": true})
			return

		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"message": "Delete operation Failed", "success": false})
			return

		}

	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Delete operation Failed", "success": true})
		return

	}

}
