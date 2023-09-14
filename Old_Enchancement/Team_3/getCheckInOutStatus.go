package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetCheckInOutStatus(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type Req struct {
		Tb_id      string `json:"tb_id"`
		Type       string `json:"type"`
		Project_id string `json:"project_id"`
		Poa_type   string `json:"poa_type"`
	}

	var request Req

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := make(map[string]interface{})

	if request.Project_id != "" && request.Tb_id != "" && request.Type != "" {
		tb_id := request.Tb_id

		t, _ := strconv.Atoi(request.Type)

		project_id, _ := strconv.Atoi(request.Project_id)

		poa_type, _ := strconv.Atoi(request.Poa_type)

		var query string

		var projExist bool
		err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM project WHERE id = ?)", request.Project_id).Scan(&projExist)
		if err != nil {

			// Print the error and return a JSON response with error message

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
			return
		}

		if !projExist {
			// Print the RoleID not found message

			// Return a JSON response if RoleID doesn't exist
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "No Project_id found"})
			return
		}

		var tbExists bool

		err1 := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM tbl_poa WHERE tb_id = ?)", request.Tb_id).Scan(&tbExists)
		if err != nil {

			// Print the error and return a JSON response with error message
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err1, "Status Code": "400 Bad Request"})
			return
		}

		if !tbExists {
			// Print the RoleID not found message
			// Return a JSON response if RoleID doesn't exist
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "No tb_id found"})
			return
		}

		if t == 1 {
			query = fmt.Sprintf("SELECT operations_manager as id from project where id = %d ", project_id)

		} else if t == 2 {

			query = fmt.Sprintf("SELECT emp_id as id from project_emps where project_id = %d AND role_id = 5", project_id)

		} else if t == 13 {

			query = fmt.Sprintf("SELECT emp_id as id from project_emps where project_id = %d AND role_id = 13", project_id)

		} else {

			query = fmt.Sprintf("SELECT emp_id as id from project_emps where project_id = %d AND role_id = 6", project_id)

		}

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		var idslice []string
		for rows.Next() {
			var id int
			err := rows.Scan(&id)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			idslice = append(idslice, strconv.Itoa(id))
		}
		ids := strings.Join(idslice, ", ")

		fields := "tbl_poa.id, IFNULL(DATE_FORMAT(tbl_poa.check_in, '%d-%m-%Y'),'') as check_in_date_day, IFNULL(TIME_FORMAT(tbl_poa.check_in, '%h:%i %p'), '') as check_in_time_day, IFNULL(DATE_FORMAT(tbl_poa.check_out, '%d-%m-%Y'),'') as check_out_date_day, IFNULL(TIME_FORMAT(tbl_poa.check_out, '%h:%i %p'), '') as check_out_time_day, IFNULL(tbl_poa.check_in_location, '') as check_in_location_day, IFNULL(tbl_poa.check_out_location, '') as check_out_location_day, IFNULL(CONCAT(emp.first_name, ' ', emp.last_name), '') as employee_name, IFNULL(TIME_FORMAT(tbl_poa.date, '%h:%i %p'), '') as actual_time, IFNULL(tbl_poa.check_in_lat, '0.00') as check_in_lat, IFNULL(tbl_poa.check_in_lon, '0.00') as check_in_lon, IFNULL(tbl_poa.check_out_lat, '0.00') as check_out_lat, IFNULL(tbl_poa.check_out_lon, '0.00') as check_out_lon, tbl_poa.primary_id, tbl_poa.type"
		var idsQry string
		if ids != "" {
			idsQry = fmt.Sprintf(" AND tbl_poa.user_id IN (%s) ", ids)

		}

		if poa_type == 1 {
			query = fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa LEFT JOIN employee emp on tbl_poa.user_id = emp.id WHERE tbl_poa.type = '1' AND tbl_poa.tb_id = %s %s ORDER BY tbl_poa.primary_id asc", fields, tb_id, idsQry)
		} else {
			query = fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa LEFT JOIN employee emp on tbl_poa.user_id = emp.id WHERE tbl_poa.type = '2' AND tbl_poa.primary_id = %s %s", fields, tb_id, idsQry)
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

		type Querydata struct {
			Id                     string `json:"id"`
			Check_in_date_day      string `json:"check_in_date_day"`
			Check_in_time_day      string `json:"check_in_time_day"`
			Check_out_date_day     string `json:"check_out_date_day"`
			Check_out_time_day     string `json:"check_out_time_da"`
			Check_in_location_day  string `json:"check_in_location_day"`
			Check_out_location_day string `json:"check_out_location_day"`
			Employee_name          string `json:"employee_name"`
			Actual_time            string `json:"actual_time"`
			Check_in_lat           string `json:"check_in_lat"`
			Check_in_lon           string `json:"check_in_lon"`
			Check_out_lat          string `json:"check_out_lat"`
			Check_out_lon          string `json:"check_out_lon"`
			Primary_id             int    `json:"primary_id"`
			Typ                    int    `json:"typ"`
		}
		var data []Querydata
		for rows.Next() {

			var q Querydata
			err := rows.Scan(&q.Id, &q.Check_in_date_day, &q.Check_in_time_day, &q.Check_out_date_day, &q.Check_out_time_day, &q.Check_in_location_day, &q.Check_out_location_day, &q.Employee_name, &q.Actual_time, &q.Check_in_lat, &q.Check_in_lon, &q.Check_out_lat, &q.Check_out_lon, &q.Primary_id, &q.Typ)
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
			data = append(data, q)
		}

		list := make(map[string]interface{})
		day := 1
		//from_tb := 0
		day_type := 0

		if tb_id != "" && poa_type == 1 {
			//check its day1 or day2
			queryDay := fmt.Sprintf("SELECT MIN(primary_id) as primary_id from tbl_poa where tb_id = %s and added = 0", tb_id)
			rows, err = DB.Query(queryDay)

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

			var primary_id int
			for rows.Next() {
				rows.Scan(&primary_id)
			}
			if len(data) > 0 {
				if primary_id == data[0].Primary_id {
					day_type = 1
				} else {
					day_type = 2
				}
			}
		}

		for key, val := range data {
			if len(data) == 2 {
				day = key + day
				list["check_in_date_day"+strconv.Itoa(day)] = val.Check_in_date_day
				list["check_in_time_day"+strconv.Itoa(day)] = val.Check_in_time_day
				list["check_out_date_day"+strconv.Itoa(day)] = val.Check_out_date_day
				list["check_out_time_day"+strconv.Itoa(day)] = val.Check_out_time_day
				list["check_in_location_day"+strconv.Itoa(day)] = val.Check_in_location_day
				list["check_out_location_day"+strconv.Itoa(day)] = val.Check_out_location_day
				list["employee_name"] = val.Employee_name
				list["actual_time_day"+strconv.Itoa(day)] = val.Actual_time
				list["check_in_lat_day"+strconv.Itoa(day)] = val.Check_in_lat
				list["check_in_lon_day"+strconv.Itoa(day)] = val.Check_in_lon
				list["check_out_lat_day"+strconv.Itoa(day)] = val.Check_out_lat
				list["check_out_lon_day"+strconv.Itoa(day)] = val.Check_out_lon

			} else {
				if val.Typ == 1 {
					if day_type == 1 {
						day = key + day
						list["check_in_date_day"+strconv.Itoa(day)] = val.Check_in_date_day
						list["check_in_time_day"+strconv.Itoa(day)] = val.Check_in_time_day
						list["check_out_date_day"+strconv.Itoa(day)] = val.Check_out_date_day
						list["check_out_time_day"+strconv.Itoa(day)] = val.Check_out_time_day
						list["check_in_location_day"+strconv.Itoa(day)] = val.Check_in_location_day
						list["check_out_location_day"+strconv.Itoa(day)] = val.Check_out_location_day
						list["employee_name"] = val.Employee_name
						list["actual_time_day"+strconv.Itoa(day)] = val.Actual_time
						list["check_in_lat_day"+strconv.Itoa(day)] = val.Check_in_lat
						list["check_in_lon_day"+strconv.Itoa(day)] = val.Check_in_lon
						list["check_out_lat_day"+strconv.Itoa(day)] = val.Check_out_lat
						list["check_out_lon_day"+strconv.Itoa(day)] = val.Check_out_lon

						list["check_in_date_day2"] = ""
						list["check_in_time_day2"] = ""
						list["check_out_date_day2"] = ""
						list["check_out_time_day2"] = ""
						list["check_in_location_day2"] = ""
						list["check_out_location_day2"] = ""
						list["actual_time_day2"] = ""
						list["check_in_lat_day2"] = ""
						list["check_in_lon_day2"] = ""
						list["check_out_lat_day2"] = ""
						list["check_out_lon_day2"] = ""

					} else {
						day = 2
						list["check_in_date_day1"] = ""
						list["check_in_time_day1"] = ""
						list["check_out_date_day1"] = ""
						list["check_out_time_day1"] = ""
						list["check_in_location_day1"] = ""
						list["check_out_location_day1"] = ""
						list["actual_time_day1"] = ""
						list["check_in_lat_day1"] = ""
						list["check_in_lon_day1"] = ""
						list["check_out_lat_day1"] = ""
						list["check_out_lon_day1"] = ""
						list["check_in_date_day"+strconv.Itoa(day)] = val.Check_in_date_day
						list["check_in_time_day"+strconv.Itoa(day)] = val.Check_in_time_day
						list["check_out_date_day"+strconv.Itoa(day)] = val.Check_out_date_day
						list["check_out_time_day"+strconv.Itoa(day)] = val.Check_out_time_day
						list["check_in_location_day"+strconv.Itoa(day)] = val.Check_in_location_day
						list["check_out_location_day"+strconv.Itoa(day)] = val.Check_out_location_day
						list["employee_name"] = val.Employee_name
						list["actual_time_day"+strconv.Itoa(day)] = val.Actual_time
						list["check_in_lat_day"+strconv.Itoa(day)] = val.Check_in_lat
						list["check_in_lon_day"+strconv.Itoa(day)] = val.Check_in_lon
						list["check_out_lat_day"+strconv.Itoa(day)] = val.Check_out_lat
						list["check_out_lon_day"+strconv.Itoa(day)] = val.Check_out_lon
					}
				} else {
					day = key + day
					list["check_in_date_day"+strconv.Itoa(day)] = val.Check_in_date_day
					list["check_in_time_day"+strconv.Itoa(day)] = val.Check_in_time_day
					list["check_out_date_day"+strconv.Itoa(day)] = val.Check_out_date_day
					list["check_out_time_day"+strconv.Itoa(day)] = val.Check_out_time_day
					list["check_in_location_day"+strconv.Itoa(day)] = val.Check_in_location_day
					list["check_out_location_day"+strconv.Itoa(day)] = val.Check_out_location_day
					list["employee_name"] = val.Employee_name
					list["actual_time_day"+strconv.Itoa(day)] = val.Actual_time
					list["check_in_lat_day"+strconv.Itoa(day)] = val.Check_in_lat
					list["check_in_lon_day"+strconv.Itoa(day)] = val.Check_in_lon
					list["check_out_lat_day"+strconv.Itoa(day)] = val.Check_out_lat
					list["check_out_lon_day"+strconv.Itoa(day)] = val.Check_out_lon
				}
			}
		}

		if len(data) == 1 {
			if poa_type != 1 {
				list["check_in_date_day2"] = ""
				list["check_in_time_day2"] = ""
				list["check_out_date_day2"] = ""
				list["check_out_time_day2"] = ""
				list["check_in_location_day2"] = ""
				list["check_out_location_day2"] = ""
				list["actual_time_day2"] = ""
				list["check_in_lat_day2"] = ""
				list["check_in_lon_day2"] = ""
				list["check_out_lat_day2"] = ""
				list["check_out_lon_day2"] = ""
			}
		}
		if len(list) == 0 {

			list["check_in_date_day1"] = ""
			list["check_in_time_day1"] = ""
			list["check_out_date_day1"] = ""
			list["check_out_time_day1"] = ""
			list["check_in_location_day1"] = ""
			list["check_out_location_day1"] = ""
			list["employee_name"] = ""
			list["actual_time_day1"] = ""
			list["check_in_lat_day1"] = ""
			list["check_in_lon_day1"] = ""
			list["check_out_lat_day1"] = ""
			list["check_out_lon_day1"] = ""
			list["check_in_date_day2"] = ""
			list["check_in_time_day2"] = ""
			list["check_out_date_day2"] = ""
			list["check_out_time_day2"] = ""
			list["check_in_location_day2"] = ""
			list["check_out_location_day2"] = ""
			list["actual_time_day2"] = ""
			list["check_in_lat_day2"] = ""
			list["check_in_lon_day2"] = ""
			list["check_out_lat_day2"] = ""
			list["check_out_lon_day2"] = ""
		}

		response["data"] = list
		response["code"] = http.StatusOK
		response["success"] = true
		response["message"] = "successfully"

		jsonData, err := json.Marshal(response)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "StatusInternalServerError",
				"success": false,
				"error":   err,
			})
			return
		}

		w.Write(jsonData)

	} else {
		jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid data", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid data",
				"success": false,
				"error":   err,
			})

			return
		}
		w.Write(jsonData)
	}
}
