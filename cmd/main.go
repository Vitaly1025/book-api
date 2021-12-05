package main

import (
	server "book-api"
	handlers "book-api/pkg/handler"
	"book-api/pkg/middleware"
	"book-api/pkg/repository"
	"book-api/pkg/service"
	"log"
	"os"

	_ "book-api/docs"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// @title Book API
// @version 1.0
// @description This is a sample server.
// @host localhost:4000

func main() {
	srv := new(server.Server)
	cfg := initConfig()
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		DBName:   cfg.Database.DBName,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		SSLMode:  cfg.Database.SSL,
	})

	if err != nil {
		log.Fatalln("failed connection to database", err.Error())
	}

	mdw := middleware.NewMiddleware()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	hndlrs := handlers.NewHandler(services, mdw)

	if err := srv.Run(cfg.Server.Port, hndlrs.InitRoutes()); err != nil {
		log.Fatalln("Can't start server", err.Error())
	}
}

func readFile(cfg *Config) {
	f, err := os.Open("./configs/config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func initConfig() Config {
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	return cfg
}

type Config struct {
	Server struct {
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host", envconfig:"DB_HOST"`
		Port     string `yaml:"port", envconfig:"DB_PORT"`
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		DBName   string `yaml:"dbname", envconfig:"DB_NAME"`
		Password string `yaml:"password", envconfig:"DB_PASSWORD"`
		SSL      string `yaml:"ssl"`
	} `yaml:"db"`
}
