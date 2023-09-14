package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type request1 struct {
	BatchID string `json:"batch_id"`
	RoleID  string `json:"role_id"`
}
type Output struct {
	Data                    Response12    `json:"data"`
	AllParticipants         []Participant `json:"all_participants"`
	Enroll_participant_id   string        `json:"enroll_participant_id"`
	Enroll_participant_name string        `json:"enroll_participant_name"`
	Enroll_village_name     string        `json:"enroll_village_name"`
	Enroll_village_contact  string        `json:"enroll_village_contact"`
	Suggested               []interface{} `json:"suggested"`
	DayGelathi              int           `json:"dayGelathi"`
	Photo                   []Photo       `json:"photos"`
	Evaluation_first        string        `json:"evaluation_first" `
	Evaluation_second       string        `json:"evaluation_second"`
	Total_participants      int           `json:"total_participants"`
	CountOfSurvey           string        `json:"countOfSurvey"`
	Code                    int           `json:"code"`
	Success                 bool          `json:"success"`
	Message                 string        `json:"message"`
}

type Response12 struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Participants   string `json:"participants"`
	ProjectID      string `json:"project_id"`
	ProjectName    string `json:"projectName"`
	PartnerName    string `json:"partnerName"`
	ContactPerson  string `json:"contact_person"`
	ContactNumber  string `json:"contact_number"`
	Dates          string `json:"dates"`
	SubVillage     string `json:"sub_village"`
	LocationID     string `json:"location_id"`
	LocationName   string `json:"location_name"`
	TalukID        string `json:"taluk_id"`
	PrimaryID      string `json:"primary_id"`
	UserID         string `json:"user_id"`
	Day1           string `json:"day1"`
	Day1Actual     string `json:"day1_actual"`
	Day1Status     int    `json:"day1Status"`
	Day1ID         string `json:"day1_id"`
	Day1CheckIn    string `json:"day1_checkIn"`
	Day1CheckOut   string `json:"day1_checkOut"`
	Day2           string `json:"day2"`
	Day2Actual     string `json:"day2_actual"`
	Day2Status     int    `json:"day2Status"`
	Day2Id         string `json:"day2_id"`
	Day2CheckIn    string `json:"day2_checkIn"`
	Day2CheckOut   string `json:"day2_checkOut"`
	BatchCompleted string `json:"batch_completed"`
	TrainerName    string `json:"trainer_name"`
	Day1Completed  string `json:"day1_completed"`
	Day2Completed  string `json:"day2_completed"`
}
type Photo struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

type Participant struct {
	// ID               int    `json:"id"`
	Day2             string `json:"day2"`
	Enroll           string `json:"enroll"`
	Final_save       string `json:"final_save"`
	GelathiRecomm    string `json:"gelathiRecomm"`
	IsSurveyDone     bool   `json:"isSurveyDone"`
	Participant_id   string `json:"participant_id"`
	Participant_name string `json:"participant_name"`
}
type Batch1 struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Participants   string `json:"participants"`
	ProjectID      string `json:"project_id"`
	ProjectName    string `json:"projectName"`
	PartnerName    string `json:"partnerName"`
	ContactPerson  string `json:"contact_person"`
	ContactNumber  string `json:"contact_number"`
	Dates          string `json:"dates"`
	SubVillage     string `json:"sub_village"`
	LocationID     string `json:"location_id"`
	LocationName   string `json:"location_name"`
	TalukID        string `json:"taluk_id"`
	PrimaryID      string `json:"primary_id"`
	UserID         string `json:"user_id"`
	Day1           string `json:"day1"`
	Day1Actual     string `json:"day1_actual"`
	Day1Status     int    `json:"day1Status"`
	Day1ID         string `json:"day1_id"`
	Day1CheckIn    string `json:"day1_checkIn"`
	Day1CheckOut   string `json:"day1_checkOut"`
	Day2           string `json:"day2"`
	BatchCompleted string `json:"batch_completed"`
	Day1Completed  string `json:"day1_completed"`
	Day2Status     int    `json:"day2Status"`
	Day2Actual     string `json:"day2_actual"`
	Day2CheckIn    string `json:"day2_checkIn"`
	Day2CheckOut   string `json:"day2_checkOut"`
	Day2Completed  string `json:"day2_completed"`
	Day2Id         string `json:"day2_id"`
	TrainerName    string `json:"trainer_name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
}

func checkDateStatus(date string) int {
	today := time.Now().Format("01/02/2006")
	todayTime, _ := time.Parse("01/02/2006", today)
	dateTime, _ := time.Parse("2006-01-02", date)

	if todayTime.Equal(dateTime) {
		return 1
	} else {
		return 0
	}
}

func checkDay(id int, primaryID int, datesVal []string) []string {
	if id == primaryID {
		return datesVal
	} else {
		day := make([]string, 8)
		day[4] = datesVal[0]
		// day[5] = datesVal[1]
		// day[6] = datesVal[2]
		// day[7] = datesVal[3]

		return day
	}
}

// func formatDate(date string) string {
// 	// Convert the date string to a `time.Time` object.
// 	t, err := time.Parse("2006-01-02", date)
// 	if err != nil {
// 		return ""
// 	}

// 	// Format the `time.Time` object to the desired format.
// 	formattedDate := t.Format("MM/DD/YYYY")

// 	// Return the formatted date string.
// 	return formattedDate
// }

func GetTrainingBatchData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var res Response12
	var pooh []Photo
	var p request1
	err1 := json.NewDecoder(r.Body).Decode(&p)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err1})
		return
	}

	id := p.BatchID
	if p.BatchID == "" {
		errorResponse := ErrorResponse{Message: "BatchID is required"}
		sendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}
	if p.RoleID == "" || p.RoleID == "0" {
		errorResponse := ErrorResponse{Message: "RoleID is required"}
		sendJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// if p.RoleID

	var roleExists bool
	err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM employee WHERE empRole = ?)", p.RoleID).Scan(&roleExists)
	if err != nil {

		// Print the error and return a JSON response with error message

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal server error"})
		return
	}

	if !roleExists {
		// Print the RoleID not found message

		// Return a JSON response if RoleID doesn't exist
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "No role_id found"})
		return
	}

	getUrl := "https://bdms.buzzwomen.org/appTest/"
	condition := ""

	fields := "COALESCE(trn_batch.tb_id,'') as id,COALESCE( UPPER(trn_batch.name),'') as name,COALESCE( trn_batch.participants,''),COALESCE( prj.id,'') as project_id, COALESCE(UPPER(prj.projectName),'') as projectName,COALESCE( IFNULL(part.partnerName, ''),'') as partnerName,COALESCE( IFNULL(trn_batch.contact_person, '') ,'')as contact_person, COALESCE(IFNULL(trn_batch.contact_number, ''),'') as contact_number,COALESCE( GROUP_CONCAT(DISTINCT CONCAT(DATE_FORMAT(trn_batch.date, '%d-%m-%Y %h:%i %p'),',', trn_batch.id,',',IF(trn_batch.check_in != '', 1, 0),',',IF(trn_batch.check_out != '', 1, 0)) ),'') as dates, COALESCE(trn_batch.sub_village,''), COALESCE(trn_batch.location_id,''), COALESCE(loc.name,'') as location_name, COALESCE(taluk.id,'') as taluk_id,COALESCE( primary_id,''), COALESCE(trn_batch.user_id,'')"
	query := fmt.Sprintf("SELECT %s FROM tbl_poa trn_batch LEFT JOIN project prj ON trn_batch.project_id = prj.id LEFT JOIN partner part ON prj.partnerID = part.partnerID LEFT JOIN location loc ON trn_batch.location_id = loc.id LEFT JOIN location taluk ON loc.parentId = taluk.id JOIN employee emp ON trn_batch.user_id = emp.id WHERE trn_batch.tb_id = '%s' AND trn_batch.type = '1' %s", fields, id, condition)

	rows, err := DB.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request....", "Message": err})

		return
	}

	defer rows.Close()
	var batch Batch1
	if rows.Next() {

		err := rows.Scan(
			&batch.ID,
			&batch.Name,
			&batch.Participants,
			&batch.ProjectID,
			&batch.ProjectName,
			&batch.PartnerName,
			&batch.ContactPerson,
			&batch.ContactNumber,
			&batch.Dates,
			&batch.SubVillage,
			&batch.LocationID,
			&batch.LocationName,
			&batch.TalukID,
			&batch.PrimaryID,
			&batch.UserID,
		)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request123", "Message": err})
			return
		}

		// Process batch data

		// datesVal := strings.Split(batch.Dates, ",")

		// Process primary_id and dates_val
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Batch not found"})
		return

		// log.Println("Batch not found")
	}

	// Get trainer name who created the batch
	trainerQuery := fmt.Sprintf("SELECT CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS trainer_name FROM tbl_poa JOIN employee emp ON tbl_poa.user_id = emp.id WHERE tbl_poa.tb_id = '%s' AND tbl_poa.added = '0'", id)

	trainerResult, err := DB.Query(trainerQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request111", "Message": err})
		return
	}
	defer trainerResult.Close()

	var trainerName string
	if trainerResult.Next() {
		err := trainerResult.Scan(&trainerName)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
			return
		}
	}

	// Retrieve batch details
	batchQuery := fmt.Sprintf("SELECT COALESCE( trn_batch.tb_id,'') as id, COALESCE(UPPER(trn_batch.name),'') as name, COALESCE(trn_batch.participants,''),COALESCE( prj.id ,'')as project_id, COALESCE(UPPER(prj.projectName),'') as projectName,COALESCE( IFNULL(part.partnerName, ''),'') as partnerName,COALESCE( IFNULL(trn_batch.contact_person, ''),'') as contact_person,COALESCE( IFNULL(trn_batch.contact_number, ''),'') as contact_number, COALESCE(GROUP_CONCAT(DISTINCT CONCAT(DATE_FORMAT(trn_batch.date, '%%d-%%m-%%Y %%h:%%i %%p'), ',', trn_batch.id, ',', IF(trn_batch.check_in != '', 1, 0), ',', IF(trn_batch.check_out != '', 1, 0))),'') as dates,COALESCE( trn_batch.sub_village,''), COALESCE(trn_batch.location_id,''), COALESCE(loc.name,'') as location_name, COALESCE(taluk.id,'') as taluk_id, COALESCE(primary_id,''), COALESCE(trn_batch.user_id,'') FROM tbl_poa trn_batch LEFT JOIN project prj ON trn_batch.project_id = prj.id LEFT JOIN partner part ON prj.partnerID = part.partnerID LEFT JOIN location loc ON trn_batch.location_id = loc.id LEFT JOIN location taluk ON loc.parentId = taluk.id JOIN employee emp ON trn_batch.user_id = emp.id WHERE trn_batch.tb_id = '%s' AND trn_batch.type = '1'", id)

	batchResult, err := DB.Query(batchQuery)
	if err != nil {
		log.Println("ERROR>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	defer batchResult.Close()

	if batchResult.Next() {
		var batch Batch1
		err := batchResult.Scan(
			&batch.ID,
			&batch.Name,
			&batch.Participants,
			&batch.ProjectID,
			&batch.ProjectName,
			&batch.PartnerName,
			&batch.ContactPerson,
			&batch.ContactNumber,
			&batch.Dates,
			&batch.SubVillage,
			&batch.LocationID,
			&batch.LocationName,
			&batch.TalukID,
			&batch.PrimaryID,
			&batch.UserID,
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
			return
		}
		datesVal := strings.Split(batch.Dates, ",")
		var dates []string
		// var batch Batch1
		var day1 int
		if id == "" {
			day1 = 0

		} else {

			day1, err = strconv.Atoi(id)
			if err != nil {
				// Handle the error if the conversion fails
				return
			}

		}

		if datesVal[0] != "" && datesVal[4] != "" {
			dates = datesVal
			if timeVal, err := time.Parse("2006-01-02", dates[0]); err == nil {
				timeVal2, err := time.Parse("2006-01-02", dates[4])
				if err == nil && timeVal.After(timeVal2) {
					pieces := make([][]string, 0)
					pieces = append(pieces, dates[:len(dates)/2])
					pieces = append(pieces, dates[len(dates)/2:])
					dates = make([]string, 0)
					for _, piece := range pieces {
						dates = append(dates, piece...)
					}
				}
			}
		} else {

			pid, err := strconv.Atoi(batch.PrimaryID)

			if err != nil {
				fmt.Print("hudcsiud:", err)
			}
			dates = checkDay(day1, pid, datesVal)
		}

		res.ID = batch.ID
		res.Name = batch.Name
		res.Participants = batch.Participants
		res.ProjectID = batch.ProjectID
		res.ProjectName = batch.ProjectName
		res.PartnerName = batch.PartnerName
		res.ContactPerson = batch.ContactPerson
		res.ContactNumber = batch.ContactNumber
		res.Dates = batch.Dates
		res.SubVillage = batch.SubVillage
		res.LocationID = batch.LocationID
		res.LocationName = batch.LocationName
		res.TalukID = batch.TalukID
		res.PrimaryID = batch.PrimaryID
		res.UserID = batch.UserID
		res.Day1 = batch.Day1
		res.Day1Actual = batch.Day1Actual

		//	strconv.Atoi(batch.Day1Status)
		res.Day1Status = batch.Day1Status
		res.Day1ID = batch.Day1ID
		res.Day1CheckIn = batch.Day1CheckIn
		res.Day1CheckOut = batch.Day1CheckOut
		res.Day2 = batch.Day2
		res.Day2Actual = batch.Day2Actual
		res.Day2Status = batch.Day2Status
		res.Day2Id = batch.Day2Id
		res.Day2CheckIn = batch.Day2CheckIn
		res.Day2CheckOut = batch.Day2CheckOut
		res.BatchCompleted = batch.BatchCompleted
		res.TrainerName = trainerName

		if len(dates) >= 8 {
			res.ID = batch.ID
			res.Name = batch.Name
			res.Participants = batch.Participants
			res.ProjectID = batch.ProjectID
			res.ProjectName = batch.ProjectName
			res.PartnerName = batch.PartnerName
			res.ContactPerson = batch.ContactPerson
			res.ContactNumber = batch.ContactNumber
			res.Dates = batch.Dates
			res.SubVillage = batch.SubVillage
			res.LocationID = batch.LocationID
			res.LocationName = batch.LocationName
			res.TalukID = batch.TalukID
			res.PrimaryID = batch.PrimaryID
			res.UserID = batch.UserID
			res.Day1 = dates[0]

			res.Day1Actual = dates[0]

			res.Day1Status = (checkDateStatus(dates[4]))
			res.Day1ID = dates[1]
			res.Day1CheckIn = dates[2]
			res.Day1CheckOut = dates[3]
			res.Day2 = dates[4]

			res.Day2Actual = dates[4]

			res.Day2Status = (checkDateStatus(dates[0]))
			res.Day2Id = dates[5]
			res.Day2CheckIn = dates[6]
			res.Day2CheckOut = dates[7]
			res.BatchCompleted = ""
			res.TrainerName = trainerName
			if dates[7] == "1" {
				res.BatchCompleted = "1"
			}
			res.Day1Completed = "0"
			res.Day2Completed = "0"

			if dates[3] == "1" && dates[7] == "0" {
				if dates[3] == "1" {
					res.Day1Completed = "1"
				}
			}

			if dates[3] == "1" && dates[7] == "1" {
				if dates[3] == "1" {
					res.Day1Completed = "1"
				}

				if dates[7] == "1" {
					res.Day2Completed = "1"
				}
			} else {
				if dates[7] == "1" {
					res.Day2Completed = "1"
				}
			}
		}

	}

	// response := make(map[string]interface{})

	// Get participants
	fields = "tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.day2, IF(tr_part.saving_amt = '' or tr_part.saving_amt is null,0,1) as final_save,tr_part.isSurveyDone,tr_part.enroll"
	queryParticipant := fmt.Sprintf("SELECT %s FROM training_participants tr_part WHERE tr_part.tb_id = %s", fields, id)
	resultParticipant, err := DB.Query(queryParticipant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer resultParticipant.Close()

	var participant_id, participant_name, gelathiRecomm, final_save string
	var day2 int
	var isSurveyDone bool
	var enroll string
	var par Participant
	var participants []Participant
	for resultParticipant.Next() {
		// var participant Participant

		err := resultParticipant.Scan(&participant_id, &participant_name, &gelathiRecomm, &day2, &final_save, &isSurveyDone, &enroll)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		par.Participant_id = participant_id
		par.Participant_name = participant_name
		par.GelathiRecomm = gelathiRecomm
		par.Day2 = fmt.Sprint(day2)
		par.Final_save = final_save
		par.IsSurveyDone = isSurveyDone
		par.Enroll = enroll
		participants = append(participants, par)
	}

	var out Output
	out.Data = res
	out.AllParticipants = participants
	// response["all_participants"] = participants
	// fmt.Printf("Response: %+v\n", response)
	// json.NewEncoder(w).Encode(response)

	// Get Enroll gelathi details
	enrollFields := "IFNULL(tr_part_enroll.id,'') as enroll_participant_id, IFNULL(tr_part_enroll.firstName, '') as enroll_participant_name, '' as enroll_village_name, '' as enroll_village_contact"
	queryParticipantEnroll := fmt.Sprintf("SELECT %s FROM training_participants tr_part_enroll WHERE tr_part_enroll.tb_id = %s AND tr_part_enroll.enroll = '1'", enrollFields, id)
	resultParticipantEnroll, err := DB.Query(queryParticipantEnroll)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer resultParticipantEnroll.Close()
	var enroll_participant_id, enroll_participant_name, enroll_village_name, enroll_village_contact string

	if resultParticipantEnroll.Next() {
		err := resultParticipantEnroll.Scan(&enroll_participant_id, &enroll_participant_name, &enroll_village_name, &enroll_village_contact)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
	}

	out.Enroll_participant_id = enroll_participant_id
	out.Enroll_participant_name = enroll_participant_name
	out.Enroll_village_name = enroll_village_name
	out.Enroll_village_contact = enroll_village_contact

	if p.RoleID == "6" {
		suggestGelathi := filterParticipants(participants, func(p Participant) bool {
			return p.GelathiRecomm == "1"
		})

		out.Suggested = make([]interface{}, len(suggestGelathi))
		for i, p := range suggestGelathi {
			out.Suggested[i] = p
		}
	}

	dayGelathi := filterParticipants(participants, func(p Participant) bool {
		return p.Day2 == "1"
	})

	out.DayGelathi = len(dayGelathi)

	jsonData, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	fmt.Println(jsonData)
	//out.DayGelathi = ""
	out.CountOfSurvey = id

	// Get training photos

	// getUrl := https://bdms.buzzwomen.org/appTest/images/

	queryTrainingPhotos := fmt.Sprintf("SELECT IFNULL(CONCAT('"+getUrl+"', '', photo1), '') as photo1, IFNULL(CONCAT('"+getUrl+"', '', photo2), '') as photo2 FROM tbl_poa tr_photos WHERE tr_photos.tb_id = %s AND tr_photos.type = 1 LIMIT 2", id)
	resultTrainingPhotos, err := DB.Query(queryTrainingPhotos)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer resultTrainingPhotos.Close()

	// photos := make([]map[string]interface{}, 0)
	var photo1, photo2 string
	for resultTrainingPhotos.Next() {

		err := resultTrainingPhotos.Scan(&photo1, &photo2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}

		var photo Photo
		photo.Photo1 = photo1
		photo.Photo2 = photo2
		pooh = append(pooh, photo)

	}

	out.Evaluation_first = "0" // Set a default value of zero

	fields = fmt.Sprintf("IF(a.training_batch_id = %s, 1, 0) AS flag_evaluation_first", id)
	queryEvaluationFirst := fmt.Sprintf("SELECT %s FROM program_evaluation_first a WHERE a.training_batch_id = %s", fields, id)
	resultEvaluationFirst, err := DB.Query(queryEvaluationFirst)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer resultEvaluationFirst.Close()

	var flag_evaluation_first int // Use an appropriate data type for the flag
	// out.Evaluation_first= "0"
	firstEvaluation := make(map[string]interface{})
	if resultEvaluationFirst.Next() {
		err := resultEvaluationFirst.Scan(&flag_evaluation_first)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		} else {
			firstEvaluation["flag_evaluation_first"] = flag_evaluation_first
		}
	}

	if val, ok := firstEvaluation["flag_evaluation_first"]; ok {
		out.Evaluation_first = fmt.Sprintf("%d", val.(int)) // Convert to string directly
	}

	// out.Evaluation_second = fmt.Sprint(firstEvaluation["flag_evaluation_second"])

	// TrainingBatch program evaluation second
	out.Evaluation_second = "0"
	var flag_evaluation_second int
	fields = fmt.Sprintf("IF(b.training_batch_id = %s, 1, 0) as flag_evaluation_second", id)
	queryEvaluationSecond := fmt.Sprintf("SELECT %s FROM program_evaluation_second b WHERE b.training_batch_id = %s", fields, id)
	resultEvaluationSecond, err := DB.Query(queryEvaluationSecond)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer resultEvaluationSecond.Close()
	// var flag_evaluation_second int
	secondEvaluation := make(map[string]interface{})
	if resultEvaluationSecond.Next() {
		err := resultEvaluationSecond.Scan(&flag_evaluation_second)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		} else {
			secondEvaluation["flag_evaluation_second"] = flag_evaluation_second
		}
	}

	if val, ok := secondEvaluation["flag_evaluation_second"]; ok {
		out.Evaluation_second = fmt.Sprintf("%d", val.(int)) // Convert to string directly
	}

	out.Total_participants = (len(participants))

	// fmt.Printf("Response: %+v\n", response)
	// json.NewEncoder(w).Encode(response)

	var countOfSurvey int
	err = DB.QueryRow(fmt.Sprintf("SELECT COUNT(id) as countOfSurvey FROM training_participants WHERE isSurveyDone = 1 and tb_id = %s", id)).Scan(&countOfSurvey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	out.Photo = pooh
	out.CountOfSurvey = fmt.Sprint(countOfSurvey)
	out.DayGelathi = len(dayGelathi)
	// out.Evaluation_first = fmt.Sprint(firstEvaluation["flag_evaluation_first"])
	out.Code = http.StatusOK
	out.Success = true
	out.Message = "Successfully"
	// out.Total_participants = len(participants)
	// response["countOfSurvey"] = countOfSurvey

	// response["code"] = 200
	// json.NewEncoder(w).Encode(response)
	json.NewEncoder(w).Encode(out)
}

func filterParticipants(participants []Participant, filterFunc func(Participant) bool) []Participant {
	filtered := make([]Participant, 0)
	for _, p := range participants {
		if filterFunc(p) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
