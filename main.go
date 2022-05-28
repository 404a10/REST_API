package main

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB
var err error

type Booking struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Members int    `json:"members"`
}

var (
	red    = color.FgRed
	green  = color.FgGreen
	blue   = color.FgBlue
	yellow = color.FgYellow
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage 👌")
	log.Println(yellow.Render("Endpoint Hit: Homepage"))
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body) //reads the request (it's json)
	if err != nil {
		return
	}

	//unmarshalls the json and adds it to the database
	booking := Booking{0, "", 0}
	err = json.Unmarshal(reqBody, &booking)
	if err != nil {
		return
	}

	db.Create(booking)

	log.Println(yellow.Render("Endpoint Hit: Creating New Booking"))
	err = json.NewEncoder(w).Encode(booking)
	if err != nil {
		return
	}
	//we return the user the json they entered

}

func returnAllBookings(w http.ResponseWriter, r *http.Request) {
	var bookings []Booking
	db.Find(&bookings)
	log.Println("Endpoint Hit: Return all bookings")
	err := json.NewEncoder(w).Encode(&bookings)
	if err != nil {
		return
	}
}

func returnSingleBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	bookings := []Booking{}
	db.Find(&bookings)
	var found bool
	found = false

	for _, booking := range bookings {
		s, err := strconv.Atoi(key)
		if err == nil {
			if booking.Id == s {
				found = true
				log.Println(booking)
				log.Println("Endpoint Hit: Booking No: ", key)
				json.NewEncoder(w).Encode(booking)
			}
		}
	}
	if found == false {
		json.NewEncoder(w).Encode("Error 404: Booking not found")
	}
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	bookings := []Booking{}
	db.Find(&bookings)

	for _, booking := range bookings {
		s, err := strconv.Atoi(key)
		if err == nil {
			if booking.Id == s {
				log.Printf(red.Render("Endpoint Hit: Deleted booking Nª", s))
				db.Delete(booking)
				json.NewEncoder(w).Encode("Deleted from the database")
			}
		}
	}
}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	reqbody, _ := ioutil.ReadAll(r.Body)
	bookingId, _ := strconv.Atoi(key)

	booking := Booking{Id: bookingId, User: "", Members: 0}
	err = json.Unmarshal(reqbody, &booking)
	if err != nil {
		return
	}
	db.Save(booking)

	log.Println(yellow.Render("Endpoint hit: Updated booking Nº", key))
	json.NewEncoder(w).Encode(booking)
}

func handleRequests() {
	log.Printf("Starting http server on port %s\n", blue.Render("10000"))
	log.Printf("Stop the server with %s\n", red.Render("CTRL + C"))
	// we create a router to redirect the connection to the webpage
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/new-booking", createBooking).Methods("POST")
	router.HandleFunc("/all-bookings", returnAllBookings)
	router.HandleFunc("/booking/{id}", returnSingleBooking)
	router.HandleFunc("/delete-booking/{id}", deleteBooking)
	router.HandleFunc("/update-booking/{id}", updateBooking)
	log.Fatal(http.ListenAndServe(":10000", router))

}

func main() {
	db, err = gorm.Open("mysql", "potato:potato@tcp(127.0.0.1:3306)/Football?charset=utf8&parseTime=True")

	if err != nil {
		log.Fatalln(red.Render("Failed to open the database."))
	} else {
		log.Printf("%s to the database\n", green.Render("Connected"))
	}

	db.AutoMigrate(Booking{})
	handleRequests()
}
