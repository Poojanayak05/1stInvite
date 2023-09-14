package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Request1 struct {
	Type      string `json:"type"`
	TbID      string `json:"tb_id"`
	EmpID     string `json:"emp_id"`
	Notes     string `json:"notes"`
	Primaryid string `json:"primary_id"`
}

func CreateNotes(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var request Request1
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}

	if request.TbID != "" { // tb_id is value there
		// var primary_id int
		// err1 := DB.QueryRow("SELECT COALESCE(primary_id,''), from tbl_poa where user_id=? and tb_id=?", request.EmpID, request.TbID).Scan(&primary_id)
		// if err1 != nil {
		// 	// if err1 == sql.ErrNoRows {
		// 	// 	// Handle the case where no rows were found
		// 	// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Primary ID(tb) Id for the given User Id", "success": false})
		// 	// } else {
		// 	// 	// Handle other errors
		// 	// 	log.Println(err1)
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		// 	//}
		// 	return
		// }

		query := "INSERT INTO notes (type, tb_id, emp_id, notes,primary_id) VALUES (?, ?, ?, ?,?)"
		type1, err := strconv.Atoi(request.Type)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Type", "success": false, "error": err})
			return
		}

		note, err := DB.Exec(query, type1, request.TbID, request.EmpID, request.Notes, request.Primaryid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid input syntax", "success": false, "error": err})
			return
		}
		fmt.Println(note)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Note Added Successfully", "success": true})

	} else {
		var tb int
		err1 := DB.QueryRow("SELECT tb_id from tbl_poa where user_id=? and primary_id=?", request.EmpID, request.Primaryid).Scan(&tb)
		if err1 != nil {
			if err1 == sql.ErrNoRows {
				// Handle the case where no rows were found
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Training Batch(tb) Id for the given User Id", "success": false})
			} else {
				// Handle other errors
				log.Println(err1)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
			}
			return
		}

		query := "INSERT INTO notes (type, tb_id, emp_id, notes,primary_id) VALUES (?, ?, ?, ?,?)"
		type1, err := strconv.Atoi(request.Type)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Type", "success": false, "error": err})
			return
		}
		note, err := DB.Exec(query, type1, tb, request.EmpID, request.Notes, request.Primaryid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid input syntax", "success": false, "error": err})
			return
		}
		fmt.Println(note)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Note Added Successfully", "success": true})

	}
	//

}
