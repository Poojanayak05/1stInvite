package Attendence

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type BuzzVyaparProgramBaseline struct {
	//id                                                             int       `json:"id"`
	PartcipantId                                                   string   `json:"partcipantId"`
	GfId                                                           string   `json:"gfId"`
	When_Was_survey_done                                           string   `json:"when_was_survey_done"`
	Name_of_the_vyapari                                            string   `json:"name_of_the_vyapari"`
	Age                                                            string   `json:"age"`
	Contact_number                                                 string   `json:"contact_number"`
	Village_id                                                     string   `json:"village_id"`
	Name_of_the_cohort                                             string   `json:"name_of_the_cohort"`
	Highter_education                                              string   `json:"highter_education"`
	Marital_status                                                 string   `json:"marital_status"`
	Number_of_people_in_the_household                              string   `json:"number_of_people_in_the_household"`
	Do_you_own_a_smart_phone                                       string   `json:"do_you_own_a_smart_phone"`
	Do_you_have_internet_connection_on_your_smart_phone            string   `json:"do_you_have_internet_connection_on_your_smart_phone"`
	Sector_type_of_business                                        string   `json:"sector_type_of_business"`
	Are_you_proficient_with_numbers                                string   `json:"are_you_proficient_with_numbers"`
	Are_you_proficient_with_written_language                       string   `json:"are_you_proficient_with_written_language"`
	Household_income_monthly                                       string   `json:"household_income_monthly"`
	Over_the_last_month_your_average_income                        string   `json:"over_the_last_month_your_average_income"`
	Your_business_profit_last_month                                string   `json:"your_business_profit_last_month"`
	How_much_monthly_income_would_you_like_to_ideally_earn         string   `json:"how_much_monthly_income_would_you_like_to_ideally_earn"`
	Amount_invested_when_the_business_started                      string   `json:"amount_invested_when_the_business_started"`
	Number_of_years_the_business_has_been_operating                string   `json:"number_of_years_the_business_has_been_operating"`
	You_stopped_hold_your_business                                 string   `json:"you_stopped_hold_your_business"`
	No_hours_engaged_business                                      string   `json:"no_hours_engaged_business"`
	License_for_existing_business                                  string   `json:"license_for_existing_business"`
	Home_based_work_from_shop                                      string   `json:"home_based_work_from_shop"`
	Why_do_you_do_business                                         string   `json:"why_do_you_do_business"`
	Tell_us_three_things_about_you_as_an_entrepreneur              []string `json:"tell_us_three_things_about_you_as_an_entrepreneur"`
	Tell_us_three_things_about_your_role_as_a_woman_at_home        string   `json:"tell_us_three_things_about_your_role_as_a_woman_at_home"`
	What_are_your_challenges_in_running_and_growing_your_business  string   `json:"what_are_your_challenges_in_running_and_growing_your_business"`
	What_is_your_plan_to_overcome_these_challenges                 string   `json:"what_is_your_plan_to_overcome_these_challenges"`
	What_are_your_skills                                           string   `json:"what_are_your_skills"`
	What_are_the_resources_available_with_you_for_your_business    string   `json:"what_are_the_resources_available_with_you_for_your_business"`
	Who_is_your_customer_Describe_them_to_us                       string   `json:"who_is_your_customer_Describe_them_to_us"`
	Please_list_down_the_various_components_of_business            []string `json:"please_list_down_the_various_components_of_business"`
	I_know_the_current_state_of_my_business_in_profit_loss_revenue string   `json:"I_know_the_current_state_of_my_business_in_profit_loss_revenue"`
	What_kind_of_books_of_accounts_do_you_maintain                 string   `json:"what_kind_of_books_of_accounts_do_you_maintain"`
	I_can_generate_ideas_to_solve_my_business_problems             string   `json:"i_can_generate_ideas_to_solve_my_business_problems"`
	Tell_us_about_one_business_problem                             string   `json:"tell_us_about_one_business_problem"`
	What_is_your_business_goal_Business_impurumenet_madodu         string   `json:"what_is_your_business_goal_Business_impurumenet_madodu"`
	Do_you_have_a_business_plan_to_reach_that_goal                 string   `json:"do_you_have_a_business_plan_to_reach_that_goal"`
	Can_you_submit_a_business_plan_for_your_goal_to_us_right_now   string   `json:"can_you_submit_a_business_plan_for_your_goal_to_us_right_now"`
	What_are_the_strenghts_of_your_business                        string   `json:"what_are_the_strenghts_of_your_business"`
	What_are_the_weaknesses_of_your_business                       string   `json:"what_are_the_weaknesses_of_your_business"`
	What_are_the_oppourtunities_for_your_business                  string   `json:"what_are_the_oppourtunities_for_your_business"`
	Are_you_able_to_raise_the_required_finance                     string   `json:"are_you_able_to_raise_the_required_finance"`
	I_have_taken_a_loan_from                                       string   `json:"i_have_taken_a_loan_from"`
	I_have_trouble_accessing_loan_for_my_business                  string   `json:"i_have_trouble_accessing_loan_for_my_business"`
	What_are_the_prerequisites_to_access_a_loan                    []string `json:"what_are_the_prerequisites_to_access_a_loan"`
	Loan_currently_availed                                         string   `json:"loan_currently_availed"`
	Need_additional_skills_business                                string   `json:"need_additional_skills_business"`
	Module1                                                        bool     `json:"Module1"`
	Module2                                                        bool     `json:"Module2"`
	Module3                                                        bool     `json:"Module3"`
	Module4                                                        bool     `json:"Module4"`
	Module5                                                        bool     `json:"Module5"`
	Module                                                         string   `json:"Module"`
}

type SpoorthiBaselineQuestionnaire struct {
	Id                                                  int    `json:"Id"`
	PartcipantId                                        int    `json:"partcipantId"`
	Email_address                                       string `json:"email_address"`
	GelathiId                                           int    `json:"GelathiId"`
	Entry_date                                          string `json:"entry_date"`
	Spoorthi_Session_Number                             string `json:"Spoorthi_Session_Number"`
	List_down_your_skills                               string `json:"list_down_your_skills"`
	Skills_to_overcome_my_challenges                    string `json:"Skills_to_overcome_my_challenges"`
	Used_skills_resources_combat_challenge              string `json:"used_skills_resources_combat_challenge"`
	Listen_paragraph                                    string `json:"listen_paragraph"`
	Summarize_main_points_paragraph                     string `json:"summarize_main_points_paragraph"`
	Ask_two_questions_help_you_understand               string `json:"ask_two_questions_help_you_understand"`
	Three_infrastructure_of_your_village                string `json:"three_infrastructure_of_your_village"`
	Know_the_need_of_my_community                       string `json:"know_the_need_of_my_community"`
	Together_community_members_community_infrastructure string `json:"together_community_members_community_infrastructure"`
	With_other_community_infrastructure                 string `json:"with_other_community_infrastructure"`
	Bring_someone_together                              string `json:"bring_someone_together"`
	Brought_people_together_incident                    string `json:"brought_people_together_incident"`
	Conflict_with_anyone_ask_position                   string `json:"conflict_with_anyone_ask_position"`
	Conflict_matters_interest_mine                      string `json:"conflict_matters_interest_mine"`
	Module1                                             bool   `json:"module1"`
	Module2                                             bool   `json:"module2"`
	Module3                                             bool   `json:"module3"`
	Module4                                             bool   `json:"module4"`
	Module5                                             bool   `json:"module5"`
	There_puja_at_my_house                              string `json:"there_puja_at_my_house"`
	Module                                              string `json:"Module"`
	FirstName                                           string `json:"firstname"`
}

//------------------------------------------- all attendence -------------------------------------------
func AllAttendence(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	type Attendence struct {
		PartcipantId int `json:"PartcipantId"`
		Type         int `json:"Type"`
		Circle_id    int `json:"Circle_id"`
	}
	var p Attendence
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// extra for not filling spoorthi form
	str := "select partcipantId from SpoorthiBaselineQuestionnaire"

	rows, err := DB.Query(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer rows.Close()
	var partcipantId int
	var result []int // creating slice

	for rows.Next() {

		err := rows.Scan(&partcipantId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		result = append(result, partcipantId)

	}

	found := false
	for _, r := range result {
		if r == p.PartcipantId {
			found = true
			break
		}
	}

	//end

	// ================ grren =========
	str1 := "select partcipantId from GreenBaselineSurvey"

	rows1, err := DB.Query(str1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer rows1.Close()
	var partcipantId1 int
	var result1 []int // creating slice

	for rows1.Next() {

		err := rows1.Scan(&partcipantId1)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		result1 = append(result1, partcipantId1)

	}

	ss := false
	for _, r := range result1 {
		if r == p.PartcipantId {
			ss = true
			break
		}
	}
	//end

	// vyapar
	// ================ grren =========
	str2 := "select partcipantId from BuzzVyaparProgramBaseline"

	rows2, err := DB.Query(str2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer rows2.Close()
	var partcipantId2 int
	var result2 []int // creating slice

	for rows2.Next() {

		err := rows2.Scan(&partcipantId2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		result2 = append(result2, partcipantId2)

	}

	founder1 := false
	for _, r := range result2 {
		if r == p.PartcipantId {
			founder1 = true
			break
		}
	}
	//end

	// ================================== spoorthi attendance ==============================
	switch {

	case p.Type == 5:
		if found {
			_, err := DB.Exec("UPDATE bdms_staff.SpoorthiBaselineQuestionnaire sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module1=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Spoorthi Module1"})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Fill the spoorthi survey form"})
		}
	case p.Type == 6:

		_, err := DB.Exec("UPDATE bdms_staff.SpoorthiBaselineQuestionnaire sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module2=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Spoorthi Module2"})

	case p.Type == 7:

		_, err := DB.Exec("UPDATE bdms_staff.SpoorthiBaselineQuestionnaire sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module3=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Spoorthi Module3"})

	case p.Type == 8:

		_, err := DB.Exec("UPDATE bdms_staff.SpoorthiBaselineQuestionnaire sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module4=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Spoorthi Module4"})

	case p.Type == 9:

		_, err := DB.Exec("UPDATE bdms_staff.SpoorthiBaselineQuestionnaire sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module5=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Spoorthi Module5"})

	// =============================== green attendance ==========================
	case p.Type == 11:
		if ss {
			_, err := DB.Exec("UPDATE bdms_staff.GreenBaselineSurvey sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module1=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Green Module1"})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Fill the green form"})
		}
	case p.Type == 12:
		_, err := DB.Exec("UPDATE bdms_staff.GreenBaselineSurvey sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module2=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Green Module2"})

	case p.Type == 13:
		_, err := DB.Exec("UPDATE bdms_staff.GreenBaselineSurvey sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module3=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Green Module3"})
	case p.Type == 14:
		_, err := DB.Exec("UPDATE bdms_staff.GreenBaselineSurvey sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module4=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Green Module4"})
	case p.Type == 15:
		_, err := DB.Exec("UPDATE bdms_staff.GreenBaselineSurvey sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module5=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Green Module5"})

	// ================================ Vyapar attendance =============================
	case p.Type == 17:
		if founder1 {
			_, err := DB.Exec("UPDATE bdms_staff.BuzzVyaparProgramBaseline sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module1=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for Vyapar Module1"})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Fill the vyapar form"})
		}
	case p.Type == 18:
		_, err := DB.Exec("UPDATE bdms_staff.BuzzVyaparProgramBaseline sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module2=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for vyapar Module2"})
	case p.Type == 19:
		_, err := DB.Exec("UPDATE bdms_staff.BuzzVyaparProgramBaseline sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module3=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for vyapar Module3"})
	case p.Type == 20:
		_, err := DB.Exec("UPDATE bdms_staff.BuzzVyaparProgramBaseline sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module4=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for vyapar Module4"})
	case p.Type == 21:
		_, err := DB.Exec("UUPDATE bdms_staff.BuzzVyaparProgramBaseline sp , bdms_staff.gelathi_circle gel_circle left join bdms_staff.training_participants tr_part ON gel_circle.gelathi_id = tr_part.id SET sp.module5=1 where sp.partcipantId=? and gel_circle.circle_id =?;", p.PartcipantId, p.Circle_id)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Attendance added successfully for vyapar Module5"})

	default:
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid Attendance"})

	}
}

//-------------------------------------------listing of spoorthi participants based on the Modules1 -----------------------------------

func Spoorthilist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var p SpoorthiBaselineQuestionnaire
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		panic(err.Error())
	}

	if p.Module == "Module1" {

		str := "SELECT train.firstName FROM training_participants train,SpoorthiBaselineQuestionnaire spoo where (spoo.partcipantId=train.id) and train.enroll=1"

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var FirstName string

		type sporti struct {
			FirstName string
		}

		var result []sporti // creating slice

		for rows.Next() {

			err := rows.Scan(&FirstName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: FirstName})

		}

		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

		//result = append(result, SpoorthiBaselineQuestionnaire{Id: Id, partcipantId: partcipantId, email_address: email_address, GelathiId: GelathiId, Spoorthi_Session_Number: Spoorthi_Session_Number, list_down_your_skills: list_down_your_skills, skills_to_overcome_my_challenges: skills_to_overcome_my_challenges, used_skills_resources_combat_challenge: used_skills_resources_combat_challenge, listen_paragraph: listen_paragraph, summarize_main_points_paragraph: summarize_main_points_paragraph, ask_two_questions_help_you_understand: ask_two_questions_help_you_understand, three_infrastructure_of_your_village: three_infrastructure_of_your_village, know_the_need_of_my_community: know_the_need_of_my_community, together_community_members_community_infrastructure: together_community_members_community_infrastructure, with_other_community_infrastructure: with_other_community_infrastructure, bring_someone_together: bring_someone_together, brought_people_together_incident: brought_people_together_incident, conflict_with_anyone_ask_position: conflict_with_anyone_ask_position, conflict_matters_interest_mine: conflict_matters_interest_mine, there_puja_at_my_house: there_puja_at_my_house, module1: module1, module2: module2, module3: module3, module4: module4, module5: module5})
	} else if p.Module == "Module2" {
		str := "SELECT firstName FROM training_participants train,SpoorthiBaselineQuestionnaire spoo where (spoo.partcipantId=train.id) and module1=1 and train.enroll=1;"

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var FirstName string

		type sporti struct {
			FirstName string
		}

		var result []sporti // creating slice

		for rows.Next() {

			err := rows.Scan(&FirstName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: FirstName})

		}

		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

		//result = append(result, SpoorthiBaselineQuestionnaire{Id: Id, partcipantId: partcipantId, email_address: email_address, GelathiId: GelathiId, Spoorthi_Session_Number: Spoorthi_Session_Number, list_down_your_skills: list_down_your_skills, skills_to_overcome_my_challenges: skills_to_overcome_my_challenges, used_skills_resources_combat_challenge: used_skills_resources_combat_challenge, listen_paragraph: listen_paragraph, summarize_main_points_paragraph: summarize_main_points_paragraph, ask_two_questions_help_you_understand: ask_two_questions_help_you_understand, three_infrastructure_of_your_village: three_infrastructure_of_your_village, know_the_need_of_my_community: know_the_need_of_my_community, together_community_members_community_infrastructure: together_community_members_community_infrastructure, with_other_community_infrastructure: with_other_community_infrastructure, bring_someone_together: bring_someone_together, brought_people_together_incident: brought_people_together_incident, conflict_with_anyone_ask_position: conflict_with_anyone_ask_position, conflict_matters_interest_mine: conflict_matters_interest_mine, there_puja_at_my_house: there_puja_at_my_house, module1: module1, module2: module2, module3: module3, module4: module4, module5: module5})
	} else if p.Module == "Module3" {
		str := "SELECT firstName FROM training_participants train,SpoorthiBaselineQuestionnaire spoo where (spoo.partcipantId=train.id) and module1=1 and module2=1;"

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var FirstName string

		type sporti struct {
			FirstName string
		}

		var result []sporti // creating slice

		for rows.Next() {

			err := rows.Scan(&FirstName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: FirstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

		//result = append(result, SpoorthiBaselineQuestionnaire{Id: Id, partcipantId: partcipantId, email_address: email_address, GelathiId: GelathiId, Spoorthi_Session_Number: Spoorthi_Session_Number, list_down_your_skills: list_down_your_skills, skills_to_overcome_my_challenges: skills_to_overcome_my_challenges, used_skills_resources_combat_challenge: used_skills_resources_combat_challenge, listen_paragraph: listen_paragraph, summarize_main_points_paragraph: summarize_main_points_paragraph, ask_two_questions_help_you_understand: ask_two_questions_help_you_understand, three_infrastructure_of_your_village: three_infrastructure_of_your_village, know_the_need_of_my_community: know_the_need_of_my_community, together_community_members_community_infrastructure: together_community_members_community_infrastructure, with_other_community_infrastructure: with_other_community_infrastructure, bring_someone_together: bring_someone_together, brought_people_together_incident: brought_people_together_incident, conflict_with_anyone_ask_position: conflict_with_anyone_ask_position, conflict_matters_interest_mine: conflict_matters_interest_mine, there_puja_at_my_house: there_puja_at_my_house, module1: module1, module2: module2, module3: module3, module4: module4, module5: module5})
	} else if p.Module == "Module4" {
		str := "SELECT firstName FROM training_participants train,SpoorthiBaselineQuestionnaire spoo where (spoo.partcipantId=train.id) and module1=1 and module2=1 and module3=1;"

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var FirstName string

		type sporti struct {
			FirstName string
		}

		var result []sporti // creating slice

		for rows.Next() {

			err := rows.Scan(&FirstName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: FirstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

		//result = append(result, SpoorthiBaselineQuestionnaire{Id: Id, partcipantId: partcipantId, email_address: email_address, GelathiId: GelathiId, Spoorthi_Session_Number: Spoorthi_Session_Number, list_down_your_skills: list_down_your_skills, skills_to_overcome_my_challenges: skills_to_overcome_my_challenges, used_skills_resources_combat_challenge: used_skills_resources_combat_challenge, listen_paragraph: listen_paragraph, summarize_main_points_paragraph: summarize_main_points_paragraph, ask_two_questions_help_you_understand: ask_two_questions_help_you_understand, three_infrastructure_of_your_village: three_infrastructure_of_your_village, know_the_need_of_my_community: know_the_need_of_my_community, together_community_members_community_infrastructure: together_community_members_community_infrastructure, with_other_community_infrastructure: with_other_community_infrastructure, bring_someone_together: bring_someone_together, brought_people_together_incident: brought_people_together_incident, conflict_with_anyone_ask_position: conflict_with_anyone_ask_position, conflict_matters_interest_mine: conflict_matters_interest_mine, there_puja_at_my_house: there_puja_at_my_house, module1: module1, module2: module2, module3: module3, module4: module4, module5: module5})
	} else if p.Module == "Module5" {

		str := "SELECT firstName FROM training_participants train,SpoorthiBaselineQuestionnaire spoo where (spoo.partcipantId=train.id) and module1=1 and module2=1 and module3=1 and module4=1;"

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var FirstName string

		type sporti struct {
			FirstName string
		}

		var result []sporti // creating slice

		for rows.Next() {

			err := rows.Scan(&FirstName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: FirstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

		//result = append(result, SpoorthiBaselineQuestionnaire{Id: Id, partcipantId: partcipantId, email_address: email_address, GelathiId: GelathiId, Spoorthi_Session_Number: Spoorthi_Session_Number, list_down_your_skills: list_down_your_skills, skills_to_overcome_my_challenges: skills_to_overcome_my_challenges, used_skills_resources_combat_challenge: used_skills_resources_combat_challenge, listen_paragraph: listen_paragraph, summarize_main_points_paragraph: summarize_main_points_paragraph, ask_two_questions_help_you_understand: ask_two_questions_help_you_understand, three_infrastructure_of_your_village: three_infrastructure_of_your_village, know_the_need_of_my_community: know_the_need_of_my_community, together_community_members_community_infrastructure: together_community_members_community_infrastructure, with_other_community_infrastructure: with_other_community_infrastructure, bring_someone_together: bring_someone_together, brought_people_together_incident: brought_people_together_incident, conflict_with_anyone_ask_position: conflict_with_anyone_ask_position, conflict_matters_interest_mine: conflict_matters_interest_mine, there_puja_at_my_house: there_puja_at_my_house, module1: module1, module2: module2, module3: module3, module4: module4, module5: module5})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "No Such Modules", "Status Code": "400 "})
	}
}

//------------------------------------------- buzz vyapar ------------------------------------------
func Buzzlist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	//var firstName string
	var p BuzzVyaparProgramBaseline
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		panic(err.Error())
	}

	type sporti struct {
		FirstName string
	}
	//var participantid int
	var firstName string
	if p.Module == "Module1" {
		str := ("SELECT train.firstName FROM training_participants train, BuzzVyaparProgramBaseline spoo where (spoo.partcipantId=train.id) and train.enroll=1;")

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

	} else if p.Module == "Module2" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.BuzzVyaparProgramBaseline spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and train.enroll=1;")
		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})
	} else if p.Module == "Module3" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.BuzzVyaparProgramBaseline spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1;")
		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all module3 "})
	} else if p.Module == "Module4" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.BuzzVyaparProgramBaseline spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1 and spoo.module3=1;")
		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all module4 "})
	} else if p.Module == "Module5" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.BuzzVyaparProgramBaseline spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1 and spoo.module3=1 and spoo.module4=1;")
		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all module5 "})
	} else {
		//json.NewEncoder(w).Encode(map[string]interface{}{"Message": " No more Modules"})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "No Such Modules", "Status Code": "400 "})
	}
}

//-------------------------------------------Green vypar survey -------------------------------------------

func Greenlist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var p BuzzVyaparProgramBaseline
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		panic(err.Error())
	}

	type sporti struct {
		FirstName string
	}

	var firstName string
	if p.Module == "Module1" {

		str := ("SELECT Distinct(train.firstName) FROM training_participants train, bdms_staff.GreenBaselineSurvey spoo where (spoo.partcipantId=train.id) and train.enroll=1;")
		//	select * from bdms_staff.training_participants where GreenMotivators=1 and gelathiRecomm=1 and enroll=1
		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})

	} else if p.Module == "Module2" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.GreenBaselineSurvey spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and train.enroll=1;")

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all "})
	} else if p.Module == "Module3" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.GreenBaselineSurvey spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1;")

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all of module3"})
	} else if p.Module == "Module4" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.GreenBaselineSurvey spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1 and spoo.module3=1;")

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all of module4"})
	} else if p.Module == "Module5" {
		str := ("SELECT firstName FROM bdms_staff.training_participants train,bdms_staff.GreenBaselineSurvey spoo where (spoo.partcipantId=train.id) and spoo.module1=1 and spoo.module2=1 and spoo.module3=1 and spoo.module4=1;")

		rows, err := DB.Query(str)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 "})
			return
		}
		defer rows.Close()
		var result []sporti // creating slice
		for rows.Next() {
			err1 := rows.Scan(&firstName)
			if err1 != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}

			result = append(result, sporti{FirstName: firstName})

		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "Message": "List all of module5"})
	} else {
		//json.NewEncoder(w).Encode(map[string]interface{}{"Message": " No more Modules"})
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "No Such Modules", "Status Code": "400 "})
	}
}

//---------------------------------------------add buzz vyapar--------------------------------------------------------
func Addbuzzvyapar(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	file, err := os.OpenFile("./logfile", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Addbuzzvyapar", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 ", "API": "Addbuzzvyapar"})
		return

	}
	log.SetOutput(file)

	var assets BuzzVyaparProgramBaseline
	err1 := json.NewDecoder(r.Body).Decode(&assets)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Message": "Invalid Input Syntax", "Status Code": "400 ", "API": "Addbuzzvyapar"})
		return
	}
	checkStatement := "SELECT COUNT(*) FROM BuzzVyaparProgramBaseline WHERE partcipantId = ?"
	var count int
	err = DB.QueryRow(checkStatement, assets.PartcipantId).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error checking participant ID:", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Internal Server Error", "Status Code": "500 ", "API": "Addbuzzvyapar"})
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Participant ID already present")
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Participant ID already present", "Status Code": "400 ", "API": "Addbuzzvyapar"})
		return
	}
	date := time.Now().Add(5*time.Hour + 30*time.Minute)
	if assets.PartcipantId != "" && assets.GfId != "" {
		//	participantid, _ := strconv.Atoi(assets.partcipantId)
		addStatement := `INSERT INTO BuzzVyaparProgramBaseline (partcipantId,GfId,entry_date,When_Was_survey_done,Name_of_the_vyapari,Age,Contact_number,Village_id,name_of_the_cohort ,Highter_education,Marital_status,Number_of_people_in_the_household,Do_you_own_a_smart_phone,Do_you_have_internet_connection_on_your_smart_phone,Sector_type_of_business,Are_you_proficient_with_numbers,Are_you_proficient_with_written_language,Household_income_monthly,Over_the_last_month_your_average_income,Your_business_profit_last_month,How_much_monthly_income_would_you_like_to_ideally_earn,Amount_invested_when_the_business_started,Number_of_years_the_business_has_been_operating,you_stopped_hold_your_business,no_hours_engaged_business,license_for_existing_business,home_based_work_from_shop,Why_do_you_do_business,Tell_us_three_things_about_you_as_an_entrepreneur, Tell_us_three_things_about_your_role_as_a_woman_at_home, What_are_your_challenges_in_running_and_growing_your_business, What_is_your_plan_to_overcome_these_challenges,What_are_your_skills,What_are_the_resources_available_with_you_for_your_business,Who_is_your_customer_Describe_them_to_us,Please_list_down_the_various_components_of_business,I_know_the_current_state_of_my_business_in_profit_loss_revenue,What_kind_of_books_of_accounts_do_you_maintain,I_can_generate_ideas_to_solve_my_business_problems,Tell_us_about_one_business_problem,What_is_your_business_goal_Business_impurumenet_madodu,Do_you_have_a_business_plan_to_reach_that_goal,Can_you_submit_a_business_plan_for_your_goal_to_us_right_now,What_are_the_strenghts_of_your_business,What_are_the_weaknesses_of_your_business,What_are_the_oppourtunities_for_your_business,Are_you_able_to_raise_the_required_finance,I_have_taken_a_loan_from,I_have_trouble_accessing_loan_for_my_business,What_are_the_prerequisites_to_access_a_loan,Loan_currently_availed,Need_additional_skills_business) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		_, err := DB.Exec(addStatement, assets.PartcipantId, assets.GfId, date, assets.When_Was_survey_done, assets.Name_of_the_vyapari, assets.Age, assets.Contact_number, assets.Village_id, assets.Name_of_the_cohort, assets.Highter_education, assets.Marital_status, assets.Number_of_people_in_the_household, assets.Do_you_own_a_smart_phone, assets.Do_you_have_internet_connection_on_your_smart_phone, assets.Sector_type_of_business, assets.Are_you_proficient_with_numbers, assets.Are_you_proficient_with_written_language, assets.Household_income_monthly, assets.Over_the_last_month_your_average_income, assets.Your_business_profit_last_month, assets.How_much_monthly_income_would_you_like_to_ideally_earn, assets.Amount_invested_when_the_business_started, assets.Number_of_years_the_business_has_been_operating, assets.You_stopped_hold_your_business, assets.No_hours_engaged_business, assets.License_for_existing_business, assets.Home_based_work_from_shop, assets.Why_do_you_do_business, strings.Join(assets.Tell_us_three_things_about_you_as_an_entrepreneur, ", "), assets.Tell_us_three_things_about_your_role_as_a_woman_at_home, assets.What_are_your_challenges_in_running_and_growing_your_business, assets.What_is_your_plan_to_overcome_these_challenges, assets.What_are_your_skills, assets.What_are_the_resources_available_with_you_for_your_business, assets.Who_is_your_customer_Describe_them_to_us, strings.Join(assets.Please_list_down_the_various_components_of_business, ", "), assets.I_know_the_current_state_of_my_business_in_profit_loss_revenue, assets.What_kind_of_books_of_accounts_do_you_maintain, assets.I_can_generate_ideas_to_solve_my_business_problems, assets.Tell_us_about_one_business_problem, assets.What_is_your_business_goal_Business_impurumenet_madodu, assets.Do_you_have_a_business_plan_to_reach_that_goal, assets.Can_you_submit_a_business_plan_for_your_goal_to_us_right_now, assets.What_are_the_strenghts_of_your_business, assets.What_are_the_weaknesses_of_your_business, assets.What_are_the_oppourtunities_for_your_business, assets.Are_you_able_to_raise_the_required_finance, assets.I_have_taken_a_loan_from, assets.I_have_trouble_accessing_loan_for_my_business, strings.Join(assets.What_are_the_prerequisites_to_access_a_loan, ", "), assets.Loan_currently_availed, assets.Need_additional_skills_business)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid input syntax for IP ", "Status Code": "400 ", "Error": err, "API": "Addbuzzvyapar"})
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Println("Recorded sucessfully")
		json.NewEncoder(w).Encode(map[string]interface{}{"Status Code": "200 OK", "Message": "Recorded sucessfully"})
	} else {
		w.WriteHeader(http.StatusOK)
		log.Println("Paricipantid and GFId is missing")
		json.NewEncoder(w).Encode(map[string]interface{}{"Status Code": "200 OK", "Message": "Paricipantid and GFId is missing "})
	}
}
