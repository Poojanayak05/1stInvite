package Team_4

//Done by Prathamesh
import (
	a "buzzstaff-go/New_Enchancement/Dashboard"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TrainingBatch1 struct {
	TrainingBatchID int    `json:"training_batch_id"`
	Name            string `json:"name"`
	Flag            int    `json:"flag"`
	Date            string `json:"date"`
}

type TrainingBatch2 struct {
	TrainingBatchID int    `json:"training_batch_id"`
	Name            string `json:"name"`
	Flag            string `json:"flag"`
	EmpID           int    `json:"emp_id"`
}
type Responselist struct {
	TrainingBatchID string `json:"training_batch_id"`
	Name            string `json:"name"`
	Flag            string `json:"flag"`
	//EmpID string `json:"emp_id"`
	Date string `json:"date"`
}

func GetBatchlist(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"Status": "404 Not Found", "Error": "Method not found"})
		return
	}
	var responselistvar []Responselist
	decoder := json.NewDecoder(r.Body)

	type requestBody struct {
		ProjectID string `json:"project_id"`
		EmpID     string `json:"emp_id"`
		PageNum   string `json:"pageNum"`
	}

	var data requestBody

	err := decoder.Decode(&data)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "400 Bad Request", "message": err.Error()})
		return
	}
	pageNo := 0
	if data.PageNum != "" {
		pageNo = 1
	}
	NoOfRecords := 100000
	offSet := (pageNo - 1) * NoOfRecords
	_ = offSet
	if data.EmpID != "" && data.ProjectID != "" {
		projectId, _ := strconv.Atoi(data.ProjectID)
		empId, _ := strconv.Atoi(data.EmpID)
		projArray, _ := a.GetAssociatedProjectList(DB, projectId)
		stringOpsIds := make([]string, len(projArray))
		for i, id := range projArray {
			stringOpsIds[i] = strconv.Itoa(id)
		}

		result := strings.Join(stringOpsIds, ",")
		fields := "COALESCE(tr_btch.tb_id,0) as training_batch_id, COALESCE(UPPER(tr_btch.name),'') as name, 0 as flag, tr_btch.date"
		query := "SELECT " + fields + " FROM tbl_poa tr_btch WHERE tr_btch.project_id in (" + result + ")  and tr_btch.type = 1  AND tr_btch.check_out is not null and  tr_btch.tb_id != tr_btch.primary_id order by name"
		rows, err := DB.Query(query)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err.Error()})
			return
		}
		var flagZerolist []TrainingBatch1
		for rows.Next() {
			var obj TrainingBatch1
			err = rows.Scan(&obj.TrainingBatchID, &obj.Name, &obj.Flag, &obj.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "Bad Request", "message": err.Error()})
				return

			}
			flagZerolist = append(flagZerolist, obj)

		}

		rows.Close()

		fields = "COALESCE(tr_btch.tb_id,0) as training_batch_id, COALESCE(UPPER(tr_btch.name),'') as name, 1 as flag,gf_btch.emp_id"

		empBatches := "select " + fields + " from gf_batches gf_btch left join tbl_poa tr_btch on gf_btch.training_batch_id = tr_btch.tb_id where gf_btch.project_id in (" + result + ") group by tr_btch.tb_id,gf_btch.emp_id"

		rows1, err1 := DB.Query(empBatches)

		if err1 != nil {

			log.Println(err1)

			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(map[string]interface{}{"status": "500 Internal Server Error", "message": err1.Error()})

			return

		}

		var batchesRow []TrainingBatch2

		for rows1.Next() {
			var obj TrainingBatch2
			err = rows1.Scan(&obj.TrainingBatchID, &obj.Name, &obj.Flag, &obj.EmpID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "Bad Request", "message": err.Error()})
				return

			}
			batchesRow = append(batchesRow, obj)
		}

		rows1.Close()
		count := 0

		for _, batch := range batchesRow {
			key := -1

			/*for i, rowBatch := range flagZerolist {
				if rowBatch.TrainingBatchID == batch.TrainingBatchID {
					key = i
					break
				}
			}*/
			for i, rowBatch := range flagZerolist {
				if rowBatch.TrainingBatchID == batch.TrainingBatchID {
					key = i
					break
				}
			}

			if key != -1 {
				// Batch found in flagZerolist
				if batch.EmpID == empId {
					flagZerolist[key].Flag = 1 // Set the flag to 1
				}
			} else {
				// Batch not found in flagZerolist
				if batch.EmpID == empId {
					flagZerolist = append(flagZerolist, TrainingBatch1{
						TrainingBatchID: batch.TrainingBatchID,
						Name:            batch.Name,
						Flag:            1,
						Date:            "", // Set the appropriate value for Date
					})
				}
			}

			count++
		}
		count1 := 0
		for _, batch := range flagZerolist {
			var obj Responselist
			if batch.Flag == 1 {
				count1++

			}
			obj.TrainingBatchID = strconv.Itoa(batch.TrainingBatchID)
			obj.Date = batch.Date
			obj.Flag = strconv.Itoa(batch.Flag)
			obj.Name = batch.Name

			responselistvar = append(responselistvar, obj)

		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"projArray":     stringOpsIds,
			"list":          responselistvar,
			"code":          http.StatusOK,
			"message":       "Successfully",
			"success":       true,
			"total_count":   len(responselistvar),
			"checked_count": count1,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Empid required", "success": true, "total_count": nil})

}
