package main

import (
	"database/sql"
	"embed"
	"financial-service/data"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "82"

var counts int64

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

// main starts the server and establishing connection to database
func main() {
	log.Println("Starting financial service")
	err := godotenv.Load("example.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		Client: &http.Client{},
	}
	app.setupRepo(conn)

	goose.SetBaseFS(embedMigrations)

	migrationsDir := os.Getenv("GOOSE_MIGRATION_DIR")

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, migrationsDir); err != nil {
		panic(err)
	}

	// Инициализация Gin
	router := gin.Default()

	// Настройка маршрутов
	app.routes(router)

	// Запуск сервера
	log.Printf("Starting server on port %s\n", webPort)
	if err := router.Run(fmt.Sprintf(":%s", webPort)); err != nil {
		log.Panic(err)
	}
}

// openDB establishes a connection to the PostgreSQL database using the provided Data Source Name (DSN)
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx/v4", dsn)
	fmt.Println(db)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB connect to Postgres with provided dsn
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

// setupRepo sets new postgres repository
func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
