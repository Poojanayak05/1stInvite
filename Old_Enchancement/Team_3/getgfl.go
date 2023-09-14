package Team_3

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee10 struct {
	Id             int    `json:"id"`
	CountryID      string `json:"countryID"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Gender         string `json:"gender"`
	Doj            string `json:"doj"`
	OfficeMailId   string `json:"officeMailId"`
	PersonalMailId string `json:"personalMailId"`
	ContactNum     string `json:"contactNum"`
	Address        string `json:"address"`
	Address2       string `json:"address2"`
	Address3       string `json:"address3"`
	Pincode        string `json:"pincode"`
	EmpRole        string `json:"empRole"`
	License_number string `json:"license_number"`
	SupervisorId   string `json:"supervisorId"`
	Profile_pic    string `json:"profile_pic"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
	CreatedBy      string `json:"createdBy"`
	LastUpdatedAt  string `json:"lastUpdatedAt"`
	LastUpdatedBy  string `json:"LastUpdatedBy"`
	WorkNum        string `json:"workNum"`
}

type Response1 struct {
	GflList []Employee10 `json:"gfl_list"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

func Getgfl(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"Status": "405 ", "Message": "Method not found"})
		return
	}
	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}

	if user_id, ok := request["user_id"].(string); ok {
		if user_id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Send the user_id Parameters", "success": false})
			return
		}
		rows, err := DB.Query("SELECT * from employee e WHERE e.empRole =13 and e.supervisorId =?", user_id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "Bad Request", "success": false})
			return
		}

		defer rows.Close()

		var project Employee10
		var Projects []Employee10
		for rows.Next() {
			err := rows.Scan(
				&project.Id,
				&project.CountryID,
				&project.First_name,
				&project.Last_name,
				&project.Gender,
				&project.Doj,
				&project.OfficeMailId,
				&project.PersonalMailId,
				&project.ContactNum,
				&project.Address,
				&project.Address2,
				&project.Address3,
				&project.Pincode,
				&project.EmpRole,
				&project.License_number,
				&project.SupervisorId,
				&project.Profile_pic,
				&project.Status,
				&project.CreatedAt,
				&project.CreatedBy,
				&project.LastUpdatedAt,
				&project.LastUpdatedBy,
				&project.WorkNum)
			if err != nil {
				log.Println(err)
				//json.NewEncoder(w).Encode(map[string]interface{}{"Status": "400 Bad Request", "Message": err})

			}

			Projects = append(Projects, project)
		}

		Response := Response1{
			GflList: Projects,
			Success: true,
			Message: "Successfully",
		}
		json.NewEncoder(w).Encode(Response)
	}
}
