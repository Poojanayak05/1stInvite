package Team_4

//Done by Dhiraj Lakhane

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetParticipantData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	type Participant struct {
		ID                 int
		FirstName          string
		Age                string
		Occupation         string
		HusbandName        string
		SavingGoal         string
		Income             string
		Savings            string
		GelathiRecomm      string
		BusinessIncome     string
		BusinessProfit     string
		SavingAmt          string
		BankAcc            string
		Education          string
		ContactNo          string
		VoterID            string
		NameOfSHG          string
		TypeOfEnterprise   string
		DOB                string
		WifeIncomeMonthly  string
		WifeSavingsMonthly string
		Caste              string
		HusbandOccupation  string
	}

	type Response struct {
		ID                 string `json:"id"`
		FirstName          string `json:"firstName"`
		Age                string `json:"age"`
		Occupation         string `json:"occupation"`
		HusbandName        string `json:"husbandName"`
		SavingGoal         string `json:"saving_goal"`
		Income             string `json:"income"`
		Savings            string `json:"savings"`
		GelathiRecomm      string `json:"gelathiRecomm"`
		BusinessIncome     string `json:"business_income"`
		BusinessProfit     string `json:"business_profit"`
		SavingAmt          string `json:"saving_amt"`
		BankAcc            string `json:"bank_acc"`
		Education          string `json:"education"`
		ContactNo          string `json:"contact_no"`
		VoterID            string `json:"voter_id"`
		NameOfSHG          string `json:"nameOfSHG"`
		TypeOfEnterprise   string `json:"typeOfEnterprise"`
		DOB                string `json:"dob"`
		WifeIncomeMonthly  string `json:"wifeIncomeMonthly"`
		WifeSavingsMonthly string `json:"wifeSavingsMonthly"`
		Caste              string `json:"caste"`
		HusbandOccupation  string `json:"husbandOccupation"`
		Code               int    `json:"code"`
		Success            bool   `json:"success"`
		Message            string `json:"message"`
	}

	type input struct {
		ParticipantId int `json:"participant_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var data input
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
		return
	}

	if data.ParticipantId != 0 {
		query := fmt.Sprintf("SELECT tr_pat.id, tr_pat.firstName, IFNULL(tr_pat.age, '') as age, IFNULL(tr_pat.occupation,'') as occupation, IFNULL(tr_pat.husbandName, '') as husbandName, IFNULL(tr_pat.saving_goal,''), IFNULL(tr_pat.income, '') as income, IFNULL(tr_pat.savings, '') as savings, tr_pat.gelathiRecomm, IFNULL(tr_pat.business_income,'') as business_income, IFNULL(tr_pat.business_profit,'') as business_profit, IFNULL(tr_pat.saving_amt,'') as saving_amt, IFNULL(tr_pat.bank_acc,'') as bank_acc, IFNULL(tr_pat.education,'') as education, IFNULL(tr_pat.contact_no,'') as contact_no, IFNULL(tr_pat.voter_id,'') as voter_id,IFNULL(tr_pat.nameOfSHG,'') as nameOfSHG, IFNULL(tr_pat.typeOfEnterprise,'') as typeOfEnterprise,IFNULL(tr_pat.dob,'') as dob,IFNULL(tr_pat.wifeIncomeMonthly,'') as wifeIncomeMonthly, IFNULL(tr_pat.wifeSavingsMonthly,'') as wifeSavingsMonthly,IFNULL(c.name,'') as caste,IFNULL(tr_pat.husbandOccupation,'') as husbandOccupation FROM training_participants tr_pat left join caste c on tr_pat.caste = c.id WHERE tr_pat.id = '%d'", data.ParticipantId)

		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error", "success": false, "error": err.Error()})
			return
		}
		defer rows.Close()
		var participant Participant

		for rows.Next() {
			err := rows.Scan(
				&participant.ID,
				&participant.FirstName,
				&participant.Age,
				&participant.Occupation,
				&participant.HusbandName,
				&participant.SavingGoal,
				&participant.Income,
				&participant.Savings,
				&participant.GelathiRecomm,
				&participant.BusinessIncome,
				&participant.BusinessProfit,
				&participant.SavingAmt,
				&participant.BankAcc,
				&participant.Education,
				&participant.ContactNo,
				&participant.VoterID,
				&participant.NameOfSHG,
				&participant.TypeOfEnterprise,
				&participant.DOB,
				&participant.WifeIncomeMonthly,
				&participant.WifeSavingsMonthly,
				&participant.Caste,
				&participant.HusbandOccupation,
			)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})

				continue
			}
		}
		resp := Response{
			ID:                 fmt.Sprint(participant.ID),
			FirstName:          participant.FirstName,
			Age:                participant.Age,
			Occupation:         participant.Occupation,
			HusbandName:        participant.HusbandName,
			SavingGoal:         participant.SavingGoal,
			Income:             participant.Income,
			Savings:            participant.Savings,
			GelathiRecomm:      participant.GelathiRecomm,
			BusinessIncome:     participant.BusinessIncome,
			BusinessProfit:     participant.BusinessProfit,
			SavingAmt:          participant.SavingAmt,
			BankAcc:            participant.BankAcc,
			Education:          participant.Education,
			ContactNo:          participant.ContactNo,
			VoterID:            participant.VoterID,
			NameOfSHG:          participant.NameOfSHG,
			TypeOfEnterprise:   participant.TypeOfEnterprise,
			DOB:                participant.DOB,
			WifeIncomeMonthly:  participant.WifeIncomeMonthly,
			WifeSavingsMonthly: participant.WifeSavingsMonthly,
			Caste:              participant.Caste,
			HusbandOccupation:  participant.HusbandOccupation,
			Code:               200,
			Success:            true,
			Message:            "Successfully",
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)

	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "status": "Not Found", "message": "Id Not found"})
		return
	}

}
