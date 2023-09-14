package Team_5

import (
	a "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Response7 struct {
	Count      int         `json:"count"`
	List       interface{} `json:"list"`
	Code       int         `json:"code"`
	TotalCount string      `json:"total_count"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
}

type Request struct {
	Search      string  `json:"search"`
	ProjectID   string  `json:"project_id"`
	EmpID       string  `json:"emp_id"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	GelathiID   string  `json:"gelathi_id"`
	PageNum     int     `json:"pageNum"`
	NoOfRecords int     `json:"no_of_records"`
	Fields      string  `json:"fields"`
	QueryRole   string  `json:"query_role"`
	ResultPage  string  `json:"result_page"`
	Query       string  `json:"query"`
	Data        []Datum `json:"data"`
}

type Datum struct {
	ID          string `json:"id"`
	GelathiName string `json:"gelathiname"`
	VillageName string `json:"villagename"`
	TbID        string `json:"tb_id"`
	Flag        string `json:"flag"`
	EnrollDate  string `json:"enroll_date"`
	EnrolledBy  string `json:"enrolled_by"`
}

func GetEnrollGelathi(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false})
		return
	}
	request := Request{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to parse request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": errorMessage, "success": false, "error": err})
		return
	}

	searchFilter := ""
	if len(request.Search) > 0 {
		searchFilter = fmt.Sprintf(" and firstName like '%%%s%%'", request.Search)
	}

	if request.ProjectID != "0" && request.EmpID != "0" {
		empID := request.EmpID
		search := ""

		if len(request.StartDate) > 0 && len(request.EndDate) > 0 {
			search = fmt.Sprintf("AND date(tr_pat.enroll_date) >= '%s' AND date(tr_pat.enroll_date) <= '%s'", request.StartDate, request.EndDate)
		}
		fieldValue := reflect.ValueOf(request).FieldByName("GelathiID")
		if fieldValue.IsValid() && fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			gelathiID := fieldValue.String()
			search += fmt.Sprintf(" AND tr_pat.gelathi_id = '%s'", gelathiID)
		}

		pageno := request.PageNum
		if pageno == 0 {
			pageno = 1
		}
		//a.GetAssociatedProjectList calling function from Dashboard folder
		var proj int
		proj, err = strconv.Atoi(request.ProjectID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := Response7{Success: false, Message: "Invalid project ID"}
			json.NewEncoder(w).Encode(response)
			return
		}
		projArrayi, err := a.GetAssociatedProjectList(DB, proj)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false})
			return
		}
		projArray := []string{}
		for _, val := range projArrayi {
			projArray = append(projArray, strconv.Itoa(val))
		}
		fields := "empRole"
		queryRole := fmt.Sprintf("SELECT %s FROM employee WHERE status = 1 AND id = %s", fields, empID)
		resultRole, err := DB.Query(queryRole)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error", "success": false})
			return
		}
		defer resultRole.Close()
		// Slice to hold multiple data roles
		var dataRoles []Datum

		for resultRole.Next() {
			// Create a new instance for each row
			dataRole := Datum{}

			err := resultRole.Scan(
				&dataRole.ID,
			)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false})
				return
			}
			// Append the dataRole to the dataRoles slice
			dataRoles = append(dataRoles, dataRole)
		}

		// Access the IDs of all fetched rows
		var roleID int // Declare the roleID variable before the loop

		for _, dataRole := range dataRoles {
			roleID, err = strconv.Atoi(dataRole.ID)
			if err != nil {
				// Handle the error if the conversion fails
				log.Println(err)
				// Continue with the next iteration or return an error response
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false})
				continue

			}
		}
		var totalPagesSQL string
		var query string
		if roleID == 6 {
			fields = "GROUP_CONCAT(gf_bat.training_batch_id) as ids"
			fields = "gf_bat.training_batch_id"
			trainingBatchIds := fmt.Sprintf("SELECT %s FROM gf_batches gf_bat WHERE gf_bat.project_id in (%s) and gf_bat.emp_id = %s", fields, strings.Join(projArray, ","), empID)
			totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM training_participants tr_pat WHERE enrolledProject in (%s) AND enroll = '1' and new_green!=1 and new_vyapar!=1 AND tb_id IN (%s) %s %s", strings.Join(projArray, ","), trainingBatchIds, search, searchFilter)
			fields = "tr_pat.id, tr_pat.firstName as gelathiname, SUBSTRING_INDEX(IFNULL(tbl_poa.name, ''), '_', 1) as villagename, tr_pat.tb_id, if((select id from gelathi_circle gelath_cr WHERE gelath_cr.gelathi_id = tr_pat.id LIMIT 1), 1, 0) as flag, IFNULL(DATE_FORMAT(tr_pat.enroll_date, '%d-%m-%Y'), '') as enroll_date, '' AS enrolled_by"
			query = fmt.Sprintf("SELECT %s FROM training_participants tr_pat LEFT JOIN tbl_poa ON tr_pat.tb_id = tbl_poa.id WHERE tr_pat.project_id in (%s) AND tr_pat.enroll = '1'  AND tr_pat.tb_id IN (%s) %s ORDER BY tr_pat.enroll_date DESC", fields, strings.Join(projArray, ","), trainingBatchIds, searchFilter)
		} else {
			totalPagesSQL = fmt.Sprintf("SELECT COUNT(*) FROM training_participants tr_pat WHERE enrolledProject in (%s) and new_green!=1 and new_vyapar!=1 AND enroll = 1 %s %s", strings.Join(projArray, ","), search, searchFilter)
			fieldsEnroll := "tr_pat.id, tr_pat.firstName as gelathiname, COALESCE(SUBSTRING_INDEX(tbl_poa.name, '_', 1),'') as villagename, '' as tb_id, '' as flag, IFNULL(DATE_FORMAT(tr_pat.enroll_date, '%d-%m-%Y'), '') as enroll_date, CONCAT(IFNULL(emp.first_name, ''),' ',IFNULL(emp.last_name, '')) AS enrolled_by"
			query = fmt.Sprintf("SELECT %s FROM training_participants tr_pat LEFT JOIN tbl_poa ON tr_pat.tb_id = tbl_poa.id LEFT JOIN employee emp ON tr_pat.gelathi_id = emp.id WHERE tr_pat.enrolledProject in (%s)  AND tr_pat.enroll = '1' %s %s ORDER BY tr_pat.enroll_date DESC", fieldsEnroll, strings.Join(projArray, ","), search, searchFilter)
		}

		// Execute the query and fetch the data
		resultPage, err := DB.Query(totalPagesSQL)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false})
			return
		}
		defer resultPage.Close()

		totalRows := "0"

		for resultPage.Next() {
			err := resultPage.Scan(&totalRows)
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error ", "success": false})
				return
			}
		}

		// totalPages := int(math.Ceil(float64(totalRows) / float64(noOfRecords)))
		result, err := DB.Query(query)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Bad Request", "success": false})
			return
		}
		defer result.Close()
		// Create a slice to hold the retrieved Datum
		data := []Datum{}
		// Iterate over the rows and append data to the slice
		for result.Next() {
			d := Datum{}
			err := result.Scan(&d.ID, &d.GelathiName, &d.VillageName, &d.TbID, &d.Flag, &d.EnrollDate, &d.EnrolledBy)
			if err != nil {
				log.Println(err)

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Bad Request", "success": false})
				return
			}
			data = append(data, d)
		}
		count1, err := strconv.Atoi(totalRows)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		//Reponse structure
		response := Response7{
			List:       data,
			Code:       200,
			TotalCount: totalRows,
			Count:      count1,
			Success:    true,
			Message:    "successfully",
		}
		json.NewEncoder(w).Encode(response)
	}

}
