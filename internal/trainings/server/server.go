package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db     *sqlx.DB
	port   string
	router *gin.Engine
}

func NewServer(db *sqlx.DB, port string, router *gin.Engine) *Server {
	return &Server{
		db:     db,
		port:   port,
		router: router,
	}
}

func (s *Server) HealthCheck(c *gin.Context) {
	err := s.db.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database connection failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (s *Server) Start() {
	s.router.Run(":" + s.port)
}
