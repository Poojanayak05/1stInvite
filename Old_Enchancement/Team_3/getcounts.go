package Team_3

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func noOfselfshakthisurvey(db *sql.DB) int {

	var gelatiCountQuery string
	var noofselfshakthi int
	gelatiCountQuery = "SELECT COUNT(id) as noofselfsakthi FROM training_participants WHERE isSurveyDone=1"
	err := db.QueryRow(gelatiCountQuery).Scan(&noofselfshakthi)
	if err != nil {
		log.Println("noOfselfshakthisurvey", err)
	}
	return noofselfshakthi
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
		log.Println("GetNoOfVyaparSurvey", err)
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
		log.Println("GetNoOfgreenSurvey", err)
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
		log.Println("GetNoOfSporthiSurvey", err)
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
		log.Println("NoofVyaparCohorts", err)
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
		log.Println("NoofGelathiCohorts", err)
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
		log.Println("NoofGreenCohorts", err)
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
		log.Println("GetNoofVyaparModuleCompleted", err)
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
		log.Println("GetNoofSporthiModuleCompleted", err)
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
		log.Println("GetNoofGreenModuleCompleted", err)
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
		log.Println("getTarget", err)
	}

	if target.Valid {
		return int(target.Int64)
	} else {
		return 0
	}

}

func getActual(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and (participant_day2 BETWEEN '%s' and '%s' )and project_id in (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		getActualsQuery = fmt.Sprintf("select count(tp.id) as actual from training_participants tp inner join project p on p.id = tp.project_id where day2 = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Println("getActual", err)
	}
	return actual
}

func getDay1Count(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var day1CountQuery string
	if len(projectArray) > 0 {
		projectIDs := make([]string, len(projectArray))
		for i, projectID := range projectArray {
			projectIDs[i] = strconv.Itoa(projectID)
		}
		day1CountQuery = "select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (" + strings.Join(projectIDs, ",") + ")" + filter
	} else {
		day1CountQuery = "select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
	}
	row := db.QueryRow(day1CountQuery)
	var day1Count int
	err := row.Scan(&day1Count)
	if err != nil {
		log.Println(err)
	}
	return day1Count
}

func getGelathi(db *sql.DB, startDate string, endDate string, projectArray []int, GalathiId string, funderId string, filter string) int {
	var gelathiCount int
	// get associated projects for each project
	for _, proj := range projectArray {
		projs, _ := GetAssociatedProjectList(proj, db)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = uniqueIntSlice(projectArray)
	// build query
	var gelatiCountQuery string
	if len(projectArray) > 0 {

		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as gelathiCount FROM training_participants tp WHERE enroll = 1 AND enrolledProject IN (%s) AND enroll_date BETWEEN '%s' AND '%s'", intSliceToString(projectArray), startDate, endDate)
		} else {
			var dateRangeQuery []string
			if funderId != "" {
				dateRangeQuery = []string{fmt.Sprintf("(SELECT min(startDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", funderId), fmt.Sprintf("(SELECT max(endDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", funderId)}
			} else {
				dateRangeQuery = []string{"(SELECT min(startDate) FROM project WHERE endDate >= CURRENT_DATE())", "(SELECT max(endDate) FROM project WHERE endDate >= CURRENT_DATE())"}
			}
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as gelathiCount FROM training_participants tp WHERE enroll = 1 AND enrolledProject IN (%s) AND enroll_date BETWEEN %s AND %s", intSliceToString(projectArray), dateRangeQuery[0], dateRangeQuery[1])
		}
		if GalathiId != "" {
			gelatiCountQuery += strings.Replace(GalathiId, "trainer_id", "gelathi_id", -1)
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as gelathiCount FROM training_participants tp INNER JOIN project p ON p.id = tp.enrolledProject WHERE enroll = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}
	// execute query
	row := db.QueryRow(gelatiCountQuery)
	row.Scan(&gelathiCount)
	return gelathiCount
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

func greenMotivators(con *sql.DB, startDate string, endDate string, projectArray []int, funderId string, filter string) int {
	if funderId != "" {
		getProj := fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %s", funderId)
		projResult, err := con.Query(getProj)
		if err != nil {
			panic(err)
		}
		defer projResult.Close()

		for projResult.Next() {
			var id int
			var startDate string
			var endDate string

			err = projResult.Scan(&id, &startDate, &endDate)
			if err != nil {
				panic(err)
			}

			projectArray = append(projectArray, id)
		}
	}

	for _, proj := range projectArray {
		projs, _ := GetAssociatedProjectList(proj, con)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = uniqueIntSlice(projectArray)

	var gelatiCountQuery string
	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN '%s' AND '%s'", sliceToString(projectArray), startDate, endDate)
		} else {
			if funderId != "" {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE()) AND (SELECT MAX(endDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", sliceToString(projectArray), funderId, funderId)
			} else {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE id IN (%s)) AND (SELECT MAX(endDate) FROM project p WHERE id IN (%s))", sliceToString(projectArray), sliceToString(projectArray), sliceToString(projectArray))
			}
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as greenMoti FROM training_participants tp INNER JOIN project p ON p.id = tp.GreenMotivatorsEnrolledProject WHERE GreenMotivators = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}

	rows, err := con.Query(gelatiCountQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var greenMoti int
		err = rows.Scan(&greenMoti)
		if err != nil {
			panic(err)
		}

		return greenMoti
	}

	return 0
}

func sliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
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
		defer projResult.Close()
	}

	for _, proj := range projectArray {
		// check if there are any associated project for each project
		projs, _ := GetAssociatedProjectList(proj, db)
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

func getTrainingBatches(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	//previously considering training batch count for villages.
	/*if (count($projectArray) > 0){
		$tbQuery = "SELECT count(*) as tbCount from
						(SELECT max(id) from tbl_poa tp where check_out is not null and `type`  = 1 and added  = 0 and project_id in (".implode(',',$projectArray).") $filter
						GROUP BY tb_id ) as tb";
	}else{
		$tbQuery = "SELECT count(*) as tbCount from
						(SELECT max(tp.id) from tbl_poa tp
								inner join project p on p.id = tp.project_id
								where check_out is not null and `type` = 1 and added = 0 and startDate >= '$startDate' and endDate <= '$endDate'
						GROUP BY tb_id ) as tb ";
	}
	return mysqli_query($db,$tbQuery)->fetch_assoc()['tbCount'];*/

	for _, proj := range projectArray {
		//check if there are any associated project for each project
		projs, _ := GetAssociatedProjectList(proj, db)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = removeDuplicates(projectArray)

	//get the count of villages
	var villageQuery string
	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as villages from tbl_poa tp where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and project_id in (%s) %s", intSliceToString(projectArray), filter)
	} else {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as villages from tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and startDate >= '%s' and endDate <= '%s' %s", startDate, endDate, filter)
	}
	var villages int
	err := db.QueryRow(villageQuery).Scan(&villages)
	if err != nil {
		log.Println("getTrainingBatches", err)
	}
	return villages
}

func removeDuplicates(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func newVillageCount(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	// var projs []int
	for _, proj := range projectArray {
		// check if there are any associated project for each project
		projs, _ := GetAssociatedProjectList(proj, db)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// projs = unique(projs)

	var villageQuery, subVillageQuery string
	var villageCount, subVillageCount int

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village='' AND tb_id != primary_id AND type = 1 AND added = 0 AND project_id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]"))
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village!='' AND tb_id != primary_id AND type = 1 AND added = 0 AND project_id IN (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]"))

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
		}

		row := db.QueryRow(villageQuery)
		row.Scan(&villageCount)
		row = db.QueryRow(subVillageQuery)
		row.Scan(&subVillageCount)
	} else {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND tb_id != primary_id AND type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND tb_id != primary_id AND type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)

		row := db.QueryRow(villageQuery)
		row.Scan(&villageCount)
		row = db.QueryRow(subVillageQuery)
		row.Scan(&subVillageCount)
	}

	return villageCount + subVillageCount
}

func getSummaryOfVillagesNew(con *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	if len(projectArray) > 0 {
		villageQuery := "SELECT COUNT(DISTINCT location_id) as `village` FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0 AND project_id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
		subVillageQuery := "SELECT COUNT(DISTINCT sub_village) as `subVillage` FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0 AND project_id IN (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
		var village int
		var subVillage int
		err := con.QueryRow(villageQuery).Scan(&village)
		if err != nil {
			log.Println("getSummaryOfVillagesNew", err)
		}
		err = con.QueryRow(subVillageQuery).Scan(&subVillage)
		if err != nil {
			log.Println("getSummaryOfVillagesNew", err)
		}
		return village + subVillage
	} else {
		var dateFilter string
		if len(startDate) > 0 {
			dateFilter = " and date >= '" + startDate + "' and date <= '" + endDate + "'"
		} else {
			dateFilter = " and date >= CURRENT_DATE()"
		}
		villageQuery := "SELECT COUNT(DISTINCT location_id) as `village` FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0" + filter + dateFilter
		subVillageQuery := "SELECT COUNT(DISTINCT sub_village) as `subVillage` FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0" + filter + dateFilter
		var village int
		var subVillage int
		err := con.QueryRow(villageQuery).Scan(&village)
		if err != nil {
			log.Println("getSummaryOfVillagesNew", err)
		}
		err = con.QueryRow(subVillageQuery).Scan(&subVillage)
		if err != nil {
			log.Println("getSummaryOfVillagesNew", err)
		}
		return village + subVillage
	}
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
		log.Println("getOpsManagers", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			log.Println("getOpsManagers", err)
		}
		ids = append(ids, id)
	}
	return ids
}

func getSupervisor(db *sql.DB, empId int) []int {
	getOpsIds := fmt.Sprintf("SELECT supervisorId FROM employee WHERE id = %d", empId)
	ids := make([]int, 0)
	res, err := db.Query(getOpsIds)
	if err != nil {
		log.Println("getSupervisor", err)
	}
	defer res.Close()

	for res.Next() {
		var supervisorId int
		if err := res.Scan(&supervisorId); err != nil {
			log.Println("getSupervisor", err)
		}
		ids = append(ids, supervisorId)
	}
	return ids
}

func getReportingOpsManagers(db *sql.DB, empId int) []int {
	ids := []int{}

	getOpsIds := "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
	res, err := db.Query(getOpsIds, empId)
	if err != nil {
		log.Println("getReportingOpsManagers", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Println("getReportingOpsManagers", err)
		}
		ids = append(ids, id)
	}

	getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 12"
	res, err = db.Query(getOpsIds, empId)
	if err != nil {
		log.Println("getReportingOpsManagers", err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Println("getReportingOpsManagers", err)
		}
		ids = append(ids, id)
		som := id
		getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
		res, err = db.Query(getOpsIds, som)
		if err != nil {
			log.Println("getReportingOpsManagers", err)
		}
		defer res.Close()

		for res.Next() {
			var id int
			err := res.Scan(&id)
			if err != nil {
				log.Println("getReportingOpsManagers", err)
			}
			ids = append(ids, id)
		}
	}

	return ids
}

func getOpProjects(db *sql.DB, empID int) []int {
	// var projectId int
	// getProjIds := fmt.Sprintf("SELECT DISTINCT tr.project_id, prj.endDate FROM bdms_staff.project prj INNER JOIN bdms_staff.training_participants tr ON tr.project_id = prj.id WHERE prj.operations_manager = %d AND tr.VyaparEnrollment = 1", empID)

	// getProjIds := fmt.Sprintf("SELECT distinct(tr.project_id) FROM bdms_staff.project prj,bdms_staff.training_participants tr WHERE prj.operations_manager = %d and tr.VyaparEnrollment=1 and tr.project_id=prj.id", empID)

	getProjIds := fmt.Sprintf("SELECT id FROM project WHERE operations_manager = %d GROUP BY id", empID)
	rows, err := db.Query(getProjIds)
	if err != nil {
		log.Println("getOpProjects", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println("getOpProjects", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Println("getOpProjects", err)
	}

	return ids
}

func getTrainerTarget(db *sql.DB, empId int, projectArray []int) int {
	targetQuery := fmt.Sprintf("SELECT sum(target) as total from project_emps pe where emp_id = %d and project_id in (%s)",
		empId, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	var total sql.NullInt64
	err := db.QueryRow(targetQuery).Scan(&total)
	if err != nil {
		log.Println("getTrainerTarget", err)
	}
	if !total.Valid {
		return 0
	}
	return int(total.Int64)
}

func getGFData(db *sql.DB, filter string, sessionType int, empId int) int {
	filter += fmt.Sprintf(" and tp.user_id = %d", empId)
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
	var visit int
	row := db.QueryRow(getVisit)
	err := row.Scan(&visit)
	if err != nil {
		log.Println("getGFData", err)
	}
	return visit
}

// func getGFDataN(db *sql.DB, filter string, sessionType int, empId []int) int {
// 	filter += fmt.Sprintf(" and tp.user_id in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(empId)), ","), "[]"))
// 	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
// 	var visit int
// 	row := db.QueryRow(getVisit)
// 	err := row.Scan(&visit)
// 	if err != nil {
// 		log.Println("getGFDataN", err)
// 	}
// 	return visit
// }

// func getGFCircle(db *sql.DB, filter string, empId int) int {
// 	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id = %d%s", empId, filter)
// 	row := db.QueryRow(getCircle)
// 	var visit int
// 	row.Scan(&visit)
// 	return visit
// }

// func getGFCircleN(db *sql.DB, filter string, empId []int) int {
// 	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id in (%s)%s", strings.Trim(strings.Replace(fmt.Sprint(empId), " ", ",", -1), "[]"), filter)
// 	row := db.QueryRow(getCircle)
// 	var visit int
// 	row.Scan(&visit)
// 	return visit
// }

// func getGfEnrolled(db *sql.DB, filter string, empID int) (int, error) {
// 	query := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled FROM training_participants tp "+
// 		"INNER JOIN project p ON tp.project_id = p.id "+
// 		"WHERE enroll = 1 AND gelathi_id = %d %s", empID, filter)

// 	row := db.QueryRow(query)
// 	var enrolled int
// 	err := row.Scan(&enrolled)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return enrolled, nil
// }

// func getGfEnrolledN(db *sql.DB, filter string, empId []int) int {
// 	empIdStr := make([]string, len(empId))
// 	for i, id := range empId {
// 		empIdStr[i] = strconv.Itoa(id)
// 	}
// 	query := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
// 		"inner join project p on tp.project_id = p.id " +
// 		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " + filter
// 	row := db.QueryRow(query)
// 	var enrolled int
// 	row.Scan(&enrolled)
// 	return enrolled
// }

// func getParticipantFilterGfEnrolled(db *sql.DB, filter string, empId int, startDate string, endDate string) int {
// 	getEnrolled := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled from training_participants tp "+
// 		"inner join project p on tp.project_id = p.id "+
// 		"where enroll = 1 and gelathi_id = %d and "+
// 		"((participant_day1 >= '%s' and participant_day1 <= '%s') "+
// 		"or (participant_day2 >= '%s' and participant_day2 <= '%s')) %s", empId, startDate, endDate, startDate, endDate, filter)
// 	row := db.QueryRow(getEnrolled)
// 	var enrolled int
// 	err := row.Scan(&enrolled)
// 	if err != nil {
// 		log.Println("getParticipantFilterGfEnrolled", err)
// 	}
// 	return enrolled
// }

// func getParticipantFilterGfEnrolledN(db *sql.DB, filter string, empId []int, startDate, endDate time.Time) int {
// 	enrolled := 0
// 	empIdStr := make([]string, len(empId))
// 	for i, v := range empId {
// 		empIdStr[i] = strconv.Itoa(v)
// 	}
// 	getEnrolled := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
// 		"inner join project p on tp.project_id = p.id " +
// 		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " +
// 		"and ((participant_day1 >= ? and participant_day1 <= ?) or (participant_day2 >= ? and participant_day2 <= ?)) " + filter
// 	row := db.QueryRow(getEnrolled, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
// 	row.Scan(&enrolled)
// 	return enrolled
// }

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
		log.Println("NoofBatches", err)
	}
	return noofbatches
}

func GetAssociatedProjectList(projId int, DB *sql.DB) ([]int, error) {
	//projId := (pid).(int)
	// projId := int(proId)
	var projArray []int
	if projId > 0 {
		getProjList := fmt.Sprintf("SELECT associatedProject FROM project_association WHERE projectId IN (%d)", projId)
		projArray = append(projArray, projId)
		rows, err := DB.Query(getProjList)
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

func getTrainerActual(db *sql.DB, startDate string, endDate string, projectArray []int, trainer_id string) int {
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

func getTrainerBatches(db *sql.DB, startDate string, endDate string, projectArray []int, trainer_id string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")
		getActualsQuery = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa WHERE user_id = %s and project_id = %s", trainer_id, projectIDs)
		if startDate != "" && endDate != "" {
			getActualsQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa WHERE user_id IN (%s)", trainer_id)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Println("getActual", err)
	}
	return actual
}

func noofsurveyssTrainer(db *sql.DB, startDate string, endDate string, projectArray []int, trainer_id string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")
		getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofselfsakthi FROM training_participants WHERE isSurveyDone=1 and trainer_id = %s and project_id = %s", trainer_id, projectIDs)
		if startDate != "" && endDate != "" {
			getActualsQuery += fmt.Sprintf(" AND participant_day2 BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofselfsakthi FROM training_participants WHERE isSurveyDone=1 and trainer_id = %s", trainer_id)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Println("getActual", err)
	}
	return actual
}
