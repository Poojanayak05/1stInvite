package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type NewResponse struct {
	Code           int       `json:"code"`
	Message        string    `json:"message"`
	Success        bool      `json:"success"`
	ID             string    `json:"id"`
	CountryID      string    `json:"countryID"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Gender         string    `json:"gender"`
	Doj            string    `json:"doj"`
	OfficeMailID   string    `json:"officeMailId"`
	PersonalMailID string    `json:"personalMailId"`
	ContactNum     string    `json:"contactNum"`
	WorkNum        string    `json:"workNum"`
	Address        string    `json:"address"`
	Address2       string    `json:"address2"`
	Address3       string    `json:"address3"`
	RoleID         string    `json:"role_id"`
	SupervisorID   string    `json:"supervisorId"`
	ProfilePic     string    `json:"profile_pic"`
	Status         string    `json:"status"`
	Pincode        string    `json:"pincode"`
	SupervisorName string    `json:"supervisorName"`
	RoleName       string    `json:"role_name"`
	LicenseNumber  string    `json:"license_number"`
	ProjectList    []Project `json:"project_list"`
}

type Project struct {
	ID          string `json:"id"`
	ProjectName string `json:"projectName"`
	ProjectID   string `json:"project_id"`
	RoleID      string `json:"role_id"`
	EmpID       string `json:"emp_id"`
	Target      string `json:"target"`
}

func GetProfileData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid input"})
		return
	}
	defer r.Body.Close()

	var request struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Failed to Unmarshal"})
		return
	}

	fields := "a.id, a.countryID, a.first_name, IFNULL(a.last_name, '') as last_name, a.gender,IFNULL(DATE_FORMAT(a.doj, '%d-%m-%Y'),'') as doj, a.officeMailId, IFNULL(a.personalMailId, '') as personalMailId, IFNULL(a.contactNum, '') as contactNum, IFNULL(a.workNum, '') as workNum, IFNULL(a.address, '') as address, IFNULL(a.address2, '') as address2, IFNULL(a.address3, '') as address3, a.empRole as role_id, IFNULL(a.supervisorId, '') as supervisorId, IFNULL(a.profile_pic, '') as profile_pic, a.status, IFNULL(a.pincode, '') as pincode, IFNULL(CONCAT(b.first_name, ' ', b.last_name), '') as supervisorName, rm.roleName as role_name, a.license_number as license_number"
	query := fmt.Sprintf("SELECT %s FROM employee a LEFT JOIN employee b ON b.id = a.supervisorId LEFT JOIN roles_Master rm ON rm.id = a.empRole WHERE a.id = %s AND a.status = 1", fields, request.ID)

	rows, err := DB.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 StatusInternalServerError", "Message": "Failed to excute query"})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusConflict, "message": "Profile Not Found", "success": false})
		return
	}

	var employee NewResponse
	if err := rows.Scan(
		&employee.ID,
		&employee.CountryID,
		&employee.FirstName,
		&employee.LastName,
		&employee.Gender,
		&employee.Doj,
		&employee.OfficeMailID,
		&employee.PersonalMailID,
		&employee.ContactNum,
		&employee.WorkNum,
		&employee.Address,
		&employee.Address2,
		&employee.Address3,
		&employee.RoleID,
		&employee.SupervisorID,
		&employee.ProfilePic,
		&employee.Status,
		&employee.Pincode,
		&employee.SupervisorName,
		&employee.RoleName,
		&employee.LicenseNumber,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := NewResponse{
			Code:           200,
			Message:        "Successfully",
			Success:        true,
			ID:             employee.ID,
			CountryID:      employee.CountryID,
			FirstName:      employee.FirstName,
			LastName:       employee.LastName,
			Gender:         employee.Gender,
			Doj:            employee.Doj,
			OfficeMailID:   employee.OfficeMailID,
			PersonalMailID: employee.PersonalMailID,
			ContactNum:     employee.ContactNum,
			WorkNum:        employee.WorkNum,
			Address:        employee.Address,
			Address2:       employee.Address2,
			Address3:       employee.Address3,
			RoleID:         employee.RoleID,
			SupervisorID:   employee.SupervisorID,
			ProfilePic:     employee.ProfilePic,
			Status:         employee.Status,
			Pincode:        employee.Pincode,
			SupervisorName: employee.SupervisorName,
			RoleName:       employee.RoleName,
			LicenseNumber:  employee.LicenseNumber,
			ProjectList:    employee.ProjectList,
		}
		fmt.Println(response)
		fmt.Println(employee.OfficeMailID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	roleID := employee.RoleID
	userID := employee.ID

	if roleID == "3" || roleID == "4" || roleID == "5" || roleID == "6" || roleID == "7" {
		var fields, queryProject string
		var projects []Project

		switch roleID {
		case "3":
			fields = "prj.id, UPPER(prj.projectName) as projectName"
			queryProject = fmt.Sprintf("SELECT %s FROM employee emp JOIN project prj ON emp.id = prj.operations_manager WHERE emp.empRole = 4 AND emp.supervisorId = %s", fields, userID)
		case "4":
			fields = "prj.id, UPPER(prj.projectName) as projectName"
			queryProject = fmt.Sprintf("SELECT %s FROM project prj WHERE prj.operations_manager = %s", fields, userID)
		case "5", "6":
			fields = "em_pr.*, UPPER(pr.projectName) as projectName"
			queryProject = fmt.Sprintf("SELECT %s FROM project_emps em_pr LEFT JOIN project pr ON em_pr.project_id = pr.id WHERE emp_id IN (%s)", fields, userID)
		case "7":
			fields = "prj.id, UPPER(prj.projectName) as projectName"
			queryProject = fmt.Sprintf("SELECT %s FROM project prj WHERE prj.driverID = %s", fields, userID)
		}

		projectRows, err := DB.Query(queryProject)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 StatusInternalServerError", "Message": "Failed to scan data"})
			return
		}
		defer projectRows.Close()

		for projectRows.Next() {
			var project Project
			var target sql.NullString
			var projectName sql.NullString

			switch roleID {
			case "3", "4", "7":
				if err := projectRows.Scan(&project.ID, &projectName); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				project.ProjectName = projectName.String
			case "5", "6":
				if err := projectRows.Scan(&project.ID, &project.EmpID, &project.ProjectID, &project.RoleID, &target, &projectName); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				project.ProjectName = projectName.String
				project.Target = target.String
			}

			projects = append(projects, project)
		}

		employee.ProjectList = projects
	}

	response := NewResponse{
		Code:           200,
		Message:        "Successfully",
		Success:        true,
		ID:             employee.ID,
		CountryID:      employee.CountryID,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		Gender:         employee.Gender,
		Doj:            employee.Doj,
		OfficeMailID:   employee.OfficeMailID,
		PersonalMailID: employee.PersonalMailID,
		ContactNum:     employee.ContactNum,
		WorkNum:        employee.WorkNum,
		Address:        employee.Address,
		Address2:       employee.Address2,
		Address3:       employee.Address3,
		RoleID:         employee.RoleID,
		SupervisorID:   employee.SupervisorID,
		ProfilePic:     employee.ProfilePic,
		Status:         employee.Status,
		Pincode:        employee.Pincode,
		SupervisorName: employee.SupervisorName,
		RoleName:       employee.RoleName,
		LicenseNumber:  employee.LicenseNumber,
		ProjectList:    employee.ProjectList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
