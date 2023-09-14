package Team_3

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Partner struct {
	PartnerID   string `json:"partnerID"`
	PartnerName string `json:"partnerName"`
}

type Partner1 struct {
	List    []Partner `json:"list"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Success bool      `json:"success"`
}

func GetPartnerList(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	query := "SELECT partnerID, partnerName FROM partner WHERE status = 1 ORDER BY partnerID DESC, partnerName"

	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var partnerList []Partner

	for rows.Next() {
		var partner Partner
		err := rows.Scan(&partner.PartnerID, &partner.PartnerName)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		partnerList = append(partnerList, partner)
	}

	response := Partner1{
		List:    partnerList,
		Code:    200,
		Success: true,
		Message: "successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
