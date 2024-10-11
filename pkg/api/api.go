package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Nathene/MiniBank/internal"
	t "github.com/Nathene/MiniBank/pkg/template"
	"github.com/Nathene/MiniBank/pkg/util"
	jwt "github.com/golang-jwt/jwt/v5" // indirect
	"github.com/gorilla/mux"
)

type APIServer struct {
	listnAddr string
	store     internal.Storage
}

func NewAPIServer(listenAddr string, store internal.Storage) *APIServer {
	return &APIServer{
		listnAddr: listenAddr,
		store:     store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Printf("JSON API server running on port %s", s.listnAddr)

	if err := http.ListenAndServe(s.listnAddr, router); err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return loginGet(w, r)
	}
	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		log.Println("Raw Request Body:", string(body))
		return loginPost(s, w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func loginGet(w http.ResponseWriter, _ *http.Request) error {
	// Serve the Go template
	tmpl, err := template.New("webpage").Parse(t.Tpl)
	if err != nil {
		return err
	}
	// Pass any necessary data to the template
	data := struct {
		Title string
		Items []string
	}{
		Title: "My Login Page",
		Items: []string{"Login Form", "Forgot Password Link"},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func loginPost(s *APIServer, w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		return err
	}

	var req util.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByUsername(req.Username)
	if err != nil {
		return err
	}

	if !acc.ValidatePassword(req.Password) {
		return errors.New("unable to authenticate")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := util.LoginResponse{
		Token:    token,
		Username: acc.Username,
	}

	return writeJson(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowed: %q", r.Method)
}

// GET /account
func (s *APIServer) handleGetAccount(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	id, err := getID(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(util.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}

	account, err := util.NewAccount(accReq.FirstName, accReq.LastName, accReq.Username, accReq.Email, accReq.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	// tokenString, err := createJWT(account)
	// if err != nil {
	// 	return err
	// }

	// // TODO: store this in a cookie or vault
	// fmt.Printf("JWT TOKEN: %s\n", tokenString)

	return writeJson(w, http.StatusCreated, accReq)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(util.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return writeJson(w, http.StatusOK, transferReq)
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			_ = writeJson(w, http.StatusBadRequest, ApiError{err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, err
}

func permissionDenied(w http.ResponseWriter) {
	_ = writeJson(w, http.StatusBadRequest, ApiError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s internal.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Auth Middleware...")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)

		if err != nil {
			// log in grafana etc.
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		UserId, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		acc, err := s.GetAccountByID(UserId)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		// standard map is a float64 for whatever reason, ugly cast back to int64
		if acc.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	// change this to vault later as well.
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func createJWT(acc *util.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": acc.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
