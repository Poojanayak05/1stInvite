package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type GFResponse struct {
	ID                           string         `json:"id"`
	GFSessionName                string         `json:"gf_session_name"`
	TbID                         string         `json:"tb_id"`
	PlanDate                     string         `json:"plan_date"`
	ProjectID                    string         `json:"project_id"`
	ProjectName                  string         `json:"projectName"`
	PartnerName                  string         `json:"partnerName"`
	TrainingBatchName            string         `json:"training_batch_name"`
	ContactPerson                string         `json:"contact_person"`
	ContactNumber                string         `json:"contact_number"`
	Type                         string         `json:"type"`
	CheckIn                      string         `json:"check_in"`
	CheckOut                     string         `json:"check_out"`
	Participants                 string         `json:"participants"`
	TypeName                     string         `json:"type_name"`
	Status                       string         `json:"status"`
	CircleID                     string         `json:"circle_id"`
	PrimaryID                    string         `json:"primary_id"`
	UserID                       string         `json:"user_id"`
	GFName                       string         `json:"gf_name"`
	TrainerName                  string         `json:"trainer_name"`
	SessionCompleted             string         `json:"session_completed"`
	AllParticipants              []Participant  `json:"all_participants"`
	TotalParticipants            int            `json:"total_participants"`
	EnrollParticipantID          string         `json:"enroll_participant_id"`
	EnrollParticipantName        string         `json:"enroll_participant_name"`
	EnrollParticipantContactNum  string         `json:"enroll_participant_contactnum"`
	EnrollParticipantVillageName string         `json:"enroll_participant_villagename"`
	Photos                       []SessionPhoto `json:"photos"`
	Code                         int            `json:"code"`
	Success                      bool           `json:"success"`
	Message                      string         `json:"message"`
	ErrorMessage                 error          `json:"error,omitempty"`
}

type Notfound struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Participant struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	GelathiRecomm   string `json:"gelathiRecomm"`
	Enroll          int    `json:"enroll"`
	GelathiStatus   string `json:"gelathi_status"`
	IsAttended      string `json:"is_attended"`
	Module1         string `json:"module1"`
	Module2         string `json:"module2"`
	Module3         string `json:"module3"`
	Module4         string `json:"module4"`
	Module5         string `json:"module5"`
}

type SessionPhoto struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}
type GFRequest struct {
	UserId      int `json:"user_id"`
	GFSessionId int `json:"gf_session_id"`
}

func GetGFSessionData1(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var not Notfound
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err, "line 85")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad GFRequest", "Message": "Failed to decode request body", "Error": err})
		return
	}
	var response GFResponse

	rawURL := "https://bdms.buzzwomen.org/appTest/"
	condition := ""
	if user_id, ok := request["user_id"]; ok {
		condition = fmt.Sprintf(" and gf_sess.user_id = %v", user_id)
	}
	//id := request["gf_session_id"].(string)

	sessionID := request["gf_session_id"].(string)

	fields := fmt.Sprintf(` gf_sess.id, gf_sess.name AS gf_session_name, gf_sess.tb_id, IFNUll(gf_sess.date, '') AS plan_date, prj.id AS project_id, UPPER(prj.projectName) AS projectName, IFNULL(part.partnerName, '') AS partnerName, IFNULL(tr_bat.name, '') AS training_batch_name, IFNULL(tr_bat.contact_person, '') AS contact_person, IFNULL(tr_bat.contact_number, '') AS contact_number,
	gf_sess.session_type AS type, IF(gf_sess.check_in != '', 1, 0) as check_in, IF(gf_sess.check_out != '', 1, 0) as check_out,	IFNUll(gf_sess.participants,''), CASE WHEN gf_sess.session_type = 1 THEN 'Circle Metting' WHEN gf_sess.session_type = 2 THEN 'Village Visit'  WHEN gf_sess.session_type = '3' THEN 'Beehive Visit'  WHEN gf_sess.session_type = '4' THEN 'Spoorthi Survey'  WHEN gf_sess.session_type = '5' THEN 'Spoorthi Module1'  WHEN gf_sess.session_type = '6' THEN 'Spoorthi Module2'  WHEN gf_sess.session_type = '7' THEN 'Spoorthi Module3'  WHEN gf_sess.session_type = '8' THEN 'Spoorthi Module4'  WHEN gf_sess.session_type = '9' THEN 'Spoorthi Module5'  WHEN gf_sess.session_type = '10' THEN 'Green Survey'  WHEN gf_sess.session_type = '11' THEN 'Green Module1'  WHEN gf_sess.session_type = '12' THEN 'Green Module2'  WHEN gf_sess.session_type = '13' THEN 'Green Module3'  WHEN gf_sess.session_type = '14' THEN 'Green Module4'  WHEN gf_sess.session_type = '15' THEN 'Green Module5'  WHEN gf_sess.session_type = '16' THEN 'Vyapar Survey'  WHEN gf_sess.session_type = '17' THEN 'Vyapar Module1'  WHEN gf_sess.session_type = '18' THEN 'Vyapar Module2'  WHEN gf_sess.session_type = '19' THEN 'Vyapar Module3'  WHEN gf_sess.session_type = '20' THEN 'Vyapar Module4'  WHEN gf_sess.session_type = '21' THEN 'Vyapar Module5' END AS type_name, gf_sess.status, gf_sess.circle_id, gf_sess.primary_id, gf_sess.user_id,CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS gf_name`)

	query := fmt.Sprintf(`SELECT %s FROM tbl_poa AS gf_sess
		LEFT JOIN project AS prj ON gf_sess.project_id = prj.id
		LEFT JOIN partner AS part ON prj.partnerID = part.partnerID
		LEFT JOIN tbl_poa AS tr_bat ON gf_sess.tb_id = tr_bat.tb_id
		JOIN employee AS emp ON gf_sess.user_id = emp.id
		WHERE gf_sess.id = %s AND gf_sess.type = 2 %s`, fields, sessionID, condition)

	stmt, err := DB.Prepare(query)
	if err != nil {

		response.Message = "Error While Preparing Statement"
		response.Code = 500
		response.ErrorMessage = err
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	row := stmt.QueryRow()
	sessionData := GFResponse{}
	trainerName := ""
	err = row.Scan(
		&sessionData.ID,
		&sessionData.GFSessionName,
		&sessionData.TbID,
		&sessionData.PlanDate,
		&sessionData.ProjectID,
		&sessionData.ProjectName,
		&sessionData.PartnerName,
		&sessionData.TrainingBatchName,
		&sessionData.ContactPerson,
		&sessionData.ContactNumber,
		&sessionData.Type,
		&sessionData.CheckIn,
		&sessionData.CheckOut,
		&sessionData.Participants,
		&sessionData.TypeName,
		&sessionData.Status,
		&sessionData.CircleID,
		&sessionData.PrimaryID,
		&sessionData.UserID,
		&sessionData.GFName,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			// response.Message = "Not Found"
			// response.Code = 404
			// json.NewEncoder(w).Encode(response)
			not.Code = 400
			not.Message = "Bad request" + fmt.Sprint(err)
			not.Success = false
			json.NewEncoder(w).Encode(not)

		}
	}

	trainerQuery := fmt.Sprintf(`
		SELECT CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS trainer_name
		FROM tbl_poa
		JOIN employee AS emp ON tbl_poa.user_id = emp.id
		WHERE tbl_poa.primary_id = %s AND tbl_poa.added = '0'
	`, sessionData.TbID)

	err = DB.QueryRow(trainerQuery).Scan(&trainerName)
	if err != nil && err != sql.ErrNoRows {

		response.Message = "Database Error"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}

	sessionData.TrainerName = trainerName
	sessionData.SessionCompleted = "0"
	if sessionData.CheckOut == "1" {
		sessionData.SessionCompleted = "1"
	}

	queryParticipant := ""
	if sessionData.Type == "2" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s from training_participants tr_part where tr_part.new_green=0 and tr_part.new_vyapar=0 and  tr_part.tb_id = %s", fields, sessionData.TbID)
	} else if sessionData.Type == "3" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s from training_participants tr_part where tr_part.new_green=0 and tr_part.new_vyapar=0 and  tr_part.tb_id = %s", fields, sessionData.TbID)

	} else if sessionData.Type == "4" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "5" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "6" {
		fields := fmt.Sprintf("distinct(tr_part.id)as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "7" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended , COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "8" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'') ", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "9" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'') ", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.SpoorthiBaselineQuestionnaire s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "10" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "11" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "12" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "13" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "14" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "15" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.GreenBaselineSurvey s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "16" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "17" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "18" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "19" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendancev WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended, COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "20" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else if sessionData.Type == "21" {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From bdms_staff.participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended ,COALESCE(module1,''),COALESCE(module2,''),COALESCE(module3,''),COALESCE(module4,''),COALESCE(module5,'')", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s FROM bdms_staff.BuzzVyaparProgramBaseline s join bdms_staff.gelathi_circle ge on ge.gelathi_id=s.partcipantId join bdms_staff.training_participants tr_part on s.partcipantId=tr_part.id where ge.circle_id= %s   GROUP BY partcipantId", fields, sessionData.CircleID)
	} else {
		fields := fmt.Sprintf("distinct(tr_part.id) as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s from gelathi_circle gel_circle left join training_participants tr_part ON gel_circle.gelathi_id = tr_part.id, bdms_staff.SpoorthiBaselineQuestionnaire b where tr_part.new_green=0 and tr_part.new_vyapar=0 and gel_circle.circle_id = %s", fields, sessionData.CircleID)
	}

	stmt, err = DB.Prepare(queryParticipant)
	if err != nil {
		response.Message = "Error While Preparing Statement"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		response.Message = "Database Error"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	participants := []Participant{}
	columns, err := rows.Columns()
	if err != nil {
		response.Message = "Database Error"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}

	for rows.Next() {
		participant := Participant{}
		values := make([]interface{}, len(columns))

		for i := range columns {
			switch columns[i] {
			case "participant_id":
				values[i] = &participant.ParticipantID
			case "participant_name":
				values[i] = &participant.ParticipantName
			case "gelathiRecomm":
				values[i] = &participant.GelathiRecomm
			case "enroll":
				values[i] = &participant.Enroll
			case "gelathi_status":
				values[i] = &participant.GelathiStatus
			case "is_attended":
				values[i] = &participant.IsAttended
			case "module1":
				values[i] = &participant.Module1
			case "module2":
				values[i] = &participant.Module2
			case "module3":
				values[i] = &participant.Module3
			case "module4":
				values[i] = &participant.Module4
			case "module5":
				values[i] = &participant.Module5
			default:
				// Ignore other columns or handle them as needed
				values[i] = new(interface{})
			}
		}

		if err := rows.Scan(values...); err != nil {

			response.Message = "Database Error"
			response.Code = 500
			response.ErrorMessage = err
			json.NewEncoder(w).Encode(response)
			return
		}

		participants = append(participants, participant)
	}
	sessionData.AllParticipants = participants
	sessionData.TotalParticipants = len(participants)

	enrollQuery := fmt.Sprintf(`
		SELECT
		tr_part.id AS participant_id,
		IFNULL(tr_part.firstName, '') AS participant_name,
		'' AS contactnum,
		IFNULL(SUBSTRING_INDEX(tbl_poa.name, '_', 1), '') AS villagename
		FROM training_participants AS tr_part
		LEFT JOIN tbl_poa ON tr_part.tb_id = tbl_poa.id
		WHERE tr_part.tb_id = %s AND tr_part.enroll = 1 
	`, sessionData.TbID)

	err = DB.QueryRow(enrollQuery).Scan(
		&sessionData.EnrollParticipantID,
		&sessionData.EnrollParticipantName,
		&sessionData.EnrollParticipantContactNum,
		&sessionData.EnrollParticipantVillageName,
	)
	if err != nil && err != sql.ErrNoRows {

		response.Message = "Database Returned No Rows"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}

	photosQuery := fmt.Sprintf("SELECT IF(IFNULL(photo1, '') = '', '', CONCAT('"+rawURL+"', photo1)) AS photo1, IF(IFNULL(photo2, '') = '', '', CONCAT('"+rawURL+"', photo2)) AS photo2 FROM tbl_poa gf_sess_photos WHERE gf_sess_photos.id = %s LIMIT 1", sessionData.ID)

	rows, err = DB.Query(photosQuery)
	if err != nil {

		response.Message = "Database Error"
		response.Code = 500
		response.ErrorMessage = err
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	photos := []SessionPhoto{}
	for rows.Next() {
		photo := SessionPhoto{}
		err := rows.Scan(
			&photo.Photo1,
			&photo.Photo2,
		)
		if err != nil {

			response.Message = "Database Error"
			response.Code = 500
			response.ErrorMessage = err
			json.NewEncoder(w).Encode(response)
			return
		}
		photos = append(photos, photo)

	}
	response = GFResponse{

		Code:                         http.StatusOK,
		Success:                      true,
		Message:                      "Success",
		ID:                           sessionData.ID,
		GFSessionName:                sessionData.GFSessionName,
		TbID:                         sessionData.TbID,
		PlanDate:                     sessionData.PlanDate,
		ProjectID:                    sessionData.ProjectID,
		ProjectName:                  sessionData.ProjectName,
		PartnerName:                  sessionData.PartnerName,
		TrainingBatchName:            sessionData.TrainingBatchName,
		ContactPerson:                sessionData.ContactPerson,
		ContactNumber:                sessionData.ContactNumber,
		Type:                         sessionData.Type,
		CheckIn:                      sessionData.CheckIn,
		CheckOut:                     sessionData.CheckOut,
		Participants:                 sessionData.Participants,
		TypeName:                     sessionData.TypeName,
		Status:                       sessionData.Status,
		CircleID:                     sessionData.CircleID,
		PrimaryID:                    sessionData.PrimaryID,
		UserID:                       sessionData.UserID,
		GFName:                       sessionData.GFName,
		TrainerName:                  sessionData.TrainerName,
		SessionCompleted:             sessionData.SessionCompleted,
		AllParticipants:              sessionData.AllParticipants,
		TotalParticipants:            sessionData.TotalParticipants,
		EnrollParticipantID:          sessionData.EnrollParticipantID,
		EnrollParticipantName:        sessionData.EnrollParticipantName,
		EnrollParticipantContactNum:  sessionData.EnrollParticipantContactNum,
		EnrollParticipantVillageName: sessionData.EnrollParticipantVillageName,
		Photos:                       photos,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {

		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
