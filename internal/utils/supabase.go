package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func UploadImageToSupabase(fileHeader *multipart.FileHeader, bucket string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if supabaseURL == "" || serviceKey == "" {
		return "", fmt.Errorf("Supabase credentials not configured")
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// Create Supabase storage URL
	storageURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, filename)

	// Create request
	req, err := http.NewRequest("POST", storageURL, bytes.NewBuffer(fileContent))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", fileHeader.Header.Get("Content-Type"))
	req.Header.Set("x-upsert", "true")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload to Supabase: %s", string(body))
	}

	// Return public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucket, filename)
	return publicURL, nil
}

func DeleteImageFromSupabase(imageURL, bucket string) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if supabaseURL == "" || serviceKey == "" {
		return fmt.Errorf("Supabase credentials not configured")
	}

	// Extract filename from URL
	// URL format: https://xwnfzjgzmznzexvtxhon.supabase.co/storage/v1/object/public/futsal_db/filename
	filename := imageURL[len(supabaseURL)+len("/storage/v1/object/public/")+len(bucket)+1:]

	// Create delete URL
	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucket, filename)

	// Create request
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete from Supabase: status %d", resp.StatusCode)
	}

	return nil
}
