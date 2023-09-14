package connection

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func Connect() *sql.DB {
	// //username:password@tcp(url)/schema
	dbUsername := "bdms_staff_admin"
	dbPassword := "sfhakjfhyiqundfgs3765827635"
	dbHost := "buzzwomendatabase-new.cixgcssswxvx.ap-south-1.rds.amazonaws.com"
	dbPort := 3306
	dbName := "bdms_staff"
	// Define the maximum number of allowed connections
	maxConnections := 100000
	// Create a new MySQL connection pool
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(err)

	}
	// Set the maximum number of open connections

	db.SetMaxOpenConns(maxConnections)
	// Test the database connection

	err = db.Ping()
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database!")
	// Periodically check the number of connections and restart if it exceeds the maximum
	go func() {
		for {
			stats := db.Stats()

			if stats.OpenConnections >= maxConnections {
				log.Println("Reached maximum connections, restarting...")
				db.Close()
				db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName))
				if err != nil {
					log.Fatal(err)
				}
				db.SetMaxOpenConns(maxConnections)
				log.Println("Database connection restarted")
			}
			time.Sleep(5 * time.Second) // Adjust the sleep duration as needed
		}
	}()
	log.Println("Database connection established")
return db
}
