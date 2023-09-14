package Team_2

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetProjectListToAssociate(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	//func Login(w http.ResponseWriter, r *http.Request) {
	//SetupCORS(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	type data struct {
		Id          string `json:"id"`
		Projectname string `json:"projectName"`
	}

	type Associate struct {
		LocationId string `json:"locationId"`
		OprMgrId   string `json:"oprMgrId"`
		FunderId   string `json:"funderId"`
		ProjectId  string `json:"projectId"`
	}

	var p Associate
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode the request body ", "success": false})

		return
	}

	rows, err := DB.Query(`select id,projectName from project p where locationID = ? and operations_manager = ? and endDate < CURRENT_DATE() and funderID = ? `, p.LocationId, p.OprMgrId, p.FunderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
		return
	}
	defer rows.Close()

	// create slice to print output formate
	var listpro []data
	var id, projectName string
	for rows.Next() {
		err := rows.Scan(&id, &projectName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
			return
		}
		listpro = append(listpro, data{Id: id, Projectname: projectName})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": listpro, "success": true, "message": "successfully"})
}
