package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

//=============================  Get Alter Bus =========================

func GetAlterBus(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	type Request struct {
		Bus_id string `json:"bus_id"`
	}

	type Responce struct {
		Bus_id          string `json:"bus_id"`
		Register_number string `json:"register_number"`
	}

	var p Request
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body", "success": false})
		return
	}

	// --------------- taking the project ids ----------------
	str := "SELECT prj.id FROM project prj where prj.busID =? limit 1"
	rows, err := DB.Query(str, p.Bus_id)
	if err != nil {
	
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

		return
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
			return
		}
	}

	//--------- responce what we need -----------------
	var assign_list []Responce // creating responce slice
	if id != 0 {
		str := "SELECT GROUP_CONCAT(DISTINCT(busID)) as ids  FROM project"
		rows, err := DB.Query(str)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

			return
		}
		defer rows.Close()

		var busid string

		for rows.Next() {

			err := rows.Scan(&busid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}
		}
		var notin string
		if busid != "" {
			notin = fmt.Sprintf(" where id NOT IN (%s)", busid)
		}

		query_bus := "SELECT id as bus_id, register_number FROM bus`" + notin + "`ORDER BY id desc"

		rows1, err := DB.Query(query_bus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

			return
		}
		defer rows1.Close()

		var bus_id, register_number string
		for rows1.Next() {
			err := rows1.Scan(&bus_id, &register_number)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}
			assign_list = append(assign_list, Responce{Bus_id: bus_id, Register_number: register_number})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"assign_list": assign_list, "code": 409, "success": false, "message": "Select alter Bus for the projects"})
	} else {

		json.NewEncoder(w).Encode(map[string]interface{}{"assign_list": assign_list, "code": 200, "success": true, "message": "Bus Changed Successfully"})
	}
}
