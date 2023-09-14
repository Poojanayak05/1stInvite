package Team_3

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Bus struct {
	ID             string `json:"id"`
	RegisterNumber string `json:"register_number"`
}

type bus1 struct {
	List    []Bus  `json:"list"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func GetBusList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	query := "SELECT id, UPPER(register_number) as register_number FROM bus ORDER BY id DESC"

	rows, err := DB.Query(query)
	if err != nil {
		// log.Printf("Failed to build content from sql rows: %v \n", err)
		log.Println("Error executing query:", err)
		// http.Error(w, "Internal Server Error", "Message" :http.StatusInternalServerError)
		// jsonResponse(w, response, http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "failed to execute query", "Message": err})
		return
	}
	defer rows.Close()

	var busList []Bus

	for rows.Next() {
		var bus Bus
		err := rows.Scan(&bus.ID, &bus.RegisterNumber)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		busList = append(busList, bus)
	}

	response := bus1{
		List:    busList,
		Code:    200,
		Success: true,
		Message: "successfully"}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
