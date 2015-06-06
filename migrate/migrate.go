package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/seven5/seven5/migrate"
)

var defn = migrate.Definitions{
	Up: map[int]migrate.MigrationFunc{
		1: oneUp,
	},
	Down: map[int]migrate.MigrationFunc{
		1: oneDown,
	},
}

func main() {

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Fprintf(os.Stderr, "failed to get DATABASE_URL from envirnoment")
	}
	m, err := migrate.NewPostgresMigrator(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to selma: %v", err)
		return
	}
	migrate.Main(&defn, m)
}

func oneUp(tx *sql.Tx) error {
	//
	// USER
	//
	_, err := tx.Exec(`
        CREATE TABLE user_record (
        user_udid CHAR(36) PRIMARY KEY,
        first_name VARCHAR(63),
        last_name VARCHAR(63),
        email_addr VARCHAR(63),
        password VARCHAR(63),
        disabled BOOLEAN DEFAULT false,
        admin BOOLEAN DEFAULT false,
        manager BOOLEAN DEFAULT false
        )`)
	if err != nil {
		return err
	}
	//
	// MOVIE
	//
	_, err = tx.Exec(`
        CREATE TABLE movie (
       	id BIGINT PRIMARY KEY,
        title VARCHAR(255),
        imdb_url VARCHAR(1023),
        poster_url VARCHAR(1023),
        blurb TEXT,
        nominated_by VARCHAR(36) REFERENCES user_record(user_udid)
        )`)
	if err != nil {
		return err
	}
	//
	// LOVE
	//
	_, err = tx.Exec(`
        CREATE TABLE love (
       	id BIGINT PRIMARY KEY,
       	movie_id BIGINT REFERENCES movie(id),
        user_udid VARCHAR(36) REFERENCES user_record(user_udid)
        )`)
	if err != nil {
		return err
	}
	//
	// HATE
	//
	_, err = tx.Exec(`
        CREATE TABLE hate (
       	id BIGINT PRIMARY KEY,
       	movie_id BIGINT REFERENCES movie(id),
        user_udid VARCHAR(36) REFERENCES user_record(user_udid)
        )`)
	if err != nil {
		return err
	}
	//
	// COMMENT
	//

	_, err = tx.Exec(`
        CREATE TABLE comment (
       	id BIGINT PRIMARY KEY,
       	movie_id BIGINT REFERENCES movie(id),
        user_udid VARCHAR(36) REFERENCES user_record(user_udid),
        comment TEXT,
        updated TIMESTAMP WITH TIME ZONE
        )`)
	if err != nil {
		return err
	}

	return nil
}

func oneDown(tx *sql.Tx) error {
	//bc of foreign keys, order of these drops is siginficant
	drops := []string{
		"DROP TABLE love",
		"DROP TABLE hate",
		"DROP TABLE comment",
		"DROP TABLE movie",
		"DROP TABLE user_record",
	}
	for _, drop := range drops {
		_, err := tx.Exec(drop)
		if err != nil {
			return err
		}
	}
	return nil
}
