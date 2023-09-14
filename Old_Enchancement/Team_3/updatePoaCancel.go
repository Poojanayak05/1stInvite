package Team_3

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Request12 struct {
	PoaID string `json:"poa_id"`
	Day   string `json:"day"`
}

func UpdatePoaCancel(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request Request12
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error decoding the body", "success": false, "error": err})
		return
	}
	//status 1 - Reschedule, 2 - Cancelled
	if request.PoaID != "" {
		id := request.PoaID
		day := request.Day

		var updatePoa string

		if day == "1" {
			updatePoa = "UPDATE tbl_poa SET status = '2' WHERE tb_id = ? AND type = '1'"
		} else {
			fields := "trn_batch.tb_id as id, trn_batch.date"
			query := "SELECT " + fields + " FROM tbl_poa trn_batch WHERE trn_batch.id = ?"

			var date string
			var tbID string

			err := db.QueryRow(query, id).Scan(&tbID, &date)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": err, "success": false, "error": err})
				return
			}

			updatePoa = "UPDATE tbl_poa SET status = '2' WHERE primary_id = ?"
		}

		_, err = db.Exec(updatePoa, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Poa Cancelled Successfully ", "success": true})
			return
		}
	} else {

		w.Header().Set("Content-Type", "application/json")
		// Return JSON response
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Fill  all the required fields", "success": false})
		return
	}

}
