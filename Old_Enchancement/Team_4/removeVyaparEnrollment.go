package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func RemoveVyaparEnrollment(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var request RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	var User int
	err1 := DB.QueryRow("SELECT id from training_participants where id=?", request.ID).Scan(&User)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid ID", "success": false})
		} else {
			// Handle other errors that may have occurred during query execution
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}
	// Check if Proj is 0, indicating no valid project ID was found
	if User == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid ID", "success": false})
		return
	}
	var Proj int
	err2 := DB.QueryRow("SELECT tb_id from training_participants where id=? and tb_id=?", request.ID, request.TBID).Scan(&Proj)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Training Batch ID for the given ID", "success": false})
		} else {
			// Handle other errors
			log.Println(err2)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}

	stmt := "UPDATE training_participants SET VyaparEnrollment='0' WHERE tb_id=? AND id=?"

	_, err = DB.Exec(stmt, request.TBID, request.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to update participant", "success": false, "error": err})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Gelathi Removed Successfully", "success": true})

}
