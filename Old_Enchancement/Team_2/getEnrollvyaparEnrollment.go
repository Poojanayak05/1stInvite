package Team_2

import (
	b "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func GetEnrollVyaparEnrollment(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	type Request struct {
		Search    string `json:"search"`
		ProjectID string `json:"project_id"`
		EmpID     string `json:"emp_id"`
		StartDate int    `json:"start_date"`
		EndDate   int    `json:"end_date"`
		GelathiId string `json:"gelathi_id"`
		PageNum   int    `json:"pageNum"`
	}

	type Participant struct {
		ID          string `json:"id"`
		GelathiName string `json:"gelathiname"`
		VillageName string `json:"villagename"`
		TbID        string `json:"tb_id"`
		Flag        string `json:"flag"`
		EnrollDate  string `json:"enroll_date"`
		EnrolledBy  string `json:"enrolled_by"`
		IsSurvey    bool   `json:"is_survey"`
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed to Read Bodyr", "success": false, "error": err.Error()})
		return
	}

	defer r.Body.Close()

	var req Request
	err = json.Unmarshal(data, &req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to unmarshal", "success": false, "error": err.Error()})
		return
	}

	// Search filter by participant/gelathi name
	searchFilter := ""
	if len(req.Search) > 0 {
		searchFilter = " and firstName like '%" + req.Search + "%'"
	}

	if req.ProjectID != "" && req.EmpID != "" {

		search := ""
		if req.StartDate != 0 && req.EndDate != 0 {
			search = fmt.Sprintf("AND date(tr_pat.VyaparEnrollmentDate) >= '%d' AND date(tr_pat.VyaparEnrollmentDate) <= '%d'", req.StartDate, req.EndDate)

		}
		if req.GelathiId != "" {

			search = fmt.Sprintf("AND tr_pat.gelathi_id = %s", req.GelathiId)
		}

		if req.PageNum == 0 {
			req.PageNum = 1
		}

		// const noOfRecords = 100
		// offset := (req.PageNum - 1) * noOfRecords
		ProjectintID, _ := strconv.Atoi(req.ProjectID)
		projArra, _ := b.GetAssociatedProjectList(DB, ProjectintID)

		// Convert the array of int to a string slice
		projArray := make([]string, len(projArra))
		for i, v := range projArra {
			projArray[i] = strconv.Itoa(v)
		}

		// Join the string slice elements into a single string
		// projArray := "[" + fmt.Sprintf("%s", strArray) + "]"
		projIDs := strings.Join(projArray, ",")

		// Remove the trailing comma from the projIDs
		projIDs = strings.TrimRight(projIDs, ",")

		fields := "empRole"

		queryRole := fmt.Sprintf("SELECT %s FROM employee WHERE status = 1 AND id = %s", fields, req.EmpID)
		row := DB.QueryRow(queryRole)

		var roleID int
		row.Scan(&roleID)

		totalPagesQuery := ""
		query := ""

		if roleID == 6 {
			fields := "gf_bat.training_batch_id as ids"

			// Prepare the list of project IDs as a comma-separated string

			query = fmt.Sprintf("SELECT %s FROM gf_batches gf_bat WHERE gf_bat.project_id IN(%s) AND gf_bat.emp_id = %s", fields, projIDs, req.EmpID)

			rows, err := DB.Query(query)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
				fmt.Print("line 104")
				return
			}
			defer rows.Close()

			var trainingBatchIDs string
			for rows.Next() {
				var batchID string
				err := rows.Scan(&batchID)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Scan Error for batchID", "success": false, "error": err.Error()})
					fmt.Print("line 114")

					return
				}
				trainingBatchIDs += batchID + ","
			}
			// Remove the trailing comma from the trainingBatchIDs
			trainingBatchIDs = strings.TrimRight(trainingBatchIDs, ",")
			if trainingBatchIDs == "" {
				trainingBatchIDs = "0"
			}

			totalPagesQuery = fmt.Sprintf("SELECT COUNT(*) FROM training_participants tr_pat WHERE enrolledProject IN (%s) AND VyaparEnrollment='1' AND tb_id IN (%s) %s", trainingBatchIDs, search, searchFilter)

			fields = "tr_pat.id, tr_pat.firstName as gelathiname, SUBSTRING_INDEX(IFNULL(tbl_poa.name, ''), '_', 1) as villagename, tr_pat.tb_id, if((select id from gelathi_circle gelath_cr WHERE gelath_cr.gelathi_id = tr_pat.id LIMIT 1), 1, 0) as flag, IFNULL(DATE_FORMAT(tr_pat.VyaparEnrollmentDate, '%d-%m-%Y'), '') as enroll_date, '' AS enrolled_by"

			extra := ""
			if req.EmpID == "322" {
				extra = "AND tr_pat.VyaparEnrollmentDate BETWEEN '2023-04-12' AND '2024-03-31'"
			}

			query = fmt.Sprintf("SELECT %s FROM training_participants tr_pat LEFT JOIN tbl_poa ON tr_pat.tb_id = tbl_poa.id WHERE tr_pat.project_id IN (%s) AND tr_pat.VyaparEnrollment='1' %s AND tr_pat.tb_id IN (%s) %s ORDER BY tr_pat.VyaparEnrollmentDate DESC", fields, projIDs, extra, trainingBatchIDs, searchFilter)

		} else {
			totalPagesQuery = fmt.Sprintf("SELECT COUNT(*) FROM training_participants tr_pat WHERE enrolledProject IN (%s) AND VyaparEnrollment = 1 %s %s", projIDs, search, searchFilter)

			fieldsEnroll := "tr_pat.id, tr_pat.firstName as gelathiname, SUBSTRING_INDEX(IFNULL(tbl_poa.name, ''), '_', 1) as villagename,tr_pat.tb_id, '' as flag, IFNULL(DATE_FORMAT(tr_pat.VyaparEnrollmentDate, '%d-%m-%Y'), '') as enroll_date, CONCAT(IFNULL(emp.first_name, ''),' ',IFNULL(emp.last_name, '')) AS enrolled_by"

			query = fmt.Sprintf("SELECT %s FROM training_participants tr_pat "+
				"LEFT JOIN tbl_poa ON tr_pat.tb_id = tbl_poa.id "+
				"LEFT JOIN employee emp ON tr_pat.gelathi_id = emp.id "+
				"WHERE tr_pat.enrolledProject IN (%s) "+
				"AND tr_pat.VyaparEnrollment='1' %s %s "+
				"ORDER BY tr_pat.VyaparEnrollmentDate DESC",
				fieldsEnroll,
				projIDs,
				search,
				searchFilter,
			)

			// rows, err := DB.Query(query)
			// if err != nil {
			// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
			// 	fmt.Print("line 158")

			// 	return
			// }
			// defer rows.Close()
		}

		var totalRows int

		DB.QueryRow(totalPagesQuery).Scan(&totalRows)

		// var total_pages int=(totalRows)/(noOfRecords)

		rows, err := DB.Query(query)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
			fmt.Print("line 184")

		}
		// defer rows.Close()

		var data []Participant
		var count int
		for rows.Next() {
			count++
			var participant Participant

			err := rows.Scan(&participant.ID,
				&participant.GelathiName,
				&participant.VillageName,
				&participant.TbID,
				&participant.Flag,
				&participant.EnrollDate,
				&participant.EnrolledBy,
			)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Scan Error for Participants", "success": false, "error": err.Error()})
				fmt.Print("line 196")

			}

			anaQuery := fmt.Sprintf("SELECT * FROM BuzzVyaparProgramBaseline WHERE partcipantId = %v", participant.ID)
			anaRows, err := DB.Query(anaQuery)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
				fmt.Print("line 202")
			}
			defer anaRows.Close()

			var isSurvey bool
			if anaRows.Next() {
				isSurvey = true
			} else {
				isSurvey = false
			}
			participant.IsSurvey = isSurvey

			data = append(data, participant)
		}

		response := make(map[string]interface{})

		response["list"] = data
		response["code"] = 200
		response["total_count"] = count
		response["success"] = true
		response["message"] = "successfully"

		jsonData, err := json.Marshal(response)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false, "error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}

}
