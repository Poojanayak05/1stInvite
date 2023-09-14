package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	a "buzzstaff-go/New_Enchancement/Attendence"
	d "buzzstaff-go/New_Enchancement/Dashboard"
	gp "buzzstaff-go/New_Enchancement/GP/Dashboard"
	gppf "buzzstaff-go/New_Enchancement/GP/Participant"
	q "buzzstaff-go/New_Enchancement/QualityAssessment"
	sppf "buzzstaff-go/New_Enchancement/SP/Participant"
	sp "buzzstaff-go/New_Enchancement/SP/SPDashboard"
	spf "buzzstaff-go/New_Enchancement/SS/Participant"
	ss "buzzstaff-go/New_Enchancement/SS/SSDashboard"
	vppf "buzzstaff-go/New_Enchancement/VP/Participant"
	vp "buzzstaff-go/New_Enchancement/VP/VPDashboard"
	t1 "buzzstaff-go/Old_Enchancement/Team_1"
	t2 "buzzstaff-go/Old_Enchancement/Team_2"
	t3 "buzzstaff-go/Old_Enchancement/Team_3"
	t4 "buzzstaff-go/Old_Enchancement/Team_4"
	t5 "buzzstaff-go/Old_Enchancement/Team_5"
	dbs "buzzstaff-go/database"
)

func HandleFunc() {

	db := dbs.Connect()

	router := mux.NewRouter()
	apiPrefix := "/appGo"
	apiRouter := router.PathPrefix(apiPrefix).Subrouter()

	//-------------------Endpoint for adding data to QAF---------------------------------------
	apiRouter.HandleFunc("/addQualityAssessmentForm", func(w http.ResponseWriter, r *http.Request) {
		q.AddQualityAssessmentForm(w, r, db)
	})
	//---------------------List Quality Assessment Form-----------------------------------------
	apiRouter.HandleFunc("/listQualityAssessmentForm", func(w http.ResponseWriter, r *http.Request) {
		q.ListQualityAssessmentForm(w, r, db)
	})

	//-------------------Endpoint to get the total dashboard data for QAF-----------------------
	apiRouter.HandleFunc("/getDashboard", func(w http.ResponseWriter, r *http.Request) {
		q.GetDashboard(w, r, db)
	})

	//---------------------Endpoint to filter the dashboard data(Date Range) for QAF------------------------
	apiRouter.HandleFunc("/filterDate", func(w http.ResponseWriter, r *http.Request) {
		q.FilterDate(w, r, db)
	})

	//-----------------------list all district(general)----------------------------------
	apiRouter.HandleFunc("/allDist", func(w http.ResponseWriter, r *http.Request) {
		q.ListAllDistrict(w, r, db)
	})
	//-----------------------list all taluk(general)----------------------------------
	apiRouter.HandleFunc("/listTaluk", func(w http.ResponseWriter, r *http.Request) {
		q.ListTalukBasedOnDist(w, r, db)
	})

	//---------------------Get Employee Data for filter for QAF-----------------------------------------
	apiRouter.HandleFunc("/getEmpData", func(w http.ResponseWriter, r *http.Request) {
		q.GetEmpData(w, r, db)
	})

	//---------------------Endpoint to filter the dashboard data(Taluk) for QAF-----------------------------------------
	apiRouter.HandleFunc("/filterTaluk", func(w http.ResponseWriter, r *http.Request) {
		q.FilterTaluk(w, r, db)
	})

	//---------------------Endpoint to filter the dashboard data(District) for QAF-----------------------------------------
	apiRouter.HandleFunc("/filterDistrict", func(w http.ResponseWriter, r *http.Request) {
		q.FilterDistrict(w, r, db)
	})

	//---------------------Endpoint to list the taluks for filter QAF----------------------------------
	apiRouter.HandleFunc("/listTaluk", func(w http.ResponseWriter, r *http.Request) {
		q.ListTaluk(w, r, db)
	})

	//---------------------Endpoint to list the districts for filter QAF----------------------------------
	apiRouter.HandleFunc("/listDistrict", func(w http.ResponseWriter, r *http.Request) {
		q.ListDistrict(w, r, db)
	})

	//---------------------listing of spoorthi participants based on the Modules1----------------------------------
	apiRouter.HandleFunc("/spoorthiList", func(w http.ResponseWriter, r *http.Request) {
		a.Spoorthilist(w, r, db)
	})

	//---------------------buzz vyapar list----------------------------------
	apiRouter.HandleFunc("/buzzList", func(w http.ResponseWriter, r *http.Request) {
		a.Buzzlist(w, r, db)
	})

	//---------------------Green survey list----------------------------------
	apiRouter.HandleFunc("/greenList", func(w http.ResponseWriter, r *http.Request) {
		a.Greenlist(w, r, db)
	})

	//---------------------attendence of spoorthi----------------------------------
	apiRouter.HandleFunc("/allAttendence", func(w http.ResponseWriter, r *http.Request) {
		a.AllAttendence(w, r, db)
	})

	//---------------------Add buzz vyapar----------------------------------
	apiRouter.HandleFunc("/addBuzzVyapar", func(w http.ResponseWriter, r *http.Request) {
		a.Addbuzzvyapar(w, r, db)
	})

	//---------------------Get POA----------------------------------
	apiRouter.HandleFunc("/getPoa", func(w http.ResponseWriter, r *http.Request) {
		q.GetPoa(w, r, db)
	})

	apiRouter.HandleFunc("/funderVyaparDashboard", func(w http.ResponseWriter, r *http.Request) {
		d.FunderVyaparD(w, r, db)
	})

	apiRouter.HandleFunc("/funderGreenDashboard", func(w http.ResponseWriter, r *http.Request) {
		d.FunderGreenD(w, r, db)
	})
	apiRouter.HandleFunc("/funderSSDashboard", func(w http.ResponseWriter, r *http.Request) {
		d.FunderSSD(w, r, db)
	})
	apiRouter.HandleFunc("/funderGelathiDashboard", func(w http.ResponseWriter, r *http.Request) {
		d.FunderGelathiD(w, r, db)
	})

	//-----------------------------------------SS filter-anand-----------------------------
	apiRouter.HandleFunc("/ssfilter", func(w http.ResponseWriter, r *http.Request) {
		spf.SSCounts(w, r, db)
	})

	//-----------------------------------------Gelathi filter-anand-----------------------------
	apiRouter.HandleFunc("/gelathifilter", func(w http.ResponseWriter, r *http.Request) {
		sppf.SPCounts(w, r, db)
	})

	//-----------------------------------------Green filter-anand-----------------------------
	apiRouter.HandleFunc("/greenfilter", func(w http.ResponseWriter, r *http.Request) {
		gppf.GPCounts(w, r, db)
	})

	//-----------------------------------------Vyapar filter-anand-----------------------------
	apiRouter.HandleFunc("/vyaparfilter", func(w http.ResponseWriter, r *http.Request) {
		vppf.GetCounts(w, r, db)
	})
	//-----------------------------------------SS Dashboard-anand-----------------------------
	apiRouter.HandleFunc("/ssDashboard", func(w http.ResponseWriter, r *http.Request) {
		ss.SelfsakthiDashboard(w, r, db)
	})
	// -----------------------------------------GP dashboard-----------------------------
	apiRouter.HandleFunc("/greenDashboard", func(w http.ResponseWriter, r *http.Request) {
		gp.GreenDashboard(w, r, db)
	})
	//---------------------vyapar dashboard----------------------------------
	apiRouter.HandleFunc("/vyaparDashboard", func(w http.ResponseWriter, r *http.Request) {
		vp.VyaparDashboard(w, r, db)
	})

	//---------------------gelathiProgram dashboard----------------------------------
	apiRouter.HandleFunc("/gelathiProgramDashboard", func(w http.ResponseWriter, r *http.Request) {
		sp.GelathiProgramDashboard1(w, r, db)
	})

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>old codes converted to go>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//_______________________________TEAM_1___________________________________

	//---------------------Check Email Exist-----------------------------------------

	apiRouter.HandleFunc("/getEmailExist", func(w http.ResponseWriter, r *http.Request) {
		t1.GetEmailExist(w, r, db)
	})

	//---------------------Check Profile Details-----------------------------------------

	apiRouter.HandleFunc("/getProfileData", func(w http.ResponseWriter, r *http.Request) {
		t1.GetProfileData(w, r, db)
	})

	//---------------------Check Circle Details-----------------------------------------

	apiRouter.HandleFunc("/getGelathiCircleDataNew", func(w http.ResponseWriter, r *http.Request) {
		t1.GetGelathiCircleDataNew(w, r, db)
	})

	//---------------------Add Employee -----------------------------------------
	apiRouter.HandleFunc("/createUser", func(w http.ResponseWriter, r *http.Request) {
		t1.CreateUser(w, r, db)
	})
	//---------------------Add Funder -----------------------------------------
	apiRouter.HandleFunc("/createFunder", func(w http.ResponseWriter, r *http.Request) {
		t1.CreateFunder(w, r, db)
	})

	//---------------------Add Parnter -----------------------------------------

	apiRouter.HandleFunc("/createPartner", func(w http.ResponseWriter, r *http.Request) {
		t1.CreatePartner(w, r, db)
	})

	//---------------------Get Gf sessions -----------------------------------------

	apiRouter.HandleFunc("/getGFSessionData1", func(w http.ResponseWriter, r *http.Request) {
		t1.GetGFSessionData1(w, r, db)
	})

	//---------------------Get Green Motivators-----------------------------------------

	apiRouter.HandleFunc("/getEnrollGreenMotivators", func(w http.ResponseWriter, r *http.Request) {
		t1.GetEnrollGreenMotivators(w, r, db)
	})

	//---------------------Get PoaTA -----------------------------------------

	apiRouter.HandleFunc("/getPoaTa", func(w http.ResponseWriter, r *http.Request) {
		t1.GetPoaTa(w, r, db)
	})

	//---------------------Get Notes  -----------------------------------------

	apiRouter.HandleFunc("/getNotes", func(w http.ResponseWriter, r *http.Request) {
		t1.GetNotes(w, r, db)
	})

	//---------------------Get  my TeamQAF -----------------------------------------

	apiRouter.HandleFunc("/getMyTeamQAF", func(w http.ResponseWriter, r *http.Request) {
		t1.GetMayTeamQAF(w, r, db)
	})

	//--------------------- Participants Attendance -----------------------------------------

	apiRouter.HandleFunc("/participantsAttendance", func(w http.ResponseWriter, r *http.Request) {
		t1.ParticipantsAttendance(w, r, db)
	})

	//_______________________________TEAM_2___________________________________

	//=============================== CHETAN API'S ================================

	apiRouter.HandleFunc("/getBuzzVyaparProgramBaseline", func(w http.ResponseWriter, r *http.Request) {
		t2.GetBuzzVyaparProgramBaseline(w, r, db)
	})
	apiRouter.HandleFunc("/getEnrollVyaparEnrollment", func(w http.ResponseWriter, r *http.Request) {
		t2.GetEnrollVyaparEnrollment(w, r, db)
	})
	apiRouter.HandleFunc("/taAttachments", func(w http.ResponseWriter, r *http.Request) {
		t2.TaAttachments(w, r, db)
	})
	apiRouter.HandleFunc("/updateRescheduleEvent", func(w http.ResponseWriter, r *http.Request) {
		t2.UpdateRescheduleEvent(w, r, db)
	})
	apiRouter.HandleFunc("/editParticipant", func(w http.ResponseWriter, r *http.Request) {
		t2.EditParticipantMehtod(w, r, db)
	})
	apiRouter.HandleFunc("/editGFSession", func(w http.ResponseWriter, r *http.Request) {
		t2.EditGFSessionMethod(w, r, db)
	})
	apiRouter.HandleFunc("/checkInOut", func(w http.ResponseWriter, r *http.Request) {
		t2.CheckInOut(w, r, db)
	})
	apiRouter.HandleFunc("/signIn", func(w http.ResponseWriter, r *http.Request) {
		t2.SignInHandleRequest(w, r, db)
	})

	//====================== vishal pimple ==============================
	apiRouter.HandleFunc("/getPOA", func(w http.ResponseWriter, r *http.Request) {
		t2.GetPoa1(w, r, db)
	})

	apiRouter.HandleFunc("/deleteEmpFromProject", func(w http.ResponseWriter, r *http.Request) {
		t2.DeleteEmpFromProject(w, r, db)
	})

	apiRouter.HandleFunc("/getDemoGraphy", func(w http.ResponseWriter, r *http.Request) {
		t2.GetDemoGraphy(w, r, db)
	})

	apiRouter.HandleFunc("/createProject", func(w http.ResponseWriter, r *http.Request) {
		t2.CreateProject(w, r, db)
	})
	apiRouter.HandleFunc("/getGFAssignedBatch", func(w http.ResponseWriter, r *http.Request) {
		t2.GetGFAssignedBatch(w, r, db)
	})

	apiRouter.HandleFunc("/createGFSessions", func(w http.ResponseWriter, r *http.Request) {
		t2.CreateGFSessions(w, r, db)
	})

	apiRouter.HandleFunc("/getGFSessionsNew", func(w http.ResponseWriter, r *http.Request) {
		t2.GetGFSessionsNew(w, r, db)
	})

	//=================================== ANEEL CODES =========================================
	apiRouter.HandleFunc("/teamMembers", func(w http.ResponseWriter, r *http.Request) {
		t2.TeamMembers(w, r, db)
	})
	apiRouter.HandleFunc("/addSurveyData", func(w http.ResponseWriter, r *http.Request) {
		t2.AddSurveydata(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF---------------------------------------
	apiRouter.HandleFunc("/getBuses", func(w http.ResponseWriter, r *http.Request) {
		t2.GetBuses(w, r, db)
	})
	///------------------------- delete Buses --------------------------
	apiRouter.HandleFunc("/deleteBus", func(w http.ResponseWriter, r *http.Request) {
		t2.DeleteBus(w, r, db)
	})
	apiRouter.HandleFunc("/deleteUser", func(w http.ResponseWriter, r *http.Request) {
		t2.DeletUser(w, r, db)
	})

	//--------------------------- getvyapar enrollment --------------------------------
	apiRouter.HandleFunc("/getVyaparEnrollment", func(w http.ResponseWriter, r *http.Request) {
		t2.GetvyaparEnrollment(w, r, db)
	})

	//--------------  end pont for get Alter bus ------------------
	apiRouter.HandleFunc("/getAlterBus", func(w http.ResponseWriter, r *http.Request) {
		t2.GetAlterBus(w, r, db)
	})

	apiRouter.HandleFunc("/getEmpProjects", func(w http.ResponseWriter, r *http.Request) {
		t2.GetEmppPojects(w, r, db)
	})

	//-----------------------getdriverlist--------------------------
	apiRouter.HandleFunc("/getDriverList", func(w http.ResponseWriter, r *http.Request) {
		t2.DriverList(w, r, db)
	})

	//-----------------------CreateProgramParticipant--------------------------
	apiRouter.HandleFunc("/createProgramParticipant", func(w http.ResponseWriter, r *http.Request) {
		t2.CreateProgramParticipant(w, r, db)
	})

	//_______________________________TEAM_3___________________________________

	//-------------------Endpoint for adding data to QAF---------------------------------------
	apiRouter.HandleFunc("/getPartnerList", func(w http.ResponseWriter, r *http.Request) {
		t3.GetPartnerList(w, r, db)
	})

	//-------------------Endpoint for adding data to QAF---------------------------------------
	apiRouter.HandleFunc("/getFunderList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetFunderList(w, r, db)
	})

	//-------------------Endpoint for adding data to QAF---------------------------------------
	apiRouter.HandleFunc("/getBusList", func(w http.ResponseWriter, r *http.Request) {
		t3.GetBusList(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getTrainingBatchData", func(w http.ResponseWriter, r *http.Request) {
		t3.GetTrainingBatchData(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/editTrainingBatch", func(w http.ResponseWriter, r *http.Request) {
		t3.UpdateTrainingBatch(w, r, db)
	})

	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/updateReschedule", func(w http.ResponseWriter, r *http.Request) {
		t3.UpdateReschedule(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getgfl", func(w http.ResponseWriter, r *http.Request) {
		t3.Getgfl(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getPeopleList", func(w http.ResponseWriter, r *http.Request) {
		t3.GetPeopleList(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/editUser", func(w http.ResponseWriter, r *http.Request) {
		t3.EditUser(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/setEnrollGelathi", func(w http.ResponseWriter, r *http.Request) {
		t3.SetEnrollGelathi(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/updatePoaCancel", func(w http.ResponseWriter, r *http.Request) {
		t3.UpdatePoaCancel(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/uploadGFSessionPhotos", func(w http.ResponseWriter, r *http.Request) {
		t3.UploadGFsessionPhotos(w, r, db)
	})
	//-------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getCheckInOutStatus", func(w http.ResponseWriter, r *http.Request) {
		t3.GetCheckInOutStatus(w, r, db)
	})
	// -------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getPeopleFilters", func(w http.ResponseWriter, r *http.Request) {
		t3.GetPeopleFilters(w, r, db)
	})
	// -------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getBusCheckList", func(w http.ResponseWriter, r *http.Request) {
		t3.GetBusCheckList(w, r, db)
	})

	// -------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getChangeRole", func(w http.ResponseWriter, r *http.Request) {
		t3.GetChangeRole(w, r, db)
	})
	// -------------------Endpoint for adding data to QAF--------------------------------
	apiRouter.HandleFunc("/getProjectData", func(w http.ResponseWriter, r *http.Request) {
		t3.GetProjectData(w, r, db)
	})
	// -------------------Endpoint for adding Bus--------------------------------
	apiRouter.HandleFunc("/createBus", func(w http.ResponseWriter, r *http.Request) {
		t3.CreateBus(w, r, db)
	})
	// -------------------Endpoint for adding getProjects--------------------------------
	apiRouter.HandleFunc("/getProjects", func(w http.ResponseWriter, r *http.Request) {
		t3.GetProjects(w, r, db)
	})
	// -------------------Endpoint for adding getTa--------------------------------
	apiRouter.HandleFunc("/getTa", func(w http.ResponseWriter, r *http.Request) {
		t3.GetTa(w, r, db)
	})

	//_______________________________TEAM_4___________________________________

	//-------------------------------Sushmitha---------------------------------
	//---------------------------Get Bus Data----------------------------------
	apiRouter.HandleFunc("/getBusData", func(w http.ResponseWriter, r *http.Request) {
		t4.GetBusData(w, r, db)
	})
	//---------------------------Get Assign Targets----------------------------------
	apiRouter.HandleFunc("/getAssignTargets", func(w http.ResponseWriter, r *http.Request) {
		t4.GetAssignTargets(w, r, db)
	})
	//---------------------------Get Occupations----------------------------------
	apiRouter.HandleFunc("/getOccupations", func(w http.ResponseWriter, r *http.Request) {
		t4.GetOccupations(w, r, db)
	})
	//---------------------------Get All Buzz Team----------------------------------
	apiRouter.HandleFunc("/getAllBuzzTeam", func(w http.ResponseWriter, r *http.Request) {
		t4.GetAllBuzzTeam(w, r, db)
	})
	//---------------------------Get OperationsManager List----------------------------------
	apiRouter.HandleFunc("/getOperationsManagerList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetOperationsManagerList(w, r, db)
	})
	//---------------------------Get Gelathi List----------------------------------
	apiRouter.HandleFunc("/getGelathiList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetGelathiList(w, r, db)
	})
	//---------------------------Get Trainers List----------------------------------
	apiRouter.HandleFunc("/getTrainersList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetTrainersList(w, r, db)
	})

	//---------------------------Create Event----------------------------------
	apiRouter.HandleFunc("/createEvent", func(w http.ResponseWriter, r *http.Request) {
		t4.CreateEvent(w, r, db)
	})
	//---------------------------Create Notes----------------------------------
	apiRouter.HandleFunc("/createNotes", func(w http.ResponseWriter, r *http.Request) {
		t4.CreateNotes(w, r, db)
	})
	//---------------------------Create Circle----------------------------------
	apiRouter.HandleFunc("/createCircle", func(w http.ResponseWriter, r *http.Request) {
		t4.CreateCircle(w, r, db)
	})
	//---------------------------Create GF Batch----------------------------------
	apiRouter.HandleFunc("/createGFBatch", func(w http.ResponseWriter, r *http.Request) {
		t4.CreateGFBatch(w, r, db)
	})
	//---------------------------roles list----------------------------------
	apiRouter.HandleFunc("/roles_list", func(w http.ResponseWriter, r *http.Request) {
		t4.Roles_list(w, r, db)
	})
	//---------------------------Create Trainer Target----------------------------------
	apiRouter.HandleFunc("/createTrainerTarget", func(w http.ResponseWriter, r *http.Request) {
		t4.CreateTrainerTarget(w, r, db)
	})
	//---------------------------New/Verify Ta----------------------------------
	apiRouter.HandleFunc("/verifyTa", func(w http.ResponseWriter, r *http.Request) {
		t4.VerifyTa(w, r, db)
	})
	//---------------------------Consume Stock----------------------------------
	apiRouter.HandleFunc("/consumeStock", func(w http.ResponseWriter, r *http.Request) {
		t4.ConsumeStock(w, r, db)
	})
	//---------------------------Remove Enroll Gelathi----------------------------------
	apiRouter.HandleFunc("/removeEnrollGelathi", func(w http.ResponseWriter, r *http.Request) {
		t4.RemoveEnrollGelathi(w, r, db)
	})
	//---------------------------Update Participant Day----------------------------------
	apiRouter.HandleFunc("/updateParticipantDay", func(w http.ResponseWriter, r *http.Request) {
		t4.UpdateParticipantDay(w, r, db)
	})
	//---------------------------New/Remove green motivators---------------------------------
	apiRouter.HandleFunc("/removeGreenMotivators", func(w http.ResponseWriter, r *http.Request) {
		t4.RemoveGreenMotivators(w, r, db)
	})
	//---------------------------New/Remove vyapar enrollment---------------------------------
	apiRouter.HandleFunc("/removeVyaparEnrollment", func(w http.ResponseWriter, r *http.Request) {
		t4.RemoveVyaparEnrollment(w, r, db)
	})
	//---------------------------Upload Event Photos---------------------------------
	apiRouter.HandleFunc("/uploadEventPhotos", func(w http.ResponseWriter, r *http.Request) {
		t4.UploadEventPhotos(w, r, db)
	})
	//---------------------------New/ApproveTa---------------------------------
	apiRouter.HandleFunc("/approveTa", func(w http.ResponseWriter, r *http.Request) {
		t4.ApproveTa(w, r, db)
	})

	//---------------------------Add Bus CheckList---------------------------------
	apiRouter.HandleFunc("/addBusCheckList", func(w http.ResponseWriter, r *http.Request) {
		t4.AddBusCheckList(w, r, db)
	})

	//---------------------------Delet GF Batch--------------------------------
	apiRouter.HandleFunc("/deleteGFBatch", func(w http.ResponseWriter, r *http.Request) {
		t4.DeleteGFBatch(w, r, db)
	})

	//---------------------------Dhiraj Lakhane---------------------------------
	//----------------------------Get Location----------------------------------
	apiRouter.HandleFunc("/getLocation", func(w http.ResponseWriter, r *http.Request) {
		t4.GetLocations(w, r, db)
	})

	//----------------------------get participant Data----------------------------------
	apiRouter.HandleFunc("/getParticipantData", func(w http.ResponseWriter, r *http.Request) {
		t4.GetParticipantData(w, r, db)
	})

	//----------------------------Keerthana------------------------------------
	//----------------------------Get Project List----------------------------------
	apiRouter.HandleFunc("/getProjectList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetProjectList(w, r, db)
	})

	//----------------------------Get Location Name----------------------------------
	apiRouter.HandleFunc("/getlocationName", func(w http.ResponseWriter, r *http.Request) {
		t4.GetLocationName(w, r, db)
	})

	//----------------------------Get Training Batch----------------------------------
	apiRouter.HandleFunc("/getTrainingBatch", func(w http.ResponseWriter, r *http.Request) {
		t4.GetTrainingBatch(w, r, db)
	})

	//----------------------------new/setGreenMotivators----------------------------------
	apiRouter.HandleFunc("/setGreenMotivators", func(w http.ResponseWriter, r *http.Request) {
		t4.SetGreenMotivator(w, r, db)
	})

	//----------------------------new/getGreenMotivators----------------------------------
	apiRouter.HandleFunc("/getGreenMotivators", func(w http.ResponseWriter, r *http.Request) {
		t4.GetGreenMotivator(w, r, db)
	})

	//------------------------------Prathamesh------------------------------------
	//----------------------------Get Event Detail----------------------------------
	apiRouter.HandleFunc("/getEventDetail", func(w http.ResponseWriter, r *http.Request) {
		t4.GetEventDetailist(w, r, db)
	})

	//----------------------------Get Training Batch List---------------------------------
	apiRouter.HandleFunc("/getTrainingBatchList", func(w http.ResponseWriter, r *http.Request) {
		t4.GetBatchlist(w, r, db)
	})

	//_______________________________TEAM_5___________________________________

	//---------------------GetCaste----------------------------------
	apiRouter.HandleFunc("/getCaste", func(w http.ResponseWriter, r *http.Request) {
		t5.GetCaste(w, r, db)
	})
	//---------------------GetEducation----------------------------------
	apiRouter.HandleFunc("/getEducation", func(w http.ResponseWriter, r *http.Request) {
		t5.GetEducation(w, r, db)
	})
	//---------------------GetVillageList----------------------------------
	apiRouter.HandleFunc("/getVillageList", func(w http.ResponseWriter, r *http.Request) {
		t5.GetVillageList(w, r, db)
	})
	//---------------------GetEnrollGelathi----------------------------------
	apiRouter.HandleFunc("/getEnrollGelathi", func(w http.ResponseWriter, r *http.Request) {
		t5.GetEnrollGelathi(w, r, db)
	})
	//---------------------AddNewTA----------------------------------
	apiRouter.HandleFunc("/addNewTA", func(w http.ResponseWriter, r *http.Request) {
		t5.AddNewTA(w, r, db)
	})
	//---------------------DeleteTa----------------------------------
	apiRouter.HandleFunc("/deleteTa", func(w http.ResponseWriter, r *http.Request) {
		t5.DeleteTa(w, r, db)
	})
	//---------------------ListTa----------------------------------
	apiRouter.HandleFunc("/listTa", func(w http.ResponseWriter, r *http.Request) {
		t5.ListTa(w, r, db)
	})
	//---------------------UpdateTa----------------------------------
	apiRouter.HandleFunc("/updateTa", func(w http.ResponseWriter, r *http.Request) {
		t5.UpdateTa(w, r, db)
	})
	//---------------------AddSpoorthiBaselineQuestionnaire----------------------------------
	apiRouter.HandleFunc("/addSpoorthiBaselineQuestionnaire", func(w http.ResponseWriter, r *http.Request) {
		t5.AddSpoorthiBaselineQuestionnaire(w, r, db)
	})
	//---------------------AddGreenBaselineSurvey----------------------------------
	apiRouter.HandleFunc("/addGreenBaselineSurvey", func(w http.ResponseWriter, r *http.Request) {
		t5.AddGreenBaselineSurvey(w, r, db)
	})
	//---------------------UploadTrainingPhotos----------------------------------
	apiRouter.HandleFunc("/uploadTrainingPhotos", func(w http.ResponseWriter, r *http.Request) {
		t5.UploadTrainingPhotos(w, r, db)
	})
	//---------------------GetTotalStocks-------------------------------------------
	apiRouter.HandleFunc("/getTotalstocks", func(w http.ResponseWriter, r *http.Request) {
		t5.GetTotalStocks(w, r, db)
	})
	//---------------------GetMyTeam-------------------------------------------
	apiRouter.HandleFunc("/getMyTeam", func(w http.ResponseWriter, r *http.Request) {
		t5.GetMyTeam(w, r, db)
	})
	//---------------------CreateParticipant-------------------------------------------
	apiRouter.HandleFunc("/createParticipant", func(w http.ResponseWriter, r *http.Request) {
		t5.CreateParticipant(w, r, db)
	})
	//---------------------GetEmailExits-------------------------------------------
	apiRouter.HandleFunc("/getEmailExits", func(w http.ResponseWriter, r *http.Request) {
		t5.GetEmailExits(w, r, db)
	})
	//---------------------UpdateEnrolledGelathi-------------------------------------------
	apiRouter.HandleFunc("/updateEnrolledGelathi", func(w http.ResponseWriter, r *http.Request) {
		t5.UpdateEnrolledGelathi(w, r, db)
	})
	//---------------------UpdateEventCancel-------------------------------------------
	apiRouter.HandleFunc("/updateEventCancel", func(w http.ResponseWriter, r *http.Request) {
		t5.UpdateEventCancel(w, r, db)
	})
	//---------------------GetAllPeople-------------------------------------------
	apiRouter.HandleFunc("/getAllPeople", func(w http.ResponseWriter, r *http.Request) {
		t5.GetAllPeople(w, r, db)
	})
	// -------------------Endpoint for getting stock itmes ---------------------------------------
	apiRouter.HandleFunc("/getStockItems", func(w http.ResponseWriter, r *http.Request) {
		t5.GetStockItems(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint for getting gelathi circle ---------------------------------------
	apiRouter.HandleFunc("/getGelathiCircle", func(w http.ResponseWriter, r *http.Request) {
		t5.GetGelathiCircle(w, r, db)
	}).Methods("POST")
	// -------------------Endpoint for getting GetGFSessionData ---------------------------------------

	apiRouter.HandleFunc("/getGFSessionData", func(w http.ResponseWriter, r *http.Request) {
		t5.GetGFSessionData(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint for Editing Bus ---------------------------------------
	apiRouter.HandleFunc("/editBus", func(w http.ResponseWriter, r *http.Request) {
		t5.EditBus(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint for Set Vyapar Enrollment ---------------------------------------
	apiRouter.HandleFunc("/setVyaparEnrollment", func(w http.ResponseWriter, r *http.Request) {
		t5.SetVyaparEnrollment(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint for Add Emp to Project ---------------------------------------
	apiRouter.HandleFunc("/addEmpToProject", func(w http.ResponseWriter, r *http.Request) {
		t5.AddEmpToProject(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint to create training batch ---------------------------------------
	apiRouter.HandleFunc("/createTrainingBatch", func(w http.ResponseWriter, r *http.Request) {
		t5.CreateTrainingBatch(w, r, db)
	}).Methods("POST")

	// -------------------Endpoint to create GF Sessions New 1---------------------------------------
	apiRouter.HandleFunc("/createGFSessionsNew1", func(w http.ResponseWriter, r *http.Request) {
		t5.CreateGFSessionNew1(w, r, db)
	}).Methods("POST")

	handler := cors.Default().Handler(apiRouter)
	log.Println(http.ListenAndServe(":8080", handler))

}
