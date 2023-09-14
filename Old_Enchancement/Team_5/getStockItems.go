package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// struct to store the stock itmess
type StockItem struct {
	ID       string `json:"stock_id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity_exist"`
}

// struct to show the response
type Response struct {
	Data    []StockItem `json:"data,omitempty"`
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

// struct to collect the payload
type Project struct {
	Project_Id string `json:"project_id"`
}

func GetStockItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var pr Project
	var response Response
	var stockItems []StockItem
	// decoding the body into project struct
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error While Decoding Body. Bad Request", "success": false, "error": err})
		return
	}

	//checking id if it is empty then showing empty data with message successfully. same as php code response.
	if pr.Project_Id == "" {
		response.Code = 406
		response.Message = "Please Send the Required Project Id."
		response.Success = false
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(response)
		return
	}

	// defining the feilds to fetch from the stock_items table & it is used in the below query.
	fields := "st_item.id as stock_id, st_item.name, IFNULL(SUM(CASE WHEN pr_stock_con.type = '0' THEN pr_stock_con.quantity ELSE -pr_stock_con.quantity END), 0) AS quantity_exist"
	query := fmt.Sprintf("SELECT %s FROM stock_items st_item LEFT JOIN project_stock_consumption pr_stock_con ON st_item.id = pr_stock_con.stock_id AND pr_stock_con.project_id = %s GROUP BY st_item.id ORDER BY st_item.name", fields, pr.Project_Id)
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNoContent, "message": "No Content", "error": err})
		return
	}
	defer rows.Close()

	for rows.Next() {
		//declaring the variable of type struct StockItem.
		var stockItem StockItem

		// scanning the stock items.
		err := rows.Scan(&stockItem.ID, &stockItem.Name, &stockItem.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error.", "error": err})
			return
		}
		// appending the stock item one by one in the list.
		stockItems = append(stockItems, stockItem)
	}
	//sending success response.
	response = Response{
		Data:    stockItems,
		Code:    200,
		Success: true,
		Message: "Successfully",
	}
	json.NewEncoder(w).Encode(response)
}
