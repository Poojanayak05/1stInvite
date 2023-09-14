package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ParticipantRequest struct {
	ParticipantId      int    `json:"participant_id"`
	ProjectID          int    `json:"project_id"`
	TbID               int    `json:"tb_id"`
	GelathiRecomm      int    `json:"gelathiRecomm"`
	FirstName          string `json:"firstName"`
	Age                string `json:"age"`
	HusbandName        string `json:"husbandName"`
	Occupation         string `json:"occupation"`
	HusbandOccupation  string `json:"husbandOccupation"`
	Dob                string `json:"dob"`
	NameOfSHG          string `json:"nameOfSHG"`
	TypeOfEnterprise   string `json:"typeOfEnterprise"`
	WifeIncomeMonthly  int    `json:"wifeIncomeMonthly"`
	WifeSavingsMonthly int    `json:"wifeSavingsMonthly"`
	SavingAmt          int    `json:"saving_amt"`
	SavingGoal         string `json:"saving_goal"`
	NumOfChildren      string `json:"numOfChildren"`
	Income             int    `json:"income"`
	Savings            string `json:"savings"`
	Education          string `json:"education"`
	Caste              string `json:"caste"`
	ContactNo          string `json:"contact_no"`
	BankAcc            int    `json:"bank_acc"`
	FinalSave          int    `json:"final_save"`
	ParticipantDay2    string `json:"participant_day2"`
}

func EditParticipantMehtod(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request ParticipantRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := ""
	// params := []interface{}{}

	if request.ProjectID != 0 {
		query += fmt.Sprintf(" project_id = '%d',", request.ProjectID)
	}

	if request.TbID != 0 {
		query += fmt.Sprintf(" tb_id = '%d',", request.TbID)
	}

	if request.GelathiRecomm != 0 {
		query += fmt.Sprintf(" gelathiRecomm = '%d',", request.GelathiRecomm)
	} else {
		request.GelathiRecomm = 0
		query += fmt.Sprintf(" gelathiRecomm = '%d',", request.GelathiRecomm)
	}

	if request.FirstName != "" {
		query += fmt.Sprintf(" firstName = '%s',", request.FirstName)
	}

	if request.Age != "" {
		query += fmt.Sprintf(" age = '%s',", request.Age)
	}

	if request.HusbandName != "" {
		query += fmt.Sprintf(" husbandName = '%s',", request.HusbandName)
	}

	if request.Occupation != "" {
		query += fmt.Sprintf(" occupation = '%s',", request.Occupation)
	}

	if request.HusbandOccupation != "" {
		query += fmt.Sprintf(" husbandOccupation = '%s',", request.HusbandOccupation)
	}

	if request.Dob != "" {
		query += fmt.Sprintf(" dob = '%s',", request.Dob)
	}

	if request.NameOfSHG != "" {
		query += fmt.Sprintf(" nameOfSHG = '%s',", request.NameOfSHG)
	}

	if request.TypeOfEnterprise != "" {
		query += fmt.Sprintf(" typeOfEnterprise = '%s',", request.TypeOfEnterprise)
	}

	if request.WifeIncomeMonthly != 0 {
		query += fmt.Sprintf(" wifeIncomeMonthly = '%d',", request.WifeIncomeMonthly)
	}

	if request.WifeSavingsMonthly != 0 {
		query += fmt.Sprintf(" wifeSavingsMonthly = '%d',", request.WifeSavingsMonthly)
	}

	if request.SavingAmt != 0 {
		query += fmt.Sprintf(" saving_amt = '%d',", request.SavingAmt)
	}

	if request.SavingGoal != "" {
		query += fmt.Sprintf(" saving_goal = '%s',", request.SavingGoal)
	}

	if request.NumOfChildren != "" {
		query += fmt.Sprintf(" numOfChildren = '%s',", request.NumOfChildren)
	}

	if request.Income != 0 {
		query += fmt.Sprintf(" income = '%d',", request.Income)
	}

	if request.Savings != "" {
		query += fmt.Sprintf(" savings = '%s',", request.Savings)
	}

	if request.Education != "" {
		query += fmt.Sprintf(" education = '%s',", request.Education)
	}

	if request.Caste != "" {
		query += fmt.Sprintf(" caste = '%s',", request.Caste)
	}

	if request.ContactNo != "" {
		query += fmt.Sprintf(" contact_no = '%s',", request.ContactNo)
	}

	if request.BankAcc != 0 {
		query += fmt.Sprintf(" bank_acc = '%d',", request.BankAcc)
	} else {
		request.BankAcc = 0
		query += fmt.Sprintf(" bank_acc = '%d',", request.BankAcc)
	}

	if request.FinalSave != 0 {
		query += fmt.Sprintf(" final_save = '%d',", request.FinalSave)
	}

	// Remove the trailing comma from the query
	query = strings.TrimRight(query, ",")

	updateQuery := fmt.Sprintf("UPDATE training_participants SET %s WHERE id = ?", query)

	// Execute the SQL query with the parameters
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to update training participants data", "success": false})

	}

	defer stmt.Close()
	sqlResult, err := stmt.Exec(request.ParticipantId)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to executes a prepared statement", "success": false})

	}
	rowsAffected, _ := sqlResult.RowsAffected()

	if rowsAffected > 0 {
		// Update successfully

		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "training participants updated successfully", "success": true})

	} else {
		// Update did not affect any rows
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 404, "message": "training participants updated Unsuccessfully", "success": false})

	}

}
