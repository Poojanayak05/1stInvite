package Team_2

import (
	a "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func GetGFSessionsNew(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type RequestBody struct {
		Filter    string `json:"filter"`
		EndDate   string `json:"end_date"`
		Search    string `json:"search"`
		ProjectID string `json:"project_id"`
		GelathiID string `json:"gelathi_id"`
		StartDate string `json:"start_date"`
		EmpID     string `json:"emp_id"`
		PageNum   string `json:"pageNum"`
	}

	decoder := json.NewDecoder(r.Body)

	var data RequestBody

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}
	var PageNum, EmpID int
	EmpID, _ = strconv.Atoi(data.EmpID)
	if data.ProjectID != "" {

		if data.PageNum == "" {
			PageNum = 1
		}

		no_of_records := 100000
		offset := (PageNum - 1) * no_of_records

		searchFilter := ""
		searchParam := data.Search
		where := ""
		_ = offset

		if len(searchParam) > 0 {
			searchFilter = fmt.Sprintf(" and name like '%%%s%%'", searchParam)
		}

		if data.Filter != "" {
			session_type, err := strconv.Atoi(data.Filter)
			if err != nil {

				session_type = 0
			}

			if session_type == 22 {
				where = " AND status = '1'"
			} else if session_type == 23 {
				where = " AND status = '2'"
			} else {
				where = fmt.Sprintf(" AND session_type = %d", session_type)
			}
		}

		if data.StartDate != "" && data.EndDate != "" {
			where += fmt.Sprintf(" AND date(date) >= '%s' AND date(date) <= '%s'", data.StartDate, data.EndDate)
		}

		if data.GelathiID != "" {
			where += fmt.Sprintf(" AND user_id = '%s'", data.GelathiID)

		}
		fields := "empRole"
		query := fmt.Sprintf("SELECT %s FROM employee WHERE status = 1 AND id = %d", fields, EmpID)

		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
			return
		}
		var role_id int
		for rows.Next() {
			err := rows.Scan(&role_id)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
				return
			}

		}
		rows.Close()

		condition := ""
		if role_id == 6 {
			condition = fmt.Sprintf(" AND user_id = %d", EmpID)
		}
		projID, err := strconv.Atoi(data.ProjectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		projArray, _ := a.GetAssociatedProjectList(DB, projID)
		stringOpsIds := make([]string, len(projArray))

		for i, id := range projArray {
			stringOpsIds[i] = strconv.Itoa(id)
		}
		result := strings.Join(stringOpsIds, ",")

		totalPagesSQL := "SELECT COUNT(DISTINCT primary_id) FROM tbl_poa WHERE project_id IN (" + result + ") AND type = '2' " + condition + " " + where + " " + searchFilter

		rows, err = DB.Query(totalPagesSQL)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
			return
		}
		var total_rows int
		for rows.Next() {
			rows.Scan(&total_rows)
		}
		total_pages := math.Ceil(float64(total_rows) / float64(no_of_records))
		_ = total_pages
		fields = "DISTINCT primary_id as gf_session_id, name as gf_session_name, tb_id, DATE_FORMAT(date, '%d-%m-%Y %h:%i %p') as plan_date, status"
		query = fmt.Sprintf("SELECT %s, date FROM tbl_poa WHERE project_id IN (%s) AND type = '2' %s %s %s ORDER BY date DESC", fields, result, condition, where, searchFilter)

		type data struct {
			GFSessionID   string `json:"gf_session_id"`
			GFSessionName string `json:"gf_session_name"`
			TBID          string `json:"tb_id"`
			PlanDate      string `json:"plan_date"`
			Status        string `json:"status"`
		}
		finalRows, err1 := DB.Query(query)
		if err1 != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
			return
		}
		var list []data
		for finalRows.Next() {
			var obj data
			var temp string
			var gf_session_id, tb_id int
			err := finalRows.Scan(&gf_session_id, &obj.GFSessionName, &tb_id, &obj.PlanDate, &obj.Status, &temp)
			obj.GFSessionID = strconv.Itoa(gf_session_id)
			obj.TBID = strconv.Itoa(tb_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
			list = append(list, obj)

		}
		type response struct {
			List       []data `json:"list"`
			Code       int    `json:"code"`
			TotalCount string `json:"total_count"`
			Success    bool   `json:"success"`
			Message    string `json:"message"`
		}
		var res response
		res.Code = 200
		res.Success = true
		res.Message = "successfully"
		res.List = list
		res.TotalCount = strconv.Itoa(total_rows)
		if len(res.List) == 0 {
			res.List = []data{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return

	} else {
		var list []interface{}
		json.NewEncoder(w).Encode(map[string]interface{}{"list": list, "code": http.StatusOK, "message": "Successfully", "success": true, "total_count": nil})

	}

}
