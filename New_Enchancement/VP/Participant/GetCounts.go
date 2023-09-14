package vyapar

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
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
		log.Fatal(err)
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
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	if len(projectArray) > 0 {
		// Build the project IDs string for the SQL query
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")

		villageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE type = 1 AND added = 0 AND project_id IN (%s) AND user_id = %s", projectIDs, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE type = 1 AND added = 0 AND project_id IN (%s) AND user_id = %s", projectIDs, user_id)

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND date BETWEEN '%s' AND '%s'", startDate, endDate)
		}
		fmt.Println(villageQuery, subVillageQuery)

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

// func removeDuplicates(intSlice []int) []int {
// 	keys := make(map[int]bool)
// 	list := []int{}
// 	for _, entry := range intSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }

func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, val := range arr {
		strArr[i] = fmt.Sprint(val)
	}
	return strings.Join(strArr, ", ")
}

func GfVyaparVillageCount(db *sql.DB, startDate string, endDate string, projectArray []int, user_id string) int {
	var villageCount, subVillageCount int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		// Build the project IDs string for the SQL query
		projectIDs := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]")

		villageQuery := fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and t.project_id IN (%s) AND tp.user_id = %s", projectIDs, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) t.project_id IN (%s) AND tp.user_id = %s", projectIDs, user_id)

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND tp.date BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND tp.date BETWEEN '%s' AND '%s'", startDate, endDate)
		}

		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
		}
	} else {
		villageQuery := fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' AND tp.user_id = %s", startDate, endDate, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' AND tp.user_id = %s", startDate, endDate, user_id)
		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
		}
	}

	return villageCount + subVillageCount
}

func getParticipantFilterTrainingBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, gfid int) int {
	if gfid > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(gfid)
	}
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
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

func getParticipantFilterTrainingGFLBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, gfid int) int {
	if gfid > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(gfid)
	}
	gfl := ""

	gfl += "and p.gfl_id = " + filter
	var villageQuery, subVillageQuery string
	if len(projectArray) > 0 {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp join project p on p.id=tp.project_id where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + gfl + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join project p on p.id=tp.project_id where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + gfl + filter
	} else {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + gfl + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + gfl + filter
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
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 12"
	res, err = db.Query(getOpsIds, empId)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
		som := id
		getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
		res, err = db.Query(getOpsIds, som)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Close()

		for res.Next() {
			var id int
			err := res.Scan(&id)
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
	}
	if !total.Valid {
		return 0
	}
	return int(total.Int64)
}

func getParticipantFilternoVyaparsurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {

	var vyaparsurvey string
	var noofvyaparsurvey int
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between '%s' and '%s' and t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	} else {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(vyaparsurvey).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Println("noofvyaparsurvey", err)
	}
	return noofvyaparsurvey

	// } else {
	// 	// Create a map to track unique project IDs and a set for unique associated project IDs
	// 	uniqueProjectIDs := make(map[int]bool)
	// 	uniqueAssociatedProjectIDs := make(map[int]bool)

	// 	for _, proj := range projectArray {
	// 		// Check if the project is already in the map
	// 		if _, exists := uniqueProjectIDs[proj]; !exists {
	// 			uniqueProjectIDs[proj] = true // Store the project in the map

	// 			// Fetch associated project IDs
	// 			projs, _ := getAssociatedProjectList(db, proj)

	// 			for _, p := range projs {
	// 				// Only add the associative project if it doesn't already exist in the set
	// 				if _, exists := uniqueProjectIDs[p]; !exists {
	// 					uniqueAssociatedProjectIDs[p] = true
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Convert the map keys (unique project IDs) into a slice
	// 	var allProjectIDs []int
	// 	for projID := range uniqueProjectIDs {
	// 		allProjectIDs = append(allProjectIDs, projID)
	// 	}

	// 	// Convert the set keys (unique associated project IDs) into a slice
	// 	var allAssociatedProjectIDs []int
	// 	for assocProjID := range uniqueAssociatedProjectIDs {
	// 		allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
	// 	}

	// 	// Convert the project IDs slice into a comma-separated string
	// 	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
	// 	// Convert the associated project IDs slice into a comma-separated string
	// 	// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

	// 	fmt.Println(projectIDs)
	// 	// fmt.Println(associatedProjectIDs)

	// 	if len(projectIDs) > 0 {
	// 		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between '%s' and '%s' and t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	// 	} else {
	// 		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	// 	}

	// 	err := db.QueryRow(vyaparsurvey).Scan(&noofvyaparsurvey)
	// 	if err != nil {
	// 		log.Println("noofvyaparsurvey", err)
	// 	}
	// 	return noofvyaparsurvey

	// }
}

func getEmployeeIDsBySupervisorID(db *sql.DB, supervisorID int) ([]int, error) {
	var employeeIDs []int

	query := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d", supervisorID)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employeeID int
		if err := rows.Scan(&employeeID); err != nil {
			return nil, err
		}
		employeeIDs = append(employeeIDs, employeeID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employeeIDs, nil
}
func getParticipantFilternoVyaparenroll(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparenroll string
	var noofvyaparenroll int

	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	} else {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
	if err != nil {
		log.Println("noofvyaparenroll", err)
	}
	return noofvyaparenroll

	// } else {
	// 	// Create a map to track unique project IDs and a set for unique associated project IDs
	// 	uniqueProjectIDs := make(map[int]bool)
	// 	uniqueAssociatedProjectIDs := make(map[int]bool)

	// 	for _, proj := range projectArray {
	// 		// Check if the project is already in the map
	// 		if _, exists := uniqueProjectIDs[proj]; !exists {
	// 			uniqueProjectIDs[proj] = true // Store the project in the map

	// 			// Fetch associated project IDs
	// 			projs, _ := getAssociatedProjectList(db, proj)

	// 			for _, p := range projs {
	// 				// Only add the associative project if it doesn't already exist in the set
	// 				if _, exists := uniqueProjectIDs[p]; !exists {
	// 					uniqueAssociatedProjectIDs[p] = true
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Convert the map keys (unique project IDs) into a slice
	// 	var allProjectIDs []int
	// 	for projID := range uniqueProjectIDs {
	// 		allProjectIDs = append(allProjectIDs, projID)
	// 	}

	// 	// Convert the set keys (unique associated project IDs) into a slice
	// 	var allAssociatedProjectIDs []int
	// 	for assocProjID := range uniqueAssociatedProjectIDs {
	// 		allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
	// 	}

	// 	// Convert the project IDs slice into a comma-separated string
	// 	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
	// 	// Convert the associated project IDs slice into a comma-separated string
	// 	// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

	// 	fmt.Println(projectIDs)
	// 	// fmt.Println(associatedProjectIDs)

	// 	if len(projectIDs) > 0 {
	// 		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	// 	} else {
	// 		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	// 	}

	// 	err := db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
	// 	if err != nil {
	// 		log.Println("noofvyaparenroll", err)
	// 	}
	// 	return noofvyaparenroll

	// }
}

// Function to remove duplicates from a slice of integers
func removeDuplicates(arr []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range arr {
		if encountered[arr[v]] == true {
			// Do not add duplicate to result
			continue
		} else {
			// Add unique element to result and mark it as encountered
			encountered[arr[v]] = true
			result = append(result, arr[v])
		}
	}

	return result
}

func getParticipantFilternoVyaparCoharts(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparcoharts string
	var noofvyaparcoharts int
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) WHERE and tbl_poa.date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparcoharts", err)
	}
	return noofvyaparcoharts

	// } else {
	// 	// Create a map to track unique project IDs and a set for unique associated project IDs
	// 	uniqueProjectIDs := make(map[int]bool)
	// 	uniqueAssociatedProjectIDs := make(map[int]bool)

	// 	for _, proj := range projectArray {
	// 		// Check if the project is already in the map
	// 		if _, exists := uniqueProjectIDs[proj]; !exists {
	// 			uniqueProjectIDs[proj] = true // Store the project in the map

	// 			// Fetch associated project IDs
	// 			projs, _ := getAssociatedProjectList(db, proj)

	// 			for _, p := range projs {
	// 				// Only add the associative project if it doesn't already exist in the set
	// 				if _, exists := uniqueProjectIDs[p]; !exists {
	// 					uniqueAssociatedProjectIDs[p] = true
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Convert the map keys (unique project IDs) into a slice
	// 	var allProjectIDs []int
	// 	for projID := range uniqueProjectIDs {
	// 		allProjectIDs = append(allProjectIDs, projID)
	// 	}

	// 	// Convert the set keys (unique associated project IDs) into a slice
	// 	var allAssociatedProjectIDs []int
	// 	for assocProjID := range uniqueAssociatedProjectIDs {
	// 		allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
	// 	}

	// 	// Convert the project IDs slice into a comma-separated string
	// 	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
	// 	// Convert the associated project IDs slice into a comma-separated string
	// 	// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

	// 	fmt.Println(projectIDs)
	// 	// fmt.Println(associatedProjectIDs)

	// 	if len(projectIDs) > 0 {
	// 		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	// 	} else {
	// 		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) WHERE and tbl_poa.date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	// 	}

	// 	err := db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	// 	if err != nil {
	// 		log.Println("noofvyaparcoharts", err)
	// 	}
	// 	return noofvyaparcoharts

	// }
}

func getParticipantFilterVyaprModule(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparmodule string
	var noofvyaparmodule int
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	} else {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(vyaparmodule).Scan(&noofvyaparmodule)
	if err != nil {
		log.Println("noofvyaparmodule", err)
	}
	return noofvyaparmodule

	// } else {
	// 	// Create a map to track unique project IDs and a set for unique associated project IDs
	// 	uniqueProjectIDs := make(map[int]bool)
	// 	uniqueAssociatedProjectIDs := make(map[int]bool)

	// 	for _, proj := range projectArray {
	// 		// Check if the project is already in the map
	// 		if _, exists := uniqueProjectIDs[proj]; !exists {
	// 			uniqueProjectIDs[proj] = true // Store the project in the map

	// 			// Fetch associated project IDs
	// 			projs, _ := getAssociatedProjectList(db, proj)

	// 			for _, p := range projs {
	// 				// Only add the associative project if it doesn't already exist in the set
	// 				if _, exists := uniqueProjectIDs[p]; !exists {
	// 					uniqueAssociatedProjectIDs[p] = true
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Convert the map keys (unique project IDs) into a slice
	// 	var allProjectIDs []int
	// 	for projID := range uniqueProjectIDs {
	// 		allProjectIDs = append(allProjectIDs, projID)
	// 	}

	// 	// Convert the set keys (unique associated project IDs) into a slice
	// 	var allAssociatedProjectIDs []int
	// 	for assocProjID := range uniqueAssociatedProjectIDs {
	// 		allAssociatedProjectIDs = append(allAssociatedProjectIDs, assocProjID)
	// 	}

	// 	// Convert the project IDs slice into a comma-separated string
	// 	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(allProjectIDs), " ", ",", -1), "[]")
	// 	// Convert the associated project IDs slice into a comma-separated string
	// 	// associatedProjectIDs := strings.Trim(strings.Replace(fmt.Sprint(allAssociatedProjectIDs), " ", ",", -1), "[]")

	// 	fmt.Println(projectIDs)
	// 	// fmt.Println(associatedProjectIDs)

	// 	if len(projectIDs) > 0 {
	// 		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), filter)
	// 	} else {
	// 		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	// 	}

	// 	err := db.QueryRow(vyaparmodule).Scan(&noofvyaparmodule)
	// 	if err != nil {
	// 		log.Println("noofvyaparmodule", err)
	// 	}
	// 	return noofvyaparmodule

	// }
}

func getParticipantFilterVyaparGfBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, gfid int) int {

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

type CacheEntry struct {
	Result     int
	Expiration time.Time
}

var cache = make(map[string]CacheEntry)
var cacheLock sync.RWMutex

func getParticipantFilternoVyaparenrollCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "getParticipantFilternoVyaparenroll")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := getParticipantFilternoVyaparenroll(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}
func getParticipantFilternoVyaparsurveyCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "getParticipantFilternoVyaparsurvey")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := getParticipantFilternoVyaparsurvey(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}

func getParticipantFilterVyaprModuleCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "getParticipantFilterVyaprModule")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := getParticipantFilterVyaprModule(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}
func getParticipantFilternoVyaparCohartsCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, role int) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "getParticipantFilternoVyaparCoharts")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := getParticipantFilternoVyaparCoharts(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}
func getParticipantFilterVyaparGfBatchesNewCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, gfid int) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%d-%s", startDate, endDate, fmt.Sprint(projectArray), filter, gfid, "getParticipantFilterVyaparGfBatchesNew")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := getParticipantFilterVyaparGfBatchesNew(db, startDate, endDate, projectArray, filter, gfid)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}

func getParticipantFilternoGflVyaparenroll(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparenroll string
	var noofvyaparenroll int

	gfid, err := getEmployeeIDsBySupervisorID(db, gflid)
	if err != nil {
		log.Println("Error retrieving employee IDs:", err)
		return 0
	}

	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' AND t.project_id IN (%s) and t.gelathi_id IN (%s)  %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
	} else {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 OR t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	}

	err = db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
	if err != nil {
		log.Println("noofvyaparenroll", err)
	}
	return noofvyaparenroll

}

func getParticipantFilternoGflVyaparsurvey(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {

	var vyaparsurvey string
	var noofvyaparsurvey int
	gfid, err := getEmployeeIDsBySupervisorID(db, gflid)
	if err != nil {
		log.Println("Error retrieving employee IDs:", err)
		return 0
	}
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between '%s' and '%s' and t.project_id IN (%s) and b.gfId IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
	} else {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err = db.QueryRow(vyaparsurvey).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Println("noofvyaparsurvey", err)
	}
	return noofvyaparsurvey
}
func getParticipantFilterGflVyaprModule(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparmodule string
	var noofvyaparmodule int
	gfid, err := getEmployeeIDsBySupervisorID(db, gflid)
	if err != nil {
		log.Println("Error retrieving employee IDs:", err)
		return 0
	}
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) and b.gfId IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
	} else {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err = db.QueryRow(vyaparmodule).Scan(&noofvyaparmodule)
	if err != nil {
		log.Println("noofvyaparmodule", err)
	}
	return noofvyaparmodule
}
func getParticipantFilternoGflVyaparCoharts(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparcoharts string
	var noofvyaparcoharts int
	gfid, err := getEmployeeIDsBySupervisorID(db, gflid)
	if err != nil {
		log.Println("Error retrieving employee IDs:", err)
		return 0
	}
	// if role == 4 {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectIDs := strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]")
	if len(projectIDs) > 0 {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) and tbl_poa.user_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectIDs), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) WHERE and tbl_poa.date BETWEEN '%s' and '%s' %s", startDate, endDate, filter)
	}

	err = db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparcoharts", err)
	}
	return noofvyaparcoharts
}
func GflVyaparnewVillageCount(db *sql.DB, startDate, endDate string, projectArray []int, gflid int, filter string) int {
	gfid, err := getEmployeeIDsBySupervisorID(db, gflid)
	if err != nil {
		log.Println("Error retrieving employee IDs:", err)
		return 0
	}

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	projectArray = uniqueIntSlice(projectArray)

	var villageQuery, subVillageQuery string

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date BETWEEN '%s' AND '%s' and t.project_id IN (%s) and tp.user_id IN (%s) %s", startDate, endDate, intArrayToString(projectArray), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date BETWEEN '%s' AND '%s' and t.project_id IN (%s) and tp.user_id IN (%s) %s", startDate, endDate, intArrayToString(projectArray), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)

	} else {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) and tp.date BETWEEN '%s' AND '%s' %s", startDate, endDate, filter)
	}

	villageResult, err := db.Query(villageQuery)
	if err != nil {
		return 0
	}
	defer villageResult.Close()

	var villageCount int
	for villageResult.Next() {
		if err := villageResult.Scan(&villageCount); err != nil {
			return 0
		}
	}

	subVillageResult, err := db.Query(subVillageQuery)
	if err != nil {
		return 0
	}
	defer subVillageResult.Close()

	var subVillageCount int
	for subVillageResult.Next() {
		if err := subVillageResult.Scan(&subVillageCount); err != nil {
			return 0
		}
	}

	return villageCount + subVillageCount
}
