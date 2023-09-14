package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetAllPeople(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	fmt.Println("inside method")
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// Check if the request method is GET
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	type Request struct {
		Search    string `json:"search"`
		User_Id   string `json:"user_id"`
		Role_Id   string `json:"role_id"`
		Filter_Id string `json:"filter_id"`
		Type      string `json:"type"`
		PageNum   string `json:"pagenum"`
	}

	type Response2 struct {
		ID             string                   `json:"id"`
		FirstName      string                   `json:"first_name"`
		LastName       string                   `json:"last_name"`
		RoleID         string                   `json:"role_id"`
		CountryID      string                   `json:"countryID"`
		PersonalMailId string                   `json:"personalMailId"`
		CreatedBy      string                   `json:"createdBy"`
		LastUpdatedBy  string                   `json:"lastUpdatedBy"`
		OfficeMailID   *string                  `json:"officeMailId"`
		ContactNum     *string                  `json:"contactNum"`
		WorkNum        *string                  `json:"workNum"`
		RoleName       string                   `json:"role_name"`
		ProfilePic     string                   `json:"profile_pic"`
		Doj            *string                  `json:"doj"`
		Address        *string                  `json:"address"`
		Address2       *string                  `json:"address2"`
		Address3       *string                  `json:"address3"`
		SupervisorName *string                  `json:"supervisorName"`
		SupervisorId   string                   `json:"supervisorId"`
		Status         string                   `json:"status"`
		Pincode        *string                  `json:"pincode"`
		License_Number *string                  `json:"license_number"`
		Gender         *string                  `json:"gender"`
		ProjectList    []map[string]interface{} `json:"project_list"`
	}

	var request Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status": "Bad Request", "Code": "400"})
		return
	}
	var pageno, no_of_records, offset, totalPages, totalRows int
	var searchFilter, role, condition_funder, condition_partner, total_pages_sql, fields, query, getIds, idsString string
	var dataSlice2 []Response2
	fmt.Println(idsString)
	fmt.Println(totalPages)

	Filter_ID, _ := strconv.Atoi(request.Filter_Id)
	Role_ID, _ := strconv.Atoi(request.Role_Id)
	User_ID, _ := strconv.Atoi(request.User_Id)

	// to handle page number
	if request.PageNum != "" {
		pageno, _ = strconv.Atoi(request.PageNum)
	} else {
		pageno = 1
	}

	no_of_records = 10000000
	offset = (pageno - 1) * no_of_records
	//to handle search parameter
	// searchParam, ok := data["search"].(string)
	// if !ok {
	// 	searchParam = ""
	// }

	searchFilter = ""
	if len(request.Search) > 0 {
		searchFilter = fmt.Sprintf("AND (emp.first_name LIKE '%%%s%%' OR emp.last_name LIKE '%%%s%%')", request.Search, request.Search)
		fmt.Println("searchFilter", searchFilter)
	}
	// if len(searchParam) > 0 {
	// 	searchFilter = fmt.Sprintf("and (emp.first_name like '%s' or emp.last_name like '%s')", searchParam, searchParam)
	// }
	//to handle different scenarios based on user and role IDs
	if request.User_Id != "" && request.Role_Id != "" {
		if request.Role_Id == "1" || request.Role_Id == "2" || request.Role_Id == "3" || request.Role_Id == "4" || request.Role_Id == "5" || request.Role_Id == "12" {
			if Filter_ID == 0 {
				// If Filter_ID is 0, retrieve employees with specific empRoles
				role = "and emp.empRole IN ('1','2','3','4','5','6','7','8','9','11','12','13')"
			} else if Filter_ID == 10 {
				// If Filter_ID is 10, retrieve employees with specific empRoles (excluding empRole 6)
				role = "and emp.empRole IN ('1','2','3','4','11','12','13')"
			} else {
				fmt.Println("line 117")
				// If Filter_ID is any other value, retrieve employees with a specific empRole
				role = fmt.Sprintf("and emp.empRole = %d", Filter_ID)
			}

			if Filter_ID == 8 || Filter_ID == 9 {

				// condition_funder = ""
				// condition_partner = ""

				fmt.Println("entering line 123")
				if Role_ID == 3 {
					//fmt.Println("line 129")
					// If Role_ID is 3 and Filter_ID is 8 or 9, retrieve project IDs for the user
					project_ids := getEmpPrj(DB, User_ID, Role_ID, Filter_ID)
					//	fmt.Println("project_ids line 125", project_ids)
					filePath := "log.txt"
					file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
					if err != nil {
						fmt.Println("Failed to open file:", err)
						return
					}

					defer file.Close()

					projectIDs := project_ids
					_, err = fmt.Fprintln(file, projectIDs)
					if err != nil {
						// fmt.Println("line 144", err)
						fmt.Println("Failed to write to file:", err)
						return
					}
					ids := projectIDs
					condition_funder = fmt.Sprintf(" and funderID IN (%s)", ids)
					condition_partner = fmt.Sprintf("and partnerID IN (%s)", ids)
				}
				fmt.Println("line 147")

				if Filter_ID == 8 {

					//fmt.Println(" inside 8 line 148 ")

					// If Filter_ID is 8, retrieve funders based on the condition
					total_pages_sql = fmt.Sprintf("SELECT COUNT(*) FROM funder where status = 1 %s", condition_funder)
					fmt.Println("total_pages_sql", total_pages_sql)
					fields = "funderID as id, funderName AS first_name, '' AS last_name, emailID as officeMailId, workPhone as contactNum, mobilePhone as workNum, 8 as role_id, 'funder' as role_name, 'https://lh3.googleusercontent.com/a/AATXAJxZ1KGlNFsyK9XMmPI4whQUhCJt9yCJYm1Qor4h=s96-c' as profile_pic, address, status, pincode, '' as gender"

					query = fmt.Sprintf("SELECT %s FROM funder where status = 1 %s ORDER BY funderName LIMIT %d, %d", fields, condition_funder, offset, no_of_records)
				} else if Filter_ID == 9 {
					//fmt.Println("inside 9")
					// If Filter_ID is 9, retrieve partners based on the condition
					total_pages_sql = fmt.Sprintf("SELECT COUNT(*) FROM partner where status = 1 %s", condition_partner)

					fields = "partnerID as id, partnerName AS first_name, '' AS last_name, emailID as officeMailId, phone as contactNum, mobilePhone as workNum, 9 as role_id, 'partner' as role_name, 'https://lh3.googleusercontent.com/a/AATXAJxZ1KGlNFsyK9XMmPI4whQUhCJt9yCJYm1Qor4h=s96-c' as profile_pic, address, status, pincode, '' as gender"

					query = fmt.Sprintf("SELECT %s FROM partner where status = 1 %s ORDER BY partnerName LIMIT %d, %d", fields, condition_partner, offset, no_of_records)
				}

				rows, err := DB.Query(total_pages_sql)
				if err != nil {

					//fmt.Println("line 170", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request11", "success": false, "error": err})
					return
				}
				var totalRows int
				for rows.Next() {
					err := rows.Scan(&totalRows)
					if err != nil {
						//fmt.Println("err line 177", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
				}
				totalPages = int(math.Ceil(float64(totalRows) / float64(no_of_records)))
				fmt.Println(totalPages)

				rows1, err := DB.Query(query)
				if err != nil {
					//fmt.Println("line 190", err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}

				defer rows.Close()

				var dataSlice []map[string]interface{}

				for rows1.Next() {
					var Id, First_name, Last_name, OfficeMailID, ContactNum, WorkNum, RoleID, RoleName, ProfilePic, Address, Status, Pincode, Gender sql.NullString

					//var data Response2
					err := rows1.Scan(&Id, &First_name, &Last_name, &OfficeMailID, &ContactNum, &WorkNum, &RoleID, &RoleName, &ProfilePic, &Address, &Status, &Pincode, &Gender)
					if err != nil {
						fmt.Println("206 line", err)
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					data := make(map[string]interface{})
					data["id"] = Id.String
					data["first_name"] = First_name.String
					data["last_name"] = Last_name.String
					data["officeMailId"] = OfficeMailID.String
					data["contactNum"] = ContactNum.String
					data["workNum"] = WorkNum.String
					data["role_id"] = RoleID.String
					data["role_name"] = RoleName.String
					data["profile_pic"] = ProfilePic.String
					data["address"] = Address.String
					data["status"] = Status.String
					data["pincode"] = Pincode.String
					data["gender"] = Gender.String

					dataSlice = append(dataSlice, data)
				}
				w.Header().Set("Content-Type", "application/json")
				// Return JSON response indicating success
				json.NewEncoder(w).Encode(map[string]interface{}{"list": dataSlice, "total_count": totalRows, "code": http.StatusOK, "message": "Successfully", "success": true})
				return
			}
			if Role_ID == 2 || Role_ID == 1 {
				total_pages_sql = fmt.Sprintf("SELECT COUNT(*) FROM employee emp where status = 1 and emp.id !=%d %s %s", User_ID, role, searchFilter)

				fields = "emp.id, emp.first_name, IFNULL(emp.last_name,'') as last_name, IF(emp.empRole = 11, 1, emp.empRole)  as role_id, emp.officeMailId, emp.contactNum, emp.workNum, rm.roleName as role_name, IFNULL(emp.profile_pic,''), DATE_FORMAT(emp.doj, '%d-%m-%Y') as doj, emp.address, emp.address2, emp.address3,CONCAT(b.first_name,' ',b.last_name) AS supervisorName,COALESCE(emp.supervisorId,''), emp.status, emp.pincode, emp.license_number, emp.gender,emp.countryID,IFNULL(emp.personalMailId,''),IFNULL(emp.createdBy,''),IFNULL(emp.lastUpdatedAt,'')"

				query = fmt.Sprintf("SELECT %s FROM employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.id !=%d %s %s ORDER BY emp.first_name LIMIT %d, %d", fields, User_ID, role, searchFilter, offset, no_of_records)

			}
			getIds = fmt.Sprintf("SELECT IFNULL(GROUP_CONCAT(Level SEPARATOR ','), '')  FROM (SELECT @Ids := (SELECT GROUP_CONCAT(`ID` SEPARATOR ',')FROM employee WHERE FIND_IN_SET(`supervisorId`, @Ids) AND status = 1) Level FROM employee JOIN (SELECT @Ids:= %d) temp1) temp2", User_ID)
			rows, err := DB.Query(getIds)
			if err != nil {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			// var ids []string
			ids := make([]string, 0)
			var idd string
			for rows.Next() {

				err := rows.Scan(&idd)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				if idd == "" {
					idd = "''"
				}
				// ids = append(ids, id)

			}
			// Split the concatenated IDs into an array
			idss := strings.Split(idd, ",")
			// Remove empty strings from the array
			for _, a := range idss {
				if a != "" {
					ids = append(ids, a)
				}
			}
			// idsString := ids
			idsString := strings.Join(ids, ",")
			idsString = strings.TrimRight(idsString, ",")
			// fmt.Println("IDs:", idsString)

			if Role_ID == 3 {
				total_pages_sql = fmt.Sprintf("SELECT count(*) from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status =1  and emp.id IN(%s) %s %s", idsString, role, searchFilter)

				fields = "emp.id, emp.first_name, IFNULL(emp.last_name,'') as last_name, IF(emp.empRole = 11, 1, emp.empRole) as role_id, emp.officeMailId, emp.contactNum, emp.workNum, rm.roleName as role_name, IFNULL(emp.profile_pic,''), emp.doj, emp.address, emp.address2, emp.address3,CONCAT(b.first_name,' ',b.last_name) AS supervisorName, emp.supervisorId, emp.status, emp.pincode, emp.license_number, emp.gender,emp.countryID,IFNULL(emp.personalMailId,''),IFNULL(emp.createdBy,''),IFNULL(emp.lastUpdatedAt,'')"

				query = fmt.Sprintf("SELECT %s from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status =1  and emp.id IN(%s) %s %s ORDER BY emp.first_name LIMIT %d, %d", fields, idsString, role, searchFilter, offset, no_of_records)

			}

			if Role_ID == 12 {
				total_pages_sql = fmt.Sprintf("SELECT count(*) from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status =1  and emp.id IN(%s) %s %s", idsString, role, searchFilter)

				fields = "emp.id, emp.first_name, IFNULL(emp.last_name,'') as last_name, IF(emp.empRole = 11, 1, emp.empRole) as role_id, emp.officeMailId, emp.contactNum, emp.workNum, rm.roleName as role_name, IFNULL(emp.profile_pic,''), emp.doj, emp.address, emp.address2, emp.address3,CONCAT(b.first_name,' ',b.last_name) AS supervisorName, emp.supervisorId, emp.status, emp.pincode, emp.license_number, emp.gender,emp.countryID,IFNULL(emp.personalMailId,''),IFNULL(emp.createdBy,''),IFNULL(emp.lastUpdatedAt,'')"

				query = fmt.Sprintf("SELECT %s from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status =1  and emp.id IN(%s) %s %s ORDER BY emp.first_name LIMIT %d, %d", fields, idsString, role, searchFilter, offset, no_of_records)

			}
			if Role_ID == 4 {
				total_pages_sql = fmt.Sprintf("SELECT count(*) from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.id IN(%s) %s %s", idsString, role, searchFilter)

				fields = "emp.id, emp.first_name, IFNULL(emp.last_name,'') as last_name, emp.empRole as role_id, emp.officeMailId, emp.contactNum, emp.workNum, rm.roleName as role_name,IFNULL(emp.profile_pic,''), emp.doj, emp.address, emp.address2, emp.address3, CONCAT(b.first_name,' ',b.last_name) AS supervisorName, emp.supervisorId, emp.status, emp.pincode, emp.license_number, emp.gender,emp.countryID,IFNULL(emp.personalMailId,''),IFNULL(emp.createdBy,''),IFNULL(emp.lastUpdatedAt,'')"

				query = fmt.Sprintf("SELECT %s from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.id IN(%s) %s %s ORDER BY emp.first_name", fields, idsString, role, searchFilter)

			}
			if Role_ID == 5 {
				if request.Type == "senior" || request.Type == "" {

					total_pages_sql = fmt.Sprintf("SELECT count(*) from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.id IN(%s) %s", idsString, searchFilter)

					fields = "emp.id, emp.first_name, IFNULL(emp.last_name,'') as last_name, emp.empRole as role_id, emp.officeMailId, emp.contactNum, emp.workNum, rm.roleName as role_name,IFNULL(emp.profile_pic,''), emp.doj, emp.address, emp.address2, emp.address3,CONCAT(b.first_name,' ',b.last_name) AS supervisorName, emp.supervisorId, emp.status, emp.pincode, emp.license_number, emp.gender,emp.countryID,IFNULL(emp.personalMailId,''),IFNULL(emp.createdBy,''),IFNULL(emp.lastUpdatedAt,'')"

					query = fmt.Sprintf("SELECT %s from employee emp left join employee b on b.id = emp.supervisorId left join roles_Master rm on emp.empRole = rm.id where emp.status = 1 and emp.id IN(%s) %s ORDER BY emp.first_name LIMIT %d, %d", fields, idsString, role, offset, no_of_records)
				}
			}

			rows2, err := DB.Query(total_pages_sql)
			fmt.Println(total_pages_sql)
			if err != nil {
				fmt.Println(err)
				//w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return

			}

			for rows2.Next() {
				err := rows2.Scan(&totalRows)
				if err != nil {
					//w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
			}
			totalPages = int(math.Ceil(float64(totalRows) / float64(no_of_records)))
			fmt.Println(totalPages)

			rows3, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows3.Close()
			for rows3.Next() {
				var data1 Response2
				err := rows3.Scan(&data1.ID, &data1.FirstName, &data1.LastName, &data1.RoleID, &data1.OfficeMailID, &data1.ContactNum, &data1.WorkNum, &data1.RoleName, &data1.ProfilePic, &data1.Doj, &data1.Address, &data1.Address2, &data1.Address3, &data1.SupervisorName, &data1.SupervisorId, &data1.Status, &data1.Pincode, &data1.License_Number, &data1.Gender, &data1.CountryID, &data1.PersonalMailId, &data1.CreatedBy, &data1.LastUpdatedBy)
				fmt.Println("err", err)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}

				dataSlice2 = append(dataSlice2, data1)
			}
			var user_id, query_project string
			var target, id, projectName, emp_id, project_id, role_id sql.NullString
			for key, val := range dataSlice2 {
				projects := make([]map[string]interface{}, 0)
				//// If RoleID is 1 or 2, set an empty project list
				if val.RoleID == "1" || val.RoleID == "2" {
					dataSlice2[key].ProjectList = projects
				}
				if val.RoleID == "3" {
					// If RoleID is 3, retrieve projects based on the user's supervisor ID
					user_id = val.ID
					fields = "prj.id, UPPER(prj.projectName) as projectName"
					query_project = fmt.Sprintf("SELECT %s from employee emp join project prj on emp.id = prj.operations_manager where emp.empRole = 4 and emp.supervisorId = %s", fields, user_id)
					projectRows, err := DB.Query(query_project)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer projectRows.Close()

					for projectRows.Next() {
						err := projectRows.Scan(&id, &projectName)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}
						temp := make(map[string]interface{})
						temp["id"] = id.String
						temp["projectName"] = projectName.String
						projects = append(projects, temp)
					}

					dataSlice2[key].ProjectList = projects
					// If RoleID is 4 or 12, retrieve projects based on the user's operations manager ID
				} else if val.RoleID == "4" || val.RoleID == "12" {
					user_id = val.ID
					fields = "prj.id, UPPER(prj.projectName) as projectName"
					query_project = fmt.Sprintf("SELECT %s from project prj where prj.operations_manager = %s", fields, user_id)
					Rows, err := DB.Query(query_project)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer Rows.Close()

					for Rows.Next() {

						err := Rows.Scan(&id, &projectName)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}

						temp1 := make(map[string]interface{})
						temp1["id"] = id.String
						temp1["projectName"] = projectName.String
						projects = append(projects, temp1)
					}

					dataSlice2[key].ProjectList = projects

					// If RoleID is 5 or 6, retrieve projects based on the user's employee ID
				} else if val.RoleID == "5" || val.RoleID == "6" {
					user_id = val.ID
					fields = "em_pr.*, UPPER(pr.projectName) as projectName"
					query_project = fmt.Sprintf("SELECT %s from project_emps em_pr left join project pr on em_pr.project_id = pr.id where emp_id IN (%s)", fields, user_id)
					Rows, err := DB.Query(query_project)
					fmt.Println(query_project)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer Rows.Close()

					for Rows.Next() {
						err := Rows.Scan(&id, &emp_id, &project_id, &role_id, &target, &projectName)
						if err != nil {
							fmt.Println(err)
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}
						temp2 := make(map[string]interface{})
						temp2["id"] = id.String
						temp2["emp_id"] = emp_id.String
						temp2["project_id"] = project_id.String
						temp2["role_id"] = role_id.String
						temp2["target"] = target.String
						temp2["projectName"] = projectName.String
						projects = append(projects, temp2)

					}

					dataSlice2[key].ProjectList = projects

					//If RoleID is 13, retrieve projects based on the user's gfl_id
				} else if val.RoleID == "13" {
					user_id = val.ID
					fields = "UPPER(pr.projectName) as projectName"
					query_project = fmt.Sprintf("SELECT %s from  project pr where gfl_id IN (%s)", fields, user_id)
					Rows, err := DB.Query(query_project)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer Rows.Close()

					for Rows.Next() {
						err := Rows.Scan(&projectName)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}
						temp3 := make(map[string]interface{})
						temp3["projectName"] = projectName.String
						projects = append(projects, temp3)
					}

					dataSlice2[key].ProjectList = projects

					// If RoleID is 7, retrieve projects based on the user's driverID
				} else if val.RoleID == "7" {
					user_id = val.ID
					fields = "prj.id, UPPER(prj.projectName) as projectName"
					query_project = fmt.Sprintf("SELECT %s from project prj left join employee emp on prj.driverID = emp.id where prj.driverID = %s", fields, user_id)
					Rows, err := DB.Query(query_project)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
						return
					}
					defer Rows.Close()

					for Rows.Next() {
						err := Rows.Scan(&id, &projectName)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
							return
						}

						temp4 := make(map[string]interface{})
						temp4["id"] = id.String
						temp4["projectName"] = projectName.String
						projects = append(projects, temp4)
					}

					dataSlice2[key].ProjectList = projects

				}

			}
		}
		count := pageno * no_of_records
		if totalRows > count {
			w.Header().Set("Content-Type", "application/json")
			// Return JSON response indicating success
			json.NewEncoder(w).Encode(map[string]interface{}{"list": dataSlice2, "hasNextPage": true, "total_count": strconv.Itoa(totalRows), "code": http.StatusOK, "message": "successfully", "success": true})
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			// Return JSON response indicating success
			json.NewEncoder(w).Encode(map[string]interface{}{"list": dataSlice2, "hasNextPage": false, "total_count": strconv.Itoa(totalRows), "code": http.StatusOK, "message": "successfully", "success": true})
			return
		}

	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Enter required fields", "success": false, "Code": 200})
		return
	}

}

func getEmpPrj(DB *sql.DB, empID int, roleID int, filter int) string {
	//where := ""
	//var DB *sql.DB
	var fields, projects string
	if filter == 8 {
		fields = "GROUP_CONCAT(DISTINCT prj.funderID) as ids"
	} else {
		fields = "GROUP_CONCAT(DISTINCT prj.partnerID) as ids"
	}

	query := fmt.Sprintf("SELECT %s FROM employee emp JOIN project prj ON emp.id = prj.operations_manager WHERE emp.empRole = 4 AND emp.supervisorId = %d GROUP BY emp.id", fields, empID)
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("ERROR>>", err)

	}
	fmt.Println(query)
	defer rows.Close()
	var Ids []string
	for rows.Next() {
		err := rows.Scan(&projects)
		if err != nil {
			fmt.Println(err)
		}
		Ids = append(Ids, projects)
	}

	ids := strings.Join(Ids, ",")
	if ids == "" {
		return "''"
	}
	fmt.Println("ids", ids)
	return ids
}
