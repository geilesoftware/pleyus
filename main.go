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

//var SessionSecret = []byte{233, 28, 242, 52, 20, 142, 126, 59, 127, 4, 148, 145, 217, 100, 197, 60, 213, 63, 74, 98, 107, 71, 42, 226, 191, 25, 51, 66, 231, 27, 21, 236, 40, 197, 39, 178, 42, 61, 219, 216, 149, 174, 146, 86, 7, 218, 168, 94, 100, 67, 60, 145, 230, 65, 181, 84, 13, 6, 188, 21, 97, 37, 48, 85, 248, 237, 12, 241, 12, 7, 213, 194, 237, 157, 27, 39, 30, 81, 231, 205, 191, 56, 183, 212, 122, 109, 63, 151, 219, 59, 79, 17, 21, 67, 118, 137, 66, 171, 240, 194, 219, 239, 124, 189, 3, 140, 168, 171, 175, 212, 81, 185, 87, 26, 69, 247, 208, 212, 192, 103, 113, 51, 169, 33, 123, 111, 27, 255, 212, 190, 168, 55, 14, 37, 86, 68, 45, 35, 83, 61, 112, 13, 234, 63, 152, 253, 122, 199, 98, 12, 13, 208, 255, 46, 173, 159, 127, 33, 89, 115, 101, 231, 153, 113, 212, 232, 214, 152, 147, 150, 185, 121, 47, 104, 157, 204, 203, 182, 255, 68, 131, 205, 13, 52, 240, 125, 184, 244, 55, 160, 118, 71, 172, 61, 253, 39, 164, 107, 38, 101, 47, 96, 66, 196, 181, 6, 3, 143, 234, 160, 65, 23, 251, 248, 74, 223, 166, 191, 50, 185, 218, 194, 171, 92, 116, 131, 160, 233, 137, 224, 178, 184, 242, 95, 216, 90, 50, 32, 1, 99, 154, 119, 179, 18, 13, 65, 58, 25, 201, 186, 215, 217, 189, 155, 162, 173}
const SessionSecret = "sdg_bsidgb_okssiubgie_898z96dfg_bibi43gi_ks"
const SessionName = "gosess"
const HostName = "localhost:8081"

var store = sessions.NewCookieStore([]byte(SessionSecret))
var hashKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
var blockKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
var db *sql.DB

func initSessionOptions() {
	store.Options = &sessions.Options{
		//Domain:   HostName, disable for localhost testing reasons
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
		log.Println(err)
	}
}

func ConnectDatabase() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=pleyus password=abc sslmode=disable")
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("\n\n__________________________________\n\n\nConnected to Database\n\n\n__________________________________\n\n")
	}
}

func GetUser(id int64) (*User, error) {
	var (
		uid        int64
		username   string
		email      string
		password   string
		hash       string
		last_login time.Time
		created    time.Time
	)
	err := db.QueryRow("SELECT id, username, email, password, hash, last_login, created FROM users WHERE id = $1", id).Scan(&uid, &username, &email, &password, &hash, &last_login, &created)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		return nil, err
	case err != nil:
		log.Println(err)
		return nil, err
	default:
		fmt.Printf("Username is %s\n", username)
		fmt.Printf("LastLogin is %s\n", last_login)
		return NewUser(uid, username, email, password, hash, nil), nil
	}
}

func GetUserMock(id int64) (*User, error) {
	return NewUser(id, "username", "email", "password", "hash", nil), nil
}

func GetUserByNameMock(name string) (*User, error) {
	return NewUser(1, name, "email", "password", "hash", nil), nil
}

func GetUserByName(name string) (*User, error) {
	var (
		uid        int64
		username   string
		email      string
		password   string
		hash       string
		last_login time.Time
		created    time.Time
	)
	err := db.QueryRow("SELECT id, username, password, hash, last_login, created FROM users WHERE username = $1", name).Scan(&uid, &username, &email, &password, &hash, &last_login, &created)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that Name.")
		return nil, err
	case err != nil:
		log.Println(err)
		return nil, err
	default:
		fmt.Printf("Username is %s\n", username)
		fmt.Printf("LastLogin is %s\n", last_login)
		return NewUser(uid, username, email, password, hash, nil), nil
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	session, err := store.Get(r, SessionName)
	switch {
	case err == securecookie.ErrMacInvalid:
		log.Println(err)
		log.Println("Recreating session")
		session, err = store.New(r, SessionName)
	default:
		log.Println(err)
	}
	secureCookieUserId, ok := session.Values["security"].(int64)
	log.Println(secureCookieUserId)
	if ok {
		if secureCookieUserId >= 1 {
			user, err := GetUserMock(secureCookieUserId)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(user)
				log.Println("Already auth")
			}
		} else {
			user, err := GetUserByNameMock(username)
			if err != nil {
				log.Println(err)
			} else {
				// log user in
				secureCookieUserId = user.Id
				if err != nil {
					log.Println(err)
				} else {
					log.Println("New auth")
				}
				session.Values["security"] = secureCookieUserId
				// Save it.
				err = session.Save(r, w)
				if err != nil {
					log.Println(err)
				}
			}
		}
	} else {
		fmt.Println("SecurityCookie: %s", "ERROR", "Writing new Session")
		secureCookieUserId = 0
		if err != nil {
			log.Println(err)
		}
		session.Values["security"] = secureCookieUserId
		// Save it.
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
		}
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, SessionName)
	switch {
	case err == securecookie.ErrMacInvalid:
		log.Println(err)
		fmt.Println("Recreating session")
		session, err = store.New(r, SessionName)
	default:
		log.Println(err)
	}

	// Set some session values.
	var defaultUserId int64
	defaultUserId = 0
	secureCookieUserId := defaultUserId
	session.Values["security"] = secureCookieUserId
	// Save it.
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Recreating session?")
		log.Println(err)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	Auth(w, r)
	//r.FormValue("password")
}

func main() {
	initSessionOptions()
	//initDatabaseConnection()
	//defer db.Close()
	//GetUser(1)
	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/login", Auth)
	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
