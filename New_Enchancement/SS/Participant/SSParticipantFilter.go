package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Target     int    `json:"target"`
	Actual     int    `json:"actual"`
	Day2       int    `json:"day2"`
	Nobatches  int    `json:"nobatches"`
	Villages   int    `json:"villages"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SelectType string `json:"select_type"`
	NoofSurvey int    `json:"noofsurvey"`
}

type Response struct {
	SummaryTarget    int       `json:"summary_target"`
	SummaryVillages  int       `json:"summary_villages"`
	SummaryActuals   int       `json:"summary_actual"`
	SummaryDay2      int       `json:"summary_day2"`
	SummaryNosurvey  int       `json:"summary_nosurvey"`
	SummaryNoBatches int       `json:"summary_nobatches"`
	Data             []Project `json:"data"`
	Code             int       `json:"code"`
	Count            int       `json:"count"`
	Success          bool      `json:"success"`
	Message          string    `json:"message"`
}

func SSCounts(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	var request struct {
		RoleID    int    `json:"role_id"`
		EmpID     int    `json:"emp_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		TrainerID string `json:"trainerId,omitempty"`
		Project   string `json:"project"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	response := Response{
		Data: make([]Project, 0),
	}
	trainerid, _ := strconv.Atoi(request.TrainerID)

	if request.RoleID == 1 || request.RoleID == 9 || request.RoleID == 3 || request.RoleID == 4 || request.RoleID == 12 || request.RoleID == 11 {
		filter := ""
		filterFn := ""
		batches := ""

		if request.RoleID == 3 {
			opsIds := getReportingOpsManagers(DB, request.EmpID)
			filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

		} else if request.RoleID == 12 {
			opsIds := getOpsManagers(DB, request.EmpID)
			filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))

		} else if request.RoleID == 4 {
			filter = " AND p.operations_manager = " + strconv.Itoa(request.EmpID)

		}

		if trainerid > 0 {
			filter += " AND tp.trainer_id = " + (request.TrainerID)
			filterFn = " AND tp.trainer_id = " + (request.TrainerID)
			batches = " AND tp.user_id = " + (request.TrainerID)
		}

		projectList := "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
			"INNER JOIN project p ON tp.project_id = p.id " +
			"WHERE tp.project_id != '' " + filter

		if request.StartDate != "" && request.EndDate != "" {
			projectList = "SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp " +
				"INNER JOIN project p ON tp.project_id = p.id " +
				"WHERE ((tp.participant_day1 BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "') OR (tp.participant_day2 BETWEEN '" + request.StartDate + "' AND '" + request.EndDate + "'))" + filter
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
		summaryDay1 := 0

		summaryVillages := 0
		summaryNoSurvey := 0

		summaryNoBatches := 0

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

			obj := Project{
				ID:         prList.ID,
				Name:       prList.Name,
				StartDate:  prList.StartDate,
				EndDate:    prList.EndDate,
				SelectType: "1",
			}
			projectArray := []int{obj.ID}

			obj.Target = getTarget(DB, request.StartDate, request.EndDate, projectArray)
			summaryTarget += obj.Target

			obj.Actual = getParticipantFilterActual(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryActuals += obj.Actual

			day1Count := getParticipantFilterDay1Count(DB, request.StartDate, request.EndDate, projectArray, filterFn)
			summaryDay1 += day1Count

			if day1Count > 0 {
				day2Turnout := float64(obj.Actual) / float64(day1Count)
				obj.Day2 = int(math.Round(day2Turnout * 100))
			} else {
				obj.Day2 = 0
			}

			obj.Nobatches = getParticipantFilternoofbatches(DB, request.StartDate, request.EndDate, projectArray, batches)
			summaryNoBatches += obj.Nobatches

			obj.NoofSurvey = getParticipantFilternoOfselfshakthisurvey(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryNoSurvey += obj.NoofSurvey
			obj.Villages = getParticipantFilterTrainingBatchesNew(DB, request.StartDate, request.EndDate, projectArray, "", trainerid)
			summaryVillages += obj.Villages
			response.Data = append(response.Data, obj)
		}

		if summaryDay1 > 0 {
			day2Turnout := float64(summaryActuals) / float64(summaryDay1)
			response.SummaryDay2 = int(math.Round(day2Turnout * 100))
		} else {
			response.SummaryDay2 = 0
		}
		response.SummaryTarget = summaryTarget
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		response.SummaryNosurvey = summaryNoSurvey
		response.SummaryNoBatches = summaryNoBatches

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
	} else if request.RoleID == 5 {

		projectList := fmt.Sprintf("SELECT DISTINCT tp.project_id AS id, p.projectName AS name, p.startDate, p.endDate FROM training_participants tp "+
			"INNER JOIN project p ON tp.project_id = p.id "+
			"WHERE ((tp.participant_day1 BETWEEN '%s' AND '%s') OR (tp.participant_day2 BETWEEN '%s' AND '%s')) AND tp.trainer_id = %d",
			request.StartDate, request.EndDate, request.StartDate, request.EndDate, request.EmpID)

		tRes, err := DB.Query(projectList)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		defer tRes.Close()
		summaryTarget := 0
		summaryActuals := 0
		summaryDay1 := 0
		summaryNoSurvey := 0
		summaryVillages := 0
		summaryNoBatches := 0
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

			obj := Project{
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

			day1Count := getParticipantFilterDay1Count(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryDay1 += day1Count

			if day1Count > 0 {
				day2Turnout := float64(obj.Actual) / float64(day1Count)
				obj.Day2 = int(math.Round(day2Turnout * 100))
			} else {
				obj.Day2 = 0
			}
			batches := " and tp.user_id = " + strconv.Itoa(request.EmpID)
			obj.Nobatches = getParticipantFilternoofbatches(DB, request.StartDate, request.EndDate, projectArray, batches)
			summaryNoBatches += obj.Nobatches
			empid := strconv.Itoa(request.EmpID)
			obj.Villages = TrainerVillageCount(DB, request.StartDate, request.EndDate, projectArray, empid)
			summaryVillages += obj.Villages

			obj.NoofSurvey = getParticipantFilternoOfselfshakthisurvey(DB, request.StartDate, request.EndDate, projectArray, filter)
			summaryNoSurvey += obj.NoofSurvey
			response.Data = append(response.Data, obj)
		}
		response.SummaryTarget = summaryTarget
		response.SummaryNosurvey = summaryNoSurvey
		response.SummaryVillages = summaryVillages
		response.SummaryActuals = summaryActuals
		if summaryDay1 > 0 {
			day2Turnout := float64(summaryActuals) / float64(summaryDay1)
			response.SummaryDay2 = int(math.Round(day2Turnout * 100))
		} else {
			response.SummaryDay2 = 0
		}
		response.SummaryNoBatches = summaryNoBatches
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

		type ProjectDataN struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			CircleMeet      int    `json:"circle_meet"`
			Circles         int    `json:"circles"`
			VillageVisit    int    `json:"villagevisit"`
			Beehive         int    `json:"beehive"`
			Enroll          int    `json:"enroll"`
			GreenMotivators int    `json:"greenMotivators"`
			Vyapar          int    `json:"vyapar"`
		}
		type aResponse struct {
			SummaryCircleMeet   int            `json:"summary_circle_meet"`
			SummaryCircles      int            `json:"summary_circles"`
			SummaryVillageVisit int            `json:"summary_villagevisit"`
			SummaryBeehive      int            `json:"summary_beehive"`
			SummaryEnroll       int            `json:"summary_enroll"`
			Data                []ProjectDataN `json:"data"`
			Code                int            `json:"code"`
			Success             bool           `json:"success"`
			Message             string         `json:"message"`
		}

		// filter := ""
		filters := ""
		// filterG := ""
		// filterV := ""

		if request.StartDate != "" && request.EndDate != "" {
			// filter = " and tp.participant_day2 BETWEEN '" + request.StartDate + "' and '" + request.EndDate + "'"
			filters = " and tp.date BETWEEN '" + request.StartDate + "' and '" + request.EndDate + "'"
			// filterG = " and tp.GreenMotivatorsDate BETWEEN '" + request.StartDate + "' and '" + request.EndDate + "'"
			// filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + request.StartDate + "' and '" + request.EndDate + "'"
		}

		summaryCircleMeet := 0
		summaryCircles := 0
		summaryVillageVisit := 0
		summaryBeehive := 0
		summaryEnroll := 0

		// Get employee IDs
		employeeIDs := []int{}
		rows, err := DB.Query("SELECT id FROM employee WHERE status = 1 AND supervisorId = ?", request.EmpID)
		if err != nil {
			log.Println("Failed to query employee IDs:", err)
			response.Code = http.StatusInternalServerError
			response.Success = false
			response.Message = "Failed to query employee IDs"

		}
		defer rows.Close()

		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				log.Println("Failed to scan employee ID:", err)
				continue
			}
			employeeIDs = append(employeeIDs, id)
		}

		// Get project IDs and names
		getProjs := "SELECT project_id AS id, projectName AS name FROM tbl_poa tp " +
			"INNER JOIN project p ON p.id = tp.project_id " +
			"WHERE p.gfl_id = ?" + filters + " GROUP BY tp.project_id"

		rows, err = DB.Query(getProjs, request.EmpID)
		if err != nil {
			log.Println("Failed to query projects:", err)
			response.Code = http.StatusInternalServerError
			response.Success = false
			response.Message = "Failed to query projects"

		}
		defer rows.Close()
		var response aResponse

		projectDataList := []ProjectDataN{}
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Println("Failed to scan project:", err)
				continue
			}

			projectData := ProjectDataN{
				ID:   id,
				Name: name,
			}

			// Circle Meet
			prjFilter := " AND p.id = " + strconv.Itoa(id)
			circleProjMeet := getGFDataN(DB, prjFilter, 1, employeeIDs)
			projectData.CircleMeet = circleProjMeet
			summaryCircleMeet += circleProjMeet

			// Circles
			projectData.Circles = getGFCircleN(DB, prjFilter, employeeIDs)
			summaryCircles += projectData.Circles

			// Village Visit
			villageProjVisit := getGFDataN(DB, prjFilter, 2, employeeIDs)
			projectData.VillageVisit = villageProjVisit
			summaryVillageVisit += villageProjVisit

			// Beehive
			beehiveProj := getGFDataN(DB, prjFilter, 3, employeeIDs)
			projectData.Beehive = beehiveProj
			summaryBeehive += beehiveProj

			// Enroll
			participantFilter := ""
			if participantFilter != "" {
				projEnrolled := getParticipantFilterGfEnrolledN(DB, prjFilter, employeeIDs, request.StartDate, request.EndDate)
				projectData.Enroll = projEnrolled
				summaryEnroll += projEnrolled
			} else {
				projEnrolled := getGfEnrolledN(DB, prjFilter, employeeIDs)
				projectData.Enroll = projEnrolled
				summaryEnroll += projEnrolled
			}
			query := "SELECT GROUP_CONCAT(DISTINCT prj.id) as ids FROM project_emps em_pr " +
				"LEFT JOIN project prj ON em_pr.project_id = prj.id WHERE emp_id = ?"
			row := DB.QueryRow(query, request.EmpID)

			var projectIDString string
			err := row.Scan(&projectIDString)
			if err != nil {
				log.Println("Failed to fetch project IDs:", err)

			}

			projectIDArr := strings.Split(projectIDString, ",")
			for _, projectIDStr := range projectIDArr {
				projectID := 0
				_, err := fmt.Sscanf(projectIDStr, "%d", &projectID)
				if err != nil {
					log.Println("Failed to parse project ID:", err)
					continue
				}
			}

			

			projectDataList = append(projectDataList, projectData)
		}

		response.SummaryCircleMeet = summaryCircleMeet
		response.SummaryCircles = summaryCircles
		response.SummaryVillageVisit = summaryVillageVisit
		response.SummaryBeehive = summaryBeehive
		response.SummaryEnroll = summaryEnroll

		response.Data = projectDataList
		response.Code = http.StatusOK
		response.Success = true
		response.Message = "Successfully"
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
