package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RequestTa struct {
	TaId          []string `json:"ta_id"`
	UserId        string   `json:"user_id"`
	ExtraComments string   `json:"extra_comments"`
	Status        string   `json:"status"`
}

func VerifyTa(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var request RequestTa
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode JSON", "success": false, "error": err})
		return
	}

	userID := request.UserId
	taIDs := request.TaId

	date := time.Now().Add(5*time.Hour + 30*time.Minute).Format("2006/01/02")
	// date := time.Now().Format("2006/01/02")
	intArray := make([]int, len(taIDs))

	for i, str := range taIDs {
		num, err := strconv.Atoi(str)
		if err != nil {
			// Handle the error if the string cannot be converted to an integer
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error converting string to int", "success": false, "error": err})
			continue
		}
		intArray[i] = num
	}
	for _, taID := range intArray {
		// taID := int(taID.(string))
		extraComments := request.ExtraComments
		status := request.Status

		_, err1 := DB.Exec("UPDATE travelling_allowance SET status = ?, verifyed_by = ?, verifyied_date = ?, extra_comments = ? WHERE id = ?", status, userID, date, extraComments, taID)
		if err1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to update travelling allowance", "success": false, "error": err1})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Successfully Updated", "success": true})
}
