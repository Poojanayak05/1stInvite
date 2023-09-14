package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// struct to store bus details coming in the body.
type EditBusRequest struct {
	BusID              string `json:"bus_id"`
	RegisterNumber     string `json:"register_number"`
	RegisterDate       string `json:"register_date"`
	EngineNumber       string `json:"engine_number"`
	ChassisNumber      string `json:"chassis_number"`
	InsuranceNumber    string `json:"insurance_number"`
	InsuranceCompany   string `json:"insurance_company"`
	InsuranceStartDate string `json:"insurance_start_date"`
	InsuranceEndDate   string `json:"insurance_end_date"`
	LastServiceDate    string `json:"last_service_date"`
	NextServiceDueDate string `json:"next_service_due_date"`
	LastUpdatedBy      string `json:"lastUpdatedBy"`
	FitnessCertificate string `json:"fitness_certificate"`
	Permit             string `json:"permit"`
	EmissionDate       string `json:"emission_date"`
}

// struct to show the response.
type EditBusResponse struct {
	Code         int    `json:"code"`
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ErrorMessage error  `json:"error,omitempty"`
}

func EditBus(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// declaring request variable to refer edit bus request struct.
	var request EditBusRequest
	// declaring response variable to refer edit bus respnse struct.
	var response EditBusResponse

	// dumping the body into edit bus request.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Message = "Error While Decoding Body"
		response.Code = 400
		response.ErrorMessage = err
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.BusID == "" {
		response.Code = 406
		response.Message = "Please Send the Required Bus Id."
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	// inner function to convert string into date format.
	// dateParseFunc := func(dateStr string) *time.Time {
	// 	if dateStr == "" {
	// 		return nil
	// 	}
	// 	date, err := time.Parse("2006-01-02", dateStr)
	// 	if err != nil {
	// 		fmt.Printf("Failed to parse date: %s", err.Error())
	// 		return nil
	// 	}
	// 	return &date
	// }

	dateParseFunc2 := func(dateStr string) *time.Time {
		if dateStr == "" {
			return nil
		}
		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			fmt.Printf("Failed to parse date: %s", err.Error())
			return nil
		}
		return &date
	}

	// inner function to convert string into date time format.
	dateTimeParseFunc := func(dateStr string) *time.Time {
		if dateStr == "" {
			return nil
		}
		date, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
		if err != nil {
			fmt.Printf("Failed to parse date: %s", err.Error())
			return nil
		}
		return &date
	}
	// assining the struct variable's value to simple variables.
	id := request.BusID
	registerNumber := request.RegisterNumber
	registerDate := dateTimeParseFunc(request.RegisterDate)
	engineNumber := request.EngineNumber
	chassisNumber := request.ChassisNumber
	insuranceNumber := request.InsuranceNumber
	insuranceCompany := request.InsuranceCompany
	insuranceStartDate := dateTimeParseFunc(request.InsuranceStartDate)
	insuranceEndDate := dateTimeParseFunc(request.InsuranceEndDate)
	lastServiceDate := dateParseFunc2(request.LastServiceDate)
	nextServiceDueDate := dateTimeParseFunc(request.NextServiceDueDate)
	lastUpdatedBy := request.LastUpdatedBy
	fitnessCertificate := dateParseFunc2(request.FitnessCertificate)
	permit := dateParseFunc2(request.Permit)
	emissionDate := dateParseFunc2(request.EmissionDate)

	// Update the bus in the database
	updateQuery := `UPDATE bus SET register_number = ?, register_date = ?, engine_number = ?, chassis_number = ?, insurance_number = ?, insurance_company = ?, insurance_start_date = ?, insurance_end_date = ?, last_service_date = ?, next_service_due_date = ?, lastUpdatedBy = ?, fitness_certificate = ?, permit = ?, emission_date = ? WHERE id = ?`
	_, err = DB.Exec(updateQuery, registerNumber, registerDate, engineNumber, chassisNumber, insuranceNumber, insuranceCompany, insuranceStartDate, insuranceEndDate, lastServiceDate, nextServiceDueDate, lastUpdatedBy, fitnessCertificate, permit, emissionDate, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := EditBusResponse{
			Code:         404,
			Message:      "Registered Number Already Exists",
			ErrorMessage: err,
			Success:      false,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// setting the response.
	response = EditBusResponse{
		Code:    200,
		Message: "Bus Updated Successfully",
		Success: true,
	}
	// sending the response.
	json.NewEncoder(w).Encode(response)

}
