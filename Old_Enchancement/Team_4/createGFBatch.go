package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type GfBatchRequest struct {
	TrainingBatchID string `json:"training_batch_id"`
	EmpID           string `json:"emp_id"`
	ProjectID       string `json:"project_id"`
}

func CreateGFBatch(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var request GfBatchRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var Proj int
	err1 := DB.QueryRow("SELECT project_id from project_emps where emp_id=? and project_id=?", request.EmpID, request.ProjectID).Scan(&Proj)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID for the given User ID", "success": false})
		} else {
			// Handle other errors
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}
	// Check if Proj is 0, indicating no valid project ID was found
	if Proj == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID for the given User ID", "success": false})
		return
	}
	var projectID string
	getProjectIDQuery := fmt.Sprintf("SELECT project_id FROM tbl_poa WHERE tb_id = %s", request.TrainingBatchID)
	err = DB.QueryRow(getProjectIDQuery).Scan(&projectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to get project ID", "success": false, "error": err})
		return
	}

	if request.ProjectID == "" {
		request.ProjectID = projectID
	}

	insertQuery := "INSERT INTO gf_batches (training_batch_id, emp_id, project_id) VALUES (?, ?, ?)"
	_, err = DB.Exec(insertQuery, request.TrainingBatchID, request.EmpID, request.ProjectID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to insert gf_batch", "success": false, "error": err})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "GF Batch Created Successfully", "success": true})
}
