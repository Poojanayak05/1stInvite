package green

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func intsToString(ints []int) string {
	strs := make([]string, len(ints))
	for i, val := range ints {
		strs[i] = fmt.Sprintf("%d", val)
	}
	return strings.Join(strs, ",")
}

func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, val := range arr {
		strArr[i] = fmt.Sprint(val)
	}
	return strings.Join(strArr, ", ")
}
func greenenrolled(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var greenenroll string
	var noofgreenenroll int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greenenroll = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants WHERE (GreenMotivators = 1 or new_green=1) AND GreenMotivatorsDate BETWEEN '%s' and '%s' and project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			greenenroll = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants WHERE (GreenMotivators = 1 or new_green=1) and project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		greenenroll = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants WHERE (GreenMotivators = 1 or new_green=1) AND GreenMotivatorsDate BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(greenenroll).Scan(&noofgreenenroll)
	if err != nil {
		log.Println("noofgreenenroll", err)
	}
	return noofgreenenroll
}

func noofgreensurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var greensurvey string
	var noofgreensurvey int

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId where g.entry_date BETWEEN '%s' and '%s' and t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId where t.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId where g.entry_date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(greensurvey).Scan(&noofgreensurvey)
	if err != nil {
		log.Println("noofgreensurvey", err)
	}
	return noofgreensurvey
}

func noofgreenmodulecompleted(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var greenmodule string
	var noofgreenmodule int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId WHERE g.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId where t.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId where g.entry_date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(greenmodule).Scan(&noofgreenmodule)
	if err != nil {
		log.Println("noofgreenmodule", err)
	}
	return noofgreenmodule
}

func noofgreencoharts(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparcoharts string
	var noofvyaparcoharts int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) and tbl_poa.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) WHERE and tbl_poa.date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparecoharts", err)
	}
	return noofvyaparcoharts
}

func GreennewVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(con, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	projectArray = uniqueIntSlice(projectArray)

	var villageQuery, subVillageQuery string

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id where (t.GreenMotivators=1 or t.new_green=1) and t.project_id IN (%s)", intArrayToString(projectArray))
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) and t.project_id IN (%s)", intArrayToString(projectArray))

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' %s", startDate, endDate, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' %s", startDate, endDate, filter)
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
func removeDuplicates(arr []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for _, item := range arr {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}

	return result
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

func GfGreenVillageCount(db *sql.DB, startDate string, endDate string, projectArray []int, user_id string) int {
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

		villageQuery := fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id where (t.GreenMotivators=1 or t.new_green=1) and t.project_id IN (%s) AND tp.user_id = %s", projectIDs, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) t.project_id IN (%s) AND tp.user_id = %s", projectIDs, user_id)

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
		villageQuery := fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' AND tp.user_id = %s", startDate, endDate, user_id)
		subVillageQuery := fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' AND tp.user_id = %s", startDate, endDate, user_id)

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

func Gflgreenenrolled(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var greenenroll string
	var noofgreenenroll int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greenenroll = fmt.Sprintf("SELECT COUNT(t.id) as greenMoti FROM training_participants t join project p on p.id=t.project_id WHERE (t.GreenMotivators = 1 or t.new_green=1) AND t.GreenMotivatorsDate BETWEEN '%s' and '%s' and p.gfl_id = %d and t.project_id IN (%s)  %s", startDate, endDate, gflid, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			greenenroll = fmt.Sprintf("SELECT COUNT(t.id) as greenMoti FROM training_participants t join project p on p.id=t.project_id WHERE (t.GreenMotivators = 1 or t.new_green=1) and t.project_id IN (%s) and p.gfl_id = %d %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), gflid, filter)
		}
	} else {
		greenenroll = fmt.Sprintf("SELECT COUNT(t.id) as greenMoti FROM training_participants t join project p on p.id=t.project_id WHERE (t.GreenMotivators = 1 or t.new_green=1) AND t.GreenMotivatorsDate BETWEEN '%s' and '%s' and p.gfl_id = %d", startDate, endDate, gflid)
	}

	err := db.QueryRow(greenenroll).Scan(&noofgreenenroll)
	if err != nil {
		log.Println("noofgreenenroll", err)
	}
	return noofgreenenroll
}

func Gflnoofgreensurvey(db *sql.DB, startDate string, endDate string, projectArray []int, glfid int, filter string) int {
	var greensurvey string
	var noofgreensurvey int

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId join project p on p.id=t.project_id where g.entry_date BETWEEN '%s' and '%s' and p.gfl_id = %d and t.project_id IN (%s) %s", startDate, endDate, glfid, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId join project p on p.id=t.project_id where p.gfl_id = %d and t.project_id IN (%s) %s", glfid, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		greensurvey = fmt.Sprintf("SELECT count(g.id) FROM GreenBaselineSurvey g join training_participants t on t.id= g.partcipantId join project p on p.id=t.project_id where g.entry_date BETWEEN '%s' and '%s' and p.gfl_id = %d", startDate, endDate, glfid)
	}

	err := db.QueryRow(greensurvey).Scan(&noofgreensurvey)
	if err != nil {
		log.Println("noofgreensurvey", err)
	}
	return noofgreensurvey
}

func Gflnoofgreenmodulecompleted(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var greenmodule string
	var noofgreenmodule int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId join project p on p.id=t.project_id WHERE g.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) and p.gfl_id = %d %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), gflid, filter)
		} else {
			greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId join project p on p.id=t.project_id where t.project_id IN (%s) and p.gfl_id = %d %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), gflid, filter)
		}
	} else {
		greenmodule = fmt.Sprintf("SELECT count(module1=1 AND module2=1 AND module3=1 AND module4=1 AND module5=1) FROM GreenBaselineSurvey g JOIN training_participants t ON t.id=g.partcipantId join project p on p.id=t.project_id where g.entry_date BETWEEN '%s' and '%s' and p.gfl_id = %d", startDate, endDate, gflid)
	}

	err := db.QueryRow(greenmodule).Scan(&noofgreenmodule)
	if err != nil {
		log.Println("noofgreenmodule", err)
	}
	return noofgreenmodule
}

func Gflnoofgreencoharts(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparcoharts string
	var noofvyaparcoharts int
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa join project p on p.id=tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) and p.gfl_id = %d %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), gflid, filter)
		} else {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa join project p on p.id=tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) and tbl_poa.project_id IN (%s) and p.gfl_id = %d  %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), gflid, filter)
		}
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofgreencoharts FROM tbl_poa tbl_poa join project p on p.id=tbl_poa.project_id WHERE tbl_poa.type = 2 AND (tbl_poa.session_type IN (10, 11, 12, 13, 14, 15)) WHERE and tbl_poa.date BETWEEN '%s' and '%s' p.gfl_id = %d", startDate, endDate, gflid)
	}

	err := db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparecoharts", err)
	}
	return noofvyaparcoharts
}

func GflGreennewVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, gflid int, filter string) (int, error) {
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(con, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	// Remove duplicates from projectArray
	projectArray = removeDuplicates(projectArray)

	projectArray = uniqueIntSlice(projectArray)

	var villageQuery, subVillageQuery string

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id where (t.GreenMotivators=1 or t.new_green=1) and t.project_id IN (%s) and p.gfl_id = %d", intArrayToString(projectArray), gflid)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) and t.project_id IN (%s) and p.gfl_id = %d", intArrayToString(projectArray), gflid)

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' and p.gfl_id = %d %s", startDate, endDate, gflid, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.GreenMotivators=1 or t.new_green=1) AND p.startDate >= '%s' AND p.endDate <= '%s' and p.gfl_id = %d %s", startDate, endDate, gflid, filter)
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
