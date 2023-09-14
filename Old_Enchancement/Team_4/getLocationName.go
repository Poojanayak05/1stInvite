package Team_4

//Done by keerthana
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetLocationName(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method Not found", "success": false})
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "Error": err})
		return
	}

	var request map[string]interface{}
	err = json.Unmarshal(data, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON format", "success": false, "Error": err})
		return
	}

	latitude := request["latitude"].(string)
	longitude := request["longitude"].(string)

	address := ""

	url := fmt.Sprintf("https://reverse.geocoder.ls.hereapi.com/6.2/reversegeocode.json?apiKey=19wEmeLqfrKkUhF3oF8nNqrdrgAh-laYG8B2RluA-Lk&app_id=TXBWynjhxFEWMTIvW0CV&prox=%s,%s&mode=retrieveAddresses&additionaldata=PreserveUnitDesignators,true", latitude, longitude)
	apiResponse, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to read request body", "success": false, "Error": err})
		return
	}
	defer apiResponse.Body.Close()

	body, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to read request body", "success": false, "Error": err})
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to read request body", "success": false, "Error": err})
		return
	}



	location := result["Response"].(map[string]interface{})["View"].([]interface{})[0].(map[string]interface{})["Result"].([]interface{})[0].(map[string]interface{})["Location"].(map[string]interface{})["Address"].(map[string]interface{})["Label"].(string)
	address = location

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(address))

}
