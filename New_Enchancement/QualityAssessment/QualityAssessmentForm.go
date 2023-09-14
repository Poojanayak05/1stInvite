package QualityAssessment

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to hold data from the table
type QualityAssessmentForm_1 struct {
	ID                                                         int       `json:"id"`
	EmpID                                                      int       `json:"emp_id"`
	RoleID                                                     int       `json:"role_id"`
	EmailAddress                                               string    `json:"email_address"`
	NameOfTheAssessor                                          string    `json:"name_of_the_assessor"`
	EntryDate                                                  time.Time `json:"entry_date"`
	ProgramAssessment                                          int       `json:"program_assessment"`
	TodayPOA                                                   string    `json:"today_poa"`
	NameOfTheDistrict                                          string    `json:"name_of_the_district"`
	NameOfTheTaluk                                             string    `json:"name_of_the_taluk"`
	NameOfTheVillageAndVenueOfMeetingOrTraining                string    `json:"name_of_the_village_and_the_venue_of_meeting_or_training"`
	Day1OrDay2                                                 string    `json:"day1_or_day2"`
	NameOfTheTrainerBeingEvaluated                             string    `json:"name_of_the_trainer_being_evaluated"`
	CheckWhichOnesTheTrainerDidNotDo                           []string  `json:"check_which_ones_the_trainer_did_not_do"`
	HowManyWomenAttendedTheTrainingSession                     int       `json:"how_many_women_attended_the_training_session"`
	CheckWhichOnesTheTrainerDidNotDo1                          []string  `json:"check_which_ones_the_trainer_did_not_do_1"`
	WereTheWomenInteractive                                    string    `json:"were_the_women_interactive"`
	DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule      string    `json:"did_any_women_leave_tring_session_dring_or_after_1st_module"`
	IfSoHowMany                                                int       `json:"if_so_how_many"`
	DidThisModuleTake20MinutesAsAllotted                       string    `json:"did_this_module_take_20_minutes_as_allotted"`
	DidAnyNewWomenAttendTheTrainingSessionDuringModule         string    `json:"did_any_new_women_attend_the_training_session_during_module"`
	IfSoHowMany1                                               int       `json:"if_so_how_many_1"`
	CheckWhichOnesTheTrainerDidNotDo2                          []string  `json:"check_which_ones_the_trainer_did_not_do_2"`
	DuringTheDebriefTheTrainerDid                              []string  `json:"during_the_debrief_the_trainer_did"`
	DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule1     string    `json:"did_any_women_leve_training_session_during_or_after_1st_module_1"`
	IfSoHowMany2                                               int       `json:"if_so_how_many_2"`
	DidThisModuleTake20MinutesAsAllotted1                      string    `json:"did_this_module_take_20_minutes_as_allotted_1"`
	DidAnyNewWomenAttendTrainingSessionDuringThisModule1       string    `json:"did_any_new_women_attend_training_session_during_this_module_1"`
	IfSoHowMany3                                               int       `json:"if_so_how_many_3"`
	CheckWhichOnesTheTrainerDidNotDo3                          []string  `json:"check_which_ones_the_trainer_did_not_do_3"`
	DuringTheDebriefsForRolePlaysTheTrainerDidNotAsk           []string  `json:"during_the_debriefs_for_role_plays_the_trainer_did_not_ask"`
	DidTheTrainerLeaveWomenToReadRolePlayCardThemselves        string    `json:"did_the_trainer_leave_women_to_read_role_play_card_themselves"`
	DidTheGroupsEngageAndInteractAmongThemselvesWell           string    `json:"did_the_groups_engage_and_interact_among_themselves_well"`
	WereTheParticipantsResponsiveDuringTheDebriefing           string    `json:"were_the_participants_responsive_during_the_debriefing"`
	DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule2 string    `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_2"`
	IfSoHowMany4                                               int       `json:"if_so_how_many_4"`
	CheckWhichOnesTheTrainerDidNotDo4                          []string  `json:"check_which_ones_the_trainer_did_not_do_4"`
	DidThisModuleTake30MinutesAsAllotted1                      string    `json:"did_this_module_take_30_minutes_as_allotted_1"`
	HowManyWomenRemainedByTheEndOfThisTrainingSession          int       `json:"how_many_women_remained_by_the_end_of_this_training_session"`
	HowManyAreLikelyToComeBack                                 int       `json:"how_many_are_likely_to_come_back"`
	DidAnyNewWomenAttendTrainingSessionDuringThisModule2       string    `json:"did_any_new_women_attend_training_session_during_this_module_2"`
	IfSoHowMany5                                               int       `json:"if_so_how_many_5"`
	DidThisModuleTake30MinutesAsAllotted                       string    `json:"did_this_module_take_30_minutes_as_allotted"`
	WasTheRecapDone                                            string    `json:"was_the_recap_done"`
	DidTheRecapTake15MinutesAsAllotted                         string    `json:"did_the_recap_take_15_minutes_as_allotted"`
	NameOfTheGF                                                string    `json:"name_of_the_gf"`
	NoOfParticipantsAtTheStartOfTheSession                     int       `json:"no_of_participants_at_the_start_of_the_session"`
	AssessmentOf                                               string    `json:"assessment_of"`
	TheGFCompletelyCarriedOutFollowingFunctions                []string  `json:"the_gf_comptetly_carried_out_following_funtions"`
	TheGFCariedFollowigFuctionsBforeTraningOrMetingStarted     []string  `json:"the_gf_caried_followig_fuctions_bfore_traning_or_meting_started"`
	HowManyStoriesOfSuccessOrChangeEmergedFromTheRecap         int       `json:"how_many_stories_of_success_or_change_emerged_from_the_recap"`
	MentNameOfGelathisSuccessStoriesAndStoryCoupleOfLines      string    `json:"ment_name_of_gelathis_success_stories_and_story_couple_of_lines"`
	CheckWhichOnesTheGFDidNotDo                                []string  `json:"check_which_ones_the_gf_did_not_do"`
	NumberOfEnrolledGelathisInTheCircle                        int       `json:"number_of_enrolled_gelathis_in_the_circle"`
	NoOfAttendedGelathis                                       int       `json:"no_of_attended_gelathis"`
	LevelOfParticipation1                                      int       `json:"level_of_participation_1"`
	LevelOfParticipation2                                      int       `json:"level_of_participation_2"`
	LevelOfParticipation3                                      int       `json:"level_of_participation_3"`
	LevelOfParticipation4                                      int       `json:"level_of_participation_4"`
	LevelOfParticipation5                                      int       `json:"level_of_participation_5"`
	TheGFCompetentlyCoveredFollowingThingsInTrainingDelivered  []string  `json:"the_gf_competently_covered_folowing_things_in_training_delivered"`
	RateTheGF1                                                 int       `json:"rate_the_gf_1"`
	RateTheGF2                                                 int       `json:"rate_the_gf_2"`
	RateTheGF3                                                 int       `json:"rate_the_gf_3"`
	RateTheGF4                                                 int       `json:"rate_the_gf_4"`
	RateTheGF5                                                 int       `json:"rate_the_gf_5"`
	RateTheGF6                                                 int       `json:"rate_the_gf_6"`
	RateTheGF7                                                 int       `json:"rate_the_gf_7"`
	RateTheGF8                                                 int       `json:"rate_the_gf_8"`
	RateTheGF9                                                 int       `json:"rate_the_gf_9"`
	RateTheGF10                                                int       `json:"rate_the_gf_10"`
	WhatWorkedInTheTraining                                    string    `json:"what_worked_in_the_training"`
	WhatCanBeBetterNextTime                                    string    `json:"what_can_be_better_next_time"`
	AnyFurtherTrainingAndUnderstandingReqiredByGFTraingModule  []string  `json:"any_futher_training_and_understding_reqired_by_gf_traing_module"`
	DidYouFindAnythingTraiingOrGFThatNeedsToWorkedPriority     string    `json:"did_you_find_anything_traiing_or_gf_that_neds_to_worked_priority"`
	DetailsOfSuccessStoriesToBeCollectedFromGelathisByGF       string    `json:"details_of_success_stories_to_be_collected_from_gelathis_by_gf"`
	DeadlineToCollectTheStories                                time.Time `json:"deadline_to_collect_the_stories"`
	EndTimeOfTheTraining                                       string    `json:"end_time_of_the_training"`
	NoOfParticipantsAtEndOfTheSession                          int       `json:"no_of_participants_at_end_of_the_session"`
	AnyOtherCommentsAboutTheGelathiFacilitator                 string    `json:"any_other_comments_about_the_gelathi_facilitator"`
	NameOfTheGelathiBeingEvaluated                             string    `json:"name_of_the_gelathi_being_evaluated"`
	DaysModules                                                string    `json:"days_modules"`
	CheckWhichOnesTheGelathiDidNotDo                           []string  `json:"check_which_ones_the_gelathi_did_not_do"`
	CheckWhichOnesTheGelathiDidNotDo1                          []string  `json:"check_which_ones_the_gelathi_did_not_do_1"`
	WasTheRecapDone1                                           string    `json:"was_the_recap_done_1"`
	DidTheDebriefDoneByGelathi                                 string    `json:"did_the_debrief_done_by_gelathi"`
	DuringTheDebriefsForRolePlaysTheGelathiDidNotAsk           []string  `json:"during_the_debriefs_for_role_plays_the_gelathi_did_not_ask"`
	RepeatTheActivityWithTheSecondVolunteer                    string    `json:"repeat_the_activity_with_the_second_volunteer"`
	DuringTheDebriefDidTheGelathiNotAsk                        []string  `json:"during_the_debrief_did_the_gelathi_not_ask"`
	TheGelathiDidNotAsk1                                       []string  `json:"the_gelathi_did_not_ask_1"`
	CheckWhichOnesTheGelathiDidNotDo2                          []string  `json:"check_which_ones_the_gelathi_did_not_do_2"`
	Filter                                                     string    `json:"filter"`
	Filter_District                                            string    `json:"filter_district"`
	Filter_Taluk                                               string    `json:"filter_taluk"`
	Filter_Id                                                  string    `json:"filter_id"`
	Filter_StartDate                                           string    `json:"filter_StartDate"`
	Filter_EndDate                                             string    `json:"filter_EndDate"`
	CheckWhichInstTrainerDidNotDo                              []string  `json:"check_which_instructions_the_trainer_did_not_do"`
	DuringDebriefDidTheTrainerNotAsk                           []string  `json:"during_the_debrief_did_the_trainer_not_ask"`
	WereTheParticipantsResponsiveDuringTheDebriefing1          string    `json:"were_the_participants_responsive_during_the_debriefing_1"`
	DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule1 string    `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_1"`
	HowManyWomenAttendedTheTrainingSession3                    int       `json:"how_many_women_attended_the_training_session_3"`
	TheTrainerDidNotAsk                                        []string  `json:"the_trainer_did_not_ask"`
	WhatDidTheTrainerNotDo                                     []string  `json:"what_did_the_trainer_not_do"`
	DuringDebriefDidTheTrainerNotAsk1                          []string  `json:"during_the_debrief_did_the_trainer_not_ask_1"`
	WereTheParticipantsResponsiveDuringTheDebriefing2          string    `json:"were_the_participants_responsive_during_the_debriefing_2"`
	DidThisModuleTake30MinutesAsAllotted2                      string    `json:"did_this_module_take_30_minutes_as_allotted_2"`
	HowManyWomenAttendedTheTrainingSession4                    int       `json:"how_many_women_attended_the_training_session_4"`
	CheckWhichOnesTheTrainerDidNotDo5                          []string  `json:"check_which_ones_the_trainer_did_not_do_5"`
	DuringDebriefDidTheTrainerNotAsk2                          []string  `json:"during_the_debrief_did_the_trainer_not_ask_2"`
	WereTheParticipantsResponsiveDuringTheDebriefing3          string    `json:"were_the_participants_responsive_during_the_debriefing_3"`
	DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule3     string    `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_3"`
	DidThisModuleTake30MinutesAsAllotted3                      string    `json:"did_this_module_take_30_minutes_as_allotted_3"`
	DuringTheDebriefDidTheGelathi                              []string  `json:"during_the_debrief_did_the_gelathi"`
	NameOfTheGelathiBeingEvaluated1                            string    `json:"name_of_the_gelathi_being_evaluated_1"`
	CheckWhichInstTheGelathiDidNotDo                           []string  `json:"check_which_instructions_the_gelathi_did_not_do"`
	DuringTheDebriefDidTheTrainerDidNotDoTheFollowing          []string  `json:"during_the_debrief_did_the_trainer_did_not_do_the_following"`
	DidAnyWomenLeaveTheTrainingSessionDuringOnAfterThisModule4 string    `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_4"`
	HowManyWomenAttendedTheTrainingSession1                    int       `json:"how_many_women_attended_the_training_session_1"`
	HowManyWomenAttendedTheTrainingSession2                    int       `json:"how_many_women_attended_the_training_session_2"`
	CheckWhichOnesTheTrainerDidNotDo6                          []string  `json:"check_which_ones_the_trainer_did_not_do_6"`
}
type ListQualityAssessmentForm_1 struct {
	Id                                                               int    `json:"id"`
	Emp_id                                                           int    `json:"emp_id"`
	Email_address                                                    string `json:"email_address"`
	Name_of_the_assessor                                             string `json:"name_of_the_assessor"`
	Entry_date                                                       string `json:"entry_date"`
	Program_assessment                                               int    `json:"program_assessment"`
	Today_poa                                                        string `json:"today_poa"`
	Name_of_the_district                                             string `json:"name_of_the_district"`
	Name_of_the_taluk                                                string `json:"name_of_the_taluk"`
	Name_of_the_village_and_the_venue_of_meeting_or_training         string `json:"name_of_the_village_and_the_venue_of_meeting_or_training"`
	Day1_or_day2                                                     string `json:"day1_or_day2"`
	Name_of_the_trainer_being_evaluated                              string `json:"name_of_the_trainer_being_evaluated"`
	Check_which_ones_the_trainer_did_not_do                          string `json:"check_which_ones_the_trainer_did_not_do"`
	How_many_women_attended_the_training_session                     int    `json:"how_many_women_attended_the_training_session"`
	Check_which_ones_the_trainer_did_not_do_1                        string `json:"check_which_ones_the_trainer_did_not_do_1"`
	Were_the_women_interactive                                       string `json:"were_the_women_interactive"`
	Did_any_women_leave_tring_session_dring_or_after_1st_module      string `json:"did_any_women_leave_tring_session_dring_or_after_1st_module"`
	If_so_how_many                                                   int    `json:"if_so_how_many"`
	Did_this_module_take_20_minutes_as_allotted                      string `json:"did_this_module_take_20_minutes_as_allotted"`
	Did_any_new_women_attend_the_training_session_during_module      string `json:"did_any_new_women_attend_the_training_session_during_module"`
	If_so_how_many_1                                                 int    `json:"if_so_how_many_1"`
	Check_which_ones_the_trainer_did_not_do_2                        string `json:"check_which_ones_the_trainer_did_not_do_2"`
	During_the_debrief_the_trainer_did                               string `json:"during_the_debrief_the_trainer_did"`
	Did_any_women_leve_training_session_during_or_after_1st_module_1 string `json:"did_any_women_leve_training_session_during_or_after_1st_module_1"`
	If_so_how_many_2                                                 int    `json:"if_so_how_many_2"`
	Did_this_module_take_20_minutes_as_allotted_1                    string `json:"did_this_module_take_20_minutes_as_allotted_1"`
	Did_any_new_women_attend_training_session_during_this_module_1   string `json:"did_any_new_women_attend_training_session_during_this_module_1"`
	If_so_how_many_3                                                 int    `json:"if_so_how_many_3"`
	Check_which_ones_the_trainer_did_not_do_3                        string `json:"check_which_ones_the_trainer_did_not_do_3"`
	During_the_debriefs_for_role_plays_the_trainer_did_not_ask       string `json:"during_the_debriefs_for_role_plays_the_trainer_did_not_ask"`
	Did_the_trainer_leave_women_to_read_role_play_card_themselves    string `json:"did_the_trainer_leave_women_to_read_role_play_card_themselves"`
	Did_the_groups_engage_and_interact_among_themselves_well         string `json:"did_the_groups_engage_and_interact_among_themselves_well"`
	Were_the_participants_responsive_during_the_debriefing           string `json:"were_the_participants_responsive_during_the_debriefing"`
	Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_2  string `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_2"`
	If_so_how_many_4                                                 int    `json:"if_so_how_many_4"`
	Check_which_ones_the_trainer_did_not_do_4                        string `json:"check_which_ones_the_trainer_did_not_do_4"`
	Did_this_module_take_30_minutes_as_allotted_1                    string `json:"did_this_module_take_30_minutes_as_allotted_1"`
	How_many_women_remained_by_the_end_of_this_training_session      int    `json:"how_many_women_remained_by_the_end_of_this_training_session"`
	How_many_are_likely_to_come_back                                 int    `json:"how_many_are_likely_to_come_back"`
	Did_any_new_women_attend_training_session_during_this_module_2   string `json:"did_any_new_women_attend_training_session_during_this_module_2"`
	If_so_how_many_5                                                 int    `json:"if_so_how_many_5"`
	Did_this_module_take_30_minutes_as_allotted                      string `json:"did_this_module_take_30_minutes_as_allotted"`
	Was_the_recap_done                                               string `json:"was_the_recap_done"`
	Did_the_recap_take_15_minutes_as_allotted                        string `json:"did_the_recap_take_15_minutes_as_allotted"`
	Name_of_the_gf                                                   string `json:"name_of_the_gf"`
	No_of_participants_at_the_start_of_the_session                   int    `json:"no_of_participants_at_the_start_of_the_session"`
	Assessment_of                                                    string `json:"assessment_of"`
	The_gf_comptetly_carried_out_following_funtions                  string `json:"the_gf_comptetly_carried_out_following_funtions"`
	The_gf_caried_followig_fuctions_bfore_traning_or_meting_started  string `json:"the_gf_caried_followig_fuctions_bfore_traning_or_meting_started"`
	How_many_stories_of_success_or_change_emerged_from_the_recap     int    `json:"how_many_stories_of_success_or_change_emerged_from_the_recap"`
	Ment_name_of_gelathis_success_stories_and_story_couple_of_lines  string `json:"ment_name_of_gelathis_success_stories_and_story_couple_of_lines"`
	Check_which_ones_the_gf_did_not_do                               string `json:"check_which_ones_the_gf_did_not_do"`
	Number_of_enrolled_gelathis_in_the_circle                        int    `json:"number_of_enrolled_gelathis_in_the_circle"`
	No_of_attended_gelathis                                          int    `json:"no_of_attended_gelathis"`
	Level_of_participation_1                                         int    `json:"level_of_participation_1"`
	Level_of_participation_2                                         int    `json:"level_of_participation_2"`
	Level_of_participation_3                                         int    `json:"level_of_participation_3"`
	Level_of_participation_4                                         int    `json:"level_of_participation_4"`
	Level_of_participation_5                                         int    `json:"level_of_participation_5"`
	The_gf_competently_covered_folowing_things_in_training_delivered string `json:"the_gf_competently_covered_folowing_things_in_training_delivered"`
	Rate_the_gf_1                                                    int    `json:"rate_the_gf_1"`
	Rate_the_gf_2                                                    int    `json:"rate_the_gf_2"`
	Rate_the_gf_3                                                    int    `json:"rate_the_gf_3"`
	Rate_the_gf_4                                                    int    `json:"rate_the_gf_4"`
	Rate_the_gf_5                                                    int    `json:"rate_the_gf_5"`
	Rate_the_gf_6                                                    int    `json:"rate_the_gf_6"`
	Rate_the_gf_7                                                    int    `json:"rate_the_gf_7"`
	Rate_the_gf_8                                                    int    `json:"rate_the_gf_8"`
	Rate_the_gf_9                                                    int    `json:"rate_the_gf_9"`
	Rate_the_gf_10                                                   int    `json:"rate_the_gf_10"`
	What_worked_in_the_training                                      string `json:"what_worked_in_the_training"`
	What_can_be_better_next_time                                     string `json:"what_can_be_better_next_time"`
	Any_futher_training_and_understding_reqired_by_gf                string `json:"any_futher_training_and_understding_reqired_by_gf_traing_module"`
	Did_you_find_anything_traiing_or_gf_that_neds_to_worked_priority string `json:"did_you_find_anything_traiing_or_gf_that_neds_to_worked_priority"`
	Details_of_success_stories_to_be_collected_from_gelathis_by_gf   string `json:"details_of_success_stories_to_be_collected_from_gelathis_by_gf"`
	Deadline_to_collect_the_stories                                  string `json:"deadline_to_collect_the_stories"`
	End_time_of_the_training                                         string `json:"end_time_of_the_training"`
	No_of_participants_at_end_of_the_session                         int    `json:"no_of_participants_at_end_of_the_session"`
	Any_other_comments_about_the_gelathi_facilitator                 string `json:"any_other_comments_about_the_gelathi_facilitator"`
	Name_of_the_gelathi_being_evaluated                              string `json:"name_of_the_gelathi_being_evaluated"`
	Days_modules                                                     string `json:"days_modules"`
	Check_which_ones_the_gelathi_did_not_do                          string `json:"check_which_ones_the_gelathi_did_not_do"`
	Check_which_ones_the_gelathi_did_not_do_1                        string `json:"check_which_ones_the_gelathi_did_not_do_1"`
	Was_the_recap_done_1                                             string `json:"was_the_recap_done_1"`
	Did_the_debrief_done_by_gelathi                                  string `json:"did_the_debrief_done_by_gelathi"`
	During_the_debriefs_for_role_plays_the_gelathi_did_not_ask       string `json:"during_the_debriefs_for_role_plays_the_gelathi_did_not_ask"`
	Repeat_the_activity_with_the_second_volunteer                    string `json:"repeat_the_activity_with_the_second_volunteer"`
	During_the_debrief_did_the_gelathi_not_ask                       string `json:"during_the_debrief_did_the_gelathi_not_ask"`
	The_gelathi_did_not_ask_1                                        string `json:"the_gelathi_did_not_ask_1"`
	Check_which_ones_the_gelathi_did_not_do_2                        string `json:"check_which_ones_the_gelathi_did_not_do_2"`
	Check_which_instructions_the_trainer_did_not_do                  string `json:"check_which_instructions_the_trainer_did_not_do"`
	During_the_debrief_did_the_trainer_not_ask                       string `json:"during_the_debrief_did_the_trainer_not_ask"`
	Were_the_participants_responsive_during_the_debriefing_1         string `json:"were_the_participants_responsive_during_the_debriefing_1"`
	Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_1  string `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_1"`
	How_many_women_attended_the_training_session_3                   int    `json:"how_many_women_attended_the_training_session_3"`
	The_trainer_did_not_ask                                          string `json:"the_trainer_did_not_ask"`
	What_did_the_trainer_not_do                                      string `json:"what_did_the_trainer_not_do"`
	During_the_debrief_did_the_trainer_not_ask_1                     string `json:"during_the_debrief_did_the_trainer_not_ask_1"`
	Were_the_participants_responsive_during_the_debriefing_2         string `json:"were_the_participants_responsive_during_the_debriefing_2"`
	Did_this_module_take_30_minutes_as_allotted_2                    string `json:"did_this_module_take_30_minutes_as_allotted_2"`
	How_many_women_attended_the_training_session_4                   int    `json:"how_many_women_attended_the_training_session_4"`
	Check_which_ones_the_trainer_did_not_do_5                        string `json:"check_which_ones_the_trainer_did_not_do_5"`
	During_the_debrief_did_the_trainer_not_ask_2                     string `json:"during_the_debrief_did_the_trainer_not_ask_2"`
	Were_the_participants_responsive_during_the_debriefing_3         string `json:"were_the_participants_responsive_during_the_debriefing_3"`
	Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_3  string `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_3"`
	Did_this_module_take_30_minutes_as_allotted_3                    string `json:"did_this_module_take_30_minutes_as_allotted_3"`
	During_the_debrief_did_the_gelathi                               string `json:"during_the_debrief_did_the_gelathi"`
	Name_of_the_gelathi_being_evaluated_1                            string `json:"name_of_the_gelathi_being_evaluated_1"`
	Check_which_instructions_the_gelathi_did_not_do                  string `json:"check_which_instructions_the_gelathi_did_not_do"`
	During_the_debrief_did_the_trainer_did_not_do_the_following      string `json:"during_the_debrief_did_the_trainer_did_not_do_the_following"`
	Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_4  string `json:"did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_4"`
	How_many_women_attended_the_training_session_1                   int    `json:"how_many_women_attended_the_training_session_1"`
	How_many_women_attended_the_training_session_2                   int    `json:"how_many_women_attended_the_training_session_2"`
	Check_which_ones_the_trainer_did_not_do_6                        string `json:"check_which_ones_the_trainer_did_not_do_6"`
}

type DashboardData struct {
	SStraining     int `json:"SStraining"`
	GelathiProgram int `json:"GelathiProgram"`
	SSbyGelathi    int `json:"SSbyGelathi"`
}
type Response struct {
	List    []ListQualityAssessmentForm_1 `json:"data"`
	Code    int                           `json:"code"`
	Success bool                          `json:"success"`
	Message string                        `json:"message"`
}
type Response1 struct {
	List    DashboardData `json:"data"`
	Code    int           `json:"code"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
}

var DB *sql.DB

// ------------------------------------Adding Data into Quality Assessment Form----------------------------------------------------------------------
func AddQualityAssessmentForm(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var form QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Failed to decode request body", "Error": err})
		return
	}
	// Fetch email and role from Employee table based on Emp_id
	var email string
	var role int
	err = DB.QueryRow(`SELECT officeMailId, empRole FROM employee WHERE id = ?`, form.EmpID).Scan(&email, &role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to fetch email and role from Employee table", "Error": err})
		return
	}
	// Check if role matches Emp_id
	if role != form.RoleID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Role does not match Emp_id"})
		return
	}

	// Check if empRole is allowed
	allowedRoles := map[int]bool{1: true, 2: true, 3: true, 4: true, 12: true, 13: true}
	if !allowedRoles[role] {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "401 Unauthorized", "Message": "empRole not allowed"})
		return
	}

	// Prepare SQL statement
	query := `
	INSERT INTO QualityAssessmentForm_1 (
		emp_id,
		email_address,
		name_of_the_assessor,
		entry_date,
		program_assessment,
		today_poa,
		name_of_the_district,
		name_of_the_taluk,
		name_of_the_village_and_the_venue_of_meeting_or_training,
		day1_or_day2,
		name_of_the_trainer_being_evaluated,
		check_which_ones_the_trainer_did_not_do,
		how_many_women_attended_the_training_session,
		check_which_ones_the_trainer_did_not_do_1,
		were_the_women_interactive,
		Did_any_women_leave_tring_session_dring_or_after_1st_module,
		if_so_how_many,
		did_this_module_take_20_minutes_as_allotted,
		did_any_new_women_attend_the_training_session_during_module,
		if_so_how_many_1,
		check_which_ones_the_trainer_did_not_do_2,
		during_the_debrief_the_trainer_did,
		did_any_women_leve_training_session_during_or_after_1st_module_1,
		if_so_how_many_2,
		did_this_module_take_20_minutes_as_allotted_1,
		did_any_new_women_attend_training_session_during_this_module_1,
		if_so_how_many_3,
		check_which_ones_the_trainer_did_not_do_3,
		during_the_debriefs_for_role_plays_the_trainer_did_not_ask,
		did_the_trainer_leave_women_to_read_role_play_card_themselves,
		did_the_groups_engage_and_interact_among_themselves_well,
		were_the_participants_responsive_during_the_debriefing,
		did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_2,
		if_so_how_many_4,
		did_this_module_take_30_minutes_as_allotted_1,
		how_many_women_remained_by_the_end_of_this_training_session,
		how_many_are_likely_to_come_back,
		did_any_new_women_attend_training_session_during_this_module_2,
		if_so_how_many_5,
		did_this_module_take_30_minutes_as_allotted,
		check_which_ones_the_trainer_did_not_do_4,
		was_the_recap_done,
		did_the_recap_take_15_minutes_as_allotted,
		name_of_the_gf,
		no_of_participants_at_the_start_of_the_session,
		assessment_of,
		the_gf_comptetly_carried_out_following_funtions,
		the_gf_caried_followig_fuctions_bfore_traning_or_meting_started,
		how_many_stories_of_success_or_change_emerged_from_the_recap,
		ment_name_of_gelathis_success_stories_and_story_couple_of_lines,
		check_which_ones_the_gf_did_not_do,
		number_of_enrolled_gelathis_in_the_circle,
		no_of_attended_gelathis,
		level_of_participation_1,
		level_of_participation_2,
		level_of_participation_3,
		level_of_participation_4,
		level_of_participation_5,
		the_gf_competently_covered_folowing_things_in_training_delivered,
		rate_the_gf_1,
		rate_the_gf_2,
		rate_the_gf_3,
		rate_the_gf_4,
		rate_the_gf_5,
		rate_the_gf_6,
		rate_the_gf_7,
		rate_the_gf_8,
		rate_the_gf_9,
		rate_the_gf_10,
		what_worked_in_the_training,
		what_can_be_better_next_time,
		any_futher_training_and_understding_reqired_by_gf_traing_module,
		did_you_find_anything_traiing_or_gf_that_neds_to_worked_priority,
		details_of_success_stories_to_be_collected_from_gelathis_by_gf,
		deadline_to_collect_the_stories,
		end_time_of_the_training,
		no_of_participants_at_end_of_the_session,
		any_other_comments_about_the_gelathi_facilitator,
		name_of_the_gelathi_being_evaluated,
		days_modules,
		check_which_ones_the_gelathi_did_not_do,
		check_which_ones_the_gelathi_did_not_do_1,
		was_the_recap_done_1,
		did_the_debrief_done_by_gelathi,
		during_the_debriefs_for_role_plays_the_gelathi_did_not_ask,
		repeat_the_activity_with_the_second_volunteer,
		during_the_debrief_did_the_gelathi_not_ask,
		the_gelathi_did_not_ask_1,
		check_which_ones_the_gelathi_did_not_do_2,
		check_which_instructions_the_trainer_did_not_do,
		during_the_debrief_did_the_trainer_not_ask,
		were_the_participants_responsive_during_the_debriefing_1,
		did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_1,
		how_many_women_attended_the_training_session_3,
		the_trainer_did_not_ask,
		what_did_the_trainer_not_do,
		during_the_debrief_did_the_trainer_not_ask_1,
		were_the_participants_responsive_during_the_debriefing_2,
		did_this_module_take_30_minutes_as_allotted_2,
		how_many_women_attended_the_training_session_4,
		check_which_ones_the_trainer_did_not_do_5,
		during_the_debrief_did_the_trainer_not_ask_2,
		were_the_participants_responsive_during_the_debriefing_3,
		did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_3,
		did_this_module_take_30_minutes_as_allotted_3,
		during_the_debrief_did_the_gelathi,
		name_of_the_gelathi_being_evaluated_1,
		check_which_instructions_the_gelathi_did_not_do,
		during_the_debrief_did_the_trainer_did_not_do_the_following,
		did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_4,
		how_many_women_attended_the_training_session_1,
		how_many_women_attended_the_training_session_2,
		check_which_ones_the_trainer_did_not_do_6)
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
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
		?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	_, err = DB.Exec(query,
		form.EmpID,
		form.EmailAddress,
		form.NameOfTheAssessor,
		form.EntryDate,
		form.ProgramAssessment,
		form.TodayPOA,
		form.NameOfTheDistrict,
		form.NameOfTheTaluk,
		form.NameOfTheVillageAndVenueOfMeetingOrTraining,
		form.Day1OrDay2,
		form.NameOfTheTrainerBeingEvaluated,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo, ","),
		form.HowManyWomenAttendedTheTrainingSession,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo1, ","),
		form.WereTheWomenInteractive,
		form.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule,
		form.IfSoHowMany,
		form.DidThisModuleTake20MinutesAsAllotted,
		form.DidAnyNewWomenAttendTheTrainingSessionDuringModule,
		form.IfSoHowMany1,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo2, ","),
		strings.Join(form.DuringTheDebriefTheTrainerDid, ","),
		form.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule1,
		form.IfSoHowMany2,
		form.DidThisModuleTake20MinutesAsAllotted1,
		form.DidAnyNewWomenAttendTrainingSessionDuringThisModule1,
		form.IfSoHowMany3,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo3, ","),
		strings.Join(form.DuringTheDebriefsForRolePlaysTheTrainerDidNotAsk, ","),
		form.DidTheTrainerLeaveWomenToReadRolePlayCardThemselves,
		form.DidTheGroupsEngageAndInteractAmongThemselvesWell,
		form.WereTheParticipantsResponsiveDuringTheDebriefing,
		form.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule2,
		form.IfSoHowMany4,
		form.DidThisModuleTake30MinutesAsAllotted1,
		form.HowManyWomenRemainedByTheEndOfThisTrainingSession,
		form.HowManyAreLikelyToComeBack,
		form.DidAnyNewWomenAttendTrainingSessionDuringThisModule2,
		form.IfSoHowMany5,
		form.DidThisModuleTake30MinutesAsAllotted,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo4, ","),
		form.WasTheRecapDone,
		form.DidTheRecapTake15MinutesAsAllotted,
		form.NameOfTheGF,
		form.NoOfParticipantsAtTheStartOfTheSession,
		form.AssessmentOf,
		strings.Join(form.TheGFCompletelyCarriedOutFollowingFunctions, ","),
		strings.Join(form.TheGFCariedFollowigFuctionsBforeTraningOrMetingStarted, ","),
		form.HowManyStoriesOfSuccessOrChangeEmergedFromTheRecap,
		form.MentNameOfGelathisSuccessStoriesAndStoryCoupleOfLines,
		strings.Join(form.CheckWhichOnesTheGFDidNotDo, ","),
		form.NumberOfEnrolledGelathisInTheCircle,
		form.NoOfAttendedGelathis,
		form.LevelOfParticipation1,
		form.LevelOfParticipation2,
		form.LevelOfParticipation3,
		form.LevelOfParticipation4,
		form.LevelOfParticipation5,
		strings.Join(form.TheGFCompetentlyCoveredFollowingThingsInTrainingDelivered, ","),
		form.RateTheGF1,
		form.RateTheGF2,
		form.RateTheGF3,
		form.RateTheGF4,
		form.RateTheGF5,
		form.RateTheGF6,
		form.RateTheGF7,
		form.RateTheGF8,
		form.RateTheGF9,
		form.RateTheGF10,
		form.WhatWorkedInTheTraining,
		form.WhatCanBeBetterNextTime,
		strings.Join(form.AnyFurtherTrainingAndUnderstandingReqiredByGFTraingModule, ","),
		form.DidYouFindAnythingTraiingOrGFThatNeedsToWorkedPriority,
		form.DetailsOfSuccessStoriesToBeCollectedFromGelathisByGF,
		form.DeadlineToCollectTheStories,
		form.EndTimeOfTheTraining,
		form.NoOfParticipantsAtEndOfTheSession,
		form.AnyOtherCommentsAboutTheGelathiFacilitator,
		form.NameOfTheGelathiBeingEvaluated,
		form.DaysModules,
		strings.Join(form.CheckWhichOnesTheGelathiDidNotDo, ","),
		strings.Join(form.CheckWhichOnesTheGelathiDidNotDo1, ","),
		form.WasTheRecapDone1,
		form.DidTheDebriefDoneByGelathi,
		strings.Join(form.DuringTheDebriefsForRolePlaysTheGelathiDidNotAsk, ","),
		form.RepeatTheActivityWithTheSecondVolunteer,
		strings.Join(form.DuringTheDebriefDidTheGelathiNotAsk, ","),
		strings.Join(form.TheGelathiDidNotAsk1, ","),
		strings.Join(form.CheckWhichOnesTheGelathiDidNotDo2, ","),
		strings.Join(form.CheckWhichInstTrainerDidNotDo, ","),
		strings.Join(form.DuringDebriefDidTheTrainerNotAsk, ","),
		form.WereTheParticipantsResponsiveDuringTheDebriefing1,
		form.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule1,
		form.HowManyWomenAttendedTheTrainingSession3,
		strings.Join(form.TheTrainerDidNotAsk, ","),
		strings.Join(form.WhatDidTheTrainerNotDo, ","),
		strings.Join(form.DuringDebriefDidTheTrainerNotAsk1, ","),
		form.WereTheParticipantsResponsiveDuringTheDebriefing2,
		form.DidThisModuleTake30MinutesAsAllotted2,
		form.HowManyWomenAttendedTheTrainingSession4,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo5, ","),
		strings.Join(form.DuringDebriefDidTheTrainerNotAsk2, ","),
		form.WereTheParticipantsResponsiveDuringTheDebriefing3,
		form.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule3,
		form.DidThisModuleTake30MinutesAsAllotted3,
		strings.Join(form.DuringTheDebriefDidTheGelathi, ","),
		form.NameOfTheGelathiBeingEvaluated1,
		strings.Join(form.CheckWhichInstTheGelathiDidNotDo, ","),
		strings.Join(form.DuringTheDebriefDidTheTrainerDidNotDoTheFollowing, ","),
		form.DidAnyWomenLeaveTheTrainingSessionDuringOnAfterThisModule4,
		form.HowManyWomenAttendedTheTrainingSession1,
		form.HowManyWomenAttendedTheTrainingSession2,
		strings.Join(form.CheckWhichOnesTheTrainerDidNotDo6, ","),
	)

	if err != nil {
		log.Println("AddQualityAssessmentForm", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Status Code": "200", "Message": "Data inserted successfully"})
}

func ListQualityAssessmentForm(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var a QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %v", err)
		return
	}
	// Fetch role from Employee table based on Emp_id
	var roleID int

	err = DB.QueryRow(`SELECT empRole FROM employee WHERE id = ?`, a.EmpID).Scan(&roleID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to fetch email and role from Employee table"})
		return
	}

	// Check if role matches Emp_id
	if roleID != a.RoleID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Role does not match Emp_id"})
		return
	}

	// Rest of your code...

	// Check if empRole is allowed
	allowedRoles := map[int]bool{1: true, 2: true, 3: true, 4: true, 12: true, 13: true}
	if !allowedRoles[roleID] {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "401 Unauthorized", "Message": "empRole not allowed"})
		return
	}

	// Prepare SQL statement
	rows, err := DB.Query(`SELECT * FROM QualityAssessmentForm_1 WHERE Emp_id=?`, a.EmpID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare response as JSON
	var quality []ListQualityAssessmentForm_1
	var EntryDate, DeadlineToCollectTheStories string
	var CheckWhichOnesTheTrainerDidNotDo, CheckWhichOnesTheTrainerDidNotDo1, CheckWhichOnesTheTrainerDidNotDo2, DuringTheDebriefTheTrainerDid, CheckWhichOnesTheTrainerDidNotDo3, DuringTheDebriefsForRolePlaysTheTrainerDidNotAsk, CheckWhichOnesTheTrainerDidNotDo4, TheGFCompletelyCarriedOutFollowingFunctions, TheGFCariedFollowigFuctionsBforeTraningOrMetingStarted, CheckWhichOnesTheGFDidNotDo, TheGFCompetentlyCoveredFollowingThingsInTrainingDelivered, AnyFurtherTrainingAndUnderstandingReqiredByGFTraingModule, CheckWhichOnesTheGelathiDidNotDo, CheckWhichOnesTheGelathiDidNotDo1, DuringTheDebriefsForRolePlaysTheGelathiDidNotAsk, DuringTheDebriefDidTheGelathiNotAsk, TheGelathiDidNotAsk1, CheckWhichOnesTheGelathiDidNotDo2, CheckWhichInstTrainerDidNotDo, DuringDebriefDidTheTrainerNotAsk, TheTrainerDidNotAsk, WhatDidTheTrainerNotDo, DuringDebriefDidTheTrainerNotAsk1, CheckWhichOnesTheTrainerDidNotDo5, DuringDebriefDidTheTrainerNotAsk2, DuringTheDebriefDidTheGelathi, CheckWhichInstTheGelathiDidNotDo, DuringTheDebriefDidTheTrainerDidNotDoTheFollowing, CheckWhichOnesTheTrainerDidNotDo6 string

	for rows.Next() {
		err := rows.Scan(
			&a.ID,
			&a.EmpID,
			&a.EmailAddress,
			&a.NameOfTheAssessor,
			&EntryDate,
			&a.ProgramAssessment,
			&a.TodayPOA,
			&a.NameOfTheDistrict,
			&a.NameOfTheTaluk,
			&a.NameOfTheVillageAndVenueOfMeetingOrTraining,
			&a.Day1OrDay2,
			&a.NameOfTheTrainerBeingEvaluated,
			&CheckWhichOnesTheTrainerDidNotDo,
			&a.HowManyWomenAttendedTheTrainingSession,
			&CheckWhichOnesTheTrainerDidNotDo1,
			&a.WereTheWomenInteractive,
			&a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule,
			&a.IfSoHowMany,
			&a.DidThisModuleTake20MinutesAsAllotted,
			&a.DidAnyNewWomenAttendTheTrainingSessionDuringModule,
			&a.IfSoHowMany1,
			&CheckWhichOnesTheTrainerDidNotDo2,
			&DuringTheDebriefDidTheTrainerDidNotDoTheFollowing,
			&DuringTheDebriefTheTrainerDid,
			&a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule1,
			&a.IfSoHowMany2,
			&a.DidThisModuleTake20MinutesAsAllotted1,
			&a.DidAnyNewWomenAttendTrainingSessionDuringThisModule1,
			&a.IfSoHowMany3,
			&CheckWhichOnesTheTrainerDidNotDo3,
			&DuringTheDebriefsForRolePlaysTheTrainerDidNotAsk,
			&a.DidTheTrainerLeaveWomenToReadRolePlayCardThemselves,
			&a.DidTheGroupsEngageAndInteractAmongThemselvesWell,
			&a.WereTheParticipantsResponsiveDuringTheDebriefing,
			&a.IfSoHowMany4,
			&a.DidThisModuleTake30MinutesAsAllotted,
			&a.HowManyWomenRemainedByTheEndOfThisTrainingSession,
			&a.HowManyAreLikelyToComeBack,
			&a.DidAnyNewWomenAttendTrainingSessionDuringThisModule2,
			&a.IfSoHowMany5,
			&a.DidThisModuleTake30MinutesAsAllotted1,
			&CheckWhichOnesTheTrainerDidNotDo4,
			&a.WasTheRecapDone,
			&a.DidTheRecapTake15MinutesAsAllotted,
			&a.NameOfTheGF,
			&a.NoOfParticipantsAtTheStartOfTheSession,
			&a.AssessmentOf,
			&TheGFCompletelyCarriedOutFollowingFunctions,
			&TheGFCariedFollowigFuctionsBforeTraningOrMetingStarted,
			&a.HowManyStoriesOfSuccessOrChangeEmergedFromTheRecap,
			&a.MentNameOfGelathisSuccessStoriesAndStoryCoupleOfLines,
			&CheckWhichOnesTheGFDidNotDo,
			&a.NumberOfEnrolledGelathisInTheCircle,
			&a.NoOfAttendedGelathis,
			&a.LevelOfParticipation1,
			&a.LevelOfParticipation2,
			&a.LevelOfParticipation3,
			&a.LevelOfParticipation4,
			&a.LevelOfParticipation5,
			&TheGFCompetentlyCoveredFollowingThingsInTrainingDelivered,
			&a.RateTheGF1,
			&a.RateTheGF2,
			&a.RateTheGF3,
			&a.RateTheGF4,
			&a.RateTheGF5,
			&a.RateTheGF6,
			&a.RateTheGF7,
			&a.RateTheGF8,
			&a.RateTheGF9,
			&a.RateTheGF10,
			&a.WhatWorkedInTheTraining,
			&a.WhatCanBeBetterNextTime,
			&AnyFurtherTrainingAndUnderstandingReqiredByGFTraingModule,
			&a.DidYouFindAnythingTraiingOrGFThatNeedsToWorkedPriority,
			&a.DetailsOfSuccessStoriesToBeCollectedFromGelathisByGF,
			&DeadlineToCollectTheStories,
			&a.EndTimeOfTheTraining,
			&a.NoOfParticipantsAtEndOfTheSession,
			&a.AnyOtherCommentsAboutTheGelathiFacilitator,
			&a.NameOfTheGelathiBeingEvaluated,
			&a.DaysModules,
			&CheckWhichOnesTheGelathiDidNotDo,
			&CheckWhichOnesTheGelathiDidNotDo1,
			&a.WasTheRecapDone1,
			&a.DidTheDebriefDoneByGelathi,
			&DuringTheDebriefsForRolePlaysTheGelathiDidNotAsk,
			&a.RepeatTheActivityWithTheSecondVolunteer,
			&DuringTheDebriefDidTheGelathiNotAsk,
			&TheGelathiDidNotAsk1,
			&CheckWhichOnesTheGelathiDidNotDo2,
			&a.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule2,
			&CheckWhichInstTrainerDidNotDo,
			&DuringDebriefDidTheTrainerNotAsk,
			&a.WereTheParticipantsResponsiveDuringTheDebriefing1,
			&a.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule1,
			&a.HowManyWomenAttendedTheTrainingSession3,
			&TheTrainerDidNotAsk,
			&WhatDidTheTrainerNotDo,
			&DuringDebriefDidTheTrainerNotAsk1,
			&a.WereTheParticipantsResponsiveDuringTheDebriefing2,
			&a.DidThisModuleTake30MinutesAsAllotted2,
			&a.HowManyWomenAttendedTheTrainingSession4,
			&CheckWhichOnesTheTrainerDidNotDo5,
			&DuringDebriefDidTheTrainerNotAsk2,
			&a.WereTheParticipantsResponsiveDuringTheDebriefing3,
			&a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule3,
			&a.DidThisModuleTake30MinutesAsAllotted3,
			&a.NameOfTheGelathiBeingEvaluated1,
			&CheckWhichInstTheGelathiDidNotDo,
			&DuringTheDebriefDidTheGelathi,
			&a.DidAnyWomenLeaveTheTrainingSessionDuringOnAfterThisModule4,
			&a.HowManyWomenAttendedTheTrainingSession1,
			&a.HowManyWomenAttendedTheTrainingSession2,
			&CheckWhichOnesTheTrainerDidNotDo6)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := time.Parse("2006-01-02 15:04:05.000000", EntryDate)
		if err != nil {
			panic(err)
		}

		EntryDate1 := t.Format("02 January 2006, 03:04 PM")

		t1, err := time.Parse("2006-01-02", DeadlineToCollectTheStories)
		if err != nil {
			panic(err)
		}

		DeadlineToCollectTheStories1 := t1.Format("02 January 2006")

		quality = append(quality, ListQualityAssessmentForm_1{
			Id:                   a.ID,
			Emp_id:               a.EmpID,
			Email_address:        a.EmailAddress,
			Name_of_the_assessor: a.NameOfTheAssessor,
			Entry_date:           EntryDate1,
			Program_assessment:   a.ProgramAssessment,
			Today_poa:            a.TodayPOA,
			Name_of_the_district: a.NameOfTheDistrict,
			Name_of_the_taluk:    a.NameOfTheTaluk,
			Name_of_the_village_and_the_venue_of_meeting_or_training: a.NameOfTheVillageAndVenueOfMeetingOrTraining,
			Day1_or_day2:                                                     a.Day1OrDay2,
			Name_of_the_trainer_being_evaluated:                              a.NameOfTheTrainerBeingEvaluated,
			Check_which_ones_the_trainer_did_not_do:                          CheckWhichOnesTheTrainerDidNotDo,
			How_many_women_attended_the_training_session:                     a.HowManyWomenAttendedTheTrainingSession,
			Check_which_ones_the_trainer_did_not_do_1:                        CheckWhichOnesTheTrainerDidNotDo1,
			Were_the_women_interactive:                                       a.WereTheWomenInteractive,
			Did_any_women_leave_tring_session_dring_or_after_1st_module:      a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule,
			If_so_how_many:                                                   a.IfSoHowMany,
			Did_this_module_take_20_minutes_as_allotted:                      a.DidThisModuleTake20MinutesAsAllotted,
			Did_any_new_women_attend_the_training_session_during_module:      a.DidAnyNewWomenAttendTheTrainingSessionDuringModule,
			If_so_how_many_1:                                                 a.IfSoHowMany1,
			Check_which_ones_the_trainer_did_not_do_2:                        CheckWhichOnesTheTrainerDidNotDo2,
			During_the_debrief_did_the_trainer_did_not_do_the_following:      DuringTheDebriefDidTheTrainerDidNotDoTheFollowing,
			During_the_debrief_the_trainer_did:                               DuringTheDebriefTheTrainerDid,
			Did_any_women_leve_training_session_during_or_after_1st_module_1: a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule1,
			If_so_how_many_2:                                                 a.IfSoHowMany2,
			Did_this_module_take_20_minutes_as_allotted_1:                    a.DidThisModuleTake20MinutesAsAllotted1,
			Did_any_new_women_attend_training_session_during_this_module_1:   a.DidAnyNewWomenAttendTrainingSessionDuringThisModule1,
			If_so_how_many_3:                                                 a.IfSoHowMany3,
			Check_which_ones_the_trainer_did_not_do_3:                        CheckWhichOnesTheTrainerDidNotDo3,
			During_the_debriefs_for_role_plays_the_trainer_did_not_ask:       DuringTheDebriefsForRolePlaysTheTrainerDidNotAsk,
			Did_the_trainer_leave_women_to_read_role_play_card_themselves:    a.DidTheTrainerLeaveWomenToReadRolePlayCardThemselves,
			Did_the_groups_engage_and_interact_among_themselves_well:         a.DidTheGroupsEngageAndInteractAmongThemselvesWell,
			Were_the_participants_responsive_during_the_debriefing:           a.WereTheParticipantsResponsiveDuringTheDebriefing,
			If_so_how_many_4:                                                 a.IfSoHowMany4,
			Did_this_module_take_30_minutes_as_allotted:                      a.DidThisModuleTake30MinutesAsAllotted,
			How_many_women_remained_by_the_end_of_this_training_session:      a.HowManyWomenRemainedByTheEndOfThisTrainingSession,
			How_many_are_likely_to_come_back:                                 a.HowManyAreLikelyToComeBack,
			Did_any_new_women_attend_training_session_during_this_module_2:   a.DidAnyNewWomenAttendTrainingSessionDuringThisModule2,
			If_so_how_many_5:                                                 a.IfSoHowMany5,
			Did_this_module_take_30_minutes_as_allotted_1:                    a.DidThisModuleTake30MinutesAsAllotted1,
			Check_which_ones_the_trainer_did_not_do_4:                        CheckWhichOnesTheTrainerDidNotDo4,
			Was_the_recap_done:                                               a.WasTheRecapDone,
			Did_the_recap_take_15_minutes_as_allotted:                        a.DidTheRecapTake15MinutesAsAllotted,
			Name_of_the_gf: a.NameOfTheGF,
			No_of_participants_at_the_start_of_the_session: a.NoOfParticipantsAtTheStartOfTheSession,
			Assessment_of: a.AssessmentOf,
			The_gf_comptetly_carried_out_following_funtions:                  TheGFCompletelyCarriedOutFollowingFunctions,
			The_gf_caried_followig_fuctions_bfore_traning_or_meting_started:  TheGFCariedFollowigFuctionsBforeTraningOrMetingStarted,
			How_many_stories_of_success_or_change_emerged_from_the_recap:     a.HowManyStoriesOfSuccessOrChangeEmergedFromTheRecap,
			Ment_name_of_gelathis_success_stories_and_story_couple_of_lines:  a.MentNameOfGelathisSuccessStoriesAndStoryCoupleOfLines,
			Check_which_ones_the_gf_did_not_do:                               CheckWhichOnesTheGFDidNotDo,
			Number_of_enrolled_gelathis_in_the_circle:                        a.NumberOfEnrolledGelathisInTheCircle,
			No_of_attended_gelathis:                                          a.NoOfAttendedGelathis,
			Level_of_participation_1:                                         a.LevelOfParticipation1,
			Level_of_participation_2:                                         a.LevelOfParticipation2,
			Level_of_participation_3:                                         a.LevelOfParticipation3,
			Level_of_participation_4:                                         a.LevelOfParticipation4,
			Level_of_participation_5:                                         a.LevelOfParticipation5,
			The_gf_competently_covered_folowing_things_in_training_delivered: TheGFCompetentlyCoveredFollowingThingsInTrainingDelivered,
			Rate_the_gf_1:                a.RateTheGF1,
			Rate_the_gf_2:                a.RateTheGF2,
			Rate_the_gf_3:                a.RateTheGF3,
			Rate_the_gf_4:                a.RateTheGF4,
			Rate_the_gf_5:                a.RateTheGF5,
			Rate_the_gf_6:                a.RateTheGF6,
			Rate_the_gf_7:                a.RateTheGF7,
			Rate_the_gf_8:                a.RateTheGF8,
			Rate_the_gf_9:                a.RateTheGF9,
			Rate_the_gf_10:               a.RateTheGF10,
			What_worked_in_the_training:  a.WhatWorkedInTheTraining,
			What_can_be_better_next_time: a.WhatCanBeBetterNextTime,
			Any_futher_training_and_understding_reqired_by_gf:                AnyFurtherTrainingAndUnderstandingReqiredByGFTraingModule,
			Did_you_find_anything_traiing_or_gf_that_neds_to_worked_priority: a.DidYouFindAnythingTraiingOrGFThatNeedsToWorkedPriority,
			Details_of_success_stories_to_be_collected_from_gelathis_by_gf:   a.DetailsOfSuccessStoriesToBeCollectedFromGelathisByGF,
			Deadline_to_collect_the_stories:                                  DeadlineToCollectTheStories1,
			End_time_of_the_training:                                         a.EndTimeOfTheTraining,
			No_of_participants_at_end_of_the_session:                         a.NoOfParticipantsAtEndOfTheSession,
			Any_other_comments_about_the_gelathi_facilitator:                 a.AnyOtherCommentsAboutTheGelathiFacilitator,
			Name_of_the_gelathi_being_evaluated:                              a.NameOfTheGelathiBeingEvaluated,
			Days_modules:                                                     a.DaysModules,
			Check_which_ones_the_gelathi_did_not_do:                          CheckWhichOnesTheGelathiDidNotDo,
			Check_which_ones_the_gelathi_did_not_do_1:                        CheckWhichOnesTheGelathiDidNotDo1,
			Was_the_recap_done_1:                                             a.WasTheRecapDone1,
			Did_the_debrief_done_by_gelathi:                                  a.DidTheDebriefDoneByGelathi,
			During_the_debriefs_for_role_plays_the_gelathi_did_not_ask:       DuringTheDebriefsForRolePlaysTheGelathiDidNotAsk,
			Repeat_the_activity_with_the_second_volunteer:                    a.RepeatTheActivityWithTheSecondVolunteer,
			During_the_debrief_did_the_gelathi_not_ask:                       DuringTheDebriefDidTheGelathiNotAsk,
			The_gelathi_did_not_ask_1:                                        TheGelathiDidNotAsk1,
			Check_which_ones_the_gelathi_did_not_do_2:                        CheckWhichOnesTheGelathiDidNotDo2,
			Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_2:  a.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule2,
			Check_which_instructions_the_trainer_did_not_do:                  CheckWhichInstTrainerDidNotDo,
			During_the_debrief_did_the_trainer_not_ask:                       DuringDebriefDidTheTrainerNotAsk,
			Were_the_participants_responsive_during_the_debriefing_1:         a.WereTheParticipantsResponsiveDuringTheDebriefing1,
			Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_1:  a.DidAnyWomenLeaveTrainingSessionDuringOrAfterThis1stModule1,
			How_many_women_attended_the_training_session_3:                   a.HowManyWomenAttendedTheTrainingSession3,
			The_trainer_did_not_ask:                                          TheTrainerDidNotAsk,
			What_did_the_trainer_not_do:                                      WhatDidTheTrainerNotDo,
			During_the_debrief_did_the_trainer_not_ask_1:                     DuringDebriefDidTheTrainerNotAsk1,
			Were_the_participants_responsive_during_the_debriefing_2:         a.WereTheParticipantsResponsiveDuringTheDebriefing2,
			Did_this_module_take_30_minutes_as_allotted_2:                    a.DidThisModuleTake30MinutesAsAllotted2,
			How_many_women_attended_the_training_session_4:                   a.HowManyWomenAttendedTheTrainingSession4,
			Check_which_ones_the_trainer_did_not_do_5:                        CheckWhichOnesTheTrainerDidNotDo5,
			During_the_debrief_did_the_trainer_not_ask_2:                     DuringDebriefDidTheTrainerNotAsk2,
			Were_the_participants_responsive_during_the_debriefing_3:         a.WereTheParticipantsResponsiveDuringTheDebriefing3,
			Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_3:  a.DidAnyWomenLeaveTrainingSessionDuringOrAfter1stModule3,
			Did_this_module_take_30_minutes_as_allotted_3:                    a.DidThisModuleTake30MinutesAsAllotted3,
			Name_of_the_gelathi_being_evaluated_1:                            a.NameOfTheGelathiBeingEvaluated1,
			Check_which_instructions_the_gelathi_did_not_do:                  CheckWhichInstTheGelathiDidNotDo,
			During_the_debrief_did_the_gelathi:                               DuringTheDebriefDidTheGelathi,
			Did_any_wmen_leave_the_trning_sesion_during_or_aftr_tis_modle_4:  a.DidAnyWomenLeaveTheTrainingSessionDuringOnAfterThisModule4,
			How_many_women_attended_the_training_session_1:                   a.HowManyWomenAttendedTheTrainingSession1,
			How_many_women_attended_the_training_session_2:                   a.HowManyWomenAttendedTheTrainingSession2,
			Check_which_ones_the_trainer_did_not_do_6:                        CheckWhichOnesTheTrainerDidNotDo6})
	}

	response := Response{
		List:    quality,
		Code:    200,
		Success: true,
		Message: "Successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func GetDashboard(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var q QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	if q.EmpID == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": "Emp_id is required"})
		return
	}
	var role int
	err = DB.QueryRow(`SELECT empRole FROM employee WHERE id = ?`, q.EmpID).Scan(&role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 Internal Server Error", "Message": "Failed to fetch role from Employee table"})
		return
	}

	// Check if role matches Emp_id
	if role != q.RoleID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Role does not match Emp_id"})
		return
	}
	var query string

	var SStraining, GelathiProgram, SSbyGelathi int
	if q.Filter == "" {
		query = `
		SELECT
    COALESCE(SUM(CASE WHEN Program_assessment = 1 THEN 1 ELSE 0 END), 0) AS SStraining,
    COALESCE(SUM(CASE WHEN Program_assessment = 2 THEN 1 ELSE 0 END), 0) AS GelathiProgram,
    COALESCE(SUM(CASE WHEN Program_assessment = 3 THEN 1 ELSE 0 END), 0) AS SSbyGelathi
FROM
    QualityAssessmentForm_1
WHERE
    Emp_id = ?`
		rows, err := DB.Query(query, q.EmpID)
		if err != nil {
			log.Println("GetDashboard", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
			if err != nil {
				log.Println("GetDashboard", err)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
				return
			}
		}
	} else if q.Filter == "Location" {
		query = `SELECT
		COALESCE(SUM(CASE WHEN Program_assessment = 1 THEN 1 ELSE 0 END), 0) AS SStraining,
		COALESCE(SUM(CASE WHEN Program_assessment = 2 THEN 1 ELSE 0 END), 0) AS GelathiProgram,
		COALESCE(SUM(CASE WHEN Program_assessment = 3 THEN 1 ELSE 0 END), 0) AS SSbyGelathi
	FROM
		QualityAssessmentForm_1
	WHERE
		Emp_id = ? and name_of_the_district in (?)and name_of_the_taluk in (?)`

		rows, err := DB.Query(query, q.EmpID, q.Filter_District, q.Filter_Taluk)
		if err != nil {
			log.Println("GetDashboard", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
			if err != nil {
				log.Println("GetDashboard", err)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
				return
			}
		}
	} else if q.Filter == "Role" {
		query = `
		SELECT
    COALESCE(SUM(CASE WHEN Program_assessment = 1 THEN 1 ELSE 0 END), 0) AS SStraining,
    COALESCE(SUM(CASE WHEN Program_assessment = 2 THEN 1 ELSE 0 END), 0) AS GelathiProgram,
    COALESCE(SUM(CASE WHEN Program_assessment = 3 THEN 1 ELSE 0 END), 0) AS SSbyGelathi
FROM
    QualityAssessmentForm_1
WHERE
    Emp_id = ?`
		rows, err := DB.Query(query, q.Filter_Id)
		if err != nil {
			log.Println("GetDashboard", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
			if err != nil {
				log.Println("GetDashboard", err)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
				return
			}
		}
	} else if q.Filter == "Date" {
		query = `SELECT
		COALESCE(SUM(CASE WHEN Program_assessment = 1 THEN 1 ELSE 0 END), 0) AS SStraining,
		COALESCE(SUM(CASE WHEN Program_assessment = 2 THEN 1 ELSE 0 END), 0) AS GelathiProgram,
		COALESCE(SUM(CASE WHEN Program_assessment = 3 THEN 1 ELSE 0 END), 0) AS SSbyGelathi
	FROM
		QualityAssessmentForm_1
	WHERE
		Emp_id = ? and entry_date between ? and ?`
		rows, err := DB.Query(query, q.EmpID, q.Filter_StartDate, q.Filter_EndDate)
		if err != nil {
			log.Println("GetDashboard", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
			if err != nil {
				log.Println("GetDashboard", err)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
				return
			}
		}
	}
	// rows, err := DB.Query(query)
	// if err != nil {
	// 	log.Println("GetDashboard", err)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
	// 	return
	// }
	// defer rows.Close()

	// var SStraining, GelathiProgram, SSbyGelathi int
	// for rows.Next() {
	// 	err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
	// 	if err != nil {
	// 		log.Println("GetDashboard", err)
	// 		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
	// 		return
	// 	}
	// }
	dashboard := DashboardData{SStraining, GelathiProgram, SSbyGelathi}

	response := Response1{
		List:    dashboard,
		Code:    200,
		Success: true,
		Message: "Successfully",
	}

	json.NewEncoder(w).Encode(response)
}

type Filter struct {
	Start_Date string `json:"Start_Date"`
	End_Date   string `json:"End_Date"`
}

func FilterDate(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	var filter Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	// Parse and validate the start and end date strings
	startDate, err := time.Parse("2006-01-02", filter.Start_Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", filter.End_Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid end date format"})
		return
	}

	// Check that the end date is greater than or equal to the start date
	if endDate.Before(startDate) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "End date must be greater than or equal to start date"})
		return
	}

	rows, err := DB.Query(`SELECT 
        (SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=1 AND Entry_date BETWEEN ? AND ?) as SStraining,
        (SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=2 AND Entry_date BETWEEN ? AND ?) as GelathiProgram,
        (SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=3 AND Entry_date BETWEEN ? AND ?) as SSbyGelathi;`,
		filter.Start_Date, filter.End_Date, filter.Start_Date, filter.End_Date, filter.Start_Date, filter.End_Date)
	if err != nil {
		log.Println("FilterDate", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}

	var SStraining, GelathiProgram, SSbyGelathi int
	for rows.Next() {
		err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
		if err != nil {
			log.Println("FilterDate", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
			return
		}
	}

	dashboard := DashboardData{SStraining, GelathiProgram, SSbyGelathi}
	ww, err := json.Marshal(dashboard)
	if err != nil {
		log.Println("FilterDate", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}

	w.Write(ww)
}

type Emp struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	EmpRole      int    `json:"empRole"`
	SupervisorID int    `json:"supervisorId"`
}

func GetEmpData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data struct {
		ID      int `json:"id"`
		Role_Id int `json:"role_id"`
	}
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Bad Request", "Code": "400", "Message": err})
		return
	}
	if data.Role_Id != 1 && data.Role_Id != 2 && data.Role_Id != 3 && data.Role_Id != 4 && data.Role_Id != 12 && data.Role_Id != 13 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Status Not Found", "Code": "404", "Message": "role id not found"})
		return
	}
	// Find the supervisor employee with the given ID
	var supervisor Emp
	err = DB.QueryRow("SELECT id, CONCAT(first_name, ' ', last_name), empRole, supervisorId FROM employee WHERE id = ? AND empRole IN (1, 2, 3, 4, 12, 13)", data.ID).Scan(&supervisor.ID, &supervisor.Name, &supervisor.EmpRole, &supervisor.SupervisorID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Status Not Found", "Code": "404", "Message": "Supervisor not found", "Error": err})
		return
	}

	// Find all employees under the supervisor
	rows, err := DB.Query("SELECT id, CONCAT(first_name, ' ', last_name), empRole, supervisorId FROM employee WHERE supervisorId = ? AND empRole IN (1, 2, 3, 4, 12, 13)", supervisor.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Internal server error", "Code": "500", "Message": "Invalid Input", "Error": err})
		log.Println(err)
		return
	}
	defer rows.Close()

	var subordinates []Emp
	for rows.Next() {
		var subordinate Emp
		err := rows.Scan(&subordinate.ID, &subordinate.Name, &subordinate.EmpRole, &subordinate.SupervisorID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// fmt.Fprintf(w, "Internal server error")
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Internal server error", "Code": "500", "Message": "Invalid Input", "Error": err})
			log.Println(err)
			return
		}
		subordinates = append(subordinates, subordinate)
	}

	// Find subordinates' subordinates recursively
	for i := 0; i < len(subordinates); i++ {
		subordinate := subordinates[i]
		rows, err := DB.Query("SELECT id, CONCAT(first_name, ' ', last_name), empRole, supervisorId FROM employee WHERE supervisorId = ? AND empRole IN (1, 2, 3, 4, 12, 13)", subordinate.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Internal server error", "Code": "500", "Message": "Invalid Input", "Error": err})
			// fmt.Fprintf(w, "Internal server error")
			log.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var subSubordinate Emp
			err := rows.Scan(&subSubordinate.ID, &subSubordinate.Name, &subSubordinate.EmpRole, &subSubordinate.SupervisorID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				// fmt.Fprintf(w, "Internal server error")
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Internal server error", "Code": "500", "Message": "INvalid Input", "Error": err})
				log.Println(err)
				return
			}
			subordinates = append(subordinates, subSubordinate)

		}

	}
	var filteredSubordinates []Emp
	for _, emp := range subordinates {
		if emp.EmpRole == data.Role_Id {
			filteredSubordinates = append(filteredSubordinates, emp)
		} else {
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Status": "StatusBadRequest", "Code": "400", "Message": "Invalid Role Id", "Error": err})
				return
			}
		}
	}

	// Encode the filteredSubordinates slice as JSON

	if len(filteredSubordinates) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "Bad Request", "Code": "400", "Message": "No Employee Data"})
		return
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "OK", "Code": "200", "Message": "Successfull", "Data": filteredSubordinates})
	}
}

func ListDistrict(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var q QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	// Verify that the employee ID matches the specified role
	var role int
	if err := DB.QueryRow("SELECT empRole FROM employee WHERE id = ?", q.EmpID).Scan(&role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if role != q.RoleID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Employee ID does not match role"})
		return
	}

	// Fetch distinct taluks from QualityAssessmentForm table
	rows, err := DB.Query("SELECT DISTINCT Name_of_the_district FROM QualityAssessmentForm_1 WHERE Emp_id = ?", q.EmpID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid Input"})
		return
	}
	defer rows.Close()

	// Collect taluk names into slice
	var districts []string
	for rows.Next() {
		var district string
		if err := rows.Scan(&district); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		districts = append(districts, district)
	}

	// Create the response struct and encode it as JSON
	jsonData, err := json.Marshal(districts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 StatusInternalServerError", "Message": "No response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func ListTaluk(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var q QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	// Verify that the employee ID matches the specified role
	var role int
	if err := DB.QueryRow("SELECT empRole FROM employee WHERE id = ?", q.EmpID).Scan(&role); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 StatusInternalServerError", "Message": "No response"})
		return
	}
	if role != q.RoleID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Employee ID does not match role"})
		return
	}

	// Fetch distinct taluks from QualityAssessmentForm table
	rows, err := DB.Query("SELECT DISTINCT Name_of_the_taluk FROM QualityAssessmentForm_1 WHERE Emp_id = ?", q.EmpID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid Input"})
		return
	}
	defer rows.Close()

	// Collect taluk names into slice
	var taluks []string
	for rows.Next() {
		var taluk string
		if err := rows.Scan(&taluk); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid Taluk ID"})
			return
		}
		taluks = append(taluks, taluk)
	}

	// Create the response struct and encode it as JSON
	jsonData, err := json.Marshal(taluks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "500 StatusInternalServerError", "Message": "No response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
func FilterDistrict(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var q QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	if q.NameOfTheDistrict == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": "Name_of_the_district is required"})
		return
	}

	rows, err := DB.Query(`SELECT 
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=1 AND Name_of_the_district=?) as SStraining,
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=2 AND Name_of_the_district=?) as GelathiProgram,
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=3 AND Name_of_the_district=?) as SSbyGelathi`, q.NameOfTheDistrict, q.NameOfTheDistrict, q.NameOfTheDistrict)
	if err != nil {
		log.Println("FilterDistrict", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
		return
	}
	defer rows.Close()

	var SStraining, GelathiProgram, SSbyGelathi int
	for rows.Next() {
		err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
		if err != nil {
			log.Println("FilterDistrict", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
	}

	// Check if any of the counts is zero, indicating that the Name_of_the_district is not in the database
	if SStraining == 0 && GelathiProgram == 0 && SSbyGelathi == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": "District not found in database"})
		return
	}

	dashboard := DashboardData{SStraining, GelathiProgram, SSbyGelathi}
	ww, err := json.Marshal(dashboard)
	if err != nil {
		log.Println("FilterDistrict", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}
	w.Write(ww)
}
func FilterTaluk(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var q QualityAssessmentForm_1
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	if q.NameOfTheTaluk == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": "Name_of_the_taluk is required"})
		return
	}

	rows, err := DB.Query(`SELECT 
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=1 AND Name_of_the_taluk=?) as SStraining,
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=2 AND Name_of_the_taluk=?) as GelathiProgram,
		(SELECT COUNT(*) FROM QualityAssessmentForm_1 WHERE Program_assessment=3 AND Name_of_the_taluk=?) as SSbyGelathi`, q.NameOfTheTaluk, q.NameOfTheTaluk, q.NameOfTheTaluk)
	if err != nil {
		log.Println("FilterTaluk", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
		return
	}
	defer rows.Close()

	var SStraining, GelathiProgram, SSbyGelathi int
	for rows.Next() {
		err := rows.Scan(&SStraining, &GelathiProgram, &SSbyGelathi)
		if err != nil {
			log.Println("FilterTaluk", err)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
			return
		}
	}

	// Check if any of the counts is zero, indicating that the Name_of_the_taluk is not in the database
	if SStraining == 0 && GelathiProgram == 0 && SSbyGelathi == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": "Taluk not found in database"})
		return
	}

	dashboard := DashboardData{SStraining, GelathiProgram, SSbyGelathi}
	ww, err := json.Marshal(dashboard)
	if err != nil {
		log.Println("FilterTaluk", err)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}
	w.Write(ww)
}

type District struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ListAllDistrict(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	rows, err := DB.Query("SELECT id, name FROM bdms_staff.location WHERE type = 3 AND parentId = 3")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "Invalid input syntax"})
		return
	}
	defer rows.Close()

	// Collect districts into a slice
	var districts []District
	for rows.Next() {
		var district District
		if err := rows.Scan(&district.ID, &district.Name); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "Invalid Scan"})
			return
		}
		districts = append(districts, district)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(districts)
}

type DistrictRequest struct {
	DistID string `json:"dist_id"`
}

func ListTalukBasedOnDist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	var distReq DistrictRequest
	err := json.NewDecoder(r.Body).Decode(&distReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	if distReq.DistID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Error": "Missing Dist_id"})
		return
	}

	// Check if parent ID exists
	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM bdms_staff.location WHERE type = 3 AND id = ? LIMIT 1)", distReq.DistID).Scan(&exists)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "Invalid INput Syntax"})
		return
	}
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Error": "Invalid Dist_id"})
		return
	}

	rows, err := DB.Query("SELECT name FROM bdms_staff.location WHERE type = 4 AND parentId = ?", distReq.DistID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "Invalid INput Syntax"})
		return
	}
	defer rows.Close()

	// Collect taluk names into slice
	var taluks []string
	for rows.Next() {
		var taluk string
		if err := rows.Scan(&taluk); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "Invalid Scan"})
			return
		}
		taluks = append(taluks, taluk)
	}

	// Create the response JSON
	responseJSON, err := json.Marshal(taluks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Status": "400 Bad Request", "Message": "No response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

//--------------------------------------------Get Today's POA--------------------------
type POARequest struct {
	Type        string `json:"type"`
	SessionType string `json:"session_type"`
}

func GetPoa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	var q POARequest
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err.Error()})
		return
	}

	// Get the current date
	currentDate := time.Now().Format("2006-01-02")

	// Construct the SQL query based on the input parameters
	query := `
		SELECT tbl_poa.name as poa_name, CONCAT(employee.first_name, ' ', employee.last_name) as full_name, employee.empRole 
		FROM tbl_poa 
		INNER JOIN employee ON tbl_poa.user_id = employee.id 
		WHERE DATE(tbl_poa.check_in) = ?`

	var args []interface{}
	args = append(args, currentDate)

	if q.Type != "" {
		query += " AND tbl_poa.type = ?"
		args = append(args, q.Type)
	}

	if q.SessionType != "" {
		query += " AND tbl_poa.session_type = ?"
		args = append(args, q.SessionType)
	}

	// Execute the SQL query
	rows, err := DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect names and additional fields into a modified version of the Emp struct
	var employees []map[string]interface{}
	for rows.Next() {
		employee := make(map[string]interface{})
		var poaName, fullName string
		var empRole int
		if err := rows.Scan(&poaName, &fullName, &empRole); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employee["name"] = poaName
		employee["full_name"] = fullName
		employee["empRole"] = empRole
		employees = append(employees, employee)
	}

	// Create the response as an array of objects
	response := employees

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
