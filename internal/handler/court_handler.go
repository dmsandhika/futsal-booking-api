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
)

type CourtHandler struct {
	Repo *repository.CourtRepository
}

func (h *CourtHandler) GetAllCourts(c *gin.Context) {
	courts, err := h.Repo.GetAllCourts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Courts retrieved successfully", "data": courts})
}

func (h *CourtHandler) CreateCourt(c *gin.Context) {
	name := c.PostForm("name")
	price := c.PostForm("price")
	
	if(name == "" || price == "") {
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
	uploadPath :="uploads/courts"
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
	
	court := model.Court{
		Name:         name,
		Location:     c.PostForm("location"),
		PricePerHour: priceInt,
		Status:       "active",
		Image:        filepath,
	}
	if err := h.Repo.CreateCourt(&court); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Court created successfully", "data": court})
}

func (h *CourtHandler) GetCourtByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	court, err := h.Repo.GetCourtByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "Court not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Court retrieved successfully", "data": court})
}

func (h *CourtHandler) UpdateCourt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	var court model.Court
	if err := c.ShouldBindJSON(&court); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	court.ID = uint(id)
	if err := h.Repo.UpdateCourt(&court); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Court updated successfully"})
}

func (h *CourtHandler) DeleteCourt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid court ID"})
		return
	}
	if err := h.Repo.DeleteCourt(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Court deleted successfully"})
}