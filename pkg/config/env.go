package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	//server
	PublicHost string
	Port       string

	//baza
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	//apikeys
	OMDbAPIKey string
	TMDbAPIKey string
	AuthAPIKey string

	MovieSourcesPath string
}

// kreiram Singleton da ne bih svaki put pozivala initConfig funkciju
// globalna promenljiva kojoj mogu da pristupim
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load("../../.env")

	return Config{
		PublicHost:       getEnv("PUBLIC_HOST", "http://localhost"),
		Port:             getEnv("PORT", "9010"),
		DBUser:           getEnv("DB_USER", "user"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBAddress:        fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3307")),
		DBName:           getEnv("DB_NAME", "movies"),
		OMDbAPIKey:       getEnv("OMDB_API_KEY", ""),
		TMDbAPIKey:       getEnv("TMDB_API_KEY", ""),
		AuthAPIKey:       getEnv("AUTH_API_KEY", "my_api_key"),
		MovieSourcesPath: getEnv("MOVIE_SOURCES_PATH", "../../movie_sources.json"),
	}
}

func getEnv(key, fallback string) string { //fallback u slucaju da vrednost kljuca ne postoji
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
