package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetChangeRole(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type Req struct {
		Emp_id  string `json:"emp_id"`
		Role_id string `json:"role_id"`
	}

	type RespData struct {
		Emp_id        int `json:"emp_id"`
		Role_id       int `json:"role_id"`
		Project_exist int `json:"project_exist"`
		Driver_exist  int `json:"driver_exist"`
	}

	type Emp_result struct {
		Id         int    `json:"id"`
		Role_id    int    `json:"role_id"`
		First_name string `json:"first_name"`
		RoleName   string `json:"role_name"`
	}

	var request Req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := make(map[string]interface{})

	// var roleid int
	// err = DB.QueryRow("select emprole from employee where id=?", request.Emp_id).Scan(&roleid)
	// if err != nil {

	// 	return
	// }

	//if roleid == request.Role_id {
	if request.Emp_id != "" && request.Role_id != "" {
		//if emp assigned as a supervisor
		query := fmt.Sprintf("SELECT emp.id FROM employee emp where emp.supervisorId = %s and  emp.status = 1", request.Emp_id)

		rows, err := DB.Query(query)
		if err != nil {
			if err == sql.ErrNoRows {
				// ID doesn't exist in tbl_poa
				response := Response{
					Code:    http.StatusNotFound,
					Msg:     "empID doesn't exist",
					Success: false,
				}
				jsonResponse(w, response, http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		var data_supervisor string
		for rows.Next() {
			rows.Scan(&data_supervisor)
		}

		// here i need to create for opeartions manager
		query = fmt.Sprintf("SELECT prj.id FROM project prj where prj.operations_manager =%s", request.Emp_id)

		rows, err = DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		var operation_manager string
		for rows.Next() {
			rows.Scan(&operation_manager)
		}

		fields := "emp.id, emp.empRole, if((select id from project_emps prj_emp WHERE prj_emp.emp_id = emp.id LIMIT 1), 1, 0) as project_exist, if((select prj.driverID from project prj WHERE prj.driverID = emp.id LIMIT 1), 1, 0) as driver_exist"

		query = fmt.Sprintf("SELECT %s FROM employee emp where emp.status = 1 and emp.id =%s", fields, request.Emp_id)

		rows, err = DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		var data RespData
		var count int
		for rows.Next() {
			rows.Scan(&data.Emp_id, &data.Role_id, &data.Project_exist, &data.Driver_exist)

			count = 1
		}
		if count == 0 {
			jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid employee id", "code": 400, "success": false})
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Write(jsonData)
			return
		}
		var roleExists bool
		err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM employee WHERE empRole = ?)", request.Role_id).Scan(&roleExists)
		if err != nil {

			// Print the error and return a JSON response with error message

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
			return
		}

		if !roleExists {
			// Print the RoleID not found message
			// Return a JSON response if RoleID doesn't exist
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "No role_id found"})
			return
		}

		if strconv.Itoa(data.Role_id) != request.Role_id && (data.Project_exist == 1 || data.Driver_exist == 1 || len(data_supervisor) > 0 || len(operation_manager) > 0) {
			emps_fields := "emp.id, emp.empRole as role_id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name,'')) AS first_name, rm.roleName"

			query_emp := fmt.Sprintf("SELECT %s FROM employee emp left join roles_Master rm on emp.empRole = rm.id where emp.empRole = %d and emp.status = 1 and emp.id NOT IN (%s) ORDER BY emp.id", emps_fields, data.Role_id, request.Emp_id)

			rows, err = DB.Query(query_emp)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			var list []Emp_result
			for rows.Next() {
				var emp Emp_result
				rows.Scan(&emp.Id, &emp.Role_id, &emp.First_name, &emp.RoleName)
				list = append(list, emp)
			}

			response["assign_list"] = list
			response["code"] = 409
			response["success"] = false
			response["message"] = "Select alter employee for the projects"
			json.NewEncoder(w).Encode(response)

		} else {
			response["code"] = http.StatusOK
			response["success"] = true
			response["message"] = "Role Changed Successfully"
			json.NewEncoder(w).Encode(response)

		}

	} else {
		jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid Role id for this Employee id", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			if err == sql.ErrNoRows {
				// ID doesn't exist in tbl_poa
				response := Response{
					Code:    http.StatusNotFound,
					Msg:     "RoleID doesn't exist",
					Success: false,
				}
				jsonResponse(w, response, http.StatusNotFound)
				return
			}

			// http.Error(w, "Failed to query database", http.StatusInternalServerError)

			return
		}
		w.Write(jsonData)
	}
}
