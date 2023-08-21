package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage interface {
	init(string)
	getLongerURL(string) (string, error)
	storeShortenedURL(string, ...string)
	removeShortenedURL(string)
}

type sqlstorage struct {
	db *sql.DB
}

type postgresstorage sqlstorage

func (s *postgresstorage) init(dataSourceName string) {
	if dataSourceName == "" {
		dataSourceName = "postgres://postgres:postgres@postgres/postgres?sslmode=disable"
	}

	db, _ := sql.Open("postgres", dataSourceName)
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping SQL data source '%s':\n%s", dataSourceName, err)
	} else {
		log.Printf("Successfully connected to SQL data source '%s'\n", dataSourceName)
		s.db = db
	}

	m, err := migrate.New("file://migrations/", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	m.Up()

}

func (s postgresstorage) getLongerURL(shortenedURL string) (string, error) {
	var longUrlId string
	err := s.db.QueryRow(fmt.Sprintf("SELECT long_url_id FROM shortened_urls WHERE short_url = '%s'", shortenedURL)).Scan(&longUrlId)
	if err != nil {
		return longUrlId, err
	}
	var longUrl string
	err = s.db.QueryRow(fmt.Sprintf("SELECT long_url FROM long_urls WHERE id = %s", longUrlId)).Scan(&longUrl)
	return longUrl, err

}

func (s postgresstorage) storeShortenedURL(longURL string, shortenedUrls ...string) {
	// Check if long_url already exists and create it otherwise
	var longUrlId int64
	fmt.Println("Checking if long_url already exists")
	err := s.db.QueryRow(fmt.Sprintf("SELECT id FROM long_urls WHERE long_url = '%s'", longURL)).Scan(&longUrlId)
	if err != nil {
		fmt.Println("long_url does not exist yet, creating it")
		res, err := s.db.Exec(fmt.Sprintf("INSERT INTO long_urls (long_url) VALUES ('%s')", longURL))
		if err == nil {
			longUrlId, _ = res.LastInsertId()
		} else {
			log.Println(err)
		}
	} else {
		fmt.Println("long_url already exists, and is ", longUrlId)
	}
	for _, shortUrl := range shortenedUrls {
		s.db.Exec(fmt.Sprintf("INSERT INTO shortened_urls (short_url, long_url_id) VALUES ('%s',%d)", shortUrl, longUrlId))
	}
}

func (s postgresstorage) removeShortenedURL(shortenedURL string) {
	_, err := s.db.Exec(fmt.Sprintf("DELETE FROM shortened_urls WHERE short_url = '%s'", shortenedURL))
	if err != nil {
		log.Println(err)
	}
}
