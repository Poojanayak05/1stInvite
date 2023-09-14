package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Request struct {
	Flag1           string `json:"flag"`
	Participant_Id1 string `json:"participant_id"`
	Tb_Poa_Id1      string `json:"tbl_poa_id"`
	Type1           string `json:"type"`
}

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ParticipantsAttendance(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "Message": "Method Not Found", "success": false})
		return
	}
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Bad Request", "success": false})
		return
	}

	Flag, err := strconv.Atoi(request.Flag1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Flag", "success": false})
		return
	}
	Participant_Id, err := strconv.Atoi(request.Participant_Id1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Participant ID", "success": false})
		return
	}
	Tb_Poa_Id, err := strconv.Atoi(request.Tb_Poa_Id1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Tb_Poa_Id", "success": false})
		return
	}
	Type, err := strconv.Atoi(request.Type1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid Type", "success": false})
		return
	}

	if Participant_Id > 0 && Tb_Poa_Id > 0 && Type > 0 {
		if Flag == 1 {
			// Check if participant ID already exists
			query := fmt.Sprintf("SELECT COUNT(*) FROM participant_attendance WHERE tbl_poa_id = %d AND participant_id = %d", Tb_Poa_Id, Participant_Id)
			var count int
			err := DB.QueryRow(query).Scan(&count)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Database error", "success": false})
				return
			}

			if count > 0 {
				w.WriteHeader(http.StatusOK)
				response := Response{
					Code:    200,
					Success: false,
					Message: "Attendance already added for the participant",
				}
				json.NewEncoder(w).Encode(response)
			} else {
				insertQuery := fmt.Sprintf("INSERT INTO participant_attendance (type, tbl_poa_id, participant_id) VALUES (%d, %d, %d)", Type, Tb_Poa_Id, Participant_Id)
				_, err := DB.Exec(insertQuery)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Failed to add attendance", "success": false})
					return
				}

				response := Response{
					Code:    200,
					Success: true,
					Message: "Attendance Added Successfully",
				}
				json.NewEncoder(w).Encode(response)
			}
		} else {
			// Check if participant ID exists
			query := fmt.Sprintf("SELECT COUNT(*) FROM participant_attendance WHERE tbl_poa_id = %d AND participant_id = %d", Tb_Poa_Id, Participant_Id)
			var count int
			err := DB.QueryRow(query).Scan(&count)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Database error", "success": false})
				return
			}

			if count == 0 {
				response := Response{
					Code:    200,
					Success: false,
					Message: "Invalid Participant ID",
				}
				json.NewEncoder(w).Encode(response)
			} else {
				deleteQuery := fmt.Sprintf("DELETE FROM participant_attendance WHERE tbl_poa_id = %d AND participant_id = %d", Tb_Poa_Id, Participant_Id)
				_, err := DB.Exec(deleteQuery)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Failed to delete attendance", "success": false})
					return
				}

				response := Response{
					Code:    200,
					Success: true,
					Message: "Attendance Deleted Successfully",
				}
				json.NewEncoder(w).Encode(response)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Required parameters not passed", "Success": false})
	}
}
