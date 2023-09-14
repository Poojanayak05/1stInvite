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

type Request struct {
	Name        string `json:"name"`
	AllDay      string `json:"all_day"`
	Date        string `json:"date"`
	Date2       string `json:"date2"`
	UserID      string `json:"user_id"`
	Description string `json:"description"`
}

func CreateEvent(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var User int
	err1 := DB.QueryRow("SELECT id from employee where id=?", request.UserID).Scan(&User)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid User ID", "success": false})
		} else {
			// Handle other errors that may have occurred during query execution
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}
	// Check if Proj is 0, indicating no valid project ID was found
	if User == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid User ID", "success": false})
		return
	}
	if request.Name == "" || request.AllDay == "" || request.Date == "" || request.Date2 == "" || request.UserID == "" || request.Description == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Missing Required Fields", "success": false})
		return
	}
	// Duplicate record for date & time validation
	if request.Date != time.Now().Format("1970-01-01 00:00:00") {
		countQuery := "SELECT COUNT(id) as count FROM tbl_poa WHERE date = ? AND user_id = ? AND status != '2'"

		var count int
		err := DB.QueryRow(countQuery, request.Date, request.UserID).Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error", "success": false, "error": err})
			return
		}

		if count > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Poa already exists for date amd time", "success": false})
			return
		}
	}

	insertQuery := "INSERT INTO tbl_poa (type, name, all_day, date, date2, user_id, description) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err3 := DB.Exec(insertQuery, 3, request.Name, request.AllDay, request.Date, request.Date2, request.UserID, request.Description)
	if err3 != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err3)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid input syntax", "success": false, "error": err3})
		return
	}

	insertIDQuery := "SELECT LAST_INSERT_ID()"
	var insertID int64
	err = DB.QueryRow(insertIDQuery).Scan(&insertID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid input syntax", "success": false, "error": err})
		return
	}

	update := fmt.Sprintf("UPDATE tbl_poa SET primary_id = '%d' WHERE id= '%d'", insertID, insertID)

	_, err2 := DB.Exec(update)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err2)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Event Added Unsuccessfully", "success": false, "error": err2})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Event Added Successfully", "success": true})

}
