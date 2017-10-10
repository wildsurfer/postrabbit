package main

import (
	"database/sql"
	"log"
)

// Config contains various config data populated from YAML

func setup(config Config) {

	db, err := sql.Open("postgres", config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS urls;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE urls (id serial primary key, url varchar);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE OR REPLACE FUNCTION notify_trigger() RETURNS trigger AS $$
BEGIN
	PERFORM pg_notify('urlwork', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TRIGGER urlbefore BEFORE INSERT ON urls
    FOR EACH ROW EXECUTE PROCEDURE notify_trigger();`)
	if err != nil {
		log.Fatal(err)
	}
}
