package spoorthi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SProject struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Target                 int    `json:"target"`
	Actual                 int    `json:"actual"`
	SpoorthiEnroll         int    `json:"spoorthienroll"`
	NoofSpoorthiCircleMeet int    `json:"noofspoortthimeeting"`
	Villages               int    `json:"villages"`
	StartDate              string `json:"start_date"`
	EndDate                string `json:"end_date"`
	SelectType             string `json:"select_type"`
	NoofSpoorthiSurvey     int    `json:"noofspoorthisurvey"`
	NoSpoorthiModule       int    `json:"noofspoorthimodule"`
	NoofSpoorthiBeehives   int    `json:"noofspoorthibeehives"`
}

type SResponse struct {
	SummaryTarget               int        `json:"summary_target"`
	SummaryVillages             int        `json:"summary_villages"`
	SummaryActuals              int        `json:"summary_actual"`
	SummarySpoorthiEnroll       int        `json:"summary_spoorthienroll"`
	SummaryNoSpoorthisurvey     int        `json:"summary_nospoorthisurvey"`
	SummaryNoSpoorthiCircleMeet int        `json:"summary_nospoorthiciclemeet"`
	SummarySpoorthiModule       int        `json:"summary_spoorthimodule"`
	SummaryNoofSpoorthiBeehives int        `json:"summary_noofspoorthibeehives"`
	Data                        []SProject `json:"data"`
	Code                        int        `json:"code"`
	Count                       int        `json:"count"`
	Success                     bool       `json:"success"`
	Message                     string     `json:"message"`
}

func SPCounts(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
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

	response := SResponse{
		Data: make([]SProject, 0),
	}
	gfid, _ := strconv.Atoi(request.GfId)

	if request.RoleID == 1 || request.RoleID == 9 || request.RoleID == 3 || request.RoleID == 4 || request.RoleID == 12 || request.RoleID == 11 {
		filter := ""
		filterFn := ""
		survey := ""

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
			filter += " AND tp.gelathi_id = " + request.GfId
			filterFn = " AND tp.user_id = " + (request.GfId)
			survey = " AND t.gelathi_id = " + (request.GfId)

		}

		projectList := "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
			"JOIN SpoorthiBaselineQuestionnaire s ON s.partcipantId = tp.id " +
			"INNER JOIN project p ON tp.project_id = p.id " +
			"WHERE tp.project_id != '' AND tp.enroll=1 " + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
				"JOIN SpoorthiBaselineQuestionnaire s ON s.partcipantId = tp.id " +
				"INNER JOIN project p ON tp.project_id = p.id " +
				"WHERE tp.enroll=1 AND ((tp.enroll_date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "')) or ((s.entry_date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter
		}

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
		summaryNoSpoorthiSurvey := 0
		summarySpoorthiEnroll := 0
		summarySpoorthimodule := 0
		summarySpoorthibeehives := 0
		summaryNoSpoorthiCircleMeet := 0

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

			obj := SProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}
			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilterActual(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryActuals += obj.Actual

			obj.SpoorthiEnroll = getParticipantFilterGelathi(DB, request.StartDate, request.EndDate, projectArray, filter)
			summarySpoorthiEnroll += obj.SpoorthiEnroll

			obj.NoofSpoorthiCircleMeet = getParticipantFilterSpoortthiCircleMeet(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryNoSpoorthiCircleMeet += obj.NoofSpoorthiCircleMeet

			obj.NoofSpoorthiSurvey = getParticipantFilterSpoorthisurvey(DB, request.StartDate, request.EndDate, projectArray, survey)
			summaryNoSpoorthiSurvey += obj.NoofSpoorthiSurvey

			obj.NoSpoorthiModule = getParticipantFilterSpoorthiModule(DB, request.StartDate, request.EndDate, projectArray, survey)
			summarySpoorthimodule += obj.NoSpoorthiModule

			obj.NoofSpoorthiBeehives = getParticipantFilterSpoortthiBeehives(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summarySpoorthibeehives += obj.NoofSpoorthiBeehives
			obj.Villages = getParticipantFilterTrainingBatchesNew(DB, request.StartDate, request.EndDate, projectArray, "", gfid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}

		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoSpoorthisurvey = summaryNoSpoorthiSurvey
		response.SummaryNoSpoorthiCircleMeet = summaryNoSpoorthiCircleMeet
		response.SummarySpoorthiEnroll = summarySpoorthiEnroll
		response.SummarySpoorthiModule = summarySpoorthimodule
		response.SummaryNoofSpoorthiBeehives = summarySpoorthibeehives
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
			"JOIN SpoorthiBaselineQuestionnaire s ON s.partcipantId = tp.id "+
			"WHERE tp.enroll=1 AND ((tp.enroll_date BETWEEN '%s' AND '%s') or (s.entry_date BETWEEN '%s' AND '%s')) AND tp.gelathi_id = %d",
			request.StartDate, request.EndDate, request.StartDate, request.EndDate, request.EmpID)

		tRes, err := DB.Query(projectList)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoSpoorthiSurvey := 0
		summarySpoorthiEnroll := 0
		summarySpoorthimodule := 0
		summarySpoorthibeehives := 0
		summaryNoSpoorthiCircleMeet := 0
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

			obj := SProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}

			projectArray := []int{obj.ID}

			obj.Target = getTrainerTarget(DB, request.EmpID, projectArray)
			summaryTarget += obj.Target
			filter := " and tp.trainer_id = " + strconv.Itoa(request.EmpID)

			obj.Actual = getParticipantFilterActual(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryActuals += obj.Actual
			greeneroll := " and tp.gelathi_id = " + strconv.Itoa(request.EmpID)
			filterFn := " and t.gelathi_id = " + strconv.Itoa(request.EmpID)
			filterC := " and tp.user_id = " + strconv.Itoa(request.EmpID)
			obj.NoofSpoorthiCircleMeet = getParticipantFilterSpoortthiCircleMeet(DB, request.StartDate, request.EndDate, projectArray, filterC)
			summaryNoSpoorthiCircleMeet += obj.NoofSpoorthiCircleMeet
			empid := strconv.Itoa(request.EmpID)
			obj.Villages = TrainerVillageCount(DB, request.StartDate, request.EndDate, projectArray, empid)
			summaryVillages += obj.Villages
			obj.NoSpoorthiModule = getParticipantFilterSpoorthiModule(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summarySpoorthimodule += obj.NoSpoorthiModule
			obj.SpoorthiEnroll = getParticipantFilterGelathi(DB, request.StartDate, request.EndDate, projectArray, greeneroll)
			summarySpoorthiEnroll += obj.SpoorthiEnroll
			obj.NoofSpoorthiBeehives = getParticipantFilterSpoortthiBeehives(DB, request.StartDate, request.EndDate, projectArray, filterC)
			summarySpoorthibeehives += obj.NoofSpoorthiBeehives
			obj.NoofSpoorthiSurvey = getParticipantFilterSpoorthisurvey(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryNoSpoorthiSurvey += obj.NoofSpoorthiSurvey
			response.Data = append(response.Data, obj)
		}
		response.SummaryTarget = summaryTarget
		response.SummaryNoSpoorthisurvey = summaryNoSpoorthiSurvey
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoSpoorthiCircleMeet = summaryNoSpoorthiCircleMeet
		response.SummarySpoorthiEnroll = summarySpoorthiEnroll
		response.SummarySpoorthiModule = summarySpoorthimodule
		response.SummaryNoofSpoorthiBeehives = summarySpoorthibeehives
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
		if gfid > 0 {
			filter += " AND tp.gelathi_id = " + request.GfId
			filterFn = " AND tp.user_id = " + (request.GfId)
			survey = " AND t.gelathi_id = " + (request.GfId)
		}

		projectList := "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
			"INNER JOIN project p ON p.id = tp.project_id " +
			"JOIN SpoorthiBaselineQuestionnaire s ON s.partcipantId = tp.id " +
			"WHERE tp.enroll = 1 AND p.gfl_id = ?" + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
				"INNER JOIN project p ON p.id = tp.project_id " +
				"JOIN SpoorthiBaselineQuestionnaire s ON s.partcipantId = tp.id " +
				"WHERE tp.enroll=1 AND ((tp.enroll_date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "')) or ((s.entry_date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter + " AND p.gfl_id = ?"
		}

		rows, err := DB.Query(projectList, request.EmpID)
		if err != nil {
			log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
		}

		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoSpoorthiSurvey := 0
		summarySpoorthiEnroll := 0
		summarySpoorthimodule := 0
		summarySpoorthibeehives := 0
		summaryNoSpoorthiCircleMeet := 0

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

			obj := SProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}
			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilterActual(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryActuals += obj.Actual

			obj.SpoorthiEnroll = getParticipantFilterGelathi(DB, request.StartDate, request.EndDate, projectArray, filter)
			summarySpoorthiEnroll += obj.SpoorthiEnroll

			obj.NoofSpoorthiCircleMeet = getParticipantFilterSpoortthiCircleMeet(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryNoSpoorthiCircleMeet += obj.NoofSpoorthiCircleMeet

			obj.NoofSpoorthiSurvey = getParticipantFilterSpoorthisurvey(DB, request.StartDate, request.EndDate, projectArray, survey)
			summaryNoSpoorthiSurvey += obj.NoofSpoorthiSurvey

			obj.NoSpoorthiModule = getParticipantFilterSpoorthiModule(DB, request.StartDate, request.EndDate, projectArray, survey)
			summarySpoorthimodule += obj.NoSpoorthiModule

			obj.NoofSpoorthiBeehives = getParticipantFilterSpoortthiBeehives(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summarySpoorthibeehives += obj.NoofSpoorthiBeehives
			obj.Villages = getParticipantFilterTrainingBatchesNew(DB, request.StartDate, request.EndDate, projectArray, "", gfid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}
		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoSpoorthisurvey = summaryNoSpoorthiSurvey
		response.SummaryNoSpoorthiCircleMeet = summaryNoSpoorthiCircleMeet
		response.SummarySpoorthiEnroll = summarySpoorthiEnroll
		response.SummarySpoorthiModule = summarySpoorthimodule
		response.SummaryNoofSpoorthiBeehives = summarySpoorthibeehives
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
