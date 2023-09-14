package green

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

func getParticipantFilternoVyaparenroll(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	var vyaparenroll string
	var noofvyaparenroll int
	if role == 4 {
		for _, proj := range projectArray {
			projs, _ := getAssociatedProjectList(db, proj)
			if len(projs) > 1 {
				projectArray = append(projectArray, projs...)
			}
		}
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
		if len(projectIDs) > 0 {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS green FROM training_participants t WHERE (t.GreenMotivators=1 or t.new_green = 1) AND t.GreenMotivatorsDate BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
		} else {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS green FROM training_participants t WHERE (t.GreenMotivators=1 or t.new_green = 1) AND t.GreenMotivatorsDate BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
		if err != nil {
			log.Println("noofgreenenroll", err)
		}
		return noofvyaparenroll

	} else {
		// Create a map to track unique project IDs and a set for unique associated project IDs
		uniqueProjectIDs := make(map[int]bool)
		uniqueAssociatedProjectIDs := make(map[int]bool)

		for _, proj := range projectArray {
			// Check if the project is already in the map
			if _, exists := uniqueProjectIDs[proj]; !exists {
				uniqueProjectIDs[proj] = true // Store the project in the map

				// Fetch associated project IDs
				projs, _ := getAssociatedProjectList(db, proj)

				for _, p := range projs {
					// Only add the associative project if it doesn't already exist in the set
					if _, exists := uniqueProjectIDs[p]; !exists {
						uniqueAssociatedProjectIDs[p] = true
					}
				}
			}
		}

		// Convert the map keys (unique project IDs) into a slice
		var allProjectIDs []int
		for projID := range uniqueProjectIDs {
			allProjectIDs = append(allProjectIDs, projID)
		}

		// Convert the set keys (unique associated project IDs) into a slice
		var allAssociatedProjectIDs []int
		for assocProjID := range uniqueAssociatedProjectIDs {
			allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
		}

		// Convert the project IDs slice into a comma-separated string
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
		// Convert the associated project IDs slice into a comma-separated string
		// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

		fmt.Println(projectIDs)
		// fmt.Println(associatedProjectIDs)

		if len(projectIDs) > 0 {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS green FROM training_participants t WHERE (t.GreenMotivators=1 or t.new_green = 1) AND t.GreenMotivatorsDate BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
		} else {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS green FROM training_participants t WHERE (t.GreenMotivators=1 or t.new_green = 1) AND t.GreenMotivatorsDate BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
		if err != nil {
			log.Println("noofgreenenroll", err)
		}
		return noofvyaparenroll

	}
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
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

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

func getParticipantFilternoGreensurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	var greensurvey string
	var noofgreensurevy int

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")

	if len(projectArray) > 0 {
		// Use the AND operator to separate different conditions in the WHERE clause.
		greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId where g.entry_date BETWEEN '%s' and '%s' and t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	} else {
		// Use the AND operator to separate different conditions in the WHERE clause.
		greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId where g.entry_date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(greensurvey).Scan(&noofgreensurevy)
	if err != nil {
		log.Println("greensurvey", err)
	}
	return noofgreensurevy

}

func getParticipantFilternoGreenCoharts(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	var greencoharts string
	var noofgreencoharts int
	if role == 4 {
		for _, proj := range projectArray {
			projs, _ := getAssociatedProjectList(db, proj)
			if len(projs) > 1 {
				projectArray = append(projectArray, projs...)
			}
		}
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")

		if len(projectArray) > 0 {

			greencoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa INNER JOIN project ON project.id = tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) AND tbl_poa.date between '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)

		} else {
			greencoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa INNER JOIN project ON project.id = tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) AND tbl_poa.date between '%s' and '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(greencoharts).Scan(&noofgreencoharts)
		if err != nil {
			log.Println("greencoharts", err)
		}
		return noofgreencoharts
	} else {
		uniqueProjectIDs := make(map[int]bool)
		uniqueAssociatedProjectIDs := make(map[int]bool)

		for _, proj := range projectArray {
			// Check if the project is already in the map
			if _, exists := uniqueProjectIDs[proj]; !exists {
				uniqueProjectIDs[proj] = true // Store the project in the map

				// Fetch associated project IDs
				projs, _ := getAssociatedProjectList(db, proj)

				for _, p := range projs {
					// Only add the associative project if it doesn't already exist in the set
					if _, exists := uniqueProjectIDs[p]; !exists {
						uniqueAssociatedProjectIDs[p] = true
					}
				}
			}
		}

		// Convert the map keys (unique project IDs) into a slice
		var allProjectIDs []int
		for projID := range uniqueProjectIDs {
			allProjectIDs = append(allProjectIDs, projID)
		}

		// Convert the set keys (unique associated project IDs) into a slice
		var allAssociatedProjectIDs []int
		for assocProjID := range uniqueAssociatedProjectIDs {
			allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
		}

		// Convert the project IDs slice into a comma-separated string
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
		// Convert the associated project IDs slice into a comma-separated string
		// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

		fmt.Println(projectIDs)
		// fmt.Println(associatedProjectIDs)

		if len(projectArray) > 0 {

			greencoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa INNER JOIN project ON project.id = tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) AND tbl_poa.date between '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)

		} else {
			greencoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa INNER JOIN project ON project.id = tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 10 OR tbl_poa.session_type = 11 OR tbl_poa.session_type = 12 OR tbl_poa.session_type = 13 OR tbl_poa.session_type = 14 OR tbl_poa.session_type = 15) AND tbl_poa.date between '%s' and '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(greencoharts).Scan(&noofgreencoharts)
		if err != nil {
			log.Println("greencoharts", err)
		}
		return noofgreencoharts
	}

}

func getParticipantFilterGreenModule(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	var greenmodule string
	var noofgreenmodule int
	if role == 4 {
		for _, proj := range projectArray {
			projs, _ := getAssociatedProjectList(db, proj)
			if len(projs) > 1 {
				projectArray = append(projectArray, projs...)
			}
		}
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")

		if len(projectArray) > 0 {
			// Use the AND operator to separate different conditions in the WHERE clause.
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId WHERE g.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
		} else {
			// Use the AND operator to separate different conditions in the WHERE clause.
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId WHERE g.entry_date BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(greenmodule).Scan(&noofgreenmodule)
		if err != nil {
			log.Println("greenmodule", err)
		}
		return noofgreenmodule
	} else {
		// Create a map to track unique project IDs and a set for unique associated project IDs
		uniqueProjectIDs := make(map[int]bool)
		uniqueAssociatedProjectIDs := make(map[int]bool)

		for _, proj := range projectArray {
			// Check if the project is already in the map
			if _, exists := uniqueProjectIDs[proj]; !exists {
				uniqueProjectIDs[proj] = true // Store the project in the map

				// Fetch associated project IDs
				projs, _ := getAssociatedProjectList(db, proj)

				for _, p := range projs {
					// Only add the associative project if it doesn't already exist in the set
					if _, exists := uniqueProjectIDs[p]; !exists {
						uniqueAssociatedProjectIDs[p] = true
					}
				}
			}
		}

		// Convert the map keys (unique project IDs) into a slice
		var allProjectIDs []int
		for projID := range uniqueProjectIDs {
			allProjectIDs = append(allProjectIDs, projID)
		}

		// Convert the set keys (unique associated project IDs) into a slice
		var allAssociatedProjectIDs []int
		for assocProjID := range uniqueAssociatedProjectIDs {
			allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
		}

		// Convert the project IDs slice into a comma-separated string
		projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
		// Convert the associated project IDs slice into a comma-separated string
		// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

		fmt.Println(projectIDs)
		// fmt.Println(associatedProjectIDs)

		if len(projectArray) > 0 {
			// Use the AND operator to separate different conditions in the WHERE clause.
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId WHERE g.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
		} else {
			// Use the AND operator to separate different conditions in the WHERE clause.
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId WHERE g.entry_date BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
		}

		err := db.QueryRow(greenmodule).Scan(&noofgreenmodule)
		if err != nil {
			log.Println("greenmodule", err)
		}
		return noofgreenmodule
	}
}

func getParticipantFilterGreenGfBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, gfid int) int {

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	if gfid > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(gfid)
	}
	var villageQuery, subVillageQuery string
	if len(projectArray) > 0 {
		villageQuery = "select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND tp.project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND tp.project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
	} else {
		villageQuery = "select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or new_vyapar=1) and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
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
