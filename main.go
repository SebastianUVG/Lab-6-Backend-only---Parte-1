package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Match struct {
	ID        int       `json:"id"`
	HomeTeam  string    `json:"homeTeam"`
	AwayTeam  string    `json:"awayTeam"`
	MatchDate time.Time `json:"matchDate"`
}

var db *pgx.Conn

// Funcion para obtener todos los partidos
func getMatch(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Query(ctx, "SELECT id, home_team, away_team, match_date FROM matches")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		matches = append(matches, m)
	}

	c.IndentedJSON(http.StatusOK, matches)
}

// Funcion para crear un nuevo partido

//Invoke-RestMethod -Uri http://localhost:8080/api/matches `
//-Method POST `
//-Headers @{"Content-Type"="application/json"} `
//-Body (Get-Content -Raw -Path .\body.json)

func createMatch(c *gin.Context) {
	var newMatch struct {
		HomeTeam  string `json:"homeTeam" binding:"required"`
		AwayTeam  string `json:"awayTeam" binding:"required"`
		MatchDate string `json:"matchDate" binding:"required"`
	}

	if err := c.BindJSON(&newMatch); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos", "error": err.Error()})
		return
	}

	// Parsear la fecha desde string
	parsedDate, err := time.Parse("2006-01-02", newMatch.MatchDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}

	ctx := context.Background()
	var id int
	err = db.QueryRow(ctx,
		"INSERT INTO matches (home_team, away_team, match_date) VALUES ($1, $2, $3) RETURNING id",
		newMatch.HomeTeam, newMatch.AwayTeam, parsedDate,
	).Scan(&id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"id":        id,
		"homeTeam":  newMatch.HomeTeam,
		"awayTeam":  newMatch.AwayTeam,
		"matchDate": newMatch.MatchDate,
	})
}

// funcion para obtener un partido mediante el uso de la iD
func matchById(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	var match Match
	ctx := context.Background()
	err = db.QueryRow(ctx,
		"SELECT id, home_team, away_team, match_date FROM matches WHERE id = $1",
		matchID,
	).Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, match)
}

// FUncino para borrar partidos

//Invoke-RestMethod -Uri "http://localhost:8080/api/matches/1" -Method DELETE

func deleteMatch(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	ctx := context.Background()
	result, err := db.Exec(ctx, "DELETE FROM matches WHERE id = $1", matchID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Partido eliminado correctamente"})
}

//Funcion para actualizar un partido existente.

//$body = @{
//    homeTeam  = "Barcelona"
//    awayTeam  = "Real Madrid"
//    matchDate = "2025-05-20"
//} | ConvertTo-Json

func updateMatch(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
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

	// Parsear la fecha desde string
	parsedDate, err := time.Parse("2006-01-02", updatedData.MatchDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}

	ctx := context.Background()
	result, err := db.Exec(ctx,
		"UPDATE matches SET home_team = $1, away_team = $2, match_date = $3 WHERE id = $4",
		updatedData.HomeTeam, updatedData.AwayTeam, parsedDate, matchID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Partido actualizado correctamente",
	})
}

func initDB() error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var err error
	var conn *pgx.Conn

	// Intenta conectarse hasta 5 veces con espera exponencial
	for i := 0; i < 5; i++ {
		ctx := context.Background()
		conn, err = pgx.Connect(ctx, connStr)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i*i) * time.Second) // Espera 0, 1, 4, 9, 16 segundos
	}

	if err != nil {
		return fmt.Errorf("no se pudo conectar a la base de datos después de 5 intentos: %v", err)
	}

	// Verifica la conexión
	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("no se pudo hacer ping a la base de datos: %v", err)
	}

	db = conn
	return nil
}

// Funcion principal
func main() {
	// Inicializar la conexión a la base de datos
	if err := initDB(); err != nil {
		fmt.Printf("Error inicializando la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	router := gin.Default()

	// Middleware CORS mejorado
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

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

	router.Run("0.0.0.0:8080")
}

// docker-compose down -v
// docker system prune -f
// docker-compose build --no-cache
// docker-compose up
