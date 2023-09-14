package Team_2

import (
	a "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetGFAssignedBatch(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	type TrainingBatch struct {
		TrainingBatchID   string `json:"training_batch_id"`
		TrainingBatchName string `json:"training_batch_name"`
	}

	type response struct {
		Asso  []string        `json:"Asso"`
		List  []TrainingBatch `json:"list"`
		Count int             `json:"count"`
		Code  int             `json:"code"`

		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"status": "404 Not Found", "Error": "Method not found"})
		return
	}

	type req struct {
		ProjectID    string `json:"project_id"`
		EmpID        string `json:"emp_id"`
		SearchFilter string `json:"search"`
	}
	var data req
	var stringOpsIds1 []string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "message": err.Error()})
		return
	}

	var searchFilter = ""
	searchParam := data.SearchFilter

	if len(data.SearchFilter) > 0 {
		searchFilter = fmt.Sprintf(" and tr_bat.name like '%%%s%%'", searchParam)
	}

	if data.ProjectID != "" && data.EmpID != "" {
		project_id, _ := strconv.Atoi(data.ProjectID)

		if project_id == 0 {
			var resp response

			resp.Asso = []string{}
			resp.List = []TrainingBatch{}
			resp.Message = "successfully"
			resp.Success = true
			resp.Count = 0
			resp.Code = 200
			json.NewEncoder(w).Encode(resp)
			return

		}
		projArray, _ := a.GetAssociatedProjectList(DB, project_id)
		stringOpsIds := make([]string, len(projArray))

		for i, id := range projArray {
			stringOpsIds[i] = strconv.Itoa(id)
		}
		result := strings.Join(stringOpsIds, ",")
		stringOpsIds1 = stringOpsIds
		var roleId int
		query := fmt.Sprintf("SELECT empRole FROM employee e WHERE e.id = %s", data.EmpID)

		err := DB.QueryRow(query).Scan(&roleId)
		if err != nil {
			log.Println(err)
			roleId = 2
		}

		var query2 string

		fields := "gf_bat.training_batch_id, tr_bat.name as training_batch_name"

		if roleId == 13 {
			query2 = "SELECT " + fields + " FROM gf_batches gf_bat " +
				"LEFT JOIN tbl_poa tr_bat ON gf_bat.training_batch_id = tr_bat.tb_id " +
				"WHERE gf_bat.project_id IN (" + result + ") " +
				"AND gf_bat.emp_id IN (SELECT id FROM employee e WHERE e.supervisorId = " + data.EmpID + ") " +
				"AND tr_bat.type = '1'" + searchFilter +
				" GROUP BY gf_bat.training_batch_id, tr_bat.name"
		} else if roleId == 6 {
			query2 = "SELECT " + fields + " FROM gf_batches gf_bat " +
				"LEFT JOIN tbl_poa tr_bat ON gf_bat.training_batch_id = tr_bat.tb_id " +
				"WHERE gf_bat.project_id IN (" + result + ") " +
				"AND gf_bat.emp_id = " + data.EmpID +
				" AND tr_bat.type = '1'" + searchFilter +
				" GROUP BY gf_bat.training_batch_id, tr_bat.name"
		} else {
			query2 = "SELECT " + fields + " FROM gf_batches gf_bat " +
				"LEFT JOIN tbl_poa tr_bat ON gf_bat.training_batch_id = tr_bat.tb_id " +
				"WHERE gf_bat.project_id IN (" + result + ") " +
				"AND tr_bat.type = '1'" + searchFilter +
				" GROUP BY gf_bat.training_batch_id, tr_bat.name"
		}

		rows, err := DB.Query(query2)
		if err != nil {
			log.Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}

		var list []TrainingBatch
		var count int
		for rows.Next() {
			var obj TrainingBatch
			var id int
			err = rows.Scan(&id, &obj.TrainingBatchName)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "sessage": err.Error()})
				return
			}
			obj.TrainingBatchID = strconv.Itoa(id)
			list = append(list, obj)
			count++
		}

		var resp response

		resp.Asso = stringOpsIds1
		resp.List = list
		resp.Message = "successfully"
		resp.Success = true
		resp.Count = count
		resp.Code = 200
		if len(list) == 0 {
			resp.List = []TrainingBatch{}
		}
		json.NewEncoder(w).Encode(resp)
		return

	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "successfully", "success": true})
		return
	}

}
