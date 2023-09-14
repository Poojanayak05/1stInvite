package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func DeletUser(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	type Request struct {
		Emp_id     string   `json:"Emp_id"`
		Delete     string   `json:"delete"`
		Project_id []string `json:"project_id"`
		Role_id    string   `json:"role_id"`
	}

	Responce := make(map[string]interface{})

	var p Request
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})

		return
	}

	//----------------- one if condition -----------------
	if p.Delete != "" {
		if p.Emp_id != "" && p.Delete == "0" {
			id := p.Emp_id

			//---------- reading data supervisor ---------------------
			query := "SELECT emp.id FROM employee emp where emp.supervisorId =? and  emp.status = 1"
			rows, err := DB.Query(query, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}
			defer rows.Close()

			var data_supervisor []map[string]interface{}
			for rows.Next() {
				empid := make(map[string]interface{})
				var id int
				err := rows.Scan(&id)
				if err != nil {
					//http.Error(w, err.Error(), http.StatusInternalServerError)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

					return
				}
				empid["id"] = id
				data_supervisor = append(data_supervisor, empid)

			}
			var result []map[string]interface{}
			//-------------------- reading operatiion  manager project id ----------------
			query = "SELECT prj.id FROM project prj where prj.operations_manager =?"
			rows1, err := DB.Query(query, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}
			defer rows1.Close()

			var operation_manager []map[string]interface{}
			for rows1.Next() {
				vrr := make(map[string]interface{})
				var prj_id string
				err := rows1.Scan(&prj_id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

					return

				}
				vrr["prj_id"] = prj_id
				operation_manager = append(operation_manager, vrr)
			}

			fields := "emp.id, emp.empRole, if((select id from project_emps prj_emp WHERE prj_emp.emp_id = emp.id LIMIT 1), 1, 0) as project_exist, if((select prj.driverID from project prj WHERE prj.driverID = emp.id LIMIT 1), 1, 0) as driver_exist"
			query = fmt.Sprintf("SELECT %s FROM employee emp where emp.status = 1 and emp.id =%s", fields, id)

			rows2, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}
			defer rows2.Close()

			//var data []map[string]interface{}
			var empid, role, project_exist, driver_exist string
			for rows2.Next() {
				data1 := make(map[string]interface{})

				err := rows2.Scan(&empid, &role, &project_exist, &driver_exist)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				data1["empid"] = empid
				data1["role"] = role
				data1["project_exist"] = project_exist
				data1["driver_exist"] = driver_exist
				//data = append(data, data1)
			}
			var query_emp string
			if project_exist == "1" || driver_exist == "1" || len(data_supervisor) > 0 || len(operation_manager) > 0 {

				emps_fields := "DISTINCT(emp.id), emp.empRole as role_id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name,'')) AS first_name, rm.roleName"

				var projectids []int

				for _, projectIDStr := range p.Project_id {
					//projectIDStr1 := strings.Join(projectIDStr, ",")
					projectid, err := strconv.Atoi(projectIDStr)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

						return
					}

					projectids = append(projectids, projectid)
				}

				var strids []string
				for _, projectIDStr1 := range projectids {
					strids = append(strids, fmt.Sprintf("%d", projectIDStr1))
				}
				joinedIDs := strings.Join(strids, ",")

				//--------------------
				if p.Role_id == "5" || p.Role_id == "6" {
					//strings.Trim(strings.Join(strings.Fields(fmt.Sprint(employee.ProjectID)), ","), "[]")
					query_emp = fmt.Sprintf("SELECT %s FROM employee emp left join roles_Master rm on emp.empRole = rm.id left join project_emps prj_emp on emp.id = prj_emp.emp_id where emp.empRole = %s and emp.status = 1 and emp.id NOT IN (%s) and prj_emp.project_id IN (%s) ", emps_fields, role, id, joinedIDs)

				} else if p.Role_id == "7" {

					var notIn string

					driverQuery := "SELECT GROUP_CONCAT(DISTINCT(driverID)) as ids  FROM project where status = 1"

					rows, err := DB.Query(driverQuery)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

						return

					}
					defer rows.Close()

					var id2 []map[string]interface{}
					id1 := make(map[string]interface{})
					var driverid string
					for rows.Next() {
						var driverid string

						err := rows.Scan(&driverid)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

							return
						}
						id1["driver"] = driverid
					}

					id2 = append(id2, id1)
					if len(id2) == 0 {
						notIn = ""
					} else {
						notIn = fmt.Sprintf(" emp.id NOT IN (%s)", driverid)
					}
					query_emp = fmt.Sprintf("SELECT %s FROM employee emp left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.empRole = 7 %s ORDER BY emp.id desc", emps_fields, notIn)
				} else {
					query_emp = fmt.Sprintf("SELECT %s FROM employee emp left join roles_Master rm on emp.empRole = rm.id where emp.empRole = %s and emp.status = 1 and emp.id NOT IN (%s) ORDER BY emp.id", emps_fields, role, id)
				}

				rows, err := DB.Query(query_emp)
				if err != nil {

					w.WriteHeader(409)
					json.NewEncoder(w).Encode(map[string]interface{}{"assign_list": rows, "code": 409, "message": "Projects Assigned To Employee Can't Change Role", "success": false})
					return
				}
				defer rows.Close()

				var empid, role_id, first_name, roleName string
				for rows.Next() {
					result1 := make(map[string]interface{})
					err := rows.Scan(&empid, &role_id, &first_name, &roleName)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}

					result1["id"] = empid
					result1["role_id"] = role_id
					result1["first_name"] = first_name
					result1["roleName"] = roleName
					result = append(result, result1)
				}
				Responce["assign_list"] = result
				Responce["code"] = 409
				Responce["message"] = "Projects Assigned To Employee Can't Change Role"
				Responce["success"] = false
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{"assign_list": result, "code": 409, "message": "Projects Assigned To Employee Can't Change Role", "success": false})
			}

			//--- another if condition ---------------
		}

		//-------------- -- another if condition -------------
		if p.Emp_id != "" && p.Delete == "1" {
			id := p.Emp_id

			deleteQuery := "UPDATE employee SET status = '0' WHERE id = ?"

			_, err := DB.Exec(deleteQuery, id)
			if err != nil {

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}
			if p.Role_id == "5" || p.Role_id == "6" {
				// Delete from project_emps
				deleteProjectEmpsQuery := fmt.Sprintf("DELETE FROM project_emps WHERE emp_id = %s", id)
				_, err := DB.Exec(deleteProjectEmpsQuery)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}

				// Update project_emps_history
				updateHistoryQuery := fmt.Sprintf("UPDATE project_emps_history SET endDate = CURDATE() WHERE emp_id = %s", id)
				_, err1 := DB.Exec(updateHistoryQuery)
				if err1 != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
			}

			Responce["code"] = 200
			Responce["message"] = "People Deleted Successfully"
			Responce["success"] = true
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "People Deleted Successfully", "success": true})

		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "Sending Empty in delete parameter", "success": true})

	}
}
