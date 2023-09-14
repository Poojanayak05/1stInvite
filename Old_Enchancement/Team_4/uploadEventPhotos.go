package Team_4

//Done by Sushmitha
import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PhotoUploadRequest struct {
	EventID string `json:"event_id"`
	Photos  string `json:"photos"`
}

func UploadEventPhotos(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusNotFound, "message": "Method not found", "success": false})
		return
	}
	var request PhotoUploadRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Request Body", "success": false})
		return
	}

	if request.EventID == "" || request.Photos == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Invalid Input passed", "success": false})
		return
	}
	// Decode base64 image data
	imgData, err := base64.StdEncoding.DecodeString(request.Photos)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to decode image data", "success": false})
		return
	}

	// Generate unique image name
	imageName := time.Now().Add(5*time.Hour+30*time.Minute).Format("20060102150405") + request.EventID + ".jpg"

	// Create or open the images folder
	err = os.MkdirAll("images", os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to create images folder", "success": false, "error": err})
		return
	}

	// Save the image file
	file, err := os.Create("images/" + imageName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to save image file", "success": false, "error": err})
		return
	}

	defer file.Close()

	_, err = file.Write(imgData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to write image file", "success": false, "error": err})
		return
	}
	imagePath := "images/" + imageName
	file1, err := os.Open(imagePath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"Message": err, "Status Code": "400 Bad Request"})
		return
	}
	defer file1.Close()
	// Update the database with the image path
	updateQuery := "UPDATE tbl_poa SET photo1 = ? WHERE id = ?"
	_, err = DB.Exec(updateQuery, "images/"+imageName, request.EventID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to update photo path in the database", "success": false, "error": err})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusOK, "message": "Photos Uploaded Successfully", "success": true})

}
