package Team_4

//Done by Keerthana
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type GreenMotivatorsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type Session struct {
	ID                           string        `json:"id"`
	GfSessionName                string        `json:"gf_session_name"`
	TbID                         string        `json:"tb_id"`
	PlanDate                     string        `json:"plan_date"`
	ProjectID                    string        `json:"project_id"`
	ProjectName                  string        `json:"projectName"`
	PartnerName                  string        `json:"partnerName"`
	TrainingBatchName            string        `json:"training_batch_name"`
	ContactPerson                string        `json:"contact_person"`
	ContactNumber                string        `json:"contact_number"`
	Type                         string        `json:"type"`
	CheckIn                      string        `json:"check_in"`
	CheckOut                     string        `json:"check_out"`
	Participants                 string        `json:"participants"`
	TypeName                     string        `json:"type_name"`
	Status                       string        `json:"status"`
	CircleID                     string        `json:"circle_id"`
	PrimaryID                    string        `json:"primary_id"`
	UserID                       string        `json:"user_id"`
	GfName                       string        `json:"gf_name"`
	TrainerName                  string        `json:"trainer_name"`
	SessionCompleted             string        `json:"session_completed"`
	AllParticipants              []Participant `json:"all_participants"`
	TotalParticipants            string        `json:"total_participants"`
	EnrollParticipantID          string        `json:"enroll_participant_id"`
	EnrollParticipantName        string        `json:"enroll_participant_name"`
	EnrollParticipantContactNum  string        `json:"enroll_participant_contactnum"`
	EnrollParticipantVillageName string        `json:"enroll_participant_villagename"`

	Photos  []Photo `json:"photos"`
	Code    int     `json:"code"`
	Success bool    `json:"success"`
	Message string  `json:"message"`
}

type Participant struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	GelathiRecomm   string `json:"gelathiRecomm"`
	GreenMotivators int    `json:"GreenMotivators"`
	GelathiStatus   string `json:"gelathi_status"`
	IsAttended      string `json:"is_attended"`
	ContactNum      string `json:"contactnum,omitempty"`
	VillageName     string `json:"villagename,omitempty"`
}

type Photo struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

func GetGreenMotivator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method not allowed", "success": false})
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request body", "success": false})
		return
	}

	var request map[string]interface{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error decoding request body", "success": false})
		return
	}

	// response := Response{}
	sessionID, ok := request["gf_session_id"].(string)

	if !ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "invalid session id", "success": false})
		return
	}

	condition := ""
	userID, ok := request["user_id"].(int)
	if ok {
		condition = fmt.Sprintf("AND gf_sess.user_id = %d", userID)
	}

	query := fmt.Sprintf(`SELECT gf_sess.id, gf_sess.name AS gf_session_name, gf_sess.tb_id, DATE_FORMAT(gf_sess.date, '%%d-%%m-%%Y %%h:%%i %%p') AS plan_date, prj.id AS project_id, UPPER(prj.projectName) AS projectName, IFNULL(part.partnerName, '') AS partnerName, IFNULL(tr_bat.name, '') AS training_batch_name, IFNULL(tr_bat.contact_person, '') AS contact_person, IFNULL(tr_bat.contact_number, '') AS contact_number, gf_sess.session_type AS type, IF(gf_sess.check_in != '', 1, 0) AS check_in, IF(gf_sess.check_out != '', 1, 0) AS check_out, gf_sess.participants, CASE WHEN gf_sess.session_type = 1 THEN 'Circle Metting' WHEN gf_sess.session_type = '2' THEN 'Village Visit' ELSE 'Beehive Visit' END AS type_name, gf_sess.status, gf_sess.circle_id, gf_sess.primary_id, gf_sess.user_id, CONCAT(emp.first_name,' ',IFNULL(emp.last_name, '')) AS gf_name FROM tbl_poa gf_sess LEFT JOIN project prj ON gf_sess.project_id = prj.id LEFT JOIN partner part ON prj.partnerID = part.partnerID LEFT JOIN tbl_poa tr_bat ON gf_sess.tb_id = tr_bat.tb_id JOIN employee emp ON gf_sess.user_id = emp.id WHERE gf_sess.id = '%s' AND gf_sess.type = '2' %s`, sessionID, condition)

	rows := db.QueryRow(query)

	var session Session

	rows.Scan(
		&session.ID,
		&session.GfSessionName,
		&session.TbID,
		&session.PlanDate,
		&session.ProjectID,
		&session.ProjectName,
		&session.PartnerName,
		&session.TrainingBatchName,
		&session.ContactPerson,
		&session.ContactNumber,
		&session.Type,
		&session.CheckIn,
		&session.CheckOut,
		&session.Participants,
		&session.TypeName,
		&session.Status,
		&session.CircleID,
		&session.PrimaryID,
		&session.UserID,
		&session.GfName,
	)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan row", "success": false})
		return
	}
	// sessions = append(sessions, session)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "session not found", "success": false})
		return
	}

	// Get trainer name who created the batch
	trainerQuery := `
		SELECT
			CONCAT(emp.first_name, ' ', IFNULL(emp.last_name, '')) AS trainer_name
		FROM
			tbl_poa
		JOIN employee AS emp ON tbl_poa.user_id = emp.id
		WHERE
			tbl_poa.primary_id = ? AND tbl_poa.added = '0'
	`

	trainerStmt, err := db.Prepare(trainerQuery)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error preparing  trainer query", "success": false})
		return
	}
	defer trainerStmt.Close()

	trainerRow := trainerStmt.QueryRow(session.TbID)

	err = trainerRow.Scan(&session.TrainerName)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error retrieving trainer name", "success": false})
		return
	}

	session.SessionCompleted = "0"
	if session.CheckOut == "1" {
		session.SessionCompleted = "1"
	}

	var queryParticipant string
	if session.Type == "1" {
		queryParticipant = fmt.Sprintf(`
			SELECT
				tr_part.id AS participant_id,
				IFNULL(tr_part.firstName, '') AS participant_name,
				tr_part.gelathiRecomm,
				coalesce(tr_part.GreenMotivators,0),
				IF(tr_part.gelathiRecomm = '1', 'Suggested', '') AS gelathi_status,
				IFNULL(
					(SELECT COUNT(1)
					FROM participant_attendance
					WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) AS is_attended
			FROM
			
				gelathi_circle AS gel_circle
			LEFT JOIN training_participants AS tr_part ON gel_circle.gelathi_id = tr_part.id
			WHERE
				tr_part.new_green IN (0, 1) AND
				tr_part.new_vyapar = 0 AND
				gel_circle.circle_id = %s`, sessionID, session.TbID)
	} else {
		queryParticipant = fmt.Sprintf(`
			SELECT
				tr_part.id AS participant_id,
				IFNULL(tr_part.firstName, '') AS participant_name,
				tr_part.gelathiRecomm,
				coalesce(tr_part.GreenMotivators,0),
				IF(tr_part.gelathiRecomm = '1', 'Suggested', '') AS gelathi_status,
				IFNULL(
					(SELECT COUNT(1)
					FROM participant_attendance
					WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) AS is_attended
			FROM
				training_participants AS tr_part
			WHERE
				tr_part.new_green IN (0, 1) AND
				tr_part.new_vyapar = 0 AND
				tr_part.tb_id = %s`, sessionID, session.TbID)
	}

	participantsData, err := db.Query(queryParticipant)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error retrieving participants", "success": false})
		return
	}
	defer participantsData.Close()

	participants := []Participant{}
	for participantsData.Next() {
		participant := Participant{}
		err := participantsData.Scan(
			&participant.ParticipantID,
			&participant.ParticipantName,
			&participant.GelathiRecomm,
			&participant.GreenMotivators,
			&participant.GelathiStatus,
			&participant.IsAttended,
		)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error scanning participants", "success": false})
			continue
		}
		// participant.GreenMotivators = map[string]string{
		// 	"0": participant.GreenMotivators,
		// }[participant.GreenMotivators]
		participants = append(participants, participant)
	}
	session.AllParticipants = participants
	session.TotalParticipants = strconv.Itoa(len(participants))

	enrollFields := `
		tr_part.id AS participant_id,
		IFNULL(tr_part.firstName, '') AS participant_name,
		'' AS contactnum,
		IFNULL(SUBSTRING_INDEX(tbl_poa.name, '_', 1), '') AS villagename
	`
	queryEnroll := `
		SELECT
			` + enrollFields + `
		FROM
			training_participants AS tr_part
		LEFT JOIN tbl_poa ON tr_part.tb_id = tbl_poa.id
		WHERE
			tr_part.tb_id = ? AND
			tr_part.GreenMotivators = 1
	`
	enrollData := db.QueryRow(queryEnroll, session.TbID)

	enrollParticipant := Participant{}
	err = enrollData.Scan(
		&enrollParticipant.ParticipantID,
		&enrollParticipant.ParticipantName,
		&enrollParticipant.ContactNum,
		&enrollParticipant.VillageName,
	)
	if err != nil && err != sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error retrieving enrolled participant", "success": false})
		return
	}

	session.EnrollParticipantID = enrollParticipant.ParticipantID
	session.EnrollParticipantName = enrollParticipant.ParticipantName
	session.EnrollParticipantContactNum = enrollParticipant.ContactNum
	session.EnrollParticipantVillageName = enrollParticipant.VillageName

	querySessionPhotos := `
		SELECT
			IF(IFNULL(photo1, '') = '', '', CONCAT(?, '', photo1)) AS photo1,
			IF(IFNULL(photo2, '') = '', '', CONCAT(?, '', photo2)) AS photo2
		FROM
			tbl_poa AS gf_sess_photos
		WHERE
			gf_sess_photos.id = ? LIMIT 1
	`
	sessionPhotos := Photo{}
	err = db.QueryRow(querySessionPhotos, "https://bdms.buzzwomen.org/appTest/", "https://bdms.buzzwomen.org/appTest/", session.ID).Scan(
		&sessionPhotos.Photo1,
		&sessionPhotos.Photo2,
	)
	if err != nil && err != sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error retrieving session photo", "success": false})
		return
	}
	session.Photos = append(session.Photos, sessionPhotos)

	session.Code = 200
	session.Message = "Successfully"
	session.Success = true
	json.NewEncoder(w).Encode(session)
}
