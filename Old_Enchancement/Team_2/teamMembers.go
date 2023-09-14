package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type TeamRequest struct {
	Empid  int    `json:"emp_id"`
	Roleid int    `json:"role_id"`
	Date   string `json:"date"`
}

type TeamEmployee struct {
	ID             int    `json:"id"`
	EmpRole        int    `json:"empRole"`
	CountryId      string `json:"CountryId"`
	Gender         string `json:"gender"`
	DOJ            string `json:"DOJ"`
	Officemailid   string `json:"officemailid"`
	Personalmailid string `json:"personalmailid"`
	Contactnum     string `json:"contactnum"`
	Address        string `json:"address"`
	Address2       string `json:"address2"`
	Address3       string `json:"address3"`
	Pincode        string `json:"pincode"`
	Licensenumber  string `json:"licensenumber"`
	Created_at     string `json:"created_at"`
	Created_by     string `json:"created_by"`
	Lastupdated_at string `json:"lastupdated_at"`
	Lastupdated_by string `json:"lastupdated_by"`
	Worknum        string `json:"Worknum"`
	SupervisorID   int    `json:"supervisorId"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePic     string `json:"profile_pic"`
	Status         string `json:"Status"`
	RoleName       string `json:"roleName"`
}

type TeamResponse struct {
	FullName    string `json:"fullName"`
	ID          string `json:"id"`
	EmpRole     string `json:"empRole"`
	Designation string `json:"designation"`
	Profile_pic string `json:"profile_pic"`
	IsApprove   bool   `json:"isApprove"`
	Message     string `json:"message"`
	Verify      bool   `json:"verify"`
	TotalTA     string `json:"totalTa"`
	Status      string `json:"status"`
	TaName      string `json:"TaName"`
}

func TeamMembers(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	var p TeamRequest
	var employee TeamEmployee
	var res TeamResponse
	var team []TeamResponse
	var isApprove = true
	var verify = true
	var message string
	var totalTa float64
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})
		return
	}

	if p.Empid > 0 && p.Roleid > 0 {
		//----------------------zeroth row--------------------
		r := "SELECT e.id,e.countryid,e.first_name,COALESCE(e.last_name,''),e.gender,COALESCE(COALESCE(e.doj,''),''),COALESCE(e.officemailid,''),COALESCE(e.personalMailId,''),COALESCE(e.contactNum,''),COALESCE(e.address,''),COALESCE(e.address2,''),COALESCE(e.address3,''),COALESCE(e.pincode,''),e.empRole,e.license_number,e.supervisorid,COALESCE(e.profile_pic,''),e.status,e.createdAt,COALESCE(e.createdBy,''),e.lastUpdatedAt,COALESCE(e.lastUpdatedby,''),COALESCE(e.workNum,''), rm.roleName from employee e join roles_Master rm on e.empRole =rm.id where e.status=1 and e.supervisorId= ?"

		rows, err := DB.Query(r, p.Empid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer rows.Close()
		for rows.Next() {

			if err := rows.Scan(&employee.ID, &employee.CountryId, &employee.FirstName, &employee.LastName, &employee.Gender, &employee.DOJ, &employee.Officemailid, &employee.Personalmailid, &employee.Contactnum, &employee.Address, &employee.Address2, &employee.Address3, &employee.Pincode, &employee.EmpRole, &employee.Licensenumber, &employee.SupervisorID, &employee.ProfilePic, &employee.Status, &employee.Created_at, &employee.Created_by, &employee.Lastupdated_at, &employee.Lastupdated_by, &employee.Worknum, &employee.RoleName); err != nil {
				log.Println(err, "line 108")
				//continue
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}
			dateFilter := fmt.Sprintf("MONTH(entry_date) = MONTH('%s') and ", p.Date)
			re := DB.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(total_ta),0) AS 'total_ta' FROM travelling_allowance WHERE %s emp_id=%d", dateFilter, employee.ID))

			if err := re.Scan(&totalTa); err != nil {
				//log.Println(err)
				//continue
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}

			var status string
			t, err := DB.Query(fmt.Sprintf("SELECT DISTINCT status FROM travelling_allowance ta WHERE %s emp_id=%d ORDER BY status", dateFilter, employee.ID))
			if err != nil {
				log.Println(err)
				//continue
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}

			defer t.Close()

			for t.Next() {
				var ttStatus string
				if err := t.Scan(&ttStatus); err != nil {
					//log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}

				switch ttStatus {
				case "0":
					status = ttStatus
				case "1":
					if status != "0" {
						status = ttStatus
					}
				case "2":
					if status != "0" && status != "1" {
						status = ttStatus
					}
				case "3":
					if status != "0" && status != "1" && status != "2" {
						status = ttStatus
					}
				case "4":
					if status != "0" && status != "1" && status != "2" && status != "3" {
						status = ttStatus
					}
				}
			}

			if employee.SupervisorID == p.Empid {
				switch p.Roleid {
				case 3:
					isApprove = true
				case 12:
					if totalTa < 15000 {
						isApprove = true
					}
				default:
					if totalTa < 10000 {
						isApprove = true
					}
				}
			}

			if !isApprove {
				if totalTa != 0 {
					if totalTa > 10000 && totalTa < 15000 {
						message = "Approval pending from senior operation manager"
					} else if totalTa > 15000 {
						message = "Approval pending from Program manager"
					}
				}
			}

			v, err := DB.Query(fmt.Sprintf("SELECT DISTINCT verifyed_by FROM travelling_allowance ta WHERE %s emp_id=%d ORDER BY verifyed_by", dateFilter, employee.ID))
			if err != nil {
				// log.Println(err)
				// return
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
				//continue
			}
			defer v.Close()

			for v.Next() {
				var verifyBy sql.NullString
				if err := v.Scan(&verifyBy); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

					return
					//continue
				}

				if !verifyBy.Valid {
					verify = true
				} else {
					verify = false
				}
			}

			if p.Roleid == 13 || p.Roleid == 5 || p.Roleid == 6 {
				isApprove = false
			}

			if totalTa != 0 {
				res.FullName = employee.FirstName + " " + employee.LastName
				res.ID = fmt.Sprint(employee.ID)
				res.EmpRole = fmt.Sprint(employee.EmpRole)
				res.Designation = employee.RoleName
				res.Profile_pic = employee.ProfilePic
				res.IsApprove = isApprove
				res.Message = message
				res.Verify = verify
				res.TotalTA = fmt.Sprint(totalTa)
				res.Status = status
				res.TaName = fmt.Sprintf("%s %s_%d", employee.FirstName, employee.LastName, int(totalTa))
				team = append(team, res)
			}

		}
		//------------------------first row---------------------
		r1, err := DB.Query("select e.id,e.countryid,e.first_name,COALESCE(e.last_name,''),e.gender,COALESCE(e.doj,''),COALESCE(e.officemailid,''),COALESCE(e.personalMailId,''),COALESCE(e.contactNum,''),COALESCE(e.address,''),COALESCE(e.address2,''),COALESCE(e.address3,''),COALESCE(e.pincode,''),e.empRole,e.license_number,e.supervisorid,COALESCE(e.profile_pic,''),e.status,e.createdAt,COALESCE(e.createdBy,''),e.lastUpdatedAt,COALESCE(e.lastUpdatedby,''),COALESCE(e.workNum,'') from employee e where status=1 and supervisorId=?", p.Empid)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
			return
		}
		defer r1.Close()
		for r1.Next() {
			err := r1.Scan(&employee.ID, &employee.CountryId, &employee.FirstName, &employee.LastName, &employee.Gender, &employee.DOJ, &employee.Officemailid, &employee.Personalmailid, &employee.Contactnum, &employee.Address, &employee.Address2, &employee.Address3, &employee.Pincode, &employee.EmpRole, &employee.Licensenumber, &employee.SupervisorID, &employee.ProfilePic, &employee.Status, &employee.Created_at, &employee.Created_by, &employee.Lastupdated_at, &employee.Lastupdated_by, &employee.Worknum)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})

				return
			}
			r2, err := DB.Query("select e.id,e.countryid,e.first_name,COALESCE(e.last_name,''),e.gender,COALESCE(e.doj,''),COALESCE(e.officemailid,''),COALESCE(e.personalMailId,''),COALESCE(e.contactNum,''),COALESCE(e.address,''),COALESCE(e.address2,''),COALESCE(e.address3,''),COALESCE(e.pincode,''),e.empRole,e.license_number,e.supervisorid,COALESCE(e.profile_pic,''),e.status,e.createdAt,COALESCE(e.createdBy,''),e.lastUpdatedAt,COALESCE(e.lastUpdatedby,''),COALESCE(e.workNum,''),rm.roleName from employee e join roles_Master rm on e.empRole =rm.id where e.status=1 and e.supervisorId= ?", employee.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
				return
			}
			defer r2.Close()
			for r2.Next() {
				if err := r2.Scan(&employee.ID, &employee.CountryId, &employee.FirstName, &employee.LastName, &employee.Gender, &employee.DOJ, &employee.Officemailid, &employee.Personalmailid, &employee.Contactnum, &employee.Address, &employee.Address2, &employee.Address3, &employee.Pincode, &employee.EmpRole, &employee.Licensenumber, &employee.SupervisorID, &employee.ProfilePic, &employee.Status, &employee.Created_at, &employee.Created_by, &employee.Lastupdated_at, &employee.Lastupdated_by, &employee.Worknum, &employee.RoleName); err != nil {

					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				dateFilter0 := fmt.Sprintf("MONTH(entry_date) = MONTH('%s') and ", p.Date)
				re0 := "SELECT COALESCE(SUM(total_ta),0) AS total_ta FROM travelling_allowance WHERE status IN (0,4) AND " + fmt.Sprint(dateFilter0) + " emp_id=" + fmt.Sprint(employee.ID)
				var rrr0 struct {
					TotalTA float64
				}
				err := DB.QueryRow(re0).Scan(&rrr0.TotalTA)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}

				status := ""
				t, err := DB.Query("SELECT DISTINCT status FROM travelling_allowance ta WHERE " + fmt.Sprint(dateFilter0) + " emp_id=" + fmt.Sprint(employee.ID) + " ORDER BY status")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				defer t.Close()
				for t.Next() {
					var tt struct {
						Status string
					}
					if err := t.Scan(&tt.Status); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}
					if tt.Status == "0" {
						status = tt.Status
					} else if tt.Status == "1" && status != "0" {
						status = tt.Status
					} else if tt.Status == "2" && status != "0" && status != "1" {
						status = tt.Status
					} else if tt.Status == "3" && status != "0" && status != "1" && status != "2" {
						status = tt.Status
					} else if tt.Status == "4" && status != "0" && status != "1" && status != "2" && status != "3" {
						status = tt.Status
					}
				}

				v, err := DB.Query("SELECT DISTINCT verifyed_by FROM travelling_allowance ta WHERE " + fmt.Sprint(dateFilter0) + " emp_id=" + fmt.Sprint(employee.ID) + " ORDER BY verifyed_by")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				defer v.Close()
				for v.Next() {
					var vv struct {
						VerifiedBy sql.NullInt64
					}
					if err := v.Scan(&vv.VerifiedBy); err != nil {
						log.Println("ERROR>>", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
						return
					}
					if !vv.VerifiedBy.Valid {
						verify = true
					} else {
						verify = false
					}
				}

				if p.Roleid == 4 {
					if rrr0.TotalTA < 10000 && rrr0.TotalTA > 0 {
						isApprove = true
					} else {
						isApprove = false
					}
				} else if p.Roleid == 12 {
					if rrr0.TotalTA < 15000 && rrr0.TotalTA > 0 {
						isApprove = true
					} else {
						isApprove = false
					}
				} else {
					if rrr0.TotalTA > 15000 && rrr0.TotalTA > 0 {
						isApprove = true
					} else {
						isApprove = false
					}
				}
				if !isApprove {
					if totalTa != 0 {
						if totalTa > 10000 && totalTa < 15000 {
							message = "Approval pending from senior operation manager"
						} else if totalTa > 15000 {
							message = "Approval pending from Program manager"
						}
					}
				} else {
					message = ""
				}
				if rrr0.TotalTA > 0 {
					team = append(team, TeamResponse{
						FullName:    employee.FirstName + " " + employee.LastName,
						ID:          fmt.Sprint(employee.ID),
						Designation: employee.RoleName,
						EmpRole:     fmt.Sprint(employee.EmpRole),
						Profile_pic: employee.ProfilePic,
						IsApprove:   isApprove,
						Message:     message,
						Verify:      verify,
						TotalTA:     fmt.Sprint(rrr0.TotalTA),
						Status:      status,
						TaName:      employee.FirstName + " " + employee.LastName + "_" + fmt.Sprint(int(rrr0.TotalTA)),
					})
				}
				//-------------------------------third row-----------------------------

				r3, err := DB.Query("SELECT e.id,e.countryid,e.first_name,COALESCE(e.last_name,''),e.gender,COALESCE(e.doj,''),COALESCE(e.officemailid,''),COALESCE(COALESCE(e.personalMailId,''),''),COALESCE(e.contactNum,''),COALESCE(e.address,''),COALESCE(e.address2,''),COALESCE(e.address3,''),COALESCE(e.pincode,''),e.empRole,e.license_number,e.supervisorid,COALESCE(e.profile_pic,''),e.status,e.createdAt,COALESCE(e.createdBy,''),e.lastUpdatedAt,COALESCE(e.lastUpdatedby,''),COALESCE(e.workNum,''),rm.roleName FROM employee e JOIN roles_Master rm ON e.empRole = rm.id WHERE e.status=1 AND e.supervisorId=?", employee.ID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
					return
				}
				defer r3.Close()
				for r3.Next() {
					if err := r3.Scan(&employee.ID, &employee.CountryId, &employee.FirstName, &employee.LastName, &employee.Gender, &employee.DOJ, &employee.Officemailid, &employee.Personalmailid, &employee.Contactnum, &employee.Address, &employee.Address2, &employee.Address3, &employee.Pincode, &employee.EmpRole, &employee.Licensenumber, &employee.SupervisorID, &employee.ProfilePic, &employee.Status, &employee.Created_at, &employee.Created_by, &employee.Lastupdated_at, &employee.Lastupdated_by, &employee.Worknum, &employee.RoleName); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}

					dateFilter1 := fmt.Sprintf("MONTH(entry_date) = MONTH('%s') and ", p.Date)
					re1 := DB.QueryRow("SELECT COALESCE(SUM(total_ta),0) AS total_ta FROM travelling_allowance WHERE status IN (0,4) AND "+dateFilter1+" emp_id=?", employee.ID)
					if err := re1.Scan(&rrr0.TotalTA); err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
						return
					}

					status := ""
					t, err := DB.Query("SELECT DISTINCT status FROM travelling_allowance ta WHERE "+dateFilter0+" emp_id=? ORDER BY status", employee.ID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}
					defer t.Close()
					for t.Next() {
						var tt struct {
							Status string
						}
						if err := t.Scan(&tt.Status); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
							return
						}
						if tt.Status == "0" {
							status = tt.Status
						} else if tt.Status == "1" && status != "0" {
							status = tt.Status
						} else if tt.Status == "2" && status != "0" && status != "1" {
							status = tt.Status
						} else if tt.Status == "3" && status != "0" && status != "1" && status != "2" {
							status = tt.Status
						} else if tt.Status == "4" && status != "0" && status != "1" && status != "2" && status != "3" {
							status = tt.Status
						}
					}

					v, err := DB.Query("SELECT DISTINCT verifyed_by FROM travelling_allowance ta WHERE "+dateFilter0+" emp_id=? ORDER BY verifyed_by", employee.ID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}
					defer v.Close()
					for v.Next() {
						var vv struct {
							VerifiedBy sql.NullInt64
						}
						if err := v.Scan(&vv.VerifiedBy); err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
							return
						}
						if !vv.VerifiedBy.Valid {
							verify = true
						} else {
							verify = false
						}
					}

					isApprove := true
					if p.Roleid == 4 {
						if rrr0.TotalTA < 10000 && rrr0.TotalTA > 0 {
							isApprove = true
						} else {
							isApprove = false
						}
					} else if p.Roleid == 12 {
						if rrr0.TotalTA < 15000 && rrr0.TotalTA > 0 {
							isApprove = true
						} else {
							isApprove = false
						}
					} else {
						if rrr0.TotalTA > 15000 && rrr0.TotalTA > 0 {
							isApprove = true
						} else {
							isApprove = false
						}
					}

					message := ""
					if !isApprove {
						if rrr0.TotalTA > 10000 && rrr0.TotalTA < 15000 {
							message = "Approval pending from senior operation manager"
						}
						if rrr0.TotalTA > 15000 {
							message = "Approval pending from senior Program manager"
						}
					}

					if rrr0.TotalTA > 0 {
						team = append(team, TeamResponse{
							FullName:    employee.FirstName + " " + employee.LastName,
							ID:          fmt.Sprint(employee.ID),
							Designation: employee.RoleName,
							EmpRole:     fmt.Sprint(employee.EmpRole),
							Profile_pic: employee.ProfilePic,
							IsApprove:   isApprove,
							Message:     message,
							Verify:      verify,
							TotalTA:     fmt.Sprint(rrr0.TotalTA),
							Status:      status,
							TaName:      employee.FirstName + " " + employee.LastName + "_" + fmt.Sprint(int(rrr0.TotalTA)),
						})
					}

					r4, err := DB.Query("SELECT e.*, rm.roleName FROM employee e JOIN roles_Master rm ON e.empRole = rm.id WHERE e.status=1 AND e.supervisorId=?", employee.ID)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
						return
					}
					defer r4.Close()
					for r4.Next() {
						if err := r4.Scan(&employee.ID, &employee.CountryId, &employee.FirstName, &employee.LastName, &employee.Gender, &employee.DOJ, &employee.Officemailid, &employee.Personalmailid, &employee.Contactnum, &employee.Address, &employee.Address2, &employee.Address3, &employee.Pincode, &employee.EmpRole, &employee.Licensenumber, &employee.SupervisorID, &employee.ProfilePic, &employee.Status, &employee.Created_at, &employee.Created_by, &employee.Lastupdated_at, &employee.Lastupdated_by, &employee.Worknum, &employee.RoleName); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
							return
						}
						dateFilter4 := fmt.Sprintf("MONTH(entry_date) = MONTH('%s') and ", p.Date)
						re4 := DB.QueryRow("SELECT COALESCE(SUM(total_ta),0) AS total_ta FROM travelling_allowance WHERE status IN (0,4) AND "+dateFilter4+" emp_id=?", employee.ID)
						if err := re4.Scan(&rrr0.TotalTA); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
							return
						}

						status = ""
						t4, err := DB.Query("SELECT DISTINCT status FROM travelling_allowance ta WHERE "+dateFilter0+" emp_id=? ORDER BY status", employee.ID)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
							return
						}
						defer t4.Close()
						for t4.Next() {
							var tt4 struct {
								Status string
							}
							if err := t4.Scan(&tt4.Status); err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
								json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
								return
							}
							if tt4.Status == "0" {
								status = tt4.Status
							} else if tt4.Status == "1" && status != "0" {
								status = tt4.Status
							} else if tt4.Status == "2" && status != "0" && status != "1" {
								status = tt4.Status
							} else if tt4.Status == "3" && status != "0" && status != "1" && status != "2" {
								status = tt4.Status
							} else if tt4.Status == "4" && status != "0" && status != "1" && status != "2" && status != "3" {
								status = tt4.Status
							}
						}

						v, err := DB.Query("SELECT DISTINCT verifyed_by FROM travelling_allowance ta WHERE "+dateFilter0+" emp_id=? ORDER BY verifyed_by", employee.ID)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
							return
						}
						defer v.Close()
						for v.Next() {
							var vv struct {
								VerifyedBy sql.NullString
							}
							if err := v.Scan(&vv.VerifyedBy); err != nil {
								w.WriteHeader(http.StatusBadRequest)
								json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
								return
							}
							if !vv.VerifyedBy.Valid {
								verify = true
							} else {
								verify = false
							}
						}

						isApprove = true
						if p.Roleid == 4 {
							if rrr0.TotalTA < 10000 && rrr0.TotalTA > 0 {
								isApprove = true
							} else {
								isApprove = false
							}
						} else if p.Roleid == 12 {
							if rrr0.TotalTA < 15000 && rrr0.TotalTA > 0 {
								isApprove = true
							} else {
								isApprove = false
							}
						} else {
							if rrr0.TotalTA > 15000 && rrr0.TotalTA > 0 {
								isApprove = true
							} else {
								isApprove = false
							}
						}

						message := ""
						if !isApprove {
							if rrr0.TotalTA > 10000 && rrr0.TotalTA < 15000 {
								message = "Approval pending from senior operation manager"
							}
							if rrr0.TotalTA > 15000 {
								message = "Approval pending from senior Program manager"
							}
						}

						if rrr0.TotalTA > 0 {
							team = append(team, TeamResponse{
								FullName:    employee.FirstName + " " + employee.LastName,
								ID:          fmt.Sprint(employee.ID),
								Designation: employee.RoleName,
								EmpRole:     fmt.Sprint(employee.EmpRole),
								Profile_pic: employee.ProfilePic,
								IsApprove:   isApprove,
								Message:     message,
								Verify:      verify,
								TotalTA:     fmt.Sprint(rrr0.TotalTA),
								Status:      status,
								TaName:      employee.FirstName + " " + employee.LastName + "_" + fmt.Sprint(int(rrr0.TotalTA)),
							})
						}
					}
				}

			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Team Data", "data": team})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Send Required Paramter", "data": team})
	}

}
