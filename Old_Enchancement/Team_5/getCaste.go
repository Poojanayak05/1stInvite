package Team_5

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Caste struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Response5 struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

func GetCaste(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	// Execute the query and fetch the data
	rows, err := DB.Query("SELECT * FROM caste")
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
		return
	}
	defer rows.Close()

	// Create a slice to hold the retrieved data1
	var data1 []Caste

	// Iterate over the rows and append data1 to the slice
	for rows.Next() {
		var row Caste
		err := rows.Scan(&row.ID, &row.Name)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false}) //frontend
			return
		}
		data1 = append(data1, row)
	}
	//Reponse structure
	response := Response5{
		Data:    data1,
		Success: true,
		Message: "Caste data sent successfully",
	}

	json.NewEncoder(w).Encode(response)

	// json.NewEncoder(w).Encode(map[string]interface{}{"data": data1, "Message": "Caste data sent successfully", "success": true})

}
