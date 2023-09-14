package Team_2

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateProject(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"status": "404 Not Found", "Error": "Method not found"})
		return
	}
	busID := ""
	locationID := r.FormValue("locationID")
	locationName := r.FormValue("location_name")
	funderID := r.FormValue("funderID")
	funderName := r.FormValue("funder_name")
	createdBy := r.FormValue("createdBy")
	lastUpdatedBy := r.FormValue("lastUpdatedBy")
	projectID := r.FormValue("project_id")
	publish := r.FormValue("publish")
	partnerID := r.FormValue("partnerID")
	driverID := r.FormValue("driverID")
	busID = r.FormValue("busID")
	operationsManagerID := r.FormValue("operations_manager_id")

	type ProjectData struct {
		ProjectID    int    `json:"project_id"`
		ProjectName  string `json:"projectName"`
		LocationID   int    `json:"locationID"`
		LocationName string `json:"location_name"`
	}
	var queryData ProjectData
	if locationID != "" && funderID != "" && projectID == "" && partnerID == "" && operationsManagerID == "" && busID == "" && driverID == "" {

		query := "INSERT INTO project (funderID, locationID, createdBy, lastUpdatedBy) values (?,?,?,?)"
		funderIDInt, _ := strconv.Atoi(funderID)
		locationIDInt, _ := strconv.Atoi(locationID)
		createdByInt, _ := strconv.Atoi(createdBy)
		lastUpdatedByInt, _ := strconv.Atoi(lastUpdatedBy)

		result, err := DB.Exec(query, funderIDInt, locationIDInt, createdByInt, lastUpdatedByInt)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}

		insertedID, err := result.LastInsertId()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		if insertedID > 0 {
			GenericProjectName := generateProjectName(locationName, funderName, insertedID)

			outPut := updateProjectID(DB, insertedID, GenericProjectName)

			if outPut {
				fields := "pr.id as project_id, pr.projectName, pr.locationID, loc.name as location_name"

				q := "SELECT " + fields + " from project pr left join location loc on pr.locationID = loc.id  where pr.id = ?"

				rows, err := DB.Query(q, insertedID)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
					return
				}
				for rows.Next() {
					err := rows.Scan(&queryData.ProjectID, &queryData.ProjectName, &queryData.LocationID, &queryData.LocationName)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
				}
				type ProjectResponse struct {
					ProjectID    string `json:"project_id"`
					ProjectName  string `json:"projectName"`
					LocationID   string `json:"locationID"`
					LocationName string `json:"location_name"`
					Success      bool   `json:"success"`
					Message      string `json:"message"`
				}

				var resp ProjectResponse

				resp.ProjectID = strconv.Itoa(queryData.ProjectID)
				resp.ProjectName = queryData.ProjectName
				resp.LocationName = queryData.LocationName
				resp.LocationID = strconv.Itoa(queryData.LocationID)
				resp.Message = "Project Added Successfully"
				resp.Success = true

				json.NewEncoder(w).Encode(resp)
				return

			}

		}
	} else {

		if !(publish != "") && projectID != "" {

			query := ""

			if partnerID := r.FormValue("partnerID"); partnerID != "" {
				query += "partnerID = '" + partnerID + "',"
			}

			if trainingTarget := r.FormValue("training_target"); trainingTarget != "" {
				query += "training_target = '" + trainingTarget + "',"
			}

			if startDate := r.FormValue("start_date"); startDate != "" {
				parsedStartDate, err := time.Parse("02-01-2006", startDate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				startDate := parsedStartDate.Format("2006-01-02")
				query += "startDate = '" + startDate + "',"
			}

			if endDate := r.FormValue("end_date"); endDate != "" {
				parsedEndDate, err := time.Parse("02-01-2006", endDate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				endDate := parsedEndDate.Format("2006-01-02")
				query += "endDate = '" + endDate + "',"
			}

			if operationsManagerID := r.FormValue("operations_manager_id"); operationsManagerID != "" {
				query += "operations_manager = '" + operationsManagerID + "',"
			}

			if gflID := r.FormValue("gfl_id"); gflID != "" {
				query += "gfl_id = '" + gflID + "',"
			}

			if busID := r.FormValue("busID"); busID != "" {
				query += "busID = '" + busID + "',"
			}

			if driverID := r.FormValue("driverID"); driverID != "" {
				query += "driverID = '" + driverID + "',"
			}

			if createdBy := r.FormValue("createdBy"); createdBy != "" {
				query += "createdBy = '" + createdBy + "',"
			}

			if lastUpdatedBy := r.FormValue("lastUpdatedBy"); lastUpdatedBy != "" {
				query += "lastUpdatedBy = '" + lastUpdatedBy + "',"
			}
			query = strings.TrimSuffix(query, ",")
			updateQuery := "UPDATE project SET " + query + " WHERE id = '" + projectID + "'"

			_, err := DB.Exec(updateQuery)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "message": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Project Updated Successfully"})
			return

		} else {

			var msg string

			if r.FormValue("start_date") == "" {
				msg = "Start Date"
			}

			if r.FormValue("end_date") == "" {
				msg = "End date"
			}

			if r.FormValue("partnerID") == "" {
				msg = "Partner"
			}

			if r.FormValue("operations_manager_id") == "" {
				msg = "Operations Manager"
			}

			if r.FormValue("busID") == "" {
				msg = "Bus"
			}

			if r.FormValue("driverID") == "" {
				msg = "Driver"
			}

			if msg != "" {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "" + msg + "is required", "response": []interface{}{}})
				return
			}

			projectID := r.FormValue("project_id")
			partnerID := r.FormValue("partnerID")
			trainingTarget := r.FormValue("training_target")

			startDateStr := r.FormValue("start_date")
			parsedStartDate, err := time.Parse("02-01-2006", startDateStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			startDate := parsedStartDate.Format("2006-01-02")

			endDateStr := r.FormValue("end_date")
			parsedEndDate, err := time.Parse("02-01-2006", endDateStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			endDate := parsedEndDate.Format("2006-01-02")

			operationsManagerID := r.FormValue("operations_manager_id")
			gflID := r.FormValue("gfl_id")
			busID := r.FormValue("busID")
			driverID := r.FormValue("driverID")
			// createdBy := r.FormValue("createdBy")
			lastUpdatedBy, _ := strconv.Atoi(r.FormValue("user_id"))
			projectStatus := "1"

			updateQuery := "UPDATE project SET partnerID = '" + partnerID + "', training_target = '" + trainingTarget + "', startDate = '" + startDate + "', endDate = '" + endDate + "', operations_manager = '" + operationsManagerID + "',gfl_id='" + gflID + "', busID = '" + busID + "', driverID = '" + driverID + "', project_status = '" + projectStatus + "', lastUpdatedBy = ? WHERE id = '" + projectID + "'"

			result, err := DB.Exec(updateQuery, lastUpdatedBy)
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
				return
			}
			AffectedRowss, err := result.RowsAffected()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			if AffectedRowss > 0 {
				associatedProjectsValues := r.Form["associatedProject"]

				if len(associatedProjectsValues) > 0 {
					var associatedProjects []int

					for _, value := range associatedProjectsValues {
						intValue, err := strconv.Atoi(value)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
							return
						}
						associatedProjects = append(associatedProjects, intValue)
					}

					for _, v := range associatedProjects {

						insertQuery := "INSERT INTO project_association (projectId, associatedProject) VALUES (?, ?)"

						rows, err := DB.Exec(insertQuery, projectID, v)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
							return
						}

						check, _ := rows.LastInsertId()
						if check > 0 {
							w.WriteHeader(http.StatusOK)
							json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Project Publishedd Successfully", "response": queryData})
							return

						}
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "An error occurred while processing the request", "response": []interface{}{}})

					}
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Project Publishedd Successfully"})
				return

			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Project already Published"})

		}
	}
}

func updateProjectID(DB *sql.DB, insertedID int64, projectName string) bool {
	updateQuery := "UPDATE project SET projectName = ? WHERE id = ?"

	_, err := DB.Exec(updateQuery, projectName, insertedID)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func generateProjectName(location, funder string, insertedID int64) string {
	funder = funder[:2]
	year := strconv.Itoa(time.Now().Year() % 100)

	projectName := location + funder + year + strconv.FormatInt(insertedID, 10)

	return projectName
}
