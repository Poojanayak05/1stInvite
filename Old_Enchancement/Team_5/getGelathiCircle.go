package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// struct to store the circle fields coming from the database.
type Circle struct {
	CircleID    string  `json:"circle_id"`
	CircleName  string  `json:"circle_name"`
	CircleDate  string  `json:"circle_date"`
	ProjectID   string  `json:"project_id"`
	ProjectName *string `json:"projectName"`
}

// struct to show the response
type Respons struct {
	List       []Circle `json:"list,omitempty"`
	Code       int      `json:"code"`
	TotalCount string   `json:"total_count,omitempty"`
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
}

// struct to store the parameters coming from the body.
type GelathiCircle struct {
	Search    string `json:"search"`
	ProjectId string `json:"project_id"`
	GelathiId string `json:"gelathi_id"`
	PageNo    int    `json:"page_no"`
}

func GetGelathiCircle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var gc GelathiCircle
	response := Respons{}
	// dumping the body into GelathiCircle Struct.
	err := json.NewDecoder(r.Body).Decode(&gc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error While Decoding Body. Bad Request", "error": err})
		return
	}

	if gc.ProjectId == "" {
		response.Code = 406
		response.Message = "Please Send the Required Project Id."
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	if gc.GelathiId == "" {
		response.Code = 406
		response.Message = "Please Send the Required Gelathi Id."
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Search by village name
	searchFilter := ""
	searchParam := gc.Search

	//if search parameter is present then making a condition to search according that string.
	if len(searchParam) > 0 {
		searchFilter = " AND circle_name LIKE '%" + searchParam + "%'"
	}

	if gc.PageNo < 1 {
		gc.PageNo = 1

	}

	noOfRecords := 100
	offset := (gc.PageNo - 1) * noOfRecords

	fields := "empRole"
	query := fmt.Sprintf("SELECT %s FROM employee WHERE status = 1 AND id = %s", fields, gc.GelathiId)
	var roleID int
	// preparing the query to execute.
	stmt, err := db.Prepare(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
		return
	}
	// executng the query.
	rows, err := stmt.Query()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
		return
	}
	for rows.Next() {
		//scanning the role id coming from employee table.
		rows.Scan(&roleID)
	}

	//declaring fields to fetch from the project association table and it is used in the below query.
	fields = "cr.id AS circle_id, cr.circle_name, DATE_FORMAT(cr.circle_date, '%d-%m-%Y') AS circle_date, cr.project_id, UPPER(prj.projectName) AS projectName"
	// based on the role id executing the query.
	if roleID == 6 {
		query := fmt.Sprintf("SELECT associatedProject FROM project_association pa WHERE pa.projectId = %s", gc.ProjectId)
		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "No Content", "error": err})
			return
		}
		defer rows.Close()

		prj := ""
		for rows.Next() {
			var associatedProject string
			//scanning the associated project id from the project association table based on project id.
			err := rows.Scan(&associatedProject)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
				return
			}

			if prj == "" {
				prj = gc.ProjectId
			}
			prj += "," + associatedProject

		}
		if prj == "" {
			prj = ""
		}

		wr := ""
		// if the project id and associated project is not empty then making a condition which is using in the below totalPagesQuery
		if prj != "" {
			wr = " project_id IN (" + prj + ") AND "
		}

		// query to get the total count of circle based on gelathi id, wr condition and search string.
		totalPagesQuery := fmt.Sprintf("SELECT COUNT(*) FROM circle WHERE %s gelathi_created_id = %s and project_id=%s %s", wr, gc.GelathiId, gc.ProjectId, searchFilter)

		// query to get the details of the circle based on gelathi id and search string.
		query2 := fmt.Sprintf("SELECT %s FROM circle cr LEFT JOIN project prj ON cr.project_id = prj.id WHERE %s cr.gelathi_created_id = %s and cr.project_id=%s %s ORDER BY cr.id DESC LIMIT %d, %d ", fields, wr, gc.GelathiId, gc.ProjectId, searchFilter, offset, noOfRecords)
		//executing the query.
		rows, err = db.Query(query2)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "No Content", "error": err})
			return
		}
		defer rows.Close()

		var totalCount string
		//scanning the total count of circle.
		err = db.QueryRow(totalPagesQuery).Scan(&totalCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
			return
		}
		var circles []Circle
		for rows.Next() {
			var circle Circle
			// scanning the circle details.
			err := rows.Scan(&circle.CircleID, &circle.CircleName, &circle.CircleDate, &circle.ProjectID, &circle.ProjectName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
				return
			}
			// appending the circle details one by one in circles list.
			circles = append(circles, circle)
		}
		// setting the response in the Respns struct.
		response.List = circles
		response.Code = 200
		response.TotalCount = totalCount
		response.Success = true
		response.Message = "successfully"
	} else {
		// if role id is not 6 then executing the below query.
		// query to get the total count of circle based on gelathi id and search string only here wr condition is not present.
		totalPagesQuery := fmt.Sprintf("SELECT COUNT(*) FROM circle WHERE project_id = %s %s", gc.ProjectId, searchFilter)

		//query to get the details of the circle based on gelathi id and search string same wr condition is not present here.
		query := fmt.Sprintf("SELECT %s FROM circle cr LEFT JOIN project prj ON cr.project_id = prj.id WHERE cr.project_id = %s %s ORDER BY cr.id DESC LIMIT %d, %d ", fields, gc.ProjectId, searchFilter, offset, noOfRecords)
		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "No Content", "error": err})
			return
		}
		defer rows.Close()
		var totalCount string
		// scanning the total count.
		err = db.QueryRow(totalPagesQuery).Scan(&totalCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
			return
		}

		var circles []Circle
		for rows.Next() {
			var circle Circle
			// scanning the circle details.
			err := rows.Scan(&circle.CircleID, &circle.CircleName, &circle.CircleDate, &circle.ProjectID, &circle.ProjectName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
				return
			}
			// appending the circle details one by one in the circles list.
			circles = append(circles, circle)
		}
		//setting the response in Respons struct.
		response.List = circles
		response.Code = 200
		response.TotalCount = totalCount
		response.Success = true
		response.Message = "successfully"
	}

	// marshaling the response into json.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
		return
	}
	// sending response.
	w.Write(jsonResponse)

}
