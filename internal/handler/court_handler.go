package handler

import (
	"futsal-booking/internal/model"
	"futsal-booking/internal/repository"
	"strconv"
	"strings"
	"path/filepath"
	"os"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourtHandler struct {
	Repo *repository.CourtRepository
}

func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	courts, total, err := h.Repo.GetAllCourtsPaginated(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var result []gin.H
	for _, court := range courts {
		result = append(result, gin.H{
			"id": court.ID,
			"name": court.Name,
			"description": court.Description,
			"price_per_hour": court.PricePerHour,
			"image": court.Image,
			"features": court.Features,
			"created_at": court.CreatedAt,
			"updated_at": court.UpdatedAt,
		})
	}
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(200, gin.H{"message": "Courts retrieved successfully", "data": result, "pagination": gin.H{"page": page, "limit": limit, "total": total, "total_pages": totalPages}})
}

func (h *CourtHandler) CreateCourt(c *gin.Context) {
	name := c.PostForm("name")
	price := c.PostForm("price")
    
	if name == "" || price == "" {
		c.JSON(400, gin.H{"error": "Name and Price are required"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Image is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExts[ext] {
		c.JSON(400, gin.H{"error": "Invalid image format"})
		return
	}
	uploadPath := "uploads/courts"
	os.MkdirAll(uploadPath, os.ModePerm)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadPath, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload image"})
		return
	}

	priceInt, err := strconv.Atoi(price)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid price format"})
		return
	}



	// Ambil dan filter features agar tidak ada nilai kosong/null
	featuresArr := c.PostFormArray("features[]")
	var filteredFeatures []string
	for _, f := range featuresArr {
		if f != "" && f != "null" {
			filteredFeatures = append(filteredFeatures, f)
		}
	}
	if len(filteredFeatures) == 0 {
		// fallback jika frontend kirim "features" saja
		singleFeature := c.PostForm("features")
		if singleFeature != "" && singleFeature != "null" {
			filteredFeatures = []string{singleFeature}
		}
	}
	if filteredFeatures == nil {
		filteredFeatures = []string{}
	}

	court := model.Court{
		ID:           uuid.New(),
		Name:         name,
		Description:  c.PostForm("description"),
		PricePerHour: priceInt,
		Image:        filepath,
		Features:     model.StringArray(filteredFeatures),
	}
	if err := h.Repo.CreateCourt(&court); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Court created successfully", "data": court})
}

func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	court, err := h.Repo.GetCourtByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Court not found"})
		return
	}
	result := gin.H{
		"id": court.ID,
		"name": court.Name,
		"description": court.Description,
		"price_per_hour": court.PricePerHour,
		"image": court.Image,
		"features": court.Features,
		"created_at": court.CreatedAt,
		"updated_at": court.UpdatedAt,
	}
	c.JSON(200, gin.H{"message": "Court retrieved successfully", "data": result})
}

func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}

	// Get the existing court for old image cleanup
	oldCourt, err := h.Repo.GetCourtByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Court not found"})
		return
	}

	name := c.PostForm("name")
	price := c.PostForm("price")
	description := c.PostForm("description")

	var imagePath string = oldCourt.Image
	file, err := c.FormFile("image")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
		if !allowedExts[ext] {
			c.JSON(400, gin.H{"error": "Invalid image format"})
			return
		}
		uploadPath := "uploads/courts"
		os.MkdirAll(uploadPath, os.ModePerm)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		imagePath = filepath.Join(uploadPath, filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload image"})
			return
		}
		// Delete old image if exists and different
		if oldCourt.Image != "" && oldCourt.Image != imagePath {
			_ = os.Remove(oldCourt.Image)
		}
	}

	priceInt, err := strconv.Atoi(price)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid price format"})
		return
	}

	featuresArr := c.PostFormArray("features[]")
	var filteredFeatures []string
	for _, f := range featuresArr {
		if f != "" && f != "null" {
			filteredFeatures = append(filteredFeatures, f)
		}
	}
	if len(filteredFeatures) == 0 {
		singleFeature := c.PostForm("features")
		if singleFeature != "" {
			filteredFeatures = []string{singleFeature}
		}
	}

	court := model.Court{
		ID:           id,
		Name:         name,
		Description:  description,
		PricePerHour: priceInt,
		Features:     model.StringArray(filteredFeatures),
		Image:        imagePath,
	}
	if err := h.Repo.UpdateCourt(&court); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Court updated successfully", "data": court})
}

func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	// Get court to find image path
	court, err := h.Repo.GetCourtByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Court not found"})
		return
	}
	// Delete court from DB
	if err := h.Repo.DeleteCourt(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if court.Image != "" {
		_ = os.Remove(court.Image)
	}
	c.JSON(200, gin.H{"message": "Court deleted successfully"})
}