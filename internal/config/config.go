package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db         DataBase   `yaml:"db"`
	HttpServer HttpServer `yaml:"server"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user" env:"DB_USER"`
	Name     string `yaml:"dbname" env:"DB_NAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	SSL      string `yaml:"ssl"`
}

type HttpServer struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func MustLoad() Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path doesn't exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Sprintf("failed to read config: %s\n%s", path, err.Error()))
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return fmt.Sprintf("%s%s", getProjectRootPath(), res)
}

func getProjectRootPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic("cannot find executable")
	}

	absPath, err := filepath.Abs(ex)
	if err != nil {
		panic("cannot find absolute path")
	}
	return filepath.Dir(absPath)
}
