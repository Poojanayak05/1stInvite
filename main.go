package main

import (
	// dbs "buzzstaff-go/database"
	h "buzzstaff-go/handler"
	// c "buzzstaff-go/job"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// db := dbs.Connect()
	// c.RunCronJobs(db)
	// fmt.Scanln()
	h.HandleFunc()

}
