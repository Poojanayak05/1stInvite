package vyapar

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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

func getReportingGFlManagers(con *sql.DB, empID int) ([]int, error) {
	ids := []int{}
	getOpsIDsQuery := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d AND (empRole = 6)", empID)
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
			subOpsIDs, err := getReportingGFlManagers(con, id)
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

func noofvyaparsurvey(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparsurvey string
	var noofvyaparsurvey int

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
			vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between '%s' and '%s' and t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where t.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(vyaparsurvey).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Println("noofvyaparsurvey", err)
	}
	return noofvyaparsurvey
}

func noofvyaparmodulecompleted(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparmodule string
	var noofvyaparmodule int
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
			vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId where t.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(vyaparmodule).Scan(&noofvyaparmodule)
	if err != nil {
		log.Println("noofvyaparmodule", err)
	}
	return noofvyaparmodule
}

func noofvyaparcoharts(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
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
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) WHERE and tbl_poa.date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparecoharts", err)
	}
	return noofvyaparcoharts
}

func VyaparnewVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
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
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and t.project_id IN (%s)", intArrayToString(projectArray))
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) and t.project_id IN (%s)", intArrayToString(projectArray))

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' %s", startDate, endDate, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' %s", startDate, endDate, filter)
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
func vyaparenrolled(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparenroll string
	var noofvyaparenroll int

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
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants WHERE (VyaparEnrollment = 1 or new_vyapar = 1) AND VyaparEnrollmentDate BETWEEN '%s' and '%s' and project_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants where (VyaparEnrollment = 1 or new_vyapar = 1) and project_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(id) AS Vyapar FROM training_participants WHERE (VyaparEnrollment = 1 or new_vyapar = 1) AND VyaparEnrollmentDate BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err := db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
	if err != nil {
		log.Println("noofvyaparenroll", err)
	}
	return noofvyaparenroll
}

type CacheEntry struct {
	Result     int
	Expiration time.Time
}

var cache = make(map[string]CacheEntry)
var cacheLock sync.RWMutex

func vyaparenrolledCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "vyaparenrolled")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := vyaparenrolled(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}

func noofvyaparcohartsCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "noofvyaparcoharts")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := noofvyaparcoharts(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}

func noofvyaparsurveyCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "noofvyaparsurvey")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := noofvyaparsurvey(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}

func noofvyaparmodulecompletedCached(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "noofvyaparmodulecompleted")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result
	}

	result := noofvyaparmodulecompleted(db, startDate, endDate, projectArray, filter)

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result
}
func VyaparnewVillageCountCached(con *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s-%s", startDate, endDate, fmt.Sprint(projectArray), filter, "VyaparnewVillageCount")

	cacheLock.RLock()
	cachedEntry, found := cache[cacheKey]
	cacheLock.RUnlock()

	if found && time.Now().Before(cachedEntry.Expiration) {
		return cachedEntry.Result, nil
	}

	result, err := VyaparnewVillageCount(con, startDate, endDate, projectArray, filter)
	if err != nil {
		return 0, err
	}

	cacheLock.Lock()
	cache[cacheKey] = CacheEntry{
		Result:     result,
		Expiration: time.Now().Add(time.Minute * 5), // Cache result for 5 minutes
	}
	cacheLock.Unlock()

	return result, nil
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

func Gflvyaparenrolled(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparenroll string
	var noofvyaparenroll int

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

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(t.id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 or t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' and '%s' and t.project_id IN (%s) and t.gelathi_id IN (%s)  %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		} else {
			vyaparenroll = fmt.Sprintf("SELECT COUNT(t.id) AS Vyapar FROM training_participants t where (t.VyaparEnrollment = 1 or t.new_vyapar = 1) and t.project_id IN (%s) and t.gelathi_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparenroll = fmt.Sprintf("SELECT COUNT(t.id) AS Vyapar FROM training_participants t WHERE (t.VyaparEnrollment = 1 or t.new_vyapar = 1) AND t.VyaparEnrollmentDate BETWEEN '%s' and '%s' ", startDate, endDate)
	}

	err = db.QueryRow(vyaparenroll).Scan(&noofvyaparenroll)
	if err != nil {
		log.Println("noofvyaparenroll", err)
	}
	return noofvyaparenroll
}

func Gflnoofvyaparsurvey(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparsurvey string
	var noofvyaparsurvey int

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

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between '%s' and '%s' and t.project_id IN (%s) and b.gfId IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		} else {
			vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where t.project_id IN (%s) and b.gfId IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparsurvey = fmt.Sprintf("SELECT count(b.id) FROM bdms_staff.BuzzVyaparProgramBaseline b join training_participants t on t.id= b.partcipantId where b.entry_date between BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err = db.QueryRow(vyaparsurvey).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Println("noofvyaparsurvey", err)
	}
	return noofvyaparsurvey
}

func Gflnoofvyaparmodulecompleted(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparmodule string
	var noofvyaparmodule int
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

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' AND '%s' AND t.project_id IN (%s) and b.gfId IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		} else {
			vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId where t.project_id IN (%s) and b.gfId IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparmodule = fmt.Sprintf("SELECT count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) FROM BuzzVyaparProgramBaseline b JOIN training_participants t ON t.id=b.partcipantId WHERE b.entry_date BETWEEN '%s' and '%s' ", startDate, endDate)
	}

	err = db.QueryRow(vyaparmodule).Scan(&noofvyaparmodule)
	if err != nil {
		log.Println("noofvyaparmodule", err)
	}
	return noofvyaparmodule
}

func Gflnoofvyaparcoharts(db *sql.DB, startDate string, endDate string, projectArray []int, gflid int, filter string) int {
	var vyaparcoharts string
	var noofvyaparcoharts int

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

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.date BETWEEN '%s' and '%s' and tbl_poa.project_id IN (%s) and tbl_poa.user_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		} else {
			vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) and tbl_poa.project_id IN (%s) and tbl_poa.user_id IN (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		}
	} else {
		vyaparcoharts = fmt.Sprintf("SELECT COUNT(tbl_poa.session_type) AS noofvyaparcohorts FROM tbl_poa tbl_poa  WHERE tbl_poa.type = 2 AND (tbl_poa.session_type = 16 OR tbl_poa.session_type = 17 OR tbl_poa.session_type = 18 OR tbl_poa.session_type = 19 OR tbl_poa.session_type = 20 OR tbl_poa.session_type = 21) WHERE and tbl_poa.date BETWEEN '%s' and '%s'", startDate, endDate)
	}

	err = db.QueryRow(vyaparcoharts).Scan(&noofvyaparcoharts)
	if err != nil {
		log.Println("noofvyaparecoharts", err)
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
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id where (t.VyaparEnrollment=1 or t.new_vyapar=1) and t.project_id IN (%s) and tp.user_id IN (%s) %s", intArrayToString(projectArray), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) and t.project_id IN (%s) and tp.user_id IN (%s) %s", intArrayToString(projectArray), strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)

		if startDate != "" && endDate != "" {
			villageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
			subVillageQuery += fmt.Sprintf(" AND `tp.date` BETWEEN '%s' AND '%s'", startDate, endDate)
		}
	} else {
		villageQuery = fmt.Sprintf("select COUNT(DISTINCT location_id) from tbl_poa tp  INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id  where (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' and tp.user_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id join training_participants t on t.project_id =tp.project_id join project p on p.id=t.project_id WHERE (t.VyaparEnrollment=1 or t.new_vyapar=1) AND p.startDate >= '%s' AND p.endDate <= '%s' and tp.user_id IN (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(gfid), " ", ",", -1), "[]"), filter)
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
