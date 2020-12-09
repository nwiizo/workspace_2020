package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/Songmu/go-ltsv"
	"github.com/kr/pretty"
	_ "github.com/mattn/go-sqlite3"
)

type log struct {
	Time    *logTime
	Id      string
	Level   string
	Method  string
	Uri     string
	Reqtime float64
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

type logTime struct {
	time.Time
}

func (lt *logTime) UnmarshalText(t []byte) error {
	ti, err := time.ParseInLocation(timeFormat, string(t), time.UTC)
	if err != nil {
		return err
	}
	lt.Time = ti
	return nil
}

func main() {
	datePtr := flag.String("logfile", "", "a string")
	flag.Parse()
	// file name
	filename := "../logs/access-" + *datePtr + ".log"
	pretty.Println(filename)

	// open file
	fp, err := os.Open(filename)
	if err != nil {
	}
	// Close when finished
	defer fp.Close()

	// sqlite3
	// Setting Connection DB
	var DbConnection *sql.DB
	// Created if DB is not opened
	DbConnection, _ = sql.Open("sqlite3", "./logs.sql")
	// Close when finished
	defer DbConnection.Close()
	// DB creation SQL command
	// CREATE TABLE logs (Id text PRIMARY KEY,uri text,date integer);
	cmd := `CREATE TABLE logs (Id text PRIMARY KEY,uri text,date integer)`
	// Run Since the execution result is not returned, set it to _
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		pretty.Println(err)
	}
	// Process line by line
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		l := log{}
		ltsv.Unmarshal(scanner.Bytes(), &l)
		// INSERT INTO logs VALUES ("********-****-****-****-a*****f","/api/v1/messages",20171102);
		cmd = "INSERT INTO logs (id, uri, date) VALUES (?, ?, ?)"
		_, err = DbConnection.Exec(cmd, l.Id, l.Uri, *datePtr)
		if err != nil {
			pretty.Println(err)
		}
		pretty.Println(l.Id, l.Uri, *datePtr)
	}
	// SELECT URI,DATE,count(*) from logs group by URI,DATE;
	cmd = "SELECT URI,DATE,count(*) from logs group by URI,DATE;"
	rows, err := DbConnection.Query(cmd)
	defer rows.Close()
	if err != nil {
		pretty.Println(err)
	}
	for rows.Next() {
		var duri string
		var ddate string
		var dcount string
		if err := rows.Scan(&duri, &ddate, &dcount); err != nil {
			pretty.Println(err)
			return
		}
		pretty.Println(duri, ddate, dcount)
	}
}
