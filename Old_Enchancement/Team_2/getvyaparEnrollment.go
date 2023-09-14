package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Rest struct {
	User_id       string `json:"user_id"`
	Gf_session_id string `json:"gf_session_id"`
}

type Vyapar struct {
	Id                             string             `json:"id"`
	Gf_session_name                string             `json:"gf_session_name"`
	Tb_id                          string             `json:"tb_id"`
	Plan_date                      string             `json:"plan_date"`
	Project_id                     string             `json:"project_id"`
	ProjectName                    string             `json:"projectName"`
	PartnerName                    string             `json:"partnerName"`
	Training_batch_name            string             `json:"training_batch_name"`
	Contact_person                 string             `json:"contact_person"`
	Contact_number                 string             `json:"contact_number"`
	Type                           string             `json:"type"`
	Check_in                       string             `json:"check_in"`
	Check_out                      string             `json:"check_out"`
	Participants                   string             `json:"participants"`
	Type_name                      string             `json:"type_name"`
	Status                         string             `json:"status"`
	Circle_id                      string             `json:"circle_id"`
	Primary_id                     string             `json:"primary_id"`
	User_id                        string             `json:"user_id"`
	Gf_name                        string             `json:"gf_name"`
	Trainer_name                   string             `json:"trainer_name"`
	Session_completed              string             `json:"session_completed"`
	All_participants               []All_participants `json:"all_participants"`
	Total_participants             int                `json:"total_participants"`
	Enroll_participant_id          string             `json:"enroll_participant_id"`
	Enroll_participant_name        string             `json:"enroll_participant_name"`
	Enroll_participant_contactnum  string             `json:"enroll_participant_contactnum"`
	Enroll_participant_villagename string             `json:"enroll_participant_villagename"`
	Photos                         []Photos           `json:"photos"`
	Code                           int                `json:"code"`
	Success                        bool               `json:"success"`
	Message                        string             `json:"message"`
}
type Photos struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

type All_participants struct {
	Participant_id   string `json:"participant_id"`
	Participant_name string `json:"participant_name"`
	GelathiRecomm    string `json:"gelathiRecomm"`
	VyaparEnrollment int    `json:"VyaparEnrollment"`
	Gelathi_status   string `json:"gelathi_status"`
	Is_attended      string `json:"is_attended"`
}

func GetvyaparEnrollment(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")


	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	var p Rest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(500) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Decoding error", "Status Code": "500"})
		return
	}

	rawURL := "https://bdms.buzzwomen.org/appTest/images"

	//================================================================

	var condition string
	var query string
	var query_participant string
	var queryEnroll string

	condition = ""
	if p.User_id != "" {
		condition = fmt.Sprintf("and gf_sess.user_id = %s", p.User_id)
	}
	id := p.Gf_session_id
	if id != "" {

		fields := "gf_sess.id, gf_sess.name as gf_session_name, gf_sess.tb_id, DATE_FORMAT(gf_sess.date, '%d-%m-%Y %h:%i %p') as plan_date, prj.id as project_id, UPPER(prj.projectName) as projectName, IFNULL(part.partnerName, '') as partnerName, IFNULL(tr_bat.name, '') as training_batch_name, IFNULL(tr_bat.contact_person, '') as contact_person, IFNULL(tr_bat.contact_number, '') as contact_number, gf_sess.session_type as type, IF(gf_sess.check_in != '', 1, 0) as check_in, IF(gf_sess.check_out != '', 1, 0) as check_out, gf_sess.participants, CASE WHEN gf_sess.session_type = 1 THEN 'Circle Metting' WHEN gf_sess.session_type = '2' THEN 'Village Visit' ELSE 'Beehive Visit' END AS type_name, gf_sess.status, gf_sess.circle_id, gf_sess.primary_id, gf_sess.user_id, CONCAT(emp.first_name,' ',IFNULL(emp.last_name, '')) AS gf_name"

		query = fmt.Sprintf("SELECT %s FROM tbl_poa gf_sess left join project prj on gf_sess.project_id = prj.id left join partner part on prj.partnerID = part.partnerID left join tbl_poa tr_bat on gf_sess.tb_id = tr_bat.tb_id join employee emp on gf_sess.user_id = emp.id where gf_sess.id = %s and gf_sess.type = '2' %s", fields, id, condition)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})

			return

		}
		defer rows.Close()

		var v Vyapar
		for rows.Next() {

			err := rows.Scan(
				&v.Id,
				&v.Gf_session_name,
				&v.Tb_id,
				&v.Plan_date,
				&v.Project_id,
				&v.ProjectName,
				&v.PartnerName,
				&v.Training_batch_name,
				&v.Contact_person,
				&v.Contact_number,
				&v.Type,
				&v.Check_in,
				&v.Check_out,
				&v.Participants,
				&v.Type_name,
				&v.Status,
				&v.Circle_id,
				&v.Primary_id,
				&v.User_id,
				&v.Gf_name,
			)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}

		}

		tbatch_id := v.Tb_id
		tbl_query := fmt.Sprintf("SELECT CONCAT(emp.first_name,' ',IFNULL(emp.last_name, '')) AS trainer_name FROM tbl_poa join employee emp on tbl_poa.user_id = emp.id where tbl_poa.primary_id = %s and tbl_poa.added = '0'", tbatch_id)

		rows1, err := DB.Query(tbl_query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
			return
		}
		defer rows1.Close()

		for rows1.Next() {

			err := rows1.Scan(&v.Trainer_name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return

			}
		}
		session_id := v.Id
		tb_id := v.Tb_id
		Type := v.Type
		circle_id := v.Circle_id

		if Type == "1" {
			query_participant = fmt.Sprintf("SELECT tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, COALESCE(tr_part.VyaparEnrollment,0), IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended from gelathi_circle gel_circle left join training_participants tr_part ON gel_circle.gelathi_id = tr_part.id where tr_part.new_green =0 and tr_part.new_vyapar in (0,1) and gel_circle.circle_id = %s", id, circle_id)

		} else {
			query_participant = fmt.Sprintf("SELECT tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, tr_part.gelathiRecomm, COALESCE(tr_part.VyaparEnrollment,0), IF(tr_part.gelathiRecomm = '1', 'Suggested', '') as gelathi_status, IFNULL((SELECT count(1) From participant_attendance WHERE tbl_poa_id = %s AND participant_id = tr_part.id LIMIT 1), 0) as is_attended from training_participants tr_part where tr_part.new_green =0 and tr_part.new_vyapar in (0,1) and tr_part.tb_id =%s", id, tb_id)
		}

		rows, err1 := DB.Query(query_participant)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
			return
		}
		defer rows.Close()
		var parti []All_participants
		for rows.Next() {
			var rw All_participants
			err := rows.Scan(&rw.Participant_id, &rw.Participant_name, &rw.GelathiRecomm, &rw.VyaparEnrollment, &rw.Gelathi_status, &rw.Is_attended)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}
			parti = append(parti, rw)
		}

		v.Total_participants = len(parti)

		if len(parti) > 0 {
			v.All_participants = parti
		} else {
			v.All_participants = []All_participants{}
		}


		//------------------ Enrolled participant --------------
		enrollfields := "tr_part.id as participant_id, IFNULL(tr_part.firstName, '') as participant_name, '' as contactnum, IFNULL(SUBSTRING_INDEX(tbl_poa.name, '_', 1), '') as villagename"

		queryEnroll = fmt.Sprintf("SELECT %s from training_participants tr_part LEFT JOIN tbl_poa ON tr_part.tb_id = tbl_poa.id where tr_part.tb_id = %s AND tr_part.VyaparEnrollment = 1", enrollfields, tb_id)

		rows, err1 = DB.Query(queryEnroll)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
			return
		}

		defer rows.Close()

		for rows.Next() {

			err := rows.Scan(&v.Enroll_participant_id, &v.Enroll_participant_name, &v.Enroll_participant_contactnum, &v.Enroll_participant_villagename)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}
		}
		//======================================== photos =====================================
		rrr := fmt.Sprintf("SELECT IF(IFNULL(photo1, '') = '', '', CONCAT('"+rawURL+"', '', photo1)) AS photo1, IF(IFNULL(photo2, '') = '', '', CONCAT('"+rawURL+"', '', photo2)) AS photo2 from tbl_poa gf_sess_photos where gf_sess_photos.id = %s LIMIT 1", session_id)

		rows1, err2 := DB.Query(rrr)
		if err2 != nil {

			w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
			return

		}
		defer rows1.Close()
		var pics []Photos
		for rows1.Next() {

			var phto Photos

			err := rows1.Scan(&phto.Photo1, &phto.Photo2)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // Return 405 Method Not Allowed.
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Scan issue ", "success": false})
				return
			}
			pics = append(pics, phto)

		}

		v.Photos = pics
		v.Code = http.StatusOK
		v.Message = "Successfully"
		v.Success = true
		jsonData1, err1 := json.Marshal(v)
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false})
			return
		}
		w.Write(jsonData1)

	}
}
