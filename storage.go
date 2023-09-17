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

// Almost like a ternary operator
func If(condition bool, trueVal any, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

type Storage interface {
	// Initializes the storage model
	init(dataSourceName string)
	// Returns the long URL associated with the provided shortened URL
	getLongerURL(shortenedURL string) (string, error)
	// Stores the long URL in the database, and maps it to the provided shortened URLs
	storeShortenedURL(longURL string, hidden bool, shortenedUrls ...string)
	// Removes the provided shortened URLs from the database
	removeShortenedURL(shortenedURLs ...string)
	// Returns all shortened urls in the form of a map of shortURL -> [longURL, hidden]
	// If publiconly is true, only returns mappings that have 'hidden' set to false
	getAllShortenedURLs(publiconly bool) (map[string][2]string, error)
	// Logs a request to shorten a URL.
	logShorteningRequest(ip, shortURL, longUrl string)
}

type sqlstorage struct {
	db *sql.DB
}

type postgresstorage sqlstorage

func (s *postgresstorage) init(dataSourceName string) {
	if dataSourceName == "" {
		dataSourceName = "postgres://postgres:postgres@postgres-dev/postgres?sslmode=disable"
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
	var longUrl string
	err := s.db.QueryRow(fmt.Sprintf(`SELECT long_urls.long_url 
	FROM shortened_urls 
	INNER JOIN long_urls ON shortened_urls.long_url_id = long_urls.id 
	WHERE shortened_urls.short_url = '%s'`, shortenedURL)).Scan(&longUrl)
	return longUrl, err
}

func (s postgresstorage) storeShortenedURL(longURL string, hidden bool, shortenedUrls ...string) {
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
		s.db.Exec(fmt.Sprintf("INSERT INTO shortened_urls (short_url, long_url_id, hidden) VALUES ('%s',%d,'%v')", shortUrl, longUrlId, hidden))
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
	rows, err := s.db.Query(fmt.Sprintf("SELECT short_url, long_url, hidden FROM shortened_urls INNER JOIN long_urls ON shortened_urls.long_url_id = long_urls.id%s", If(publiconly, " WHERE hidden=FALSE", "")))
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	shortenedURLs := map[string][2]string{}
	for rows.Next() {
		var shortURL string
		var longURL string
		var hidden bool
		err = rows.Scan(&shortURL, &longURL, &hidden)
		if err != nil {
			return nil, err
		}
		shortenedURLs[shortURL] = [2]string{longURL, fmt.Sprintf("%v", hidden)}
	}
	return shortenedURLs, nil
}

func (s postgresstorage) logShorteningRequest(ip, shortURL, longUrl string) {
	s.db.Exec(fmt.Sprintf("INSERT INTO log (id, timestamp, ip, success, shortURL, redirectURL) VALUES (DEFAULT, DEFAULT, '%s', %v, '%s', '%s')", ip, longUrl != "", shortURL, longUrl))
}
