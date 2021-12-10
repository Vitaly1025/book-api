package main

import (
	server "book-api"
	handlers "book-api/handler"
	"book-api/middleware"
	"book-api/repository"
	"book-api/service"
	"io/ioutil"
	"log"
	"os"

	_ "book-api/docs"

	"github.com/joho/godotenv"
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
		Host:     cfg.Db.Host,
		Port:     cfg.Db.Port,
		DBName:   cfg.Db.Name,
		Username: cfg.Db.Username,
		Password: cfg.Db.Password,
		SSLMode:  cfg.Db.SSL,
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

func OSReadDir(root string)  {
    var files []string
    f, err := os.Open(root)
    if err != nil {
        log.Println(err)
    }
    fileInfo, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Println(err)
    }

    for _, file := range fileInfo {
		log.Println(files, file.Name())
    }
    
}

func readFile(cfg *Config) {
	files, err := ioutil.ReadDir(".")
    if err != nil {
        log.Fatal(err)
    }
 
    for _, f := range files {
            log.Println(f.Name())
    }

	log.Println(os.Getwd())
	OSReadDir(".")

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
	if err := godotenv.Load("./configs/config.env"); err != nil {
		log.Print("No .env file found")
	}
	var cfg Config
	// readFile(&cfg)
	readEnv(&cfg)
	return cfg
}

type Config struct {
	Server struct {
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Db struct {
		Host     string `yaml:"host", envconfig:"DB_HOST"`
		Port     string `yaml:"port", envconfig:"DB_PORT"`
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		Name   string `yaml:"dbname", envconfig:"DB_NAME"`
		Password string `yaml:"password", envconfig:"DB_PASSWORD"`
		SSL      string `yaml:"ssl", envconfig:"DB_SSL"`
	} `yaml:"db"`
}
