package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type BuzzVyaparProgramBaseline struct {
	PartcipantID int `json:"partcipantId"`
}

type Querydata struct {
	ID                                                   string `json:"id"`
	PartcipantID                                         string `json:"partcipantId"`
	GfID                                                 string `json:"gfId"`
	EntryDate                                            string `json:"entry_date"`
	WhenWasSurveyDone                                    string `json:"when_was_survey_done"`
	NameOfTheVyapari                                     string `json:"name_of_the_vyapari"`
	Age                                                  string `json:"age"`
	ContactNumber                                        string `json:"contact_number"`
	VillageID                                            string `json:"village_id"`
	LocationCircle                                       string `json:"location_circle"`
	HighterEducation                                     string `json:"highter_education"`
	MaritalStatus                                        string `json:"marital_status"`
	NumberOfPeopleInTheHousehold                         string `json:"number_of_people_in_the_household"`
	DoYouOwnASmartPhone                                  string `json:"do_you_own_a_smart_phone"`
	DoYouHaveInternetConnectionOnYourSmartPhone          string `json:"do_you_have_internet_connection_on_your_smart_phone"`
	SectorTypeOfBusiness                                 string `json:"sector_type_of_business"`
	AreYouProficientWithNumbers                          string `json:"are_you_proficient_with_numbers"`
	AreYouProficientWithWrittenLanguage                  string `json:"are_you_proficient_with_written_language"`
	HouseholdIncomeMonthly                               string `json:"household_income_monthly"`
	OverTheLastMonthYourAverageIncome                    string `json:"over_the_last_month_your_average_income"`
	YourBusinessProfitLastMonth                          string `json:"your_business_profit_last_month"`
	HowMuchMonthlyIncomeWouldYouLikeToIdeallyEarn        string `json:"how_much_monthly_income_would_you_like_to_ideally_earn"`
	AmountInvestedWhenTheBusinessStarted                 string `json:"amount_invested_when_the_business_started"`
	NumberOfYearsTheBusinessHasBeenOperating             string `json:"number_of_years_the_business_has_been_operating"`
	WhyDoYouDoBusiness                                   string `json:"why_do_you_do_business"`
	TellUsThreeThingsAboutYouAsAnEntrepreneur            string `json:"tell_us_three_things_about_you_as_an_entrepreneur"`
	TellUsThreeThingsAboutYourRoleAsAWomanAtHome         string `json:"tell_us_three_things_about_your_role_as_a_woman_at_home"`
	WhatAreYourChallengesInRunningAndGrowingYourBusiness string `json:"what_are_your_challenges_in_running_and_growing_your_business"`
	WhatIsYourPlanToOvercomeTheseChallenges              string `json:"what_is_your_plan_to_overcome_these_challenges"`
	WhatAreYourSkills                                    string `json:"what_are_your_skills"`
	WhatAreTheResourcesAvailableWithYouForYourBusiness   string `json:"what_are_the_resources_available_with_you_for_your_business"`
	WhoIsYourCustomerDescribeThemToUs                    string `json:"who_is_your_customer_Describe_them_to_us"`
	PleaseListDownTheVariousComponentsOfBusiness         string `json:"please_list_down_the_various_components_of_business"`
	IKnowTheCurrentStateOfMyBusinessInProfitLossRevenue  string `json:"I_know_the_current_state_of_my_business_in_profit_loss_revenue"`
	WhatKindOfBooksOfAccountsDoYouMaintain               string `json:"what_kind_of_books_of_accounts_do_you_maintain"`
	ICanGenerateIdeasToSolveMyBusinessProblems           string `json:"i_can_generate_ideas_to_solve_my_business_problems"`
	TellUsAboutOneBusinessProblem                        string `json:"tell_us_about_one_business_problem"`
	WhatIsYourBusinessGoalBusinessImpurumenetMadodu      string `json:"what_is_your_business_goal_Business_impurumenet_madodu"`
	DoYouHaveABusinessPlanToReachThatGoal                string `json:"do_you_have_a_business_plan_to_reach_that_goal"`
	CanYouSubmitABusinessPlanForYourGoalToUsRightNow     string `json:"can_you_submit_a_business_plan_for_your_goal_to_us_right_now"`
	WhatAreTheStrengthsOfYourBusiness                    string `json:"what_are_the_strengths_of_your_business"`
	WhatAreTheWeaknessesOfYourBusiness                   string `json:"what_are_the_weaknesses_of_your_business"`
	WhatAreTheOpportunitiesForYourBusiness               string `json:"what_are_the_opportunities_for_your_business"`
	AreYouAbleToRaiseTheRequiredFinance                  string `json:"are_you_able_to_raise_the_required_finance"`
	IHaveTakenALoanFrom                                  string `json:"i_have_taken_a_loan_from"`
	IHaveTroubleAccessingLoanForMyBusiness               string `json:"i_have_trouble_accessing_loan_for_my_business"`
	WhatAreThePrerequisitesToAccessALoan                 string `json:"what_are_the_prerequisites_to_access_a_loan"`
	TalukDistrict                                        string `json:"taluk_district"`
	NameOfTheCohort                                      string `json:"name_of_the_cohort"`
	YouStoppedHoldYourBusiness                           string `json:"you_stopped_hold_your_business"`
	NoHoursEngagedBusiness                               string `json:"no_hours_engaged_business"`
	LicenseForExistingBusiness                           string `json:"license_for_existing_business"`
	HomeBasedWorkFromShop                                string `json:"home_based_work_from_shop"`
	LoanCurrentlyAvailed                                 string `json:"loan_currently_availed"`
	RelationWhoBorrowed                                  string `json:"relation_who_borrowed"`
	LoanTotalAmount                                      string `json:"loan_total_amount"`
	LoanSource                                           string `json:"loan_source"`
	LoanRepaymentTillDate                                string `json:"loan_repayment_till_date"`
	NeedAdditionalSkillsBusiness                         string `json:"need_additional_skills_business"`
	SkillsWhatAreThose                                   string `json:"skills_what_are_those"`
	Module1                                              any    `json:"module1"`
	Module2                                              any    `json:"module2"`
	Module3                                              any    `json:"module3"`
	Module4                                              any    `json:"module4"`
	Module5                                              any    `json:"module5"`
}

// var data []Querydata

func GetBuzzVyaparProgramBaseline(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed to Read Bodyr", "success": false, "error": err.Error()})

		return
	}

	defer r.Body.Close()

	var req BuzzVyaparProgramBaseline
	err = json.Unmarshal(data, &req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to unmarshal", "success": false, "error": err.Error()})
		return
	}

	var query string
	if req.PartcipantID != 0 {

		query = fmt.Sprintf("SELECT * FROM BuzzVyaparProgramBaseline WHERE partcipantId  = %d", req.PartcipantID)
		rows, err := db.Query(query)

		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Query Error", "success": false, "error": err.Error()})
			return
		}
		defer rows.Close()

		var response []Querydata

		for rows.Next() {

			var queryData Querydata

			err := rows.Scan(
				&queryData.ID,
				&queryData.PartcipantID,
				&queryData.GfID,
				&queryData.EntryDate,
				&queryData.WhenWasSurveyDone,
				&queryData.NameOfTheVyapari,
				&queryData.Age,
				&queryData.ContactNumber,
				&queryData.VillageID,
				&queryData.LocationCircle,
				&queryData.HighterEducation,
				&queryData.MaritalStatus,
				&queryData.NumberOfPeopleInTheHousehold,
				&queryData.DoYouOwnASmartPhone,
				&queryData.DoYouHaveInternetConnectionOnYourSmartPhone,
				&queryData.SectorTypeOfBusiness,
				&queryData.AreYouProficientWithNumbers,
				&queryData.AreYouProficientWithWrittenLanguage,
				&queryData.HouseholdIncomeMonthly,
				&queryData.OverTheLastMonthYourAverageIncome,
				&queryData.YourBusinessProfitLastMonth,
				&queryData.HowMuchMonthlyIncomeWouldYouLikeToIdeallyEarn,
				&queryData.AmountInvestedWhenTheBusinessStarted,
				&queryData.NumberOfYearsTheBusinessHasBeenOperating,
				&queryData.WhyDoYouDoBusiness,
				&queryData.TellUsThreeThingsAboutYouAsAnEntrepreneur,
				&queryData.TellUsThreeThingsAboutYourRoleAsAWomanAtHome,
				&queryData.WhatAreYourChallengesInRunningAndGrowingYourBusiness,
				&queryData.WhatIsYourPlanToOvercomeTheseChallenges,
				&queryData.WhatAreYourSkills,
				&queryData.WhatAreTheResourcesAvailableWithYouForYourBusiness,
				&queryData.WhoIsYourCustomerDescribeThemToUs,
				&queryData.PleaseListDownTheVariousComponentsOfBusiness,
				&queryData.IKnowTheCurrentStateOfMyBusinessInProfitLossRevenue,
				&queryData.WhatKindOfBooksOfAccountsDoYouMaintain,
				&queryData.ICanGenerateIdeasToSolveMyBusinessProblems,
				&queryData.TellUsAboutOneBusinessProblem,
				&queryData.WhatIsYourBusinessGoalBusinessImpurumenetMadodu,
				&queryData.DoYouHaveABusinessPlanToReachThatGoal,
				&queryData.CanYouSubmitABusinessPlanForYourGoalToUsRightNow,
				&queryData.WhatAreTheStrengthsOfYourBusiness,
				&queryData.WhatAreTheWeaknessesOfYourBusiness,
				&queryData.WhatAreTheOpportunitiesForYourBusiness,
				&queryData.AreYouAbleToRaiseTheRequiredFinance,
				&queryData.IHaveTakenALoanFrom,
				&queryData.IHaveTroubleAccessingLoanForMyBusiness,
				&queryData.WhatAreThePrerequisitesToAccessALoan,
				&queryData.TalukDistrict,
				&queryData.NameOfTheCohort,
				&queryData.YouStoppedHoldYourBusiness,
				&queryData.NoHoursEngagedBusiness,
				&queryData.LicenseForExistingBusiness,
				&queryData.HomeBasedWorkFromShop,
				&queryData.LoanCurrentlyAvailed,
				&queryData.RelationWhoBorrowed,
				&queryData.LoanTotalAmount,
				&queryData.LoanSource,
				&queryData.LoanRepaymentTillDate,
				&queryData.NeedAdditionalSkillsBusiness,
				&queryData.SkillsWhatAreThose,
				&queryData.Module1,
				&queryData.Module2,
				&queryData.Module3,
				&queryData.Module4,
				&queryData.Module5,
			)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Scan Error", "success": false, "error": err.Error()})
				return
			}

			response = append(response, queryData)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": response, "success": true, "message": "BuzzVyaparProgramBaseline"})

	}

}
