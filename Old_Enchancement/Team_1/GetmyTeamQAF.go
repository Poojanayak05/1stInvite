package Team_1

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProfilePic   string `json:"profile_pic"`
	EmployeeRole int    `json:"empRole"`
	MyTeam       string `json:"myteam"`
	EmpID        int    `json:"emp_id"`
}

func GetMayTeamQAF(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Failed to read request body"})
		return
	}

	// Define a struct to hold the request parameters
	type Request struct {
		EmpID string `json:"emp_id"`
	}

	// Parse the request body into the Request struct
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Failed to read request body"})
		return
	}

	data := make([]Employee, 0)
	empID := req.EmpID
	if empID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Missing emp_id parameter"})
		return
	}

	id, err := strconv.Atoi(empID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid emp_id parameter"})
		return
	}
	query := `SELECT emp.empRole AS role
				FROM employee emp
				WHERE emp.id = ?`

	var role int
	err = DB.QueryRow(query, id).Scan(&role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
		return
	}

	var fields string
	switch role {
	case 12:
		fields = "emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id"
		query = `SELECT ` + fields + `
					FROM employee emp
					INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
					WHERE emp.supervisorId = ? AND emp.status = 1
					GROUP BY emp.id
					ORDER BY emp.first_name`
		rows, err := DB.Query(query, id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)

			if emp.EmployeeRole == 4 {
				query = `SELECT ` + fields + `
							FROM employee emp
							INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
							WHERE emp.empRole IN (13) AND emp.supervisorId = ? AND emp.status = 1
							GROUP BY emp.id
							ORDER BY emp.first_name`
				rows, err := DB.Query(query, emp.ID)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
					return
				}
				defer rows.Close()
				for rows.Next() {
					var subEmp Employee
					err := rows.Scan(&subEmp.ID, &subEmp.Name, &subEmp.ProfilePic, &subEmp.EmployeeRole, &subEmp.MyTeam, &subEmp.EmpID)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
						return
					}
					data = append(data, subEmp)
				}
			}
		}
	case 3:

		fields = "emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id"
		query = `SELECT ` + fields + `
					FROM employee emp
					INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
					WHERE emp.supervisorId = ? AND emp.status = 1
					GROUP BY emp.id
					ORDER BY emp.first_name`
		rows, err := DB.Query(query, id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)

			if emp.EmployeeRole == 4 || emp.EmployeeRole == 12 {
				query = `SELECT ` + fields + `
							FROM employee emp
							INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
							WHERE emp.empRole IN (4, 12, 13) AND emp.supervisorId = ? AND emp.status = 1
							GROUP BY emp.id
							ORDER BY emp.first_name`
				rows, err := DB.Query(query, emp.ID)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
					return
				}
				defer rows.Close()
				for rows.Next() {
					var subEmp Employee
					err := rows.Scan(&subEmp.ID, &subEmp.Name, &subEmp.ProfilePic, &subEmp.EmployeeRole, &subEmp.MyTeam, &subEmp.EmpID)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
						return
					}
					data = append(data, subEmp)

					if subEmp.EmployeeRole == 4 {
						query = `SELECT ` + fields + `
									FROM employee emp
									INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
									WHERE emp.empRole IN (13) AND emp.supervisorId = ? AND emp.status = 1
									GROUP BY emp.id
									ORDER BY emp.first_name`
						rows, err := DB.Query(query, subEmp.ID)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
							return
						}
						defer rows.Close()
						for rows.Next() {
							var subSubEmp Employee
							err := rows.Scan(&subSubEmp.ID, &subSubEmp.Name, &subSubEmp.ProfilePic, &subSubEmp.EmployeeRole, &subSubEmp.MyTeam, &subSubEmp.EmpID)
							if err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
								return
							}
							data = append(data, subSubEmp)
						}
					}
				}
			}
		}
	case 1:
		query := `SELECT emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id
			  FROM employee emp
			  INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
			  WHERE emp.empRole IN (3, 4, 12, 13) AND emp.status = 1
			  GROUP BY emp.id
			  ORDER BY emp.first_name`

		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)
		}
	case 4:

		query := `SELECT emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id
				  FROM employee emp
				  INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
				  WHERE emp.supervisorId = ? AND emp.status = 1
				  GROUP BY emp.id
				  ORDER BY emp.first_name`

		rows, err := DB.Query(query, id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)

			if emp.EmployeeRole == 13 {
				querys := `SELECT emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id
						   FROM employee emp
						   INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
						   WHERE emp.empRole IN (6) AND emp.supervisorId = ? AND emp.status = 1
						   GROUP BY emp.id
						   ORDER BY emp.first_name`

				rows2, err := DB.Query(querys, emp.ID)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
					return
				}
				defer rows2.Close()

				for rows2.Next() {
					var emp2 Employee
					err := rows2.Scan(&emp2.ID, &emp2.Name, &emp2.ProfilePic, &emp2.EmployeeRole, &emp2.MyTeam, &emp2.EmpID)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
						return
					}
					data = append(data, emp2)
				}
			}
		}

	case 13:
		fields := "emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id"
		if role == 3 {
			fields = "emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id"
		}

		query := `SELECT ` + fields + `
				  FROM employee emp
				  INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
				  WHERE emp.supervisorId = ? AND emp.status = 1
				  GROUP BY emp.id
				  ORDER BY emp.first_name`

		rows, err := DB.Query(query, id) // Assuming 'id' is a variable with the supervisor ID value
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)
		}
	case 11:
		fields := "emp.id, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS name, IFNULL(emp.profile_pic, '') AS profile_pic, emp.empRole, '' AS myteam, qaf.emp_id"
		query := `SELECT ` + fields + `
				  FROM employee emp
				  INNER JOIN QualityAssessmentForm_1 qaf ON emp.id = qaf.emp_id
				  WHERE emp.empRole = 4 AND emp.status = 1
				  GROUP BY emp.id
				  ORDER BY emp.first_name`

		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var emp Employee
			err := rows.Scan(&emp.ID, &emp.Name, &emp.ProfilePic, &emp.EmployeeRole, &emp.MyTeam, &emp.EmpID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
				return
			}
			data = append(data, emp)
		}

	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "Message": "Invalid emp_id parameter or no data found", "success": false, "data": data})
		return
	}

	response := map[string]interface{}{
		"data":    data,
		"code":    200,
		"message": "Successfully",
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500", "Message": "StatusInternalServerError"})
		return
	}
}
