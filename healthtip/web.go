package healthtip

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/auth0-community/auth0"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	joseutil "github.com/square/go-jose"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func Run(dbCon *sql.DB) error {

	httpAddr := os.Getenv("ADDR")

	mux := makeMuxRouter(dbCon)

	log.Println("Listening on ", httpAddr)

	corsSites := strings.Split(os.Getenv("CROSS_ORIGIN"), ",")

	log.Println("Cross origin sites set to: ", corsSites)

	c := cors.New(cors.Options{
		Debug:          true,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: corsSites,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // Allowing GET, POST, PUT
	})

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        c.Handler(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func authMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sharedKeyPath := os.Getenv("AUTH0_PUBLIC_KEY_PATH")
		audience := []string{os.Getenv("AUTH0_AUDIENCE")}
		sharedKey, err := ioutil.ReadFile(sharedKeyPath)
		if err != nil {
			log.Println("PEM file at " + sharedKeyPath + " could not be read")
			handleError(w, r, http.StatusInternalServerError, "Internal error")
			return
		}
		secret, err := joseutil.LoadPublicKey(sharedKey)
		if err != nil {
			log.Println("Could not load public key from PEM file.")
			handleError(w, r, http.StatusInternalServerError, "Internal error")
			return
		}
		secretProvider := auth0.NewKeyProvider(secret)
		configuration := auth0.NewConfiguration(secretProvider, audience, "https://"+os.Getenv("AUTH0_DOMAIN")+".auth0.com/", jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		// Determine user Id.
		claims := jwt.Claims{}
		token.Claims(secret, &claims)
		userId := strings.Split(claims.Subject, "|")[1]
		ctx := context.WithValue(r.Context(), "userId", userId)

		if err != nil {
			log.Println(err)
			log.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next(w, r.WithContext(ctx))
		}
	}
}

func makeMuxRouter(dbCon *sql.DB) http.Handler {
	wrap := func(f func(w http.ResponseWriter, r *http.Request, dbCon *sql.DB)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			f(w, r, dbCon)
		}
	}

	apiAuth := func(f func(w http.ResponseWriter, r *http.Request, dbCon *sql.DB)) func(w http.ResponseWriter, r *http.Request) {
		return authMiddleware(wrap(f))
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/api/records", apiAuth(handleRecords)).Methods("GET")
	muxRouter.HandleFunc("/api/records", apiAuth(handleRecords)).Methods("POST")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("GET")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("PUT")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}/tip", apiAuth(handleRecordTip)).Methods("POST")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}/approval", apiAuth(handleInsuranceApproval)).Methods("POST")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("DELETE")
	muxRouter.HandleFunc("/api/companies", apiAuth(handleCompanies)).Methods("GET")
	muxRouter.HandleFunc("/api/companies/{companyId:[0-9]+}/procedures/{procedureId:[0-9]+}/policy",
		apiAuth(handleCompanyPolicy)).Methods("GET")
	muxRouter.HandleFunc("/api/procedures", apiAuth(handleProcedures)).Methods("GET")

	return muxRouter
}

func handleError(w http.ResponseWriter, r *http.Request, code int, message string) {
	http.Error(w, message, code)
	log.Println(r.Method, r.URL, ": HTTP", code, ": ", message)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
