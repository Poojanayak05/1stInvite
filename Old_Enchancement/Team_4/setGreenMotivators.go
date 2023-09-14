package Team_4

//Done by Keerthana
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type req struct {
	ID        string `json:"id"`
	GelathiID string `json:"gelathi_id"`
	TbID      string `json:"tb_id"`
	ProjectID string `json:"projectId"`
}

func SetGreenMotivator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid request body", "success": false})
			return
		}

		var request req
		err = json.Unmarshal(data, &request)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid JSON data", "success": false})
			return
		}

		id, _ := strconv.Atoi(request.ID)
		if id == 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "invalid id", "success": false})
			return
		}

		tbID, _ := strconv.Atoi(request.TbID)
		if tbID == 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "invalid tb_id", "success": false})
			return
		}

		gelathiID, _ := strconv.Atoi(request.GelathiID)

		enrolledDate := time.Now().Add(5*time.Hour + 30*time.Minute).Format("2006-01-02 15:04:05")

		projectID, _ := strconv.Atoi(request.ProjectID)

		checkIfEnrolledPreviously1 := fmt.Sprintf("SELECT id FROM training_participants tp WHERE GreenMotivators = 1 AND tb_id = %d", tbID)
		checkRes1, err := db.Query(checkIfEnrolledPreviously1)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false})
			return
		}

		count := 0
		for checkRes1.Next() {
			count++
		}
		checkRes1.Close()

		if count < 2 {
			checkIfEnrolledPreviously := fmt.Sprintf("SELECT id FROM training_participants tp WHERE id = '%d' AND GreenMotivators = 1 AND tb_id = %d", tbID, tbID)
			checkRes, err := db.Query(checkIfEnrolledPreviously)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false})
				return
			}

			if checkRes.Next() {
				var participantID int
				err := checkRes.Scan(&participantID)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false})
					return
				}

				updateReset := fmt.Sprintf("UPDATE training_participants SET GreenMotivators = '0', GreenMotivatorsDate = '%s' WHERE id = %d", enrolledDate, participantID)
				_, err = db.Exec(updateReset)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute Update SQL query", "success": false})
					return
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Green Motivators Updated Successfully", "success": true})
			} else {
				updateSet := fmt.Sprintf("UPDATE training_participants SET GreenMotivators = '1', gelathi_id = %d, GreenMotivatorsDate = '%s', GreenMotivatorsEnrolledProject = %d WHERE tb_id = %d AND id = '%d'", gelathiID, enrolledDate, int(projectID), tbID, int(id))
				_, err := db.Exec(updateSet)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute Update SQL query", "success": false})
					return
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Green Motivators Updated Successfully", "success": true})
			}
		} else {
			checkIfEnrolledPreviously := fmt.Sprintf("SELECT id FROM training_participants tp WHERE id = '%d' AND GreenMotivators = 1 AND tb_id = %d", id, tbID)
			checkRes, err := db.Query(checkIfEnrolledPreviously)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false})
				return
			}

			if checkRes.Next() {
				var participantID int
				err := checkRes.Scan(&participantID)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan query results", "success": false})
					return
				}

				updateReset := fmt.Sprintf("UPDATE training_participants SET GreenMotivators = '0', GreenMotivatorsDate = '%s' WHERE id = %d", enrolledDate, participantID)
				_, err = db.Exec(updateReset)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to execute SQL query", "success": false})
					return
				}

				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Green Motivators Updated Successfully", "success": true})
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Max Selection is Two", "success": true})

			}
		}
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusMethodNotAllowed, "message": "Method not allowed", "success": false})
	}
}
