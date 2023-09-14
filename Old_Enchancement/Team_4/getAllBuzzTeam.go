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

type Employee struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	ProfilePic string `json:"profile_pic"`
}

type Response3 struct {
	List    []Employee `json:"list"`
	Code    int        `json:"code"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

func GetAllBuzzTeam(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	roleID := r.PostFormValue("role_id")
	name := r.PostFormValue("name")

	//role_master
	//1 = CEO, 2 = Admin, 3 = Program Manager, 4 = operations Manager, 5 = Trainer, 6 = Gelathi Facilitor, 7 = Diver,
	// 8 = Funder, 9 = Partner

	var roleCondition string

	switch roleID {
	case "0":
		roleCondition = ""
	case "2", "11":
		roleCondition = "and empRole IN (1) and first_name LIKE '%" + name + "%'"
	case "3":
		roleCondition = "and empRole IN (1) and first_name LIKE '%" + name + "%'"
	case "4":
		roleCondition = "and empRole IN (3,12) and first_name LIKE '%" + name + "%'"
	case "12":
		roleCondition = "and empRole IN (3) and first_name LIKE '%" + name + "%'"
	case "6":
		roleCondition = "and empRole IN (4,13) and first_name LIKE '%" + name + "%'"
	case "5", "13":
		roleCondition = "and empRole IN (4) and first_name LIKE '%" + name + "%'"
	case "7":
		roleCondition = "and empRole IN (5) and first_name LIKE '%" + name + "%'"
	default:
		roleCondition = ""
	}

	fields := "id, CONCAT(first_name,' ',IFNULL(last_name, '')) AS name, empRole as role, profile_pic"
	query := fmt.Sprintf("SELECT %s FROM employee WHERE status = 1 %s ORDER BY first_name", fields, roleCondition)

	rows, err := DB.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Role, &employee.ProfilePic)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}
		employees = append(employees, employee)
	}

	response := Response3{
		List:    employees,
		Code:    200,
		Success: true,
		Message: "Successfully",
	}

	json.NewEncoder(w).Encode(response)

}
