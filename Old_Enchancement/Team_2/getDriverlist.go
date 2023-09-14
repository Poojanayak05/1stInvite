package Team_2

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

//-------------------------------- get Driver List ---------------------

func DriverList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	//func Login(w http.ResponseWriter, r *http.Request) {
	//SetupCORS(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	type listdriver struct {
		Id         string `json:"id"`
		First_name string `json:"first_name"`
	}

	rows, err := DB.Query(`SELECT id, CONCAT(first_name, " ", IFNULL(last_name, "")) as first_name FROM employee where status = 1 and empRole = 7  ORDER BY id desc`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error ", "success": false})
		return
	}
	defer rows.Close()

	var driver []listdriver
	var id, firstname string
	for rows.Next() {
		err := rows.Scan(&id, &firstname)

		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error ", "success": false})
			return

		}
		driver = append(driver, listdriver{Id: id, First_name: firstname})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"list": driver, "code": 200, "success": true, "message": "successfully"})
}
