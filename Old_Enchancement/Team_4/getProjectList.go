package Team_4

//Done by Keerthana
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Employee1 struct {
	ID           int `json:"id"`
	EmpRole      int `json:"empRole"`
	SupervisorID int `json:"supervisorId"`
}

type Project1 struct {
	ID                 string `json:"id"`
	CountryID          string `json:"countryID"`
	ProjectName        string `json:"projectName"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	FunderID           string `json:"funderID"`
	PartnerID          string `json:"partnerID"`
	Training_target    string `json:"training_target"`
	LocationID         string `json:"locationID"`
	Operations_manager string `json:"operations_manager"`
	BusID              string `json:"busID"`
	DriverID           string `json:"driverID"`
	Project_status     string `json:"project_status"`
	Status             string `json:"status"`
	CreatedAt          string `json:"createdAt"`
	CreatedBy          string `json:"createdBy"`
	LastUpdatedAt      string `json:"lastUpdatedAt"`
	LastUpdatedBy      string `json:"lastUpdatedBy"`
	Gfl_id             string `json:"gfl_id"`
}
type Project2 struct {
	ID          string `json:"id"`
	ProjectName string `json:"projectName"`
}

type Response struct {
	List    []Project1 `json:"list"`
	Code    int        `json:"code"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}
type Response11 struct {
	List    []Project2 `json:"list"`
	Code    int        `json:"code"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

func GetProjectList(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method Not found", "success": false})
		return
	}
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "error": err})
		return
	}

	var request map[string]interface{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "error": err})
		return
	}

	managerID, ok := request["manager_id"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid manager_id in request", "success": false, "error": err})
		return
	}
	var managerIDInt int
	if managerID != "" {
		managerIDInt, err = strconv.Atoi(managerID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Unable to convert string to int", "success": false, "error": err})
			return
		}
	} else if managerID == "" {
		managerIDInt = 0
	}
	var employee Employee1
	query := fmt.Sprintf("SELECT id,empRole ,supervisorId  FROM employee WHERE id = %d", managerIDInt)
	err = db.QueryRow(query).Scan(&employee.ID, &employee.EmpRole, &employee.SupervisorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to fetch employee data", "success": false, "error": err})
		return
	}

	fields := "id, UPPER(projectName) as projectName"
	if employee.EmpRole == 13 {
		query = fmt.Sprintf("SELECT %s FROM project where operations_manager =%d ORDER BY projectName", fields, employee.SupervisorID)
		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false, "error": err})
			return
		}
		defer rows.Close()
		projects := []Project2{}
		for rows.Next() {
			var project Project2
			err := rows.Scan(&project.ID, &project.ProjectName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false, "error": err})
				return
			}
			projects = append(projects, project)
		}

		response := Response11{projects, http.StatusOK, true, "successfully"}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to encode response JSON", "success": false, "error": err})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else if employee.EmpRole != 3 {
		query = fmt.Sprintf("SELECT %s FROM project where operations_manager = %d ORDER BY projectName", fields, managerIDInt)
		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false, "error": err})
			return
		}
		defer rows.Close()
		projects := []Project2{}
		for rows.Next() {
			var project Project2
			err := rows.Scan(&project.ID, &project.ProjectName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false, "error": err})
				return
			}
			projects = append(projects, project)
		}

		response := Response11{projects, http.StatusOK, true, "successfully"}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to encode response JSON", "success": false, "error": err})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else {
		query = fmt.Sprintf("SELECT id,projectName,COALESCE(countryId,''),COALESCE(StartDate,''),COALESCE(EndDate,''),COALESCE(FunderId,''),COALESCE(PartnerId,''),COALESCE(Training_target,''),COALESCE(LocationID,''),COALESCE(Operations_manager,''),COALESCE(BusID,''),COALESCE(DriverID,''),COALESCE(Project_status,''),COALESCE(Status,''),COALESCE(CreatedAt,''),COALESCE(CreatedBy,''),COALESCE(LastUpdatedAt,''),COALESCE(LastUpdatedBy,''),COALESCE(Gfl_id,'') from project p WHERE operations_manager in(SELECT id from employee e where e.supervisorId=%d) ORDER BY projectName", managerIDInt)

		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false, "error": err})
			return
		}
		defer rows.Close()
		projects := []Project1{}
		for rows.Next() {
			var project Project1
			err := rows.Scan(&project.ID, &project.ProjectName, &project.CountryID, &project.StartDate, &project.EndDate, &project.FunderID, &project.PartnerID, &project.Training_target, &project.LocationID, &project.Operations_manager, &project.BusID, &project.DriverID, &project.Project_status, &project.Status, &project.CreatedAt, &project.CreatedBy, &project.LastUpdatedAt, &project.LastUpdatedBy, &project.Gfl_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false, "error": err})
				return
			}
			projects = append(projects, project)
		}

		response := Response{projects, http.StatusOK, true, "successfully"}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to encode response JSON", "success": false, "error": err})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}

}
