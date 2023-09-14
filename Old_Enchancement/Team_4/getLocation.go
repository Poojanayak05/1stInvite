package Team_4

//done by Dhiraj Lakhane
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetLocations(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	type Location struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Type     int    `json:"type"`
		ParentID int    `json:"parent_id"`
	}

	type Response struct {
		List    []Location `json:"list"`
		Code    int        `json:"code"`
		Success bool       `json:"success"`
		Message string     `json:"message"`
	}

	type Loc struct {
		CountryID  string `json:"country_id"`
		StateID    string `json:"state_id"`
		DistrictID string `json:"district_id"`
	}

	response := Response{}

	decoder := json.NewDecoder(r.Body)
	var data Loc
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err.Error()})
		return
	}
	var resList []Location

	var res Location

	if data.CountryID != "" && data.StateID == "" && data.DistrictID == "" {
		countryID, err1 := strconv.Atoi(data.CountryID)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err1.Error()})
			return
		}
		query := fmt.Sprintf("SELECT id,name,type,parentId FROM location WHERE parentId = %d AND type = 2 ORDER BY name", countryID)

		rows, err := DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error", "success": false, "error": err.Error()})
			return
		}
		defer rows.Close()
		// var res Location

		for rows.Next() {
			// location := Location{}

			err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.ParentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error occured while scanning in country id func", "success": false, "error": err.Error()})
				return
			}
			resList = append(resList, res)

		}

	}

	if data.StateID != "" && data.DistrictID == "" {
		stateID, err1 := strconv.Atoi(data.StateID)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err1.Error()})
			return
		}
		query := fmt.Sprintf("SELECT id,name,type,parentId FROM location WHERE parentId = %d AND type = 3 ORDER BY name", stateID)
		rows, err := DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error", "success": false, "error": err.Error()})
			return
		}
		defer rows.Close()

		// var res1 Location
		for rows.Next() {
			// location := Location{}

			err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.ParentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error occured while scanning in state id func", "success": false, "error": err})
				return
			}
			resList = append(resList, res)

		}
	}

	if data.DistrictID != "" {
		districtID, err1 := strconv.Atoi(data.DistrictID)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err1.Error()})
			return
		}
		query := fmt.Sprintf("SELECT id,name,type,parentId FROM location WHERE parentId = %d AND type = 4 ORDER BY name", districtID)

		rows, err := DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal server error", "success": false, "error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(&res.ID, &res.Name, &res.Type, &res.ParentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error occured while scanning in district id func", "success": false, "error": err.Error()})
				return
			}

			resList = append(resList, res)

		}
	}
	if len(resList) == 0 {
		resList = []Location{}
	}

	response.List = resList
	response.Code = http.StatusOK
	response.Success = true
	response.Message = "successfully"

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "error while encoding response", "success": false, "error": err.Error()})
		return
	}
}
