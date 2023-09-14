package Team_1

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getAssociatedProjectList(db *sql.DB, projId int) ([]int, error) {
	projArray := []int{}
	if projId > 0 {
		getProjList := fmt.Sprintf("SELECT associatedProject FROM project_association WHERE projectId IN (%d)", projId)
		projArray = append(projArray, projId)
		rows, err := db.Query(getProjList)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var associatedProject int
			err := rows.Scan(&associatedProject)
			if err != nil {
				return nil, err
			}
			projArray = append(projArray, associatedProject)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return projArray, nil
}

func GetNoOfVyaparSurvey(db *sql.DB, startDate string, endDate string, gfId string) int {
	var getActualsQuery string
	var noofvyaparsurvey int

	// if len(projectArray) > 0 {
	getActualsQuery = "select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline"
	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline where entry_date BETWEEN '%s' AND '%s'", startDate, endDate)
	} else {
		if gfId != "" {
			getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline tp WHERE entry_date BETWEEN '%s' AND '%s' and gfid '%s'", startDate, endDate, gfId)

		}
	}
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofvyaparsurvey

}
func GetNoOfgreenSurvey(db *sql.DB, startDate string, endDate string, gfId string) int {
	var getActualsQuery string
	var noofgreensurvey int

	// if len(projectArray) > 0 {
	getActualsQuery = "select count(id) as noofgreensurvey from GreenBaselineSurvey"
	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("select count(id) as noofgreensurvey from GreenBaselineSurvey where entry_date BETWEEN '%s' AND '%s'", startDate, endDate)
	} else {
		if gfId != "" {
			getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofgreensurvey from GreenBaselineSurvey tp WHERE entry_date BETWEEN '%s' AND '%s' and gfid '%s'", startDate, endDate, gfId)

		}
	}
	err := db.QueryRow(getActualsQuery).Scan(&noofgreensurvey)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofgreensurvey

}
func GetNoOfSporthiSurvey(db *sql.DB, startDate string, endDate string, gfId string) int {
	var getActualsQuery string
	var noofsporthisurvey int

	// if len(projectArray) > 0 {
	getActualsQuery = "select count(id) as noofvyaparsurvey from SpoorthiBaselineQuestionnaire"
	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from SpoorthiBaselineQuestionnaire where entry_date BETWEEN '%s' AND '%s'", startDate, endDate)
	} else {
		if gfId != "" {
			getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofvyaparsurvey from SpoorthiBaselineQuestionnaire tp WHERE entry_date BETWEEN '%s' AND '%s' and gfid '%s'", startDate, endDate, gfId)

		}
	}
	err := db.QueryRow(getActualsQuery).Scan(&noofsporthisurvey)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofsporthisurvey

}

func NoofVyaparCohorts(db *sql.DB, startDate string, endDate string, project_id string) int {
	var getActualsQuery string
	var noofvyaparcohorts, funderId int

	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf(`SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts, project.funderId as funderId 
		FROM tbl_poa 
		INNER JOIN project on project.id=tbl_poa.project_id 
		WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and date BETWEEN '%s' AND '%s' group by project.funderID order by project.funderID`, startDate, endDate)
	} else if project_id != "" {
		getActualsQuery = fmt.Sprintf(`SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts, project.funderId as funderId 
		FROM tbl_poa 
		INNER JOIN project on project.id=tbl_poa.project_id 
		WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) 
		
		 and project_id '%s' group by project.funderID order by project.funderID`, project_id)

	} else {
		getActualsQuery = `SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts, project.funderId as funderId 
		FROM tbl_poa 
		INNER JOIN project on project.id=tbl_poa.project_id 
		WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) 
		GROUP BY project.funderID 
		ORDER BY project.funderID`

	}
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparcohorts, &funderId)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofvyaparcohorts

}

func NoofGelathiCohorts(db *sql.DB, startDate string, endDate string, project_id string) int {
	var getActualsQuery string
	var noofGelathicohorts, funderId int

	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofspoorthicohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 4 OR tbl_poa.session_type = 5 OR tbl_poa.session_type = 6 OR tbl_poa.session_type = 7 OR tbl_poa.session_type = 8 OR tbl_poa.session_type = 9) and date BETWEEN '%s' AND '%s' group by project.funderID order by project.funderID", startDate, endDate)
	} else if project_id != "" {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofspoorthicohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 4 OR tbl_poa.session_type = 5 OR tbl_poa.session_type = 6 OR tbl_poa.session_type = 7 OR tbl_poa.session_type = 8 OR tbl_poa.session_type = 9) and project_id '%s' group by project.funderID order by project.funderID", project_id)

	} else {
		getActualsQuery = "SELECT COUNT(tbl_poa.session_type) AS noofspoorthicohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 4 OR tbl_poa.session_type = 5 OR tbl_poa.session_type = 6 OR tbl_poa.session_type = 7 OR tbl_poa.session_type = 8 OR tbl_poa.session_type = 9) group by project.funderID order by project.funderID"

	}
	err := db.QueryRow(getActualsQuery).Scan(&noofGelathicohorts, &funderId)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofGelathicohorts

}
func NoofGreenCohorts(db *sql.DB, startDate string, endDate string, project_id string) int {
	var getActualsQuery string
	var noofgreencohorts, funderId int

	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) and date BETWEEN '%s' AND '%s' group by project.funderID order by project.funderID", startDate, endDate)
	} else if project_id != "" {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) and project_id '%s' group by project.funderID order by project.funderID", project_id)

	} else {
		getActualsQuery = "SELECT COUNT(tbl_poa.session_type) AS noofgreencohorts, project.funderId as funderId FROM tbl_poa INNER JOIN project on project.id=tbl_poa.project_id where tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) group by project.funderID order by project.funderID"

	}
	err := db.QueryRow(getActualsQuery).Scan(&noofgreencohorts, &funderId)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofgreencohorts

}

func GetNoofVyaparModuleCompleted(db *sql.DB) int {
	var getActualsQuery string
	var noofvyaparmodulecompleted int

	// if len(projectArray) > 0 {
	// if startDate != "" && endDate != "" {
	getActualsQuery = "select count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) from BuzzVyaparProgramBaseline"
	//  } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	// }
	//  }
	// } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	// }
	// }
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparmodulecompleted)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofvyaparmodulecompleted

}

func GetNoofSporthiModuleCompleted(db *sql.DB) int {
	var getActualsQuery string
	var noofsporthimodulecompleted int

	// if len(projectArray) > 0 {
	// if startDate != "" && endDate != "" {
	getActualsQuery = "select count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) from SpoorthiBaselineQuestionnaire"
	//  } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	// }
	//  }
	// } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	// }
	// }
	err := db.QueryRow(getActualsQuery).Scan(&noofsporthimodulecompleted)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofsporthimodulecompleted

}
func GetNoofGreenModuleCompleted(db *sql.DB) int {
	var getActualsQuery string
	var noofgreenmodulecompleted int

	// if len(projectArray) > 0 {
	// if startDate != "" && endDate != "" {
	getActualsQuery = "select count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) from GreenBaselineSurvey"
	//  } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	// }
	//  }
	// } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	// }
	// }
	err := db.QueryRow(getActualsQuery).Scan(&noofgreenmodulecompleted)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofgreenmodulecompleted

}

func getTarget(db *sql.DB, startDate string, endDate string, projectArray []int) int {
	var getTargetQuery string
	var target sql.NullInt64

	if len(projectArray) > 0 {
		getTargetQuery = fmt.Sprintf("SELECT COALESCE(sum(training_target), 0) as target from project p where id in (%s)", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	} else {
		getTargetQuery = fmt.Sprintf("SELECT sum(training_target) as target from project p where (startDate >= '%s' and endDate <= '%s')", startDate, endDate)
	}

	// ...

	err := db.QueryRow(getTargetQuery).Scan(&target)
	if err != nil {
		log.Println("ERROR>>", err)
	}

	if target.Valid {
		return int(target.Int64)
	} else {
		return 0
	}

}

func getParticipantFilterDay1Count(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var day1CountQuery string
	if len(projectArray) > 0 {
		projectIDs := make([]string, len(projectArray))
		for i, projectID := range projectArray {
			projectIDs[i] = strconv.Itoa(projectID)
		}
		day1CountQuery = "select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (" + strings.Join(projectIDs, ",") + ") and  participant_day1 >= '" + startDate + "' and participant_day1 <= '" + endDate + "'" + filter
	} else {
		day1CountQuery = "select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and participant_day1 >= '" + startDate + "' and participant_day1 <= '" + endDate + "'" + filter
	}
	row := db.QueryRow(day1CountQuery)
	var day1Count int
	err := row.Scan(&day1Count)
	if err != nil {
		log.Println(err)
	}
	return day1Count
}

// helper function to remove duplicates from an int slice
func uniqueIntSlice(slice []int) []int {
	keys := make(map[int]bool)
	var unique []int
	for _, val := range slice {
		if _, value := keys[val]; !value {
			keys[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}

// helper function to convert an int slice to a string of comma-separated values
func intSliceToString(slice []int) string {
	var strSlice []string
	for _, val := range slice {
		strSlice = append(strSlice, strconv.Itoa(val))
	}
	return strings.Join(strSlice, ", ")
}

func Vyapar(db *sql.DB, startDate string, endDate string, projectArray []int, funderId string, filter string) int {
	if funderId != "" {
		getProj := fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %s", funderId)
		projResult, _ := db.Query(getProj)
		for projResult.Next() {
			var id int
			var startDate string
			var endDate string
			projResult.Scan(&id, &startDate, &endDate)
			projectArray = append(projectArray, id)
		}
	}

	for _, proj := range projectArray {
		// check if there are any associated project for each project
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	projectArray = uniqueIntSlice(projectArray)
	var gelatiCountQuery string
	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN '%s' AND '%s'", intSliceToString(projectArray), startDate, endDate)
		} else {
			if funderId != "" {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT min(startDate) from project p where funderID = %s and endDate >= CURRENT_DATE()) and (SELECT max(endDate) from project p where funderID = %s and endDate >= CURRENT_DATE())", intSliceToString(projectArray), funderId, funderId)
			} else {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT min(startDate) from project p where id IN (%s)) and (SELECT max(endDate) from project p where id IN (%s))", intSliceToString(projectArray), intSliceToString(projectArray), intSliceToString(projectArray))
			}
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as Vyapar FROM training_participants tp INNER JOIN project p ON p.id = tp.VyaparEnrollmentEnrolledProject WHERE VyaparEnrollment = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}

	rows, _ := db.Query(gelatiCountQuery)
	defer rows.Close()
	var VyaparCount int
	for rows.Next() {
		rows.Scan(&VyaparCount)
	}
	return VyaparCount
}

func ParticipantFiltergreenMotivators(db *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	var gelatiCountQuery string
	var greenMoti int
	var err error

	if len(projectArray) > 0 {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp where GreenMotivators = 1 and GreenMotivatorsEnrolledProject in (%s) and (date(GreenMotivatorsDate) >= '%s' and date(GreenMotivatorsDate) <= '%s') %s", intsToString(projectArray), startDate, endDate, filter)
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as greenMoti FROM training_participants tp inner join project p on p.id = tp.GreenMotivatorsEnrolledProject where GreenMotivators = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	err = db.QueryRow(gelatiCountQuery).Scan(&greenMoti)
	if err != nil {
		return 0, err
	}

	return greenMoti, nil
}

func TrainerVillageCount(db *sql.DB, startDate string, endDate string, projectArray []int, user_id string) int {
	var villageCount, subVillageCount int

	if len(projectArray) > 0 {
		// Build the project IDs string for the SQL query
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")

		villageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE type = 1 AND added = 0 AND project_id IN (%s) AND user_id = %s", projectIDs, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE type = 1 AND added = 0 AND project_id IN (%s) AND user_id = %s", projectIDs, user_id)

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
		}

		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			log.Println("TrainerVillageCount", err)
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			log.Println("TrainerVillageCount", err)
		}
	} else {
		villageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' AND user_id = %s", startDate, endDate, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' AND user_id = %s", startDate, endDate, user_id)

		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			log.Println("TrainerVillageCount", err)
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			log.Println("TrainerVillageCount", err)
		}
	}

	return villageCount + subVillageCount
}

func getParticipantFilterTrainingBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, trainerId int) int {
	if trainerId > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(trainerId)
	}
	var villageQuery, subVillageQuery string
	if len(projectArray) > 0 {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
	} else {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
	}
	villageResult, err := db.Query(villageQuery)
	if err != nil {
		panic(err.Error())
	}
	subVillageResult, err := db.Query(subVillageQuery)
	if err != nil {
		panic(err.Error())
	}
	var village, subVillage int
	for villageResult.Next() {
		err := villageResult.Scan(&village)
		if err != nil {
			panic(err.Error())
		}
	}
	for subVillageResult.Next() {
		err := subVillageResult.Scan(&subVillage)
		if err != nil {
			panic(err.Error())
		}
	}
	return village + subVillage
}

func intsToString(ints []int) string {
	var stringSlice []string
	for _, i := range ints {
		stringSlice = append(stringSlice, strconv.Itoa(i))
	}
	return strings.Join(stringSlice, ",")
}

func getOpsManagers(db *sql.DB, empId int) []int {
	getOpsIds := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d AND empRole = 4", empId)
	ids := make([]int, 0)
	res, err := db.Query(getOpsIds)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			log.Println("ERROR>>", err)
		}
		ids = append(ids, id)
	}
	return ids
}

func getReportingOpsManagers(db *sql.DB, empId int) []int {
	ids := []int{}

	getOpsIds := "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
	res, err := db.Query(getOpsIds, empId)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Println("ERROR>>", err)
		}
		ids = append(ids, id)
	}

	getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 12"
	res, err = db.Query(getOpsIds, empId)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Println("ERROR>>", err)
		}
		ids = append(ids, id)
		som := id
		getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
		res, err = db.Query(getOpsIds, som)
		if err != nil {
			log.Println("ERROR>>", err)
		}
		defer res.Close()

		for res.Next() {
			var id int
			err := res.Scan(&id)
			if err != nil {
				log.Println("ERROR>>", err)
			}
			ids = append(ids, id)
		}
	}

	return ids
}

func getTrainerTarget(db *sql.DB, empId int, projectArray []int) int {
	targetQuery := fmt.Sprintf("SELECT sum(target) as total from project_emps pe where emp_id = %d and project_id in (%s)",
		empId, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	var total sql.NullInt64
	err := db.QueryRow(targetQuery).Scan(&total)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	if !total.Valid {
		return 0
	}
	return int(total.Int64)
}

func getGFDataN(db *sql.DB, filter string, sessionType int, empId []int) int {
	filter += fmt.Sprintf(" and tp.user_id in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(empId)), ","), "[]"))
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
	var visit int
	row := db.QueryRow(getVisit)
	err := row.Scan(&visit)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return visit
}

func getGFCircleN(db *sql.DB, filter string, empId []int) int {
	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id in (%s)%s", strings.Trim(strings.Replace(fmt.Sprint(empId), " ", ",", -1), "[]"), filter)
	row := db.QueryRow(getCircle)
	var visit int
	row.Scan(&visit)
	return visit
}

func getGfEnrolledN(db *sql.DB, filter string, empId []int) int {
	empIdStr := make([]string, len(empId))
	for i, id := range empId {
		empIdStr[i] = strconv.Itoa(id)
	}
	query := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
		"inner join project p on tp.project_id = p.id " +
		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " + filter
	row := db.QueryRow(query)
	var enrolled int
	row.Scan(&enrolled)
	return enrolled
}

func Kann(db *sql.DB, filter string, empId []int, startDate, endDate time.Time) int {
	enrolled := 0
	empIdStr := make([]string, len(empId))
	for i, v := range empId {
		empIdStr[i] = strconv.Itoa(v)
	}
	getEnrolled := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
		"inner join project p on tp.project_id = p.id " +
		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " +
		"and ((participant_day1 >= ? and participant_day1 <= ?) or (participant_day2 >= ? and participant_day2 <= ?)) " + filter
	row := db.QueryRow(getEnrolled, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	row.Scan(&enrolled)
	return enrolled
}

func NoofBatches(db *sql.DB, emp_id int, project_id int) int {
	var getActualsQuery string
	var noofbatches int

	if emp_id != 0 {
		getActualsQuery = fmt.Sprintf(`SELECT COUNT(DISTINCT tb_id) FROM tbl_poa WHERE user_id = %d`, emp_id)
	} else if project_id != 0 {
		getActualsQuery = fmt.Sprintf(`SELECT COUNT(DISTINCT tb_id) FROM tbl_poa WHERE project_id = %d`, project_id)

	} else if emp_id != 0 && project_id != 0 {
		getActualsQuery = fmt.Sprintf(`SELECT COUNT(DISTINCT tb_id) FROM tbl_poa WHERE user_id = %d AND project_id = %d`, emp_id, project_id)

	} else {
		getActualsQuery = `SELECT COUNT(DISTINCT tb_id) FROM tbl_poa`

	}
	err := db.QueryRow(getActualsQuery).Scan(&noofbatches)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return noofbatches
}

func getParticipantFilterGfEnrolledN(db *sql.DB, filter string, empIDs []int, startDate, endDate string) int {
	empIDStr := make([]string, len(empIDs))
	for i, id := range empIDs {
		empIDStr[i] = strconv.Itoa(id)
	}
	empIDList := strings.Join(empIDStr, ",")

	query := "SELECT COUNT(tp.id) as enrolled FROM training_participants tp " +
		"INNER JOIN project p ON tp.project_id = p.id " +
		"WHERE enroll = 1 AND gelathi_id IN (" + empIDList + ") " +
		"AND ((participant_day1 >= ? AND participant_day1 <= ?) OR (participant_day2 >= ? AND participant_day2 <= ?)) " + filter

	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Println("Failed to parse start date:", err)
		return 0
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Println("Failed to parse end date:", err)
		return 0
	}

	rows, err := db.Query(query, startTime, endTime, startTime, endTime)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return 0
	}
	defer rows.Close()

	var enrolledCount int
	if rows.Next() {
		err = rows.Scan(&enrolledCount)
		if err != nil {
			log.Println("Failed to scan result:", err)
			return 0
		}
	}

	return enrolledCount
}
func getParticipantFilterActual(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) and participant_day2 >= '%s' and participant_day2 <= '%s' %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), startDate, endDate, filter)
	} else {
		getActualsQuery = fmt.Sprintf("select count(tp.id) as actual from training_participants tp inner join project p on p.id = tp.project_id where day2 = 1 and participant_day2 >= '%s' and participant_day2 <= '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Println("ERROR>>", err)
	}
	return actual
}

func TrainerActual(db *sql.DB, startDate string, endDate string, projectArray []int, trainer_id string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")
		getActualsQuery = fmt.Sprintf("SELECT COUNT(id) AS actual FROM bdms_staff.training_participants WHERE trainer_id IN (%s) AND project_id IN (%s)", trainer_id, projectIDs)
		if startDate != "" && endDate != "" {
			getActualsQuery += fmt.Sprintf(" AND participant_day2 BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(id) AS actual FROM bdms_staff.training_participants WHERE trainer_id IN (%s)", trainer_id)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Println("getActual", err)
	}
	return actual
}

func getParticipantFilternoofbatches(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getnoofbatches string
	var noofbatches int

	if len(projectArray) > 0 {

		getnoofbatches = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa tp JOIN project t3 ON tp.project_id = t3.id WHERE date BETWEEN '%s' AND '%s' AND project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)

	} else {
		getnoofbatches = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa tp JOIN project t3 ON tp.project_id = t3.id WHERE date BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(getnoofbatches).Scan(&noofbatches)
	if err != nil {
		log.Println("getActual", err)
	}

	return noofbatches
}

func getParticipantFilternoOfselfshakthisurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var gelatiCountQuery string
	var noofselfshakthi int

	if len(projectArray) > 0 {

		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) AS actual FROM training_participants tp JOIN project p ON p.id = tp.project_id WHERE isSurveyDone = 1 AND participant_day1 BETWEEN '%s' AND '%s' AND tp.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)

	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) AS actual FROM training_participants tp JOIN project p ON p.id = tp.project_id WHERE isSurveyDone = 1 AND participant_day1 BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(gelatiCountQuery).Scan(&noofselfshakthi)
	if err != nil {
		log.Println("getActual", err)
	}
	return noofselfshakthi
}
