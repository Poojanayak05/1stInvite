package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DataItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TalukName    string `json:"taluk_name"`
	ProjectName  string `json:"project_name"`
	Time         string `json:"time"`
	Date         string `json:"date"`
	Type         string `json:"type"`
	EmpName      string `json:"emp_name"`
	RoleName     string `json:"roleName"`
	Flag         string `json:"flag"`
	UserID       string `json:"user_id"`
	PrimaryID    string `json:"primary_id"`
	ActualDate   string `json:"actual_date"`
	Status       string `json:"status"`
	Date2        string `json:"date2"`
	Time2        string `json:"time2"`
	Description  string `json:"description"`
	AllDay       string `json:"all_day"`
	LocationName string `json:"location_name"`
	CheckIn      string `json:"check_in"`
	CheckOut     string `json:"check_out"`
}

func GetPoa1(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data struct {
		Date    string `json:"date"`
		For     int    `json:"for"`
		Team    string `jason:"team"`
		Emp_id  string `jason:"emp_id"`
		PageNum string `jason:"pageNum"`
	}

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "message": err.Error()})
		return
	}
	var pageNo int

	if data.PageNum != "" {
		pageNo, _ = strconv.Atoi(data.PageNum)
	} else {
		pageNo = 1
	}

	var no_of_records int = 1000000
	offset := (pageNo - 1) * no_of_records
	_ = offset

	var dateVar string
	var where string
	var id int
	var empIds string
	var total_pages_sql, fields, query string

	if data.For >= 0 {

		if data.Date != "" {
			dateVar = data.Date
		} else {
			dateVar = time.Now().Format("2006-01-02")
		}

		if data.For == 1 { //week
			where = fmt.Sprintf("WHERE yearweek(DATE(tbl_poa.date), 1) = yearweek('%s', 1)", dateVar)

		} else if data.For == 2 { //month
			where = fmt.Sprintf("WHERE MONTH(tbl_poa.date) = MONTH('%s') AND YEAR(tbl_poa.date) = YEAR('%s')", dateVar, dateVar)

		} else { //today
			where = fmt.Sprintf("WHERE (DATE_FORMAT(tbl_poa.date,'%%Y-%%m-%%d') = '%s')", dateVar)
		}
	}

	if data.Emp_id != "" {

		id, _ = strconv.Atoi(data.Emp_id)

		if data.Team == "1" {
			empIds = getIds(DB, id)

		} else {
			empIds = strconv.Itoa(id)
		}

		total_pages_sql = fmt.Sprintf("SELECT count(*) FROM tbl_poa tbl_poa left join location loc on tbl_poa.location_id = loc.id left join project prj on tbl_poa.project_id = prj.id %s AND tbl_poa.user_id IN (%s)", where, empIds)

		fields = fmt.Sprintf("IF(tbl_poa.type = '1', tbl_poa.tb_id, tbl_poa.id) as id, UPPER(tbl_poa.name) as name, IFNULL(loc.name, '') as taluk_name, IFNULL(UPPER(prj.projectName), '') as project_name, TIME_FORMAT(tbl_poa.date, '%%h:%%i %%p') as time, DATE_FORMAT(tbl_poa.date, '%%W %%b %%e %%Y') as date, tbl_poa.type, CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS emp_name, tbl_poa.type, rm.roleName, if((select b.id from tbl_poa b WHERE b.type = tbl_poa.type and b.user_id = %d and b.date = tbl_poa.date LIMIT 1), IFNULL((select b.id from tbl_poa b WHERE b.type = tbl_poa.type and b.user_id = %d and b.date = tbl_poa.date and b.primary_id = tbl_poa.primary_id LIMIT 1), 0), 0) as flag, tbl_poa.user_id, tbl_poa.id as primary_id, DATE(tbl_poa.date) as actual_date, tbl_poa.status, IFNULL(DATE(tbl_poa.date2), '') as date2, IFNULL(TIME_FORMAT(tbl_poa.date2, '%%h:%%i %%p'), '') as time2, IFNULL(tbl_poa.description, '') as description, tbl_poa.all_day, IFNULL(district2.name, '') as location_name, IF(tbl_poa.check_in IS NULL, 0, 1) as check_in, IF(tbl_poa.check_out IS NULL, 0, 1) as check_out", id, id)

		query = fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa left join employee emp on emp.id = tbl_poa.user_id left join location loc on tbl_poa.location_id = loc.id left join roles_Master rm on rm.id = emp.empRole left join project prj on tbl_poa.project_id = prj.id left join location district on prj.locationID = district.id left join location district2 on district.parentId = district2.id %s AND tbl_poa.user_id IN (%s) ORDER BY tbl_poa.date", fields, where, empIds)

		var countNum string

		rows, err := DB.Query(total_pages_sql)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			rows.Scan(&countNum)
		}

		rows.Close()

		rows1, err := DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dataList [][]DataItem
		list := make(map[string][]DataItem)
		var count int
		for rows1.Next() {
			count++

			var (
				id           int
				name         sql.NullString
				talukName    string
				projectName  string
				time         sql.NullString
				date         sql.NullString
				typ          sql.NullInt64
				empName      sql.NullString
				roleName     sql.NullString
				flag         int
				userID       sql.NullInt64
				primaryID    int
				actualDate   sql.NullString
				status       int
				date2        sql.NullString
				time2        sql.NullString
				description  string
				allDay       int
				locationName sql.NullString
				checkIn      int
				checkOut     int
			)

			var a sql.NullInt64
			err := rows1.Scan(&id, &name, &talukName, &projectName, &time, &date, &typ, &empName, &a, &roleName, &flag, &userID, &primaryID, &actualDate, &status, &date2, &time2, &description, &allDay, &locationName, &checkIn, &checkOut)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			var nameVar, dateVar, empNameVar, timeVar, roleNameVar, actualDateVar string
			var typeVar, user_id int
			if name.Valid {
				nameVar = name.String
			}
			if date.Valid {
				dateVar = date.String
			}
			if empName.Valid {
				empNameVar = empName.String
			}
			if time.Valid {
				timeVar = time.String
			}
			if roleName.Valid {
				roleNameVar = roleName.String
			}
			if actualDate.Valid {
				actualDateVar = actualDate.String
			}
			if userID.Valid {
				user_id = int(userID.Int64)
			}
			if typ.Valid {
				typeVar = int(typ.Int64)
			} else {
				typeVar = 0
			}
			var date2Var, time2Var, locationNameVar string
			if date2.Valid {
				date2Var = date2.String
			}

			if time2.Valid {
				time2Var = time2.String
			}
			if locationName.Valid {
				locationNameVar = locationName.String
			}

			item := DataItem{
				ID:           strconv.Itoa(id),
				Name:         nameVar,
				TalukName:    talukName,
				ProjectName:  projectName,
				Time:         timeVar,
				Date:         dateVar,
				Type:         strconv.Itoa(typeVar),
				EmpName:      empNameVar,
				RoleName:     roleNameVar,
				Flag:         strconv.Itoa(flag),
				UserID:       strconv.Itoa(user_id),
				PrimaryID:    strconv.Itoa(primaryID),
				ActualDate:   actualDateVar,
				Status:       strconv.Itoa(status),
				Date2:        date2Var,
				Time2:        time2Var,
				Description:  description,
				AllDay:       strconv.Itoa(allDay),
				LocationName: locationNameVar,
				CheckIn:      strconv.Itoa(checkIn),
				CheckOut:     strconv.Itoa(checkOut),
			}

			dateKey := item.Date
			list[dateKey] = append(list[dateKey], item)
		}

		//=========================
		for items := range list {
			length := len(list[items])
			reversedSlice := make([]DataItem, length)
			for i, j := 0, length-1; i < length; i, j = i+1, j-1 {
				reversedSlice[i] = list[items][j]

			}
			dataList = append(dataList, reversedSlice)
		}
		///=================
		w.Header().Set("Content-Type", "application/json")

		if len(dataList) == 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"count":   countNum,
				"data":    []interface{}{},
				"code":    200,
				"success": true,
				"message": "Successfully",
			})

		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"count":   countNum,
				"data":    dataList,
				"code":    200,
				"success": true,
				"message": "Successfully",
			})
		}

		//}

	}

}
func getIds(db *sql.DB, empID int) string {
	var ids []int

	rows, err := db.Query("SELECT emp.id FROM employee emp WHERE emp.supervisorId = ? AND emp.status = 1", empID)
	if err != nil {
		log.Println("ERROR>>", err)

	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println("ERROR>>", err)

		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return "''"
	}

	roleQuery := "SELECT empRole FROM employee emp WHERE emp.id = ? AND emp.status = 1"
	var empRole int
	err = db.QueryRow(roleQuery, empID).Scan(&empRole)
	if err != nil {
		log.Println("ERROR>>", err)
	}

	if empRole == 4 {
		ids = append(ids, empID)
	}

	idStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	return idStr
}
