package Team_4

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Funder struct {
	FunderID   string `json:"funderID"`
	FunderName string `json:"funderName"`
}

type Funder1 struct {
	List    []Funder `json:"list"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Success bool     `json:"success"`
}

func GetFunderList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	query := "SELECT funderID, funderName FROM funder WHERE status = 1 ORDER BY funderID DESC, funderName"

	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var funderList []Funder

	for rows.Next() {
		var funder Funder
		err := rows.Scan(&funder.FunderID, &funder.FunderName)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		funderList = append(funderList, funder)
	}

	response := Funder1{
		List:    funderList,
		Code:    200,
		Success: true,
		Message: "successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
