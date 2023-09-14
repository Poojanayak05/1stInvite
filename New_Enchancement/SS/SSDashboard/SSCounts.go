package New

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func getParticipantFilterGfEnrolled(db *sql.DB, filter string, empId int, startDate string, endDate string) int {
	getEnrolled := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled from training_participants tp "+
		"inner join project p on tp.project_id = p.id "+
		"where enroll = 1 and gelathi_id = %d and "+
		"((participant_day1 >= '%s' and participant_day1 <= '%s') "+
		"or (participant_day2 >= '%s' and participant_day2 <= '%s')) %s", empId, startDate, endDate, startDate, endDate, filter)
	row := db.QueryRow(getEnrolled)
	var enrolled int
	err := row.Scan(&enrolled)
	if err != nil {
		log.Println("getParticipantFilterGfEnrolled", err)
	}
	return enrolled
}
func getGFCircle(db *sql.DB, filter string, empId int) int {
	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id = %d%s", empId, filter)
	row := db.QueryRow(getCircle)
	var visit int
	row.Scan(&visit)
	return visit
}
func getGfEnrolled(db *sql.DB, filter string, empID int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled FROM training_participants tp "+
		"INNER JOIN project p ON tp.project_id = p.id "+
		"WHERE enroll = 1 AND gelathi_id = %d %s", empID, filter)

	row := db.QueryRow(query)
	var enrolled int
	err := row.Scan(&enrolled)
	if err != nil {
		return 0, err
	}
	return enrolled, nil
}
func getParticipantFilterGfEnrolledN(con *sql.DB, filter string, empIDs []int, startDate, endDate string) int {
	empIDsStr := make([]string, len(empIDs))
	for i, id := range empIDs {
		empIDsStr[i] = strconv.Itoa(id)
	}

	empIDsJoined := strings.Join(empIDsStr, ",")

	getEnrolled := fmt.Sprintf(`SELECT COUNT(tp.id) as enrolled FROM training_participants tp
		INNER JOIN project p ON tp.project_id = p.id
		WHERE enroll = 1 AND gelathi_id IN (%s) AND ((participant_day1 >= '%s' AND participant_day1 <= '%s') OR (participant_day2 >= '%s' AND participant_day2 <= '%s')) %s`,
		empIDsJoined, startDate, endDate, startDate, endDate, filter)

	row := con.QueryRow(getEnrolled)

	var enrolled int
	err := row.Scan(&enrolled)
	if err != nil {
		log.Println("ERROR>>", err)
	}

	return enrolled
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

func getGFCircleN(db *sql.DB, filter string, empId []int) int {
	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id in (%s)%s", strings.Trim(strings.Replace(fmt.Sprint(empId), " ", ",", -1), "[]"), filter)
	row := db.QueryRow(getCircle)
	var visit int
	row.Scan(&visit)
	return visit
}
func getGFDataN(db *sql.DB, filter string, sessionType int, empId []int) int {
	filter += fmt.Sprintf(" and tp.user_id in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(empId)), ","), "[]"))
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
	var visit int
	row := db.QueryRow(getVisit)
	err := row.Scan(&visit)
	if err != nil {
		log.Println("getGFDataN", err)
	}
	return visit
}
func getTarget(con *sql.DB, startDate, endDate string, projectArray []int) (int, error) {
	var getTargetQuery string
	if len(projectArray) > 0 {
		// getTargetQuery = fmt.Sprintf("SELECT sum(training_target) as target from project p where (startDate and endDate BETWEEN '%s' and '%s') and id in (%s)", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
		getTargetQuery = fmt.Sprintf("SELECT sum(training_target) as target from project p where id in (%s)", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	} else {
		getTargetQuery = fmt.Sprintf("SELECT sum(training_target) as target from project p where (startDate >= '%s' and endDate <= '%s')", startDate, endDate)
	}

	// Execute the query
	rows, err := con.Query(getTargetQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Fetch the result
	var target int
	for rows.Next() {
		if err := rows.Scan(&target); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return target, nil
}
func noOfselfshakthisurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var gelatiCountQuery string
	var noofselfshakthi int

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp2.id) AS actual FROM training_participants tp2 JOIN project p ON p.id = tp2.project_id WHERE isSurveyDone = 1 AND participant_day1 BETWEEN '%s' AND '%s' AND tp2.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp2.id) AS actual FROM training_participants tp2 JOIN project p ON p.id = tp2.project_id WHERE isSurveyDone = 1 AND tp2.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp2.id) AS actual FROM training_participants tp2 JOIN project p ON p.id = tp2.project_id WHERE isSurveyDone = 1 AND participant_day1 BETWEEN '%s' AND '%s'", startDate, endDate)
	}

	err := db.QueryRow(gelatiCountQuery).Scan(&noofselfshakthi)
	if err != nil {
		log.Println("getActual", err)
	}
	return noofselfshakthi
}
func noofbatches(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getnoofbatches string
	var noofbatches int

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			getnoofbatches = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa tp JOIN project t3 ON tp.project_id = t3.id WHERE date BETWEEN '%s' AND '%s' AND project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			getnoofbatches = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa tp JOIN project t3 ON tp.project_id = t3.id WHERE project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		getnoofbatches = fmt.Sprintf("SELECT COUNT(DISTINCT tb_id) FROM tbl_poa tp JOIN project t3 ON tp.project_id = t3.id WHERE date BETWEEN '%s' AND '%s'", startDate, endDate)
	}

	err := db.QueryRow(getnoofbatches).Scan(&noofbatches)
	if err != nil {
		log.Println("getActual", err)
	}

	return noofbatches
}
func getActual(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	var getActualsQuery string

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			// getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and (participant_day2 BETWEEN '%s' and '%s') and project_id in (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and (participant_day2 BETWEEN '%s' and '%s') and project_id in (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			// getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		getActualsQuery = fmt.Sprintf("select count(tp.id) as actual from training_participants tp inner join project p on p.id = tp.project_id where day2 = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	// Execute the query
	rows, err := con.Query(getActualsQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Fetch the result
	var actual int
	for rows.Next() {
		if err := rows.Scan(&actual); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return actual, nil
}

func getDay1Count(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	var day1CountQuery string

	if len(projectArray) > 0 {
		day1CountQuery = fmt.Sprintf("select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	} else {
		day1CountQuery = fmt.Sprintf("select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	// Execute the query
	rows, err := con.Query(day1CountQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Fetch the result
	var day1Count int
	for rows.Next() {
		if err := rows.Scan(&day1Count); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return day1Count, nil
}

func getParticipantFilterDay1Count(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	var day1CountQuery string

	if len(projectArray) > 0 {
		day1CountQuery = fmt.Sprintf("select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (%s) and participant_day1 >= '%s' and participant_day1 <= '%s' %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), startDate, endDate, filter)
	} else {
		day1CountQuery = fmt.Sprintf("select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and participant_day1 >= '%s' and participant_day1 <= '%s' %s", startDate, endDate, filter)
	}

	// Execute the query
	rows, err := con.Query(day1CountQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Fetch the result
	var day1Count int
	for rows.Next() {
		if err := rows.Scan(&day1Count); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return day1Count, nil
}

// Helper function to convert an int array to a string
func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, val := range arr {
		strArr[i] = fmt.Sprint(val)
	}
	return strings.Join(strArr, ", ")
}

func uniqueIntSlice(slice []int) []int {
	keys := make(map[int]bool)
	unique := []int{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			unique = append(unique, entry)
		}
	}
	return unique
}

func Vyapar(con *sql.DB, startDate, endDate string, projectArray []int, funderID int, filter string) (int, error) {
	if funderID != 0 {
		getProj := fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %d", funderID)
		projResult, err := con.Query(getProj)
		if err != nil {
			return 0, err
		}
		defer projResult.Close()

		for projResult.Next() {
			var projectID int
			if err := projResult.Scan(&projectID); err != nil {
				return 0, err
			}
			projectArray = append(projectArray, projectID)
		}

		if err := projResult.Err(); err != nil {
			return 0, err
		}
	}

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(con, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	projectArray = uniqueIntSlice(projectArray)

	var gelatiCountQuery string

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN '%s' AND '%s'", intArrayToString(projectArray), startDate, endDate)
		} else {
			if funderID != 0 {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE funderID = %d AND endDate >= CURRENT_DATE()) AND (SELECT MAX(endDate) FROM project p WHERE funderID = %d AND endDate >= CURRENT_DATE())", intArrayToString(projectArray), funderID, funderID)
			} else {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE id IN (%s)) AND (SELECT MAX(endDate) FROM project p WHERE id IN (%s))", intArrayToString(projectArray), intArrayToString(projectArray), intArrayToString(projectArray))
			}
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) AS Vyapar FROM training_participants tp INNER JOIN project p ON p.id = tp.VyaparEnrollmentEnrolledProject WHERE VyaparEnrollment = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}

	rows, err := con.Query(gelatiCountQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var VyaparCount int
	for rows.Next() {
		if err := rows.Scan(&VyaparCount); err != nil {
			return 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return VyaparCount, nil
}
func newVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(con, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	projectArray = uniqueIntSlice(projectArray)

	var villageQuery, subVillageQuery string

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village='' AND tb_id != primary_id AND `type` = 1 AND added = 0 AND project_id IN (%s)", intArrayToString(projectArray))
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village !='' AND tb_id != primary_id AND `type` = 1 AND added = 0 AND project_id IN (%s)", intArrayToString(projectArray))

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND `date` BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND `date` BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND tb_id != primary_id AND `type` = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' %s", startDate, endDate, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village !='' AND tb_id != primary_id AND `type` = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' %s", startDate, endDate, filter)
	}

	villageResult, err := con.Query(villageQuery)
	if err != nil {
		return 0, err
	}
	defer villageResult.Close()

	var villageCount int
	for villageResult.Next() {
		if err := villageResult.Scan(&villageCount); err != nil {
			return 0, err
		}
	}

	subVillageResult, err := con.Query(subVillageQuery)
	if err != nil {
		return 0, err
	}
	defer subVillageResult.Close()

	var subVillageCount int
	for subVillageResult.Next() {
		if err := subVillageResult.Scan(&subVillageCount); err != nil {
			return 0, err
		}
	}

	return villageCount + subVillageCount, nil
}

func getOpsManagers(con *sql.DB, empID int) ([]int, error) {
	getOpsIDsQuery := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d AND empRole = 4", empID)
	opsIDsResult, err := con.Query(getOpsIDsQuery)
	if err != nil {
		return nil, err
	}
	defer opsIDsResult.Close()

	var opsIDs []int
	for opsIDsResult.Next() {
		var opsID int
		if err := opsIDsResult.Scan(&opsID); err != nil {
			return nil, err
		}
		opsIDs = append(opsIDs, opsID)
	}

	return opsIDs, nil
}
func getSupervisor(con *sql.DB, empID int) ([]int, error) {
	getSupervisorIDQuery := fmt.Sprintf("SELECT supervisorId FROM employee WHERE id = %d", empID)
	supervisorIDResult, err := con.Query(getSupervisorIDQuery)
	if err != nil {
		return nil, err
	}
	defer supervisorIDResult.Close()

	var supervisorIDs []int
	for supervisorIDResult.Next() {
		var supervisorID int
		if err := supervisorIDResult.Scan(&supervisorID); err != nil {
			return nil, err
		}
		supervisorIDs = append(supervisorIDs, supervisorID)
	}

	return supervisorIDs, nil
}

func getReportingOpsManagers(con *sql.DB, empID int) ([]int, error) {
	ids := []int{}
	getOpsIDsQuery := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d AND (empRole = 4 OR empRole = 12)", empID)
	res, err := con.Query(getOpsIDsQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)

		if id != empID {
			subOpsIDs, err := getReportingOpsManagers(con, id)
			if err != nil {
				return nil, err
			}
			ids = append(ids, subOpsIDs...)
		}
	}

	return ids, nil
}
func getOpProjects(con *sql.DB, empID int) ([]int, error) {
	ids := []int{}
	getProjIDsQuery := fmt.Sprintf("SELECT id FROM project WHERE operations_manager = %d GROUP BY id", empID)
	res, err := con.Query(getProjIDsQuery)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
func getTrainerTarget(con *sql.DB, empID int, projectArray []int) (int, error) {
	targetQuery := fmt.Sprintf("SELECT SUM(target) AS total FROM project_emps WHERE emp_id = %d AND project_id IN (%s)", empID, intsToString(projectArray))
	res := con.QueryRow(targetQuery)

	var total int
	if err := res.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func intsToString(ints []int) string {
	strs := make([]string, len(ints))
	for i, val := range ints {
		strs[i] = fmt.Sprintf("%d", val)
	}
	return strings.Join(strs, ",")
}

func getGFData(con *sql.DB, filter string, sessionType int, empID int) (int, error) {
	filter += fmt.Sprintf(" AND tp.user_id = %d", empID)
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) AS visit FROM tbl_poa tp "+
		"INNER JOIN project p ON p.id = tp.project_id "+
		"WHERE `type` = 2 AND session_type = %d AND tp.check_out IS NOT NULL %s", sessionType, filter)
	res := con.QueryRow(getVisit)

	var visit int
	if err := res.Scan(&visit); err != nil {
		return 0, err
	}

	return visit, nil
}

func getAssociatedProjectList(con *sql.DB, projID int) ([]int, error) {
	projArray := []int{}
	if projID != 0 {
		getProjList := fmt.Sprintf("SELECT associatedProject FROM project_association WHERE projectId IN (%d)", projID)
		projArray = append(projArray, projID)
		res, err := con.Query(getProjList)
		if err != nil {
			return nil, err
		}
		defer res.Close()

		for res.Next() {
			var associatedProject int
			if err := res.Scan(&associatedProject); err != nil {
				return nil, err
			}
			projArray = append(projArray, associatedProject)
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
			// Handle the error, e.g., log or return an error
			log.Println(err)
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
			log.Println(err)
		}
	} else {
		villageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' AND user_id = %s", startDate, endDate, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE type = 1 AND added = 0 AND startDate >= '%s' AND endDate <= '%s' AND user_id = %s", startDate, endDate, user_id)

		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
			log.Println(err)
		}

		err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
		if err != nil {
			// Handle the error, e.g., log or return an error
			log.Println(err)
		}
	}

	return villageCount + subVillageCount
}

func TrainerTarget(db *sql.DB, empID int, projectArray []int) int {
	projectIDs := make([]string, len(projectArray))
	for i, id := range projectArray {
		projectIDs[i] = fmt.Sprintf("%d", id)
	}

	query := fmt.Sprintf("SELECT SUM(target) AS total FROM project_emps WHERE emp_id = %d AND project_id IN (%s)", empID, strings.Join(projectIDs, ","))

	var total sql.NullInt64
	err := db.QueryRow(query).Scan(&total)
	if err != nil {
		log.Println("ERROR>>", err)
	}

	if total.Valid {
		return int(total.Int64)
	}

	return 0
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
			getActualsQuery += fmt.Sprintf(" AND participant_day1 BETWEEN '%s' AND '%s'", startDate, endDate)
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
