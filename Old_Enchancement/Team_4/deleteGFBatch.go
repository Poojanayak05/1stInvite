package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type BatchRequest struct {
	EmpID           string `json:"emp_id"`
	TrainingBatchID string `json:"training_batch_id"`
	ProjectID       string `json:"project_id"`
}

func DeleteGFBatch(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var req BatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	if req.EmpID != "" && req.ProjectID != "" {
		var proj_id int
		err := DB.QueryRow("select project_id from tbl_poa where tb_id = ?", req.TrainingBatchID).Scan(&proj_id)
		if err != nil {
			if err == sql.ErrNoRows {
				// Handle the case where no rows were found
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Tb ID", "success": false})
			} else {
				// Handle other errors that may have occurred during query execution
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
			}
			return
		}
		// Check if Proj is 0, indicating no valid project ID was found
		if proj_id == 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Tb ID", "success": false})
			return
		}
		_, err1 := DB.Exec("DELETE FROM gf_batches WHERE emp_id = ? and training_batch_id = ?", req.EmpID, req.TrainingBatchID)
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "GF Batch Deleted Successfully", "success": true})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
}
