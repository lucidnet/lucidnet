package vertex

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

type ServerConfig struct {
	Host     string
	Port     string
	DataPath string
}

func (s *Server) Start() {
	router := gin.Default()

	log.Println("starting vertex server on host and port: ", s.config.Host, s.config.Port)

	err := router.Run(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))

	if err != nil {
		log.Fatalln(err)
	}
}
