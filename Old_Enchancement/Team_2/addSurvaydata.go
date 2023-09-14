package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Participant struct represents the participant data
type Participant struct {
	ID                         string   `json:"participantId"`
	PersonOfWorth              string   `json:"personOfWorth"`
	GoodQuality                string   `json:"goodQuality"`
	AmFailure                  string   `json:"amFailure"`
	HaveGoal                   string   `json:"haveGoal"`
	Goal                       string   `json:"goal"`
	PathwayToGoal              string   `json:"pathwayToGoal"`
	DisheartenedToProblems     string   `json:"disheartenedToProblems"`
	SolutionToProblems         string   `json:"solutionToProblems"`
	ImplementationPlan         string   `json:"implementationPlan"`
	NoChoiceForSolution        string   `json:"noChoiceForSolution"`
	ProfitForSarees            string   `json:"profitForSarees"`
	AnnualLoanInterest         string   `json:"annualLoanInterest"`
	SaveRegularly              string   `json:"saveRegularly"`
	WhereSaveMoney             string   `json:"whereSaveMoney"`
	FrequencyOfSaving          string   `json:"frequencyOfSaving"`
	OwnAsset                   string   `json:"ownAsset"`
	SeparateFinancialAsset     string   `json:"separateFinancialAsset"`
	SpendMoney                 string   `json:"spendMoney"`
	HaveLoan                   string   `json:"haveLoan"`
	LoanOnWhoseName            string   `json:"loanOnWhoseName"`
	GovernmentBank             string   `json:"governmentBank"`
	PrivateBank                string   `json:"privateBank"`
	LocalMFI                   string   `json:"localMFI"`
	NGO                        string   `json:"NGO"`
	MoneyLender                string   `json:"moneyLender"`
	Middleman                  string   `json:"middleman"`
	AgroProcessors             string   `json:"agroProcessors"`
	Parents                    string   `json:"parents"`
	Relatives                  string   `json:"relatives"`
	Neighbours                 string   `json:"neighbours"`
	Friends                    string   `json:"friends"`
	SocialWelfareDepartments   string   `json:"socialWelfareDepartments"`
	Coorperatives              string   `json:"coorperatives"`
	Others                     string   `json:"others"`
	OwnIncomeGeneration        string   `json:"ownIncomeGeneration"`
	FamilyIncomeGeneration     string   `json:"familyIncomeGeneration"`
	EducationReason            string   `json:"educationReason"`
	FutureEmployment           string   `json:"futureEmployment"`
	OwnMarriage                string   `json:"ownMarriage"`
	BrotherMarriage            string   `json:"brotherMarriage"`
	PersonalExpenses           string   `json:"personalExpenses"`
	HouseholdUse               string   `json:"householdUse"`
	HouseRepair                string   `json:"houseRepair"`
	Medical                    string   `json:"medical"`
	Festival                   string   `json:"festival"`
	ReasonOthersToBorrowLoan   []string `json:"reasonOthersToBorrowLoan"`
	SpecificGoalForSavings     string   `json:"specificGoalForSavings"`
	HowMuchSaveToAchieve       string   `json:"howMuchSaveToAchieve"`
	EducationDecision          string   `json:"educationDecision"`
	AccessToHealtcare          string   `json:"accessToHealtcare"`
	AccessToCredit             string   `json:"accessToCredit"`
	SavingMoney                string   `json:"savingMoney"`
	AssetPurchase              string   `json:"assetPurchase"`
	DayTodayExpenditure        string   `json:"dayTodayExpenditure"`
	Livelihood                 string   `json:"livelihood"`
	PartOfCollective           string   `json:"partOfCollective"`
	ImportantToShareTheirProb  string   `json:"importantToShareTheirProb"`
	ShareLearningWithCommunity string   `json:"shareLearningWithCommunity"`
	HouseholdBooksAccounts     string   `json:"household_books_accounts"`
	AccountsForSelfEnterprises string   `json:"accounts_for_Self_Enterprises"`
	AccountHousehold           string   `json:"account_household"`
	AccountBusiness            string   `json:"account_business"`
}

func AddSurveydata(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	var p Participant
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})
		return
	}

	updateQuery := `
		UPDATE training_participants SET
			personOfWorth = ?,
			goodQuality = ?,
			amFailure = ?,
			haveGoal = ?,
			goal = ?,
			pathwayToGoal = ?,
			disheartenedToProblems = ?,
			solutionToProblems = ?,
			implementationPlan = ?,
			noChoiceForSolution = ?,
			profitForSarees = ?,
			annualLoanInterest = ?,
			saveRegularly = ?,
			whereSaveMoney = ?,
			frequencyOfSaving = ?,
			ownAsset = ?,
			separateFinancialAsset = ?,
			spendMoney = ?,
			haveLoan = ?,
			loanOnWhoseName = ?,
			governmentBank = ?,
			privateBank = ?,
			localMFI = ?,
			NGO = ?,
			moneyLender = ?,
			middleman = ?,
			agroProcessors = ?,
			parents = ?,
			relatives = ?,
			neighbours = ?,
			friends = ?,
			socialWelfareDepartments = ?,
			coorperatives = ?,
			others = ?,
			ownIncomeGeneration = ?,
			familyIncomeGeneration = ?,
			educationReason = ?,
			futureEmployment = ?,
			ownMarriage = ?,
			brotherMarriage = ?,
			personalExpenses = ?,
			householdUse = ?,
			houseRepair = ?,
			medical = ?,
			festival = ?,
			reasonOthersToBorrowLoan = ?,
			specificGoalForSavings = ?,
			howMuchSaveToAchieve = ?,
			educationDecision = ?,
			accessToHealtcare = ?,
			accessToCredit = ?,
			savingMoney = ?,
			assetPurchase = ?,
			dayTodayExpenditure = ?,
			livelihood = ?,
			partOfCollective = ?,
			importantToShareTheirProb = ?,
			shareLearningWithCommunity = ?,
			household_books_accounts = ?,
			accounts_for_Self_Enterprises = ?,
			account_household = ?,
			account_business = ?
		WHERE id = ?
	`

	stmt, err := DB.Prepare(updateQuery)
	if err != nil {
		res := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to prepare update query",
			Success: false,
		}
		//	w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		p.PersonOfWorth,
		p.GoodQuality,
		p.AmFailure,
		p.HaveGoal,
		p.Goal,
		p.PathwayToGoal,
		p.DisheartenedToProblems,
		p.SolutionToProblems,
		p.ImplementationPlan,
		p.NoChoiceForSolution,
		p.ProfitForSarees,
		p.AnnualLoanInterest,
		p.SaveRegularly,
		p.WhereSaveMoney,
		p.FrequencyOfSaving,
		p.OwnAsset,
		p.SeparateFinancialAsset,
		p.SpendMoney,
		p.HaveLoan,
		p.LoanOnWhoseName,
		p.GovernmentBank,
		p.PrivateBank,
		p.LocalMFI,
		p.NGO,
		p.MoneyLender,
		p.Middleman,
		p.AgroProcessors,
		p.Parents,
		p.Relatives,
		p.Neighbours,
		p.Friends,
		p.SocialWelfareDepartments,
		p.Coorperatives,
		p.Others,
		p.OwnIncomeGeneration,
		p.FamilyIncomeGeneration,
		p.EducationReason,
		p.FutureEmployment,
		p.OwnMarriage,
		p.BrotherMarriage,
		p.PersonalExpenses,
		p.HouseholdUse,
		p.HouseRepair,
		p.Medical,
		p.Festival,
		strings.Join(p.ReasonOthersToBorrowLoan, ","),
		p.SpecificGoalForSavings,
		p.HowMuchSaveToAchieve,
		p.EducationDecision,
		p.AccessToHealtcare,
		p.AccessToCredit,
		p.SavingMoney,
		p.AssetPurchase,
		p.DayTodayExpenditure,
		p.Livelihood,
		p.PartOfCollective,
		p.ImportantToShareTheirProb,
		p.ShareLearningWithCommunity,
		p.HouseholdBooksAccounts,
		p.AccountsForSelfEnterprises,
		p.AccountHousehold,
		p.AccountBusiness,
		p.ID,
	)
	if err != nil {
		res := Response{

			Code:    http.StatusInternalServerError,
			Message: "Failed to update survey data",
			Success: false,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle the error
		res := Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve rows affected",
			Success: false,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if rowsAffected > 0 {
		updateSurveyDone := fmt.Sprintf("UPDATE training_participants SET isSurveyDone = 1 WHERE id = %s AND personOfWorth IS NOT NULL AND personOfWorth != '' AND goodQuality IS NOT NULL AND goodQuality != ''", p.ID)
		_, err := DB.Exec(updateSurveyDone)
		if err != nil {
			// Handle the error
			res := Response{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update survey status",
				Success: false,
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		res := Response{
			Code:    http.StatusOK,
			Message: "Survey data added successfully",
			Success: true,
		}
		json.NewEncoder(w).Encode(res)
		return
	} else {
		res := Response{
			Code:    http.StatusBadRequest,
			Message: "Failed to add survey data",
			Success: false,
		}
		json.NewEncoder(w).Encode(res)
	}
}