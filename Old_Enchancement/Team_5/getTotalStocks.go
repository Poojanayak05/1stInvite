package Team_5

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetTotalStocks retrieves the total stocks based on the given parameters.

func GetTotalStocks(w http.ResponseWriter, r *http.Request, DB *sql.DB) {

	// Check if the request method is GET
	if r.Method != http.MethodPost {
		w.WriteHeader(405) // Return 405 Method Not Allowed.
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
		return
	}

	// Define the request structure
	type Request struct {
		From_Date   string `json:"from_date"`
		To_Date     string `json:"to_date"`
		User_id     string `json:"user_id"`
		Role_id     string `json:"role_id"`
		Project_id  string `json:"project_id"`
		Taluk_id    string `json:"taluk_id"`
		District_id string `json:"district_id"`
		Funder_id   string `json:"funder_id"`
	}

	decoder := json.NewDecoder(r.Body)
	var data Request
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false, "error": err})
		return
	}
	User_id, _ := strconv.Atoi(data.User_id)
	Role_id, _ := strconv.Atoi(data.Role_id)

	if User_id != 0 && Role_id != 0 {

		var condition, condition_role, join string

		if Role_id == 1 || Role_id == 2 || Role_id == 3 || Role_id == 4 || Role_id == 5 || Role_id == 6 || Role_id == 12 || Role_id == 13 {
			if data.Project_id != "" {
				condition = fmt.Sprintf("AND prj.id = %s", data.Project_id)
				condition_role = fmt.Sprintf("prj.id = %s", data.Project_id)
			}

			if data.Funder_id != "" {
				condition = fmt.Sprintf("AND prj.funderID = %s", data.Funder_id)
				condition_role = fmt.Sprintf("prj.funderID = %s", data.Funder_id)
			}

			if data.District_id != "" {
				condition = fmt.Sprintf("AND loc.parentId = %s", data.District_id)
				condition_role = fmt.Sprintf("loc.parentId = %s", data.District_id)
				join = fmt.Sprintf("LEFT JOIN location district ON district.parentId = %s", data.District_id)
			}

			if data.Taluk_id != "" {
				condition = fmt.Sprintf("AND prj.locationID = %s", data.Taluk_id)
				condition_role = fmt.Sprintf("prj.locationID = %s", data.Taluk_id)
			}

			// checking and setting conditions based on role and provided filters...
			if data.From_Date != "" && data.To_Date != "" {
				where_stock := fmt.Sprintf("AND DATE(pr_stock_con.CreatedAt) BETWEEN %s AND %s", data.From_Date, data.To_Date)
				fmt.Println(where_stock)
			}

			var fields string
			var query_project string

			// Construct the query based on the role and conditions
			if Role_id == 1 || Role_id == 2 {
				fields = "DISTINCT prj.id"

				if condition_role == "" {
					query_project = fmt.Sprintf("SELECT %s FROM project prj %s", fields, join)
				} else {
					query_project = fmt.Sprintf("SELECT %s FROM project prj %s WHERE %s", fields, join, condition_role)
				}
			} else if Role_id == 3 {
				fields = "DISTINCT prj.id"
				query_project = fmt.Sprintf("SELECT %s FROM employee emp JOIN project prj ON emp.id = prj.operations_manager %s WHERE emp.empRole = 4 AND emp.supervisorId = %d %s", fields, join, User_id, condition)
			} else if Role_id == 12 {
				fields = "DISTINCT prj.id"
				query_project = fmt.Sprintf("SELECT %s FROM employee emp JOIN project prj ON emp.id = prj.operations_manager %s WHERE emp.empRole = 4 AND emp.supervisorId = %d %s", fields, join, User_id, condition)
			} else if Role_id == 4 {
				fields = "DISTINCT prj.id"
				query_project = fmt.Sprintf("SELECT %s FROM project prj %s WHERE prj.operations_manager = %d %s", fields, join, User_id, condition)
				//strconv.Itoa(User_id)
			} else if Role_id == 5 || Role_id == 6 || Role_id == 13 {
				fields = "em_pr.project_id AS id"
				query_project = fmt.Sprintf("SELECT %s FROM project_emps em_pr LEFT JOIN project prj ON em_pr.project_id = prj.id %s WHERE emp_id IN (%d) %s", fields, join, User_id, condition)
			} else if Role_id == 7 {
				fields = "prj.id"
				query_project = fmt.Sprintf("SELECT %s FROM project prj LEFT JOIN employee emp ON prj.driverID = emp.id WHERE prj.driverID = %d", fields, User_id)
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Role ID not found", "success": false})
				return
			}

			// Query the database to retrieve project IDs based on the constructed query
			row, err := DB.Query(query_project)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No records found", "success": false, "error": err})
				return
			}
			defer row.Close()

			var projectIDs []string

			// Process the retrieved project IDs
			for row.Next() {
				var projectID string
				err := row.Scan(&projectID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No records found", "success": false, "error": err})
					return
				}

				projectIDs = append(projectIDs, projectID)
			}
			projectIDsStr := strings.Join(projectIDs, ",")
			if projectIDsStr == "" {
				projectIDsStr = "''"
			}

			stockConsumed := "''"
			stockAdded := "''"

			// Check if from_date and to_date parameters are provided
			if fromdate, todate := data.From_Date, data.To_Date; fromdate != "" && todate != "" {
				conditionStock := fmt.Sprintf(" WHERE DATE(CreatedAt) BETWEEN '%s' AND '%s' AND type = '1' AND project_id IN(%s) AND stock_id = st_item.id", fromdate, todate, projectIDsStr)
				stockConsumed = fmt.Sprintf("(SELECT IFNULL(SUM(quantity), 0) FROM project_stock_consumption %s)", conditionStock)

				// conditionAdded := fmt.Sprintf(" WHERE DATE(CreatedAt) BETWEEN '%s' AND '%s' AND type = '0' AND project_id IN(%s) AND stock_id = st_item.id", fromdate, todate, projectIDsStr)
				// stockAdded = fmt.Sprintf("(SELECT IFNULL(SUM(quantity), 0) FROM project_stock_consumption %s)", conditionAdded)
			}

			fields = "st_item.id AS stock_id, st_item.name, IFNULL(SUM(CASE WHEN pr_stock_con.type = '0' THEN pr_stock_con.quantity ELSE -pr_stock_con.quantity END), 0) AS current_stock, " + stockConsumed + " AS consumed, " + stockAdded + " AS stock_added"
			query := "SELECT " + fields + " FROM stock_items st_item LEFT JOIN project_stock_consumption pr_stock_con ON st_item.id = pr_stock_con.stock_id AND pr_stock_con.project_id IN(" + projectIDsStr + ") GROUP BY st_item.id ORDER BY st_item.name"

			rows, err := DB.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No records found", "success": false, "error": err})
				return
			}
			defer rows.Close()

			type StockItem struct {
				Stock_id      string `json:"stock_id"`
				Name          string `json:"name"`
				Current_stock string `json:"current_stock"`
				Consumed      string `json:"consumed"`
				Stock_added   string `json:"stock_added"`
			}

			var stocks []StockItem

			// Process the retrieved stock items
			for rows.Next() {
				var stock StockItem
				err := rows.Scan(&stock.Stock_id, &stock.Name, &stock.Current_stock, &stock.Consumed, &stock.Stock_added)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "No records found", "success": false, "error": err})
					return
				}

				stocks = append(stocks, stock)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"data": stocks, "code": 200, "success": true, "message": "Successfully"})

		}
	}
}
