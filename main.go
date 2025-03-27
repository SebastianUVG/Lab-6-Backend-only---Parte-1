package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

var matches = []Match{
	{ID: 1, HomeTeam: "Barcelona", AwayTeam: "Real Madrid", MatchDate: "2025-04-01"},
}

func getMatch(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, matches)
}

func createMatch(c *gin.Context) {
	var newMatch Match

	if err := c.BindJSON(&newMatch); err != nil {
		return
	}

	matches = append(matches, newMatch)
	c.IndentedJSON(http.StatusCreated, newMatch)

}

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/matches", getMatch)
		api.POST("/matches", createMatch)
	}

	router.Run("localhost:8080")
}
