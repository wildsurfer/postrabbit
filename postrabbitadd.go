package main

import (
	"database/sql"
	"log"
)

// Config contains various config data populated from YAML

func add(config Config) {

	db, err := sql.Open("postgres", config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO urls(url) VALUES($1)`, urlarg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
