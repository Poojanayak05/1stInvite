package Team_2

import (
	a "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ListDataDemo struct {
	AdminCount   string `json:"adminCount"`
	ProjectCount string `json:"project_count"`
	OpmCount     string `json:"opmCount"`
	PartnerCount string `json:"partnerCount"`
	FunderCount  string `json:"funderCount"`
	DriverCount  string `json:"driverCount"`
	BusCount     string `json:"busCount"`
	GelathiCount string `json:"gelathiCount"`
	TrainerCount string `json:"trainerCount"`
	PmCount      string `json:"pmCount"`
	BhCount      string `json:"bh_count"`
	VvCount      string `json:"vv_count"`
	CmCount      string `json:"cm_count"`
	TbCount      string `json:"tb_count"`
	SomCount     string `json:"somCount"`
	GflCount     string `json:"gflCount"`
}

func GetDemoGraphy(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	decoder := json.NewDecoder(r.Body)

	type Data struct {
		ProjectID  string `json:"project_id"`
		TalukID    string `json:"taluk_id"`
		DistrictID string `json:"district_id"`
		FunderID   string `json:"funder_id"`
		EmpID      int    `json:"emp_id"`
		RoleID     int    `json:"role_id"`
	}
	var data Data

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	condition := ""

	if data.ProjectID != "" {
		condition = fmt.Sprintf(" AND prj.id = %s", data.ProjectID)
	}

	if data.FunderID != "" {
		condition = fmt.Sprintf(" AND prj.funderID = %s", data.FunderID)
	}

	if data.TalukID != "" {
		condition = fmt.Sprintf(" AND prj.locationID = %s", data.TalukID)
	}

	if data.EmpID != 0 && data.RoleID != 0 {
		opsIds := a.GetOpsManagers(DB, data.EmpID)
		stringOpsIds := make([]string, len(opsIds))

		for i, id := range opsIds {
			stringOpsIds[i] = strconv.Itoa(id)
		}
		if len(opsIds) > 0 {
			condition = fmt.Sprintf(" AND prj.operations_manager IN (%s)", strings.Join(stringOpsIds, ","))
		} else {
			condition = " AND prj.operations_manager = 0"
		}
	}

	leftJoin := ""
	if data.DistrictID != "" {
		condition = fmt.Sprintf(" AND loc.parentId = %s", data.DistrictID)

		leftJoin = "LEFT JOIN location district ON loc.parentId = district.id"
	}

	queryPrj := fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as id, IFNULL(count(prj.id), 0) as project_count FROM project prj LEFT JOIN location loc ON prj.locationID = loc.id %s WHERE prj.project_status = '1' %s", leftJoin, condition)

	var projectsINValue sql.NullString
	var projectCountValue sql.NullString

	err = DB.QueryRow(queryPrj).Scan(&projectsINValue, &projectCountValue)
	if err != nil {
		log.Println("line 320")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}
	var projectsIN, projectCount string

	if projectsINValue.Valid {
		projectsIN = projectsINValue.String
	} else {
		projectsIN = ""
	}
	if projectCountValue.Valid {
		projectCount = projectCountValue.String
	} else {
		projectCount = ""
	}

	if projectsIN == "" {
		projectsIN = "''"
	}

	queryAdmin := fmt.Sprintf("SELECT count(id) as adminCount, '%s' as project_count FROM employee emp WHERE emp.empRole = '2' AND emp.status = '1'", projectCount)

	var adminCount int
	err = DB.QueryRow(queryAdmin).Scan(&adminCount, &projectCount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}
	query := fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.operations_manager) as opm_id FROM project prj where prj.id IN (%s)", projectsIN)

	var opmIdVal sql.NullString
	err = DB.QueryRow(query).Scan(&opmIdVal)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}
	var opmId string
	if opmIdVal.Valid {
		opmId = opmIdVal.String
	} else {
		opmId = ""
	}

	if opmId == "" {
		opmId = "''"
	}

	//Pm count
	queryPM := "SELECT count(*) as pmCount FROM employee emp WHERE empRole = 3"

	var pmCount int
	err = DB.QueryRow(queryPM).Scan(&pmCount)
	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//trainer count
	queryTrainer := "SELECT count(*) as trainerCount FROM employee e WHERE empRole = 5 AND status = 1"

	var trainerCount int
	err = DB.QueryRow(queryTrainer).Scan(&trainerCount)
	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//driver count
	queryDriver := "SELECT COUNT(*) as driverCount FROM employee e WHERE empRole = 7 AND status = 1"

	var driverCount int
	err = DB.QueryRow(queryDriver).Scan(&driverCount)
	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//emp count
	queryEmp := fmt.Sprintf("SELECT (SELECT COUNT(DISTINCT prj_opm.operations_manager) as opmCount FROM project prj_opm WHERE prj_opm.id IN (%s)) as opmCount, (SELECT COUNT(DISTINCT prj_partner.partnerID) as partnerCount FROM project prj_partner WHERE prj_partner.id IN (%s)) as partnerCount, (SELECT COUNT(DISTINCT prj_funder.funderID) as funderCount FROM project prj_funder WHERE prj_funder.id IN (%s)) as funderCount, (SELECT COUNT(DISTINCT prj.busID) as busCount FROM project prj WHERE prj.id IN (%s)) as busCount, (SELECT COUNT(DISTINCT prj_gelathi.emp_id) as gelathiCount FROM project_emps prj_gelathi WHERE prj_gelathi.role_id = '6' AND prj_gelathi.project_id IN (%s)) as gelathiCount", projectsIN, projectsIN, projectsIN, projectsIN, projectsIN)

	type EmpCounts struct {
		OpmCount     int `json:"opmCount"`
		PartnerCount int `json:"partnerCount"`
		FunderCount  int `json:"funderCount"`
		BusCount     int `json:"busCount"`
		GelathiCount int `json:"gelathiCount"`
	}

	var empCounts EmpCounts
	err = DB.QueryRow(queryEmp).Scan(&empCounts.OpmCount, &empCounts.PartnerCount, &empCounts.FunderCount, &empCounts.BusCount, &empCounts.GelathiCount)
	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//somCount

	somResult := fmt.Sprintf("SELECT COUNT(DISTINCT emp_id) as somCount FROM project_emps pe WHERE pe.role_id = 12 AND pe.project_id IN (%s)", projectsIN)

	var somCount int
	err = DB.QueryRow(somResult).Scan(&somCount)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//gfl Count
	queryGFL := "SELECT COUNT(*) as gflCount FROM employee e WHERE empRole = 13 AND status = 1"

	var gflCount int
	err = DB.QueryRow(queryGFL).Scan(&gflCount)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}

	//vv cm count
	queryTBL := fmt.Sprintf(`SELECT 
    (SELECT COUNT(DISTINCT primary_id) FROM tbl_poa WHERE project_id IN (%s) AND session_type = '3' AND check_out != '') as bh_count,
    (SELECT COUNT(DISTINCT primary_id) FROM tbl_poa WHERE project_id IN (%s) AND session_type = '2' AND check_out != '') as vv_count,
    (SELECT COUNT(DISTINCT primary_id) FROM tbl_poa WHERE project_id IN (%s) AND session_type = '1' AND check_out != '') as cm_count,
    COUNT(DISTINCT tb_id) as tb_count
FROM tbl_poa
WHERE project_id IN (%s) AND id IN (SELECT MAX(id) FROM tbl_poa WHERE type = 1 AND added = 0 GROUP BY tb_id) AND check_out != ''`,
		projectsIN, projectsIN, projectsIN, projectsIN)

	var bhCount, vvCount, cmCount, tbCount int

	err = DB.QueryRow(queryTBL).Scan(&bhCount, &vvCount, &cmCount, &tbCount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": err.Error()})
		return
	}
	var respList1 ListDataDemo
	respList1.AdminCount = strconv.Itoa(adminCount)
	respList1.BusCount = strconv.Itoa(empCounts.BusCount)
	respList1.BhCount = strconv.Itoa(bhCount)
	respList1.PartnerCount = strconv.Itoa(empCounts.PartnerCount)
	respList1.FunderCount = strconv.Itoa(empCounts.FunderCount)
	respList1.DriverCount = strconv.Itoa(driverCount)
	respList1.GelathiCount = strconv.Itoa(empCounts.GelathiCount)
	respList1.TrainerCount = strconv.Itoa(trainerCount)
	respList1.PmCount = strconv.Itoa(pmCount)
	respList1.VvCount = strconv.Itoa(vvCount)
	respList1.CmCount = strconv.Itoa(cmCount)
	respList1.TbCount = strconv.Itoa(tbCount)
	respList1.SomCount = strconv.Itoa(somCount)
	respList1.GflCount = strconv.Itoa(gflCount)
	respList1.ProjectCount = projectCount

	var respList = map[string]interface{}{
		"adminCount":    strconv.Itoa(adminCount),
		"project_count": projectCount,
		"opmCount":      strconv.Itoa(empCounts.OpmCount),
		"partnerCount":  strconv.Itoa(empCounts.PartnerCount),
		"funderCount":   strconv.Itoa(empCounts.FunderCount),
		"driverCount":   strconv.Itoa(driverCount),
		"busCount":      strconv.Itoa(empCounts.BusCount),
		"gelathiCount":  strconv.Itoa(empCounts.GelathiCount),
		"trainerCount":  strconv.Itoa(trainerCount),
		"pmCount":       strconv.Itoa(pmCount),
		"bh_count":      strconv.Itoa(bhCount),
		"vv_count":      strconv.Itoa(vvCount),
		"cm_count":      strconv.Itoa(cmCount),
		"tb_count":      strconv.Itoa(tbCount),
		"somCount":      strconv.Itoa(somCount),
		"gflCount":      strconv.Itoa(gflCount),
	}
	_ = respList

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{"list": respList1, "message": "successfully", "code": 200, "success": true})

}
