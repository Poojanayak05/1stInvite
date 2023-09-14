package Team_2

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func DeleteBus(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	type Request struct {
		Bus_id string `json:"bus_id"`
	}

	var p Request
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})
		return
	}

	if p.Bus_id != "" {
		str := "DELETE FROM bus WHERE id = ?"
		rows, err := DB.Query(str, p.Bus_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
			return
		}
		defer rows.Close()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "success": true, "message": "Bus Deleted Successfully"})
	} else {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 404, "success": false, "message": "Bus Deleted Unsuccessfully"})
	}
}
