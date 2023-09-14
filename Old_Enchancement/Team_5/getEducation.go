package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// var DB *sql.DB

type Education struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Response6 struct {
	List       interface{} `json:"list"`
	TotalCount string      `json:"total_count"`
	Code       int         `json:"code"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
}

func GetEducation(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	// Execute the query and fetch the data
	rows, err := DB.Query("SELECT * FROM education")
	if err != nil {
		
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
		return
	}
	defer rows.Close()

	// Create a slice to hold the retrieved data
	var data1 []Education

	// Iterate over the rows and append data to the slice
	for rows.Next() {
		var row Education
		err := rows.Scan(&row.ID, &row.Name)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}
		data1 = append(data1, row)
	}
	var request struct {
		PageNum int `json:"pageNum"`
	}
	pageNum := request.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}

	noOfRecords := 1000000
	offset := (pageNum - 1) * noOfRecords
	var totalRows string
	err = DB.QueryRow("SELECT COUNT(*) FROM education").Scan(&totalRows)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
		return
	}
	totalRowsInt, err := strconv.Atoi(totalRows) //converting string to int
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
		return
	}

	totalPages := totalRowsInt / noOfRecords
	if totalRowsInt%noOfRecords != 0 {
		totalPages++
	}
	fields := "id, name"
	query := fmt.Sprintf("SELECT %s FROM education ORDER BY id LIMIT %d, %d", fields, offset, noOfRecords)
	//Executing the query
	rows1, err := DB.Query(query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false})
		return
	}
	defer rows1.Close()
	//Reponse structure
	response := Response6{
		List:       data1,
		TotalCount: totalRows,
		Code:       200,
		Success:    true,
		Message:    "successfully",
	}

	// json.NewEncoder(w).Encode(map[string]interface{}{"list": data1, "total_count": totalRows, "code": 200, "Message": "Successfully", "success": true})
	json.NewEncoder(w).Encode(response)
}
