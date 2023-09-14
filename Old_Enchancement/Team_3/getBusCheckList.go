package Team_3

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetBusCheckList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type Req struct {
		Bus_id string `json:"bus_id"`
		Date   string `json:"date"`
	}

	type RespData struct {
		Id                  int    `json:"id"`
		Sufficient_diesel   int    `json:"sufficient_diesel"`
		Tyres_air_condition int    `json:"tyres_air_condition"`
		Battery             int    `json:"battery"`
		Switch_box          int    `json:"switch_box"`
		Led_tv_remote       int    `json:"led_tv_remote"`
		Fan                 int    `json:"fan"`
		Tube_lights         int    `json:"tube_lights"`
		Mic                 int    `json:"mic"`
		Jock_liver          int    `json:"jock_liver"`
		Wheel_spanner       int    `json:"wheel_spanner"`
		Canopy_unfold_rod   int    `json:"canopy_unfold_rod"`
		Speaker             int    `json:"speaker"`
		Chairs              int    `json:"chairs"`
		Tent_sidewall       int    `json:"tent_sidewall"`
		Tray                int    `json:"tray"`
		Vehicle_clean       int    `json:"vehicle_clean"`
		Final_save          int    `json:"final_save"`
		Code                int    `json:"code"`
		Success             bool   `json:"success"`
		Message             string `json:"message"`
	}

	var request Req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})
		return
	}

	if request.Date == "" {
		request.Date = time.Now().Format("2006-01-02")
	}
	count := 0
	err = DB.QueryRow("select count(bus_id) from bus_checklist where bus_id=?", request.Bus_id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "busID doesn't exist",
				"success": true,
			})
		}
		return

	}
	if count != 0 && request.Bus_id != "" {

		fields := "bus_list.id, bus_list.sufficient_diesel, bus_list.tyres_air_condition, bus_list.battery, bus_list.switch_box, bus_list.led_tv_remote, bus_list.fan, bus_list.tube_lights, bus_list.mic, bus_list.jock_liver, bus_list.wheel_spanner, bus_list.canopy_unfold_rod, bus_list.speaker, bus_list.chairs, bus_list.tent_sidewall , bus_list.tray, bus_list.vehicle_clean,bus_list.final_save"

		query := fmt.Sprintf("SELECT %s FROM bus_checklist bus_list WHERE bus_list.bus_id = %s AND bus_list.date_checked = '%s'", fields, request.Bus_id, request.Date)

		rows, err := DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "Failed to query database1",
				"success": true,
			})
			return
		}
		var response RespData
		for rows.Next() {
			rows.Scan(&response.Id, &response.Sufficient_diesel, &response.Tyres_air_condition, &response.Battery, &response.Switch_box, &response.Led_tv_remote, &response.Fan, &response.Tube_lights, &response.Mic, &response.Jock_liver, &response.Wheel_spanner, &response.Canopy_unfold_rod, &response.Speaker, &response.Chairs, &response.Tent_sidewall, &response.Tray, &response.Vehicle_clean, &response.Final_save)

		}
		response.Code = http.StatusOK
		response.Message = "successfully"
		response.Success = true

		json.NewEncoder(w).Encode(response)

	} else {
		jsonData, err := json.Marshal(map[string]interface{}{"message": "invalid bus id", "code": http.StatusBadRequest, "success": false})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": err, "success": false})
			return
		}
		w.Write(jsonData)
	}
}
