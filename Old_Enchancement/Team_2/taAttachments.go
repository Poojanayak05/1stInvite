package Team_2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Attachment struct {
	ID        string `json:"id"`
	ImageName string `json:"image_name"`
	TAID      string `json:"ta_id"`
	EmpID     string `json:"emp_id"`
}

func TaAttachments(w http.ResponseWriter, r *http.Request, DB *sql.DB) {
	empID := r.FormValue("emp_id")
	targetDir := "./uploads/"
	// allowTypes := []string{"jpg", "png", "jpeg", "gif", "JPG", "PNG", "JPEG", "GIF"}

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // Set max file size to 10MB
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to parse multipart form", "success": false, "error": err.Error()})
		return
	}

	files := r.MultipartForm.File["file[]"]

	var responseFiles []string
	var errorUpload string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to returns the FileHeader's associated File.", "success": false, "error": err.Error()})
			errorUpload += fileHeader.Filename + " | "
			continue
		}

		defer file.Close()

		fileName := filepath.Base(fileHeader.Filename)


		responseFiles = append(responseFiles, fileName)
		targetFilePath := filepath.Join(targetDir, fileName)


		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in creating directory:", "success": false, "error": err.Error()})
			return
		}

		out, err := os.Create(targetFilePath)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in creating File:", "success": false, "error": err.Error()})
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in copies from src to dst", "success": false, "error": err.Error()})
			errorUpload += fileHeader.Filename + " | "
		}

		// Perform database insert here
		insertQuery := fmt.Sprintf("INSERT INTO taAttachments (emp_id, image_name, ta_id) VALUES (%s, 'https://bdms.buzzwomen.org/appTest/new/uploads/%s', 0)",
			empID, fileName)
		_, err = DB.Exec(insertQuery)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Error in Database insertQuery", "success": false, "error": err.Error()})
		}

		// Print uploaded file names
	}

	// Execute the SELECT query
	query := "SELECT * FROM taAttachments WHERE emp_id = ? AND ta_id =0"
	rows, err := DB.Query(query, empID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to execute SELECT query", "success": false, "error": err.Error()})
		return
	}
	defer rows.Close()

	// Iterate over the result rows and build the response data
	var attachments []Attachment
	for rows.Next() {

		var attachment Attachment
		err := rows.Scan(&attachment.ID, &attachment.ImageName, &attachment.TAID, &attachment.EmpID)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to scan attachment data", "success": false, "error": err.Error()})
		}
		attachments = append(attachments, attachment)
	}

	// Prepare response data

	response := make(map[string]interface{})
	if responseFiles != nil {
		response["files"] = responseFiles
		response["error"] = "Success"
		response["data"] = attachments
		response["success"] = true
		response["message"] = "Success"

	} else {
		response["files"] = responseFiles
		response["success"] = true
		response["message"] = "Success"

	}

	// Convert response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": http.StatusBadRequest, "message": "Failed to marshal response", "success": false, "error": err.Error()})
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send response
	w.Write(responseJSON)

}
