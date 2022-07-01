package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

type score struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
	Location string `json:"location"`
}

// var scores = []score{}
var db *gorm.DB

func main() {

	db, _ = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	db.AutoMigrate(&score{})

	router := gin.Default()
	router.GET("/highscore", getScores)
	router.POST("/highscore", addScore)
	router.GET("/highscore/:username", getScoreByUsername)

	router.Run("localhost:8080")
}

func getScores(c *gin.Context) {
	startStr := c.DefaultQuery("start", "0")
	countStr := c.DefaultQuery("count", "10")

	start, _ := strconv.Atoi(startStr)
	count, _ := strconv.Atoi(countStr)

	var scores []score

	db.Order("Score desc").Limit(count).Offset(start).Find(&scores)

	c.IndentedJSON(http.StatusOK, scores)
}

func addScore(c *gin.Context) {
	var newScore score

	if err := c.BindJSON(&newScore); err != nil {
		return
	}

	db.Create(&newScore)

	// scores = append(scores, newScore)
	c.IndentedJSON(http.StatusCreated, newScore)
}

func getScoreByUsername(c *gin.Context) {
	username := c.Param("username")

	var score score
	result := db.Where("Username = ?", username).First(&score)
	if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, score)
	return
}

