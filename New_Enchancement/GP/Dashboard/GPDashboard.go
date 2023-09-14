package green

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Project struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	StartDate                string `json:"start_date"`
	EndDate                  string `json:"end_date"`
	Target                   int    `json:"target"`
	Actual                   int    `json:"actual"`
	NoOfgreensurvey          int    `json:"nofgreensurvey"`
	Villages                 int    `json:"villages"`
	SelectType               string `json:"select_type"`
	NoOfGreenCoharts         int    `json:"noofgreencoharts"`
	NoofGreenModuleCompleted int    `json:"noofgreenmodulecompleted"`
	GreenEnrolled            int    `json:"greenenrolled"`
}

type Response struct {
	SummaryTarget                   int       `json:"summary_target"`
	SummaryActual                   int       `json:"summary_actual"`
	SummaryVillages                 int       `json:"summary_villages"`
	SummaryNoofGreensurvey          int       `json:"summary_NoofGreenrsurvey"`
	SummaryNoofGreencoharts         int       `json:"summary_NoofGreencoharts"`
	SummaryNoofGreenmodulecomoleted int       `json:"summary_NoofGreenmodulecomoleted"`
	SummaryGreenenrolled            int       `json:"summary_Greenenrolled"`
	Data                            []Project `json:"data"`
	Code                            int       `json:"code"`
	Success                         bool      `json:"success"`
	Message                         string    `json:"message"`
}

func GreenDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
	gflID, _ := getStringValue(request, "gflid")
	opsManager, _ := getStringValue(request, "opsmanager")
	somID, _ := getStringValue(request, "somid")
	gfID, _ := getStringValue(request, "gfid")

	gfid, _ := strconv.Atoi(gfID)
	som, _ := strconv.Atoi(somID)
	opsid, _ := strconv.Atoi(opsManager)
	gflid, _ := strconv.Atoi(gflID)
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

	data := []interface{}{}

	summaryVillages := 0
	summaryTarget := 0
	summaryActuals := 0

	summaryGreenrenolled := 0
	summaryNoofGreensurvey := 0
	summaryNoofGreenmodulecompleted := 0
	summayNoofGreencoharts := 0

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 11 || roleid == 12 {
		filter := ""
		summaryFilter := ""
		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", empid)
			if err != nil {
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
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
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", talukid)
					if err != nil {
						fmt.Println(err)
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
						//fmt.Println(response)
						if err != nil {
							log.Println("GreenDashboard", err)
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
					//fmt.Println(response)
					if err != nil {
						log.Println("GreenDashboard", err)
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
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				defer talukRes.Close()
				for talukRes.Next() {
					var tlk int
					err := talukRes.Scan(&tlk)
					if err != nil {
						log.Println("GreenDashboard", err)
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
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT funderID AS id, funderName AS name FROM funder WHERE funderID = %d ", funderid)
				summaryFilter = fmt.Sprintf(" AND p.funderID = %d", funderid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid funderid"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if projectid == 0 && gfid == 0 && opsid == 0 && som == 0 && gflid == 0 && !isDateFilterApplied && roleid != 4 {
			// Role 4 OpsManager Default should be project list
			funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p " +
				"INNER JOIN funder ON p.funderId = funder.funderID " +
				"WHERE " + dateFilter + filter
		}

		if len(funderListQuery) > 0 {
			rows, err := DB.Query(funderListQuery)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var funderID int
				var funderName string
				err = rows.Scan(&funderID, &funderName)
				if err != nil {
					log.Fatal(err)
				}
				getProj := ""

				projectArray := make([]int, 0)

				if startDate != "" && endDate != "" {
					getProj = fmt.Sprintf("SELECT distinct p.id, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and p.funderID = %d AND '%s' BETWEEN p.startDate AND p.endDate AND '%s' BETWEEN p.startDate AND p.endDate", funderID, startDate, endDate)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Fatal(err)
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						var startDate, endDate string
						err = projResult.Scan(&projectID, &startDate, &endDate)
						if err != nil {
							log.Fatal(err)
						}

						fmt.Println(projectID)

						projectArray = append(projectArray, projectID)
					}

				} else {
					getProj = fmt.Sprintf("SELECT distinct p.id from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and p.funderID = %d and %s %s", funderID, dateFilter, filter)
					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Fatal(err)
					}
					defer projResult.Close()

					for projResult.Next() {
						var projectID int
						err = projResult.Scan(&projectID)
						if err != nil {
							log.Fatal(err)
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
					obj["greenenrolled"] = 0
					obj["nofgreensurvey"] = 0
					obj["noofgreencoharts"] = 0
					obj["noofgreenmodulecompleted"] = 0
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
				obj["actual"] = greenenrolled(DB, startDate, endDate, projectArray, "")
				summaryActuals += obj["actual"].(int)
				obj["greenenrolled"] = greenenrolled(DB, startDate, endDate, projectArray, "")
				summaryNoofGreenmodulecompleted += obj["greenenrolled"].(int)
				obj["noofgreenmodulecompleted"] = noofgreenmodulecompleted(DB, startDate, endDate, projectArray, "")
				summaryNoofGreenmodulecompleted += obj["noofgreenmodulecompleted"].(int)
				obj["nofgreensurvey"] = noofgreensurvey(DB, startDate, endDate, projectArray, "")
				summaryNoofGreensurvey += obj["nofgreensurvey"].(int)
				obj["noofgreencoharts"] = noofgreencoharts(DB, startDate, endDate, projectArray, "")
				summayNoofGreencoharts += obj["noofgreencoharts"].(int)
				obj["villages"], _ = GreennewVillageCount(DB, startDate, endDate, projectArray, "")
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = ""
				obj["endDate"] = ""
				obj["select_type"] = "2"
				data = append(data, obj)
			}
		}
		fmt.Println(summaryFilter)

		var projectList string

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectid)
			if err != nil {
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = fmt.Sprintf(" AND t.GreenMotivatorsDate >= '%s' and '%s'", startDate, endDate)
				}

				projectList = fmt.Sprintf("SELECT distinct p.id, p.projectName AS name, p.startDate, p.endDate FROM project p join training_participants t on p.id=t.project_id where (t.GreenMotivators=1 or t.new_green=1) and id = %d%s%s", projectid, filter, dateFilterNew)

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		} else if gfid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id = ?", gfid)
			if err != nil {
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and t.user_id = %d and %s GROUP BY t.project_id`, gfid, dateFilter)
				summaryFilter = fmt.Sprintf(" and t.user_id = %d", gfid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (startDate == "" && endDate == "") {
					projectList = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) AND p.operations_manager = %d and %s %s GROUP BY id`, opsid, dateFilter, filter)
				} else {
					projectList = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) AND p.operations_manager = %d and p.startDate BETWEEN %s and %s
					GROUP BY p.id`, opsid, startDate, endDate)
				}
				summaryFilter = fmt.Sprintf(" and p.operations_manager = %d", opsid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate
				from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id 
				where (tp.GreenMotivators=1 or tp.new_green=1) and p.id in (
					SELECT DISTINCT project_id
					from tbl_poa tp
					where %s
                                ) and operations_manager in (
                                    SELECT id
                                    from employee e
                                    where e.supervisorId = %d
                                ) %s GROUP BY id`, dateFilters, som, filter)
				summaryFilter = fmt.Sprintf(" and p.operations_manager in (SELECT id from employee e where e.supervisorId = %d)", som)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
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
				fmt.Println(err)
			}
			defer rows.Close()
			if rows.Next() {
				projectList = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and operations_manager in (
					SELECT supervisorId
					from employee e
					where e.id = %d
				) and %s %s GROUP BY id`, gflid, dateFilter, filter)
				summaryFilter = fmt.Sprintf(" and p.operations_manager in (SELECT supervisorId from employee e where e.id = %d)", gflid)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)
				//fmt.Println(response)
				if err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0 || roleid == 4 && distid == 0 {
			projectList = fmt.Sprintf("SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and %s %s", dateFilter, filter)
		}

		fmt.Println(projectList)
		if len(projectList) > 0 {
			rows, err := DB.Query(projectList)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var prList Project
				err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
				if err != nil {
					log.Fatal(err)
				}
				tpFilter := ""
				tbFilter := ""
				projectArray := []int{prList.ID}
				if gfid > 0 {
					prList.Target, _ = getTrainerTarget(DB, gfid, projectArray)
					summaryTarget += prList.Target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", gfid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", gfid)
				} else {
					prList.Target, _ = getTarget(DB, startDate, endDate, projectArray)
					summaryTarget += prList.Target

				}
				prList.Actual = greenenrolled(DB, startDate, endDate, projectArray, tpFilter)
				summaryActuals += prList.Actual
				prList.NoOfgreensurvey = noofgreensurvey(DB, startDate, endDate, projectArray, "")
				summaryNoofGreensurvey += prList.NoOfgreensurvey
				prList.NoofGreenModuleCompleted = noofgreenmodulecompleted(DB, startDate, endDate, projectArray, "")
				summaryNoofGreenmodulecompleted += prList.NoofGreenModuleCompleted
				prList.GreenEnrolled = greenenrolled(DB, startDate, endDate, projectArray, "")
				summaryGreenrenolled += prList.GreenEnrolled
				prList.NoOfGreenCoharts = noofgreencoharts(DB, startDate, endDate, projectArray, "")
				summayNoofGreencoharts += prList.NoOfGreenCoharts
				prList.Villages, _ = GreennewVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += prList.Villages
				prList.SelectType = "1"

				data = append(data, prList)
			}
			if err := rows.Err(); err != nil {
				log.Fatal(err)
			}
		}
		var projects []Project
		for _, item := range data {
			project, ok := item.(Project)
			if ok {
				projects = append(projects, project)
			}
		}

		fmt.Println(dateFilters)

		response := make(map[string]interface{})
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals

		response["summary_NoofGreencoharts"] = summayNoofGreencoharts
		response["summary_NoofGreenrsurvey"] = summaryNoofGreensurvey
		filter += summaryFilter
		response["summary_villages"] = summaryVillages

		response["summary_NoofGreenmodulecomoleted"] = summaryNoofGreenmodulecompleted
		response["summary_Greenenrolled"] = summaryGreenrenolled
		response["data"] = data
		response["code"] = 200
		response["success"] = true
		response["message"] = "Successfully"

		js, err := json.Marshal(response)
		if err != nil {
			log.Println("GreenDashboard", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return

	} else if roleid == 6 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", empID)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		if rows.Next() {
			filter := ""
			filters := ""
			if isDateFilterApplied {
				filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
				filters = fmt.Sprintf(" and GreenMotivatorsDate >= '%s' and GreenMotivatorsDate <= '%s'", startDate, endDate)
			} else {
				filter = " and p.endDate >= CURRENT_DATE()"
				filters = " and GreenMotivatorsDate >= CURRENT_DATE()"
			}
			fmt.Println(filters)

			getProjs := fmt.Sprintf(`SELECT project_id as id, p.projectName as name,p.startDate,p.endDate FROM tbl_poa tp
							INNER JOIN project p ON p.id = tp.project_id
							WHERE tp.user_id = %d %s GROUP BY tp.project_id
							UNION
							SELECT tp.project_id as id, p.projectName as name,p.startDate,p.endDate FROM training_participants tp
							INNER JOIN project p ON tp.project_id = p.id
							WHERE (GreenMotivators = 1 or new_green=1) AND gelathi_id = %d %s`, empid, filter, empid, filter)

			if projectid > 0 {
				getProjs = fmt.Sprintf(`SELECT project_id as id, p.projectName as name,p.startDate,p.endDate FROM tbl_poa tp
							INNER JOIN project p ON p.id = tp.project_id
							WHERE tp.project_id = %d AND tp.user_id = %d GROUP BY tp.project_id
							UNION
							SELECT tp.project_id as id, p.projectName as name,p.startDate,p.endDate FROM training_participants tp
							INNER JOIN project p ON tp.project_id = p.id
							WHERE (GreenMotivators = 1 or new_green=1) AND gelathi_id = %d AND tp.project_id = %d`, projectid, empid, empid, projectid)
			}
			fmt.Println(getProjs)
			projectsList, err := DB.Query(getProjs)
			if err != nil {
				log.Fatal(err)
			}
			defer projectsList.Close()
			var data []Project
			for projectsList.Next() {
				var prList Project
				err := projectsList.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
				if err != nil {
					log.Fatal(err)
				}

				projectArray := []int{prList.ID}
				prList.Target, _ = getTrainerTarget(DB, empid, projectArray)
				summaryTarget += prList.Target
				filter := fmt.Sprintf(" and gelathi_id = %d", empid)
				coharts := fmt.Sprintf(" and tbl_poa.user_id = %d", empid)
				prList.Actual = greenenrolled(DB, startDate, endDate, projectArray, filter)
				summaryActuals += prList.Actual
				prList.NoOfgreensurvey = noofgreensurvey(DB, startDate, endDate, projectArray, filter)
				summaryNoofGreensurvey += prList.NoOfgreensurvey
				prList.NoofGreenModuleCompleted = noofgreenmodulecompleted(DB, startDate, endDate, projectArray, filter)
				summaryNoofGreenmodulecompleted += prList.NoofGreenModuleCompleted
				prList.GreenEnrolled = greenenrolled(DB, startDate, endDate, projectArray, filter)
				summaryGreenrenolled += prList.GreenEnrolled
				prList.NoOfGreenCoharts = noofgreencoharts(DB, startDate, endDate, projectArray, coharts)
				summayNoofGreencoharts += prList.NoOfGreenCoharts
				prList.Villages = GfGreenVillageCount(DB, startDate, endDate, projectArray, empID)
				summaryVillages += prList.Villages
				prList.SelectType = "1"
				data = append(data, prList)
			}

			response := make(map[string]interface{})
			response["data"] = data
			response["summary_target"] = summaryTarget
			response["summary_actual"] = summaryActuals
			response["summary_Greensurvey"] = summaryNoofGreensurvey
			response["summary_Greencoharts"] = summayNoofGreencoharts
			response["summary_villages"] = summaryVillages
			response["summary_noofGreenmodulecompleted"] = summaryNoofGreenmodulecompleted
			response["summary_Greenenrolled"] = summaryGreenrenolled
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
			////fmt.Println(response)
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
			fmt.Println(err)
		}
		defer rows.Close()
		if rows.Next() {
			var filter string
			if isDateFilterApplied {
				filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", startDate, endDate)
			} else {
				filter = " and p.endDate >= CURRENT_DATE()"
			}

			filters := ""
			gf := ""
			if projectid != 0 {
				filters += fmt.Sprintf(" and p.id = %d", projectid)
			}
			if gfid != 0 {
				gf += fmt.Sprintf(" and tp.gelathi_id = %d", gfid)
			}

			getProjs := fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and p.gfl_id = %d %s GROUP BY tp.project_id`, empid, filter)
			if projectid > 0 {
				getProjs = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and p.gfl_id = %d %s GROUP BY tp.project_id`, empid, filters)
			}
			if gfid > 0 {
				getProjs = fmt.Sprintf(`SELECT distinct p.id, projectName as name, p.startDate, p.endDate from project p join tbl_poa t on t.project_id = p.id join training_participants tp on tp.tb_id=t.tb_id where (tp.GreenMotivators=1 or tp.new_green=1) and p.gfl_id = %d %s %s GROUP BY tp.project_id`, empid, gf, filter)
			}

			rows, err = DB.Query(getProjs)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var prList Project
				err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
				if err != nil {
					log.Fatal(err)
				}
				projectArray := []int{prList.ID}
				prList.Target, _ = getTrainerTarget(DB, empid, projectArray)
				summaryTarget += prList.Target
				filter := fmt.Sprintf(" and gelathi_id = %d", gfid)
				coharts := fmt.Sprintf(" and tbl_poa.user_id = %d", gfid)
				village := fmt.Sprintf(" and tp.user_id = %d", gfid)
				prList.Actual = Gflgreenenrolled(DB, startDate, endDate, projectArray, empid, filter)
				summaryActuals += prList.Actual
				prList.NoOfgreensurvey = Gflgreenenrolled(DB, startDate, endDate, projectArray, empid, filter)
				summaryNoofGreensurvey += prList.NoOfgreensurvey
				prList.NoofGreenModuleCompleted = Gflnoofgreenmodulecompleted(DB, startDate, endDate, projectArray, empid, filter)
				summaryNoofGreenmodulecompleted += prList.NoofGreenModuleCompleted
				prList.GreenEnrolled = Gflgreenenrolled(DB, startDate, endDate, projectArray, empid, filter)
				summaryGreenrenolled += prList.GreenEnrolled
				prList.NoOfGreenCoharts = Gflnoofgreencoharts(DB, startDate, endDate, projectArray, empid, coharts)
				summayNoofGreencoharts += prList.NoOfGreenCoharts
				prList.Villages, _ = GflGreennewVillageCount(DB, startDate, endDate, projectArray, empid, village)
				summaryVillages += prList.Villages
				prList.SelectType = "1"
				data = append(data, prList)
			}
			response := make(map[string]interface{})
			response["data"] = data
			response["summary_target"] = summaryTarget
			response["summary_actual"] = summaryActuals
			response["summary_noofVyaparsurvey"] = summaryNoofGreensurvey
			response["summary_noofvyaparcoharts"] = summayNoofGreencoharts
			response["summary_villages"] = summaryVillages
			response["summary_noofvyaparmodulecomoleted"] = summaryNoofGreenmodulecompleted
			response["summary_vyaparenrolled"] = summaryGreenrenolled
			response["success"] = true
			response["message"] = "Successfully"

			// Convert response to JSON
			jsonData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
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
