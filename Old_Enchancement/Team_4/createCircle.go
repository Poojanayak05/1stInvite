package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type CircleRequest struct {
	ProjectID        string    `json:"project_id"`
	CircleName       string    `json:"circle_name"`
	CircleDate       string    `json:"circle_date"`
	GelathiCreatedID string    `json:"gelathi_created_id"`
	EnrolledGelathi  []Gelathi `json:"gelathi"`
}
type Gelathi struct {
	ID         string `json:"id"`
	IsSelected string `json:"isSelected"`
}

func CreateCircle(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var request CircleRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	var Proj int
	err1 := DB.QueryRow("SELECT project_id from project_emps where emp_id=? and project_id=?", request.GelathiCreatedID, request.ProjectID).Scan(&Proj)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID for the given User ID", "success": false})
		} else {
			// Handle other errors
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}

	// Check if Proj is 0, indicating no valid project ID was found
	if Proj == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID for the given User ID", "success": false})
		return
	}
	// if err1 != nil {

	// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project Id for the given User Id", "success": false})
	// 	return
	// }
	insertCircleQuery := "INSERT INTO circle (project_id, circle_name, circle_date, gelathi_created_id) VALUES (?, ?, ?, ?)"
	result, err := DB.Exec(insertCircleQuery, request.ProjectID, request.CircleName, request.CircleDate, request.GelathiCreatedID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to insert circle", "success": false, "error": err})
		return
	}
	circleID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to get last inserted ID", "success": false, "error": err})
		return
	}

	insertGelathiCircleQuery := "INSERT INTO gelathi_circle (circle_id, gelathi_id) VALUES (?, ?)"
	for _, gelathi := range request.EnrolledGelathi {
		Selected, err := strconv.ParseBool(gelathi.IsSelected)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to get last inserted ID", "success": false, "error": err})
			return
		}
		if Selected {
			_, err = DB.Exec(insertGelathiCircleQuery, circleID, gelathi.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to insert gelathi_circle", "success": false, "error": err})
				return
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Circle Created Successfully", "success": true})
}
