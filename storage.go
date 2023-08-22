package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Almost like a ternary operator
func If(condition bool, trueVal any, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

type Storage interface {
	init(string)
	getLongerURL(string) (string, error)
	storeShortenedURL(string, bool, ...string)
	removeShortenedURL(...string)
	getAllShortenedURLs(bool) (map[string][2]string, error)
	// Logs a request to shorten a URL.
	// The first argument is the request object.
	//
	// The second argument is the shortened URL. Can be empty.
	logShorteningRequest(*http.Request, string)
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

func (s postgresstorage) storeShortenedURL(longURL string, public bool, shortenedUrls ...string) {
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
		s.db.Exec(fmt.Sprintf("INSERT INTO shortened_urls (short_url, long_url_id, private) VALUES ('%s',%d,'%v')", shortUrl, longUrlId, !public))
	}
}

func (s postgresstorage) removeShortenedURL(shortenedURLs ...string) {
	conditions := ""
	for i, shortUrl := range shortenedURLs {
		if i != 0 {
			conditions += " OR "
		}
		conditions += fmt.Sprintf("short_url = '%s'", shortUrl)
	}
	_, err := s.db.Exec(fmt.Sprintf("DELETE FROM shortened_urls WHERE %s", conditions))
	if err != nil {
		log.Println(err)
	}
}

func (s postgresstorage) getAllShortenedURLs(publiconly bool) (map[string][2]string, error) {
	rows, err := s.db.Query(fmt.Sprintf("SELECT short_url, long_url, private FROM shortened_urls INNER JOIN long_urls ON shortened_urls.long_url_id = long_urls.id%s", If(publiconly, " WHERE private=FALSE", "")))
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	shortenedURLs := map[string][2]string{}
	for rows.Next() {
		var shortURL string
		var longURL string
		var private bool
		err = rows.Scan(&shortURL, &longURL, &private)
		if err != nil {
			return nil, err
		}
		shortenedURLs[shortURL] = [2]string{longURL, fmt.Sprintf("%v", private)}
	}
	return shortenedURLs, nil
}

func (s postgresstorage) logShorteningRequest(r *http.Request, shortenedURL string) {

}
