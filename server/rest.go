package server

import (
	"net/http"

	"context"
	"log"

	"beautifulthings/store"

	"io/ioutil"

	"encoding/json"

	"github.com/gorilla/mux"
)

type RestServer struct {
	s Server
}

type SignInResponse struct {
	EncryptedToken []byte
}

type SetRequest struct {
	Date string
	Ct   []byte
}

func errorOut(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("Error processing %s: %+v", r.RequestURI, err)
}

func (rs *RestServer) signUp(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorOut(w, r, err)
		return
	}
	err = rs.s.SignUp(b)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rs *RestServer) signIn(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorOut(w, r, err)
		return
	}
	token, err := rs.s.SignIn(b)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	resp, err := json.Marshal(SignInResponse{
		EncryptedToken: token,
	})
	if err != nil {
		errorOut(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (rs *RestServer) set(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	sr := &SetRequest{}
	err = json.Unmarshal(b, sr)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	err = rs.s.Set(vars["token"], sr.Date, sr.Ct)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rs *RestServer) enumerate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bts, err := rs.s.Enumerate(vars["token"], vars["from"], vars["to"])
	if err != nil {
		errorOut(w, r, err)
		return
	}

	resp, err := json.Marshal(bts)
	if err != nil {
		errorOut(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func ServeRest(ctx context.Context, addr string, store store.ServerStore) (func(), error) {
	rs := &RestServer{
		s: New(store),
	}

	r := mux.NewRouter()
	r.HandleFunc("/signup", rs.signUp).Methods("POST")
	r.HandleFunc("/signin", rs.signIn).Methods("POST")
	r.HandleFunc("/things", rs.set).Methods("POST").Queries("token", "{token}")
	r.HandleFunc("/things/{from}/{to}", rs.enumerate).Methods("GET").Queries("token", "{token}")

	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	cancel := func() {
		log.Printf("Shutting down server")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	return cancel, nil
}
