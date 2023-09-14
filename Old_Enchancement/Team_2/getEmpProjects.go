package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project_list struct {
	User_id string `json:"user_id"`
	Role_id string `json:"role_id"`
}

type list struct {
	Id              string `json:"id"`
	Projectname     string `json:"projectName"`
	Bus_id          string `json:"bus_id"`
	Register_number string `json:"register_number"`
}

func GetEmppPojects(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	var p Project_list
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})

		return
	}

	var project_list []list
	var query_project string
	if p.User_id != "" && p.Role_id != "" {
		if p.Role_id == "1" || p.Role_id == "2" || p.Role_id == "3" || p.Role_id == "4" || p.Role_id == "5" || p.Role_id == "6" || p.Role_id == "7" || p.Role_id == "12" {

			if p.Role_id == "1" || p.Role_id == "2" {
				fields := "prj.id, COALESCE(UPPER(prj.projectName),'') as projectName, COALESCE(prj.busID,'') as bus_id, COALESCE(bus.register_number,'')"
				query_project = fmt.Sprintf("SELECT %s from project prj LEFT JOIN bus on prj.busID = bus.id where prj.project_status = '1' AND prj.status = 1", fields)

			}

			if p.Role_id == "3" || p.Role_id == "12" {

				fields := "prj.id, COALESCE(UPPER(prj.projectName),'') as projectName, COALESCE(prj.busID,'') as bus_id, COALESCE(bus.register_number,'')"
				query_project = fmt.Sprintf("SELECT %s from employee emp join project prj on emp.id = prj.operations_manager  LEFT JOIN bus on prj.busID = bus.id where emp.empRole = 4 and emp.supervisorId = %s", fields, p.User_id)
			}

			if p.Role_id == "4" {
				fields := "prj.id, COALESCE(UPPER(prj.projectName),'') as projectName, COALESCE(prj.busID,'') as bus_id, COALESCE(bus.register_number,'')"
				query_project = fmt.Sprintf("SELECT %s from project prj LEFT JOIN bus on prj.busID = bus.id where prj.operations_manager = %s", fields, p.User_id)
			}

			if p.Role_id == "5" || p.Role_id == "6" {
				fields := "em_pr.project_id as id, COALESCE(UPPER(pr.projectName),'') as projectName, COALESCE(pr.busID,'') as bus_id,COALESCE(bus.register_number,'')"
				query_project = fmt.Sprintf("SELECT %s from project_emps em_pr left join project pr on em_pr.project_id = pr.id LEFT JOIN bus on pr.busID = bus.id where emp_id IN (%s)", fields, p.User_id)
			}

			if p.Role_id == "7" {
				fields := "prj.id, COALESCE(UPPER(prj.projectName),'') as projectName, COALESCE(prj.busID,'') as bus_id, COALESCE(bus.register_number,'')"
				query_project = fmt.Sprintf("SELECT %s from project prj left join employee emp on prj.driverID = emp.id LEFT JOIN bus on prj.busID = bus.id where prj.driverID =  (%s)", fields, p.User_id)
			}

			rows, err := DB.Query(query_project)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
			}
			defer rows.Close()

			for rows.Next() {
				var id, projectname, busid, register_number string
				err := rows.Scan(&id, &projectname, &busid, &register_number)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				}

				project_list = append(project_list, list{Id: id, Projectname: projectname, Bus_id: busid, Register_number: register_number})
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"project_list": project_list, "success": true, "message": "Successfully"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"project_list": project_list, "success": true, "message": "Successfully"})
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Send Required Parameter"})
	}
}
