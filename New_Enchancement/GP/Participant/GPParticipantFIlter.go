package green

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type GProject struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Target           int    `json:"target"`
	Actual           int    `json:"actual"`
	GreenEnroll      int    `json:"greenenroll"`
	NoofGreenCohorts int    `json:"noofgreencoharts"`
	Villages         int    `json:"villages"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	SelectType       string `json:"select_type"`
	NoofGreenSurvey  int    `json:"noofgreensurvey"`
	NoGreenModule    int    `json:"noofgreenmodule"`
}

type GResponse struct {
	SummaryTarget         int        `json:"summary_target"`
	SummaryVillages       int        `json:"summary_villages"`
	SummaryActuals        int        `json:"summary_actual"`
	SummaryGreenEnroll    int        `json:"summary_greenenroll"`
	SummaryNoGreensurvey  int        `json:"summary_nogreensurvey"`
	SummaryNoGreenCoharts int        `json:"summary_nofgreencoharts"`
	SummaryGreenModule    int        `json:"summary_greenmodule"`
	Data                  []GProject `json:"data"`
	Code                  int        `json:"code"`
	Count                 int        `json:"count"`
	Success               bool       `json:"success"`
	Message               string     `json:"message"`
}

func GPCounts(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var request struct {
		RoleID    int    `json:"role_id"`
		EmpID     int    `json:"emp_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		GfId      string `json:"gfid,omitempty"`
		Project   string `json:"project"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	response := GResponse{
		Data: make([]GProject, 0),
	}
	gfid, _ := strconv.Atoi(request.GfId)

	if request.RoleID == 1 || request.RoleID == 9 || request.RoleID == 3 || request.RoleID == 4 || request.RoleID == 12 || request.RoleID == 11 {
		filter := ""

		survey := ""
		greencoharts := ""

		if request.RoleID == 3 {
			opsIds := getReportingOpsManagers(DB, request.EmpID)
			filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

		} else if request.RoleID == 12 {
			opsIds := getOpsManagers(DB, request.EmpID)
			filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

		} else if request.RoleID == 4 {
			filter = " AND p.operations_manager = " + strconv.Itoa(request.EmpID)

		}

		if gfid > 0 {

			filter += " AND tp.user_id = " + request.GfId
			survey = " AND t.gelathi_id = " + (request.GfId)
			greencoharts = " AND tbl_poa.user_id = " + (request.GfId)

		}

		projectList := "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
			"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.GreenMotivators=1 or t.new_green = 1) and tp.project_id != ''" + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
				"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.GreenMotivators=1 or t.new_green = 1) and ((tp.date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter
		}
		fmt.Println(projectList)
		rows, err := DB.Query(projectList)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
		defer rows.Close()

		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoGreenSurvey := 0
		summaryGreenEnroll := 0
		summaryGreenmodule := 0

		summaryNoGreenCoharts := 0

		for rows.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}

			obj := GProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}

			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryActuals += obj.Actual

			obj.GreenEnroll = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryGreenEnroll += obj.GreenEnroll

			obj.NoofGreenCohorts = getParticipantFilternoGreenCoharts(DB, request.StartDate, request.EndDate, projectArray, greencoharts, request.RoleID)
			summaryNoGreenCoharts += obj.NoofGreenCohorts

			obj.NoofGreenSurvey = getParticipantFilternoGreensurvey(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryNoGreenSurvey += obj.NoofGreenSurvey

			obj.NoGreenModule = getParticipantFilterGreenModule(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryGreenmodule += obj.NoGreenModule
			obj.Villages = getParticipantFilterGreenGfBatchesNew(DB, request.StartDate, request.EndDate, projectArray, "", gfid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}

		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoGreensurvey = summaryNoGreenSurvey
		response.SummaryNoGreenCoharts = summaryNoGreenCoharts
		response.SummaryGreenEnroll = summaryGreenEnroll
		response.SummaryGreenModule = summaryGreenmodule

		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
	} else if request.RoleID == 6 {

		projectList := fmt.Sprintf("SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp "+
			"INNER JOIN project p ON tp.project_id = p.id "+
			"JOIN GreenBaselineSurvey g ON g.partcipantId = tp.id "+
			"WHERE (tp.GreenMotivators=1 or tp.new_green = 1) AND ((tp.GreenMotivatorsDate BETWEEN '%s' AND '%s') or (g.entry_date BETWEEN '%s' AND '%s')) AND tp.gelathi_id = %d",
			request.StartDate, request.EndDate, request.StartDate, request.EndDate, request.EmpID)

		tRes, err := DB.Query(projectList)
		fmt.Println(projectList)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoGreenSurvey := 0
		summaryGreenEnroll := 0
		summaryGreenmodule := 0

		summaryNoGreenCoharts := 0
		for tRes.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := tRes.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}

			obj := GProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}

			projectArray := []int{obj.ID}

			obj.Target = getTrainerTarget(DB, request.EmpID, projectArray)
			summaryTarget += obj.Target
			greeneroll := " and tp.gelathi_id = " + strconv.Itoa(request.EmpID)
			filterFn := " and t.gelathi_id = " + strconv.Itoa(request.EmpID)
			greencoharts := " AND tbl_poa.user_id = " + strconv.Itoa(request.EmpID)
			obj.Actual = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, greeneroll, request.RoleID)
			summaryActuals += obj.Actual
			obj.NoofGreenCohorts = getParticipantFilternoGreenCoharts(DB, request.StartDate, request.EndDate, projectArray, greencoharts, request.RoleID)
			summaryNoGreenCoharts += obj.NoofGreenCohorts
			empid := strconv.Itoa(request.EmpID)
			obj.Villages = TrainerVillageCount(DB, request.StartDate, request.EndDate, projectArray, empid)
			summaryVillages += obj.Villages
			obj.NoGreenModule = getParticipantFilterGreenModule(DB, request.StartDate, request.EndDate, projectArray, filterFn, request.RoleID)
			summaryGreenmodule += obj.NoGreenModule
			obj.GreenEnroll = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, greeneroll, request.RoleID)
			summaryGreenEnroll += obj.GreenEnroll
			obj.NoofGreenSurvey = getParticipantFilternoGreensurvey(DB, request.StartDate, request.EndDate, projectArray, filterFn, request.RoleID)
			summaryNoGreenSurvey += obj.NoofGreenSurvey
			response.Data = append(response.Data, obj)
		}
		response.SummaryTarget = summaryTarget
		response.SummaryNoGreensurvey = summaryNoGreenSurvey
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoGreenCoharts = summaryNoGreenCoharts
		response.SummaryGreenEnroll = summaryGreenEnroll
		response.SummaryGreenModule = summaryGreenmodule
		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
	} else if request.RoleID == 13 {
		filter := ""
		filterFn := ""
		survey := ""
		greencoharts := ""

		if gfid > 0 {
			filter += " AND tp.gelathi_id = " + request.GfId
			filterFn = " AND tp.gelathi_id = " + (request.GfId)
			survey = " AND t.gelathi_id = " + (request.GfId)
			greencoharts = " AND tbl_poa.user_id = " + (request.GfId)

		}

		projectList := "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
			"INNER JOIN project p ON p.id = tp.project_id " +
			"JOIN GreenBaselineSurvey g ON g.partcipantId = tp.id " +
			"WHERE (tp.GreenMotivators=1 or tp.new_green = 1) AND p.gfl_id = ?" + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
				"INNER JOIN project p ON p.id = tp.project_id " +
				"JOIN GreenBaselineSurvey g ON g.partcipantId = tp.id " +
				"WHERE (tp.GreenMotivators=1 or tp.new_green = 1) AND ((tp.GreenMotivatorsDate BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "')) or ((g.entry_date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter + " AND p.gfl_id = ?"
		}

		pr, err := DB.Query(projectList, request.EmpID)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}

		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoGreenSurvey := 0
		summaryGreenEnroll := 0
		summaryGreenmodule := 0

		summaryNoGreenCoharts := 0

		for pr.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := pr.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}

			obj := GProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}

			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, filterFn, request.RoleID)
			summaryActuals += obj.Actual

			obj.GreenEnroll = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, filterFn, request.RoleID)
			summaryGreenEnroll += obj.GreenEnroll

			obj.NoofGreenCohorts = getParticipantFilternoGreenCoharts(DB, request.StartDate, request.EndDate, projectArray, greencoharts, request.RoleID)
			summaryNoGreenCoharts += obj.NoofGreenCohorts

			obj.NoofGreenSurvey = getParticipantFilternoGreensurvey(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryNoGreenSurvey += obj.NoofGreenSurvey

			obj.NoGreenModule = getParticipantFilterGreenModule(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryGreenmodule += obj.NoGreenModule
			obj.Villages = getParticipantFilterTrainingBatchesNew(DB, request.StartDate, request.EndDate, projectArray, "", gfid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}
		defer pr.Close()

		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoGreensurvey = summaryNoGreenSurvey
		response.SummaryNoGreenCoharts = summaryNoGreenCoharts
		response.SummaryGreenEnroll = summaryGreenEnroll
		response.SummaryGreenModule = summaryGreenmodule

		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("ERROR>>", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
			return
		}
	}

}
