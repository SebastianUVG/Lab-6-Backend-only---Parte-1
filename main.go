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

func main() {
	router := gin.Default()
	router.GET("/matches", getMatch)
	router.Run("localhost:8080")
}
