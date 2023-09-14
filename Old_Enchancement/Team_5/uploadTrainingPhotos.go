package Team_5

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Requestb struct {
	ProjectID string   `json:"project_id"`
	TbID      string   `json:"tb_id"`
	TrainerID string   `json:"trainer_id"`
	Day       string   `json:"day"`
	Photos    []string `json:"photos"`
}

type Photo struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

type Response9 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func UploadTrainingPhotos(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}

	var req Requestb
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	// Validate request fields
	if req.ProjectID == "" || req.TbID == "" || req.TrainerID == "" || req.Day == "" || len(req.Photos) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid input passed")
		return
	}

	// Define condition for ordering
	condition := ""
	if req.Day == "1" {
		condition = "ORDER BY tbl_poa.id ASC LIMIT 1"
		// Define table fields and SQL query
		tableFields := "tbl_poa.id AS tblPoa_id,COALESCE(tbl_poa.photo1,'')"
		query := fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa WHERE tbl_poa.project_id = %s AND tbl_poa.tb_id = %s AND tbl_poa.type = '1' AND tbl_poa.user_id = %s %s", tableFields, req.ProjectID, req.TbID, req.TrainerID, condition)

		// Execute the SQL query
		row := DB.QueryRow(query)

		// Fetch data from the query result
		// var tblData struct {
		// 	TblPoaID int    `json:"tblPoa_id"`
		// 	Photo1   string `json:"photo1"`
		// }
		var tbl_poa_id int
		var photo Photo
		err = row.Scan(&tbl_poa_id, &photo.Photo1)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows found
				w.WriteHeader(http.StatusNotFound)
				resp := Response9{
					Code:    http.StatusNotFound,
					Message: "No data found",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
			} else {
				// Other error occurred
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to fetch data:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to fetch data",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
			}
			return
		}
		// Create or open the images folder
		err = os.MkdirAll("images", os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := Response9{
				Code:    http.StatusBadRequest,
				Message: "Failed to create images folder",
				Success: false,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		var count int
		err := DB.QueryRow("SELECT count(photo1) from tbl_poa where tb_id=? and tb_id=primary_id", req.TbID).Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "photo 1 already exists", "success": false})
			return
		}
		// Process and save photos
		photos := req.Photos
		if photos == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid photo", "success": false})
			return
		}
		if count == 0 {
			img := photos[0]
			// Decode base64 image data
			img = strings.Replace(img, "data:image/png;base64,", "", 1)
			img = strings.Replace(img, "data:image/jpeg;base64,", "", 1)
			img = strings.Replace(img, " ", "+", -1)
			data, err := base64.StdEncoding.DecodeString(img)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to decode base64 image:", err)
				return
			}

			// Generate unique image name
			imageName := fmt.Sprintf("%s%s.jpg", time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405"), req.TbID)
			// file, err := os.Create("images/" + imageName)
			// if err != nil {
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to save image file", "success": false, "error": err})
			// 	return
			// }

			// defer file.Close()
			// Save image to disk
			imagePath := "images/" + imageName
			err = ioutil.WriteFile(imagePath, data, 0644)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to save image:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to save image",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return

			}

			// Update photo field in the database

			updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo1 = '%s' WHERE id = '%d'", imagePath, tbl_poa_id)
			_, err = DB.Exec(updateQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to update photo field:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update photo field",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		} else {
			img := photos[0]
			// Decode base64 image data
			img = strings.Replace(img, "data:image/png;base64,", "", 1)
			img = strings.Replace(img, "data:image/jpeg;base64,", "", 1)
			img = strings.Replace(img, " ", "+", -1)
			data, err := base64.StdEncoding.DecodeString(img)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to decode base64 image:", err)
				return
			}

			// Generate unique image name
			imageName := fmt.Sprintf("%s%s.jpg", time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405"), req.TbID)
			file, err := os.Create("images/" + imageName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to save image file", "success": false, "error": err})
				return
			}

			defer file.Close()
			// Save image to disk
			imagePath := "images/" + imageName
			err = ioutil.WriteFile(imagePath, data, 0644)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to save image:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to save image",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return

			}

			// Update photo field in the database

			updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo2 = '%s' WHERE id = '%d'", imagePath, tbl_poa_id)
			_, err = DB.Exec(updateQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to update photo field:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update photo field",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
		// Send success response
		resp := Response9{
			Code:    200,
			Message: "Photos Uploaded Successfully",
			Success: true,
		}

		// Convert response to JSON
		respJSON, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Failed to marshal response JSON:", err)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write response JSON
		_, err = w.Write(respJSON)
		if err != nil {
			log.Println("Failed to write response:", err)
		}

	} else {
		condition = "ORDER BY tbl_poa.id DESC LIMIT 1"
		tableFields := "tbl_poa.id AS tblPoa_id,COALESCE(tbl_poa.photo1,'')"
		query := fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa WHERE tbl_poa.project_id = %s AND tbl_poa.tb_id = %s AND tbl_poa.type = '1' AND tbl_poa.user_id = %s %s", tableFields, req.ProjectID, req.TbID, req.TrainerID, condition)

		// Execute the SQL query
		row := DB.QueryRow(query)

		// Fetch data from the query result
		// var tblData struct {
		// 	TblPoaID int    `json:"tblPoa_id"`
		// 	Photo1   string `json:"photo1"`
		// }
		var tbl_poa_id int
		var photo Photo
		err = row.Scan(&tbl_poa_id, &photo.Photo1)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows found
				w.WriteHeader(http.StatusNotFound)
				resp := Response9{
					Code:    http.StatusNotFound,
					Message: "No data found",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
			} else {
				// Other error occurred
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to fetch data:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to fetch data",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
			}
			return
		}
		err = os.MkdirAll("images", os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := Response9{
				Code:    http.StatusBadRequest,
				Message: "Failed to create images folder",
				Success: false,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		var count int
		err := DB.QueryRow("SELECT count(photo1) from tbl_poa where tb_id=? and tb_id!=primary_id", req.TbID).Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "photo 1 already exists", "success": false})
			return
		}
		// Process and save photos
		photos := req.Photos
		if photos == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid photo", "success": false})
			return
		}
		fmt.Println(photo.Photo1)
		fmt.Println("'''''''", photos[0])
		if count == 0 {
			img := photos[0]
			// Decode base64 image data
			img = strings.Replace(img, "data:image/png;base64,", "", 1)
			img = strings.Replace(img, "data:image/jpeg;base64,", "", 1)
			img = strings.Replace(img, " ", "+", -1)
			data, err := base64.StdEncoding.DecodeString(img)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to decode base64 image:", err)
				return
			}
			fmt.Println("Hello")
			fmt.Println(">>>>>", data)
			// Generate unique image name
			imageName := fmt.Sprintf("%s%s.jpg", time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405"), req.TbID)
			file, err := os.Create("images/" + imageName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to save image file", "success": false, "error": err})
				return
			}

			defer file.Close()
			// Save image to disk
			imagePath := "images/" + imageName
			err = ioutil.WriteFile(imagePath, data, 0644)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to save image:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to save image",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return

			}

			// Update photo field in the database

			updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo1 = '%s' WHERE id = '%d'", imagePath, tbl_poa_id)
			_, err = DB.Exec(updateQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to update photo field:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update photo field",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		} else {
			img := photos[0]
			// Decode base64 image data
			img = strings.Replace(img, "data:image/png;base64,", "", 1)
			img = strings.Replace(img, "data:image/jpeg;base64,", "", 1)
			img = strings.Replace(img, " ", "+", -1)
			data, err := base64.StdEncoding.DecodeString(img)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to decode base64 image:", err)
				return
			}

			imageName := fmt.Sprintf("%s%s.jpg", time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405"), req.TbID)
			file, err := os.Create("images/" + imageName)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to save image file", "success": false, "error": err})
				return
			}

			defer file.Close()
			// Save image to disk
			imagePath := "images/" + imageName
			err = ioutil.WriteFile(imagePath, data, 0644)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to save image:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to save image",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return

			}

			// Update photo field in the database

			updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo2 = '%s' WHERE id = '%d'", imagePath, tbl_poa_id)
			_, err = DB.Exec(updateQuery)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Failed to update photo field:", err)
				resp := Response9{
					Code:    http.StatusInternalServerError,
					Message: "Failed to update photo field",
					Success: false,
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
		// Send success response
		resp := Response9{
			Code:    200,
			Message: "Photos Uploaded Successfully",
			Success: true,
		}

		// Convert response to JSON
		respJSON, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Failed to marshal response JSON:", err)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write response JSON
		_, err = w.Write(respJSON)
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	}
}
