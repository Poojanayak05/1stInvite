package Team_3

import (
	// p "buzzstaff/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetPeopleList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	type Employee struct {
		ID               string `json:"id"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		RoleID           string `json:"role_id"`
		Flag             string `json:"flag"`
		VillagesAssigned string `json:"villages_assigned"`
	}

	type Response struct {
		List       []Employee `json:"list"`
		ProjArray  []string   `json:"projArray"`
		TotalCount string     `json:"total_count"`
		Code       int        `json:"code"`
		Message    string     `json:"message"`
		Success    bool       `json:"success"`
	}

	type req struct {
		RoleID               string `json:"role_id"`
		ProjectID            string `json:"project_id"`
		Operation_manager_id string `json:"operation_manager_id"`
		Assign               string `json:"assign"`
		PageNum              string `json:"pageNum"`
		Name                 string `json:"name"`
	}
	var request req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	RoleID, err := strconv.Atoi(request.RoleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid ID", "success": false, "error": err})
		return
	}
	ProjectID, err := strconv.Atoi(request.ProjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID", "success": false, "error": err})
		return
	}
	Operation_manager_id, err := strconv.Atoi(request.Operation_manager_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Operation Manager ID", "success": false, "error": err})
		return
	}
	// Assign, err := strconv.Atoi(request.Assign)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Assign value", "success": false, "error": err})
	// 	return
	// }
	PageNum, _ := strconv.Atoi(request.PageNum)
	if PageNum == 0 {
		PageNum = 1
	}

	noOfRecords := 100000
	// offset := (PageNum - 1) * noOfRecords
	offset := (noOfRecords * PageNum) - noOfRecords

	search := ""
	if request.Name != "" {
		Name := request.Name
		search = fmt.Sprintf("AND first_name LIKE '%%%s%%'", Name)
	}

	if RoleID != 0 && ProjectID != 0 && Operation_manager_id != 0 {
		ProjArray, err := GetAssociatedProjectList(ProjectID, DB)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while getting associated project list", "success": false, "error": err})
			return
		}

		stringOpsIds := make([]string, len(ProjArray))
		for i, id := range ProjArray {
			stringOpsIds[i] = strconv.Itoa(id)
		}
		result := strings.Join(stringOpsIds, ",")
	

		OperationManagerID := Operation_manager_id

		if RoleID == 13 {
			RoleID = 6
			t, err := DB.Query(fmt.Sprintf("SELECT id FROM employee WHERE supervisorId IN (%d)", OperationManagerID))
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while executing the query", "success": false, "error": err})
				return
			}
			defer t.Close()

			var operationManagerIDs string
			for t.Next() {
				var managerID string
				if err := t.Scan(&managerID); err != nil {
					log.Println(err)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while reading manager ID", "success": false, "error": err})
					return
				}
				operationManagerIDs += "," + managerID
			}
		}

		totalPagesQuery := fmt.Sprintf("SELECT COUNT(*) FROM employee WHERE status = 1 AND empRole = %d AND supervisorId = %d %s", RoleID, Operation_manager_id, search)
		var totalRows int
		err = DB.QueryRow(totalPagesQuery).Scan(&totalRows)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while executing the totalpagesQuery", "success": false, "error": err})
			return
		}

		fields := " DISTINCT emp.id as id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) as first_name, IFNULL(emp.last_name,'') as last_name, emp.empRole as role_id, if((select COUNT(id) from project_emps prj_emp WHERE prj_emp.emp_id = emp.id and prj_emp.project_id IN (" + result + ")), 1, 0) as flag, (select count(id) from gf_batches WHERE emp_id = emp.id and project_id IN (" + result + ")) as villages_assigned"
		x := 0
		if RoleID == 12 || RoleID == 4 {
			x = 1
		}
		// if request.RoleID == "5" || request.RoleID == "6" {
		// 	x = 2
		// }
		
		if request.Assign != "" {
			if request.Assign == "1" {
				RoleID = 6
			}

			if x == 1 {
				query := fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN project_emps ON emp.id = project_emps.emp_id WHERE emp.status = 1 AND emp.empRole = %d AND project_emps.project_id IN ("+result+") ORDER BY first_name LIMIT %d offset %d",
					fields, RoleID, noOfRecords, offset)
				
				rows, err := DB.Query(query)
				if err != nil {
					log.Println(err)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while executing the query", "success": false, "error": err})
					return
				}
				defer rows.Close()

				employees := make([]Employee, 0)
				for rows.Next() {
					var emp Employee
					err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.RoleID, &emp.Flag, &emp.VillagesAssigned)
					if err != nil {
						log.Println(err)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "Error while scanning", "success": false, "error": err})
						return
					}
					employees = append(employees, emp)
				}
				
				response := Response{
					List:       employees,
					ProjArray:  stringOpsIds,
					TotalCount: strconv.Itoa(totalRows),
					Code:       http.StatusOK,
					Message:    "Successfully",
					Success:    true,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)

			} else {
				query := fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN project_emps ON emp.id = project_emps.emp_id WHERE emp.status = 1 AND emp.empRole = %d AND emp.supervisorId IN (%d) AND project_emps.project_id IN ("+result+") %s ORDER BY first_name LIMIT %d offset %d",
					fields, RoleID, OperationManagerID, search, noOfRecords, offset)
				
				rows, err := DB.Query(query)
				if err != nil {
					log.Println(err)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while executing the query", "success": false, "error": err})
					return
				}
				defer rows.Close()

				employees1 := make([]Employee, 0)
				for rows.Next() {
					var emp Employee
					err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.RoleID, &emp.Flag, &emp.VillagesAssigned)
					if err != nil {
						log.Println(err)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "Error while scanning", "success": false, "error": err})
						return
					}
					employees1 = append(employees1, emp)
					
				}

				
				response := Response{
					List:       employees1,
					ProjArray:  stringOpsIds,
					TotalCount: strconv.Itoa(totalRows),
					Code:       http.StatusOK,
					Message:    "Successfully",
					Success:    true,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		} else {
			query := fmt.Sprintf("SELECT %s FROM employee emp WHERE emp.status = 1 AND emp.empRole = %d AND emp.supervisorId IN (%d) %s ORDER BY first_name LIMIT %d offset %d ",
				fields, RoleID, OperationManagerID, search, noOfRecords, offset)
		
			rows, err := DB.Query(query)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error while executing the query", "success": false, "error": err})
				return
			}
			defer rows.Close()

			employees2 := make([]Employee, 0)
			for rows.Next() {
				var emp Employee
				err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.RoleID, &emp.Flag, &emp.VillagesAssigned)
				if err != nil {
					log.Println(err)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "Error while scanning", "success": false, "error": err})
					return
				}
				employees2 = append(employees2, emp)
			
			}

			response := Response{
				List:       employees2,
				ProjArray:  stringOpsIds,
				TotalCount: strconv.Itoa(totalRows),
				Code:       http.StatusOK,
				Message:    "Successfully",
				Success:    true,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter Valid Role ID, Project ID, or Operation Manager ID", "success": false})
		return
	}
}
