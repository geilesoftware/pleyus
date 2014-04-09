package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const SessionSecret = "sdngosngpasp4oa7ta8473280gtgadfgh"
const SessionName = "gosess"
const HostName = "localhost"

var store = sessions.NewCookieStore([]byte(SessionSecret))
var hashKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
var blockKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
var db *sql.DB

func initSessionOptions() {
	store.Options = &sessions.Options{
		Domain:   HostName,
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func initDatabaseConnection() {
	ConnectDatabase()
	PingDatabase()
}

func PingDatabase() {
	log.Printf("\n\n__________________________________\n\n\nPinging Database\n\n\n__________________________________\n\n")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDatabase() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=pleyus password=abc sslmode=disable")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("\n\n__________________________________\n\n\nConnected to Database\n\n\n__________________________________\n\n")
	}
}

func GetUser(id int64) {
	var (
		uid        int64
		username   string
		password   string
		hash       string
		last_login time.Time
		created    time.Time
	)
	err := db.QueryRow("SELECT id, username, password, hash, last_login, created FROM users WHERE id = $1", id).Scan(&uid, &username, &password, &hash, &last_login, &created)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Username is %s\n", username)
		fmt.Printf("LastLogin is %s\n", last_login)
	}
}

func IsAuth(r *http.Request) bool {
	session, _ := store.Get(r, SessionName)
	var s = securecookie.New(hashKey, blockKey)
	var secureCookie string
	sessionValue, ok := session.Values["security"].(string)
	if ok {
		s.Decode("security", sessionValue, &secureCookie)
		fmt.Printf("SecurityCookie: %s", secureCookie)
		return true
	} else {
		fmt.Println("SecurityCookie: %s", "ERROR")
		return false
	}
	return false
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, SessionName)

	var s = securecookie.New(hashKey, blockKey)

	// Set some session values.
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	secureCookie, err := s.Encode("security", time.Now().Format(layout))
	if err != nil {
		fmt.Println(err)
	}

	session.Values["security"] = secureCookie
	// Save it.
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, SessionName)

	var s = securecookie.New(hashKey, blockKey)

	// Set some session values.
	var secureCookie string
	sessionValue, ok := session.Values["security"].(string)
	if ok {
		s.Decode("security", sessionValue, &secureCookie)
		fmt.Fprintf(w, "SecurityCookie: %s", secureCookie)
	} else {
		fmt.Fprintf(w, "SecurityCookie: %s", "ERROR")
	}
}

func main() {
	initSessionOptions()
	initDatabaseConnection()
	defer db.Close()
	GetUser(1)
	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/login", LoginHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
