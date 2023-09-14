package Team_3

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	ProjectID   string   `json:"project_id"`
	GFSessionID string   `json:"gf_session_id"`
	GelathiID   string   `json:"gelathi_id"`
	Photos      []string `json:"photos"`
}

type Photo1 struct {
	Photo1 string `json:"photo1"`
	Photo2 string `json:"photo2"`
}

func UploadGFsessionPhotos(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method Not found", "success": false})
		return
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Bad Request", "success": false, "error": err})
		return
	}
	ProjectID, err := strconv.Atoi(req.ProjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid ProjectID", "success": false, "error": err})
		return
	}
	GFSessionID, err := strconv.Atoi(req.GFSessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid GFSessionID", "success": false, "error": err})
		return
	}
	GelathiID, err := strconv.Atoi(req.GelathiID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid GelathiID", "success": false, "error": err})
		return
	}

	photos := req.Photos
	if photos == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Enter valid photo", "success": false})
		return
	}

	// Decode base64 image data
	// imgData, err := base64.StdEncoding.DecodeString(pooh.Photos[])
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode image data", "success": false})
	// 	return
	// }
	// var resp Response

	if ProjectID != 0 && GFSessionID != 0 && GelathiID != 0 {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
			return
		}
		// Project_ID := req.ProjectID
		GFSession_ID := GFSessionID
		// Gelathi_ID := req.GelathiID
		tblFields := "tbl_poa.id as tblPoa_id,COALESCE(tbl_poa.photo1, '') AS photo"
		query := fmt.Sprintf("SELECT %s FROM tbl_poa tbl_poa WHERE tbl_poa.id = %d", tblFields, GFSession_ID)
		rows, err := DB.Query(query)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusInternalServerError, "message": "Failed to encode response", "success": false})
			return
		}
		defer rows.Close()
		//Fetching the result
		var tblPoaID string
		var photo1 string
		if rows.Next() {
			err := rows.Scan(&tblPoaID, &photo1)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
				return
			}
		}
		// Project_ID := req.ProjectID
		// Create or open the images folder
		err = os.MkdirAll("images", os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to create images folder", "success": false, "error": err})
			return
		}
		if len(photos) == 1 {
			img := photos[0]
			img = strings.Replace(img, "data:image/png;base64,", "", 1)
			img = strings.Replace(img, "data:image/jpeg;base64,", "", 1)
			img = strings.Replace(img, " ", "+", -1)
			data, err := base64.StdEncoding.DecodeString(img)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			fmt.Println(data)
			// query1:=fmt.Sprintf("SELECT id from tbl_poa where")
			imageName := fmt.Sprintf("%s%d.jpg", time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405"), GFSession_ID)
			imagePath := "images/" + imageName
			err = ioutil.WriteFile(imagePath, data, 0644)
			if err != nil {
				log.Println("ERROR>>", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
				return
			}
			if photo1 == "" {
				// Save image file
				err := ioutil.WriteFile("images/"+imageName, data, 0644)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while writing the binary data for a file named imagepath", "success": false})
					return
				}

				// Update tbl_poa with photo1 path
				updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo1 = '%s' WHERE id = %d", imagePath, GFSession_ID)
				_, err = DB.Exec(updateQuery)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while updating table poa with imagePath and GFSession_ID", "success": false})
					return
				}
			} else {
				// Save image file
				err := ioutil.WriteFile("images/"+imageName, data, 0644)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while writing the binary data for a file named imageName", "success": false})
					return
				}

				// Update tbl_poa with photo2 path
				updateQuery := fmt.Sprintf("UPDATE tbl_poa SET photo2 = '%s' WHERE id = %d", imagePath, GFSession_ID)
				_, err = DB.Exec(updateQuery)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while updating table poa with imagePath and GFSession_ID", "success": false})
					return

				}
			}
			w.Header().Set("Content-Type", "application/json")
			// Return JSON response indicating success
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Photos Uploaded Successfully", "success": true})
			return

		} else {

			img1 := photos[0]
			img2 := photos[1]

			// Process and save Image 1
			img1 = strings.Replace(img1, "data:image/png;base64,", "", 1)
			img1 = strings.Replace(img1, " ", "+", -1)
			data1, err := base64.StdEncoding.DecodeString(img1)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while decoding base64", "success": false})
				return
			}

			imageName1 := fmt.Sprintf("%d%d1.jpg", time.Now().Unix(), GFSession_ID)
			imagePath1 := "images/" + imageName1

			err = ioutil.WriteFile(imagePath1, data1, 0644)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while writing the binary data for a file named imagepath1 ", "success": false})
				return
			}

			// Update tbl_poa with photo1 path
			updateQuery1 := fmt.Sprintf("UPDATE tbl_poa SET photo1 = '%s' WHERE id = %d", imagePath1, GFSession_ID)
			_, err = DB.Exec(updateQuery1)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while updating table poa with imagePath1 and GFSession_ID", "success": false})
				return
			}

			// Process and save Image 2
			img2 = strings.Replace(img2, "data:image/png;base64,", "", 1)
			img2 = strings.Replace(img2, " ", "+", -1)
			data2, err := base64.StdEncoding.DecodeString(img2)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while decoding base64", "success": false})
				return
			}

			imageName2 := fmt.Sprintf("%d%d2.jpg", time.Now().Unix(), GFSession_ID)
			imagePath2 := "images/" + imageName2

			err = ioutil.WriteFile(imagePath2, data2, 0644)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while writing the binary data for a file named imagepath2 ", "success": false})
				return
			}

			// Update tbl_poa with photo2 path
			updateQuery2 := fmt.Sprintf("UPDATE tbl_poa SET photo2 = '%s' WHERE id = '%s'", imagePath2, tblPoaID)
			_, err = DB.Exec(updateQuery2)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error while updating table poa with imagePath2 and tblPoaID", "success": false})
				return
			}

			// Log image paths to a file
			// fp, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			// if err != nil {
			// 	log.Println("ERROR>>", err)
			// }
			// defer fp.Close()

			// logLine := fmt.Sprintf("%s,,, %s\n", imagePath1, imagePath2)
			// _, err = fp.WriteString(logLine)
			// if err != nil {
			// 	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "error in logLine Part", "success": false})
			// 	return
			// }
			w.Header().Set("Content-Type", "application/json")
			// Return JSON response indicating success
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Successfully", "success": true})
			return
		}
	}
}
