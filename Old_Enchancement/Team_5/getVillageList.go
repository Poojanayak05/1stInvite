package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response8 struct {
	List    []Location `json:"list"`
	Code    int        `json:"code"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

func GetVillageList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	locations := []Location{}
	//Decoding the json  structure
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"list": locations, "message": "successfully", "success": true, "code": 200})
		return
	}

	searchParam, ok := data["search"].(string)
	if !ok {
		searchParam = ""
	}

	searchFilter := ""
	if len(searchParam) > 0 {
		searchFilter = " and name like '%" + searchParam + "%'"
	}

	talukIDFloat, ok := data["taluk_id"].(float64)
	if !ok {
		http.Error(w, "Invalid taluk_id", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

		return
	}
	//converting talukIDFloat to int
	talukID := int(talukIDFloat)

	var query string
	if talukID != 194 {
		query = fmt.Sprintf("SELECT id, name FROM location WHERE type = 5 AND parentId = %d %s ORDER BY name", talukID, searchFilter)
	} else {
		query = fmt.Sprintf("SELECT id, name FROM location WHERE type = 5 AND parentId IN (%d, 187, 189) %s ORDER BY name", talukID, searchFilter)
	}
	// Execute the query and fetch the data
	rows, err := DB.Query(query)
	if err != nil {
		
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
		return
	}
	defer rows.Close()
	// Create a slice to hold the retrieved Location

	// Iterate over the rows and append ocation to the slice
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Location ID doesn't exist", "success": false}) //frontend
			return
		}
		locations = append(locations, location)
	}
	//Reponse structure
	response := Response8{
		List:    locations,
		Code:    200,
		Success: true,
		Message: "successfully",
	}

	// json.NewEncoder(w).Encode(map[string]interface{}{"list": locations, "code": 200, "success": true, "message": "successfully"})
	json.NewEncoder(w).Encode(response)

}
