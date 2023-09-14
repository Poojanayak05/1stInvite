package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TrainingBatchResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type TrainingBatchRequest struct {
	BatchName            string `json:"batch_name"`
	SubVillage           string `json:"sub_village"`
	TrainerID            string `json:"trainer_id"`
	ProjectID            string `json:"project_id"`
	LocationID           string `json:"location_id"`
	NumberOfParticipants string `json:"number_of_participants"`
	ContactPerson        string `json:"contact_person"`
	ContactNumber        string `json:"contact_number"`
	Day1                 string `json:"day1"`
	Day2                 string `json:"day2"`
}

func CreateTrainingBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// Parse request data
	var request TrainingBatchRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Error While Decoding Body. Bad Request", "code": http.StatusBadRequest, "error": err})
		return
	}

	if request.TrainerID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Trainer Id", "success": false})
		return
	}
	if request.ProjectID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Project Id", "success": false})
		return
	}
	response := TrainingBatchResponse{}
	role := 5
	id := 0
	if role == 5 {

		err = db.QueryRow(`SELECT  id FROM employee WHERE id = ? and empRole = ?`, request.TrainerID, role).Scan(&id)

		if err != nil {
			response.Code = http.StatusInternalServerError
			response.Message = "Trainer id does not exist."
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return

		}
		validateDateTime := func(dateStr string) (bool, string) {
			date, err := time.Parse("2006/01/2 3:04 PM", dateStr)
			outputFormat := "2006-01-02 15:04:05"
			output := date.Format(outputFormat)
			return err == nil, output
		}
		// Check date time format
		isValied, day1 := validateDateTime(request.Day1)
		if !isValied {
			response.Code = 400
			response.Message = "Invalid date format for day1"
			response.Success = false
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if day1 already exists
		var count int
		err = db.QueryRow("SELECT COUNT(id) as count FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'", day1, request.TrainerID).Scan(&count)
		if err != nil {
			response.Code = 500
			response.Message = "Error While Getting Count For Day1."
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		if count > 0 {
			response.Code = 409
			response.Message = "Day1 already exists for date & time"
			response.Success = false
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check date time format for day2
		isValied, day2 := validateDateTime(request.Day2)
		if !isValied {
			response.Code = 400
			response.Message = "Invalid date format for day2"
			response.Success = false
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if day2 already exists
		err = db.QueryRow("SELECT COUNT(id) FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'", day2, request.TrainerID).Scan(&count)
		if err != nil {
			response.Code = 500
			response.Message = "Error While Getting Count For Day2."
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		if count > 0 {
			response.Code = 409
			response.Message = "Day2 already exists for date & time"
			response.Success = false
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Retrieve project end date
		var projEndDate string
		err = db.QueryRow("SELECT endDate FROM project WHERE id = ?", request.ProjectID).Scan(&projEndDate)
		if err != nil {
			response.Code = 500
			response.Message = "Error While Getting End Date Of Project."
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		projEndDate += " 23:59:59"
		today := time.Now().Format("2006-01-02 15:04:05")

		// Check if project is completed
		if today > projEndDate {
			response.Code = 409
			response.Message = "Cannot create a training batch on completed projects"
			response.Success = false
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if day2 is beyond the project end date
		if day2 > projEndDate {
			response.Code = 409
			response.Message = "Cannot create a training batch beyond the end date of the project"
			response.Success = false
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Insert into tbl_poa for day1
		insertStmt1, err := db.Prepare("INSERT INTO tbl_poa (type, sub_village, user_id, project_id, location_id, participants, contact_person, contact_number, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer insertStmt1.Close()

		res1, err := insertStmt1.Exec("1", request.SubVillage, request.TrainerID, request.ProjectID, request.LocationID, request.NumberOfParticipants, request.ContactPerson, request.ContactNumber, day1)
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		insertedID, err := res1.LastInsertId()
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Generate training batch name
		name := request.BatchName
		if request.SubVillage != "" {
			name = request.SubVillage
		}
		name = fmt.Sprintf("%s_TB%d", name, insertedID)

		// Update tbl_poa with training batch name and tb_id for day1
		updateStmt, err := db.Prepare("UPDATE tbl_poa SET name = ?, tb_id = ?, primary_id = ? WHERE id = ?")
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer updateStmt.Close()

		_, err = updateStmt.Exec(name, insertedID, insertedID, insertedID)
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Insert into tbl_poa for day2
		insertStmt2, err := db.Prepare("INSERT INTO tbl_poa (type, tb_id, name, sub_village, user_id, project_id, location_id, participants, contact_person, contact_number, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer insertStmt2.Close()

		res2, err := insertStmt2.Exec("1", insertedID, name, request.SubVillage, request.TrainerID, request.ProjectID, request.LocationID, request.NumberOfParticipants, request.ContactPerson, request.ContactNumber, day2)
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		insertedID2, err := res2.LastInsertId()
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Update tbl_poa with primary_id for day2
		_, err = db.Exec("UPDATE tbl_poa SET primary_id = ? WHERE id = ?", insertedID2, insertedID2)
		if err != nil {
			response.Code = 500
			response.Message = err.Error()
			response.Success = false
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response.Code = 200
		response.Message = "Training Batch Added Successfully"
		response.Success = true
		json.NewEncoder(w).Encode(response)
	}
}
