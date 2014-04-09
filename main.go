package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

const SessionSecret = "sdngosngpasp4oa7ta8473280gtgadfgh"
const SessionName = "gosess"
const HostName = "localhost"

var store = sessions.NewCookieStore([]byte(SessionSecret))
var hashKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
var blockKey = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}

func initSessionOptions() {
	store.Options = &sessions.Options{
		Domain:   HostName,
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
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

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if IsAuth(r) == false {
	}
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
	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/login", LoginHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8081", nil)
}
