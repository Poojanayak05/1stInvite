package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Occupation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetOccupations(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	pageno := 1

	noOfRecords := 100000
	offset := (pageno - 1) * noOfRecords

	totalPagesQuery := "SELECT COUNT(*) FROM occupations"
	var totalRows int
	err := DB.QueryRow(totalPagesQuery).Scan(&totalRows)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to get total rows", "success": false, "error": err})
		return
	}

	// totalPages := int(totalRows/noOfRecords) + 1

	fields := "id, name"
	query := "SELECT " + fields + " FROM occupations ORDER BY name LIMIT ?, ?"
	rows, err := DB.Query(query, offset, noOfRecords)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to fetch occupations", "success": false, "error": err})
		return
	}
	defer rows.Close()

	var occupations []Occupation
	for rows.Next() {
		var occupation Occupation
		err := rows.Scan(&occupation.ID, &occupation.Name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
			continue
		}
		occupations = append(occupations, occupation)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"list": occupations, "total_count": totalRows, "code": http.StatusOK, "message": "Successfully", "success": true})
}
