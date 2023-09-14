package New

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Target      int    `json:"target"`
	Actual      int    `json:"actual"`
	Day2        int    `json:"day2"`
	Nosurvey    int    `json:"noofselfshakthisurvey"`
	Noofbatches int    `json:"noofbatches"`
	Villages    int    `json:"villages"`
	SelectType  string `json:"select_type"`
}

type Response struct {
	SummaryTarget   int       `json:"summary_target"`
	SummaryWomen    int       `json:"summary_women"`
	SummaryVillages int       `json:"summary_villages"`
	SummaryActual   int       `json:"summary_actual"`
	SummaryDay2     int       `json:"summary_day2"`
	SummaryEnrolled int       `json:"summary_enrolled"`
	SummaryGreen    int       `json:"summary_green"`
	SummaryVyapar   int       `json:"summary_vyapar"`
	Data            []Project `json:"data"`
	Code            int       `json:"code"`
	Success         bool      `json:"success"`
	Message         string    `json:"message"`
}

func SelfsakthiDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	data := []interface{}{}
	summaryTarget := 0
	summaryVillages := 0
	summaryActuals := 0
	summaryDay1 := 0
	summaryNoofbatches := 0

	summaryonselfshakthi := 0
	var day2Turnout float64

	// Parse the request body
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve values from the request
	roleID, err := getStringValue(request, "roleid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	empID, err := getStringValue(request, "emp_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	talukID, _ := getStringValue(request, "taluk")
	districtID, _ := getStringValue(request, "dist")
	startDate, _ := getStringValue(request, "start_date")
	endDate, _ := getStringValue(request, "end_date")
	funderID, _ := getStringValue(request, "funder_id")
	partnerID, _ := getStringValue(request, "partner_id")
	projectID, _ := getStringValue(request, "project_id")
	trainerID, _ := getStringValue(request, "trainer_id")
	opsManager, _ := getStringValue(request, "opsmanager")
	somID, _ := getStringValue(request, "somid")
	gflID, _ := getStringValue(request, "gflid")

	gflid, _ := strconv.Atoi(gflID)
	som, _ := strconv.Atoi(somID)
	opsid, _ := strconv.Atoi(opsManager)
	trainid, _ := strconv.Atoi(trainerID)
	projectid, _ := strconv.Atoi(projectID)
	funderid, _ := strconv.Atoi(funderID)
	talukid, _ := strconv.Atoi(talukID)
	distid, _ := strconv.Atoi(districtID)
	empid, _ := strconv.Atoi(empID)
	roleid, _ := strconv.Atoi(roleID)
	partnerid, _ := strconv.Atoi(partnerID)

	var isDateFilterApplied = false

	if startDate == "" || endDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 11 || roleid == 12 {
		filter := ""
		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}

		} else if roleid == 3 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				// var opsIds []int
				if som != 0 {
					opsIds, _ := getReportingOpsManagers(DB, som)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else if gflid != 0 {
					opsIds, _ := getSupervisor(DB, gflid)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					opsIds, _ := getReportingOpsManagers(DB, empid)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}

		} else if roleid == 12 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds, _ := getOpsManagers(DB, empid)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				}
			} else {
				filter = " and p.operations_manager in (0)"
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}
		} else if roleid == 4 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectIDs, _ := getOpProjects(DB, empid)
				if len(projectIDs) > 0 {
					filter = fmt.Sprintf(" AND p.operations_manager = %d", empid)
				} else {

					return
				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}
		} else if roleid == 9 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		}

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
			dateFilters = "date >= '" + startDate + "' and date <= '" + endDate + "'"

		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", partnerid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p INNER JOIN funder ON funder.funderID = p.funderID WHERE p.partnerID = %d AND %s %s", partnerid, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", partnerid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if distid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", distid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", talukid)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p INNER JOIN funder ON funder.funderID = p.funderID WHERE locationID = %d AND %s %s GROUP BY p.funderID", talukid, dateFilter, filter)
						filter += fmt.Sprintf(" AND locationID = %d", talukid)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)

						if err != nil {
							log.Println("SelfSakthiDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")

						w.Write(js)
						return

					}
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid locationid"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")

					w.Write(js)
					return

				}
			} else {
				getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + strconv.Itoa(distid)
				talukArray := []int{}
				talukRes, err := DB.Query(getTaluk)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				defer talukRes.Close()
				for talukRes.Next() {
					var tlk int
					err := talukRes.Scan(&tlk)
					if err != nil {
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					talukArray = append(talukArray, tlk)
				}
				funderListQuery = "SELECT p.funderID as id, funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ") and " + dateFilter + filter + " GROUP by p.funderID"
				filter += " and locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ")"
			}
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", funderid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT funderID AS id, funderName AS name FROM funder WHERE funderID = %d ", funderid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid funderid"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if projectid == 0 && trainid == 0 && opsid == 0 && som == 0 && gflid == 0 && !isDateFilterApplied && roleid != 4 {
			// Role 4 OpsManager Default should be project list
			funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p " +
				"INNER JOIN funder ON p.funderId = funder.funderID " +
				"WHERE " + dateFilter + filter
		}

		if len(funderListQuery) > 0 {
			rows, err := DB.Query(funderListQuery)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var funderID int
				var funderName string
				err = rows.Scan(&funderID, &funderName)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				getProj := ""

				projectArray := make([]int, 0)

				if startDate != "" && endDate != "" {
					getProj = fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %d AND p.startDate BETWEEN '%s' AND '%s'", funderID, startDate, endDate)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("ERROR>>", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						var startDate, endDate string
						err = projResult.Scan(&projectID, &startDate, &endDate)
						if err != nil {
							log.Println("ERROR>>", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
							return
						}

						projectArray = append(projectArray, projectID)
					}

				} else {
					getProj = fmt.Sprintf("SELECT id FROM project p WHERE funderID = %d and %s %s", funderID, dateFilter, filter)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("ERROR>>", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						err = projResult.Scan(&projectID)
						if err != nil {
							log.Println("ERROR>>", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
							return
						}

						projectArray = append(projectArray, projectID)
					}
				}

				if len(projectArray) == 0 {
					obj := make(map[string]interface{})
					obj["id"] = funderID
					obj["name"] = funderName
					obj["target"] = 0
					obj["actual"] = 0
					obj["day2"] = 0
					obj["noofselfshakthisurvey"] = 0
					obj["noofbatches"] = 0
					obj["villages"] = 0
					obj["startDate"] = ""
					obj["endDate"] = ""
					obj["select_type"] = "2"
					data = append(data, obj)
					continue
				}

				obj := make(map[string]interface{})
				obj["id"] = funderID
				obj["name"] = funderName
				obj["target"], _ = getTarget(DB, startDate, endDate, projectArray)
				summaryTarget += obj["target"].(int)
				obj["actual"], _ = getActual(DB, startDate, endDate, projectArray, "")
				summaryActuals += obj["actual"].(int)

				day1Count, _ := getDay1Count(DB, startDate, endDate, projectArray, "")
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout = float64(obj["actual"].(int)) / float64(day1Count)
					obj["day2"] = math.Round(day2Turnout * 100)
				} else {
					obj["day2"] = 0
				}
				obj["noofselfshakthisurvey"] = noOfselfshakthisurvey(DB, startDate, endDate, projectArray, "")
				summaryonselfshakthi += obj["noofselfshakthisurvey"].(int)
				obj["noofbatches"] = noofbatches(DB, startDate, endDate, projectArray, "")
				summaryNoofbatches += obj["noofbatches"].(int)
				obj["villages"], _ = newVillageCount(DB, startDate, endDate, projectArray, "")
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = ""
				obj["endDate"] = ""
				obj["select_type"] = "2"
				data = append(data, obj)
			}
		}

		var projectList string

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = fmt.Sprintf(" AND startDate BETWEEN '%s' AND '%s'", startDate, endDate)
				}

				projectList = fmt.Sprintf("SELECT id, projectName AS name, p.startDate, p.endDate FROM project p WHERE id = %d%s%s", projectid, filter, dateFilterNew)

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		} else if trainid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", trainid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
                                from tbl_poa tp
                                inner join project p on p.id = tp.project_id 
                                where user_id = %d and %s %s GROUP BY project_id`, trainid, dateFilter, filter)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if opsid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", opsid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (startDate == "" && endDate == "") {
					projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where operations_manager = %d and %s %s GROUP BY id`, opsid, dateFilter, filter)
				} else {
					projectList = fmt.Sprintf(`SELECT p.id, p.projectName as name, p.startDate, p.endDate
                            FROM project p
                            JOIN training_participants tp ON p.id = tp.project_id
                            WHERE p.operations_manager = %d AND p.startDate BETWEEN '%s' AND '%s'
                            GROUP BY p.id`, opsid, startDate, endDate)

				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if som > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", som)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where p.id in (
                                    SELECT DISTINCT project_id
                                    from tbl_poa tp
                                    where %s
                                ) and operations_manager in (
                                    SELECT id
                                    from employee e
                                    where e.supervisorId = %d
                                ) %s GROUP BY id`, dateFilters, som, filter)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if gflid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", gflid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where operations_manager in (
                                    SELECT supervisorId
                                    from employee e
                                    where e.id = %d
                                ) and %s %s GROUP BY id`, gflid, dateFilter, filter)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0 || roleid == 4 && distid == 0 {
			projectList = fmt.Sprintf("SELECT id, projectName as name, p.startDate, p.endDate from project p where %s %s", dateFilter, filter)
		}

		if len(projectList) > 0 {
			rows, err := DB.Query(projectList)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var prList Project
				err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				tpFilter := ""
				tbFilter := ""

				projectArray := []int{prList.ID}
				if trainid > 0 {
					prList.Target, _ = getTrainerTarget(DB, trainid, projectArray)
					summaryTarget += prList.Target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					prList.Target, _ = getTarget(DB, startDate, endDate, projectArray)
					summaryTarget += prList.Target

				}
				prList.Actual, _ = getActual(DB, startDate, endDate, projectArray, tpFilter)
				summaryActuals += prList.Actual
				day1Count, _ := getParticipantFilterDay1Count(DB, startDate, endDate, projectArray, tbFilter)
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout = float64(prList.Actual) / float64(day1Count)
					prList.Day2 = int(day2Turnout * 100)
				} else {
					prList.Day2 = 0
				}
				prList.Nosurvey = noOfselfshakthisurvey(DB, startDate, endDate, projectArray, "")
				summaryonselfshakthi += prList.Nosurvey
				prList.Noofbatches = noofbatches(DB, startDate, endDate, projectArray, "")
				summaryNoofbatches += prList.Noofbatches
				prList.Villages, _ = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += prList.Villages
				prList.SelectType = "1"

				data = append(data, prList)
			}
			if err := rows.Err(); err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
		}

		// Convert data to JSON and send response
		response := make(map[string]interface{})
		response["data"] = data
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals
		response["summary_Noofselfshakthisurvey"] = summaryonselfshakthi
		if summaryDay1 > 0 {
			day2Turnout = float64(summaryActuals) / float64(summaryDay1)
			response["summary_day2"] = int(math.Round(day2Turnout * 100))
		} else {
			response["summary_day2"] = 0
		}

		response["summary_NoofBatches"] = summaryNoofbatches
		response["summary_villages"] = summaryVillages

		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	} else if roleid == 5 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			var dateFilter string
			if isDateFilterApplied {
				dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			query := fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
			FROM tbl_poa tp
			INNER JOIN project p ON p.id = tp.project_id
			WHERE user_id = %d%s GROUP BY project_id`, empid, dateFilter)

			if projectid > 0 {
				query = fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
				FROM tbl_poa tp
				INNER JOIN project p ON p.id = tp.project_id
				WHERE user_id = %d AND tp.project_id = %d GROUP BY tp.project_id`, empid, projectid)
			}

			rows, err := DB.Query(query)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			var data []Project

			for rows.Next() {
				var prj Project
				err := rows.Scan(&prj.ID, &prj.Name, &prj.StartDate, &prj.EndDate)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				projectArray := []int{prj.ID}
				prj.Target, _ = getTrainerTarget(DB, empid, projectArray)
				summaryTarget += prj.Target
				filter := fmt.Sprintf(" and tp.trainer_id = %d", empid)

				prj.Actual, _ = getActual(DB, startDate, endDate, projectArray, filter)
				summaryActuals += prj.Actual
				day1Count, _ := getDay1Count(DB, startDate, endDate, projectArray, filter)
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout := float64(prj.Actual) / float64(day1Count)
					prj.Day2 = int(day2Turnout * 100)
				} else {
					prj.Day2 = 0
				}

				prj.Noofbatches = getTrainerBatches(DB, startDate, endDate, projectArray, empID)
				summaryNoofbatches += prj.Noofbatches

				prj.Villages = TrainerVillageCount(DB, startDate, endDate, projectArray, empID)
				summaryVillages += prj.Villages
				prj.SelectType = "1"

				prj.Nosurvey = noofsurveyssTrainer(DB, startDate, endDate, projectArray, empID)
				summaryonselfshakthi += prj.Nosurvey
				data = append(data, prj)
			}

			response := make(map[string]interface{})
			response["data"] = data
			response["summary_target"] = summaryTarget
			response["summary_actual"] = summaryActuals
			response["summary_Noofselfshakthisurvey"] = summaryonselfshakthi
			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				response["summary_day2"] = 0
			}

			response["summary_NoofBatches"] = summaryNoofbatches
			response["summary_villages"] = summaryVillages

			// Convert response to JSON
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}
	} else if roleid == 13 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {

			// filter := ""
			filters := ""
			// filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			// filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			var participantFilter = false

			if startDate != "" && endDate != "" {
				// filter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s'", startDate, endDate)
				filters = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s'", startDate, endDate)
				// filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate BETWEEN '%s' and '%s'", startDate, endDate)
				// filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate BETWEEN '%s' and '%s'", startDate, endDate)
			}

			if projectID != "" {
				filters += fmt.Sprintf(" and p.id = %s", projectID)
			}

			f := ""
			if gfID, ok := request["gfId"].(string); ok && gfID != "" {
				f = fmt.Sprintf(" and id=%s", gfID)
			}

			var (
				summaryCircleMeet   int
				summaryCircles      int
				summaryVillageVisit int
				summaryBeehive      int
				summaryEnroll       int
			)

			employeeIDs := make([]int, 0)
			rows, err := DB.Query("SELECT id FROM employee WHERE status = 1 AND supervisorId = ?"+f, empID)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()
			for rows.Next() {
				var id int
				err := rows.Scan(&id)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				employeeIDs = append(employeeIDs, id)
			}
			err = rows.Err()
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}

			getProjs := fmt.Sprintf(`SELECT project_id AS id, p.projectName AS name FROM tbl_poa tp
			INNER JOIN project p ON p.id = tp.project_id
			WHERE p.gfl_id = %d %s GROUP BY tp.project_id`, empid, filters)

			rows, err = DB.Query(getProjs)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var name string
				err := rows.Scan(&id, &name)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				obj := make(map[string]interface{})
				obj["id"] = id
				obj["name"] = name

				prjFilter := ""
				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				obj["circlemeet"] = getGFDataN(DB, prjFilter, 1, employeeIDs)
				summaryCircleMeet += obj["circlemeet"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and p.endDate BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				obj["circle"] = getGFCircleN(DB, prjFilter, employeeIDs)
				summaryCircles += obj["circle"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				}
				obj["villagevisit"] = getGFDataN(DB, prjFilter, 2, employeeIDs)
				summaryVillageVisit += obj["villagevisit"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)

				}
				obj["beehive"] = getGFDataN(DB, prjFilter, 3, employeeIDs)
				summaryBeehive += obj["beehive"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				projEnrolled := getGfEnrolledN(DB, prjFilter, employeeIDs)

				if participantFilter {
					if startDate != "" && endDate != "" {
						prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
						projEnrolled = getParticipantFilterGfEnrolledN(DB, prjFilter, employeeIDs, startDate, endDate)
					}
				}

				obj["enroll"] = projEnrolled
				summaryEnroll += obj["enroll"].(int)

				projectArray := ""
				// project_result := db.QueryRow("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id IN (?)", empId)
				project_result := DB.QueryRow("SELECT GROUP_CONCAT(DISTINCT id) as ids FROM project prj WHERE gfl_id IN (?)", empID)
				err = project_result.Scan(&projectArray)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				data = append(data, obj)
			}
			response := make(map[string]interface{})
			response["data"] = data
			response["summary_circle_meet"] = summaryCircleMeet
			response["summary_circles"] = summaryCircles
			response["summary_villagevisit"] = summaryVillageVisit

			response["summary_beehive"] = summaryBeehive
			response["summary_enroll"] = summaryEnroll

			// Convert response to JSON
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}
	} else if roleid == 6 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			participantFilter := ""
			filter := ""
			filters := ""
			// filterG := ""
			// filterV := ""
			if projectid > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", projectid)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
					filters = fmt.Sprintf(" and tp.enroll_date >= '%s' and tp.enroll_date <= '%s'", startDate, endDate)
					// filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", startDate, endDate)
					// filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", startDate, endDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					// filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					// filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet, _ := getGFData(DB, filter, 1, empid)
			villageVisit, _ := getGFData(DB, filter, 2, empid)
			beehive, _ := getGFData(DB, filter, 3, empid)
			enrolled, _ := getGfEnrolled(DB, filters, empid)
			circleVisit := getGFCircle(DB, filter, empid)

			data := make([]map[string]interface{}, 0)

			getProjs := fmt.Sprintf(`SELECT project_id as id, p.projectName as name FROM tbl_poa tp
						INNER JOIN project p ON p.id = tp.project_id
						WHERE tp.user_id = %d %s GROUP BY tp.project_id
						UNION
						SELECT tp.project_id as id, p.projectName as name FROM training_participants tp
						INNER JOIN project p ON tp.project_id = p.id
						WHERE enroll = 1 AND gelathi_id = %d %s`, empid, filter, empid, filter)

			if projectid > 0 {
				getProjs = fmt.Sprintf(`SELECT project_id as id, p.projectName as name FROM tbl_poa tp
						INNER JOIN project p ON p.id = tp.project_id
						WHERE tp.project_id = %d AND tp.user_id = %d GROUP BY tp.project_id
						UNION
						SELECT tp.project_id as id, p.projectName as name FROM training_participants tp
						INNER JOIN project p ON tp.project_id = p.id
						WHERE enroll = 1 AND gelathi_id = %d AND tp.project_id = %d`, projectid, empid, empid, projectid)
			}

			projectsList, err := DB.Query(getProjs)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer projectsList.Close()

			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				obj := make(map[string]interface{})
				obj["id"] = id

				obj["name"] = name

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet, _ := getGFData(DB, prjFilter, 1, empid)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(DB, prjFilter, empid)
				villageProjvisit, _ := getGFData(DB, prjFilter, 2, empid)
				obj["villagevisit"] = villageProjvisit
				beehiveProj, _ := getGFData(DB, prjFilter, 3, empid)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(DB, prjFilter, empid)

				if participantFilter != "" {
					projEnrolled = getParticipantFilterGfEnrolled(DB, prjFilter, empid, startDate, endDate)
				}
				obj["enroll"] = projEnrolled

				projectResult := DB.QueryRow("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id IN (?)", empID)
				var projectArray sql.NullString
				err = projectResult.Scan(&projectArray)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				// obj["greenMotivators"] = greenMotivators(con, startDate, endDate, projectArray, "", filterG)
				// obj["vyapar"] = Vyapar(con, startDate, endDate, projectArray, "", filterV)
				data = append(data, obj)
			}

			response := make(map[string]interface{})
			response["summary_circle_meet"] = circleMeet
			response["summary_circles"] = circleVisit
			response["summary_villagevisit"] = villageVisit
			response["summary_beehive"] = beehive
			response["summary_enroll"] = enrolled
			// response["summary_green"] = summaryGreen
			// response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}

	} else {
		w.WriteHeader(http.StatusCreated)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid role id"
		json.NewEncoder(w).Encode(response)
	}
}

func GelathiDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	data := []interface{}{}
	summaryTarget := 0
	summaryVillages := 0
	summaryActuals := 0
	summaryDay1 := 0
	summaryNoofbatches := 0

	summaryonselfshakthi := 0
	var day2Turnout float64

	// Parse the request body
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve values from the request
	roleID, err := getStringValue(request, "roleid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	empID, err := getStringValue(request, "emp_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	talukID, _ := getStringValue(request, "taluk")
	districtID, _ := getStringValue(request, "dist")
	startDate, _ := getStringValue(request, "start_date")
	endDate, _ := getStringValue(request, "end_date")
	funderID, _ := getStringValue(request, "funder_id")
	partnerID, _ := getStringValue(request, "partner_id")
	projectID, _ := getStringValue(request, "project_id")
	trainerID, _ := getStringValue(request, "trainer_id")
	opsManager, _ := getStringValue(request, "opsmanager")
	somID, _ := getStringValue(request, "somid")
	gflID, _ := getStringValue(request, "gflid")

	gflid, _ := strconv.Atoi(gflID)
	som, _ := strconv.Atoi(somID)
	opsid, _ := strconv.Atoi(opsManager)
	trainid, _ := strconv.Atoi(trainerID)
	projectid, _ := strconv.Atoi(projectID)
	funderid, _ := strconv.Atoi(funderID)
	talukid, _ := strconv.Atoi(talukID)
	distid, _ := strconv.Atoi(districtID)
	empid, _ := strconv.Atoi(empID)
	roleid, _ := strconv.Atoi(roleID)
	partnerid, _ := strconv.Atoi(partnerID)

	var isDateFilterApplied = false

	if startDate == "" || endDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 12 {
		filter := ""
		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}

		} else if roleid == 3 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				// var opsIds []int
				if som != 0 {
					opsIds, _ := getReportingOpsManagers(DB, som)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else if gflid != 0 {
					opsIds, _ := getSupervisor(DB, gflid)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					opsIds, _ := getReportingOpsManagers(DB, empid)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}

		} else if roleid == 12 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds, _ := getOpsManagers(DB, empid)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				}
			} else {
				filter = " and p.operations_manager in (0)"
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}
		} else if roleid == 4 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectIDs, _ := getOpProjects(DB, empid)
				if len(projectIDs) > 0 {
					filter = fmt.Sprintf(" AND p.operations_manager = %d", empid)
				} else {

					return
				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusOK)
				w.Write(js)
				return

			}
		} else if roleid == 9 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", empid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		}

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
			dateFilters = "date >= '" + startDate + "' and date <= '" + endDate + "'"

		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", partnerid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p INNER JOIN funder ON funder.funderID = p.funderID WHERE p.partnerID = %d AND %s %s", partnerid, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", partnerid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if distid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", distid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", talukid)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p INNER JOIN funder ON funder.funderID = p.funderID WHERE locationID = %d AND %s %s GROUP BY p.funderID", talukid, dateFilter, filter)
						filter += fmt.Sprintf(" AND locationID = %d", talukid)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)
						if err != nil {
							log.Println("SelfSakthiDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")

						w.Write(js)
						return

					}
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid locationid"
					js, err := json.Marshal(response)
					if err != nil {
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")

					w.Write(js)
					return

				}
			} else {
				getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + strconv.Itoa(distid)
				talukArray := []int{}
				talukRes, err := DB.Query(getTaluk)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				defer talukRes.Close()
				for talukRes.Next() {
					var tlk int
					err := talukRes.Scan(&tlk)
					if err != nil {
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					talukArray = append(talukArray, tlk)
				}
				funderListQuery = "SELECT p.funderID as id, funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ") and " + dateFilter + filter + " GROUP by p.funderID"
				filter += " and locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ")"
			}
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", funderid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT funderID AS id, funderName AS name FROM funder WHERE funderID = %d ", funderid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid funderid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if projectid == 0 && trainid == 0 && opsid == 0 && som == 0 && gflid == 0 && !isDateFilterApplied && roleid != 4 {
			// Role 4 OpsManager Default should be project list
			funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p " +
				"INNER JOIN funder ON p.funderId = funder.funderID " +
				"WHERE " + dateFilter + filter
		}

		if len(funderListQuery) > 0 {
			rows, err := DB.Query(funderListQuery)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var funderID int
				var funderName string
				err = rows.Scan(&funderID, &funderName)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				getProj := ""

				projectArray := make([]int, 0)

				if startDate != "" && endDate != "" {
					getProj = fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %d AND '%s' BETWEEN startDate AND endDate AND '%s' BETWEEN startDate AND endDate", funderID, startDate, endDate)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("ERROR>>", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						var startDate, endDate string
						err = projResult.Scan(&projectID, &startDate, &endDate)
						if err != nil {
							log.Println("ERROR>>", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
							return
						}

						projectArray = append(projectArray, projectID)
					}

				} else {
					getProj = fmt.Sprintf("SELECT id FROM project p WHERE funderID = %d and %s %s", funderID, dateFilter, filter)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("ERROR>>", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						err = projResult.Scan(&projectID)
						if err != nil {
							log.Println("ERROR>>", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
							return
						}

						projectArray = append(projectArray, projectID)
					}
				}

				if len(projectArray) == 0 {
					obj := make(map[string]interface{})
					obj["id"] = funderID
					obj["name"] = funderName
					obj["target"] = 0
					obj["actual"] = 0
					obj["day2"] = 0
					obj["noofselfshakthisurvey"] = 0
					obj["noofbatches"] = 0
					obj["villages"] = 0
					obj["startDate"] = ""
					obj["endDate"] = ""
					obj["select_type"] = "2"
					data = append(data, obj)
					continue
				}

				obj := make(map[string]interface{})
				obj["id"] = funderID
				obj["name"] = funderName
				obj["target"], _ = getTarget(DB, startDate, endDate, projectArray)
				summaryTarget += obj["target"].(int)
				obj["actual"], _ = getActual(DB, startDate, endDate, projectArray, "")
				summaryActuals += obj["actual"].(int)

				day1Count, _ := getDay1Count(DB, startDate, endDate, projectArray, "")
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout = float64(obj["actual"].(int)) / float64(day1Count)
					obj["day2"] = math.Round(day2Turnout * 100)
				} else {
					obj["day2"] = 0
				}
				obj["noofselfshakthisurvey"] = noOfselfshakthisurvey(DB, startDate, endDate, projectArray, "")
				summaryonselfshakthi += obj["noofselfshakthisurvey"].(int)
				obj["noofbatches"] = noofbatches(DB, startDate, endDate, projectArray, "")
				summaryNoofbatches += obj["noofbatches"].(int)
				obj["villages"], _ = newVillageCount(DB, startDate, endDate, projectArray, "")
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = ""
				obj["endDate"] = ""
				obj["select_type"] = "2"
				data = append(data, obj)
			}
		}

		var projectList string

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = fmt.Sprintf(" AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
				}

				projectList = fmt.Sprintf("SELECT id, projectName AS name, p.startDate, p.endDate FROM project p WHERE id = %d%s%s", projectid, filter, dateFilterNew)

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		} else if trainid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", trainid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
                                from tbl_poa tp
                                inner join project p on p.id = tp.project_id 
                                where user_id = %d and %s %s GROUP BY project_id`, trainid, dateFilter, filter)

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if opsid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", opsid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (startDate == "" && endDate == "") {
					projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where operations_manager = %d and %s %s GROUP BY id`, opsid, dateFilter, filter)
				} else {
					projectList = fmt.Sprintf(`SELECT p.id, p.projectName as name, p.startDate, p.endDate
                                from project p
                                join training_participants tp on p.id = tp.project_id
                                where p.operations_manager = %d and p.startDate BETWEEN %s and %s
                                GROUP BY p.id`, opsid, startDate, endDate)
				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if som > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", som)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where p.id in (
                                    SELECT DISTINCT project_id
                                    from tbl_poa tp
                                    where %s
                                ) and operations_manager in (
                                    SELECT id
                                    from employee e
                                    where e.supervisorId = %d
                                ) %s GROUP BY id`, dateFilters, som, filter)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if gflid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", gflid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT id, projectName as name, p.startDate, p.endDate
                                from project p
                                where operations_manager in (
                                    SELECT supervisorId
                                    from employee e
                                    where e.id = %d
                                ) and %s %s GROUP BY id`, gflid, dateFilter, filter)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0 || roleid == 4 && distid == 0 {
			projectList = fmt.Sprintf("SELECT id, projectName as name, p.startDate, p.endDate from project p where %s %s", dateFilter, filter)
		}

		if len(projectList) > 0 {
			rows, err := DB.Query(projectList)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var prList Project
				err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				tpFilter := ""
				tbFilter := ""

				projectArray := []int{prList.ID}

				if trainid > 0 {
					prList.Target, _ = getTrainerTarget(DB, trainid, projectArray)
					summaryTarget += prList.Target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					prList.Target, _ = getTarget(DB, startDate, endDate, projectArray)
					summaryTarget += prList.Target

				}
				prList.Actual, _ = getActual(DB, startDate, endDate, projectArray, tpFilter)
				summaryActuals += prList.Actual
				day1Count, _ := getDay1Count(DB, startDate, endDate, projectArray, tpFilter)
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout = float64(prList.Actual) / float64(day1Count)
					prList.Day2 = int(day2Turnout * 100)
				} else {
					prList.Day2 = 0
				}
				prList.Nosurvey = noOfselfshakthisurvey(DB, startDate, endDate, projectArray, "")
				summaryonselfshakthi += prList.Nosurvey
				prList.Noofbatches = noofbatches(DB, startDate, endDate, projectArray, "")
				summaryNoofbatches += prList.Noofbatches
				prList.Villages, _ = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += prList.Villages
				prList.SelectType = "1"

				data = append(data, prList)
			}
			if err := rows.Err(); err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
		}

		// Convert data to JSON and send response
		response := make(map[string]interface{})
		response["data"] = data
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals
		response["summary_Noofselfshakthisurvey"] = summaryonselfshakthi
		if summaryDay1 > 0 {
			day2Turnout = float64(summaryActuals) / float64(summaryDay1)
			response["summary_day2"] = int(math.Round(day2Turnout * 100))
		} else {
			response["summary_day2"] = 0
		}

		response["summary_NoofBatches"] = summaryNoofbatches
		response["summary_villages"] = summaryVillages

		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	} else if roleid == 5 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			var dateFilter string
			if isDateFilterApplied {
				dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			query := fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
			FROM tbl_poa tp
			INNER JOIN project p ON p.id = tp.project_id
			WHERE user_id = %d%s GROUP BY project_id`, empid, dateFilter)

			if projectid > 0 {
				query = fmt.Sprintf(`SELECT project_id as id, projectName as name, p.startDate, p.endDate
				FROM tbl_poa tp
				INNER JOIN project p ON p.id = tp.project_id
				WHERE user_id = %d AND tp.project_id = %d GROUP BY tp.project_id`, empid, projectid)
			}

			rows, err := DB.Query(query)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			var data []Project

			for rows.Next() {
				var prj Project
				err := rows.Scan(&prj.ID, &prj.Name, &prj.StartDate, &prj.EndDate)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				projectArray := []int{prj.ID}
				prj.Target, _ = getTrainerTarget(DB, empid, projectArray)
				summaryTarget += prj.Target
				filter := fmt.Sprintf(" and tp.trainer_id = %d", empid)

				prj.Actual, _ = getActual(DB, startDate, endDate, projectArray, filter)
				summaryActuals += prj.Actual
				day1Count, _ := getDay1Count(DB, startDate, endDate, projectArray, filter)
				summaryDay1 += day1Count
				if day1Count > 0 {
					day2Turnout := float64(prj.Actual) / float64(day1Count)
					prj.Day2 = int(day2Turnout * 100)
				} else {
					prj.Day2 = 0
				}

				prj.Noofbatches = getTrainerBatches(DB, startDate, endDate, projectArray, empID)
				summaryNoofbatches += prj.Noofbatches

				prj.Villages = TrainerVillageCount(DB, startDate, endDate, projectArray, empID)
				summaryVillages += prj.Villages
				prj.SelectType = "1"

				prj.Nosurvey = noofsurveyssTrainer(DB, startDate, endDate, projectArray, empID)
				summaryonselfshakthi += prj.Nosurvey
				data = append(data, prj)
			}

			response := make(map[string]interface{})
			response["data"] = data
			response["summary_target"] = summaryTarget
			response["summary_actual"] = summaryActuals
			response["summary_Noofselfshakthisurvey"] = summaryonselfshakthi
			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				response["summary_day2"] = 0
			}

			response["summary_NoofBatches"] = summaryNoofbatches
			response["summary_villages"] = summaryVillages

			// Convert response to JSON
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}
	} else if roleid == 13 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {

			// filter := ""
			filters := ""
			// filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			// filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			var participantFilter = false

			if startDate != "" && endDate != "" {
				// filter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s'", startDate, endDate)
				filters = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s'", startDate, endDate)
				// filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate BETWEEN '%s' and '%s'", startDate, endDate)
				// filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate BETWEEN '%s' and '%s'", startDate, endDate)
			}

			if projectID != "" {
				filters += fmt.Sprintf(" and p.id = %s", projectID)
			}

			f := ""
			if gfID, ok := request["gfId"].(string); ok && gfID != "" {
				f = fmt.Sprintf(" and id=%s", gfID)
			}

			var (
				summaryCircleMeet   int
				summaryCircles      int
				summaryVillageVisit int
				summaryBeehive      int
				summaryEnroll       int
			)

			employeeIDs := make([]int, 0)
			rows, err := DB.Query("SELECT id FROM employee WHERE status = 1 AND supervisorId = ?"+f, empID)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()
			for rows.Next() {
				var id int
				err := rows.Scan(&id)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}
				employeeIDs = append(employeeIDs, id)
			}
			err = rows.Err()
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}

			getProjs := fmt.Sprintf(`SELECT project_id AS id, p.projectName AS name FROM tbl_poa tp
			INNER JOIN project p ON p.id = tp.project_id
			WHERE p.gfl_id = %d %s GROUP BY tp.project_id`, empid, filters)

			rows, err = DB.Query(getProjs)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var id int
				var name string
				err := rows.Scan(&id, &name)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				obj := make(map[string]interface{})
				obj["id"] = id
				obj["name"] = name

				prjFilter := ""
				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				obj["circlemeet"] = getGFDataN(DB, prjFilter, 1, employeeIDs)
				summaryCircleMeet += obj["circlemeet"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and p.endDate BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				obj["circle"] = getGFCircleN(DB, prjFilter, employeeIDs)
				summaryCircles += obj["circle"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				}
				obj["villagevisit"] = getGFDataN(DB, prjFilter, 2, employeeIDs)
				summaryVillageVisit += obj["villagevisit"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)

				}
				obj["beehive"] = getGFDataN(DB, prjFilter, 3, employeeIDs)
				summaryBeehive += obj["beehive"].(int)

				if startDate != "" && endDate != "" {
					prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
				} else {
					prjFilter = fmt.Sprintf(" and p.id = %d", id)
				}

				projEnrolled := getGfEnrolledN(DB, prjFilter, employeeIDs)

				if participantFilter {
					if startDate != "" && endDate != "" {
						prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", startDate, endDate, id)
						projEnrolled = getParticipantFilterGfEnrolledN(DB, prjFilter, employeeIDs, startDate, endDate)
					}
				}

				obj["enroll"] = projEnrolled
				summaryEnroll += obj["enroll"].(int)

				projectArray := ""
				// project_result := db.QueryRow("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id IN (?)", empId)
				project_result := DB.QueryRow("SELECT GROUP_CONCAT(DISTINCT id) as ids FROM project prj WHERE gfl_id IN (?)", empID)
				err = project_result.Scan(&projectArray)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				data = append(data, obj)
			}
			response := make(map[string]interface{})
			response["data"] = data
			response["summary_circle_meet"] = summaryCircleMeet
			response["summary_circles"] = summaryCircles
			response["summary_villagevisit"] = summaryVillageVisit

			response["summary_beehive"] = summaryBeehive
			response["summary_enroll"] = summaryEnroll

			// Convert response to JSON
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}
	} else if roleid == 6 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", empID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			participantFilter := ""
			filter := ""
			filters := ""
			// filterG := ""
			// filterV := ""
			if projectid > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", projectid)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
					filters = fmt.Sprintf(" and tp.enroll_date >= '%s' and tp.enroll_date <= '%s'", startDate, endDate)
					// filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", startDate, endDate)
					// filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", startDate, endDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					// filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					// filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet, _ := getGFData(DB, filter, 1, empid)
			villageVisit, _ := getGFData(DB, filter, 2, empid)
			beehive, _ := getGFData(DB, filter, 3, empid)
			enrolled, _ := getGfEnrolled(DB, filters, empid)
			circleVisit := getGFCircle(DB, filter, empid)

			data := make([]map[string]interface{}, 0)

			getProjs := fmt.Sprintf(`SELECT project_id as id, p.projectName as name FROM tbl_poa tp
						INNER JOIN project p ON p.id = tp.project_id
						WHERE tp.user_id = %d %s GROUP BY tp.project_id
						UNION
						SELECT tp.project_id as id, p.projectName as name FROM training_participants tp
						INNER JOIN project p ON tp.project_id = p.id
						WHERE enroll = 1 AND gelathi_id = %d %s`, empid, filter, empid, filter)

			if projectid > 0 {
				getProjs = fmt.Sprintf(`SELECT project_id as id, p.projectName as name FROM tbl_poa tp
						INNER JOIN project p ON p.id = tp.project_id
						WHERE tp.project_id = %d AND tp.user_id = %d GROUP BY tp.project_id
						UNION
						SELECT tp.project_id as id, p.projectName as name FROM training_participants tp
						INNER JOIN project p ON tp.project_id = p.id
						WHERE enroll = 1 AND gelathi_id = %d AND tp.project_id = %d`, projectid, empid, empid, projectid)
			}

			projectsList, err := DB.Query(getProjs)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			defer projectsList.Close()

			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				obj := make(map[string]interface{})
				obj["id"] = id

				obj["name"] = name

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet, _ := getGFData(DB, prjFilter, 1, empid)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(DB, prjFilter, empid)
				villageProjvisit, _ := getGFData(DB, prjFilter, 2, empid)
				obj["villagevisit"] = villageProjvisit
				beehiveProj, _ := getGFData(DB, prjFilter, 3, empid)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(DB, prjFilter, empid)

				if participantFilter != "" {
					projEnrolled = getParticipantFilterGfEnrolled(DB, prjFilter, empid, startDate, endDate)
				}
				obj["enroll"] = projEnrolled

				projectResult := DB.QueryRow("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id IN (?)", empID)
				var projectArray sql.NullString
				err = projectResult.Scan(&projectArray)
				if err != nil {
					log.Println("ERROR>>", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
					return
				}

				// obj["greenMotivators"] = greenMotivators(con, startDate, endDate, projectArray, "", filterG)
				// obj["vyapar"] = Vyapar(con, startDate, endDate, projectArray, "", filterV)
				data = append(data, obj)
			}

			response := make(map[string]interface{})
			response["summary_circle_meet"] = circleMeet
			response["summary_circles"] = circleVisit
			response["summary_villagevisit"] = villageVisit
			response["summary_beehive"] = beehive
			response["summary_enroll"] = enrolled
			// response["summary_green"] = summaryGreen
			// response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		}

	} else {
		w.WriteHeader(http.StatusCreated)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid role id"
		json.NewEncoder(w).Encode(response)
	}
}

func getStringValue(request map[string]interface{}, key string) (string, error) {
	value, ok := request[key].(string)
	if !ok {
		return "", fmt.Errorf("invalid value for %s", key)
	}
	return value, nil
}
