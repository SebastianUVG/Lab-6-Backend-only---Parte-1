// @title LaLigaTracker API
// @version 1.0
// @description API para gestión de partidos de fútbol
// @contact.name Soporte API
// @contact.email soporte@sebastian_laliga.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Match define la estructura de un partido de fútbol
// @Description Información completa sobre un partido de fútbol
type Match struct {
	ID          int       `json:"id"`
	HomeTeam    string    `json:"homeTeam"`
	AwayTeam    string    `json:"awayTeam"`
	MatchDate   time.Time `json:"matchDate"`
	Goals       int       `json:"goals,omitempty"` // Cambiado de GoalsHome/GoalsAway
	YellowCards int       `json:"yellowCards,omitempty"`
	RedCards    int       `json:"redCards,omitempty"`
	ExtraTime   int       `json:"extraTime,omitempty"`
}

var db *pgx.Conn

// getMatch godoc
// @Summary Obtener todos los partidos
// @Description Retorna una lista de todos los partidos registrados
// @Tags matches
// @Accept json
// @Produce json
// @Success 200 {array} Match
// @Failure 500 {object} map[string]string
// @Router /matches [get]
func getMatch(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Query(ctx, `
        SELECT id, home_team, away_team, match_date, 
            yellow_cards, red_cards, extra_time 
        FROM matches`)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate,
			&m.YellowCards, &m.RedCards, &m.ExtraTime)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		matches = append(matches, m)
	}

	c.IndentedJSON(http.StatusOK, matches)
}

// createMatch godoc
// @Summary Crear un nuevo partido
// @Description Crea un nuevo partido con los datos proporcionados
// @Tags matches
// @Accept json
// @Produce json
// @Param match body object{homeTeam=string,awayTeam=string,matchDate=string} true "Datos del partido"
// @Success 201 {object} Match
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches [post]
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

// matchById godoc
// @Summary Obtener un partido por ID
// @Description Retorna un partido específico según su ID
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} Match
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id} [get]
func matchById(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	var match Match
	ctx := context.Background()
	err = db.QueryRow(ctx, `
        SELECT id, home_team, away_team, match_date, 
            yellow_cards, red_cards, extra_time 
        FROM matches WHERE id = $1`,
		matchID,
	).Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.MatchDate,
		&match.YellowCards, &match.RedCards, &match.ExtraTime)

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

// deleteMatch godoc
// @Summary Eliminar un partido
// @Description Elimina un partido según su ID
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id} [delete]
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

// updateMatch godoc
// @Summary Actualizar un partido
// @Description Actualiza todos los campos de un partido existente
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Param match body object{homeTeam=string,awayTeam=string,matchDate=string} true "Datos actualizados del partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id} [put]
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

// @Summary Registrar un gol
// @Description Incrementa el contador general de goles
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id}/goals [patch]
func registerGoal(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	ctx := context.Background()
	result, err := db.Exec(ctx,
		"UPDATE matches SET goals = goals + 1 WHERE id = $1",
		matchID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Gol registrado correctamente"})
}

// registerYellowCard godoc
// @Summary Registrar tarjeta amarilla
// @Description Incrementa el contador de tarjetas amarillas del partido
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id}/yellowcards [patch]
func registerYellowCard(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	ctx := context.Background()
	result, err := db.Exec(ctx,
		"UPDATE matches SET yellow_cards = yellow_cards + 1 WHERE id = $1",
		matchID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Tarjeta amarilla registrada"})
}

// registerRedCard godoc
// @Summary Registrar tarjeta roja
// @Description Incrementa el contador de tarjetas rojas del partido
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id}/redcards [patch]
func registerRedCard(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	ctx := context.Background()
	result, err := db.Exec(ctx,
		"UPDATE matches SET red_cards = red_cards + 1 WHERE id = $1",
		matchID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Tarjeta roja registrada"})
}

// setExtraTime godoc
// @Summary Incrementar tiempo extra
// @Description Incrementa en 1 minuto el tiempo extra del partido (hasta máximo 30 minutos)
// @Tags matches
// @Accept json
// @Produce json
// @Param id path int true "ID del Partido"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /matches/{id}/extratime [patch]
func setExtraTime(c *gin.Context) {
	id := c.Param("id")
	matchID, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID debe ser un número"})
		return
	}

	ctx := context.Background()

	// Primero obtenemos el valor actual del tiempo extra
	var currentExtraTime int
	err = db.QueryRow(ctx, "SELECT extra_time FROM matches WHERE id = $1", matchID).Scan(&currentExtraTime)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Calculamos el nuevo valor (incrementamos 1 minuto, máximo 30)
	newExtraTime := currentExtraTime + 1
	if newExtraTime > 30 {
		newExtraTime = 30
	}

	// Actualizamos en la base de datos
	result, err := db.Exec(ctx,
		"UPDATE matches SET extra_time = $1 WHERE id = $2",
		newExtraTime, matchID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partido no encontrado"})
		return
	}

	// Mensaje de respuesta
	message := fmt.Sprintf("Tiempo extra incrementado a %d minutos", newExtraTime)
	if newExtraTime >= 30 {
		message = "Tiempo extra alcanzó el máximo de 30 minutos"
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":   message,
		"extraTime": newExtraTime,
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

	for i := 0; i < 5; i++ {
		ctx := context.Background()
		conn, err = pgx.Connect(ctx, connStr)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i*i) * time.Second)
	}

	if err != nil {
		return fmt.Errorf("no se pudo conectar a la base de datos después de 5 intentos: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("no se pudo hacer ping a la base de datos: %v", err)
	}

	db = conn
	return nil
}

// @title LaLigaTracker API
// @version 1.0
// @description API para gestión de partidos de fútbol
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soporte@laliga.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api
func main() {
	if err := initDB(); err != nil {
		fmt.Printf("Error inicializando la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Configuración de Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.GET("/matches", getMatch)
		api.POST("/matches", createMatch)
		api.GET("/matches/:id", matchById)
		api.DELETE("/matches/:id", deleteMatch)
		api.PUT("/matches/:id", updateMatch)

		api.PATCH("/matches/:id/goals", registerGoal)
		api.PATCH("/matches/:id/yellowcards", registerYellowCard)
		api.PATCH("/matches/:id/redcards", registerRedCard)
		api.PATCH("/matches/:id/extratime", setExtraTime)
	}

	router.Run("0.0.0.0:8080")
}

// docker-compose down -v
// docker system prune -f
// docker-compose build --no-cache
// docker-compose up
