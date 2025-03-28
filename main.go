package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Match struct {
	ID        string `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

var matches = []Match{
	{ID: "1", HomeTeam: "Barcelona", AwayTeam: "Real Madrid", MatchDate: "2025-04-01"},
}

// Funcion para obtener todos los partidos
func getMatch(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, matches)
}

// Funcion para crear un nuevo partido

//Invoke-RestMethod -Uri http://localhost:8080/api/matches `
//-Method POST `
//-Headers @{"Content-Type"="application/json"} `
//-Body (Get-Content -Raw -Path .\body.json)

func createMatch(c *gin.Context) {
	var newMatch Match

	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos"})
		return
	}

	newMatch.ID = generateID()
	matches = append(matches, newMatch)
	c.IndentedJSON(http.StatusCreated, newMatch)
}

func generateID() string {
	if len(matches) == 0 {
		return "1"
	}
	lastID, _ := strconv.Atoi(matches[len(matches)-1].ID)
	return strconv.Itoa(lastID + 1)

}

// funcion para obtener un partido mediante el uso de la iD
func matchById(c *gin.Context) {
	id := c.Param("id")
	match, err := getMatchID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, match)

}

func getMatchID(id string) (*Match, error) {
	for i, m := range matches {
		if m.ID == id {
			return &matches[i], nil
		}
	}

	return nil, errors.New("partido no encontrado")
}

// FUncino para borrar partidos

//Invoke-RestMethod -Uri "http://localhost:8080/api/matches/1" -Method DELETE

func deleteMatch(c *gin.Context) {
	id := c.Param("id")
	match, err := getMatchID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	//
	for i, m := range matches {
		if m.ID == id {
			matches = append(matches[:i], matches[i+1:]...)
			break
		}
	}

	// Devolvemos el partido eliminado (opcional)
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Partido eliminado correctamente",
		"deleted": match,
	})
}

//Funcion para actualizar un partido existente.

//$body = @{
//    homeTeam  = "Barcelona"
//    awayTeam  = "Real Madrid"
//    matchDate = "2025-05-20"
//} | ConvertTo-Json

func updateMatch(c *gin.Context) {
	id := c.Param("id")
	existingMatch, err := getMatchID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	var updatedData struct {
		HomeTeam  string `json:"homeTeam" binding:"required"`
		AwayTeam  string `json:"awayTeam" binding:"required"`
		MatchDate string `json:"matchDate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Debes enviar homeTeam, awayTeam y matchDate",
			"error":   err.Error(),
		})
		return
	}

	// actualizamos los datos
	existingMatch.HomeTeam = updatedData.HomeTeam
	existingMatch.AwayTeam = updatedData.AwayTeam
	existingMatch.MatchDate = updatedData.MatchDate

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Partido actualizado correctamente",
		"match":   existingMatch,
	})
}

// Funcion principal
func main() {
	router := gin.Default()

	// Middleware CORS (permite solicitudes desde el frontend)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := router.Group("/api")
	{
		api.GET("/matches", getMatch)
		api.POST("/matches", createMatch)
		api.GET("/matches/:id", matchById)
		api.DELETE("/matches/:id", deleteMatch)
		api.PUT("/matches/:id", updateMatch)
	}

	router.Run("localhost:8080")
}
