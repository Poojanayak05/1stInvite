package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Bus struct {
	BusID              string         `json:"bus_id"`
	RegisterNumber     string         `json:"register_number"`
	RegisterDate       string         `json:"register_date"`
	EngineNumber       string         `json:"engine_number"`
	ChassisNumber      string         `json:"chassis_number"`
	FitnessCertificate string         `json:"fitness_certificate"`
	Permit             string         `json:"permit"`
	InsuranceNumber    string         `json:"insurance_number"`
	InsuranceCompany   string         `json:"insurance_company"`
	InsuranceStartDate string         `json:"insurance_start_date"`
	InsuranceEndDate   string         `json:"insurance_end_date"`
	LastServiceDate    string         `json:"last_service_date"`
	NextServiceDueDate string         `json:"next_service_due_date"`
	EmissionDate       string         `json:"emission_date"`
	BusCheckList       []BusCheckList `json:"bus_checkList"`
	Code               int            `json:"code"`
	Success            bool           `json:"success"`
	Message            string         `json:"message"`
}

type BusCheckList struct {
	DateChecked       string `json:"date_checked"`
	SufficientDiesel  string `json:"sufficient_diesel"`
	TyresAirCondition string `json:"tyres_air_condition"`
	Battery           string `json:"battery"`
	SwitchBox         string `json:"switch_box"`
	LedTVRemote       string `json:"led_tv_remote"`
	Fan               string `json:"fan"`
	TubeLights        string `json:"tube_lights"`
	Mic               string `json:"mic"`
	JockLiver         string `json:"jock_liver"`
	WheelSpanner      string `json:"wheel_spanner"`
	CanopyUnfoldRod   string `json:"canopy_unfold_rod"`
	Speaker           string `json:"speaker"`
	Chairs            string `json:"chairs"`
	TentSidewall      string `json:"tent_sidewall"`
	Tray              string `json:"tray"`
	VehicleClean      string `json:"vehicle_clean"`
	CheckedCount      string `json:"checked_count"`
}

func GetBusData(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	var request map[string]interface{}
	var busChecklist []BusCheckList
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode request body", "success": false, "error": err})
		return
	}
	if busID, ok := request["bus_id"].(string); ok {
		fields := " id as bus_id, UPPER(register_number) as register_number, DATE_FORMAT(register_date, '%d-%m-%Y') as register_date, engine_number, chassis_number, DATE_FORMAT(fitness_certificate, '%d-%m-%Y') as fitness_certificate, DATE_FORMAT(permit, '%d-%m-%Y') as permit, insurance_number, insurance_company, DATE_FORMAT(insurance_start_date, '%d-%m-%Y') as insurance_start_date,DATE_FORMAT(insurance_end_date, '%d-%m-%Y')as insurance_end_date, DATE_FORMAT(last_service_date, '%d-%m-%Y') as last_service_date, DATE_FORMAT(next_service_due_date, '%d-%m-%Y')as next_service_due_date, DATE_FORMAT(emission_date, '%d-%m-%Y') as emission_date"
		query := fmt.Sprintf("SELECT %s FROM bus WHERE id = ?", fields)
		row := DB.QueryRow(query, busID)

		var bus Bus

		err = row.Scan(
			&bus.BusID,
			&bus.RegisterNumber,
			&bus.RegisterDate,
			&bus.EngineNumber,
			&bus.ChassisNumber,
			&bus.FitnessCertificate,
			&bus.Permit,
			&bus.InsuranceNumber,
			&bus.InsuranceCompany,
			&bus.InsuranceStartDate,
			&bus.InsuranceEndDate,
			&bus.LastServiceDate,
			&bus.NextServiceDueDate,
			&bus.EmissionDate,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Response", "success": false, "error": err})
			return
		}

		busID := bus.BusID
		checked := fmt.Sprintf("IFNULL((SELECT check_count From bus_checklist WHERE date_checked = bus_chk.date_checked AND bus_id = %s ), 0) as checked_count", busID)
		var date string
		if fromDate, ok := request["fromDate"].(string); ok {
			if toDate, ok := request["toDate"].(string); ok {
				if fromDate != "" && toDate != "" {
					date = fmt.Sprintf("BETWEEN '%s' AND '%s'", fromDate, toDate)
				} else if fromDate != "" && toDate == "" {
					date = fmt.Sprintf("BETWEEN '%s' AND '%s'", fromDate, fromDate)
				} else {
					today := time.Now().Add(5*time.Hour + 30*time.Minute).Format("2006-01-02")
					date = fmt.Sprintf("BETWEEN '%s' AND '%s'", today, today)
					checked = fmt.Sprintf("IFNULL((SELECT check_count From bus_checklist WHERE date_checked = '%s' AND bus_id = %s), 0) as checked_count", today, busID)

				}
			}
		}
		bus_fields := "DATE_FORMAT(bus_chk.date_checked, '%d-%m-%Y')as date_checked, bus_chk.sufficient_diesel, bus_chk.tyres_air_condition, bus_chk.battery, bus_chk.switch_box, bus_chk.led_tv_remote, bus_chk.fan, bus_chk.tube_lights, bus_chk.mic, bus_chk.jock_liver, bus_chk .wheel_spanner, bus_chk.canopy_unfold_rod, bus_chk.speaker, bus_chk.chairs, bus_chk.tent_sidewall, bus_chk.tray, bus_chk.vehicle_clean, " + fmt.Sprint(checked)

		if date != "" {
			query := "SELECT " + fmt.Sprint(bus_fields) + " FROM bus_checklist AS bus_chk WHERE bus_chk.date_checked " + fmt.Sprint(date) + " AND bus_chk.bus_id = ? ORDER BY bus_chk.id DESC"
			rows, err := DB.Query(query, busID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
				return
			}
			defer rows.Close()

			for rows.Next() {
				var busCheck BusCheckList
				err := rows.Scan(
					&busCheck.DateChecked,
					&busCheck.SufficientDiesel,
					&busCheck.TyresAirCondition,
					&busCheck.Battery,
					&busCheck.SwitchBox,
					&busCheck.LedTVRemote,
					&busCheck.Fan,
					&busCheck.TubeLights,
					&busCheck.Mic,
					&busCheck.JockLiver,
					&busCheck.WheelSpanner,
					&busCheck.CanopyUnfoldRod,
					&busCheck.Speaker,
					&busCheck.Chairs,
					&busCheck.TentSidewall,
					&busCheck.Tray,
					&busCheck.VehicleClean,
					&busCheck.CheckedCount,
				)

				if err != nil {
					log.Println(err)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
					return
				}
				busChecklist = append(busChecklist, busCheck)

			}
		}
		if len(busChecklist) > 0 {
			bus.BusCheckList = busChecklist
		} else {
			bus.BusCheckList = []BusCheckList{}
		}
		bus.Code = 200
		bus.Success = true
		bus.Message = "Successfully"

		jsonData, err := json.Marshal(bus)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}

}
