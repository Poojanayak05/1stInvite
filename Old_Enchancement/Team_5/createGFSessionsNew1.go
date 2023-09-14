package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GFSessionNew1Request struct {
	ProjectID            string `json:"project_id"`
	UserID               string `json:"user_id"`
	LocationID           string `json:"locationId"`
	CircleID             string `json:"circle_id"`
	TbName               string `json:"tb_name"`
	NumberOfParticipants int    `json:"numOfParticipants"`
	TbID                 int    `json:"tb_id"`
	GfSessionType        string `json:"gf_session_type"`
	PlanDate             string `json:"plan_date"`
}
type GFSessionNew1Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func CreateGFSessionNew1(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	request := GFSessionNew1Request{}
	response := GFSessionNew1Response{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Error While Decoding Body. Bad Request", "code": http.StatusBadRequest, "error": err})
		return
	}

	GfSessionName := func(name, gfSessionType string) string {

		var gfType string
		switch gfSessionType {
		case "1":
			gfType = "CM"
		case "2":
			gfType = "VV"
		case "3":
			gfType = "BV"
		case "4":
			gfType = "SPS"
		case "5":
			gfType = "SPM1"
		case "6":
			gfType = "SPM2"
		case "7":
			gfType = "SPM3"
		case "8":
			gfType = "SPM4"
		case "9":
			gfType = "SPM5"
		case "10":
			gfType = "GPS"
		case "11":
			gfType = "GPM1"
		case "12":
			gfType = "GPM2"
		case "13":
			gfType = "GPM3"
		case "14":
			gfType = "GPM4"
		case "15":
			gfType = "GPM5"
		case "16":
			gfType = "VPS"
		case "17":
			gfType = "VPM1"
		case "18":
			gfType = "VPM2"
		case "19":
			gfType = "VPM3"
		case "20":
			gfType = "VPM4"
		case "21":
			gfType = "VPM5"
		default:
			gfType = ""
		}
		nameArr := strings.Split(name, "_")

		return fmt.Sprintf("%s_%s", nameArr[0], gfType)
	}

	typ := "2"
	planDate := time.Time{}
	planDateWithoutSec := time.Time{}
	name := GfSessionName(request.TbName, request.GfSessionType)
	dateFormat := "2006-01-02 15:04:05"

	if request.PlanDate == "" {
		response.Code = 406
		response.Message = "Plan date is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.LocationID == "" {
		response.Code = 406
		response.Message = "locationId is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.UserID == "" {
		response.Code = 406
		response.Message = "user_id is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.ProjectID == "" {
		response.Code = 406
		response.Message = "project_id is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.TbID > 0 {
		response.Code = 406
		response.Message = "Tb id is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.CircleID == "" {
		response.Code = 406
		response.Message = "Circle id is required in order to create a GFsession"
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.PlanDate != "" {
		planDate, err = time.Parse(dateFormat, request.PlanDate)
		if err != nil || planDate.IsZero() {
			response.Code = http.StatusBadRequest
			response.Message = "Invalid date format"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			response.Success = false
			return
		}
		planDateWithoutSecStr := planDate.Format("2006-01-02 15:04")
		planDateWithoutSec, _ = time.Parse("2006-01-02 15:04", planDateWithoutSecStr)

		if !planDate.Equal(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)) {
			query := fmt.Sprintf("SELECT COUNT(id) as count FROM tbl_poa WHERE date = '%s' AND user_id = %s AND status != '2'", planDateWithoutSec, request.UserID)

			row := db.QueryRow(query)

			var count int
			err = row.Scan(&count)
			if err != nil {
				response.Code = 500
				response.Message = "Error While Getting Count. " + err.Error()
				response.Success = false
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			if count > 0 {
				response.Code = 409
				response.Message = "Poa already exists for date & time."
				response.Success = false
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	stmt, err := db.Prepare("INSERT INTO tbl_poa (type, tb_id, name, user_id, project_id, location_id, participants, date, session_type, circle_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Message = "Error While Inserting Record in tbl_poa. " + err.Error()
		response.Success = false
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Execute the insert statement
	result, err := stmt.Exec(typ, request.TbID, name, request.UserID, request.ProjectID, request.LocationID, request.NumberOfParticipants, planDateWithoutSec, request.GfSessionType, request.CircleID)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Message = "Error While Executing Insert Query in tbl_poa. " + err.Error()
		response.Success = false
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Message = "Error While getting Last Inserted Id. " + err.Error()
		response.Success = false
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	gfSessionName := fmt.Sprintf("%s_%d", GfSessionName(request.TbName, request.GfSessionType), insertedID)
	_, err = db.Exec("UPDATE tbl_poa SET name = ?, primary_id = ? WHERE id = ?", gfSessionName, insertedID, insertedID)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Message = "Error While updating tbl_poa. " + err.Error()
		response.Success = false
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = true
	response.Code = 200
	response.Message = "GF Sessions Added Successfully"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
