package Team_3

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Responsef struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Employee struct {
	ID             int      `json:"id"`
	CountryID      string   `json:"countryID"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Gender         string   `json:"gender"`
	Doj            string   `json:"doj"`
	OfficeMailID   string   `json:"officeMailId"`
	PersonalMailID string   `json:"personalMailId"`
	ContactNum     string   `json:"contactNum"`
	WorkNum        string   `json:"workNum"`
	Address        string   `json:"address"`
	Address2       string   `json:"address2"`
	Address3       string   `json:"address3"`
	Pincode        string   `json:"pincode"`
	EmpRole        string   `json:"empRole"`
	Role           string   `json:"role"`
	SupervisorID   string   `json:"supervisorId"`
	ProfilePic     string   `json:"profile_pic"`
	Status         string   `json:"status"`
	CreatedBy      string   `json:"createdBy"`
	LastUpdatedBy  string   `json:"lastUpdatedBy"`
	LicenseNumber  string   `json:"license_number"`
	ProjectList    []string `json:"project_list"`
}

func EditUser(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		showResponse(w, "Invalid Request Body", false)
		return
	}

	var request Employee
	err = json.Unmarshal(data, &request)
	if err != nil {
		showResponse(w, "Invalid JSON", false)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		showResponse(w, "checkkk error ", false)
		return
	}

	updateQuery := `
		UPDATE employee
		SET countryID = ?, first_name = ?, last_name = ?, gender = ?, doj = ?, officeMailId = ?,
			personalMailId = ?, contactNum = ?, workNum = ?, address = ?, address2 = ?,
			address3 = ?, pincode = ?, empRole = ?, supervisorId = ?, profile_pic = ?, status = ?,
			createdBy = ?, lastUpdatedBy = ?, license_number = ?
		WHERE id = ?
	`

	_, err = tx.Exec(updateQuery, request.CountryID, request.FirstName, request.LastName, request.Gender, request.Doj,
		request.OfficeMailID, request.PersonalMailID, request.ContactNum, request.WorkNum, request.Address,
		request.Address2, request.Address3, request.Pincode, request.Role, request.SupervisorID,
		request.ProfilePic, request.Status, request.CreatedBy, request.LastUpdatedBy, request.LicenseNumber,
		request.ID)
	if err != nil {
		tx.Rollback()
		showResponse(w, "Update Error", false)
		return
	}

	if len(request.ProjectList) > 0 {
		projectExist := []string{}
		query := "SELECT project_id FROM project_emps WHERE emp_id = ?"
		rows, err := tx.Query(query, request.ID)
		if err != nil {
			tx.Rollback()
			showResponse(w, "Fetch Projects Error", false)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var projectID string
			err := rows.Scan(&projectID)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Scan Projects Error", false)
				return
			}
			projectExist = append(projectExist, projectID)
		}

		if len(projectExist) > 0 && projectExist[0] != request.Role {
			updateRole := "UPDATE project_emps SET role_id = ? WHERE emp_id = ?"
			_, err = tx.Exec(updateRole, request.Role, request.ID)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Update Role Error", false)
				return
			}

			projectIDsStr := "'" + strings.Join(projectExist, "','") + "'"
			updateRoleHistory := `
				UPDATE project_emps_history
				SET role_id = ?
				WHERE projectId IN (` + projectIDsStr + `) AND emp_id = ? AND endDate IS NULL
			`
			_, err = tx.Exec(updateRoleHistory, request.Role, request.ID)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Update Role History Error", false)
				return
			}
		}

		deleteResult := sliceDiff(projectExist, request.ProjectList)
		updateResult := sliceDiff(request.ProjectList, projectExist)
		
		for _, projectID := range deleteResult {
			deleteQuery := "DELETE FROM project_emps WHERE project_id = ? AND emp_id = ?"
			_, err = tx.Exec(deleteQuery, projectID, request.ID)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Delete Project Error", false)
				return
			}

			updateHistory := `
				UPDATE project_emps_history
				SET endDate = CURDATE()
				WHERE projectId = ? AND emp_id = ? AND role_id = ?
			`
			_, err = tx.Exec(updateHistory, projectID, request.ID, request.Role)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Update Project History Error", false)
				return
			}
		}

		for _, projectID := range updateResult {
			empInsertQuery := "INSERT INTO project_emps (emp_id, project_id, role_id) VALUES (?, ?, ?)"
			_, err = tx.Exec(empInsertQuery, request.ID, projectID, request.Role)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Insert Project Error", false)
				return
			}

			empHistInsertQuery := "INSERT INTO project_emps_history (emp_id, projectId, role_id) VALUES (?, ?, ?)"
			_, err = tx.Exec(empHistInsertQuery, request.ID, projectID, request.Role)
			if err != nil {
				tx.Rollback()
				showResponse(w, "Insert Project History Error", false)
				return
			}
		}
	} else {
		deleteQuery := "DELETE FROM project_emps WHERE emp_id = ?"
		_, err = tx.Exec(deleteQuery, request.ID)
		if err != nil {
			tx.Rollback()
			showResponse(w, "Delete Projects Error", false)
			return
		}

		updateHistory := "UPDATE project_emps_history SET endDate = CURDATE() WHERE emp_id = ?"
		_, err = tx.Exec(updateHistory, request.ID)
		if err != nil {
			tx.Rollback()
			showResponse(w, "Update Project History Error", false)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		showResponse(w, "Commit Transaction Error", false)
		return
	}

	showResponse(w, "Updated successfully", true)
}

func sliceDiff(slice1, slice2 []string) []string {
	// Create a map to store elements from slice2
	// for faster lookups
	slice2Map := make(map[string]bool)
	for _, value := range slice2 {
		slice2Map[value] = true
	}

	// Create a new slice to store the difference
	var diff []string

	// Iterate over slice1 and check if each element is present in slice2
	// If not present, add it to the difference slice
	for _, value := range slice1 {
		if _, found := slice2Map[value]; !found {
			diff = append(diff, value)
		}
	}

	return diff
}
func showResponse(w http.ResponseWriter, message string, success bool) {
	response := Responsef{
		Message: message,
		Success: success,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal JSON response:", err)
		http.Error(w, "JSON Encoding Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
