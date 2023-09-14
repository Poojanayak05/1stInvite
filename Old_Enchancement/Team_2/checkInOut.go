package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CheckInOutReq struct {
	LocationName string `json:"location_name"`
	UserID       string `json:"user_id"`
	Lon          string `json:"lon"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Lat          string `json:"lat"`
}

func CheckInOut(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to Read request Body", "success": false, "error": err.Error()})

		return
	}

	defer r.Body.Close()

	var request CheckInOutReq
	err = json.Unmarshal(data, &request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to Unmarshal request", "success": false, "error": err.Error()})
		return
	}

	result := make(map[string]interface{})

	if request.ID != "" && request.UserID != "" && request.Type != "" {

		id := request.ID
		userID := request.UserID
		reqType := request.Type
		lat := request.Lat
		lon := request.Lon
		location := ""
		if request.LocationName == "" {
			location, err = getAddress(request.Lat, request.Lon)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to get address", "success": false, "error": err.Error()})
				return
			}

		} else {
			location = request.LocationName
		}

		// Assuming checkInOut is fetched from the request or any other source
		// checkInOut := request["check_in_out"]

		checkInOutTime := time.Now().Add(5*time.Hour + 30*time.Minute)
		if checkInOutTime != (time.Time{}) {
			checkInOutTimeString := checkInOutTime.Format("2006-01-02 15:04:05")
			parsedTime, err := time.Parse("2006-01-02 15:04:05", checkInOutTimeString)
			if err == nil {
				checkInOutTime = parsedTime
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to parsed Time", "success": false, "error": err.Error()})
				return
			}
		}

		year := checkInOutTime.Year()
		if year == 1970 {
			checkInOutTime = time.Now().Add(5*time.Hour + 30*time.Minute)
		}

		msg := ""
		if reqType == "1" {
			msg = "Check In"
		} else {
			msg = "Check Out"
		}


		var updateQuery string

		// Construct update query

		switch reqType {
		case "1":
			
			updateQuery = fmt.Sprintf("UPDATE tbl_poa SET check_in = '%s', check_in_lat = '%s', check_in_lon = '%s', check_in_location = '%s' WHERE id = %s AND user_id = %s",
				checkInOutTime.Format("2006-01-02 15:04:05"), lat, lon, location, id, userID)
		case "2":
			
			updateQuery = fmt.Sprintf("UPDATE tbl_poa SET check_out = '%s', check_out_lat = '%s', check_out_lon = '%s', check_out_location = '%s' WHERE id = %s AND user_id = %s",
				checkInOutTime.Format("2006-01-02 15:04:05"), lat, lon, location, id, userID)
		default:
			// Handle any other case
		}

		// Execute the update query
		row, err := DB.Exec(updateQuery)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to updateQuery", "success": false, "error": err.Error()})

		}

		if row != nil {
			result["code"] = 200
			result["success"] = true
			result["message"] = msg + " Successfully"

		} else {

			result["code"] = 200
			result["success"] = false
			result["message"] = msg + " Unsuccessfully"

		}

		// Convert response to JSON

		responseJSON, err := json.Marshal(result)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false, "error": err.Error()})
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send response
		_, err = w.Write(responseJSON)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to write response", "success": false, "error": err.Error()})
		}
	} else {

		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Send Required Parameters", "success": false})

	}
}

func getAddress(latitude, longitude string) (string, error) {
	url := fmt.Sprintf("https://reverse.geocoder.ls.hereapi.com/6.2/reversegeocode.json?apiKey=19wEmeLqfrKkUhF3oF8nNqrdrgAh-laYG8B2RluA-Lk&app_id=TXBWynjhxFEWMTIvW0CV&prox=%s,%s&mode=retrieveAddresses&additionaldata=PreserveUnitDesignators,true", latitude, longitude)

	response, err := http.Get(url)

	if err != nil {

		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	

	location := data["Response"].(map[string]interface{})["View"].([]interface{})[0].(map[string]interface{})["Result"].([]interface{})[0].(map[string]interface{})["Location"].(map[string]interface{})["Address"].(map[string]interface{})["Label"].(string)
	

	return location, err
}
