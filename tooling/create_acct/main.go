package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/lib/pq"
	s5 "github.com/seven5/seven5"
)

var (
	admin, mgr, disabled bool
)

func init() {
	flag.BoolVar(&admin, "admin", false, "set the admin flag on this account")
	flag.BoolVar(&mgr, "mgr", false, "set the mgr flag on this account")
	flag.BoolVar(&disabled, "disabled", false, "set the disabled flag on this account")
}

func main() {

	flag.Parse()
	if flag.NArg() != 4 {
		log.Fatalf("required args: first_name last_name email password")
	}
	first := flag.Arg(0)
	last := flag.Arg(1)
	email := flag.Arg(2)
	password := flag.Arg(3)

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatalf("failed to get DATABASE_URL from envirnoment")
	}
	opts, err := pq.ParseURL(url)
	if err != nil {
		log.Fatalf("failed to parse DATABASE_URL: %v", err)
	}

	db, err := sql.Open("postgres", opts)
	if err != nil {
		log.Fatalf("failed to open database (%s): %v", url, err)
	}

	_, err = db.Exec("INSERT INTO user_record "+
		"(user_udid, first_name, last_name, email_addr,password,disabled,admin,manager) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8)",
		s5.UDID(), first, last, email, password, disabled, admin, mgr)
	if err != nil {
		log.Fatalf("failed on INSERT: %v", err)
	}
}
