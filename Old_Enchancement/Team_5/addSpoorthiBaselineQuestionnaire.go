package Team_5

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type addSpoorthiBaselineQuestionnaire struct {
	ParticipantID                          string `json:"partcipantId"`
	GelathiID                              string `json:"GelathiId"`
	Email                                  string `json:"email_address"`
	SpoorthiSessionNumber                  string `json:"Spoorthi_Session_Number"`
	ListDownYourSkills                     string `json:"list_down_your_skills"`
	SkillsToOvercomeMyChallenges           string `json:"skills_to_overcome_my_challenges"`
	UsedSkillsResourcesCombatChallenge     string `json:"used_skills_resources_combat_challenge"`
	ListenParagraph                        string `json:"listen_paragraph"`
	SummarizeMainPointsParagraph           string `json:"summarize_main_points_paragraph"`
	AskTwoQuestionsHelpYouUnderstand       string `json:"ask_two_questions_help_you_understand"`
	ThreeInfrastructureOfYourVillage       string `json:"three_infrastructure_of_your_village"`
	KnowTheNeedOfMyCommunity               string `json:"know_the_need_of_my_community"`
	TogetherCommunityMembersCommunityInfra string `json:"together_community_members_community_infrastructure"`
	WithOtherCommunityInfrastructure       string `json:"with_other_community_infrastructure"`
	BringSomeoneTogether                   string `json:"bring_someone_together"`
	BroughtPeopleTogetherIncident          string `json:"brought_people_together_incident"`
	ConflictWithAnyoneAskPosition          string `json:"conflict_with_anyone_ask_position"`
	ConflictMattersInterestMine            string `json:"conflict_matters_interest_mine"`
	TherePujaAtMyHouse                     string `json:"There_puja_at_my_house"`
}

type Response3 struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func AddSpoorthiBaselineQuestionnaire(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	// Set response headers
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := Response3{Success: false, Message: "Method Not Allowed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var request addSpoorthiBaselineQuestionnaire
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response3{Success: false, Message: "Invalid request body"}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()
	checkStatement := "SELECT COUNT(*) FROM SpoorthiBaselineQuestionnaire WHERE partcipantId = ?"
	var count int
	err = DB.QueryRow(checkStatement, request.ParticipantID).Scan(&count)
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

	if request.ParticipantID != "" && request.GelathiID != "" {
		date := time.Now().Add(5*time.Hour + 30*time.Minute)

		addStatement := `INSERT INTO SpoorthiBaselineQuestionnaire (
			partcipantId,
			email_address,
			GelathiId,
			entry_date,
			Spoorthi_Session_Number,
			list_down_your_skills,
			skills_to_overcome_my_challenges,
			used_skills_resources_combat_challenge,
			listen_paragraph,
			summarize_main_points_paragraph,
			ask_two_questions_help_you_understand,
			three_infrastructure_of_your_village,
			know_the_need_of_my_community,
			together_community_members_community_infrastructure,
			with_other_community_infrastructure,
			bring_someone_together,
			brought_people_together_incident,
			conflict_with_anyone_ask_position,
			conflict_matters_interest_mine,
			There_puja_at_my_house
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

		_, err := DB.Exec(addStatement,
			request.ParticipantID,
			request.Email,
			request.GelathiID,
			date,
			request.SpoorthiSessionNumber,
			request.ListDownYourSkills,
			request.SkillsToOvercomeMyChallenges,
			request.UsedSkillsResourcesCombatChallenge,
			request.ListenParagraph,
			request.SummarizeMainPointsParagraph,
			request.AskTwoQuestionsHelpYouUnderstand,
			request.ThreeInfrastructureOfYourVillage,
			request.KnowTheNeedOfMyCommunity,
			request.TogetherCommunityMembersCommunityInfra,
			request.WithOtherCommunityInfrastructure,
			request.BringSomeoneTogether,
			request.BroughtPeopleTogetherIncident,
			request.ConflictWithAnyoneAskPosition,
			request.ConflictMattersInterestMine,
			request.TherePujaAtMyHouse,
		)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			response := Response3{Success: false, Message: "Failed to add questionnaire"}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := Response3{Success: true, Message: "Added Successfully"}
		json.NewEncoder(w).Encode(response)
	} else {
		response := Response3{Success: false, Message: "Send Required Parameters"}
		json.NewEncoder(w).Encode(response)
	}
}
