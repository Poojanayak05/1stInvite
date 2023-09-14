package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Function get Projects based on emps projects
func getEmpPrjs(empID string, roleID string, DB *sql.DB) string {

	var queryProject string
	var idslice []string

	if roleID == "1" || roleID == "2" {
		_, err := DB.Exec(" SET SESSION group_concat_max_len = 1000000; ")
		if err != nil {
			log.Println("ERROR>>", err)

		}

		queryProject = "SELECT COALESCE(GROUP_CONCAT(DISTINCT prj.id),0) as ids FROM project prj WHERE prj.project_status = '1'"

	} else if roleID == "3" || roleID == "12" {
		if roleID == "3" {
			_, err := DB.Exec(" SET SESSION group_concat_max_len = 1000000; ")
			if err != nil {
				log.Println("ERROR>>", err)

			}
			rows, err := DB.Query("SELECT id FROM employee e WHERE supervisorId = ? AND empRole = 12", empID)
			if err != nil {
				log.Println("ERROR>>", err)

			}
			defer rows.Close()

			var somID int
			if rows.Next() {
				err := rows.Scan(&somID)
				if err != nil {
					log.Println("ERROR>>", err)

				}
			}

			whr := fmt.Sprintf(" emp.supervisorId IN (%s, %d)", empID, somID)
			queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM employee emp JOIN project prj ON emp.id = prj.operations_manager WHERE emp.empRole = 4 AND %s", whr)
		} else {
			queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM employee emp JOIN project prj ON emp.id = prj.operations_manager WHERE emp.empRole = 4 AND emp.supervisorId = %s", empID)
		}
	} else if roleID == "4" || roleID == "13" {
		if roleID == "13" {
			_, err := DB.Exec(" SET SESSION group_concat_max_len = 1000000; ")
			if err != nil {
				log.Println("ERROR>>", err)

			}
			queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project prj WHERE prj.gfl_id = %s", empID)
		} else {
			_, err := DB.Exec(" SET SESSION group_concat_max_len = 1000000; ")
			if err != nil {
				log.Println("ERROR>>", err)

			}
			queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project prj WHERE prj.operations_manager = %s", empID)
		}
	} else if roleID == "5" || roleID == "6" { // gelathi and trainer
		queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id IN (%s)", empID)
	} else if roleID == "7" { // driver
		queryProject = fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project prj LEFT JOIN employee emp ON prj.driverID = emp.id WHERE prj.driverID = %s", empID)
	}

	rows, err := DB.Query(queryProject)
	if err != nil && err != sql.ErrNoRows {
		log.Println("ERROR>>", err)

	}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		idslice = append(idslice, id)
	}
	ids := strings.Join(idslice, ",")
	if ids == "" {
		ids = "''"
	}

	return ids
}

func GetPeopleFilters(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type req struct {
		Role_id     string `json:"role_id"`
		Filter_type string `json:"filter_type"`
		PageNum     int    `json:"pageNum"`
		Emp_id      string `json:"emp_id"`
	}

	var request req

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageno := request.PageNum
	if err != nil {
		pageno = 1
	}

	// noOfRecords := 25
	// // offset := (PageNum - 1) * noOfRecords
	// offset := (noOfRecords * PageNum) - noOfRecords

	noOfRecords := 1000000
	offset := (noOfRecords * pageno) - noOfRecords
	var roleid string
	err = DB.QueryRow("select empRole from employee where id=?", request.Emp_id).Scan(&roleid)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid RoleID",
			"success": false,
			"error":   err,
		})
		return
	}

	//if request.Emp_id != 0 && request.Role_id != 0 {
	if roleid == request.Role_id {

		projects := getEmpPrjs(request.Emp_id, request.Role_id, DB)

		var totalPagesSQL string
		var query string

		if request.Filter_type != "" {
			//partner = 1, funder = 2, project = 3, opm = 4, trainer = 5, gelathi = 6 SOM=12 GFl=13
			if request.Filter_type == "1" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj LEFT JOIN partner part ON prj.partnerID = part.partnerID WHERE prj.id IN (%s)", projects)
				query = fmt.Sprintf("SELECT DISTINCT part.partnerID AS id, part.partnerName AS name, '' AS role_id FROM project prj LEFT JOIN partner part ON prj.partnerID = part.partnerID WHERE part.partnerID != '' AND prj.id IN (%s) LIMIT %d, %d", projects, offset, noOfRecords)
			} else if request.Filter_type == "2" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj LEFT JOIN funder fund ON prj.funderID = fund.funderID WHERE prj.id IN (%s)", projects)
				query = fmt.Sprintf("SELECT DISTINCT fund.funderID AS id, fund.funderName AS name, '' AS role_id FROM project prj LEFT JOIN funder fund ON prj.funderID = fund.funderID WHERE prj.id IN (%s) LIMIT %d, %d", projects, offset, noOfRecords)
			} else if request.Filter_type == "3" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj WHERE prj.id IN (%s)", projects)
				query = fmt.Sprintf("SELECT DISTINCT prj.id AS id, prj.projectName AS name, '' AS role_id FROM project prj WHERE prj.id IN (%s) LIMIT %d, %d", projects, offset, noOfRecords)
			} else if request.Filter_type == "4" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(DISTINCT emp.id) FROM project prj LEFT JOIN employee emp ON prj.operations_manager = emp.id WHERE prj.id IN (%s)", projects)

				query = fmt.Sprintf("SELECT DISTINCT emp.id AS id, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') AS name, emp.empRole AS role_id FROM project prj LEFT JOIN employee emp ON prj.operations_manager = emp.id WHERE prj.id IN (%s) LIMIT %d, %d", projects, offset, noOfRecords)

			} else if request.Filter_type == "12" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(DISTINCT emp.id) FROM project prj LEFT JOIN employee emp ON prj.operations_manager = emp.id WHERE prj.id IN (%s)", projects)
				query = "SELECT DISTINCT emp.id AS id, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') AS name, emp.empRole AS role_id FROM employee emp WHERE emp.empRole = 12"
			} else if request.Filter_type == "5" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project_emps prj_emps LEFT JOIN employee emp ON prj_emps.emp_id = emp.id WHERE prj_emps.project_id IN (%s) AND prj_emps.role_id = 5", projects)
				query = fmt.Sprintf("SELECT DISTINCT emp.id AS id, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') AS name, emp.empRole AS role_id FROM project_emps prj_emps LEFT JOIN employee emp ON prj_emps.emp_id = emp.id WHERE prj_emps.project_id IN (%s) AND prj_emps.role_id = 5 LIMIT %d, %d", projects, offset, noOfRecords)
			} else if request.Filter_type == "6" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project_emps prj_emps LEFT JOIN employee emp ON prj_emps.emp_id = emp.id WHERE prj_emps.project_id IN (%s) AND prj_emps.role_id = 6", projects)
				query = fmt.Sprintf("SELECT DISTINCT emp.id AS id, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') AS name, emp.empRole AS role_id FROM project_emps prj_emps LEFT JOIN employee emp ON prj_emps.emp_id = emp.id WHERE prj_emps.project_id IN (%s) AND prj_emps.role_id = 6 LIMIT %d, %d", projects, offset, noOfRecords)
			} else if request.Filter_type == "13" {
				totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project_emps prj_emps LEFT JOIN employee emp ON prj_emps.emp_id = emp.id WHERE prj_emps.project_id IN (%s) AND prj_emps.role_id = 13", projects)
				query = "SELECT DISTINCT emp.id AS id, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') AS name, emp.empRole AS role_id FROM employee emp WHERE emp.empRole = 13"
			} else {
				jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid filter_type", "code": http.StatusBadRequest, "success": false})
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)

					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    http.StatusBadRequest,
						"message": "invalid filter_type",
						"success": false,
						"error":   err,
					})
					return
				}
				w.Write(jsonData)
			}
		} else {
			jsonData, err := json.Marshal(map[string]interface{}{"message": "no filter_type given", "code": http.StatusBadRequest, "success": false})
			if err != nil {
				// http.Error(w, err.Error(), http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)

				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "Failed to execute the query: 16400",
					"success": false,
					"error":   err,
				})

				return
			}
			w.Write(jsonData)
		}

		rows, err := DB.Query(totalPagesSQL)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Failed to execute the query: 164",
				"success": false,
				"error":   err,
			})
			return
		}
		count := 0
		type res struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Role_id string `json:"role_id"`
		}

		for rows.Next() {
			rows.Scan(&count)

		}

		rows, err = DB.Query(query)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Failed to execute the query",
				"success": false,
				"error":   err,
			})
			return

		}
		var data []res

		for rows.Next() {
			var temp res
			rows.Scan(&temp.Id, &temp.Name, &temp.Role_id)

			data = append(data, temp)
		}

		type Response struct {
			Data    []res  `json:"data"`
			Count   int    `json:"count"`
			Code    int    `json:"code"`
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		var resp Response
		resp.Data = data
		resp.Count = count
		resp.Code = http.StatusOK
		resp.Success = true
		resp.Message = "successfully"

		json.NewEncoder(w).Encode(resp)

	} else {
		jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid Role id or emp_id", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid Role id or Employee id",
				"success": false,
				"error":   err,
			})
			return
		}
		w.Write(jsonData)
	}
}
