package Team_5

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Request2 struct {
	TaID string `json:"ta_id"`
}

type Response4 struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func DeleteTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	var request Request2
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request body", "success": false, "error": err})
		// response := Responsed{Success: true, Message: "Send Required Parameters"}
		// json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()

	if request.TaID != "" {
		var exists bool
		err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM travelling_allowance WHERE id = ?)", request.TaID).Scan(&exists)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to check the existence of the travelling allowance", "success": false})
			return
		}

		if !exists {
			// The record does not exist
			response := Response4{Success: false, Message: "Travelling allowance does not exist"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Delete the travelling allowance record with the specified ID
		_, err = DB.Exec("DELETE FROM travelling_allowance WHERE id = ?", request.TaID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to delete the travelling allowance", "success": false})
			return
		}

		response := Response4{Message: "Delete Successfully", Success: true}
		json.NewEncoder(w).Encode(response)
	} else {
		response := Response4{Success: false, Message: "Send Required Parameters"}
		json.NewEncoder(w).Encode(response)
	}
}
