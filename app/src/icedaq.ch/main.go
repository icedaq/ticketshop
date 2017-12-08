package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// Read the type of the used database.
	var dbtype string
	value := os.Getenv("DBBACKEND")
	if len(value) == 0 {
		dbtype = "mysql"
	} else {
		dbtype = value
	}

	// Read the connection string
	var dbConString string
	value = os.Getenv("DBCONSTRING")
	if len(value) == 0 {
		if dbtype == "mysql" {
			dbConString = "root:notsosecretpassword@tcp(127.0.0.1:3306)/"
		} else {
			dbConString = "projects/default-1296/instances/ticketshop/databases/ticketshop"
		}
	} else {
		dbConString = value
	}

	var db mydatabase
	if dbtype == "mysql" {
		db = mysql{dbConString}
	} else if dbtype == "spanner" {
		db = NewSpanner(dbConString)
	} else {
		fmt.Println("Unknown database backend. Exiting.")
		os.Exit(1)
	}

	// Init database
	db.create()
	seedDatabase(db)

	// Init stackdriver
	myMon := createMon()
	myMon.init()

	go updateTicketsSold(db, myMon)

	// // Write a TimeSeries value for that metric
	// metricType := "custom.googleapis.com/tickets_sold"
	// if err := myMon.writeTimeSeriesValue(42, metricType); err != nil {
	// 	log.Fatal(err)
	// }

	// MUX Router
	rtr := mux.NewRouter()

	// Handler
	h := &handler{&db}

	// run the webserver.
	rtr.HandleFunc("/", h.rootHandler).Methods("GET")
	rtr.HandleFunc("/api/artists", h.artistsHandler).Methods("GET")
	rtr.HandleFunc("/api/shows/{id:[0-9]+}", h.showHandler).Methods("GET")
	rtr.HandleFunc("/api/tickets/{id:[0-9]+}", h.ticketHandler).Methods("GET")
	rtr.HandleFunc("/api/buy/{id:[0-9]+}", h.buyHandler).Methods("GET")
	// static files
	//http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/", rtr)

	fmt.Print("Webserver started.")
	if err := http.ListenAndServe("0.0.0.0:9001", nil); err != nil {
		fmt.Println("Starting webserver failed!")
	}
}

// This function should be database agnostic. Maybe we need to have a database object.
func seedDatabase(db mydatabase) {

	// Put in some artists.
	db.insert(Artist{Id: 1, Name: "Airbourne"})
	db.insert(Artist{Id: 2, Name: "Bloody Beetroots"})
	db.insert(Artist{Id: 3, Name: "Angerfist"})
	db.insert(Artist{Id: 4, Name: "Dieselboy"})
	db.insert(Artist{Id: 5, Name: "SXTN"})
	db.insert(Artist{Id: 6, Name: "Steel Panther"})
	db.insert(Artist{Id: 7, Name: "Rob Zombie"})
	db.insert(Artist{Id: 8, Name: "N.W.A"})
	db.insert(Artist{Id: 9, Name: "NERO"})
	db.insert(Artist{Id: 10, Name: "Heaven Shall Burn"})

	// Put in a show for every artist.
	// We are lazy and just reuse values.

	name := "Live!"
	stimeInt := 1510516800
	etimeInt := 1510524000
	price := 40.50
	maxTickets := 10000000

	stime := time.Unix(int64(stimeInt), 0)
	etime := time.Unix(int64(etimeInt), 0)

	db.insert(Show{Id: 1, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 1})
	db.insert(Show{Id: 2, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 2})
	db.insert(Show{Id: 3, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 3})
	db.insert(Show{Id: 4, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 4})
	db.insert(Show{Id: 5, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 5})
	db.insert(Show{Id: 6, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 6})
	db.insert(Show{Id: 7, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 7})
	db.insert(Show{Id: 8, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 8})
	db.insert(Show{Id: 9, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 9})
	db.insert(Show{Id: 10, Name: name, StartTime: stime, EndTime: etime, Price: price, MaxTickets: int64(maxTickets), ArtistId: 10})

	// Current status: Fuck the user table. We just buy tickets with random user ids.
}

// Update the stackdriver counter of sold tickets. We do some random sleeping so we will not kill stackdriver.
func updateTicketsSold(db mydatabase, m *mon) {

	metricType := "custom.googleapis.com/tickets_sold"
	for {
		ticketsSold := db.ticketsSold()

		m.writeTimeSeriesValue(ticketsSold, metricType)

		timeout := random(1000, 10000)
		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
