package Team_1

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Funder struct {
	CountryID     string `json:"countryID"`
	FunderName    string `json:"name"`
	WorkPhone     string `json:"phone"`
	MobilePhone   string `json:"mobilePhone"`
	EmailID       string `json:"emailID"`
	Address       string `json:"address"`
	Status        string `json:"status"`
	City          string `json:"city"`
	District      string `json:"district"`
	State         string `json:"state"`
	Pincode       string `json:"pincode"`
	Designation   string `json:"designation"`
	CreatedBy     string `json:"createdBy"`
	LastUpdatedBy string `json:"lastUpdatedBy"`
}

func CreateFunder(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var funder Funder
	err := json.NewDecoder(r.Body).Decode(&funder)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid input"})
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	// Check if email ID already exists
	checkEmailQuery := `SELECT COUNT(*) FROM funder WHERE emailID = ?`
	checkEmailStmt, err := tx.Prepare(checkEmailQuery)
	if err != nil {
		tx.Rollback()
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	defer checkEmailStmt.Close()

	var count int
	err = checkEmailStmt.QueryRow(funder.EmailID).Scan(&count)
	if err != nil {
		tx.Rollback()
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	if count > 0 {
		tx.Rollback()
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    409,
			"message": "Funder already exists with the given email ID",
			"success": false,
		})
		return
	}

	// Insert the funder record
	insertQuery := `INSERT INTO funder(countryID, funderName, workPhone, mobilePhone, emailID, address, status, city, district, state, pincode, designation, createdBy, lastUpdatedBy) 
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		tx.Rollback()
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(funder.CountryID, funder.FunderName, funder.WorkPhone, funder.MobilePhone, funder.EmailID, funder.Address, funder.Status, funder.City, funder.District, funder.State, funder.Pincode, funder.Designation, funder.CreatedBy, funder.LastUpdatedBy)
	if err != nil {
		tx.Rollback()
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Funder Added Successfully",
		"success": true,
	})
}
