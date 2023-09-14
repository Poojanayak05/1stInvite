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

type StockRequest struct {
	ProjectID string  `json:"project_id"`
	BusID     string  `json:"bus_id"`
	Type      string  `json:"type"`
	StockList []Stock `json:"stock_list"`
}

type Stock struct {
	StockId  string `json:"stock_id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

func ConsumeStock(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}

	var req StockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}
	var User int
	err1 := DB.QueryRow("SELECT id from project where id=?", req.ProjectID).Scan(&User)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// Handle the case where no rows were found
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID", "success": false})
		} else {
			// Handle other errors that may have occurred during query execution
			log.Println(err1)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Database error", "success": false})
		}
		return
	}
	// Check if Proj is 0, indicating no valid project ID was found
	if User == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Project ID", "success": false})
		return
	}

	// Iterate over the stock_list
	for _, stock := range req.StockList {
		stockID := stock.StockId
		// quantity := stock.Quantity
		var quantity int
		if stock.Quantity == "" {
			quantity = 0
		} else {
			quantity, err = strconv.Atoi(stock.Quantity)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Database query error", "success": false, "error": err})
				return
			}
		}
		// Check quantity_exist from the database
		var quantityExist int
		err = DB.QueryRow("SELECT IFNULL(SUM(CASE WHEN prj_con.type = '0' THEN prj_con.quantity ELSE - prj_con.quantity END), 0) AS quantity_exist FROM project_stock_consumption prj_con WHERE prj_con.project_id = ? AND prj_con.stock_id = ?", req.ProjectID, stockID).Scan(&quantityExist)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Database query error", "success": false, "error": err})
			return
		}

		if req.Type != "0" && quantity > 0 && quantityExist < quantity {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Quantity should not be greater than Quantity Exist", "success": false, "error": err})
			return
		}
		if quantity != 0 {
			_, err = DB.Exec("INSERT INTO project_stock_consumption (project_id, bus_id, stock_id, quantity, type) VALUES (?, ?, ?, ?, ?)", req.ProjectID, req.BusID, stockID, quantity, req.Type)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Database insert error", "success": false, "error": err})
				return
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Stock Updated Successfully", "success": true})

}
