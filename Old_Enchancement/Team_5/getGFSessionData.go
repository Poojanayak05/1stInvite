package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// struct to combine all the data including participants and photos.
type SessionData struct {
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

// struct to store the participant details.
type Participant struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	GelathiRecomm   string `json:"gelathiRecomm"`
	Enroll          int    `json:"enroll"`
	GelathiStatus   string `json:"gelathi_status"`
	IsAttended      string `json:"is_attended"`
}

// struct to store the session photos.
type SessionPhoto struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

// struct to store the parameters coming in body.
type GFSessionDataRequest struct {
	UserId      string `json:"user_id"`
	GFSessionId string `json:"gf_session_id"`
	RoleId      string `json:"role_id"`
}

func GetGFSessionData(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// req variable refers the GFSessionDataRequest struct.
	var req GFSessionDataRequest

	// dumping the body parameters into the GFSessionDataRequest struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error While Decoding Body", "success": false, "error": err})
		return
	}

	if req.UserId == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required User Id", "success": false})
		return
	}
	if req.GFSessionId == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Gf Session Id", "success": false})
		return
	}
	// if req.RoleId == "" {
	// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotAcceptable, "message": "Please Send the Required Role Id", "success": false})
	// 	return
	// }
	// defining the simple condition variable with empty string.
	condition := ""
	if req.UserId != "" {
		// condition if user id present.
		condition = fmt.Sprintf(" AND gf_sess.user_id = %s ", req.UserId)
	} else if req.UserId == "" && req.GFSessionId == "" {
		// setting and returning the response if user id and GF session id is not present.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Not Found", "success": false})
		return
	}

	sessionID := req.GFSessionId

	// declaring fields to collect from the tbl_poa.
	fields := fmt.Sprintf(` gf_sess.id, gf_sess.name AS gf_session_name, gf_sess.tb_id, DATE_FORMAT(gf_sess.date, '%%d-%%m-%%Y %%h:%%i %%p') AS plan_date, prj.id AS project_id, UPPER(prj.projectName) AS projectName, IFNULL(part.partnerName, '') AS partnerName, IFNULL(tr_bat.name, '') AS training_batch_name, IFNULL(tr_bat.contact_person, '') AS contact_person, IFNULL(tr_bat.contact_number, '') AS contact_number, 
	gf_sess.session_type AS type, IF(gf_sess.check_in != '', 1, 0) as check_in, IF(gf_sess.check_out != '', 1, 0) as check_out,	gf_sess.participants, CASE WHEN gf_sess.session_type = 1 THEN 'Circle Metting' WHEN gf_sess.session_type = 2 THEN 'Village Visit' ELSE 'Beehive Visit' END AS type_name, gf_sess.status, gf_sess.circle_id, gf_sess.primary_id, gf_sess.user_id,CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS gf_name`)

	// query to fetch the above fields data from the tbl_poa based on GF Session id and the condition.
	query := fmt.Sprintf(`SELECT %s FROM tbl_poa AS gf_sess
		LEFT JOIN project AS prj ON gf_sess.project_id = prj.id
		LEFT JOIN partner AS part ON prj.partnerID = part.partnerID
		LEFT JOIN tbl_poa AS tr_bat ON gf_sess.tb_id = tr_bat.tb_id
		JOIN employee AS emp ON gf_sess.user_id = emp.id
		WHERE gf_sess.id = %s AND gf_sess.type = 2 %s`, fields, sessionID, condition)

	// preparing the query to execute.
	stmt, err := db.Prepare(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(response)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Error While Preparing Statement", "success": false, "error": err})
		return
	}

	// executing the query.
	row := stmt.QueryRow()

	// declaring the sessionData variable to refer the SessionData struct.
	sessionData := SessionData{}
	trainerName := ""
	err = row.Scan(
		// scanning the data from the tbl_poa.
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
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Not Found", "success": false})

		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}

	// query to fetch the trainer name from tbl_poa.
	trainerQuery := fmt.Sprintf(`
		SELECT CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS trainer_name
		FROM tbl_poa
		JOIN employee AS emp ON tbl_poa.user_id = emp.id
		WHERE tbl_poa.primary_id = %s AND tbl_poa.added = '0'
	`, sessionData.TbID)

	// scanning the trainer name.
	err = db.QueryRow(trainerQuery).Scan(&trainerName)
	if err != nil && err != sql.ErrNoRows {

		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
		return
	}

	//setting the trainer name to session data struct field.
	sessionData.TrainerName = trainerName

	// setting session completed to 0.
	sessionData.SessionCompleted = "0"
	// if checkout is 1 then setting session completed as 1.
	if sessionData.CheckOut == "1" {
		sessionData.SessionCompleted = "1"
	}

	queryParticipant := ""

	// if type is 1 then fetching data from gelathi circle and training participant table otherwise fetching data from training paricipants table only.
	if sessionData.Type == "1" {
		fields := fmt.Sprintf("tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s from gelathi_circle gel_circle left join training_participants tr_part ON gel_circle.gelathi_id = tr_part.id where tr_part.new_green=0 and tr_part.new_vyapar=0 and gel_circle.circle_id = %s", fields, sessionData.CircleID)
	} else {
		fields := fmt.Sprintf("tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, tr_part.enroll, IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended", sessionData.ID)
		queryParticipant = fmt.Sprintf("SELECT %s from training_participants tr_part where tr_part.new_green=0 and tr_part.new_vyapar=0 and  tr_part.tb_id = %s", fields, sessionData.TbID)
	}

	// preparing the query to execute.
	stmt, err = db.Prepare(queryParticipant)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Error While Preparing Statement", "success": false, "error": err})
		return
	}
	defer stmt.Close()
	// executing the query.
	rows, err := stmt.Query()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
		return
	}
	defer rows.Close()

	// declaring a list variable to store multiple particpants details.
	participants := []Participant{}

	for rows.Next() {
		participant := Participant{}
		// scanning the participant.
		err := rows.Scan(
			&participant.ParticipantID,
			&participant.ParticipantName,
			&participant.GelathiRecomm,
			&participant.Enroll,
			&participant.GelathiStatus,
			&participant.IsAttended,
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
			return
		}

		// appending the participant into the participants list.
		participants = append(participants, participant)
	}

	// setting up the variables to the struct variables.
	sessionData.AllParticipants = participants
	sessionData.TotalParticipants = len(participants)

	// query to fetch data of enrolled participant.
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

	// executing the query.
	err = db.QueryRow(enrollQuery).Scan(
		// scanning the data.
		&sessionData.EnrollParticipantID,
		&sessionData.EnrollParticipantName,
		&sessionData.EnrollParticipantContactNum,
		&sessionData.EnrollParticipantVillageName,
	)
	if err != nil && err != sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Returned No Rows", "success": false, "error": err})
		return
	}
	rawURL := "https://bdms.buzzwomen.org/appTest/"

	// query to fetch photos from tbl_poa based on session id.
	photosQuery := fmt.Sprintf("SELECT IF(IFNULL(photo1, '') = '', '', CONCAT('"+rawURL+"', photo1)) AS photo1, IF(IFNULL(photo2, '') = '', '', CONCAT('"+rawURL+"', photo2)) AS photo2 FROM tbl_poa gf_sess_photos WHERE gf_sess_photos.id = %s LIMIT 1", sessionData.ID)

	// executing the query.
	rows, err = db.Query(photosQuery)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
		return
	}
	defer rows.Close()

	// variable to store photos in the list.
	photos := []SessionPhoto{}
	for rows.Next() {
		photo := SessionPhoto{}
		err := rows.Scan(
			// scanning the photos.
			&photo.Photo1,
			&photo.Photo2,
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database Error", "success": false, "error": err})
			return
		}
		// appending the photos in the list one by one.
		photos = append(photos, photo)
	}

	sessionData.Photos = photos

	// setting the response.
	sessionData.Code = http.StatusOK
	sessionData.Message = "Successfully"
	sessionData.Success = true

	// sending the response.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessionData)

}
