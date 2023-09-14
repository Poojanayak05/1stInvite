package Team_1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	NotesId string `json:"notes_id"`
	Notes   string `json:"notes"`
	Emp_Id  string `json:"emp_id"`
	Date    string `json:"date"`
	Name    string `json:"name"`
}

type NRequest struct {
	Type1      string `json:"type"`
	Tb_id1     string `json:"tb_id"`
	Primary_id string `json:"primary_id"`
}

type NResponse struct {
	Notes   []Note `json:"notes"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func GetNotes(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var request NRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Bad Request", "success": false})
		return
	}

	Type, err := strconv.Atoi(request.Type1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Type", "success": false})
		return
	}

	var notes []Note
	//================================ tb_id is there value
	if request.Tb_id1 != "" {
		Tb_id, err := strconv.Atoi(request.Tb_id1)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Tb_id1", "success": false})
			return
		}

		rows, err := DB.Query("SELECT nt.id AS notes_id, nt.notes, nt.emp_id, nt.date, emp.first_name AS name FROM notes nt LEFT JOIN employee emp ON nt.emp_id = emp.id WHERE nt.type = ? AND nt.tb_id = ? ORDER BY nt.id DESC", Type, Tb_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to retrieve notes", "success": false})
			return
		}

		defer rows.Close()

		for rows.Next() {
			var note Note
			err := rows.Scan(&note.NotesId, &note.Notes, &note.Emp_Id, &note.Date, &note.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to scan row", "success": false})
				return
			}

			//
			// Parse date string to time.Time
			t, err := time.Parse("2006-01-02 15:04:05", note.Date)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to parse date", "success": false})
				return
			}

			// Format time as 12-hour time
			note.Date = t.Format("02-01-2006 03:04 PM")

			notes = append(notes, note)
		}

		// if len(notes) == 0 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Tb_id or Type", "success": false})
		// 	return
		// }

		response := NResponse{
			Notes: notes,
			Code:  200,

			Success: true,
			Message: "Successfully",
		}

		w.Header().Set("Content_Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to encode the response", "success": false})
			return
		}
	} else {
		primary_id, err := strconv.Atoi(request.Primary_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Primary id", "success": false})
			return
		}
		rows, err := DB.Query("SELECT nt.id AS notes_id, nt.notes, nt.emp_id, nt.date, emp.first_name AS name FROM notes nt LEFT JOIN employee emp ON nt.emp_id = emp.id WHERE nt.type = ?  AND nt.primary_id=? ORDER BY nt.id DESC", Type, primary_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to retrieve notes", "success": false})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var note Note
			err := rows.Scan(&note.NotesId, &note.Notes, &note.Emp_Id, &note.Date, &note.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to scan row", "success": false})
				return
			}

			// Parse date string to time.Time
			t, err := time.Parse("2006-01-02 15:04:05", note.Date)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to parse date", "success": false})
				return
			}

			// Format time as 12-hour time
			note.Date = t.Format("02-01-2006 03:04 PM")

			notes = append(notes, note)
		}

		// if len(notes) == 0 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Tb_id or Type", "success": false})
		// 	return
		// }

		response := NResponse{
			Notes:   notes,
			Code:    200,
			Success: true,
			Message: "Successfully",
		}
		w.Header().Set("Content_Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to encode the response", "success": false})
			return
		}

	}
}
