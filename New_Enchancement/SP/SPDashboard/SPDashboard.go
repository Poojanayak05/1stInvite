package spoorthi

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
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Target               int    `json:"target"`
	Actual               int    `json:"actual"`
	NoOfsporthisurvey    int    `json:"Noofsporthisurvey"`
	Villages             int    `json:"villages"`
	SelectType           string `json:"select_type"`
	NoofCircleMeeting    int    `json:"NoofCircleMeeting"`
	Noofsporthicompleted int    `json:"Noofsporthicompleted"`
	Gelathienrolled      int    `json:"Gelathienrolled"`
	Noofbeehives         int    `json:"Noofbeehives"`
}

type Response struct {
	SummaryTarget               int       `json:"summary_target"`
	SummaryActual               int       `json:"summary_actual"`
	SummaryVillages             int       `json:"summary_villages"`
	SummaryNoofSporthisurvey    int       `json:"summary_sporthisurvey"`
	SummaryNoofCircleMeeting    int       `json:"summary_NoofCircleMeeting"`
	SummaryNoofsporthicompleted int       `json:"summary_Noofsporthicompleted"`
	SummaryGelathienrolled      int       `json:"summary_Gelathienrolled"`
	SummaryNoofbeehives         int       `json:"summary_Noofbeehives"`
	Data                        []Project `json:"data"`
	Code                        int       `json:"code"`
	Success                     bool      `json:"success"`
	Message                     string    `json:"message"`
}

func GelathiProgramDashboard1(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
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
	summaryTarget := 0
	summaryVillages := 0
	summaryActuals := 0
	summarycirclemeeting := 0
	summarysporthisurvey := 0
	summarysporthicompleted := 0
	summaryGelathiEnrolled := 0
	summaryofbeehives := 0

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 11 || roleid == 12 {
		filter := ""
		summaryFilter := ""

		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", empID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", empID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				// var opsIds []int
				if som != 0 {
					opsIds := getReportingOpsManagers(DB, som)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else if gfid != 0 {
					opsIds := getSupervisor(DB, gfid)
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					opsIds := getReportingOpsManagers(DB, empid)
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
					log.Println("GelathiProgramDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", empID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds := GetOpsManagers(DB, empid)
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
					log.Println("GelathiProgramDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", empID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				// Ops Manager
				projectIds := getOpProjects(DB, empid)
				if len(projectIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager = %d", empid)
				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid employe id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", empID)
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
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		}
		circleMeetCountQuery := ""
		var circleMeetCount int
		var noofsporthimodulecompleted, noofsporthisurvey, noofenrollgelathi int

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
			dateFilters = "date >= '" + startDate + "' and date <= '" + endDate + "'"
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		// CEO Dashboard
		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", partnerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id ,funderName as name FROM project p inner join funder on funder.funderID = p.funderID where p.partnerID = " + (partnerID) + " and " + dateFilter + filter
				filter += " and p.partnerID = " + (partnerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
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
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", talukID)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = "SELECT p.funderID as id,funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID = " + (talukID) + " and " + dateFilter + filter + " GROUP by p.funderID"
						filter += " and locationID = " + (talukID)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)
						if err != nil {
							log.Println("GelathiProgramDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")

						w.Write(js)
						return

					}
				} else {
					// get taluk of specified dist
					getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + (districtID)
					talukArray := []int{}
					talukRes, err := DB.Query(getTaluk)
					if err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer talukRes.Close()
					for talukRes.Next() {
						var tlk int
						err := talukRes.Scan(&tlk)
						if err != nil {
							log.Println("GelathiProgramDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						talukArray = append(talukArray, tlk)
					}
					funderListQuery = "SELECT p.funderID as id, funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ") and " + dateFilter + filter + " GROUP by p.funderID"
					filter += " and locationID in (" + strings.Trim(strings.Replace(fmt.Sprint(talukArray), " ", ",", -1), "[]") + ")"
				}
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid locationid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", funderID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder f where funderID = %d", funderid)
				summaryFilter = fmt.Sprintf(" and p.funderID = %d", funderid)

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid funderid"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if partnerid == 0 && gflid == 0 && opsid == 0 && som == 0 && gfid == 0 && !isDateFilterApplied && roleid != 4 {
			// role 4 OpsManager Default should be project list
			funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) as id,funderName as name FROM project p inner join funder on p.funderId = funder.funderID where %s%s", dateFilter, filter)
		}
		if len(funderListQuery) > 0 {
			if projectID == "" {

				res, err := DB.Query(funderListQuery)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				defer res.Close()

				for res.Next() {
					var projectArray []int
					var funderId int
					var funderName string

					if err := res.Scan(&funderId, &funderName); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					circleMeetCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS NoofCircleMeeting
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (4, 5, 6, 7, 8, 9)
										WHERE p.funderId IS NOT NULL
											AND p.funderId = ?
										GROUP BY p.funderId`

					if err := DB.QueryRow(circleMeetCountQuery, funderId).Scan(&circleMeetCount); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					getmodulecompletedquery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE (t2.module1=1 and t2.module2=1 and t2.module3=1 and t2.module4=1 and t2.module5=1)
				AND t3.funderId = ?`

					if err := DB.QueryRow(getmodulecompletedquery, funderId).Scan(&noofsporthimodulecompleted); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					sporthisurvey := `SELECT count(t2.id) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.funderId = ?`
					if err := DB.QueryRow(sporthisurvey, funderId).Scan(&noofsporthisurvey); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					enrollgelathi := `SELECT COUNT(tp.id) as gelathiCount FROM training_participants tp INNER JOIN project p ON p.id = tp.enrolledProject WHERE enroll = 1 and p.funderID = ?`

					if err := DB.QueryRow(enrollgelathi, funderId).Scan(&noofenrollgelathi); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					var noofgelathibeehive int
					gelathibeehive := `SELECT count(tbl_poa.id) as noofbeehvees FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id   where tbl_poa.type = 2 AND tbl_poa.session_type = 3 and project.funderID = ?`
					if err := DB.QueryRow(gelathibeehive, funderId).Scan(&noofgelathibeehive); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					var villages, subvillages int
					newvillages := `SELECT COUNT(distinct tp.location_id) as villages FROM bdms_staff.tbl_poa tp INNER JOIN training_participants tr on tr.tb_id=tp.tb_id join project p on p.id=tr.project_id where (enroll=1) and p.funderID=?`
					if err := DB.QueryRow(newvillages, funderId).Scan(&villages); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					subvillages1 := `SELECT COUNT(distinct tp.sub_village) as sub_villages FROM bdms_staff.tbl_poa tp INNER JOIN training_participants tr on tr.tb_id=tp.tb_id join project p on p.id=tr.project_id where (enroll=1) and p.funderID=?`
					if err := DB.QueryRow(subvillages1, funderId).Scan(&subvillages); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					var noofactual int
					actual := `select count(tp.id) as actual from training_participants tp join project p on p.id=tp.project_id where tp.day2 = 1 and tp.enroll=1 and p.funderID=?`
					if err := DB.QueryRow(actual, funderId).Scan(&noofactual); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
					if startDate != "" && endDate != "" {
						getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + startDate + "' BETWEEN startDate and endDate and '" + endDate + "' BETWEEN startDate and endDate"
					}

					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var id int
						if err := projResult.Scan(&id); err != nil {
							log.Println("GelathiProgramDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						projectArray = append(projectArray, id)
					}

					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":                   funderId,
							"name":                 funderName,
							"target":               0,
							"actual":               0,
							"NoofCircleMeeting":    0,
							"villages":             0,
							"Gelathienrolled":      0,
							"Noofsporthisurvey":    0,
							"Noofsporthicompleted": 0,
							"Noofbeehives":         0,
							"start_date":           "",
							"end_date":             "",
							"select_type":          "2",
						}
						data = append(data, obj)
						continue
					}
					// str := strconv.Itoa(reqBody.FunderId)
					stringSlice := make([]string, len(projectArray))

					var obj map[string]interface{}
					for i, val := range projectArray {
						stringSlice[i] = strconv.Itoa(val)

						obj = map[string]interface{}{
							"id":                   funderId,
							"name":                 funderName,
							"target":               getTarget(DB, startDate, endDate, projectArray),
							"actual":               noofactual,
							"NoofCircleMeeting":    circleMeetCount,
							"villages":             villages + subvillages, // New village count function anas
							"Gelathienrolled":      noofenrollgelathi,
							"Noofsporthisurvey":    noofsporthisurvey,
							"Noofsporthicompleted": noofsporthimodulecompleted,
							"Noofbeehives":         noofgelathibeehive,
							"start_date":           "",
							"end_date":             "",
							"select_type":          "2",
						}
					}
					summaryTarget += obj["target"].(int)
					summaryActuals += obj["actual"].(int)
					summaryVillages += (villages + subvillages)
					summaryGelathiEnrolled += obj["Gelathienrolled"].(int)
					summarycirclemeeting += obj["NoofCircleMeeting"].(int)
					summarysporthisurvey += obj["Noofsporthisurvey"].(int)
					summarysporthicompleted += obj["Noofsporthicompleted"].(int)
					summaryofbeehives += noofgelathibeehive

					data = append(data, obj)
				}
			}
		}
		projectList := ""

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + (projectID) + filter + dateFilterNew
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if gfid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id = ?", gfID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + (gfID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + (gfID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if opsid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", opsManager)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (startDate == "" && endDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + (opsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + (opsManager) + " and tp.participant_day2 >= '" + startDate + "' and tp.participant_day2 <= '" + endDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + (opsManager)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if som > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", somID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + (somID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + (somID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if gflid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", gflID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + (gflID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + (gflID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if (isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0) || (roleid == 4 && distid == 0) {
			//role 4 - OpsManager Default should be project list without location filter
			projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
		}

		if len(projectList) > 0 {
			res, err := DB.Query(projectList)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			defer res.Close()
			var projectArray []int
			for res.Next() {
				var obj = make(map[string]interface{})

				var id int
				var name string
				var startDate1 string
				var endDate1 string

				err := res.Scan(&id, &name, &startDate1, &endDate1)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}

				obj["id"] = id
				obj["name"] = name

				projectArray = append(projectArray, id)

				// var tbFilter string
				circleMeetCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS NoofCircleMeeting
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (4, 5, 6, 7, 8, 9)
										WHERE p.id IS NOT NULL
											AND p.id = ?
										GROUP BY p.id`

				if err := DB.QueryRow(circleMeetCountQuery, id).Scan(&circleMeetCount); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				getmodulecompletedquery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE (t2.module1=1 and t2.module2=1 and t2.module3=1 and t2.module4=1 and t2.module5=1)
				AND t3.id = ?`

				if err := DB.QueryRow(getmodulecompletedquery, id).Scan(&noofsporthimodulecompleted); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				sporthisurvey := `SELECT count(t2.id) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.id = ?`

				if err := DB.QueryRow(sporthisurvey, id).Scan(&noofsporthisurvey); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				enrollgelathi := `SELECT COUNT(tp.id) as gelathiCount FROM training_participants tp INNER JOIN project p ON p.id = tp.enrolledProject WHERE enroll = 1 and p.id = ?`

				if err := DB.QueryRow(enrollgelathi, id).Scan(&noofenrollgelathi); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				var noofactual int
				actual := `select count(tp.id) as actual from training_participants tp where day2 = 1 and enroll=1 and project_id=?`
				if err := DB.QueryRow(actual, id).Scan(&noofactual); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				var villages, subvillages int
				newvillages := `SELECT COUNT(distinct tp.location_id) as villages FROM bdms_staff.tbl_poa tp INNER JOIN training_participants tr on tr.tb_id=tp.tb_id join project p on p.id=tr.project_id where (enroll=1) and p.id=?`
				if err := DB.QueryRow(newvillages, id).Scan(&villages); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				subvillages1 := `SELECT COUNT(distinct tp.sub_village) as sub_villages FROM bdms_staff.tbl_poa tp INNER JOIN training_participants tr on tr.tb_id=tp.tb_id join project p on p.id=tr.project_id where (enroll=1) and  p.id=?`
				if err := DB.QueryRow(subvillages1, id).Scan(&subvillages); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				//var tpFilter string
				if gfid > 0 {
					target := getTrainerTarget(DB, gfid, projectArray)
					obj["target"] = target
					summaryTarget += target
				} else {
					target := getTarget(DB, startDate, endDate, projectArray)
					obj["target"] = target
					summaryTarget += target
				}

				obj["actual"] = noofactual
				summaryActuals += noofactual

				Noofcirclemeeting := circleMeetCount
				obj["NoofCircleMeeting"] = Noofcirclemeeting
				summarycirclemeeting += Noofcirclemeeting

				obj["Noofsporthisurvey"] = noofsporthisurvey
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = noofsporthimodulecompleted
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				obj["villages"] = villages + subvillages
				summaryVillages += obj["villages"].(int)

				obj["start_date"] = startDate
				obj["end_date"] = endDate
				obj["select_type"] = "1"

				obj["Gelathienrolled"] = noofenrollgelathi
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				// obj["Noofbeehives"] = getGFData(DB, "", 3, empid)
				// summaryofbeehives += obj["Noofbeehives"].(int)

				data = append(data, obj)

			}
		}

		fmt.Println(dateFilters)

		response := make(map[string]interface{})
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals
		response["summary_NoofCircleMeeting"] = summarycirclemeeting
		response["summary_sporthisurvey"] = summarysporthisurvey
		filter += summaryFilter
		response["summary_villages"] = summaryVillages
		response["summary_Noofsporthicompleted"] = summarysporthicompleted
		response["summary_Gelathienrolled"] = summaryGelathiEnrolled
		response["summary_Noofbeehives"] = summaryofbeehives
		response["data"] = data
		response["code"] = 200
		response["success"] = true
		response["message"] = "Successfully"

		js, err := json.Marshal(response)
		if err != nil {
			log.Println("GelathiProgramDashboard", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return

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
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + startDate + "' and p.endDate <= '" + endDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()
				if rows.Next() {
					query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
						"from tbl_poa tp " +
						"inner join project p on p.id = tp.project_id " +
						"where user_id = " + (empID) + " and tp.project_id = " + (projectID) +
						dateFilter +
						" GROUP by tp.project_id"
					// summaryProjectsArray = append(summaryProjectsArray, projectID)
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
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (empID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// circleMeetCountQuery := ""
			// var circleMeetCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate1, endDate1 string

				err := res.Scan(&id, &name, &startDate1, &endDate1)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["start_date"] = startDate1
				obj["end_date"] = endDate1
				obj["select_type"] = "1"

				// loop through each element in intSlice and convert to string

				obj["Gelathienrolled"] = Spoorthienrolledgelathi(DB, startDate, endDate, gfID, projectArray, empID)
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = GetNoOfspoorthibevee(DB, startDate, endDate, projectArray, gfID, empID)
				summaryofbeehives += obj["Noofbeehives"].(int)
				// Noofcirclemeeting := circleMeetCount
				// obj["NoofCircleMeeting"] = Noofcirclemeeting
				// summarycirclemeeting += Noofcirclemeeting
				NoofCircleMeeting := NoofCircleMeeting(DB, startDate, endDate, projectArray, gfID, empID)
				obj["NoofCircleMeeting"] = NoofCircleMeeting
				summarycirclemeeting += NoofCircleMeeting

				obj["Noofsporthisurvey"] = GetNoOfspoorthiSurvey(DB, startDate, endDate, gfID, projectArray, empID)
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB, startDate, endDate, projectArray, gfID, empID)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				Actual := Getspoorthiactual(DB, startDate, endDate, projectArray, gfID, empID)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["villages"] = newVillageCount(DB, startDate, endDate, gfID, projectArray, empID)
				summaryVillages += obj["villages"].(int)

				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target

				data = append(data, obj)
			}

			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			//tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			// intSlice := []int{}

			// loop through each element in the []interface{} slice
			// for _, v := range summaryProjectsArray {

			// 	if i, ok := v.(int); ok {

			// 		intSlice = append(intSlice, i)
			// 	}
			// }

			//getSummaryOfVillagesNew(DB, startDate, endDate, intSlice, tbFilter)
			response["summary_villages"] = summaryVillages
			response["summary_actual"] = summaryActuals
			response["summary_target"] = summaryTarget
			response["summary_NoofCircleMeeting"] = summarycirclemeeting
			response["summary_sporthisurvey"] = summarysporthisurvey
			response["summary_Noofsporthicompleted"] = summarysporthicompleted
			response["summary_Noofbeehives"] = summaryofbeehives

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
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
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + startDate + "' and p.endDate <= '" + endDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()
				if rows.Next() {
					query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
						"from tbl_poa tp " +
						"inner join project p on p.id = tp.project_id " +
						"where user_id = " + (empID) + " and tp.project_id = " + (projectID) +
						dateFilter +
						" GROUP by tp.project_id"
					// summaryProjectsArray = append(summaryProjectsArray, projectID)
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
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (empID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// circleMeetCountQuery := ""
			// var circleMeetCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate1, endDate1 string

				err := res.Scan(&id, &name, &startDate1, &endDate1)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["start_date"] = startDate1
				obj["end_date"] = endDate1
				obj["select_type"] = "1"

				obj["Gelathienrolled"] = Spoorthienrolledgelathi(DB, startDate, endDate, gfID, projectArray, empID)
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)
				obj["Noofbeehives"] = GetNoOfspoorthibevee(DB, startDate, endDate, projectArray, gfID, empID)
				summaryofbeehives += obj["Noofbeehives"].(int)

				NoofCircleMeeting := NoofCircleMeeting(DB, startDate, endDate, projectArray, gfID, empID)
				obj["noofVyaparCohorts"] = NoofCircleMeeting
				summarycirclemeeting += NoofCircleMeeting

				obj["Noofsporthisurvey"] = GetNoOfspoorthiSurvey(DB, startDate, endDate, gfID, projectArray, empID)
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)

				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB, startDate, endDate, projectArray, gfID, empID)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				Actual := Getspoorthiactual(DB, startDate, endDate, projectArray, gfID, empID)
				obj["actual"] = Actual
				summaryActuals += Actual

				// tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				obj["villages"] = newVillageCount(DB, startDate, endDate, gfID, projectArray, empID)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			// intSlice := []int{}

			// // loop through each element in the []interface{} slice
			// for _, v := range summaryProjectsArray {

			// 	if i, ok := v.(int); ok {

			// 		intSlice = append(intSlice, i)
			// 	}
			// }

			response["summary_villages"] = summaryVillages
			response["summary_actual"] = summaryActuals
			response["summary_target"] = summaryTarget
			response["summary_NoofCircleMeeting"] = summarycirclemeeting
			response["summary_sporthisurvey"] = summarysporthisurvey
			response["summary_Noofsporthicompleted"] = summarysporthicompleted
			response["summary_Noofbeehives"] = summaryofbeehives

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		} else {

			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
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
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + startDate + "' and p.endDate <= '" + endDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", projectID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()
				if rows.Next() {
					query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
						"from tbl_poa tp " +
						"inner join project p on p.id = tp.project_id " +
						"where user_id = " + (empID) + " and tp.project_id = " + (projectID) +
						dateFilter +
						" GROUP by tp.project_id"
					// summaryProjectsArray = append(summaryProjectsArray, projectID)
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
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (empID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// circleMeetCountQuery := ""
			// var circleMeetCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate1, endDate1 string

				err := res.Scan(&id, &name, &startDate1, &endDate1)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["start_date"] = startDate1
				obj["end_date"] = endDate1
				obj["select_type"] = "1"

				// summaryEnrolled += obj["enrolled"].(int)
				// tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				// strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string

				obj["Gelathienrolled"] = Spoorthienrolledgelathi(DB, startDate, endDate, gfID, projectArray, empID)
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = GetNoOfspoorthibevee(DB, startDate, endDate, projectArray, gfID, empID)
				summaryofbeehives += obj["Noofbeehives"].(int)

				NoofCircleMeeting := NoofCircleMeeting(DB, startDate, endDate, projectArray, gfID, empID)
				obj["noofVyaparCohorts"] = NoofCircleMeeting
				summarycirclemeeting += NoofCircleMeeting

				obj["Noofsporthisurvey"] = GetNoOfspoorthiSurvey(DB, startDate, endDate, gfID, projectArray, empID)
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)

				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB, startDate, endDate, projectArray, gfID, empID)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				Actual := Getspoorthiactual(DB, startDate, endDate, projectArray, gfID, empID)
				obj["actual"] = Actual
				summaryActuals += Actual

				// tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				obj["villages"] = newVillageCount(DB, startDate, endDate, gfID, projectArray, empID)
				summaryVillages += obj["villages"].(int)

				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			// intSlice := []int{}

			// // loop through each element in the []interface{} slice
			// for _, v := range summaryProjectsArray {

			// 	if i, ok := v.(int); ok {

			// 		intSlice = append(intSlice, i)
			// 	}
			// }
			response["summary_villages"] = summaryVillages
			response["summary_actual"] = summaryActuals
			response["summary_target"] = summaryTarget
			response["summary_NoofCircleMeeting"] = summarycirclemeeting
			response["summary_sporthisurvey"] = summarysporthisurvey
			response["summary_Noofsporthicompleted"] = summarysporthicompleted
			response["summary_Noofbeehives"] = summaryofbeehives

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(js)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println("GelathiProgramDashboard", err)
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
