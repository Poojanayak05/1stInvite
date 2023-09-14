package vyapar

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type VProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Target            int    `json:"target"`
	Actual            int    `json:"actual"`
	VyaparEnroll      int    `json:"vyaparenroll"`
	NoofVyaparCohorts int    `json:"noofvyaparcoharts"`
	Villages          int    `json:"villages"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	SelectType        string `json:"select_type"`
	NoofVyaparSurvey  int    `json:"noofvyaparsurvey"`
	NoVyaparModule    int    `json:"noofvyaparmodule"`
}

type VResponse struct {
	SummaryTarget          int        `json:"summary_target"`
	SummaryVillages        int        `json:"summary_villages"`
	SummaryActuals         int        `json:"summary_actual"`
	SummaryVyaparEnroll    int        `json:"summary_vyaparenroll"`
	SummaryNoVyaparsurvey  int        `json:"summary_novyaparsurvey"`
	SummaryNoVyaparCoharts int        `json:"summary_nofvyaparcoharts"`
	SummaryVyaparModule    int        `json:"summary_vyaparmodule"`
	Data                   []VProject `json:"data"`
	Code                   int        `json:"code"`
	Count                  int        `json:"count"`
	Success                bool       `json:"success"`
	Message                string     `json:"message"`
}

func GetCounts(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
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
		log.Fatal(err)
	}

	response := VResponse{
		Data: make([]VProject, 0),
	}
	gfid, _ := strconv.Atoi(request.GfId)

	if request.RoleID == 1 || request.RoleID == 9 || request.RoleID == 3 || request.RoleID == 4 || request.RoleID == 12 || request.RoleID == 11 {
		filter := ""
		filterFn := ""
		survey := ""

		vyaparcoharts := ""

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
			filterFn = " AND t.gelathi_id = " + (request.GfId)
			survey = " AND b.gfId = " + (request.GfId)
			vyaparcoharts = " AND tbl_poa.user_id = " + (request.GfId)

		}

		projectList := "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
			"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.VyaparEnrollment=1 or t.new_vyapar=1)  and tp.project_id != ''" + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
				"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.session_type in (16,17,18,19,20,21) and ((tp.date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter
		}
		fmt.Println(projectList)

		rows, err := DB.Query(projectList)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// summaryProjectsArray := make([]int, 0)
		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoVyaparSurvey := 0
		summaryVyaparEnroll := 0
		summaryVyaparmodule := 0

		summaryNoVyaparCoharts := 0

		for rows.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Fatal(err)
			}

			obj := VProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}
			// summaryProjectsArray = append(summaryProjectsArray, obj.ID)

			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilternoVyaparenrollCached(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryActuals += obj.Actual

			obj.VyaparEnroll = getParticipantFilternoVyaparenrollCached(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryVyaparEnroll += obj.VyaparEnroll

			obj.NoofVyaparCohorts = getParticipantFilternoVyaparCohartsCached(DB, request.StartDate, request.EndDate, projectArray, vyaparcoharts, request.RoleID)
			summaryNoVyaparCoharts += obj.NoofVyaparCohorts

			obj.NoofVyaparSurvey = getParticipantFilternoVyaparsurveyCached(DB, request.StartDate, request.EndDate, projectArray, survey, request.RoleID)
			summaryNoVyaparSurvey += obj.NoofVyaparSurvey

			obj.NoVyaparModule = getParticipantFilterVyaprModuleCached(DB, request.StartDate, request.EndDate, projectArray, survey)
			summaryVyaparmodule += obj.NoVyaparModule

			obj.Villages = getParticipantFilterVyaparGfBatchesNewCached(DB, request.StartDate, request.EndDate, projectArray, "", gfid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}

		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoVyaparsurvey = summaryNoVyaparSurvey
		response.SummaryNoVyaparCoharts = summaryNoVyaparCoharts
		response.SummaryVyaparEnroll = summaryVyaparEnroll
		response.SummaryVyaparModule = summaryVyaparmodule

		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	} else if request.RoleID == 6 {
		filter := ""
		empid := strconv.Itoa(request.EmpID)

		filter += " AND tp.user_id = " + empid

		projectList := "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
			"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.session_type in (16,17,18,19,20,21) and ((tp.date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter

		fmt.Println(projectList)

		tRes, err := DB.Query(projectList)

		if err != nil {
			log.Fatal(err)
		}
		// summaryProjectsArray := make([]int, 0)
		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoVyaparSurvey := 0
		summaryVyaparEnroll := 0
		summaryVyaparmodule := 0

		summaryNoVyaparCoharts := 0
		for tRes.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := tRes.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Fatal(err)
			}

			obj := VProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}

			// summaryProjectsArray = append(summaryProjectsArray, obj.ID)
			projectArray := []int{obj.ID}

			obj.Target = getTrainerTarget(DB, request.EmpID, projectArray)
			summaryTarget += obj.Target

			vyaparenroll := " and t.gelathi_id = " + strconv.Itoa(request.EmpID)
			filterFn := " AND b.gfId = " + strconv.Itoa(request.EmpID)
			vyaparcoharts := " AND tbl_poa.user_id = " + strconv.Itoa(request.EmpID)
			obj.Actual = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, vyaparenroll)
			summaryActuals += obj.Actual
			obj.NoofVyaparCohorts = getParticipantFilternoVyaparCoharts(DB, request.StartDate, request.EndDate, projectArray, vyaparcoharts)
			summaryNoVyaparCoharts += obj.NoofVyaparCohorts
			// empid := strconv.Itoa(request.EmpID)
			obj.Villages = GfVyaparVillageCount(DB, request.StartDate, request.EndDate, projectArray, empid)
			summaryVillages += obj.Villages
			obj.NoVyaparModule = getParticipantFilterVyaprModule(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryVyaparmodule += obj.NoVyaparModule
			obj.VyaparEnroll = getParticipantFilternoVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, vyaparenroll)
			summaryVyaparEnroll += obj.VyaparEnroll
			obj.NoofVyaparSurvey = getParticipantFilternoVyaparsurvey(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryNoVyaparSurvey += obj.NoofVyaparSurvey
			response.Data = append(response.Data, obj)
		}
		response.SummaryTarget = summaryTarget
		response.SummaryNoVyaparsurvey = summaryNoVyaparSurvey
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoVyaparCoharts = summaryNoVyaparCoharts
		response.SummaryVyaparEnroll = summaryVyaparEnroll
		response.SummaryVyaparModule = summaryVyaparmodule
		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	} else if request.RoleID == 13 {
		filter := ""
		filterFn := ""
		survey := ""
		vyaparcoharts := ""
		village := ""

		if gfid > 0 {
			filter += " AND tp.user_id = " + request.GfId
			filterFn = " AND t.gelathi_id = " + (request.GfId)
			survey = " AND b.gfId = " + (request.GfId)
			vyaparcoharts = " AND tbl_poa.user_id = " + (request.GfId)
			village = " and tp.user_id = " + request.GfId

		}
		gflid := ""
		gflid += " AND p.gfl_id = " + strconv.Itoa(request.EmpID)

		projectList := "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
			"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.VyaparEnrollment=1 or t.new_vyapar=1)  and tp.project_id != ''" + gflid

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT COALESCE(tp.project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
				"from tbl_poa tp inner join project p on p.id = tp.project_id join training_participants t on t.tb_id=tp.tb_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.session_type in (16,17,18,19,20,21) and ((tp.date BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + gflid
		}
		fmt.Println(projectList)

		rows, err := DB.Query(projectList)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// summaryProjectsArray := make([]int, 0)
		summaryTarget := 0
		summaryActuals := 0

		summaryVillages := 0
		summaryNoVyaparSurvey := 0
		summaryVyaparEnroll := 0
		summaryVyaparmodule := 0

		summaryNoVyaparCoharts := 0

		for rows.Next() {
			var prList struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				StartDate string `json:"start_date"`
				EndDate   string `json:"end_date"`
			}
			err := rows.Scan(&prList.ID, &prList.Name, &prList.StartDate, &prList.EndDate)
			if err != nil {
				log.Fatal(err)
			}

			obj := VProject{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}
			fmt.Println(prList.ID)

			// summaryProjectsArray = append(summaryProjectsArray, obj.ID)
			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilternoGflVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, filterFn)
			summaryActuals += obj.Actual

			obj.VyaparEnroll = getParticipantFilternoGflVyaparenroll(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, filterFn)
			summaryVyaparEnroll += obj.VyaparEnroll

			obj.NoofVyaparCohorts = getParticipantFilternoGflVyaparCoharts(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, vyaparcoharts)
			summaryNoVyaparCoharts += obj.NoofVyaparCohorts

			obj.NoofVyaparSurvey = getParticipantFilternoGflVyaparsurvey(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, survey)
			summaryNoVyaparSurvey += obj.NoofVyaparSurvey

			obj.NoVyaparModule = getParticipantFilterGflVyaprModule(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, survey)
			summaryVyaparmodule += obj.NoVyaparModule
			obj.Villages = GflVyaparnewVillageCount(DB, request.StartDate, request.EndDate, projectArray, request.EmpID, village)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}
		defer rows.Close()

		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNoVyaparsurvey = summaryNoVyaparSurvey
		response.SummaryNoVyaparCoharts = summaryNoVyaparCoharts
		response.SummaryVyaparEnroll = summaryVyaparEnroll
		response.SummaryVyaparModule = summaryVyaparmodule

		response.Code = 200
		response.Count = len(response.Data)
		response.Success = true
		response.Message = "Successfully"

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	}

}
