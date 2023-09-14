package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// struct to store the respnse.
type AddEmpToProjectResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// struct to store the parameters coming into body.
type AddEmpToProjectRequest struct {
	EmpId     string `json:"emp_id"`
	ProjectId string `json:"project_id"`
	RoleId    string `json:"role_id"`
}

func AddEmpToProject(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// declaring req variable to refer the request struct's fields.
	var req AddEmpToProjectRequest
	// declaring res variable to refer the response struct's fields.
	var res AddEmpToProjectResponse

	// dumping body data into struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error While Decoding Body. Bad Request", "error": err})
		return
	}

	if req.EmpId == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Emp Id", "success": false})
		return
	}

	if req.ProjectId == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Project Id", "success": false})
		return
	}

	if req.RoleId == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Role Id", "success": false})
		return
	}
	// query to check employee is already assigned to a project or not.
	checkAddedPreviously := "Select emp_id,project_id, role_id from project_emps where emp_id = ? and project_id = ? and role_id = ?"
	// executing the query.
	rows, err := db.Query(checkAddedPreviously, req.EmpId, req.ProjectId, req.RoleId)
	if err != nil {
		res.Message = "Failed to execute query"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var emp_id, project_id, role_id string
	i := 0
	for rows.Next() {
		i++
		// scanning the data.
		rows.Scan(&emp_id, &project_id, &role_id)
	}
	// if employee already assigned to a project then showing the response.
	if i > 0 {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusAccepted, "message": "Employee Already Registered to this project."})
		return
	}

	// query to assigning the employee to a project.
	projectEmpInsertQuery := fmt.Sprintf("INSERT INTO project_emps(emp_id, project_id, role_id) VALUES ('%s', '%s', '%s')", req.EmpId, req.ProjectId, req.RoleId)
	// executing the query.
	_, err = db.Exec(projectEmpInsertQuery)
	if err != nil {
		res.Message = "Failed to execute query"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// query to insert the data into project_emps_history table.
	projectEmpHistInsertQuery := fmt.Sprintf("INSERT INTO project_emps_history(emp_id, projectId, role_id) VALUES ('%s', '%s', '%s')", req.EmpId, req.ProjectId, req.RoleId)
	// executing query.
	_, err = db.Exec(projectEmpHistInsertQuery)
	if err != nil {
		res.Message = "Failed to execute query"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// setting the response.
	res.Message = "Added Successfully."
	res.Success = true
	// sending the response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
