package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var DB *sql.DB

type ProjectRequest[T any] struct {
	PartnerID  string `json:"partner_id"`
	Dist       string `json:"dist"`
	Taluk      string `json:"taluk"`
	Filter     int    `json:"filter"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	FunderId   string `json:"funder_id"`
	ProjectID  string `json:"project_id"`
	TrainerID  string `json:"trainer_id"`
	OpsManager string `json:"opsmanager"`
	SOMID      string `json:"somid"`
	GFLID      string `json:"gflid"`
	RoleID     string `json:"roleid"`
	GalathiID  string `json:"galathi_id"`
	EmpID      string `json:"emp_id"`
	GfID       string `json:"gf_id"`
}

// ---------------------vyapar dashboard----------------------------------
func VyaparDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method == "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return

	}

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	gfid, _ := strconv.Atoi(reqBody.GFLID)
	som, _ := strconv.Atoi(reqBody.SOMID)
	opsid, _ := strconv.Atoi(reqBody.OpsManager)
	trainid, _ := strconv.Atoi(reqBody.TrainerID)
	projectid, _ := strconv.Atoi(reqBody.ProjectID)
	funderid, _ := strconv.Atoi(reqBody.FunderId)
	talukid, _ := strconv.Atoi(reqBody.Taluk)
	distid, _ := strconv.Atoi(reqBody.Dist)
	empid, _ := strconv.Atoi(reqBody.EmpID)
	partnerid, _ := strconv.Atoi(reqBody.PartnerID) //converting string to int
	roleid, _ := strconv.Atoi(reqBody.RoleID)
	isDateFilterApplied := false
	if reqBody.StartDate == "" || reqBody.StartDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	data := []interface{}{}

	summaryProjectsArray := []interface{}{}

	summaryVillages := 0
	summaryTarget := 0
	summaryActuals := 0

	summaryVyaparenolled := 0
	summaryNoofvyaparsurvey := 0
	summaryNoofvyaparmodulecompleted := 0
	summayNoofvyaparcoharts := 0

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 12 {
		filter := ""
		summaryFilter := ""
		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", reqBody.EmpID)
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
					log.Println("VyaparDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", reqBody.EmpID)
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
					log.Println("VyaparDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", reqBody.EmpID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds := getOpsManagers(DB, empid)
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
					log.Println("VyaparDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", reqBody.EmpID)
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
					log.Println("VyaparDashboard", err)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", reqBody.EmpID)
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
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}

		}
		cohortCountQuery := ""
		var cohortCount int
		var noofvyaparmodulecompleted, noofvyaparsurvey, noofenrollvyapar int

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
			dateFilters = "date >= '" + reqBody.StartDate + "' and date <= '" + reqBody.EndDate + "'"
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		// CEO Dashboard
		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", reqBody.PartnerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id ,funderName as name FROM project p inner join funder on funder.funderID = p.funderID where p.partnerID = " + (reqBody.PartnerID) + " and " + dateFilter + filter
				filter += " and p.partnerID = " + (reqBody.PartnerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if distid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Dist)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Taluk)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = "SELECT p.funderID as id,funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID = " + (reqBody.Taluk) + " and " + dateFilter + filter + " GROUP by p.funderID"
						filter += " and locationID = " + (reqBody.Taluk)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)

						if err != nil {
							//log.Println(err)

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
					getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + (reqBody.Dist)
					talukArray := []int{}
					talukRes, err := DB.Query(getTaluk)
					if err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer talukRes.Close()
					for talukRes.Next() {
						var tlk int
						err := talukRes.Scan(&tlk)
						if err != nil {
							log.Println("VyaparDashboard", err)
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
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", reqBody.FunderId)
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
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if partnerid == 0 && trainid == 0 && opsid == 0 && som == 0 && gfid == 0 && !isDateFilterApplied && roleid != 4 {
			// role 4 OpsManager Default should be project list
			funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) as id,funderName as name FROM project p inner join funder on p.funderId = funder.funderID where %s%s", dateFilter, filter)
		}
		if len(funderListQuery) > 0 {
			if reqBody.ProjectID == "" {
				res, err := DB.Query(funderListQuery)
				if err != nil {
					log.Println("VyaparDashboard", err)
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
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					// if len(funderListQuery) > 0 {
					cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (16, 17, 18, 19, 20, 21)
										WHERE p.funderId IS NOT NULL
											AND p.funderId = ?
										GROUP BY p.funderId`

					if err := DB.QueryRow(cohortCountQuery, funderId).Scan(&cohortCount); err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					// } else if reqBody.StartDate != "" && reqBody.EndDate != "" {
					// 	cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
					// 					FROM project p
					// 					LEFT JOIN tbl_poa tp ON p.id = tp.project_id
					// 						AND tp.type = 2
					// 						AND tp.session_type IN (16, 17, 18, 19, 20, 21)
					// 					WHERE p.funderId IS NOT NULL
					// 						AND p.funderId = ? AND date between ? AND ?
					// 					GROUP BY p.funderId`

					// 	if err := DB.QueryRow(cohortCountQuery, funderId, reqBody.StartDate, reqBody.EndDate).Scan(&cohortCount); err != nil {
					// 		log.Println("VyaparDashboard", err)
					// 		w.WriteHeader(http.StatusBadRequest)
					// 		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					// 		return
					// 	}

					// }
					getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN BuzzVyaparProgramBaseline t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.funderId = ? `

					if err := DB.QueryRow(getActualsQuery, funderId).Scan(&noofvyaparmodulecompleted); err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					vyaprsurvey := `SELECT count(t2.id) FROM training_participants t1
				JOIN BuzzVyaparProgramBaseline t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.funderId = ?`

					if err := DB.QueryRow(vyaprsurvey, funderId).Scan(&noofvyaparsurvey); err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					enrollvyapar := `SELECT COUNT(tp.id) as Vyapar FROM training_participants tp INNER JOIN project p ON p.id = tp.VyaparEnrollmentEnrolledProject WHERE VyaparEnrollment = 1 and p.funderID= ? `

					if err := DB.QueryRow(enrollvyapar, funderId).Scan(&noofenrollvyapar); err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}

					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("VyaparDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var id int
						if err := projResult.Scan(&id); err != nil {
							log.Println("VyaparDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						projectArray = append(projectArray, id)
					}

					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":                        funderId,
							"name":                      funderName,
							"target":                    0,
							"actual":                    0,
							"noofvyaparcoharts":         0,
							"nofvyaparsurvey":           0,
							"noofvyaparmodulecompleted": 0,
							"villages":                  0,
							"startDate":                 "",
							"endDate":                   "",
							"select_type":               "2",
						}
						data = append(data, obj)
						continue
					}
					// str := strconv.Itoa(reqBody.FunderId)
					stringSlice := make([]string, len(projectArray))

					for i, val := range projectArray {
						stringSlice[i] = strconv.Itoa(val)
					}

					obj := map[string]interface{}{
						"id":                funderId,
						"name":              funderName,
						"target":            getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":            getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"noofvyaparcoharts": cohortCount,

						"nofvyaparsurvey":           noofvyaparsurvey,
						"noofvyaparmodulecompleted": noofvyaparmodulecompleted,
						"villages":                  newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas

						"startDate":      "",
						"endDate":        "",
						"select_type":    "2",
						"vyaparenrolled": noofenrollvyapar,
					}

					summaryTarget += obj["target"].(int)
					summaryActuals += obj["actual"].(int)
					summaryNoofvyaparsurvey += obj["nofvyaparsurvey"].(int)
					summaryNoofvyaparmodulecompleted += obj["noofvyaparmodulecompleted"].(int)

					summayNoofvyaparcoharts += obj["noofvyaparcoharts"].(int)
					summaryVillages += obj["villages"].(int)
					summaryVyaparenolled += obj["vyaparenrolled"].(int)

					data = append(data, obj)
				}
			}
		}

		projectList := ""

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + (reqBody.ProjectID) + filter + dateFilterNew

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if trainid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", reqBody.TrainerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + (reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + (reqBody.TrainerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				if err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if opsid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", reqBody.OpsManager)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + (reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + (reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + (reqBody.OpsManager)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if som > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", reqBody.SOMID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if gfid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", reqBody.GFLID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)

				if err != nil {
					log.Println("VyaparDashboard", err)
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
				log.Println("VyaparDashboard", err)
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
				var startDate string
				var endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}

				obj["id"] = id
				obj["name"] = name

				projectArray = append(projectArray, id)

				var tpFilter string
				var tbFilter string
				// if len(funderListQuery) > 0 {
				cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
				FROM project p
				LEFT JOIN tbl_poa tp ON p.id = tp.project_id
					AND tp.type = 2
					AND tp.session_type IN (16, 17, 18, 19, 20, 21)
				WHERE p.id IS NOT NULL
					AND p.id = ?
				GROUP BY p.id`

				if err := DB.QueryRow(cohortCountQuery, id).Scan(&cohortCount); err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				// } else if reqBody.StartDate != "" && reqBody.EndDate != "" {
				// 	cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
				// 					FROM project p
				// 					LEFT JOIN tbl_poa tp ON p.id = tp.project_id
				// 						AND tp.type = 2
				// 						AND tp.session_type IN (16, 17, 18, 19, 20, 21)
				// 					WHERE p.id IS NOT NULL
				// 						AND p.id = ? AND date between ? AND ?
				// 					GROUP BY p.id`

				// 	if err := DB.QueryRow(cohortCountQuery, id, reqBody.StartDate, reqBody.EndDate).Scan(&cohortCount); err != nil {
				// 		log.Println("VyaparDashboard", err)
				// 		w.WriteHeader(http.StatusBadRequest)
				// 		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				// 		return
				// 	}

				// }

				// cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
				// 						FROM project p
				// 						LEFT JOIN tbl_poa tp ON p.id = tp.project_id
				// 							AND tp.type = 2
				// 							AND tp.session_type IN (16, 17, 18, 19, 20, 21)
				// 						WHERE p.id IS NOT NULL
				// 							AND p.id = ?
				// 						GROUP BY p.id`

				// if err := DB.QueryRow(cohortCountQuery, id).Scan(&cohortCount); err != nil {
				// 	log.Println("VyaparDashboard", err)
				// 	w.WriteHeader(http.StatusBadRequest)
				// 	json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				// 	return
				// }
				getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN BuzzVyaparProgramBaseline t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.id = ?`

				if err := DB.QueryRow(getActualsQuery, id).Scan(&noofvyaparmodulecompleted); err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				vyaprsurvey := `SELECT count(t2.id) FROM training_participants t1
				JOIN BuzzVyaparProgramBaseline t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.id = ?`

				if err := DB.QueryRow(vyaprsurvey, id).Scan(&noofvyaparsurvey); err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				enrollvyapar := `SELECT COUNT(tp.id) as Vyapar FROM training_participants tp INNER JOIN project p ON p.id = tp.VyaparEnrollmentEnrolledProject WHERE VyaparEnrollment = 1 and p.id = ? `

				if err := DB.QueryRow(enrollvyapar, id).Scan(&noofenrollvyapar); err != nil {
					log.Println("VyaparDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				if trainid > 0 {
					target := getTrainerTarget(DB, trainid, projectArray)
					obj["target"] = target
					summaryTarget += target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					target := getTarget(DB, startDate, endDate, projectArray)
					obj["target"] = target
					summaryTarget += target
				}

				actual := getActual(DB, startDate, endDate, projectArray, tpFilter)
				obj["actual"] = actual
				summaryActuals += actual

				obj["NoOfvyaparcoharts"] = cohortCount
				summayNoofvyaparcoharts += obj["NoOfvyaparcoharts"].(int)

				obj["noOfvyaparsurvey"] = noofvyaparsurvey
				summaryNoofvyaparsurvey += obj["noOfvyaparsurvey"].(int)
				obj["vyaparenrolled"] = noofenrollvyapar
				summaryVyaparenolled += obj["vyaparenrolled"].(int)
				obj["noofvyaparmodulecompleted"] = noofvyaparmodulecompleted
				summaryVyaparenolled += obj["noofvyaparmodulecompleted"].(int)

				obj["villages"] = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				data = append(data, obj)

			}
		}

		fmt.Println(dateFilters)

		response := make(map[string]interface{})
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals

		response["summary_Noofvyaparcoharts"] = summayNoofvyaparcoharts
		response["summary_NoofVyaparsurvey"] = summaryNoofvyaparsurvey
		filter += summaryFilter
		response["summary_villages"] = summaryVillages

		response["summary_Noofvyaparmodulecomoleted"] = summaryNoofvyaparmodulecompleted
		response["summary_vyaparenrolled"] = summaryVyaparenolled
		response["data"] = data
		response["code"] = 200
		response["success"] = true
		response["message"] = "Successfully"

		js, err := json.Marshal(response)

		if err != nil {
			log.Println("VyaparDashboard", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return

	} else if roleid == 5 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)
					if err != nil {
						log.Println("VyaparDashboard", err)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println("VyaparDashboard", err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println("VyaparDashboard", err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofVyaparCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofVyaparCohorts
				summayNoofvyaparcoharts += NoofVyaparCohorts

				NoOfvyaparsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfvyaparsurvey
				summaryNoofvyaparsurvey += NoOfvyaparsurvey

				noofvyaparmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofvyaparmodulcompleted
				summaryNoofvyaparmodulecompleted += noofvyaparmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.trainer_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["vyapar"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryVyaparenolled += obj["vyapar"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_vyaparcoharts"] = summayNoofvyaparcoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_vyaparsurvey"] = summaryNoofvyaparsurvey
			response["summary_noofmodulecompleted"] = summaryNoofvyaparmodulecompleted
			response["summary_vyparenrolled"] = summaryVyaparenolled

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("VyaparDashboard", err)
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
				log.Println("VyaparDashboard", err)
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
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println("VyaparDashboard", err)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println("VyaparDashboard", err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println("VyaparDashboard", err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofVyaparCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofVyaparCohorts
				summayNoofvyaparcoharts += NoofVyaparCohorts

				NoOfvyaparsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfvyaparsurvey
				summaryNoofvyaparsurvey += NoOfvyaparsurvey

				noofvyaparmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofvyaparmodulcompleted
				summaryNoofvyaparmodulecompleted += noofvyaparmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["vyapar"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryVyaparenolled += obj["vyapar"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_vyaparcoharts"] = summayNoofvyaparcoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_vyaparsurvey"] = summaryNoofvyaparsurvey
			response["summary_noofmodulecompleted"] = summaryNoofvyaparmodulecompleted
			response["summary_vyparenrolled"] = summaryVyaparenolled

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("VyaparDashboard", err)
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
				log.Println("VyaparDashboard", err)
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

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println("VyaparDashboard", err)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println("VyaparDashboard", err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println("VyaparDashboard", err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofVyaparCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofVyaparCohorts
				summayNoofvyaparcoharts += NoofVyaparCohorts

				NoOfvyaparsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfvyaparsurvey
				summaryNoofvyaparsurvey += NoOfvyaparsurvey

				noofvyaparmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofvyaparmodulcompleted
				summaryNoofvyaparmodulecompleted += noofvyaparmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["vyapar"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryVyaparenolled += obj["vyapar"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_vyaparcoharts"] = summayNoofvyaparcoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_vyaparsurvey"] = summaryNoofvyaparsurvey
			response["summary_noofmodulecompleted"] = summaryNoofvyaparmodulecompleted
			response["summary_vyparenrolled"] = summaryVyaparenolled

			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("VyaparDashboard", err)
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
				log.Println("VyaparDashboard", err)
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

// ---------------------green dashboard----------------------------------
func GreenDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method == "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return

	}
	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	gfid, _ := strconv.Atoi(reqBody.GFLID)
	som, _ := strconv.Atoi(reqBody.SOMID)
	opsid, _ := strconv.Atoi(reqBody.OpsManager)
	trainid, _ := strconv.Atoi(reqBody.TrainerID)
	projectid, _ := strconv.Atoi(reqBody.ProjectID)
	funderid, _ := strconv.Atoi(reqBody.FunderId)
	talukid, _ := strconv.Atoi(reqBody.Taluk)
	distid, _ := strconv.Atoi(reqBody.Dist)
	empid, _ := strconv.Atoi(reqBody.EmpID)
	partnerid, _ := strconv.Atoi(reqBody.PartnerID) //converting string to int
	roleid, _ := strconv.Atoi(reqBody.RoleID)
	isDateFilterApplied := false
	if reqBody.StartDate == "" || reqBody.StartDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	data := []interface{}{}

	summaryProjectsArray := []interface{}{}

	summaryVillages := 0
	summaryTarget := 0
	summaryActuals := 0

	summaryGreenrenolled := 0
	summaryNoofGreensurvey := 0
	summaryNoofGreenmodulecompleted := 0
	summayNoofGreencoharts := 0

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 12 {
		filter := ""
		summaryFilter := ""
		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", reqBody.EmpID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds := getOpsManagers(DB, empid)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", reqBody.EmpID)
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
		cohortCountQuery := ""
		var cohortCount, noofenrollgreen int

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
			dateFilters = "date >= '" + reqBody.StartDate + "' and date <= '" + reqBody.EndDate + "'"
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		// CEO Dashboard
		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", reqBody.PartnerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id ,funderName as name FROM project p inner join funder on funder.funderID = p.funderID where p.partnerID = " + (reqBody.PartnerID) + " and " + dateFilter + filter
				filter += " and p.partnerID = " + (reqBody.PartnerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)

				if err != nil {
					//log.Println(err)

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if distid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Dist)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Taluk)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = "SELECT p.funderID as id,funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID = " + (reqBody.Taluk) + " and " + dateFilter + filter + " GROUP by p.funderID"
						filter += " and locationID = " + (reqBody.Taluk)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)

						if err != nil {
							//log.Println(err)

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
					getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + (reqBody.Dist)
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
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid locationid"
				js, err := json.Marshal(response)

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
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", reqBody.FunderId)
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
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if partnerid == 0 && trainid == 0 && opsid == 0 && som == 0 && gfid == 0 && !isDateFilterApplied && roleid != 4 {
			// role 4 OpsManager Default should be project list
			funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) as id,funderName as name FROM project p inner join funder on p.funderId = funder.funderID where %s%s", dateFilter, filter)
		}
		if len(funderListQuery) > 0 {
			if reqBody.ProjectID == "" {

				res, err := DB.Query(funderListQuery)
				if err != nil {
					log.Println("GreenDashboard", err)
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
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgreencohorts
                                        FROM project p
                                        LEFT JOIN tbl_poa tp ON p.id = tp.project_id
                                            AND tp.type = 2
                                            AND tp.session_type IN (10, 11, 12, 13, 14, 15)
                                        WHERE p.funderId IS NOT NULL
                                            AND p.funderId = ?
                                        GROUP BY p.funderId`

					if err := DB.QueryRow(cohortCountQuery, funderId).Scan(&cohortCount); err != nil {
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
                JOIN GreenBaselineSurvey t2 ON t1.id = t2.partcipantId
                JOIN project t3 ON t1.project_id = t3.id
                WHERE t1.id = t2.partcipantId
                AND t3.funderId = ?`

					var noofgreenmodulecompleted int
					if err := DB.QueryRow(getActualsQuery, funderId).Scan(&noofgreenmodulecompleted); err != nil {
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					greensurvey := `SELECT count(t2.id) FROM training_participants t1
                JOIN GreenBaselineSurvey t2 ON t1.id = t2.partcipantId
                JOIN project t3 ON t1.project_id = t3.id
                WHERE t1.id = t2.partcipantId
                AND t3.funderId = ?`

					var noofgreensurvey int
					if err := DB.QueryRow(greensurvey, funderId).Scan(&noofgreensurvey); err != nil {
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					enrollgreen := `SELECT COUNT(tp.id) as greenMoti FROM training_participants tp INNER JOIN project p ON p.id = tp.GreenMotivatorsEnrolledProject WHERE GreenMotivators = 1 and p.funderID = ? `

					if err := DB.QueryRow(enrollgreen, funderId).Scan(&noofenrollgreen); err != nil {
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}

					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("GreenDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var id int
						if err := projResult.Scan(&id); err != nil {
							log.Println("GreenDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						projectArray = append(projectArray, id)
					}

					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":                       funderId,
							"name":                     funderName,
							"target":                   0,
							"actual":                   0,
							"noofgreencoharts":         0,
							"nofgreensurvey":           0,
							"noofgreenmodulecompleted": 0,
							"villages":                 0,
							"startDate":                "",
							"endDate":                  "",
							"select_type":              "2",
						}
						data = append(data, obj)
						continue
					}
					// str := strconv.Itoa(reqBody.FunderId)
					stringSlice := make([]string, len(projectArray))

					for i, val := range projectArray {
						stringSlice[i] = strconv.Itoa(val)
					}

					obj := map[string]interface{}{
						"id":               funderId,
						"name":             funderName,
						"target":           getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":           getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"noofgreencoharts": cohortCount,

						"nofgreensurvey":           noofgreensurvey,
						"noofgreenmodulecompleted": noofgreenmodulecompleted,
						"villages":                 newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas

						"startDate":     "",
						"endDate":       "",
						"select_type":   "2",
						"greenenrolled": noofenrollgreen,
					}

					summaryTarget += obj["target"].(int)
					summaryActuals += obj["actual"].(int)
					summaryNoofGreensurvey += obj["nofgreensurvey"].(int)
					summaryNoofGreenmodulecompleted += obj["noofgreenmodulecompleted"].(int)

					summayNoofGreencoharts += obj["noofgreencoharts"].(int)
					summaryVillages += obj["villages"].(int)
					summaryGreenrenolled += obj["greenenrolled"].(int)

					data = append(data, obj)
				}
			}
		}

		projectList := ""

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + (reqBody.ProjectID) + filter + dateFilterNew

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)

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
		} else if trainid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", reqBody.TrainerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + (reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + (reqBody.TrainerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)

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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", reqBody.OpsManager)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + (reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + (reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + (reqBody.OpsManager)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)

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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", reqBody.SOMID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)

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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", reqBody.GFLID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)

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
		} else if (isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0) || (roleid == 4 && distid == 0) {
			//role 4 - OpsManager Default should be project list without location filter
			projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
		}

		if len(projectList) > 0 {

			res, err := DB.Query(projectList)
			if err != nil {
				log.Println("GreenDashboard", err)
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
				var startDate string
				var endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}

				obj["id"] = id
				obj["name"] = name

				projectArray = append(projectArray, id)

				var noofgreensurvey int

				var tpFilter string
				var tbFilter string
				greensurvey := `SELECT count(t2.id) FROM training_participants t1
            JOIN GreenBaselineSurvey t2 ON t1.id = t2.partcipantId
            JOIN project t3 ON t1.project_id = t3.id
            WHERE t1.id = t2.partcipantId
            AND t3.id = ?`

				if err := DB.QueryRow(greensurvey, id).Scan(&noofgreensurvey); err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgreencohorts
                                    FROM project p
                                    LEFT JOIN tbl_poa tp ON p.id = tp.project_id
                                        AND tp.type = 2
                                        AND tp.session_type IN (10, 11, 12, 13, 14, 15)
                                    WHERE p.id IS NOT NULL
                                        AND p.id = ?
                                    GROUP BY p.id`

				if err := DB.QueryRow(cohortCountQuery, id).Scan(&cohortCount); err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
            JOIN GreenBaselineSurvey t2 ON t1.id = t2.partcipantId
            JOIN project t3 ON t1.project_id = t3.id
            WHERE t1.id = t2.partcipantId
            AND t3.id = ?`

				var noofgreenmodulecompleted int
				if err := DB.QueryRow(getActualsQuery, id).Scan(&noofgreenmodulecompleted); err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				enrollgreen := `SELECT COUNT(tp.id) as greenMoti FROM training_participants tp INNER JOIN project p ON p.id = tp.GreenMotivatorsEnrolledProject WHERE GreenMotivators = 1 and p.id = ? `

				if err := DB.QueryRow(enrollgreen, id).Scan(&noofenrollgreen); err != nil {
					log.Println("GreenDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}

				if trainid > 0 {
					target := getTrainerTarget(DB, trainid, projectArray)
					obj["target"] = target
					summaryTarget += target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					target := getTarget(DB, startDate, endDate, projectArray)
					obj["target"] = target
					summaryTarget += target
				}

				actual := getActual(DB, startDate, endDate, projectArray, tpFilter)
				obj["actual"] = actual
				summaryActuals += actual

				obj["NoOfgreencoharts"] = cohortCount
				summayNoofGreencoharts += obj["NoOfgreencoharts"].(int)

				obj["noOfgreensurvey"] = noofgreensurvey
				summaryNoofGreensurvey += obj["noOfgreensurvey"].(int)
				obj["greenenrolled"] = noofenrollgreen
				summaryGreenrenolled += obj["greenenrolled"].(int)
				obj["noofgreenmodulecompleted"] = noofgreenmodulecompleted
				summaryNoofGreenmodulecompleted += obj["noofgreenmodulecompleted"].(int)

				obj["villages"] = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				data = append(data, obj)

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

	} else if roleid == 5 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)

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
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofGreenCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGreenCohorts
				summayNoofGreencoharts += NoofGreenCohorts

				NoOfGreenrsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfGreenrsurvey
				summaryNoofGreensurvey += NoOfGreenrsurvey

				noofGreenmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofGreenmodulcompleted
				summaryNoofGreenmodulecompleted += noofGreenmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.trainer_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["Greenenrolled"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryGreenrenolled += obj["Greenenrolled"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Greencoharts"] = summayNoofGreencoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_Greensurvey"] = summaryNoofGreensurvey
			response["summary_noofGreenmodulecompleted"] = summaryNoofGreenmodulecompleted
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
		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

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

	} else if roleid == 6 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)

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
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofGreenCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGreenCohorts
				summayNoofGreencoharts += NoofGreenCohorts

				NoOfGreenrsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfGreenrsurvey
				summaryNoofGreensurvey += NoOfGreenrsurvey

				noofGreenmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofGreenmodulcompleted
				summaryNoofGreenmodulecompleted += noofGreenmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["Greenenrolled"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryGreenrenolled += obj["Greenenrolled"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Greencoharts"] = summayNoofGreencoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_Greensurvey"] = summaryNoofGreensurvey
			response["summary_noofGreenmodulecompleted"] = summaryNoofGreenmodulecompleted
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

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

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

	} else if roleid == 13 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
				} else {
					// showNoProj()
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid project id"
					js, err := json.Marshal(response)

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
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				str := (reqBody.EmpID)

				NoofGreenCohorts := NoofVyaparCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGreenCohorts
				summayNoofGreencoharts += NoofGreenCohorts

				NoOfGreenrsurvey := GetNoOfVyaparSurvey(DB, str, "", "")
				obj["noOfvyaparsurvey"] = NoOfGreenrsurvey
				summaryNoofGreensurvey += NoOfGreenrsurvey

				noofGreenmodulcompleted := GetNoofVyaparModuleCompleted(DB)
				obj["noofvyaparmodulcompleted"] = noofGreenmodulcompleted
				summaryNoofGreenmodulecompleted += noofGreenmodulcompleted

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				obj["Greenenrolled"] = Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)

				summaryGreenrenolled += obj["Greenenrolled"].(int)
				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Greencoharts"] = summayNoofGreencoharts

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_Greensurvey"] = summaryNoofGreensurvey
			response["summary_noofGreenmodulecompleted"] = summaryNoofGreenmodulecompleted
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

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid employe id"
			js, err := json.Marshal(response)

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
	} else {
		w.WriteHeader(http.StatusCreated)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid role id"
		json.NewEncoder(w).Encode(response)
	}

}

// ---------------------selfSakthi dashboard----------------------------------
func SelfSakthiDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	gfid, _ := strconv.Atoi(reqBody.GFLID)
	som, _ := strconv.Atoi(reqBody.SOMID)
	opsid, _ := strconv.Atoi(reqBody.OpsManager)
	trainid, _ := strconv.Atoi(reqBody.TrainerID)
	projectid, _ := strconv.Atoi(reqBody.ProjectID)
	funderid, _ := strconv.Atoi(reqBody.FunderId)
	talukid, _ := strconv.Atoi(reqBody.Taluk)
	distid, _ := strconv.Atoi(reqBody.Dist)
	empid, _ := strconv.Atoi(reqBody.EmpID)
	roleid, _ := strconv.Atoi(reqBody.RoleID)
	partnerid, _ := strconv.Atoi(reqBody.PartnerID) //converting string to int

	isDateFilterApplied := false
	if reqBody.StartDate == "" || reqBody.StartDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	data := []interface{}{}
	summaryProjectsArray := []interface{}{}
	summaryTarget := 0
	summaryVillages := 0
	summaryActuals := 0
	summaryDay1 := 0
	summaryNoofBatches := 0
	summaryNoofSelfShakthisurvey := 0
	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 12 {
		filter := ""
		summaryFilter := ""

		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", reqBody.EmpID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds := getOpsManagers(DB, empid)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", reqBody.EmpID)
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

			dateFilter = "startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
			dateFilters = "date >= '" + reqBody.StartDate + "' and date <= '" + reqBody.EndDate + "'"
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		// CEO Dashboard
		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", reqBody.PartnerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id ,funderName as name FROM project p inner join funder on funder.funderID = p.funderID where p.partnerID = " + (reqBody.PartnerID) + " and " + dateFilter + filter
				filter += " and p.partnerID = " + (reqBody.PartnerID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Dist)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Taluk)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = "SELECT p.funderID as id,funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID = " + (reqBody.Taluk) + " and " + dateFilter + filter + " GROUP by p.funderID"
						filter += " and locationID = " + (reqBody.Taluk)
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
					// get taluk of specified dist
					getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + (reqBody.Dist)
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
		} else if funderid > 0 {
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", reqBody.FunderId)
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
					log.Println("SelfSakthiDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")

				w.Write(js)
				return

			}
		} else if partnerid == 0 && trainid == 0 && opsid == 0 && som == 0 && gfid == 0 && !isDateFilterApplied && roleid != 4 {
			// role 4 OpsManager Default should be project list
			funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) as id,funderName as name FROM project p inner join funder on p.funderId = funder.funderID where %s%s", dateFilter, filter)
		}
		// funderList := []map[string]interface{}{}
		if len(funderListQuery) > 0 {
			if reqBody.ProjectID == "" {
				res, err := DB.Query(funderListQuery)
				if err != nil {
					log.Println("SelfSakthiDashboard", err)
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
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}

					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}

					projResult, err := DB.Query(getProj)
					if err != nil {
						log.Println("SelfSakthiDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer projResult.Close()

					for projResult.Next() {
						var id int
						if err := projResult.Scan(&id); err != nil {
							log.Println("SelfSakthiDashboard", err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						projectArray = append(projectArray, id)
					}

					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":     funderId,
							"name":   funderName,
							"target": 0,
							"actual": 0,
							"day2":   0,

							"villages":              0,
							"noofbatches":           0,
							"noofselfshakthisurvey": 0,
							"startDate":             "",
							"endDate":               "",
							"select_type":           "2",
						}
						data = append(data, obj)
						continue
					}

					obj := map[string]interface{}{
						"id":                    funderId,
						"name":                  funderName,
						"target":                getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":                getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"day2":                  0,
						"villages":              newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
						"noofbatches":           getTrainingBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"noofselfshakthisurvey": noOfselfshakthisurvey(DB),

						"startDate":   "",
						"endDate":     "",
						"select_type": "2",
					}

					day1Count := getDay1Count(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "")
					if day1Count > 0 {
						day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
						obj["day2"] = int(math.Round(day2Turnout * 100))
					}

					summaryTarget += obj["target"].(int)
					summaryActuals += obj["actual"].(int)
					summaryNoofBatches += obj["noofbatches"].(int)
					summaryNoofSelfShakthisurvey += obj["noofselfshakthisurvey"].(int)

					summaryDay1 += day1Count
					summaryVillages += obj["villages"].(int)

					data = append(data, obj)
				}
			}
		}
		projectList := ""

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + (reqBody.ProjectID) + filter + dateFilterNew

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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", reqBody.TrainerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + (reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + (reqBody.TrainerID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", reqBody.OpsManager)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + (reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + (reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + (reqBody.OpsManager)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", reqBody.SOMID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ")"
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
		} else if gfid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", reqBody.GFLID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ")"
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
		} else if (isDateFilterApplied && partnerid == 0 && distid == 0 && funderid == 0) || (roleid == 4 && distid == 0) {
			//role 4 - OpsManager Default should be project list without location filter
			projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
		}

		if len(projectList) > 0 {
			res, err := DB.Query(projectList)
			if err != nil {
				log.Println("SelfSakthiDashboard", err)
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
				var startDate string
				var endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}

				obj["id"] = id
				obj["name"] = name

				projectArray = append(projectArray, id)

				var tpFilter string
				var tbFilter string

				if trainid > 0 {
					target := getTrainerTarget(DB, trainid, projectArray)
					obj["target"] = target
					summaryTarget += target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					target := getTarget(DB, startDate, endDate, projectArray)
					obj["target"] = target
					summaryTarget += target
				}

				actual := getActual(DB, startDate, endDate, projectArray, tpFilter)
				obj["actual"] = actual
				summaryActuals += actual

				day1Count := getDay1Count(DB, startDate, endDate, projectArray, tpFilter)
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2Turnout := float64(actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2Turnout * 100))
				} else {
					obj["day2"] = 0
				}
				obj["NoOfBatches"] = getTrainingBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "")
				summaryNoofBatches += obj["NoOfBatches"].(int)

				obj["noOfselfshakthisurvey"] = noOfselfshakthisurvey(DB)
				summaryNoofSelfShakthisurvey += obj["noOfselfshakthisurvey"].(int)

				obj["villages"] = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				data = append(data, obj)

			}
		}

		fmt.Println(dateFilters)

		response := make(map[string]interface{})
		response["summary_target"] = summaryTarget
		response["summary_Noofselfshakthisurvey"] = summaryNoofSelfShakthisurvey
		response["summary_NoofBatches"] = summaryNoofBatches
		filter += summaryFilter
		response["summary_villages"] = summaryVillages
		response["summary_actual"] = summaryActuals
		var day2Turnout float64
		if summaryDay1 > 0 {
			day2Turnout = float64(summaryActuals) / float64(summaryDay1)
			response["summary_day2"] = int(math.Round(day2Turnout * 100))
		} else {
			response["summary_day2"] = 0
		}
		response["data"] = data
		response["code"] = 200
		response["success"] = true
		response["message"] = "Successfully"

		js, err := json.Marshal(response)

		if err != nil {
			log.Println("SelfSakthiDashboard", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return

	} else if roleid == 5 {
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			var summaryTarget, summayActual, summaryVillages, summaryofBatches, summayofSelfShakthi int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target

				Actual := getTrainerActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, reqBody.EmpID)
				obj["actual"] = Actual
				summayActual += Actual

				NoOfBatches := getTrainerBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, reqBody.EmpID)
				obj["Noofbatches"] = NoOfBatches
				summaryofBatches += NoOfBatches

				NoselfShakthisurvey := noofsurveyssTrainer(DB, reqBody.StartDate, reqBody.EndDate, projectArray, reqBody.EmpID)
				obj["noOfselfshakthisurvey"] = NoselfShakthisurvey
				summayofSelfShakthi += NoselfShakthisurvey

				// tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				Novillage := TrainerVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, reqBody.EmpID)
				obj["noofvilleges"] = Novillage
				summaryVillages += Novillage
				obj["select_type"] = "1"
				day1Count := getDay1Count(DB, startDate, endDate, projectArray, "")
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2Turnout := float64(Actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2Turnout * 100))
				} else {
					obj["day2"] = 0
				}

				data = append(data, obj)

			}

			response := make(map[string]interface{})

			response["summary_Target"] = summaryTarget

			// tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					_ = append(intSlice, i)
				}
			}
			response["summary_villages"] = summaryVillages

			response["summary_actual"] = summayActual
			response["summary_batches"] = summaryofBatches
			response["summay_selfShakthiSurvey"] = summayofSelfShakthi
			var day2Turnout float64
			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				response["summary_day2"] = 0
			}
			// response["summary_day2"] = summayday2
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
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
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			var summaryTarget, summayActual, summaryVillages, summaryofBatches, summayofSelfShakthi int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summayActual += Actual

				NoOfBatches := getTrainingBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "")
				obj["Noofbatches"] = NoOfBatches
				summaryofBatches += NoOfBatches

				NoselfShakthisurvey := noOfselfshakthisurvey(DB)
				obj["noOfselfshakthisurvey"] = NoselfShakthisurvey
				summayofSelfShakthi += NoselfShakthisurvey

				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)

				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"
				day1Count := getDay1Count(DB, startDate, endDate, projectArray, "")
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2Turnout := float64(Actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2Turnout * 100))
				} else {
					obj["day2"] = 0
				}

				data = append(data, obj)

			}

			response := make(map[string]interface{})

			response["summary_Target"] = summaryTarget

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)

			response["summary_actual"] = summayActual
			response["summary_batches"] = summaryofBatches
			response["summay_selfShakthiSurvey"] = summayofSelfShakthi
			var day2Turnout float64
			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				response["summary_day2"] = 0
			}
			// response["summary_day2"] = summayday2
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
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

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			var summaryTarget, summayActual, summaryVillages, summaryofBatches, summayofSelfShakthi int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.gelathi_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summayActual += Actual

				NoOfBatches := getTrainingBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "")
				obj["Noofbatches"] = NoOfBatches
				summaryofBatches += NoOfBatches

				NoselfShakthisurvey := noOfselfshakthisurvey(DB)
				obj["noOfselfshakthisurvey"] = NoselfShakthisurvey
				summayofSelfShakthi += NoselfShakthisurvey

				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)

				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"
				day1Count := getDay1Count(DB, startDate, endDate, projectArray, "")
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2Turnout := float64(Actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2Turnout * 100))
				} else {
					obj["day2"] = 0
				}

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Target"] = summaryTarget

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)

			response["summary_actual"] = summayActual
			response["summary_batches"] = summaryofBatches
			response["summay_selfShakthiSurvey"] = summayofSelfShakthi
			var day2Turnout float64
			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				response["summary_day2"] = 0
			}
			// response["summary_day2"] = summayday2
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			js, err := json.Marshal(response)

			if err != nil {
				log.Println("SelfSakthiDashboard", err)
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

// ---------------------gelathiProgram dashboard----------------------------------
func GelathiProgramDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	gfid, _ := strconv.Atoi(reqBody.GFLID)
	som, _ := strconv.Atoi(reqBody.SOMID)
	opsid, _ := strconv.Atoi(reqBody.OpsManager)
	trainid, _ := strconv.Atoi(reqBody.TrainerID)
	projectid, _ := strconv.Atoi(reqBody.ProjectID)
	funderid, _ := strconv.Atoi(reqBody.FunderId)
	talukid, _ := strconv.Atoi(reqBody.Taluk)
	distid, _ := strconv.Atoi(reqBody.Dist)
	empid, _ := strconv.Atoi(reqBody.EmpID)
	roleid, _ := strconv.Atoi(reqBody.RoleID)
	partnerid, _ := strconv.Atoi(reqBody.PartnerID) //converting string to int

	isDateFilterApplied := false
	if reqBody.StartDate == "" || reqBody.StartDate == " " {
		// do nothing
	} else {
		isDateFilterApplied = true
	}

	data := []interface{}{}
	summaryProjectsArray := []interface{}{}
	summaryTarget := 0
	summaryVillages := 0
	summaryActuals := 0
	summarycirclemeeting := 0
	summarysporthisurvey := 0
	summarysporthicompleted := 0
	summaryGelathiEnrolled := 0
	summaryofbeehives := 0

	if roleid == 1 || roleid == 9 || roleid == 3 || roleid == 4 || roleid == 12 {
		filter := ""
		summaryFilter := ""

		if roleid == 1 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=1 and id= ?", reqBody.EmpID)
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
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=3 and id= ?", reqBody.EmpID)
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id= ?", reqBody.EmpID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				opsIds := getOpsManagers(DB, empid)
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
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id= ?", reqBody.EmpID)
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
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=9 and id= ?", reqBody.EmpID)
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
				//
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
		cohortCountQuery := ""
		var cohortCount int
		var noofsporthimodulecompleted, noofsporthisurvey, noofenrollgelathi int

		dateFilter := ""
		dateFilters := ""
		if isDateFilterApplied {

			dateFilter = "startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
			dateFilters = "date >= '" + reqBody.StartDate + "' and date <= '" + reqBody.EndDate + "'"
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
			dateFilters = "date >= CURRENT_DATE()"
		}

		// CEO Dashboard
		funderListQuery := ""
		if partnerid > 0 {
			rows, err := DB.Query("SELECT partnerID FROM bdms_staff.project where partnerID= ?", reqBody.PartnerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id ,funderName as name FROM project p inner join funder on funder.funderID = p.funderID where p.partnerID = " + (reqBody.PartnerID) + " and " + dateFilter + filter
				filter += " and p.partnerID = " + (reqBody.PartnerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid partnerid"
				js, err := json.Marshal(response)
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Dist)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if talukid > 0 {
					rows, err := DB.Query("SELECT id FROM bdms_staff.location where id= ?", reqBody.Taluk)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}
					defer rows.Close()
					if rows.Next() {
						funderListQuery = "SELECT p.funderID as id,funderName as name from project p inner join funder on funder.funderID = p.funderID where locationID = " + (reqBody.Taluk) + " and " + dateFilter + filter + " GROUP by p.funderID"
						filter += " and locationID = " + (reqBody.Taluk)
					} else {
						// showNoProj()
						w.WriteHeader(http.StatusNotFound)
						response := make(map[string]interface{})
						response["success"] = false
						response["message"] = "Invalid locationid"
						js, err := json.Marshal(response)
						//
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
					getTaluk := "SELECT id from location l where `type` = 4 and parentId = " + (reqBody.Dist)
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
				//
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
			rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder where funderID = ?", reqBody.FunderId)
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
				//
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
		} else if partnerid == 0 && trainid == 0 && opsid == 0 && som == 0 && gfid == 0 && !isDateFilterApplied && roleid != 4 {
			// role 4 OpsManager Default should be project list
			funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) as id,funderName as name FROM project p inner join funder on p.funderId = funder.funderID where %s%s", dateFilter, filter)
		}
		if len(funderListQuery) > 0 {
			if reqBody.ProjectID == "" {

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
					cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgelathicohorts
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (4, 5, 6, 7, 8, 9)
										WHERE p.funderId IS NOT NULL
											AND p.funderId = ?
										GROUP BY p.funderId`

					if err := DB.QueryRow(cohortCountQuery, funderId).Scan(&cohortCount); err != nil {
						log.Println("GelathiProgramDashboard", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.funderId = ?`

					if err := DB.QueryRow(getActualsQuery, funderId).Scan(&noofsporthimodulecompleted); err != nil {
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

					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
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
							"NoofGelathiCohorts":   0,
							"villages":             0,
							"Gelathienrolled":      0,
							"Noofsporthisurvey":    0,
							"Noofsporthicompleted": 0,
							"Noofbeehives":         0,
							"startDate":            "",
							"endDate":              "",
							"select_type":          "2",
						}
						data = append(data, obj)
						continue
					}
					// str := strconv.Itoa(reqBody.FunderId)
					stringSlice := make([]string, len(projectArray))

					for i, val := range projectArray {
						stringSlice[i] = strconv.Itoa(val)
					}

					obj := map[string]interface{}{
						"id":                   funderId,
						"name":                 funderName,
						"target":               getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":               getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"NoofGelathiCohorts":   cohortCount,
						"villages":             newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
						"Gelathienrolled":      noofenrollgelathi,
						"Noofsporthisurvey":    noofsporthisurvey,
						"Noofsporthicompleted": noofsporthimodulecompleted,
						"Noofbeehives":         getGFData(DB, "", 3, empid),
						"startDate":            "",
						"endDate":              "",
						"select_type":          "2",
					}

					summaryTarget += obj["target"].(int)
					summaryActuals += obj["actual"].(int)
					summaryVillages += obj["villages"].(int)
					summaryGelathiEnrolled += obj["Gelathienrolled"].(int)
					summarycirclemeeting += obj["NoofGelathiCohorts"].(int)
					summarysporthisurvey += obj["Noofsporthisurvey"].(int)
					summarysporthicompleted += obj["Noofsporthicompleted"].(int)
					summaryofbeehives += obj["Noofbeehives"].(int)

					data = append(data, obj)
				}
			}
		}
		projectList := ""

		if projectid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + (reqBody.ProjectID) + filter + dateFilterNew

			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid project id"
				js, err := json.Marshal(response)
				//
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
		} else if trainid > 0 {
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id = ?", reqBody.TrainerID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + (reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + (reqBody.TrainerID)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid trainer id"
				js, err := json.Marshal(response)
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=4 and id = ?", reqBody.OpsManager)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + (reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + (reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + (reqBody.OpsManager)
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid operation_manager id"
				js, err := json.Marshal(response)
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=12 and id = ?", reqBody.SOMID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + (reqBody.SOMID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid senior operation_manager id"
				js, err := json.Marshal(response)
				//
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
			rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id = ?", reqBody.GFLID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			defer rows.Close()
			if rows.Next() {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + (reqBody.GFLID) + ")"
			} else {
				// showNoProj()
				w.WriteHeader(http.StatusNotFound)
				response := make(map[string]interface{})
				response["success"] = false
				response["message"] = "Invalid gelathi facilitator lead id"
				js, err := json.Marshal(response)
				//
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
				var startDate string
				var endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}

				obj["id"] = id
				obj["name"] = name

				projectArray = append(projectArray, id)

				var tpFilter string
				var tbFilter string
				cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgelathicohorts
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (4, 5, 6, 7, 8, 9)
										WHERE p.id IS NOT NULL
											AND p.id = ?
										GROUP BY p.id`

				if err := DB.QueryRow(cohortCountQuery, id).Scan(&cohortCount); err != nil {
					log.Println("GelathiProgramDashboard", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
					return
				}
				getActualsQuery := `SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM training_participants t1
				JOIN SpoorthiBaselineQuestionnaire t2 ON t1.id = t2.partcipantId
				JOIN project t3 ON t1.project_id = t3.id
				WHERE t1.id = t2.partcipantId
				AND t3.id = ?`

				if err := DB.QueryRow(getActualsQuery, id).Scan(&noofsporthimodulecompleted); err != nil {
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

				if trainid > 0 {
					target := getTrainerTarget(DB, trainid, projectArray)
					obj["target"] = target
					summaryTarget += target
					tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainid)
					tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainid)
				} else {
					target := getTarget(DB, startDate, endDate, projectArray)
					obj["target"] = target
					summaryTarget += target
				}

				actual := getActual(DB, startDate, endDate, projectArray, tpFilter)
				obj["actual"] = actual
				summaryActuals += actual

				Noofcirclemeeting := cohortCount
				obj["NoofGelathiCohorts"] = Noofcirclemeeting
				summarycirclemeeting += Noofcirclemeeting

				obj["Noofsporthisurvey"] = noofsporthisurvey
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = noofsporthimodulecompleted
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				obj["villages"] = newVillageCount(DB, startDate, endDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)

				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				obj["Gelathienrolled"] = noofenrollgelathi
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = getGFData(DB, "", 3, empid)
				summaryofbeehives += obj["Noofbeehives"].(int)

				data = append(data, obj)

			}
		}

		fmt.Println(dateFilters)

		response := make(map[string]interface{})
		response["summary_target"] = summaryTarget
		response["summary_actual"] = summaryActuals
		response["summary_NoofGelathiCohorts"] = summarycirclemeeting
		response["summary_sporthisurvey"] = summarysporthisurvey
		filter += summaryFilter
		response["summary_villages"] = summaryVillages
		response["summary_Noofsporthicompleted"] = summarysporthicompleted
		response["summary_Gelathienrolled"] = summaryGelathiEnrolled
		response["summary_Noofbeehives"] = summaryGelathiEnrolled
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
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=5 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// cohortCountQuery := ""
			// var cohortCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				obj["Gelathienrolled"] = getGelathi(DB, startDate, endDate, projectArray, "", "", "")
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = getGFData(DB, "", 3, empid)
				summaryofbeehives += obj["Noofbeehives"].(int)
				// Noofcirclemeeting := cohortCount
				// obj["NoofGelathiCohorts"] = Noofcirclemeeting
				// summarycirclemeeting += Noofcirclemeeting
				NoofGelathiCohorts := NoofGelathiCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGelathiCohorts
				summarycirclemeeting += NoofGelathiCohorts

				obj["Noofsporthisurvey"] = GetNoOfSporthiSurvey(DB, reqBody.StartDate, reqBody.EndDate, "")
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.trainer_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_NoofGelathiCohorts"] = summarycirclemeeting
			response["summary_noofsporthisurvey"] = summarysporthisurvey
			response["summary_noofsporthicompleted"] = summarysporthicompleted
			response["summary_noofbeehives"] = summaryofbeehives

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
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=6 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// cohortCountQuery := ""
			// var cohortCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				obj["Gelathienrolled"] = getGelathi(DB, startDate, endDate, projectArray, "", "", "")
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = getGFData(DB, "", 3, empid)
				summaryofbeehives += obj["Noofbeehives"].(int)
				// Noofcirclemeeting := cohortCount
				// obj["NoofGelathiCohorts"] = Noofcirclemeeting
				// summarycirclemeeting += Noofcirclemeeting
				NoofGelathiCohorts := NoofGelathiCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGelathiCohorts
				summarycirclemeeting += NoofGelathiCohorts

				obj["Noofsporthisurvey"] = GetNoOfSporthiSurvey(DB, reqBody.StartDate, reqBody.EndDate, "")
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.trainer_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_NoofGelathiCohorts"] = summarycirclemeeting
			response["summary_noofsporthisurvey"] = summarysporthisurvey
			response["summary_noofsporthicompleted"] = summarysporthicompleted
			response["summary_noofbeehives"] = summaryofbeehives

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
			//
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
		rows, err := DB.Query("SELECT id FROM bdms_staff.employee where empRole=13 and id= ?", reqBody.EmpID)
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
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if projectid > 0 {
				rows, err := DB.Query("SELECT id FROM bdms_staff.project where id= ?", reqBody.ProjectID)
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
						"where user_id = " + (reqBody.EmpID) + " and tp.project_id = " + (reqBody.ProjectID) +
						dateFilter +
						" GROUP by tp.project_id"
					summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
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
					"where user_id = " + (reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := DB.Query(query)

			if err != nil {
				log.Println(err)
			}
			defer res.Close()
			// cohortCountQuery := ""
			// var cohortCount int
			// var summary_Noofvyaparcoharts, summaryVyaparenolled, summaryVillages, summaryofvyaparsurvey, summary_vyaparmodulecompleted int
			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Println(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"
				obj["Gelathienrolled"] = getGelathi(DB, startDate, endDate, projectArray, "", "", "")
				summaryGelathiEnrolled += obj["Gelathienrolled"].(int)

				obj["Noofbeehives"] = getGFData(DB, "", 3, empid)
				summaryofbeehives += obj["Noofbeehives"].(int)
				// Noofcirclemeeting := cohortCount
				// obj["NoofGelathiCohorts"] = Noofcirclemeeting
				// summarycirclemeeting += Noofcirclemeeting
				NoofGelathiCohorts := NoofGelathiCohorts(DB, reqBody.StartDate, reqBody.EndDate, "")
				obj["noofVyaparCohorts"] = NoofGelathiCohorts
				summarycirclemeeting += NoofGelathiCohorts

				obj["Noofsporthisurvey"] = GetNoOfSporthiSurvey(DB, reqBody.StartDate, reqBody.EndDate, "")
				summarysporthisurvey += obj["Noofsporthisurvey"].(int)
				obj["Noofsporthicompleted"] = GetNoofSporthiModuleCompleted(DB)
				summarysporthicompleted += obj["Noofsporthicompleted"].(int)

				// summaryEnrolled += obj["enrolled"].(int)
				tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["select_type"] = "1"

				Target := getTrainerTarget(DB, empid, projectArray)
				obj["target"] = Target
				summaryTarget += Target
				filter := " and tp.trainer_id = " + (reqBody.EmpID)

				Actual := getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = Actual
				summaryActuals += Actual

				data = append(data, obj)

			}
			response := make(map[string]interface{})

			response["summary_Gelathienrolled"] = summaryGelathiEnrolled

			tbFilter := fmt.Sprintf(" and tp.user_id = %d", empid)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {

				if i, ok := v.(int); ok {

					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(DB, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			response["summary_Target"] = summaryTarget
			response["summary_NoofGelathiCohorts"] = summarycirclemeeting
			response["summary_noofsporthisurvey"] = summarysporthisurvey
			response["summary_noofsporthicompleted"] = summarysporthicompleted
			response["summary_noofbeehives"] = summaryofbeehives

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
			//
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

// --------------------------------------Funder Vyapar----------------------------------------------------------------
func FunderVyaparD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	// funderid, _ := strconv.Atoi(reqBody.FunderId)
	roleid, _ := strconv.Atoi(reqBody.RoleID)

	empid, _ := strconv.Atoi(reqBody.EmpID) //converting string to int

	data := []interface{}{}
	if roleid == 8 {
		funderListQuery := ""

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee WHERE empRole = 8 AND id = ?", empid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()

		if rows.Next() {
			if empid > 0 {
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE funderID = ?", empid)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder f WHERE funderID = %d", empid)
				} else {
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid funderid"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					return
				}

				if len(funderListQuery) > 0 {
					res, err := DB.Query(funderListQuery)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer res.Close()
					cohortCountQuery := ""
					var cohortCount int

					for res.Next() {
						var projectArray []int
						var funderID int
						var funderName string

						if err := res.Scan(&funderID, &funderName); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofvyaparcohorts
                                            FROM project p
                                            LEFT JOIN tbl_poa tp ON p.id = tp.project_id
                                                AND tp.type = 2
                                                AND tp.session_type IN (16, 17, 18, 19, 20, 21)
                                            WHERE p.funderId IS NOT NULL
                                                AND p.funderId = ?
                                            GROUP BY p.funderId`

						if err := DB.QueryRow(cohortCountQuery, funderID).Scan(&cohortCount); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}

						getProj := "SELECT id FROM project p WHERE funderID = " + strconv.Itoa(funderID)

						projResult, err := DB.Query(getProj)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						defer projResult.Close()

						for projResult.Next() {
							var id int
							if err := projResult.Scan(&id); err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							projectArray = append(projectArray, id)
						}

						if len(projectArray) == 0 {
							obj := map[string]interface{}{
								"id":                        funderID,
								"name":                      funderName,
								"target":                    0,
								"actual":                    0,
								"noofvyaparcoharts":         0,
								"nofvyaparsurvey":           0,
								"noofvyaparmodulecompleted": 0,
								"villages":                  0,
								"startDate":                 "",
								"endDate":                   "",
								"select_type":               "2",
							}
							data = append(data, obj)
							fmt.Println(data)
							continue
						}

						stringSlice := make([]string, len(projectArray))
						for i, val := range projectArray {
							stringSlice[i] = strconv.Itoa(val)
						}

						obj := map[string]interface{}{
							"id":                        funderID,
							"name":                      funderName,
							"target":                    getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":                    getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"noofvyaparcoharts":         cohortCount,
							"nofvyaparsurvey":           GetNoOfVyaparSurvey(DB, reqBody.StartDate, reqBody.EndDate, ""),
							"noofvyaparmodulecompleted": GetNoofVyaparModuleCompleted(DB),
							"villages":                  newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
							"startDate":                 "",
							"endDate":                   "",
							"select_type":               "2",
							"vyaparenrolled":            Vyapar(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						}

						result := map[string]interface{}{
							"actual":                    obj["actual"].(int),
							"endDate":                   obj["endDate"].(string),
							"id":                        obj["id"].(int),
							"name":                      obj["name"].(string),
							"nofvyaparsurvey":           obj["nofvyaparsurvey"].(int),
							"noofvyaparcoharts":         obj["noofvyaparcoharts"].(int),
							"noofvyaparmodulecompleted": obj["noofvyaparmodulecompleted"].(int),
							"select_type":               obj["select_type"].(string),
							"startDate":                 obj["startDate"].(string),
							"target":                    obj["target"].(int),
							"villages":                  obj["villages"].(int),
							"vyaparenrolled":            obj["vyaparenrolled"].(int),
						}

						js, err := json.Marshal(result)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
					}
				}
			}

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid empid"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid roleid"
		js, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
}

// ----------------------------------------------------Funder Gelathi--------------------------------------------------
func FunderGelathiD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	// funderid, _ := strconv.Atoi(reqBody.FunderId)
	roleid, _ := strconv.Atoi(reqBody.RoleID)

	empid, _ := strconv.Atoi(reqBody.EmpID) //converting string to int

	data := []interface{}{}
	if roleid == 8 {
		funderListQuery := ""

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee WHERE empRole = 8 AND id = ?", empid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()

		if rows.Next() {
			if empid > 0 {
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE funderID = ?", empid)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder f WHERE funderID = %d", empid)
				} else {
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid funderid"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					return
				}

				if len(funderListQuery) > 0 {
					res, err := DB.Query(funderListQuery)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer res.Close()
					cohortCountQuery := ""
					var cohortCount int

					for res.Next() {
						var projectArray []int
						var funderID int
						var funderName string

						if err := res.Scan(&funderID, &funderName); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgelathicohorts
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (10, 11, 12, 13, 14, 15)
										WHERE p.funderId IS NOT NULL
											AND p.funderId = ?
										GROUP BY p.funderId`

						if err := DB.QueryRow(cohortCountQuery, funderID).Scan(&cohortCount); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}

						getProj := "SELECT id FROM project p WHERE funderID = " + strconv.Itoa(funderID)

						projResult, err := DB.Query(getProj)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						defer projResult.Close()

						for projResult.Next() {
							var id int
							if err := projResult.Scan(&id); err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							projectArray = append(projectArray, id)
						}

						if len(projectArray) == 0 {
							obj := map[string]interface{}{
								"id":                   funderID,
								"name":                 funderName,
								"target":               0,
								"actual":               0,
								"NoofGelathiCohorts":   0,
								"villages":             0,
								"Gelathienrolled":      0,
								"Noofsporthisurvey":    0,
								"Noofsporthicompleted": 0,
								"Noofbeehives":         0,
								"startDate":            "",
								"endDate":              "",
								"select_type":          "2",
							}
							data = append(data, obj)
							fmt.Println(data)
							continue
						}

						stringSlice := make([]string, len(projectArray))
						for i, val := range projectArray {
							stringSlice[i] = strconv.Itoa(val)
						}

						obj := map[string]interface{}{
							"id":                   funderID,
							"name":                 funderName,
							"target":               getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":               getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"NoofGelathiCohorts":   cohortCount,
							"villages":             newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
							"Gelathienrolled":      getGelathi(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""),
							"Noofsporthisurvey":    GetNoOfSporthiSurvey(DB, reqBody.StartDate, reqBody.EndDate, ""),
							"Noofsporthicompleted": GetNoofSporthiModuleCompleted(DB),
							"Noofbeehives":         getGFData(DB, "", 3, empid),
							"startDate":            "",
							"endDate":              "",
							"select_type":          "2",
						}

						result := map[string]interface{}{
							"actual":               obj["actual"].(int),
							"endDate":              obj["endDate"].(string),
							"id":                   obj["id"].(int),
							"name":                 obj["name"].(string),
							"Noofsporthisurvey":    obj["Noofsporthisurvey"].(int),
							"NoofGelathiCohorts":   obj["NoofGelathiCohorts"].(int),
							"Noofsporthicompleted": obj["Noofsporthicompleted"].(int),
							"Noofbeehives":         obj["Noofbeehives"].(int),
							"select_type":          obj["select_type"].(string),
							"startDate":            obj["startDate"].(string),
							"target":               obj["target"].(int),
							"villages":             obj["villages"].(int),
							"Gelathienrolled":      obj["Gelathienrolled"].(int),
						}

						js, err := json.Marshal(result)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
					}
				}
			}

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid empid"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid roleid"
		js, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
}

// ---------------------------------------------funder SS------------------------------------------------------------
func FunderSSD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	// funderid, _ := strconv.Atoi(reqBody.FunderId)
	roleid, _ := strconv.Atoi(reqBody.RoleID)

	empid, _ := strconv.Atoi(reqBody.EmpID) //converting string to int

	data := []interface{}{}
	if roleid == 8 {
		funderListQuery := ""

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee WHERE empRole = 8 AND id = ?", empid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()

		if rows.Next() {
			if empid > 0 {
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE funderID = ?", empid)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder f WHERE funderID = %d", empid)
				} else {
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid funderid"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					return
				}

				if len(funderListQuery) > 0 {
					res, err := DB.Query(funderListQuery)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer res.Close()
					for res.Next() {
						var projectArray []int
						var funderID int
						var funderName string

						if err := res.Scan(&funderID, &funderName); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}

						getProj := "SELECT id FROM project p WHERE funderID = " + strconv.Itoa(funderID)

						//   err,_ := DB.Query(getProj)
						//    if err != nil {
						// 		log.Println(err)
						// 		w.WriteHeader(http.StatusBadRequest)
						// 		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						// 		return
						// 	}

						projResult, err := DB.Query(getProj)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						defer projResult.Close()

						for projResult.Next() {
							var id int
							if err := projResult.Scan(&id); err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							projectArray = append(projectArray, id)
						}

						if len(projectArray) == 0 {
							obj := map[string]interface{}{
								"id":     funderID,
								"name":   funderName,
								"target": 0,
								"actual": 0,
								"day2":   0,

								"villages":              0,
								"noofbatches":           0,
								"noofselfshakthisurvey": 0,
								"startDate":             "",
								"endDate":               "",
								"select_type":           "2",
							}
							data = append(data, obj)
							fmt.Println(data)
							continue
						}

						stringSlice := make([]string, len(projectArray))
						for i, val := range projectArray {
							stringSlice[i] = strconv.Itoa(val)
						}

						obj := map[string]interface{}{
							"id":                    funderID,
							"name":                  funderName,
							"target":                getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":                getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"day2":                  0,
							"villages":              newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
							"noofbatches":           getTrainingBatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"noofselfshakthisurvey": noOfselfshakthisurvey(DB),

							"startDate":   "",
							"endDate":     "",
							"select_type": "2",
						}
						result := map[string]interface{}{
							"actual":                obj["actual"].(int),
							"endDate":               obj["endDate"].(string),
							"id":                    obj["id"].(int),
							"name":                  obj["name"].(string),
							"noofselfshakthisurvey": obj["noofselfshakthisurvey"].(int),
							"day2":                  obj["day2"].(int),
							"select_type":           obj["select_type"].(string),
							"startDate":             obj["startDate"].(string),
							"target":                obj["target"].(int),
							"villages":              obj["villages"].(int),
							"noofbatches":           obj["noofbatches"].(int),
						}

						js, err := json.Marshal(result)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
					}
				}
			}

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid empid"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid roleid"
		js, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
}

// --------------------------------------------------------Funder Green----------------------------------------------
func FunderGreenD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var reqBody ProjectRequest[any]
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	// funderid, _ := strconv.Atoi(reqBody.FunderId)
	roleid, _ := strconv.Atoi(reqBody.RoleID)

	empid, _ := strconv.Atoi(reqBody.EmpID) //converting string to int

	data := []interface{}{}
	if roleid == 8 {
		funderListQuery := ""

		rows, err := DB.Query("SELECT id FROM bdms_staff.employee WHERE empRole = 8 AND id = ?", empid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()

		if rows.Next() {
			if empid > 0 {
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE funderID = ?", empid)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder f WHERE funderID = %d", empid)
				} else {
					w.WriteHeader(http.StatusNotFound)
					response := make(map[string]interface{})
					response["success"] = false
					response["message"] = "Invalid funderid"
					js, err := json.Marshal(response)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					return
				}

				if len(funderListQuery) > 0 {
					res, err := DB.Query(funderListQuery)
					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
						return
					}
					defer res.Close()
					cohortCountQuery := ""
					var cohortCount int

					for res.Next() {
						var projectArray []int
						var funderID int
						var funderName string

						if err := res.Scan(&funderID, &funderName); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						cohortCountQuery = `SELECT COALESCE(COUNT(tp.session_type), 0) AS noofgreencohorts
										FROM project p
										LEFT JOIN tbl_poa tp ON p.id = tp.project_id
											AND tp.type = 2
											AND tp.session_type IN (10, 11, 12, 13, 14, 15)
										WHERE p.funderId IS NOT NULL
											AND p.funderId = ?
										GROUP BY p.funderId`

						if err := DB.QueryRow(cohortCountQuery, funderID).Scan(&cohortCount); err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}

						getProj := "SELECT id FROM project p WHERE funderID = " + strconv.Itoa(funderID)

						projResult, err := DB.Query(getProj)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						defer projResult.Close()

						for projResult.Next() {
							var id int
							if err := projResult.Scan(&id); err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							projectArray = append(projectArray, id)
						}

						if len(projectArray) == 0 {
							obj := map[string]interface{}{
								"id":                       funderID,
								"name":                     funderName,
								"target":                   0,
								"actual":                   0,
								"noofgreencoharts":         0,
								"nofgreensurvey":           0,
								"noofgreenmodulecompleted": 0,
								"villages":                 0,
								"startDate":                "",
								"endDate":                  "",
								"select_type":              "2",
							}
							data = append(data, obj)
							fmt.Println(data)
							continue
						}

						stringSlice := make([]string, len(projectArray))
						for i, val := range projectArray {
							stringSlice[i] = strconv.Itoa(val)
						}

						obj := map[string]interface{}{
							"id":               funderID,
							"name":             funderName,
							"target":           getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":           getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"noofgreencoharts": cohortCount,

							"nofgreensurvey":           GetNoOfgreenSurvey(DB, reqBody.StartDate, reqBody.EndDate, ""),
							"noofgreenmodulecompleted": GetNoofGreenModuleCompleted(DB),
							"villages":                 newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas

							"startDate":     "",
							"endDate":       "",
							"select_type":   "2",
							"greenenrolled": greenMotivators(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						}

						result := map[string]interface{}{
							"actual":                   obj["actual"].(int),
							"endDate":                  obj["endDate"].(string),
							"id":                       obj["id"].(int),
							"name":                     obj["name"].(string),
							"nofgreensurvey":           obj["nofgreensurvey"].(int),
							"noofgreencoharts":         obj["noofgreencoharts"].(int),
							"noofgreenmodulecompleted": obj["noofgreenmodulecompleted"].(int),
							"select_type":              obj["select_type"].(string),
							"startDate":                obj["startDate"].(string),
							"target":                   obj["target"].(int),
							"villages":                 obj["villages"].(int),
							"greenenrolled":            obj["greenenrolled"].(int),
						}

						js, err := json.Marshal(result)
						if err != nil {
							log.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(js)
					}
				}
			}

		} else {
			w.WriteHeader(http.StatusNotFound)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid empid"
			js, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		response := make(map[string]interface{})
		response["success"] = false
		response["message"] = "Invalid roleid"
		js, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
}
