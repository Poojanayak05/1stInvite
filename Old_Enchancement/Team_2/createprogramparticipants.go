package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Struct fields
type ParticipantProgramRequest struct {
	ID         string `json:"id"`
	ProjectID  string `json:"project_id"`
	Tb_ID      string `json:"tb_id"`
	Trainer_ID string `json:"trainer_id"`

	FirstName   string `json:"firstName"`
	Age         string `json:"age"`
	HusbandName string `json:"husbandName"`
	// SavingGoal         string `json:"saving_goal"`
	// Day1               string `json:"day1"`
	// Day2               string `json:"day2"`
	// Participant_Day1   string `json:"participant_day1"`
	// Participant_Day2   string `json:"participant_day2"`
	// NumOfChildren      string `json:"numOfChildren"`
	// Income             string `json:"income"`
	//Savings            string `json:"savings"`
	Education  string `json:"education"`
	Contact_No string `json:"contact_no"`
	//WifeIncomeMonthly  string `json:"wifeIncomeMonthly"`
	//WifeSavingsMonthly string `json:"wifeSavingsMonthly"`
	//Saving_amt         string `json:"saving_amt"`
	Bank_acc   string `json:"bank_acc"`
	Caste      string `json:"caste"`
	Gelathi_id string `json:"gelathi_id"`
	NameOfSHG  string `json:"nameOfSHG"`
	//TypeOfEnterprise string `json:"typeOfEnterprise"`
	//DoB              string `json:"dob"`
	Type string `json:"type"`
}

func CreateProgramParticipant(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var request ParticipantProgramRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	Tb_ID, _ := strconv.Atoi(request.Tb_ID)
	//GelathiRecomm, _ := strconv.Atoi(request.GelathiRecomm)
	Caste, _ := strconv.Atoi(request.Caste)
	Contact_No, _ := strconv.Atoi(request.Contact_No)
	//Day1, _ := strconv.Atoi(request.Day1)
	//fmt.Println(Day1)
	//Day2, _ := strconv.Atoi(request.Day2)
	Age, _ := strconv.Atoi(request.Age)
	gelathiid, _ := strconv.Atoi(request.Gelathi_id)
	ID, _ := strconv.Atoi(request.ID)
	ProjectID, _ := strconv.Atoi(request.ProjectID)
	//Trainer_ID, _ := strconv.Atoi(request.Trainer_ID)

	//check  Duplicate participant data
	duplicateCheck := fmt.Sprintf("SELECT * FROM training_participants WHERE project_id = %d AND tb_id = %d AND gelathi_id = %d AND firstName = '%s' AND husbandName = '%s' ", ProjectID, Tb_ID, gelathiid, request.FirstName, request.HusbandName)
	fmt.Println(duplicateCheck)
	rows, err := DB.Query(duplicateCheck)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	fmt.Println("duplicateCheck", duplicateCheck)
	defer rows.Close()
	rowCount := 0
	for rows.Next() {
		rowCount++
	}
	fmt.Println("first count", rowCount)

	if rowCount != 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Duplicate Participant", "success": false})
		return
	}

	// Participant_Day1, _ := parseDateTime(request.Participant_Day1)
	// Participant_Day2, _ := parseDateTime(request.Participant_Day2)

	//Day1 = 1
	insertQueryPrep := "INSERT INTO training_participants (id, project_id, tb_id, firstName, age,husbandName,  education, contact_no,  caste,gelathi_id, nameOfSHG) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	result, err := DB.Exec(insertQueryPrep, ID, ProjectID, Tb_ID, request.FirstName, Age, request.HusbandName, request.Education, Contact_No, Caste, gelathiid, request.NameOfSHG)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to insert participant", "success": false, "error": err})
		return
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		// Handle the error
		fmt.Println("Failed to get the number of rows affected:", err)
		return
	}

	// Check if the query executed successfully
	if rowsAffected != 0 {
		fmt.Println("success insert")
		query1 := fmt.Sprintf("select id from training_participants where project_id = %d and tb_id = %d and gelathi_id = %d and firstName = '%s' and husbandName = '%s' and contact_no = %d", ProjectID, Tb_ID, gelathiid, request.FirstName, request.HusbandName, Contact_No)
		rows1, err1 := DB.Query(query1)
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}
		fmt.Println("uuuuuuuuu", query1)
		defer rows1.Close()
		var projectid int
		for rows1.Next() {
			err := rows1.Scan(&projectid)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
		}
		fmt.Println("projectid", projectid)
		var query string
		if request.Type != "" {
			fmt.Println("entering type")
			if request.Type == "vyapar" {
				fmt.Println("entering vyapar")
				query = fmt.Sprintf("update training_participants set VyaparEnrollment =1,new_vyapar=1,VyaparEnrollmentDate= '%s'  where id = %d", time.Now().Format("2006/01/02"), projectid)
				rows, err := DB.Query(query)
				fmt.Println("entering type", query)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in updating", "success": false, "error": err})
					return
				}
				defer rows.Close()
			} else if request.Type == "green" {
				fmt.Println("entering green")
				query = fmt.Sprintf("update training_participants set GreenMotivators =1,new_green=1,GreenMotivatorsDate= '%s' where id = %d", time.Now().Format("2006/01/02"), projectid)
				rows, err := DB.Query(query)
				fmt.Println("entering type", query)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in updating", "success": false, "error": err})
					return
				}
				defer rows.Close()
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "Participant Added Successfully", "success": true})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 404, "message": "Participant Added Unsuccessfully", "success": false})
	}

}

// func parseDateTime1(datetime string) (time.Time, error) {
// 	t, err := time.Parse("2006-01-02", datetime)

// 	if err != nil {
// 		fmt.Println("line 248", err)
// 		return time.Time{}, err
// 	}

// 	fmt.Println("time: ", t)
// 	return t, nil
// }
