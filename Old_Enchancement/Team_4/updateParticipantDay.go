package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ParticipantRequest struct {
	TbID           string `json:"tb_id"`
	ParticipantID  string `json:"participant_id"`
	Status         string `json:"status"`
	ParticipantDay string `json:"participant_day2"`
}

func UpdateParticipantDay(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var request ParticipantRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var User int
	err1 := DB.QueryRow("SELECT id from training_participants where id=?", request.ParticipantID).Scan(&User)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Participant ID", "success": false})
		} else {
			// Handle other errors that may have occurred during query execution
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}
	// Check if Proj is 0, indicating no valid project ID was found
	if User == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Participant ID", "success": false})
		return
	}

	var Proj int
	err2 := DB.QueryRow("SELECT tb_id from training_participants where id=? and tb_id=?", request.ParticipantID, request.TbID).Scan(&Proj)
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

	stmt, err := DB.Prepare("UPDATE training_participants SET day2=?, participant_day2=? WHERE tb_id=? AND id=?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to prepare update statement", "success": false, "error": err})
		return

	}
	defer stmt.Close()

	_, err = stmt.Exec(request.Status, request.ParticipantDay, request.TbID, request.ParticipantID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to update participant", "success": false, "error": err})
		return

	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Participant Updated Successfully", "success": true})

}
