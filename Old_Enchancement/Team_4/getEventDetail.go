package Team_4

//Done by Prathamesh
import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	EventID          string `json:"event_id"`
	Name             string `json:"name"`
	Date1            string `json:"date1"`
	Date2            string `json:"date2"`
	CheckIn          string `json:"check_in"`
	CheckOut         string `json:"check_out"`
	Description      string `json:"description"`
	Photo1           string `json:"photo1"`
	CheckInLat       string `json:"check_in_lat"`
	CheckInLon       string `json:"check_in_lon"`
	CheckInLocation  string `json:"check_in_location"`
	CheckOutLat      string `json:"check_out_lat"`
	CheckOutLon      string `json:"check_out_lon"`
	CheckOutLocation string `json:"check_out_location"`
	EventCompleted   string `json:"event_completed"`
}

type EventResponse struct {
	EventID          string `json:"event_id"`
	Name             string `json:"name"`
	Date1            string `json:"date1"`
	Date2            string `json:"date2"`
	CheckIn          string `json:"check_in"`
	CheckOut         string `json:"check_out"`
	Description      string `json:"description"`
	Photo1           string `json:"photo1"`
	CheckInLat       string `json:"check_in_lat"`
	CheckInLon       string `json:"check_in_lon"`
	CheckInLocation  string `json:"check_in_location"`
	CheckOutLat      string `json:"check_out_lat"`
	CheckOutLon      string `json:"check_out_lon"`
	CheckOutLocation string `json:"check_out_location"`
	EventCompleted   string `json:"event_completed"`
	Code             int    `json:"code"`
	Success          bool   `json:"success"`
	Message          string `json:"message"`
}

type EventProject struct {
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
}

func GetEventDetailist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method Not allowed", "success": false})
		return
	}

	var pr EventProject
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:    http.StatusBadRequest,
			Message: "Error while decoding request body",
			Success: false,
		})
		return
	}
	rawURL := "https://bdms.buzzwomen.org/appTest/"
	row := DB.QueryRow("SELECT id as event_id, COALESCE(name,''), COALESCE(DATE_FORMAT(date, '%d-%m-%Y %h:%i %p'),'') as date1, COALESCE(DATE_FORMAT(date2, '%d-%m-%Y %h:%i %p'),'') as date2,COALESCE( DATE_FORMAT(check_in, '%d-%m-%Y %h:%i %p'),'') as check_in, COALESCE(DATE_FORMAT(check_out, '%d-%m-%Y %h:%i %p'),'') as check_out, IFNULL(description, '') as description, IFNULL(CONCAT('"+rawURL+"', '', photo1), '') as photo1, IFNULL(check_in_lat, '') as check_in_lat, IFNULL(check_in_lon, '') as check_in_lon, IFNULL(check_in_location, '') as check_in_location, IFNULL(check_out_lat, '') as check_out_lat, IFNULL(check_out_lon, '') as check_out_lon, IFNULL(check_out_location, '') as check_out_location, IF(check_out != '', '1', '0') as event_completed FROM tbl_poa WHERE id = ?", pr.EventID)

	var event Event
	err = row.Scan(&event.EventID, &event.Name, &event.Date1, &event.Date2, &event.CheckIn, &event.CheckOut, &event.Description, &event.Photo1, &event.CheckInLat, &event.CheckInLon, &event.CheckInLocation, &event.CheckOutLat, &event.CheckOutLon, &event.CheckOutLocation, &event.EventCompleted)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No data for the given Event Id", "success": false, "error": err})
		return
	}
	response := EventResponse{
		EventID:          event.EventID,
		Name:             event.Name,
		Date1:            event.Date1,
		Date2:            event.Date2,
		CheckIn:          event.CheckIn,
		CheckOut:         event.CheckOut,
		Description:      event.Description,
		Photo1:           event.Photo1,
		CheckInLat:       event.CheckInLat,
		CheckInLon:       event.CheckInLon,
		CheckInLocation:  event.CheckInLocation,
		CheckOutLat:      event.CheckOutLat,
		CheckOutLon:      event.CheckOutLon,
		CheckOutLocation: event.CheckOutLocation,
		EventCompleted:   event.EventCompleted,
		Code:             http.StatusOK,
		Success:          true,
		Message:          "Successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
