package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

var matches = []Match{
	{ID: 1, HomeTeam: "Real Madrid", AwayTeam: "Barcelona", MatchDate: "2023-11-01"},
	{ID: 2, HomeTeam: "Atlético Madrid", AwayTeam: "Sevilla", MatchDate: "2023-11-02"},
}
var lastID = 2

func main() {
	router := gin.Default()

	router.GET("/api/matches", getMatches)
	router.GET("/api/matches/:id", getMatchByID)
	router.POST("/api/matches", createMatch)
	router.PUT("/api/matches/:id", updateMatch)
	router.DELETE("/api/matches/:id", deleteMatch)

	router.Run(":8080")
}

func getMatches(c *gin.Context) {
	c.JSON(http.StatusOK, matches)
}

func getMatchByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for _, match := range matches {
		if match.ID == id {
			c.JSON(http.StatusOK, match)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
}

func createMatch(c *gin.Context) {
	var newMatch Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastID++
	newMatch.ID = lastID
	matches = append(matches, newMatch)

	c.JSON(http.StatusCreated, newMatch)
}

func updateMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedMatch Match
	if err := c.ShouldBindJSON(&updatedMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, match := range matches {
		if match.ID == id {
			updatedMatch.ID = id
			matches[i] = updatedMatch
			c.JSON(http.StatusOK, updatedMatch)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
}

func deleteMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, match := range matches {
		if match.ID == id {
			matches = append(matches[:i], matches[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
}
