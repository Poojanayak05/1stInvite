package Dashboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
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

//--------------------------------------Funder Vyapar----------------------------------------------------------------
func FunderVyaparD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
				var email string
				err := DB.QueryRow("SELECT officeMailId FROM employee where id=?", empid).Scan(&email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE emailID = ?", email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder WHERE emailID = '%s'", email)
				} else {
					w.WriteHeader(http.StatusBadRequest)
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

							js, err := json.Marshal(data)
							if err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							w.Write(js)
							continue
						} else {
							fmt.Println(">>>>>>>>>>", data)

							fmt.Println(">>>>>>>>>>.,", projectArray)
							obj := map[string]interface{}{
								"id":                        funderID,
								"name":                      funderName,
								"target":                    getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
								"actual":                    getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"noofvyaparcoharts":         noofvyaparcoharts(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"nofvyaparsurvey":           noofvyaparsurvey(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"noofvyaparmodulecompleted": noofvyaparmodulecompleted(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"villages":                  VyaparnewVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
								"startDate":                 "",
								"endDate":                   "",
								"select_type":               "2",
								"vyaparenrolled":            vyaparenrolled(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
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

//----------------------------------------------------Funder Gelathi--------------------------------------------------
func FunderGelathiD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
				var email string
				err := DB.QueryRow("SELECT officeMailId FROM employee where id=?", empid).Scan(&email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE emailID = ?", email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder WHERE emailID = '%s'", email)
				} else {
					w.WriteHeader(http.StatusBadRequest)
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

							js, err := json.Marshal(data)
							if err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							w.Write(js)
							continue
						}

						obj := map[string]interface{}{
							"id":                   funderID,
							"name":                 funderName,
							"target":               getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":               getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"NoofGelathiCohorts":   NoofGreenCohorts(DB, reqBody.StartDate, reqBody.EndDate, ""),
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

//---------------------------------------------funder SS------------------------------------------------------------
func FunderSSD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
				var email string
				err := DB.QueryRow("SELECT officeMailId FROM employee where id=?", empid).Scan(&email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE emailID = ?", email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder WHERE emailID = '%s'", email)
				} else {
					w.WriteHeader(http.StatusBadRequest)
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
						fmt.Println(data)

						fmt.Println(projectArray)
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
							js, err := json.Marshal(data)
							if err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							w.Write(js)
							continue
						} else {

							obj := map[string]interface{}{
								"id":                    funderID,
								"name":                  funderName,
								"target":                getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
								"actual":                getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"villages":              newVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas
								"noofbatches":           noofbatches(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
								"noofselfshakthisurvey": noOfselfshakthisurvey(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),

								"startDate":   "",
								"endDate":     "",
								"select_type": "2",
							}
							var summaryDay1 int
							var day2Turnout float64
							day1Count, _ := getDay1Count(DB, reqBody.StartDate, reqBody.EndDate, projectArray, "")
							summaryDay1 += day1Count
							if day1Count > 0 {
								day2Turnout = float64(obj["actual"].(int)) / float64(day1Count)
								obj["day2"] = int(math.Round(day2Turnout * 100))
							} else {
								obj["day2"] = 0
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

							w.Write(js)
						}
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

//--------------------------------------------------------Funder Green----------------------------------------------
func FunderGreenD(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

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
				var email string
				err := DB.QueryRow("SELECT officeMailId FROM employee where id=?", empid).Scan(&email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				rows, err := DB.Query("SELECT funderID FROM bdms_staff.funder WHERE emailID = ?", email)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
					return
				}
				defer rows.Close()

				if rows.Next() {
					funderListQuery = fmt.Sprintf("SELECT funderID as id ,funderName as name FROM funder WHERE emailID = '%s'", email)
				} else {
					w.WriteHeader(http.StatusBadRequest)
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

							js, err := json.Marshal(data)
							if err != nil {
								log.Println(err)
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
								return
							}
							w.Write(js)
							continue
						}

						obj := map[string]interface{}{
							"id":               funderID,
							"name":             funderName,
							"target":           getTarget(DB, reqBody.StartDate, reqBody.EndDate, projectArray),
							"actual":           getActual(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"noofgreencoharts": noofgreencoharts(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),

							"nofgreensurvey":           noofgreensurvey(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"noofgreenmodulecompleted": noofgreenmodulecompleted(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
							"villages":                 GreennewVillageCount(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""), // New village count function anas

							"startDate":     "",
							"endDate":       "",
							"select_type":   "2",
							"greenenrolled": greenenrolled(DB, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
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
