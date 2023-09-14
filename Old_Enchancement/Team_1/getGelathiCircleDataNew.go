package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Circle struct {
	CircleID    int       `json:"circle_id"`
	ProjectID   int       `json:"project_id"`
	EmpID       int       `json:"emp_id"`
	CircleName  string    `json:"circle_name"`
	CircleDate  string    `json:"circle_date"`
	ProjectName string    `json:"project_name"`
	Gelathis    []Gelathi `json:"gelathis"`
}

type Gelathi struct {
	GelathiID      int    `json:"gelathi_id"`
	FirstName      string `json:"firstName"`
	VillageName    string `json:"villagename"`
	Flag           int    `json:"flag"`
	TBID           int    `json:"tb_id"`
	IsSurvey       bool   `json:"is_survey"`
	IsGreenSurvey  bool   `json:"is_green_survey"`
	IsVyaparSurvey bool   `json:"is_vyapar_survey"`
}

func GetGelathiCircleDataNew(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var request map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	circleID, ok := request["circle_id"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Circle ID", "success": false})
		return
	}

	projectID, ok := request["project_id"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid project_id",
			Success: false,
		})
		return
	}

	empID, ok := request["emp_id"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid emp_id",
			Success: false,
		})
		return
	}
	rows, err := DB.Query("SELECT id FROM bdms_staff.circle where id= ?", circleID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer rows.Close()
	if rows.Next() {

		circleFields := "cr.id as circle_id, cr.circle_name, DATE_FORMAT(cr.circle_date, '%d-%m-%Y') as circle_date, gelt_cr.gelathi_id, tr_part.firstName, UPPER(prj.projectName) as project_name, SUBSTRING_INDEX(tbl_poa.name, '_', 1) as villagename, if((select id from gelathi_circle gelath_cr WHERE gelath_cr.gelathi_id = gelt_cr.gelathi_id LIMIT 1), 1, 0) as flag, tr_part.tb_id"
		circleQuery := fmt.Sprintf("SELECT %s FROM circle cr LEFT JOIN gelathi_circle gelt_cr ON cr.id = gelt_cr.circle_id LEFT JOIN training_participants tr_part ON gelt_cr.gelathi_id = tr_part.id LEFT JOIN project prj ON cr.project_id = prj.id LEFT JOIN tbl_poa ON tr_part.tb_id = tbl_poa.id WHERE cr.id = %s", circleFields, circleID)

		circleRows, err := DB.Query(circleQuery)
		if err != nil {
			log.Println("Failed to execute query:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Failed to fetch circle details",
				"success": false,
			})
			return
		}
		defer circleRows.Close()

		var circle Circle
		var gelathis []Gelathi
		// var response1 map[string]interface{}
		// var pooja map[string]interface{}

		for circleRows.Next() {
			var gelathiID sql.NullInt64
			var firstName, villageName, projectName sql.NullString
			var flag, tbID sql.NullInt64
			// var circle_id string

			err := circleRows.Scan(&circle.CircleID, &circle.CircleName, &circle.CircleDate, &gelathiID, &firstName, &projectName, &villageName, &flag, &tbID)
			if err != nil {
				log.Println("Failed to execute query:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Failed to fetch circle details",
					"success": false,
				})
				return
			}
			// response1["circle_id"] = pooja["circle_id"]

			circle.ProjectName = projectName.String

			if gelathiID.Valid {
				gelathi := Gelathi{
					GelathiID:   int(gelathiID.Int64),
					FirstName:   firstName.String,
					VillageName: villageName.String,
					Flag:        int(flag.Int64),
					TBID:        int(tbID.Int64),
				}
				gelathis = append(gelathis, gelathi)
			}
		}

		circle.Gelathis = gelathis

		// Fetch enrolled gelathis
		batchIDsQuery := fmt.Sprintf("SELECT COALESCE(GROUP_CONCAT(gf_bat.training_batch_id),'') as ids FROM gf_batches gf_bat WHERE gf_bat.project_id = %s and gf_bat.emp_id = %s", projectID, empID)

		var trainingBatchIds string
		err = DB.QueryRow(batchIDsQuery).Scan(&trainingBatchIds)
		if err != nil {
			// Handle the error
			return
		}

		if trainingBatchIds == "" {
			trainingBatchIds = "''"
		}

		// Query to get gelathi data
		fields := "tr_pat.id as gelathi_id, tr_pat.firstName, SUBSTRING_INDEX(tbl_poa.name, '_', 1) as villagename, tr_pat.tb_id, IF((SELECT id FROM gelathi_circle gelath_cr WHERE gelath_cr.gelathi_id = tr_pat.id LIMIT 1), 1, 0) as flag"

		enrolledQuery := fmt.Sprintf(`
    SELECT %s 
    FROM training_participants tr_pat 
    LEFT JOIN tbl_poa ON tr_pat.tb_id = tbl_poa.id 
    LEFT JOIN gelathi_circle gl_cr ON tr_pat.id = gl_cr.gelathi_id AND gl_cr.gelathi_id != tr_pat.id 
    WHERE tr_pat.project_id = %s AND tr_pat.enroll = '1' AND tr_pat.tb_id IN (%s) AND tr_pat.id NOT IN (SELECT gelathi_id FROM gelathi_circle) 
    GROUP BY tr_pat.id 
    ORDER BY tr_pat.firstName
`, fields, projectID, trainingBatchIds)
		enrolledRows, err := DB.Query(enrolledQuery)
		if err != nil {
			log.Println("Failed to execute query:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch enrolled gelathis",
				Success: false,
			})
			return
		}
		defer enrolledRows.Close()

		for enrolledRows.Next() {
			var gelathiID sql.NullInt64
			var firstName, villageName sql.NullString
			var tbID, flag sql.NullInt64

			err := enrolledRows.Scan(&gelathiID, &firstName, &villageName, &tbID, &flag)
			if err != nil {
				log.Println("Failed to scan row:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Response{
					Code:    http.StatusInternalServerError,
					Message: "Failed to fetch enrolled gelathis",
					Success: false,
				})
				return
			}

			gelathi := Gelathi{
				GelathiID:   int(gelathiID.Int64),
				FirstName:   firstName.String,
				VillageName: villageName.String,
				Flag:        int(flag.Int64),
				TBID:        int(tbID.Int64),
			}
			gelathis = append(gelathis, gelathi)
		}

		circle.Gelathis = append(circle.Gelathis, gelathis...)

		// Add survey checks
		for i := range circle.Gelathis {
			// SpoorthiBaselineQuestionnaire check
			var count int
			err = DB.QueryRow("SELECT COUNT(*) FROM SpoorthiBaselineQuestionnaire WHERE partcipantId = ?", circle.Gelathis[i].GelathiID).Scan(&count)
			if err != nil {
				log.Println("Failed to execute query:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Failed to fetch circle details",
					"success": false,
				})
				return
			}
			if count > 0 {
				circle.Gelathis[i].IsSurvey = true
			} else {
				circle.Gelathis[i].IsSurvey = false
			}

			// GreenBaselineQuestionnaire check
			err = DB.QueryRow("SELECT COUNT(*) FROM GreenBaselineSurvey WHERE partcipantId = ?", circle.Gelathis[i].GelathiID).Scan(&count)
			if err != nil {
				log.Println("Failed to execute query:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Failed to fetch circle details",
					"success": false,
				})
				return
			}
			if count > 0 {
				circle.Gelathis[i].IsGreenSurvey = true
			} else {
				circle.Gelathis[i].IsGreenSurvey = false
			}

			// VyaparBaselineQuestionnaire check
			err = DB.QueryRow("SELECT COUNT(*) FROM BuzzVyaparProgramBaseline WHERE partcipantId = ?", circle.Gelathis[i].GelathiID).Scan(&count)
			if err != nil {
				log.Println("Failed to execute query:", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "Failed to fetch circle details",
					"success": false,
				})
				return
			}
			if count > 0 {
				circle.Gelathis[i].IsVyaparSurvey = true
			} else {
				circle.Gelathis[i].IsVyaparSurvey = false
			}
		}

		response := map[string]interface{}{
			"code":         200,
			"message":      "Successfully",
			"success":      true,
			"gelathis":     gelathis,
			"circle_id":    strconv.Itoa(circle.CircleID),
			"project_name": circle.ProjectName,
			"circle_name":  circle.CircleName,
			"circle_date":  circle.CircleDate,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	} else {
		response := map[string]interface{}{
			"code":    404,
			"message": "Invalid circle id",
			"success": true,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}
