package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Employee represents the structure of an employee
type EmployeeRequest struct {
	Emp_id string `json:"emp_id"`
	Team   string `json:"team"`
}
type EmployeeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Profile_Pic string `json:"profile_pic"`
	Emp_Role    string `json:"empRole"`
	MyTeam      string `json:"myteam"`
}

// GetMyTeam retrieves the team information for a given employee ID

func GetMyTeam(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	//It checks if the HTTP method is GET and returns a 404 Not Found error if it's not.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	var emp EmployeeRequest
	var emp1 EmployeeResponse

	//The function decodes the JSON request body into an EmployeeRequest struct.
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	//It performs some conversions on the request data using strconv.Atoi to convert string values to integers.
	Emp_id, _ := strconv.Atoi(emp.Emp_id)

	//If the Emp_id is 0, it is set to a default value of 15.
	if Emp_id == 0 {
		Emp_id = 15
	}

	if Emp_id != 0 {
		var fields, query string

		// Query the database to retrieve the employee's role
		query = fmt.Sprintf("SELECT emp.empRole as role FROM employee emp where emp.id = %d", Emp_id)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}
		defer rows.Close()
		var role int
		for rows.Next() {
			err := rows.Scan(&role)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}

		}
		var data []EmployeeResponse

		// Based on the role, query and retrieve team information
		if role == 12 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.supervisorId = %d AND emp.status = 1 ORDER BY emp.first_name", fields, Emp_id)
			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)

				if emp1.Emp_Role == "4" {
					query = fmt.Sprintf("select %s from employee emp where emp.empRole in (13) and emp.supervisorId = %s AND emp.status = 1 ORDER BY emp.first_name", fields, emp1.ID)
					rows, err := DB.Query(query)
					if err != nil {

						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer rows.Close()

					for rows.Next() {
						err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
						if err != nil {

							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}
						//The retrieved data is appended to the data slice of EmployeeResponse structs.
						data = append(data, emp1)
					}
				}

			}

		}

		if role == 3 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.supervisorId = %d AND emp.status = 1 ORDER BY emp.first_name", fields, Emp_id)
			rows, err := DB.Query(query)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}

				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)
				if emp1.Emp_Role == "4" || emp1.Emp_Role == "12" {

					query = fmt.Sprintf("select %s from employee emp where emp.empRole in (4,12,13) and emp.supervisorId = %s  AND emp.status = 1 ORDER BY emp.first_name", fields, emp1.ID)
					rows, err := DB.Query(query)
					if err != nil {

						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer rows.Close()

					for rows.Next() {
						err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
						if err != nil {

							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}

						data = append(data, emp1)
						if emp1.Emp_Role == "4" {
							query = fmt.Sprintf("select %s from employee emp where emp.empRole in (13) and emp.supervisorId = %s AND emp.status = 1 ORDER BY emp.first_name", fields, emp1.ID)
							rows, err := DB.Query(query)
							if err != nil {

								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
								return
							}
							defer rows.Close()

							for rows.Next() {
								err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
								if err != nil {

									w.WriteHeader(http.StatusBadRequest)
									json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
									return
								}
								//The retrieved data is appended to the data slice of EmployeeResponse structs.
								data = append(data, emp1)
							}
						}
					}

				}
			}

		}

		if role == 1 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.empRole in (3,4,12,13) AND emp.status = 1 ORDER BY emp.first_name", fields)
			rows, err := DB.Query(query)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)
			}
		}
		if role == 4 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.supervisorId = %d AND emp.status = 1 ORDER BY emp.first_name", fields, Emp_id)
			rows, err := DB.Query(query)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)
				if emp1.Emp_Role == "13" {
					query = fmt.Sprintf("select %s from employee emp where emp.empRole in (6) and emp.supervisorId = %s AND emp.status = 1 ORDER BY emp.first_name", fields, emp1.Emp_Role)
					rows, err := DB.Query(query)
					if err != nil {

						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer rows.Close()

					for rows.Next() {
						err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
						if err != nil {

							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}
						//The retrieved data is appended to the data slice of EmployeeResponse structs.
						data = append(data, emp1)
					}

				}
			}

		}
		if role == 13 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			if role == 3 {
				fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			}
			//fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.supervisorId = %d AND emp.status = 1 ORDER BY emp.first_name", fields, Emp_id)

			rows, err := DB.Query(query)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)
			}

		} else if role == 11 {
			fields = "emp.id, CONCAT(emp.first_name,' ', IFNULL(emp.last_name,'')) AS name, IFNULL(emp.profile_pic,'') AS profile_pic, emp.empRole, '' as myteam"
			query = fmt.Sprintf("select %s from employee emp where emp.empRole = 4 AND emp.status = 1 ORDER BY emp.first_name", fields)
			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&emp1.ID, &emp1.Name, &emp1.Profile_Pic, &emp1.Emp_Role, &emp1.MyTeam)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				//The retrieved data is appended to the data slice of EmployeeResponse structs.
				data = append(data, emp1)
			}
		}
		// Encode the retrieved data and send it as a response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": data, "code": 200, "success": true, "message": "Successfully"})

	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Enter required fields", "success": false, "Code": 200})
		return
	}
}
