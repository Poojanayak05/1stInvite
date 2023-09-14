package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CreateGFSessions(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type RequestBody struct {
		ProjectID         string `json:"project_id"`
		UserID            string `json:"user_id"`
		LocationID        string `json:"locationId"`
		TBName            string `json:"tb_name"`
		NumOfParticipants string `json:"numOfParticipants"`
		TBID              string `json:"tb_id"`
		GFSessionType     int    `json:"gf_session_type"`
		PlanDate          string `json:"plan_date"`
		GFSessionName     string `json:"gf_session_name"`
		CircleID          string `json:"circle_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var data RequestBody

	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "message": err.Error()})
		return
	}

	typeVar := 2
	tbId := data.TBID
	_, _ = typeVar, tbId

	getGFSessionName := func(name string, gfSessionType int) string {
		nameParts := strings.Split(name, "_")

		var sessionType string
		switch gfSessionType {
		case 1:
			sessionType = "CM"
		case 3:
			sessionType = "BV"
		default:
			sessionType = "VV"
		}

		return nameParts[0] + "_" + sessionType
	}
	gfSessionType := data.GFSessionType

	name := getGFSessionName(data.TBName, gfSessionType)
	userID := data.UserID
	projectID := data.ProjectID
	participants := data.NumOfParticipants
	date := data.PlanDate
	sessionType := data.GFSessionType
	circleID := data.CircleID
	locationId := data.LocationID

	// tx, err := DB.Begin()
	var dateStr string
	if date != "" {
		//to get time.time date
		parseDateTime := func(dateStr string) (time.Time, error) {
			layout := "2006-1-2 3:04 PM"
			t, err := time.Parse(layout, dateStr)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
		// to get string date

		dateFormat, err := parseDateTime(date)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}

		dateStr = dateFormat.Format("2006-01-02 15:04:05")

		if !dateFormat.Equal(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)) {
			// query:="SELECT count(id) as count FROM tbl_poa where date = "+ dateFormat +" and user_id = $user_id and status != '2'"
			query := fmt.Sprintf("SELECT count(id) as count FROM tbl_poa where date ='%s' and user_id = %s and status != '2'", dateStr, userID)

			var count int
			rows, err := DB.Query(query)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
				return
			}
			for rows.Next() {
				rows.Scan(&count)
			}

			if count > 0 {
				json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Poa already exist for date & time"})
				return

			}

		}

	}

	if data.TBName == "" || data.TBName == "0" {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "village name is empty", "success": false})

	}

	insertQuery := "INSERT INTO tbl_poa (type, tb_id, name, user_id, project_id,location_id, participants, date, session_type, circle_id)VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := DB.Prepare(insertQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
		return
	}

	result, err := stmt.Exec(typeVar, tbId, name, userID, projectID, locationId, participants, dateStr, sessionType, circleID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)

	}

	gfSessionName := name + strconv.FormatInt(insertedID, 10)

	updateQuery := "UPDATE tbl_poa SET name = ?, primary_id = ? WHERE id = ?"
	stmt, err = DB.Prepare(updateQuery)
	if err != nil {

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
		return
	}

	_, err = stmt.Exec(gfSessionName, insertedID, insertedID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "GF Sessions Added Successfully"})

}
