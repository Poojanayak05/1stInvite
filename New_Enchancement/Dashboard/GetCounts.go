package Dashboard

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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

func GreennewVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, filter string) int {
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
		return 0
	}
	defer villageResult.Close()

	var villageCount int
	for villageResult.Next() {
		if err := villageResult.Scan(&villageCount); err != nil {
			return 0
		}
	}

	subVillageResult, err := con.Query(subVillageQuery)
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

func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, val := range arr {
		strArr[i] = fmt.Sprint(val)
	}
	return strings.Join(strArr, ", ")
}

func VyaparnewVillageCount(con *sql.DB, startDate, endDate string, projectArray []int, filter string) int {
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
		return 0
	}
	defer villageResult.Close()

	var villageCount int
	for villageResult.Next() {
		if err := villageResult.Scan(&villageCount); err != nil {
			return 0
		}
	}

	subVillageResult, err := con.Query(subVillageQuery)
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

func getGelathi(db *sql.DB, startDate string, endDate string, projectArray []int, GalathiId string, funderId string, filter string) int {
	var gelathiCount int
	// get associated projects for each project
	for _, proj := range projectArray {
		projs, _ := GetAssociatedProjectList(db, proj)
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
		projs, _ := GetAssociatedProjectList(con, proj)
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
		projs, _ := GetAssociatedProjectList(db, proj)
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

func newVillageCount(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	// var projs []int
	for _, proj := range projectArray {
		// check if there are any associated project for each project
		projs, _ := GetAssociatedProjectList(db, proj)
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

func intsToString(ints []int) string {
	var stringSlice []string
	for _, i := range ints {
		stringSlice = append(stringSlice, strconv.Itoa(i))
	}
	return strings.Join(stringSlice, ",")
}

func GetOpsManagers(db *sql.DB, empId int) []int {
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

func GetAssociatedProjectList(DB *sql.DB, projId int) ([]int, error) {
	projArray := []int{}
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
