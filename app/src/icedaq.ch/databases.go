package main

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Artist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Show struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	StartTime  time.Time `json:"starttime"`
	EndTime    time.Time `json:"endtime"`
	Price      float64   `json:"price"`
	MaxTickets int64     `json:"maxtickets"`
	ArtistId   int64     `json:"artistid"`
}

type mydatabase interface {
	create()
	insert(class interface{})
	get(class string) []interface{}
	buy(showId int64) bool
	ticketsSold() int64
}

type mysql struct {
	connectionString string
}

type myspanner struct {
	connectionString string
	ctx              context.Context
	client           *spanner.Client
	adminClient      *database.DatabaseAdminClient
}

func (m mysql) create() {

	name := "ticketshop"

	db, err := sql.Open("mysql", m.connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + name)
	if err != nil {
		panic(err)
	}

	// Now create the tables.
	tableArtists := `
			   create table Artists (
				ArtistId int not null auto_increment,
				Name varchar(1024) not null,
				primary key (ArtistId))`
	tableShows := `
			 create table Shows (
				ShowId int not null auto_increment,
				Name varchar(1024) not null,
				StartTime timestamp not null default CURRENT_TIMESTAMP,
				EndTime timestamp not null default CURRENT_TIMESTAMP,
				Price float not null,
				MaxTickets int not null,
				ArtistId int not null,
				primary key (ShowId))`
	tableTickets := `
			create table Tickets (
				ShowId int not null,
				UserId int not null,
				primary key (ShowId, UserId))`
	tableUsers := `
			create table Users (
				UserId int not null auto_increment,
				 EMail varchar(1024) not null,
				 Password varchar(1024) not null,
				 primary key (UserId))`

	_, err = db.Exec(tableArtists)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(tableShows)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(tableTickets)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(tableUsers)
	if err != nil {
		panic(err)
	}

}

func (m mysql) insert(obj interface{}) {
	db, err := sql.Open("mysql", m.connectionString+"ticketshop")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if c, ok := obj.(Artist); ok { // For the Artists
		stmtIns, err := db.Prepare("INSERT INTO Artists VALUES( ?, ? )") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		stmtIns.Exec(c.Id, c.Name)
	}

	if c, ok := obj.(Show); ok { // For the Shows
		stmtIns, err := db.Prepare("INSERT INTO Shows VALUES( ?, ?, ?, ?, ?, ?, ?)") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		stmtIns.Exec(c.Id, c.Name, c.StartTime, c.EndTime, c.Price, c.MaxTickets, c.ArtistId)
	}
}

func (m mysql) get(class string) []interface{} {

	var result []interface{}

	db, err := sql.Open("mysql", m.connectionString+"ticketshop?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	switch class {
	case "artists":
		rows, err := db.Query("SELECT * FROM Artists")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		for rows.Next() {
			a := new(Artist)
			err := rows.Scan(&a.Id, &a.Name)
			if err != nil {
				panic(err.Error())
			}
			result = append(result, a)
		}
		rows.Close()
	case "shows":
		rows, err := db.Query("SELECT * FROM Shows")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		for rows.Next() {
			s := new(Show)
			err := rows.Scan(&s.Id, &s.Name, &s.StartTime, &s.EndTime, &s.Price, &s.MaxTickets, &s.ArtistId)
			if err != nil {
				panic(err.Error())
			}
			result = append(result, s)
		}
		rows.Close()
	}
	return result
}

func (m mysql) buy(showId int64) bool {

	// Buy a ticket for a certain show for a random user.

	db, err := sql.Open("mysql", m.connectionString+"ticketshop?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return false
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		err = tx.Commit()
	}()

	var maxNullId sql.NullInt64
	var maxId int64
	err = db.QueryRow("SELECT MAX(UserId) FROM Tickets").Scan(&maxNullId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if maxNullId.Valid {
		maxId = maxNullId.Int64 + 1
	} else {
		maxId = 1
	}

	// count sold tickets for show
	var ticketCount int64
	err = db.QueryRow("SELECT COUNT(*) FROM Tickets WHERE ShowId = ?", showId).Scan(&ticketCount)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// compare to max.
	var maxTickets int64
	err = db.QueryRow("SELECT MaxTickets FROM Shows WHERE ShowId = ?", showId).Scan(&maxTickets)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if ticketCount == maxTickets {
		return false
	}

	// buy ticket
	stmtIns, err := db.Prepare("INSERT INTO Tickets VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		return false
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	if _, err = stmtIns.Exec(showId, maxId); err != nil {
		return false
	}

	return true
}

func (m mysql) ticketsSold() int64 {
	db, err := sql.Open("mysql", m.connectionString+"ticketshop?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ticketCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM Tickets").Scan(&ticketCount)
	if err != nil {
		fmt.Println(err)
	}

	return int64(ticketCount)
}

func NewSpanner(connectionString string) myspanner {
	thespanner := myspanner{}

	ctx := context.Background()
	adminClient, dataClient := createSpannerClients(ctx, connectionString)

	thespanner.connectionString = connectionString
	thespanner.client = dataClient
	thespanner.adminClient = adminClient
	thespanner.ctx = ctx

	return thespanner
}

func (s myspanner) create() {
	createSpannerDatabase(s.ctx, s.adminClient, s.connectionString)
}

func (s myspanner) insert(obj interface{}) {

	if c, ok := obj.(Artist); ok { // For the Artists
		artistColumns := []string{"ArtistId", "ArtistHid", "Name"}
		artistHid := shaStringFromInt(c.Id)
		m := []*spanner.Mutation{
			spanner.InsertOrUpdate("Artists", artistColumns, []interface{}{c.Id, artistHid, c.Name}),
		}

		_, err := s.client.Apply(s.ctx, m)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	if c, ok := obj.(Show); ok { // For the Shows
		showColumns := []string{"ShowId", "ShowHid", "Name", "StartTime", "EndTime", "Price", "MaxTickets", "ArtistId"}
		showHid := shaStringFromInt(c.Id)
		m := []*spanner.Mutation{
			spanner.InsertOrUpdate("Shows", showColumns, []interface{}{c.Id, showHid, c.Name, c.StartTime, c.EndTime, c.Price, c.MaxTickets, c.ArtistId}),
		}

		_, err := s.client.Apply(s.ctx, m)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func (s myspanner) get(class string) []interface{} {

	var result []interface{}

	switch class {
	case "artists":
		stmt := spanner.Statement{SQL: `SELECT ArtistId, Name FROM Artists`}
		iter := s.client.Single().Query(s.ctx, stmt)
		defer iter.Stop()
		row, err := iter.Next()
		for err != iterator.Done {
			if err != nil {
				panic(err.Error())
			}
			a := new(Artist)
			if err := row.Columns(&a.Id, &a.Name); err != nil {
				panic(err.Error())
			}
			result = append(result, a)
			row, err = iter.Next()
			if err != nil {
				break
			}
		}
	case "shows":
		stmt := spanner.Statement{SQL: `SELECT ShowId, Name, StartTime, EndTime, Price, MaxTickets, ArtistId FROM Shows`}
		iter := s.client.Single().Query(s.ctx, stmt)
		defer iter.Stop()
		row, err := iter.Next()
		for err != iterator.Done {
			if err != nil {
				panic(err.Error())
			}
			sh := new(Show)
			if err := row.Columns(&sh.Id, &sh.Name, &sh.StartTime, &sh.EndTime, &sh.Price, &sh.MaxTickets, &sh.ArtistId); err != nil {
				panic(err.Error())
			}

			result = append(result, sh)
			row, err = iter.Next()
			if err != nil {
				break
			}
		}
	}

	return result
}

func (s myspanner) buy(showId int64) bool {

	// Since we do not have users, we just hash a random number and use this as the random user. This might fail but YOLO.

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	userHid := shaStringFromInt(int64(r1.Int()))
	showHid := shaStringFromInt(showId)

	// Buy a ticket for a certain show for a random user.
	_, err := s.client.ReadWriteTransaction(s.ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		// count sold tickets for show
		var ticketCount int64

		stmt := spanner.Statement{
			SQL: `SELECT COUNT(*) FROM Tickets WHERE ShowHid = @showid`,
			Params: map[string]interface{}{
				"showid": showHid,
			},
		}
		iter := s.client.Single().Query(s.ctx, stmt)
		defer iter.Stop()
		row, err := iter.Next()
		if err == iterator.Done {
			return err
		}
		if err != nil {
			return err
		}
		if err := row.Columns(&ticketCount); err != nil {
			return err
		}

		// check max tickets for a show.
		// Problem:
		row, err = txn.ReadRow(ctx, "Shows", spanner.Key{showHid}, []string{"MaxTickets", "ArtistId"})
		if err != nil {
			return err
		}
		var maxTickets int64
		var artistId int64
		err = row.Column(0, &maxTickets)
		err = row.Column(0, &artistId)
		if err != nil {
			return err
		}

		if ticketCount == maxTickets {
			return errors.New("max tickets reached for show: " + strconv.FormatInt(showId, 10))
		}

		artistHid := shaStringFromInt(artistId)
		cols := []string{"ArtistHid", "ShowHid", "UserHid"}
		txn.BufferWrite([]*spanner.Mutation{
			spanner.InsertOrUpdate("Tickets", cols, []interface{}{artistHid, showHid, userHid}),
		})

		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

func createSpannerClients(ctx context.Context, db string) (*database.DatabaseAdminClient, *spanner.Client) {

	var adminClient *database.DatabaseAdminClient
	var dataClient *spanner.Client

	if _, err := os.Stat("/var/run/secret/cloud.google.com/service-account.json"); !os.IsNotExist(err) {
		adminClient, err = database.NewDatabaseAdminClient(ctx, option.WithServiceAccountFile("/var/run/secret/cloud.google.com/service-account.json"))
		if err != nil {
			log.Fatal(err)
		}

		dataClient, err = spanner.NewClient(ctx, db, option.WithServiceAccountFile("/var/run/secret/cloud.google.com/service-account.json"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		adminClient, err = database.NewDatabaseAdminClient(ctx)
		if err != nil {
			log.Fatal(err)
		}

		dataClient, err = spanner.NewClient(ctx, db)
		if err != nil {
			log.Fatal(err)
		}
	}

	return adminClient, dataClient
}

func createSpannerDatabase(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`create table Artists (
				ArtistId int64 not null,
				ArtistHid string(128) not null,
				Name string(1024) not null
			) primary key (ArtistHid)`,
			`create table Shows (
				ShowId int64 not null,
				ShowHid string(128) not null,				
				Name string(1024) not null,
				StartTime timestamp not null,
				EndTime timestamp not null,
				Price float64 not null,
				MaxTickets int64 not null,
				ArtistId int64 not null
			) primary key (ShowHid)`,
			`create table Tickets (
				ArtistHid string(128) not null,
				ShowHid string(128) not null,
				UserHid string(128) not null
			) primary key (ArtistHid, ShowHid, UserHid)`,
			`create table Users (
				UserId int64 not null,
				UserHid string(128) not null,
				EMail string(1024) not null,
				Password string(1024) not null
			) primary key (UserHid)`,
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Created database " + db)
	return nil
}

func (s myspanner) ticketsSold() int64 {

	var ticketCount int64
	ticketCount = 0
	stmt := spanner.Statement{SQL: `SELECT COUNT(*) FROM Tickets`}
	iter := s.client.Single().Query(s.ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == iterator.Done {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := row.Columns(&ticketCount); err != nil {
		log.Fatal(err)
	}

	return int64(ticketCount)
}

func shaStringFromInt(number int64) string {

	numberString := strconv.FormatInt(number, 10)
	hash := sha256.Sum256([]byte(numberString))
	return fmt.Sprintf("%x\n", hash)

}
