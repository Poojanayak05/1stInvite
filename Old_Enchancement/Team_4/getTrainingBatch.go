package Team_4

//Done by Keerthana

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Response2 struct {
	List       []TrainingBatch `json:"list"`
	Code       int             `json:"code"`
	TotalCount string          `json:"total_count"`
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
}

type TrainingBatch struct {
	TrainingBatchID string `json:"training_batch_id"`
	BatchName       string `json:"batch_name"`
	Day1            string `json:"day1"`
	Day2            string `json:"day2"`
	Day1Status      string `json:"day1_status"`
	Day2Status      string `json:"day2_status"`
}

type response struct {
	SearchParam string `json:"search"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	FilterType  string `json:"filter_type"`
	TrainerID   string `json:"trainer_id"`
	// pageno := int(1)
	// if val, ok := request["pageNum"].(float64); ok {
	// 	pageno = int(val)
	// }
	ProjectID string `json:"project_id"`
	EmpID     string `json:"emp_id"`
}

func GetTrainingBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method Not found", "success": false})
		return
	}
	var request response
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "error": err})
		return
	}
	// var request map[string]interface{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "error": err})
		return
	}
	EmpID, _ := strconv.Atoi(request.EmpID)
	ProjectID, _ := strconv.Atoi(request.ProjectID)
	var searchFilter string
	if len(request.SearchParam) > 0 {
		searchFilter = fmt.Sprintf(" and name like '%%%s%%'", request.SearchParam)
	}

	var search string
	if len(request.StartDate) > 0 && len(request.EndDate) > 0 {
		search = fmt.Sprintf(" and date(tr_btch.date) >= '%s' AND date(tr_btch.date) <= '%s'", request.StartDate, request.EndDate)
	}

	if len(request.FilterType) > 0 {
		search += fmt.Sprintf(" and tr_btch.status = '%s'", request.FilterType)
	}

	if len(request.TrainerID) > 0 {
		search += fmt.Sprintf(" and tr_btch.user_id = '%s'", request.TrainerID)
	}

	noOfRecords := 100000
	// offset := (pageno - 1) * noOfRecords

	var condition string
	empType, err := getEmployeeType(db, int(EmpID), w)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query get employee type", "success": false, "error": err})
		return
	}
	if EmpID > 0 && empType == 5 {
		condition = fmt.Sprintf(" AND tr_btch.user_id = '%v'", int(EmpID))
	}

	projectIDStr := fmt.Sprint(ProjectID)

	totalRows, err := getTotalPagesAndRows(db, projectIDStr, condition, search, searchFilter, noOfRecords, w)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query for get row count", "success": false, "error": err})
		return
	}

	fields := "tr_btch.tb_id as training_batch_id, UPPER(tr_btch.name) as batch_name, " +
		"GROUP_CONCAT(DISTINCT CONCAT(DATE_FORMAT(tr_btch.date, '%Y-%m-%d %H:%i:%s'))) as dates, " +
		"GROUP_CONCAT(CONCAT(tr_btch.status)) as status"

	query := fmt.Sprintf("SELECT %s FROM tbl_poa tr_btch WHERE tr_btch.project_id = '%s' AND tr_btch.type = '1' %s %s %s GROUP BY tr_btch.tb_id ORDER BY tr_btch.tb_id",
		fields, projectIDStr, condition, search, searchFilter)

	rows, err := db.Query(query)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false, "error": err})
		return
	}
	defer rows.Close()

	list := make([]TrainingBatch, 0)
	var day1, day1Status string
	for rows.Next() {
		var (
			trainingBatchID int
			batchName       string
			dates           string
			status          string
		)
		err := rows.Scan(&trainingBatchID, &batchName, &dates, &status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false, "error": err})
			return
		}

		statusArr := strings.Split(status, ",")
		dateArr := strings.Split(dates, ",")

		if len(dateArr) > 0 {
			day1 = dateArr[0]
		}
		if len(statusArr) > 0 {
			day1Status = getStatus(statusArr[0])
		}

		list = append(list, TrainingBatch{
			TrainingBatchID: fmt.Sprint(trainingBatchID),
			BatchName:       batchName,
			Day1:            day1,
			Day2:            getDate(dateArr, 1),
			Day1Status:      day1Status,
			Day2Status:      getStatus(statusArr[1]),
		})
	}

	response := Response2{
		List:       list,
		Code:       http.StatusOK,
		TotalCount: fmt.Sprint(totalRows),
		Success:    true,
		Message:    "successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to encode response", "success": false, "error": err})
		return
	}

}
func getTotalPagesAndRows(db *sql.DB, projectID string, condition string, search string, searchFilter string, noOfRecords int, w http.ResponseWriter) (int, error) {
	totalRows := 0

	totalPagesQuery := fmt.Sprintf("SELECT COUNT(DISTINCT tr_btch.tb_id) FROM tbl_poa tr_btch WHERE tr_btch.project_id = '%s' AND tr_btch.type = '1' %s %s %s",
		projectID, condition, search, searchFilter)

	err := db.QueryRow(totalPagesQuery).Scan(&totalRows)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "error": err})
		return 0, err
	}
	return int(totalRows), nil
}

func getDate(dateArr []string, index int) string {
	if len(dateArr) > index {
		return dateArr[index]
	}
	return ""
}

func getStatus(status string) string {
	switch status {
	case "1":
		return "Reschedule"
	case "2":
		return "Cancelled"
	default:
		return ""
	}
}

func getEmployeeType(db *sql.DB, empId int, w http.ResponseWriter) (empType int, err error) {
	fields := "empRole"
	empTypeQuery := fmt.Sprintf("SELECT %s FROM employee where status = 1 and id = %v", fields, empId)
	err = db.QueryRow(empTypeQuery).Scan(&empType)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "error": err})
		return 0, err
	}
	return empType, nil

}
