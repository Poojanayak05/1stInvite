package Team_3

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Request100 struct {
	RegisterNumber     string `json:"register_number"`
	RegisterDate       string `json:"register_date"`
	EngineNumber       string `json:"engine_number"`
	ChassisNumber      string `json:"chassis_number"`
	FitnessCertificate string `json:"fitness_certificate"`
	Permit             string `json:"permit"`
	InsuranceNumber    string `json:"insurance_number"`
	InsuranceCompany   string `json:"insurance_company"`
	InsuranceStartDate string `json:"insurance_start_date"`
	InsuranceEndDate   string `json:"insurance_end_date"`
	LastServiceDate    string `json:"last_service_date"`
	NextServiceDueDate string `json:"next_service_due_date"`
	CreatedBy          string `json:"createdBy"`
	LastUpdatedBy      string `json:"lastUpdatedBy"`
}

var request Request100

func CreateBus(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return

	}

	// var request Request1
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var register_date11 time.Time
	register_date1 := "2006-01-02"
	if request.RegisterDate == "" {
		register_date11, err = time.Parse(register_date1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid RegisterDate format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		// Layout to parse the date string
		register_date11, err = time.Parse(register_date1, request.RegisterDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid RegisterDate format",
				"success": false,
				"error":   err,
			})
			return
		}
	}

	fitness_certificate1 := "2006-01-02" // Layout to parse the date string
	var fitness_certificate11 time.Time
	if request.FitnessCertificate == "" {
		fitness_certificate11, err = time.Parse(fitness_certificate1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid fitness_certificate format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		fitness_certificate11, err = time.Parse(fitness_certificate1, request.FitnessCertificate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid fitness_certificate format",
				"success": false,
				"error":   err,
			})
			return
		}
	}

	insurance_end_date1 := "2006-01-02T15:04:05.999Z" // Layout to parse the date string
	var insurance_end_date11 time.Time
	if request.InsuranceEndDate == "" {
		insurance_end_date11, err = time.Parse(insurance_end_date1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid insurance_end_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		insurance_end_date11, err = time.Parse(insurance_end_date1, request.InsuranceEndDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid insurance_end_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	}

	last_service_date1 := "2006-01-02T15:04:05.999Z" // Layout to parse the date string
	var last_service_date11 time.Time
	if request.LastServiceDate == "" {
		last_service_date11, err = time.Parse(last_service_date1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid last_service_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		last_service_date11, err = time.Parse(last_service_date1, request.LastServiceDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid last_service_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	}
	next_service_due_date1 := "2006-01-02T15:04:05.999Z" // Layout to parse the date string
	var next_service_due_date11 time.Time
	if request.NextServiceDueDate == "" {
		next_service_due_date11, err = time.Parse(next_service_due_date1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid next_service_due_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		next_service_due_date11, err = time.Parse(next_service_due_date1, request.NextServiceDueDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid next_service_due_date format",
				"success": false,
				"error":   err,
			})
			return
		}
	}

	permit1 := "2006-01-02" // Layout to parse the date string
	var permit11 time.Time
	if request.Permit == "" {
		permit11, err = time.Parse(permit1, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid permit format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		permit11, err = time.Parse(permit1, request.Permit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid permit format",
				"success": false,
				"error":   err,
			})
			return
		}
	}
	layout := "2006-01-02T15:04:05.999Z" // Layout to parse the date string
	var startDate time.Time
	if request.InsuranceStartDate == "" {
		startDate, err = time.Parse(layout, "0000-00-00")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid InsuranceStartDate format",
				"success": false,
				"error":   err,
			})
			return
		}
	} else {
		startDate, err = time.Parse(layout, request.InsuranceStartDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid InsuranceStartDate format",
				"success": false,
				"error":   err,
			})
			return
		}
	}

	query := "INSERT INTO bus(register_number,register_date,engine_number,chassis_number,fitness_certificate,permit,insurance_number,insurance_company,insurance_start_date,insurance_end_date,last_service_date,next_service_due_date,createdBy,lastUpdatedBy) Values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	_, err4 := DB.Exec(query,
		request.RegisterNumber,
		register_date11,
		request.EngineNumber,
		request.ChassisNumber,
		fitness_certificate11,
		permit11,
		request.InsuranceNumber,
		request.InsuranceCompany,
		startDate,
		insurance_end_date11,
		last_service_date11,
		next_service_due_date11,
		request.CreatedBy,
		request.LastUpdatedBy,
	)
	if err4 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid input syntax", "success": false, "error": err4})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "bus Added Successfully", "success": true})

}
