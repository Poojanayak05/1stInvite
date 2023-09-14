package Team_1

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ReportingManager struct {
	Label string `json:"label"`
	ID    string `json:"id"`
	Role  string `json:"role"`
}

type Role struct {
	ID       int    `json:"id"`
	RoleName string `json:"roleName"`
}

type Employe struct {
	Role             Role             `json:"role"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	ContactNum       string           `json:"contactNum"`
	WorkNum          string           `json:"workNum"`
	OfficeMailID     string           `json:"officeMailId"`
	Address          string           `json:"address"`
	Address3         string           `json:"address3"`
	Address2         string           `json:"address2"`
	Pincode          string           `json:"pincode"`
	Gender           string           `json:"gender"`
	DOJ              string           `json:"doj"`
	ReportingManager ReportingManager `json:"reportingManager"`
	LicenseNumber    string           `json:"license_number"`
	ProjectList      []string         `json:"project_list"`
	EmpID            string           `json:"emp_id"`
	EmpRole          string           `json:"empRole"`
	SupervisorID     string           `json:"supervisorId"`
	ProfilePic       string           `json:"profile_pic"`
	Status           string           `json:"status"`
	CreatedBy        string           `json:"createdBy"`
	LastUpdatedBy    string           `json:"lastUpdatedBy"`
	PresentStatus    bool             `json:"present_status"`
	CountryID        string           `json:"countryId"`
}

func CreateUser(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var emp Employe
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin()
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return
	}

	// Check if email IDs already exist
	checkEmailQuery := `SELECT COUNT(*) FROM employee WHERE officeMailId = ?`
	checkEmailStmt, err := tx.Prepare(checkEmailQuery)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return

	}
	defer checkEmailStmt.Close()

	var count int
	err = checkEmailStmt.QueryRow(emp.OfficeMailID).Scan(&count)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return

	}

	if count > 0 {
		tx.Rollback()
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    409,
			"message": "Employee already exists with the given email IDs",
			"success": false,
		})
		return
	}

	// Insert the employee record
	insertQuery := `INSERT INTO employee(countryID, first_name, last_name, gender, doj, officeMailId, personalMailId, contactNum, workNum, address, address2, address3, pincode, empRole, license_number, supervisorId, profile_pic, status, createdBy, lastUpdatedBy) 
						VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return

	}
	defer stmt.Close()

	_, err = stmt.Exec(emp.FirstName, emp.LastName, emp.Gender, emp.DOJ, emp.OfficeMailID, emp.OfficeMailID, emp.ContactNum, emp.WorkNum, emp.Address, emp.Address2, emp.Address3, emp.Pincode, emp.EmpRole, emp.LicenseNumber, emp.SupervisorID, emp.ProfilePic, emp.Status, emp.CreatedBy, emp.LastUpdatedBy)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return

	}

	empID := getLastInsertID(tx)

	if len(emp.ProjectList) > 0 {
		projectEmpInsertQuery := `INSERT INTO project_emps(emp_id, project_id, role_id) VALUES (?, ?, ?)`
		projectEmpStmt, err := tx.Prepare(projectEmpInsertQuery)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		defer projectEmpStmt.Close()

		projectEmpHistInsertQuery := `INSERT INTO project_emps_history(emp_id, projectId, role_id) VALUES (?, ?, ?)`
		projectEmpHistStmt, err := tx.Prepare(projectEmpHistInsertQuery)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		defer projectEmpHistStmt.Close()

		for _, projectID := range emp.ProjectList {
			_, err = projectEmpStmt.Exec(empID, projectID, emp.EmpRole)
			if err != nil {
				tx.Rollback()

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
				return

			}

			_, err = projectEmpHistStmt.Exec(empID, projectID, emp.EmpRole)
			if err != nil {
				tx.Rollback()

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
				return

			}
		}
	}

	//create funder
	if emp.EmpRole == "8" {
		tx, err := DB.Begin()
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		insertQuery := `INSERT INTO funder(countryID, funderName, workPhone, mobilePhone, emailID, address, status, city, district, state, pincode, designation, createdBy, lastUpdatedBy) 
						VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		defer stmt.Close()

		_, err = stmt.Exec(emp.FirstName+""+emp.LastName, emp.WorkNum, emp.ContactNum, emp.OfficeMailID, emp.Address, emp.Status, emp.Address, emp.Address2, emp.Address3, emp.Pincode, emp.Role.RoleName, emp.CreatedBy, emp.LastUpdatedBy)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}

		err = tx.Commit()
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Funder Added Successfully",
			"success": true,
		})

	} else if emp.EmpRole == "9" {
		tx, err := DB.Begin()
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		// Insert the partner record
		insertQuery := `INSERT INTO partner (countryID, partnerName, phone, mobilePhone, emailID, address, status, city, district, state, pincode, designation, createdBy, lastUpdatedBy) 
	VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		stmt, err := tx.Prepare(insertQuery)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		defer stmt.Close()

		_, err = stmt.Exec(emp.FirstName+""+emp.LastName, emp.WorkNum, emp.ContactNum, emp.OfficeMailID, emp.Address, emp.Status, emp.Address, emp.Address2, emp.Address3, emp.Pincode, emp.Role.RoleName, emp.CreatedBy, emp.LastUpdatedBy)
		if err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}

		err = tx.Commit()
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
			return

		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Partner Added Successfully",
			"success": true,
		})
	}

	err = tx.Commit()
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return

	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "People Added Successfully",
		"success": true,
	})
}

func getLastInsertID(tx *sql.Tx) int64 {
	var empID int64
	err := tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&empID)
	if err != nil {
		log.Println(err)

	}
	return empID
}
