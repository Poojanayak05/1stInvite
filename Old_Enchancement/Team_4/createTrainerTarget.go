package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Target struct {
	EmpID     string `json:"emp_id"`
	EmpTarget string `json:"emp_target"`
}

func CreateTrainerTarget(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var arr map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&arr)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode JSON", "success": false, "error": err})
		return
	}
	targetList := arr["target_list"].([]interface{})
	projectID := arr["project_id"].(string)

	for _, target := range targetList {
		targetMap := target.(map[string]interface{})
		empID := targetMap["emp_id"].(string)
		empTarget, err := strconv.Atoi(targetMap["emp_target"].(string))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Emp Target", "success": false, "error": err})
			return
		}
		var Proj int
		err1 := DB.QueryRow("SELECT project_id from project_emps where emp_id=? and project_id=?", empID, projectID).Scan(&Proj)
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
		if Proj == 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID for the given User ID", "success": false})
			return
		}
		_, err2 := DB.Exec("UPDATE project_emps SET target = ? WHERE emp_id = ? AND project_id = ?", empTarget, empID, projectID)
		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to update target", "success": false, "error": err})
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Trainer Assigned Successfully", "success": true})
}
