package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Project struct {
	ProjectID      int              `json:"project_id"`
	ProjectName    string           `json:"project_name"`
	LocationName   string           `json:"location_name"`
	PartnerName    string           `json:"partnerName"`
	TrainingTarget string           `json:"training_target"`
	TargetList     []TargetEmployee `json:"target_list"`
}

type TargetEmployee struct {
	EmpID     int    `json:"emp_id"`
	EmpName   string `json:"emp_name"`
	EmpTarget string `json:"emp_target"`
	RoleID    int    `json:"role_id"`
}

type Response1 struct {
	List    Project `json:"list"`
	Code    int     `json:"code"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

func GetAssignTargets(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}

	if projectID, ok := request["project_id"].(string); ok {
		// id := int(projectID)

		fields := "prj.id as project_id, UPPER(prj.projectName) as project_name, par.partnerName, prj.training_target, loc.name as location_name, prj_emp.emp_id, emp.first_name, IFNULL(prj_emp.target, '') as target, prj_emp.role_id as role_id"
		query := "SELECT " + fmt.Sprint(fields) + " FROM project prj LEFT JOIN partner par ON prj.partnerID = par.partnerID LEFT JOIN project_emps prj_emp ON prj.id = prj_emp.project_id LEFT JOIN location loc ON prj.locationID = loc.id LEFT JOIN employee emp ON prj_emp.emp_id = emp.id WHERE prj.id = " + fmt.Sprint(projectID) + " AND prj_emp.role_id = 5 GROUP BY prj_emp.emp_id ORDER BY projectName"

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}
		defer rows.Close()

		var project Project
		var response Response1
		var target TargetEmployee
		var targetList []TargetEmployee

		for rows.Next() {

			err := rows.Scan(
				&project.ProjectID,
				&project.ProjectName,
				&project.PartnerName,
				&project.TrainingTarget,
				&project.LocationName,
				&target.EmpID,
				&target.EmpName,
				&target.EmpTarget,
				&target.RoleID,
			)

			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})

			}
			targetList = append(targetList, target)
		}

		response.Code = 200
		response.Success = true
		response.Message = "Successfully"
		project.TargetList = targetList
		response.List = project
		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
}
