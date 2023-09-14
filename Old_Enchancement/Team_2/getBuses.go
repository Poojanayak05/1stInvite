package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//-------------------------------- Get Buses working but Not getting  ------------------------------------
type Bus struct {
	Date        string `json:"date"`
	Role_id     int    `json:"Role_id"`
	Project_id  string `json:"project_id"`
	Taluk_id    string `json:"taluk_id"`
	District_id string `json:"district_id"`
	Funder_id   string `json:"funder_id"`
	Emp_id      int    `json:"emp_id"`
	PageNum     int    `json:"pageNum"`
	Search      string `json:"search"`
}

type busdata struct {
	Bus_id          string `json:"bus_id"`
	Register_number string `json:"register_number"`
	Project_id      string `json:"project_id"`
	Project_name    string `json:"project_name"`
	Flag            string `json:"flag"`
	Final_save      string `json:"final_save"`
	Checked_count   string `json:"checked_count"`
}

func removeDuplicates(elements []interface{}) []interface{} {
	encountered := map[interface{}]bool{}
	result := []interface{}{}

	for _, element := range elements {
		if !encountered[element] {
			encountered[element] = true
			result = append(result, element)
		}
	}

	return result
}

func GetBuses(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	//func Login(w http.ResponseWriter, r *http.Request) {
	//SetupCORS(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var checkCount int
	var totalRows int
	checkCount = 0

	var p Bus
	var data []busdata
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})

		return
	}

	//----------- serach paramseter -------------
	searchFilter := ""
	searchParam := p.Search
	if len(searchParam) > 0 {
		searchFilter = fmt.Sprintf(" and register_number like '%s'", searchParam)
	}

	// roleid, err := strconv.Atoi(p.Role_id)
	// empid, err := strconv.Atoi(p.Emp_id)

	var query string
	//------------
	if p.Role_id > 0 && p.Emp_id > 0 {

		var date string
		var other_roles string
		var condition string
		var join string
		var join_add string

		if date != "" {
			date = p.Date
		}

		//------------- if project id is there ---------------
		if p.Project_id != "" {
			join = " LEFT JOIN project prj ON bus.id = prj.busID "
			condition = fmt.Sprintf(" WHERE prj.id = %s", p.Project_id)
			other_roles = fmt.Sprintf(" AND prj.id = %s", p.Project_id)
		}

		//-------------------funder id ------------------

		if p.Funder_id != "" {
			join = " LEFT JOIN project prj ON bus.id = prj.busID"
			condition = fmt.Sprintf(" WHERE prj.funderID = %s", p.Funder_id)
			other_roles = fmt.Sprintf(" AND prj.funderID = %s", p.Funder_id)
		}

		//----------------- district id ------------------

		if p.District_id != "" && p.Taluk_id == "" {
			join = " LEFT JOIN project prj ON bus.id = prj.busID"
			join_add = " LEFT JOIN location loc ON prj.locationID = loc.id"
			condition = fmt.Sprintf(" WHERE loc.parentId = %s", p.District_id)
			join_add += " LEFT JOIN location district ON loc.parentId = district.id"
			other_roles = fmt.Sprintf(" AND loc.parentId = %s", p.District_id)
		}

		//----------------- taluk Id -------------------
		if p.Taluk_id != "" {
			join_add = ""
			join = ""
			join = " LEFT JOIN project prj ON bus.id = prj.busID"
			condition = fmt.Sprintf(" WHERE prj.locationID = %s", p.Taluk_id)
			other_roles = fmt.Sprintf(" AND prj.locationID = %s", p.Taluk_id)
		}

		//---------- page num -----------------
		var pageno int
		if p.PageNum == 0 {
			pageno = 1
		} else {
			pageno = p.PageNum
		}
		// totalRows1, err := strconv.Atoi(totalRows)
		// if err != nil {
		// }
		//var query string
		no_of_records := 25
		offset := (pageno - 1) * no_of_records

		//var total_pages float64
		if p.Role_id == 1 || p.Role_id == 2 {

			if len(condition) == 0 && len(searchFilter) > 0 {
				searchFilter = fmt.Sprintf(" where register_number like '%s'", searchParam)
			}

			total_pages_sql := fmt.Sprintf(" SELECT COUNT(*) FROM bus" + join + join_add + condition + searchFilter)

			result, err := DB.Query(total_pages_sql)
			if err != nil {

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}

			if result.Next() {
				err := result.Scan(&totalRows)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

					return
				}
			}

			// total_pages := int(math.Ceil(float64(totalRows) / float64(no_of_records)))

			fields := "bus.id as bus_id, UPPER(bus.register_number) as register_number, COALESCE(prj.id,'') as project_id, IFNULL(UPPER(prj.projectName), '') as project_name, IFNULL((SELECT IF(check_count > 0, 1, 0) From bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as flag, IFNULL((SELECT final_save From bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as final_save, IFNULL((SELECT check_count From bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as checked_count"

			query = "SELECT " + fields + " FROM bus bus LEFT JOIN project prj ON bus.id = prj.busID " + join_add + condition + searchFilter + " ORDER BY bus.id desc LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(no_of_records)

		}
		projects := []map[string]interface{}{}
		var queryProject string
		//var totalRows int
		if p.Role_id == 3 || p.Role_id == 4 || p.Role_id == 5 || p.Role_id == 12 || p.Role_id == 13 {
			if p.Role_id == 12 {

				sEmpidRows, err := DB.Query("SELECT id FROM employee e WHERE supervisorId = ?", p.Emp_id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				defer sEmpidRows.Close()

				for sEmpidRows.Next() {
					var rowID int
					err = sEmpidRows.Scan(&rowID)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

						return
					}

					fields := "prj.id"
					queryProject := "SELECT " + fields + " FROM employee emp JOIN project prj ON emp.id = prj.operations_manager " + join_add + " WHERE emp.empRole = 4 AND emp.id = ?"
					projectResult, err := DB.Query(queryProject, rowID)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

						return
					}
					defer projectResult.Close()

					projectsRow := []map[string]interface{}{}
					for projectResult.Next() {
						projectData := make(map[string]interface{})
						var projectID int
						err = projectResult.Scan(&projectID)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

							return
						}
						projectData["id"] = projectID
						projectsRow = append(projectsRow, projectData)
					}
					projects = append(projects, projectsRow...)
				}
			}

			//---------------------------- role id ==3 --------------------------
			if p.Role_id == 3 {
				fields := "prj.id"
				queryProject = fmt.Sprintf("SELECT %s FROM employee emp JOIN project prj ON emp.id = prj.operations_manager %s WHERE emp.empRole = 4 AND emp.supervisorId = %d %s", fields, join_add, p.Emp_id, other_roles)

			}

			//------------------ role id ==  13 and  ------------------------
			if p.Role_id == 4 || p.Role_id == 13 {
				fields := "COALESCE(prj.id,0)"
				queryProject = fmt.Sprintf("SELECT %s FROM project prj %s WHERE prj.operations_manager = %d %s", fields, join_add, p.Emp_id, other_roles)

			}
			if p.Role_id == 13 {
				fields := "COALESCE(prj.id,0)"
				queryProject = fmt.Sprintf("SELECT %s FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id %s WHERE emp_id IN (%d) %s", fields, join_add, p.Emp_id, other_roles)
			}
			if p.Role_id == 5 {
				fields := "COALESCE(prj.id,0)"
				queryProject = fmt.Sprintf("SELECT %s FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id %s WHERE emp_id IN (%d) %s", fields, join_add, p.Emp_id, other_roles)
			}
			//======================= login is not there ===========================
			if p.Role_id == 7 || p.Role_id == 6 {
				fields := "COALESCE(prj.id,0)"
				queryProject = fmt.Sprintf("SELECT %s FROM project prj LEFT JOIN employee emp ON prj.driverID = emp.id WHERE prj.driverID = %d", fields, p.Emp_id)
			}

			projs := []interface{}{} // declaring as slice
			//----------------------- role id != 12 ---------------------------

			if p.Role_id != 12 {
				projectResult, err := DB.Query(queryProject)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

					return
				}
				defer projectResult.Close()

				projects := []map[string]interface{}{}
				for projectResult.Next() {
					projectData := make(map[string]interface{})
					var projectID int
					err = projectResult.Scan(&projectID)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

						return
					}
					projectData["id"] = projectID
					projects = append(projects, projectData)
				}

				for _, vv := range projects {

					if vv != nil {
						if projectID, ok := vv["id"].(int); ok {
							projs = append(projs, projectID)
						}
					}
				}
			}

			//--------------------- role id ==  12 ---------------------------------
			if p.Role_id == 12 {
				for _, vv := range projects {

					if vv != nil {
						if projectID, ok := vv["id"].(int); ok {
							projs = append(projs, projectID)
						}
					}

				}

			}
			//--------------------------reomving duplicate projet ids ----------------
			projs = removeDuplicates(projs)
			// ============== converting interface{} to string =================

			projStrings := make([]string, len(projs))
			for i, val := range projs {
				projStrings[i] = fmt.Sprint(val)
			}
			projIDs := strings.TrimRight(strings.Join(projStrings, ","), ",")
			if projIDs == "" {
				projIDs = "''"
			}

			//------------------- counting total pages -----------------
			totalPagesSQL := fmt.Sprintf("SELECT COUNT(*) FROM bus bus LEFT JOIN project prj ON bus.id = prj.busID WHERE prj.id IN (%s) %s ", projIDs, searchFilter)

			err := DB.QueryRow(totalPagesSQL).Scan(&totalRows)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}

			// total_pages := int(math.Ceil(float64(totalRows) / float64(no_of_records)))

			//------------- queries  ------------------------------
			fields := "bus.id as bus_id, UPPER(bus.register_number) as register_number, COALESCE(prj.id,'') as project_id, IFNULL(UPPER(prj.projectName), '') as project_name, IFNULL((SELECT IF(check_count > 0, 1, 0) FROM bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as flag, IFNULL((SELECT final_save FROM bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as final_save, IFNULL((SELECT check_count FROM bus_checklist WHERE date_checked = '" + date + "' AND bus_id = bus.id), 0) as checked_count"
			query = fmt.Sprintf("SELECT %s FROM bus bus LEFT JOIN project prj ON bus.id = prj.busID WHERE prj.id IN (%s) %s ORDER BY bus.id DESC LIMIT %d, %d", fields, projIDs, searchFilter, offset, no_of_records)
		}

		// ------------------------------ end program -------------------------

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"list": data, "checked_count": checkCount, "total_count": totalRows, "code": 200, "success": true, "message": "Successfully"})
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			// json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

			return
		}
		defer rows.Close()

		//checkCount := 0
		for rows.Next() {

			var bus_id, register_number, project_id, project_name, flag, final_save, checked_count string
			err := rows.Scan(&bus_id, &register_number, &project_id, &project_name, &flag, &final_save, &checked_count)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}

			if flag == "1" {
				checkCount++
			}
			data = append(data, busdata{Bus_id: bus_id, Register_number: register_number, Project_id: project_id, Project_name: project_name, Flag: flag, Final_save: final_save, Checked_count: checked_count})

		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"list": data, "checked_count": checkCount, "total_count": totalRows, "code": 200, "success": true, "message": "Successfully"})

	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"list": data, "checked_count": checkCount, "total_count": totalRows, "code": 200, "success": true, "message": "Successfully"})
	}
}
