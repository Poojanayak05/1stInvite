package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetProjects(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type Request struct {
		EndDate    string `json:"end_date"`
		StartDate  string `json:"start_date"`
		Search     string `json:"search"`
		ID         string `json:"id"`
		RoleID     string `json:"role_id"`
		FilterID   int    `json:"filter_id"`
		Type       string `json:"type"`
		PageNum    int    `json:"pageNum"`
		Count      string `json:"count"`
		TalukID    string `json:"taluk_id"`
		DistrictID string `json:"district_id"`
		FunderID   string `json:"funder_id"`
		OpsManager string `json:"operations_manager_id"`
		TrainerID  string `json:"trainer_id"`
		GelathiID  string `json:"gelathi_id"`
		PartnerID  string `json:"partner_id"`
	}

	type Project struct {
		TrainerType       string `json:"trainer_type"`
		EmpID             string `json:"emp_id"`
		ID                string `json:"id"`
		ProjectID         string `json:"project_id"`
		Name              string `json:"name"`
		StartDate         string `json:"start_date"`
		EndDate           string `json:"end_date"`
		ProjectStatusName string `json:"project_status_name"`
		ProjectStatus     string `json:"project_status"`
		LocationName      string `json:"location_name"`
	}

	type Response struct {
		List    []Project `json:"list"`
		Count   string    `json:"count"`
		Code    int       `json:"code"`
		Success bool      `json:"success"`
		Message string    `json:"message"`
	}

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}

	var searchFilter string
	if request.Search != "" {
		searchFilter = fmt.Sprintf(" and projectName like '%%%s%%'", request.Search)

	}

	search := ""

	if request.StartDate != "" && request.EndDate != "" {

		search += fmt.Sprintf("AND prj.startDate >= '%s' AND prj.endDate <= '%s'", request.StartDate, request.EndDate)
	}

	if request.FunderID != "" {

		search += fmt.Sprintf(" AND prj.funderID = '%s'", request.FunderID)
	}

	if request.OpsManager != "" {

		search += fmt.Sprintf(" AND prj.operations_manager = '%s'", request.OpsManager)
	}

	join := ""
	where := ""

	if request.DistrictID != "" {

		where += fmt.Sprintf("AND loc.parentId = '%s'", request.DistrictID)
	}

	if request.TalukID != "" {

		where += fmt.Sprintf("AND loc.id = '%s'", request.TalukID)
	}

	// If searched by district/taluk ends here

	// If searched by trainer/gelathi starts here

	if request.TrainerID != "" {

		join += "LEFT JOIN project_emps pr_employee ON prj.id = pr_employee.project_id"

		where += fmt.Sprintf(" AND pr_employee.emp_id = '%s' AND pr_employee.role_id = 5", request.TrainerID)
	}

	if request.GelathiID != "" {

		join += " LEFT JOIN project_emps pr_empp ON prj.id = pr_empp.project_id"

		where += fmt.Sprintf(" AND pr_empp.emp_id = '%s' AND pr_empp.role_id = 6", request.GelathiID)
	}

	var roleid string
	var validId bool

	if request.ID != "" {
		query := fmt.Sprintf("SELECT empRole FROM employee where status = 1 and id = %s", request.ID)
		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		for rows.Next() {
			rows.Scan(&roleid)
			validId = true
		}

	}
	if !validId {

		jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid Employee id", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(jsonData)
		return
	}

	if roleid == "11" {
		roleid = "1"
	}

	var totalPagesSQL string
	var fields string
	var query string
	var proj []Project
	var projectArray []string
	existingIntegers := make(map[string]bool)
	var totalRows string

	if roleid == "2" || roleid == "1" { // Admin = 2 / CEO
		if roleid == "1" {
			totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id WHERE prj.status = 1 AND prj.project_status = '1' %s %s %s", join, search, where, searchFilter)

			fields = "prj.id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

			query = fmt.Sprintf("SELECT %s FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE prj.status = 1 AND prj.project_status = '1' %s %s %s ORDER BY prj.endDate", fields, join, search, where, searchFilter)

			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			for rows.Next() {
				var ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
				rows.Scan(&ID, &ProjectID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

				if !existingIntegers[ID] {
					proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

					existingIntegers[ID] = true
				}
			}

		} else {
			totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id WHERE prj.status = 1 %s %s %s", join, search, where, searchFilter)

			fields = "prj.id, prj.id AS project_id, UPPER(prj.projectName) AS name, CASE WHEN prj.project_status = '0' THEN 'In Progress' WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

			query = fmt.Sprintf("SELECT %s FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE prj.status = 1 %s %s %s ORDER BY prj.endDate", fields, join, search, where, searchFilter)

			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			for rows.Next() {
				var ID, ProjectID, Name, ProjectStatusName, ProjectStatus, LocationName string
				rows.Scan(&ID, &ProjectID, &Name, &ProjectStatusName, &ProjectStatus, &LocationName)

				if !existingIntegers[ID] {
					proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

					existingIntegers[ID] = true
				}
			}

		}
	} else if roleid == "3" { // Program Manager = 3
		totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM employee emp %s LEFT JOIN project prj ON prj.operations_manager = emp.id WHERE emp.supervisorId = %s AND emp.status = 1 AND prj.status = 1 AND prj.project_status = '1' %s %s %s", join, request.ID, search, where, searchFilter)

		fields = "emp.id AS emp_id, prj.id AS id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

		query = fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN project prj ON prj.operations_manager = emp.id %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE emp.supervisorId = %s AND emp.status = 1 AND prj.status = 1 AND prj.project_status = '1' %s %s %s ORDER BY prj.endDate DESC", fields, join, request.ID, search, where, searchFilter)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		for rows.Next() {
			var EmpID, ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
			rows.Scan(&EmpID, &ID, &ProjectID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

			//proj = append(proj, Project{EmpID: EmpID, ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})
			if !existingIntegers[ID] {
				proj = append(proj, Project{EmpID: EmpID, ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				existingIntegers[ID] = true
			}
		}

	} else if roleid == "12" { // SOM
		totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM employee emp %s LEFT JOIN project prj ON prj.operations_manager = emp.id WHERE emp.supervisorId = %s AND emp.status = 1 AND prj.status = 1 AND prj.project_status = '1' %s %s %s", join, request.ID, search, where, searchFilter)

		fields = "emp.id AS emp_id, prj.id AS id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

		query = fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN project prj ON prj.operations_manager = emp.id %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE emp.supervisorId = %s AND emp.status = 1 AND prj.status = 1 AND prj.project_status = '1' %s %s %s ORDER BY prj.endDate DESC", fields, join, request.ID, search, where, searchFilter)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		for rows.Next() {
			var EmpID, ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
			rows.Scan(&EmpID, &ID, &ProjectID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

			if !existingIntegers[ID] {
				proj = append(proj, Project{EmpID: EmpID, ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				existingIntegers[ID] = true
			}
		}

	} else if roleid == "4" { // Operation Manager = 4
		// totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj %s WHERE prj.operations_manager = %s AND prj.status = '1' AND prj.project_status = '1' %s %s %s", join, request.ID, search, where, searchFilter)

		// fields = "prj.id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

		totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id WHERE prj.operations_manager = %s AND prj.status = '1' AND prj.project_status = '1' %s %s %s", join, request.ID, search, where, searchFilter)

		fields = "prj.id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"
		query = fmt.Sprintf("SELECT %s FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE prj.operations_manager = %s AND prj.status = '1' AND prj.project_status = '1' %s %s %s ORDER BY prj.endDate DESC", fields, join, request.ID, search, where, searchFilter)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		for rows.Next() {
			var ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
			rows.Scan(&ID, &ProjectID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

			if !existingIntegers[ID] {
				proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				existingIntegers[ID] = true
			}
		}

	} else if roleid == "13" { // Operation Manager = 4
		totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project prj %s WHERE prj.gfl_id = %s AND prj.status = '1' AND prj.project_status = '1' %s %s %s", join, request.ID, search, where, searchFilter)

		fields = "prj.id, prj.id AS project_id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

		query = fmt.Sprintf("SELECT %s FROM project prj %s LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE prj.gfl_id = %s AND prj.status = '1' AND prj.project_status = '1' %s %s %s ORDER BY prj.endDate DESC", fields, join, request.ID, search, where, searchFilter)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		for rows.Next() {
			var ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
			rows.Scan(&ID, &ProjectID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

			//proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

			if !existingIntegers[ID] {
				proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				existingIntegers[ID] = true
			}
		}

	} else if roleid == "5" || roleid == "6" || roleid == "7" { // Trainer = 5, GelathiFacilitator = 6, Driver = 7
		fields = ""
		if roleid == "5" {
			query = fmt.Sprintf("SELECT id FROM employee emp WHERE emp.status = 1 AND emp.supervisorId = %s", request.ID)

			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			count := 0
			for rows.Next() {
				count++
				//fields += "'senior' AS trainer_type,"
			}

			if count > 0 {
				fields += "'senior' AS trainer_type,"
			} else {
				fields += "'' AS trainer_type,"
			}

		}

		totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM project_emps pr_emp %s LEFT JOIN project prj ON pr_emp.project_id = prj.id WHERE pr_emp.emp_id = %s %s %s %s", join, request.ID, search, where, searchFilter)

		fields += "pr_emp.project_id, prj.id, UPPER(prj.projectName) AS name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name AS location_name"

		query = fmt.Sprintf("SELECT %s FROM project_emps pr_emp %s LEFT JOIN project prj ON pr_emp.project_id = prj.id LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE pr_emp.emp_id = %s %s %s %s ORDER BY prj.endDate DESC", fields, join, request.ID, search, where, searchFilter)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		for rows.Next() {
			if roleid == "5" {
				var TrainerType, ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
				rows.Scan(&TrainerType, &ProjectID, &ID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

				//proj = append(proj, Project{TrainerType: TrainerType, ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				if !existingIntegers[ID] {

					proj = append(proj, Project{TrainerType: TrainerType, ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})
					existingIntegers[ID] = true
				}

			} else {
				var ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
				rows.Scan(&ProjectID, &ID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

				//proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				if roleid == "6" {
					if !existingIntegers[ID] && ID != "" {
						projectArray = append(projectArray, ID)
						proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

						existingIntegers[ID] = true
					}
				}
				if !existingIntegers[ID] {
					proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

					existingIntegers[ID] = true
				}
			}

		}

	}

	rows, err := DB.Query(totalPagesSQL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	for rows.Next() {
		rows.Scan((&totalRows))
	}

	if roleid == "6" {
		if len(projectArray) > 0 {
			query = fmt.Sprintf("SELECT prj.id as project_id, prj.id, UPPER(prj.projectName) as name, prj.startDate, prj.endDate, CASE WHEN prj.endDate >= CURRENT_DATE() THEN 'Published' ELSE 'Completed' END AS project_status_name, project_status, district.name as location_name FROM project_association pa INNER JOIN project prj ON prj.id = pa.projectId LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN location district ON loc.parentId = district.id WHERE pa.associatedProject IN (%s) AND pa.projectId NOT IN (%s)",
				strings.Join(projectArray, ","),
				strings.Join(projectArray, ","),
			)

			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			for rows.Next() {
				var ID, ProjectID, Name, StartDate, EndDate, ProjectStatusName, ProjectStatus, LocationName string
				rows.Scan(&ProjectID, &ID, &Name, &StartDate, &EndDate, &ProjectStatusName, &ProjectStatus, &LocationName)

				//proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

				if !existingIntegers[ID] {
					proj = append(proj, Project{ID: ID, ProjectID: ProjectID, Name: Name, StartDate: StartDate, EndDate: EndDate, ProjectStatusName: ProjectStatusName, ProjectStatus: ProjectStatus, LocationName: LocationName})

					existingIntegers[ID] = true
				}
			}
		}
		totalRows = strconv.Itoa(len(proj))
	}

	response := Response{
		List:    proj,
		Count:   totalRows,
		Code:    http.StatusOK,
		Success: true,
		Message: "successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
