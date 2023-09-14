package Team_5

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// struct to store the response.
type VyaparResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// struct to store the parameter coming into the body.
type VyaparRequest struct {
	ID           string `json:"id"`
	TbID         string `json:"tb_id"`
	GelathiID    string `json:"gelathi_id"`
	EnrolledDate string `json:"enrolledDate"`
	ProjectID    string `json:"projectId"`
}

func SetVyaparEnrollment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// declaring req variable to refer fields of VyaparRequest struct.
	var req VyaparRequest
	// declarig res variable to refer fields of VyaparResponse struct.
	var res VyaparResponse

	// dumping the body parameters into VyaparRequest struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error While Decoding Body. Bad Request", "error": err})
		return
	}

	if req.ID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Id", "success": false})
		return
	}
	if req.TbID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Tb Id", "success": false})
		return
	}
	if req.GelathiID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Gelathi Id", "success": false})
		return
	}
	if req.ProjectID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Project Id", "success": false})
		return
	}

	// query to fetch id from training_participants table based on VyaparEnrollment , id and tb_id.
	checkIfEnrolledPreviously := "SELECT id FROM training_participants WHERE VyaparEnrollment = 1 AND id = ? AND tb_id = ?"
	var participantID int

	// executing and scanning the participant id.
	err = db.QueryRow(checkIfEnrolledPreviously, req.ID, req.TbID).Scan(&participantID)
	if err != nil && err != sql.ErrNoRows {

		res.Message = "Failed to execute database query."
		res.Code = 500
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var updateSet string
	enrolledDate := time.Now().Format("2006-01-02 15:04:05")
	// if participant id is present then if updatequery will be executing otherwise else updatequery will execute.
	if participantID > 0 {
		updateSet = "UPDATE training_participants SET VyaparEnrollment = '0', gelathi_id = ?, VyaparEnrollmentDate = ?, VyaparEnrollmentEnrolledProject = ? WHERE tb_id = ? AND id = ?"
	} else {
		updateSet = "UPDATE training_participants SET VyaparEnrollment = '1', gelathi_id = ?, VyaparEnrollmentDate = ?, VyaparEnrollmentEnrolledProject = ? WHERE tb_id = ? AND id = ?"
	}
	//executing query.
	_, err = db.Exec(updateSet, req.GelathiID, enrolledDate, req.ProjectID, req.TbID, req.ID)
	if err != nil {
		res.Message = "Failed to update database."
		res.Code = 500
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	//setting the response.
	res.Code = 200
	res.Success = true
	res.Message = "Vyapar Enrollment Updated Successfully."
	// sending the response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
