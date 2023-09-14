package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfilePic   string `json:"profile_pic"`
	Role         string `json:"role"`
	RoleName     string `json:"role_name"`
	SupervisorID string `json:"supervisor_id"`
	TrainerType  string `json:"trainer_type,omitempty"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	Success      bool   `json:"success"`
}

//	type ResponseField struct {
//		Key   string
//		Value interface{}
//	}
type signInRequest struct {
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}

func SignInHandleRequest(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Invalid Method", "success": false})
		return
	}

	var request signInRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Invalid Request", "success": false, "error": err.Error()})

		return
	}

	Email := strings.TrimSpace(request.Email)

	// Update profile picture if provided
	if request.ProfilePic != "" {
		profileQuery := fmt.Sprintf("UPDATE employee SET profile_pic = '%s' WHERE officeMailId = '%s'", request.ProfilePic, Email)
		_, err := DB.Exec(profileQuery)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to update profile picture", "success": false, "error": err.Error()})
			return
		}
	}

	fields := "emp.id, emp.first_name, COALESCE(emp.last_name,'') as last_name, COALESCE(emp.profile_pic,''), COALESCE(emp.empRole,'') as role, COALESCE(rm.roleName,'') as role_name,emp.supervisorId"

	query := fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN roles_Master rm ON emp.empRole = rm.id WHERE emp.officeMailId = '%s' AND emp.status = 1", fields, Email)

	rows, err := DB.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
		return
	}
	defer rows.Close()

	if rows.Next() {
		var employee Employee
		err := rows.Scan(
			&employee.ID,
			&employee.FirstName,
			&employee.LastName,
			&employee.ProfilePic,
			&employee.Role,
			&employee.RoleName,
			&employee.SupervisorID,
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Scan Error", "success": false, "error": err.Error()})

			return
		}

		// Check senior trainer
		if employee.Role == "5" {
			query = fmt.Sprintf("SELECT id FROM employee WHERE status = 1 AND supervisorId = %s", employee.ID)
			rows, err := DB.Query(query)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			if rows.Next() {
				employee.TrainerType = "senior"
			}
		}

		var response map[string]interface{}

		if employee.TrainerType == "" {
			response = map[string]interface{}{

				// // var employee Employee
				"id":            employee.ID,
				"first_name":    employee.FirstName,
				"last_name":     employee.LastName,
				"profile_pic":   employee.ProfilePic,
				"role":          employee.Role,
				"role_name":     employee.RoleName,
				"supervisor_id": employee.SupervisorID,
				"code":          http.StatusOK,
				"message":       "successfully",
				"success":       true}
		} else {
			response = map[string]interface{}{

				// // var employee Employee
				"id":            employee.ID,
				"first_name":    employee.FirstName,
				"last_name":     employee.LastName,
				"profile_pic":   employee.ProfilePic,
				"role":          employee.Role,
				"role_name":     employee.RoleName,
				"supervisor_id": employee.SupervisorID,
				"trainer_type":  employee.TrainerType,
				"code":          http.StatusOK,
				"message":       "successfully",
				"success":       true}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "EmailID Not Found", "success": false})
	}
}
