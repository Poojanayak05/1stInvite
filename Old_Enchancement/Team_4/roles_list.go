package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"roleName"`
}

func Roles_list(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	query := "SELECT * FROM roles_Master WHERE id NOT IN ('1', '10') ORDER BY id"
	rows, err := DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to execute query", "success": false, "error": err})
		return
	}
	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan row", "success": false, "error": err})
			return
		}
		roles = append(roles, role)
	}

	err = rows.Err()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error occurred while iterating over rows", "success": false, "error": err})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"list": roles, "code": http.StatusOK, "message": "Successfully", "success": true})
}
