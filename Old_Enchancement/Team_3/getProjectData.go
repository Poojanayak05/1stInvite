package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetProjectData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type Request struct {
		Project_id string `json:"project_id"`
		Role_id    string `json:"role_id"`
		Emp_id     string `json:"emp_id"`
	}

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var roleid string
	err = DB.QueryRow("select emprole from employee where id=?", request.Emp_id).Scan(&roleid)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Failed to execute the query 30",
			"success": false,
			"error":   err,
		})
		//return
	}

	var projectExists int
	err = DB.QueryRow("select count(*) from project where id =?", request.Project_id).Scan(&projectExists)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Failed to execute the query 50",
			"success": false,
			"error":   err,
		})
		//return
	}

	if roleid == request.Role_id && projectExists != 0 {
		var target string
		if request.Role_id == "5" || request.Role_id == "6" || request.Role_id == "13" {
			target = "IFNULL(prj_emp.target,\"\") as training_target"
		} else {
			target = "IFNULL(prj.training_target,\"\") as training_target"
		}

		//fields := "prj.id as project_id, UPPER(prj.projectName) as project_name, IFNULL(prj.partnerID,\"\") as partner_id, IFNULL(par.partnerName,\"\") as partnerName, IFNULL(DATE_FORMAT(prj.startDate, \"%d-%m-%Y\"), \"\") as startDate, IFNULL(DATE_FORMAT(prj.endDate, \"%d-%m-%Y\"),\"\") as endDate, ' . $target . ', prj_emp.emp_id, prj_emp.role_id, IFNULL(prj.operations_manager,\"\") as operations_manager_id, IFNULL(CONCAT(emp_a.first_name, \" \", emp_a.last_name), \"\") as operations_manager_name,IFNULL(prj.gfl_id,\"\") as gfl_id, IFNULL(CONCAT(emp_ans.first_name, \" \", emp_ans.last_name), \"\") as gfl_name, IFNULL(prj.driverId,\"\") as driverId, IFNULL(CONCAT(emp_b.first_name, \" \", emp_b.last_name), \"\") as driver_name, CONCAT(emp_c.first_name, \" \", IFNULL(emp_c.last_name, \"\")) as emp_name, loc.id as location_id, IFNULL(prj.funderID,\"\") as funder_id, IFNULL(prj.funderID,\"\") as funder_id, funder.funderName as funder_name, IFNULL(prj.busID,\"\") as bus_id, IFNULL(bus.register_number,\"\") as bus_number, CASE WHEN prj.endDate >= CURRENT_DATE() THEN \"Published\" ELSE \"Completed\" END AS project_status, district.name as location_name"

		type Project struct {
			ProjectID           int    `json:"project_id"`
			ProjectName         string `json:"project_name"`
			PartnerID           string `json:"partner_id"`
			PartnerName         string `json:"partnerName"`
			StartDate           string `json:"start_date"`
			EndDate             string `json:"end_date"`
			TrainingTarget      string `json:"training_target"`
			EmpID               int    `json:"emp_id"`
			RoleID              int    `json:"role_id"`
			OperationsManagerID string `json:"operations_manager_id"`
			OperationsManager   string `json:"operations_manager_name"`
			GflID               string `json:"gfl_id"`
			GflName             string `json:"gfl_name"`
			DriverID            string `json:"driver_id"`
			DriverName          string `json:"driver_name"`
			EmpName             string `json:"emp_name"`
			LocationID          string `json:"location_id"`
			FunderID            string `json:"funder_id"`
			FunderName          string `json:"funder_name"`
			BusID               string `json:"bus_id"`
			BusNumber           string `json:"bus_number"`
			ProjectStatus       string `json:"project_status"`
			LocationName        string `json:"location_name"`
		}

		var data []Project

		fields := `
    prj.id as project_id,
    UPPER(prj.projectName) as project_name,
    IFNULL(prj.partnerID, "") as partner_id,
    IFNULL(par.partnerName, "") as partnerName,
    IFNULL(DATE_FORMAT(prj.startDate, "%d-%m-%Y"), "") as startDate,
    IFNULL(DATE_FORMAT(prj.endDate, "%d-%m-%Y"), "") as endDate,
    ` + target + `,
    prj_emp.emp_id,
    prj_emp.role_id,
    IFNULL(prj.operations_manager, "") as operations_manager_id,
    IFNULL(CONCAT(emp_a.first_name, " ", emp_a.last_name), "") as operations_manager_name,
    IFNULL(prj.gfl_id, "") as gfl_id,
    IFNULL(CONCAT(emp_ans.first_name, " ", emp_ans.last_name), "") as gfl_name,
    IFNULL(prj.driverId, "") as driverId,
    IFNULL(CONCAT(emp_b.first_name, " ", emp_b.last_name), "") as driver_name,
    CONCAT(emp_c.first_name, " ", IFNULL(emp_c.last_name, "")) as emp_name,
    loc.id as location_id,
    IFNULL(prj.funderID, "") as funder_id,
    IFNULL(funder.funderName,'') as funder_name,
    IFNULL(prj.busID, "") as bus_id,
    IFNULL(bus.register_number, "") as bus_number,
	IFNULL(CASE
        WHEN prj.project_status = "0" THEN "In Progress"
        WHEN prj.project_status = "1" THEN "Published"
        ELSE "Completed"
    END,"") AS project_status,
    district.name as location_name
`

		query := fmt.Sprintf(`
    SELECT %s
    FROM project prj
    LEFT JOIN partner par ON prj.partnerID = par.partnerID
    LEFT JOIN project_emps prj_emp ON prj.id = prj_emp.project_id
    LEFT JOIN employee emp_a ON prj.operations_manager = emp_a.id
    LEFT JOIN employee emp_b ON prj.driverId = emp_b.id
    LEFT JOIN employee emp_c ON prj_emp.emp_id = emp_c.id
    LEFT JOIN employee emp_ans ON prj.gfl_id = emp_ans.id
    LEFT JOIN location loc ON prj.locationID = loc.id
    LEFT JOIN location district ON loc.parentId = district.id
    LEFT JOIN funder ON prj.funderID = funder.funderID
    LEFT JOIN bus ON prj.busID = bus.id
    WHERE prj.id = %s
    GROUP BY prj_emp.emp_id
    ORDER BY projectName
`, fields, request.Project_id)
		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Failed to execute the query 79",
				"success": false,
				"error":   err,
			})
			return
		}

		defer rows.Close()

		for rows.Next() {
			var project Project
			var empid, roleid sql.NullInt64
			var empname sql.NullString
			var lactionname string
			err = rows.Scan(
				&project.ProjectID,
				&project.ProjectName,
				&project.PartnerID,
				&project.PartnerName,
				&project.StartDate,
				&project.EndDate,
				&project.TrainingTarget,
				//&project.EmpID,
				&empid,
				//&project.RoleID,
				&roleid,
				&project.OperationsManagerID,
				&project.OperationsManager,
				&project.GflID,
				&project.GflName,
				&project.DriverID,
				&project.DriverName,
				//&project.EmpName,
				&empname,
				&project.LocationID,
				&project.FunderID,
				&project.FunderName,
				&project.BusID,
				&project.BusNumber,
				&project.ProjectStatus,
				&lactionname,
			)

			if err != nil {
				log.Println("Error scanning row:", err)
				// handle error
			}

			project.LocationName = lactionname
			if empid.Valid {
				project.EmpID = int(empid.Int64)
			} else {
				project.EmpID = 0
			}

			if roleid.Valid {
				project.RoleID = int(roleid.Int64)
			} else {
				project.RoleID = 0
			}

			if empname.Valid {
				project.EmpName = empname.String
			} else {
				project.EmpName = ""
			}

			data = append(data, project)
		}

		type Trainer struct {
			EmpID int    `json:"emp_id"`
			Name  string `json:"name"`
		}

		type GelathiFacilitator struct {
			EmpID int    `json:"emp_id"`
			Name  string `json:"name"`
		}

		type Employee struct {
			EmpID   string `json:"emp_id"`
			EmpName string `json:"emp_name"`
			Role    string `json:"role"`
		}

		type View struct {
			Trainers                []Trainer            `json:"trainers"`
			GelathiFacilitator      []GelathiFacilitator `json:"gelathiFacilitator"`
			ProjectPeoplesList      []Employee           `json:"projectPeoplesList"`
			TrainersCount           int                  `json:"trainers_count"`
			GelathiFacilitatorCount int                  `json:"gelathiFacilitator_count"`
			ProjectID               string               `json:"project_id"`
			ProjectName             string               `json:"project_name"`
			PartnerID               string               `json:"partner_id"`
			PartnerName             string               `json:"partnerName"`
			BusID                   string               `json:"bus_id"`
			BusNumber               string               `json:"bus_number"`
			EndDate                 string               `json:"endDate"`
			StartDate               string               `json:"startDate"`
			FunderID                string               `json:"funder_id"`
			FunderName              string               `json:"funder_name"`
			ProjectStatus           string               `json:"project_status"`
			TrainingTarget          string               `json:"training_target"`
			OperationsManagerID     string               `json:"operations_manager_id"`
			OperationsManagerName   string               `json:"operations_manager_name"`
			GflID                   string               `json:"gfl_id"`
			GflName                 string               `json:"gfl_name"`
			DriverID                string               `json:"driverId"`
			DriverName              string               `json:"driver_name"`
			LocationName            string               `json:"location_name"`
			LocationID              string               `json:"location_id"`
			IsAssociated            bool                 `json:"isAssociated"`
			Flag                    string               `json:"flag"`
			FinalSave               string               `json:"final_save"`
		}

		view := View{}
		view.Trainers = make([]Trainer, 0)
		view.GelathiFacilitator = make([]GelathiFacilitator, 0)
		peopleIds := make([]string, 0)
		view.ProjectPeoplesList = make([]Employee, 0)
		view.TrainersCount = 0
		view.GelathiFacilitatorCount = 0

		if len(data) > 0 {
			for _, val := range data {
				view.ProjectID = strconv.Itoa(val.ProjectID)
				projID := view.ProjectID
				view.ProjectName = val.ProjectName
				view.PartnerID = val.PartnerID
				view.PartnerName = val.PartnerName
				view.StartDate = val.StartDate
				view.EndDate = val.EndDate
				view.TrainingTarget = val.TrainingTarget
				view.OperationsManagerID = val.OperationsManagerID
				view.OperationsManagerName = val.OperationsManager
				view.GflID = val.GflID
				view.GflName = val.GflName
				view.DriverID = val.DriverID
				view.DriverName = val.DriverName
				view.LocationID = val.LocationID
				view.FunderID = val.FunderID
				view.FunderName = val.FunderName
				view.BusID = (val.BusID)
				view.BusNumber = val.BusNumber
				view.ProjectStatus = val.ProjectStatus
				view.LocationName = val.LocationName

				checkForAssociation := fmt.Sprintf("SELECT * FROM project_association pa WHERE projectId = %s", projID)

				rows, err = DB.Query(checkForAssociation)
				if err != nil {
					return
				}
				c := 0
				for rows.Next() {
					c++
				}
				if c > 0 {
					view.IsAssociated = true
				} else {
					view.IsAssociated = false
				}

				if val.RoleID == 5 {

					trainer := Trainer{
						EmpID: val.EmpID,
						Name:  val.EmpName,
					}
					view.Trainers = append(view.Trainers, trainer)
					view.TrainersCount++
				}
				if val.RoleID == 6 {

					gelathiFacilitator := GelathiFacilitator{
						EmpID: val.EmpID,
						Name:  val.EmpName,
					}
					view.GelathiFacilitator = append(view.GelathiFacilitator, gelathiFacilitator)
					view.GelathiFacilitatorCount++
				}

				if val.OperationsManagerID != "" {
					peopleIds = append(peopleIds, val.OperationsManagerID)
				}

				if val.GflID != "" {
					peopleIds = append(peopleIds, val.GflID)
				}

				if val.DriverID != "" {
					peopleIds = append(peopleIds, val.DriverID)
				}

				peopleIds = append(peopleIds, strconv.Itoa(val.EmpID))

				//peopleIds = append(peopleIds, val.OperationsManagerID, val.GflID, strconv.Itoa(val.EmpID), val.DriverID)

			}
		}

		employees := strings.Join(peopleIds, ",")

		fields = "emp.id as emp_id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) as emp_name, rm.roleName as role"

		query = fmt.Sprintf("SELECT %s FROM employee emp LEFT JOIN roles_Master rm ON emp.empRole = rm.id WHERE emp.id IN (%s) ORDER BY emp.empRole", fields, employees)

		rows, err = DB.Query(query)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Failed to execute the query:244",
				"success": false,
				"error":   err,
			})
			return
		}
		for rows.Next() {
			var employee Employee
			rows.Scan(&employee.EmpID, &employee.EmpName, &employee.Role)
			view.ProjectPeoplesList = append(view.ProjectPeoplesList, employee)

		}

		busFields := "IFNULL((SELECT IF(check_count > 0, 1, 0) FROM bus_checklist WHERE date_checked = CURDATE() AND bus_id = bus.id), 0) as flag, IFNULL((SELECT final_save FROM bus_checklist WHERE date_checked = CURDATE() AND bus_id = bus.id), 0) as final_save"

		busID, _ := strconv.Atoi(view.BusID)
		queryBus := fmt.Sprintf("SELECT %s FROM bus bus LEFT JOIN project prj ON bus.id = prj.busID WHERE bus.id = %d", busFields, busID)

		rows, err = DB.Query(queryBus)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Failed to execute the query : 261",
				"success": false,
				"error":   err,
			})
			return
		}
		for rows.Next() {
			rows.Scan(&view.Flag, &view.FinalSave)

		}

		if request.Role_id == "5" {
			queryTarget := fmt.Sprintf("SELECT target FROM project_emps WHERE project_id = %s AND emp_id = %s", request.Project_id, request.Emp_id)

			rows, err = DB.Query(queryTarget)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)

				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "Failed to execute the query:274",
					"success": false,
					"error":   err,
				})
				return
			}
			for rows.Next() {
				rows.Scan(&view.TrainingTarget)

			}
		}

		type Response struct {
			List    View   `json:"list"`
			Code    int    `json:"code"`
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		response := Response{
			List:    view,
			Code:    200,
			Success: true,
			Message: "successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	} else {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(map[string]interface{}{"message": "Invalid data given by user", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid data given by user",
				"success": false,
				"error":   err,
			})
			return
		}
		w.Write(jsonData)
	}
}
