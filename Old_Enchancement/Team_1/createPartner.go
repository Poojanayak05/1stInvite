package Team_1

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Partner struct {
	CountryID     string `json:"countryID"`
	PartnerName   string `json:"name"`
	Phone         string `json:"phone"`
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

func CreatePartner(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var partner Partner
	err := json.NewDecoder(r.Body).Decode(&partner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	// Check if email ID already exists
	checkEmailQuery := `SELECT COUNT(*) FROM partner WHERE emailID = ?`
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
	err = checkEmailStmt.QueryRow(partner.EmailID).Scan(&count)
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
			"message": "Partner already exists with the given email ID",
			"success": false,
		})
		return
	}

	// Insert the partner record
	insertQuery := `INSERT INTO partner (countryID, partnerName, phone, mobilePhone, emailID, address, status, city, district, state, pincode, designation, createdBy, lastUpdatedBy) 
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

	_, err = stmt.Exec(partner.CountryID, partner.PartnerName, partner.Phone, partner.MobilePhone, partner.EmailID, partner.Address, partner.Status, partner.City, partner.District, partner.State, partner.Pincode, partner.Designation, partner.CreatedBy, partner.LastUpdatedBy)
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
		"message": "Partner Added Successfully",
		"success": true,
	})
}
