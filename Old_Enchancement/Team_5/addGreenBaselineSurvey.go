package Team_5

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type GreenBaselineSurvey struct {
	PartcipantID                                       string   `json:"partcipantId"`
	Email                                              string   `json:"Email"`
	NameOfTheSurveyor                                  string   `json:"Name_of_the_surveyor"`
	NameOfTheRespondent                                string   `json:"Name_of_the_respondent"`
	VillageName                                        string   `json:"Village_Name"`
	PhoneNumber                                        string   `json:"Phone_number"`
	NaturalResources                                   []string `json:"natural_resources"`
	NaturalResourcesImpactingYourLife                  []string `json:"natural_resources_impacting_your_life"`
	NaturalWealth                                      string   `json:"Natural_Wealth"`
	ClimateChange                                      string   `json:"climate_change"`
	WhatDoYouKnowAboutIt                               string   `json:"What_do_you_know_about_it"`
	ChangeInTheWeatherClimate                          string   `json:"change_in_the_weather_climate"`
	ChangesHappenedToTheClimate                        []string `json:"changes_happened_to_the_climate"`
	ClimateChangeThreatensPersonalFamilyHealthSafety   string   `json:"climate_change_threatens_personal_family_health_safety"`
	DoneToTackleClimateChange                          string   `json:"done_to_tackle_climate_change"`
	DoSomethingToTackleClimateChange                   string   `json:"do_something_to_tackle_climate_change"`
	MainSourceOfWater                                  []string `json:"main_source_of_water"`
	ShownBelowDoYouAgreeWith                           string   `json:"shown_below_do_you_agree_with"`
	HowConcernedLocalWaterQuality                      string   `json:"How_concerned_local_water_quality"`
	PersonalActionsCanAffectWaterQuality               string   `json:"personal_actions_can_affect_water_quality"`
	TakeWaterConservationMeasures                      string   `json:"take_water_conservation_measures"`
	IfYesWhatKindOfMeasures                            string   `json:"If_yes_what_kind_of_measures"`
	ListDownImpactOfClimateChange                      string   `json:"list_down_impact_of_climate_change"`
	GiftingHis6AcresLandToChildren                     string   `json:"gifting_his_6_acres_land_to_children"`
	WhichOneAccordingToYouIsRight                      string   `json:"Which_one_according_to_you_is_right"`
	BelieveConnectionBetweenTheFoodHealthClimate       string   `json:"believe_connection_between_the_food_health_climate"`
	NativeFoodYouBelieveIsEnvironmentallyFriendly      string   `json:"native_food_you_believe_is_envionmentally_friendly"`
	HouseholdActivityPollutesNaturalResources          string   `json:"household_activity_pollutes_natural_resources"`
	AlternativesHouseholdMaterialsCausePollution       string   `json:"alternatives_household_materials_cause_pollution"`
	IfYesWhatAreThey                                   string   `json:"If_yes_what_are_they"`
	EcoFriendlyProductsAndActivities                   string   `json:"eco_friendly_products_and_activities"`
	LittleMoreThanWhatYouPayForTheChemicals            string   `json:"little_more_than_what_you_pay_for_the_chemicals"`
	ThisSwitchToEcoFriendlyProducts                    string   `json:"this_switch_to_eco_friendly_products"`
	ClimateChangeIsALotOfEffort                        string   `json:"climate_change_is_a_lot_of_effort"`
	ActionOutOfConcernForClimateChange                 string   `json:"action_out_of_concern_for_climate_change"`
	IfYesWhatDidYouDoAreYouDoing                       string   `json:"If_yes_what_did_you_do_are_you_doing"`
	NaturalResourceCommunityImmediateAttentionMeasures string   `json:"natural_resource_community_immediate_attention_measures"`
	IfYesWhatIsThatResource                            string   `json:"If_yes_what_is_that_resource"`
	AchieveWithRegardToNaturalResourceConservation     string   `json:"achieve_with_regard_to_natural_resource_conservation"`
	InitiativeToConserveTheEnvironment                 string   `json:"initiative_to_conserve_the_environment"`
	CommunityTogetherAchieveMyConservationGoal         string   `json:"community_together_achieve_my_conservation_goal"`
}
type Response1 struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func AddGreenBaselineSurvey(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Set response headers
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := Response1{Success: false, Message: "Method Not Allowed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	var request GreenBaselineSurvey
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response1{Success: false, Message: "Invalid request body"}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()

	checkStatement := "SELECT COUNT(*) FROM GreenBaselineSurvey WHERE partcipantId = ?"
	var count int
	err = DB.QueryRow(checkStatement, request.PartcipantID).Scan(&count)
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

	if request.PartcipantID != "" {
		//date := time.Now().Format("2006/01/02")
		date := time.Now().Add(5*time.Hour + 30*time.Minute)
		stmt := `INSERT INTO GreenBaselineSurvey (
			partcipantId,
			Email,
			Name_of_the_surveyor,
			Name_of_the_respondent,
			Village_Name,
			Phone_number,
			natural_resources,
			natural_resources_impacting_your_life,
			Natural_Wealth,
			climate_change,
			What_do_you_know_about_it,
			change_in_the_weather_climate,
			changes_happened_to_the_climate,
			climate_change_threatens_personal_family_health_safety,
			done_to_tackle_climate_change,
			do_something_to_tackle_climate_change,
			main_source_of_water,
			shown_below_do_you_agree_with,
			How_concerned_local_water_quality,
			personal_actions_can_affect_water_quality,
			take_water_conservation_measures,
			If_yes_what_kind_of_measures,
			list_down_impact_of_climate_change,
			gifting_his_6_acres_land_to_children,
			Which_one_according_to_you_is_right,
			believe_connection_between_the_food_health_climate,
			native_food_you_believe_is_envionmentally_friendly,
			household_activity_pollutes_natural_resources,
			alternatives_household_materials_cause_pollution,
			If_yes_what_are_they,
			eco_friendly_products_and_activities,
			little_more_than_what_you_pay_for_the_chemicals,
			this_switch_to_eco_friendly_products,
			climate_change_is_a_lot_of_effort,
			action_out_of_concern_for_climate_change,
			If_yes_what_did_you_do_are_you_doing,
			natural_resource_community_immediate_attention_measures,
			If_yes_what_is_that_resource,
			achieve_with_regard_to_natural_resource_conservation,
			initiative_to_conserve_the_environment,
			community_together_achieve_my_conservation_goal,
			entry_date
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)`

		_, err = DB.Exec(stmt,
			request.PartcipantID,
			request.Email,
			request.NameOfTheSurveyor,
			request.NameOfTheRespondent,
			request.VillageName,
			request.PhoneNumber,
			strings.Join(request.NaturalResources, ","),
			strings.Join(request.NaturalResourcesImpactingYourLife, ","),
			request.NaturalWealth,
			request.ClimateChange,
			request.WhatDoYouKnowAboutIt,
			request.ChangeInTheWeatherClimate,
			strings.Join(request.ChangesHappenedToTheClimate, ","),
			request.ClimateChangeThreatensPersonalFamilyHealthSafety,
			request.DoneToTackleClimateChange,
			request.DoSomethingToTackleClimateChange,
			strings.Join(request.MainSourceOfWater, ","),
			request.ShownBelowDoYouAgreeWith,
			request.HowConcernedLocalWaterQuality,
			request.PersonalActionsCanAffectWaterQuality,
			request.TakeWaterConservationMeasures,
			request.IfYesWhatKindOfMeasures,
			request.ListDownImpactOfClimateChange,
			request.GiftingHis6AcresLandToChildren,
			request.WhichOneAccordingToYouIsRight,
			request.BelieveConnectionBetweenTheFoodHealthClimate,
			request.NativeFoodYouBelieveIsEnvironmentallyFriendly,
			request.HouseholdActivityPollutesNaturalResources,
			request.AlternativesHouseholdMaterialsCausePollution,
			request.IfYesWhatAreThey,
			request.EcoFriendlyProductsAndActivities,
			request.LittleMoreThanWhatYouPayForTheChemicals,
			request.ThisSwitchToEcoFriendlyProducts,
			request.ClimateChangeIsALotOfEffort,
			request.ActionOutOfConcernForClimateChange,
			request.IfYesWhatDidYouDoAreYouDoing,
			request.NaturalResourceCommunityImmediateAttentionMeasures,
			request.IfYesWhatIsThatResource,
			request.AchieveWithRegardToNaturalResourceConservation,
			request.InitiativeToConserveTheEnvironment,
			request.CommunityTogetherAchieveMyConservationGoal,
			date,
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			response := Response1{Success: false, Message: "Failed to add questionnaire"}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := Response1{Success: true, Message: "Added Successfully"}
		json.NewEncoder(w).Encode(response)
	} else {
		response := Response1{Success: false, Message: "Send Required Parameters"}
		json.NewEncoder(w).Encode(response)
	}
}
