package Team_1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Response struct for the JSON response
type PResponse struct {
	RatePerKM    []RatePerKM    `json:"Rate_per_KM"`
	Data         []Data         `json:"data"`
	ModeOfTravel []ModeOfTravel `json:"Mode_of_Travel"`
	DA           []DA           `json:"DA"`
	Telephone    []Telephone    `json:"telephone"`
	Code         int            `json:"code"`
	Message      string         `json:"message"`
	Success      bool           `json:"success"`
}

// RatePerKM struct for the rate per kilometer
type RatePerKM struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

// Data struct for the POA data
type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

// ModeOfTravel struct for the mode of travel
type ModeOfTravel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DA struct for the DA (Daily Allowance)
type DA struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// Telephone struct for telephone expenses
type Telephone struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func GetPoaTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Failed to decode request body", "Error": err})
		return
	}

	response := PResponse{}
	if empID, ok := request["emp_id"].(string); ok {
		currentYear := time.Now().Year()
		currentMonth := int(time.Now().Month())

		response.RatePerKM = []RatePerKM{
			{
				ID:     1,
				Name:   "bike",
				Amount: "Rs 3.50/km",
			},
			{
				ID:     2,
				Name:   "car",
				Amount: "Rs 10/km",
			},
		}

		response.ModeOfTravel = []ModeOfTravel{
			{
				ID:   1,
				Name: "Bus",
			},
			{
				ID:   2,
				Name: "Car",
			},
			{
				ID:   3,
				Name: "Bike",
			},
			{
				ID:   4,
				Name: "Auto",
			},
			{
				ID:   5,
				Name: "Train/Metro",
			},
			{
				ID:   6,
				Name: "Taxi",
			},
		}

		response.DA = []DA{
			{
				ID:     1,
				Name:   "Bangalore - Rs 200",
				Amount: 200,
			},
			{
				ID:     2,
				Name:   "Other District - Rs 150",
				Amount: 150,
			},
			{
				ID:     3,
				Name:   "Others - Rs 100",
				Amount: 100,
			},
			{
				ID:     4,
				Name:   "No Food Expense - Rs 0",
				Amount: 0,
			},
		}

		response.Telephone = []Telephone{
			{
				ID:     1,
				Name:   "Monthly telephone Expense- Rs 249",
				Amount: 249,
			},
			{
				ID:     2,
				Name:   "No telephone Expense - Rs 0",
				Amount: 0,
			},
		}

		query := fmt.Sprintf("SELECT id, name, date FROM tbl_poa WHERE YEAR(`date`) = %d AND MONTH(`date`) = %d AND tb_id != primary_id AND user_id = %s", currentYear, currentMonth, empID)
		res, err := DB.Query(query)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Failed to execute query", "success": false})

			return
		}
		defer res.Close()

		for res.Next() {
			poa := Data{}
			err = res.Scan(&poa.ID, &poa.Name, &poa.Date)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Failed to scan row", "success": false})

				return
			}

			response.Data = append(response.Data, poa)
		}

		err = res.Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "Message": "Failed to fetch data", "success": false})
			return
		}

		response.Code = 200
		response.Success = true
		response.Message = "TA Data"

		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Data Not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": "Invalid input"})
		return
	}
}
