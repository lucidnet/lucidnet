package vertex

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lucidnet/lucidnet/internal/app/vertex/admin"
	"github.com/lucidnet/lucidnet/internal/app/vertex/config"
	"github.com/lucidnet/lucidnet/internal/app/vertex/health"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
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
	err := os.MkdirAll(s.config.DataPath, os.ModePerm)

	if err != nil {
		log.Fatalln("failed to instantiate data directory: ", err)
	}

	db, err := gorm.Open(sqlite.Open(path.Join(s.config.DataPath, "vertex.db")), &gorm.Config{})

	if err != nil {
		log.Fatalln("failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&config.Config{}, &admin.Admin{})

	if err != nil {
		log.Fatalln("failed to migrate database models: ", err)
	}

	router := gin.Default()

	configManager := config.NewManager(db)
	err = configManager.Init()

	if err != nil {
		log.Fatalln("failed to initialize config manager: ", err)
	}

	health.NewHandler(router).Register()

	log.Println("starting vertex server on host and port: ", s.config.Host, s.config.Port)

	err = router.Run(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))

	if err != nil {
		log.Fatalln(err)
	}
}
