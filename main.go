package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type score struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
	Location string `json:"location"`
}

var scores = []score{}

func main() {
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

	if start > len(scores) {
		c.IndentedJSON(http.StatusOK, []score{})
		return
	}
	end := start + count
	if (start + count) > len(scores) {
		end = len(scores)
	}
	c.IndentedJSON(http.StatusOK, scores[start:end])
}

func addScore(c *gin.Context) {
	var newScore score

	if err := c.BindJSON(&newScore); err != nil {
		return
	}

	scores = append(scores, newScore)
	c.IndentedJSON(http.StatusCreated, newScore)
}

func getScoreByUsername(c *gin.Context) {
	username := c.Param("username")

	for _, s := range scores {
		if s.Username == username {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not found"})
}
