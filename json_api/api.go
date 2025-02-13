package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
type APIError struct {	
	Error string `json:"error"`
}
func makeHTTPHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
func withJWT(fn http.HandlerFunc,s Storage) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
	 fmt.Println("calling JWT middleware")
	 tokenStr:=r.Header.Get("Authorization")
	 
	 token,err:=validateJWTToken(tokenStr)
	 if err != nil {
		 WriteJSON(w, http.StatusUnauthorized, APIError{Error: "access denied"})
		 return
	 }
	 if !token.Valid {
		 WriteJSON(w, http.StatusUnauthorized, APIError{Error: "access denied"})
		 return
	 }
	 idStr:=mux.Vars(r)["id"]
	 id,err:=strconv.Atoi(idStr)
	 if err != nil {
		 WriteJSON(w, http.StatusBadRequest, APIError{Error: "access denied"})
		 return
	 }
	 acc,err:=s.GetAccountByID(id)
	 if err != nil {
		 WriteJSON(w, http.StatusNotFound, APIError{Error: "account not found"})
		 return
	 }
	 if (acc.FirstName+" "+acc.LastName)!= token.Claims.(jwt.MapClaims)["username"] {
		 WriteJSON(w, http.StatusUnauthorized, APIError{Error: "access denied"})
		 return
	}

	 fn(w, r)
	}
}
func createJWTToken(acc *Account)(string,error) {
	// Define claims (payload)
	secretKey:=os.Getenv("JWT_SECRET")
	// payload
	claims := jwt.MapClaims{
		"username": acc.FirstName + " " + acc.LastName,
		// "number": acc.Number,
		"exp":time.Now().Add(time.Hour * 1).Unix(), 
	}

	// Create token with claims,alg:HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
// JWTtoken: for ASA
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkzNTI3MjMsIm51bWJlciI6ODkyMzc1NzY1Njk5MTE2ODI0MX0.a_H-GtLLKbZDIZkX3ZTYBi04La1x2OZG5TSVKBLYaQ8
func validateJWTToken(tokenStr string)(*jwt.Token,error) {
	secretKey:=os.Getenv("JWT_SECRET")
	token,err:=jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey),nil
	})

	return token,err
}
type APIServer struct {
	listenAddr string
	store 	Storage
}


func NewAPIServer(listenAddr string,store Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr, store: store}
}
func (s *APIServer) Start() error {
	router:=mux.NewRouter()	
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWT(makeHTTPHandlerFunc(s.handleGetAccountByID),s.store)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleDeleteAccount)).Methods("DELETE")
	log.Printf("Starting server on %s", s.listenAddr)	
	// Don't use GET for privacy reasons,browser history can be seen
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.handleTransaction))
	http.ListenAndServe(s.listenAddr, router)
	return nil
}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {	
	idStr:=mux.Vars(r)["id"]
	// convert id to int
	id, err:=strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %s", idStr)
	}
	

	acc,err:=s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	acc_req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(acc_req); err != nil {
		return err
	}
	acc := NewAccount(acc_req.FirstName, acc_req.LastName)
	if err:=s.store.CreateAccount(acc); err != nil {
		return err
	}
	tokenStr,err:=createJWTToken(acc)
	if err != nil {
		return err
	}
	fmt.Printf("JWTtoken: %s\n", tokenStr)
	return WriteJSON(w, http.StatusOK, acc)
}


func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr:=mux.Vars(r)["id"]
	id,err:=strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id %s", idStr)
	}
	if err:=s.store.DeleteAccount(id); err != nil {
		return fmt.Errorf("error deleting account: %v", err)
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": idStr})
}	

func (s* APIServer) handleTransaction(w http.ResponseWriter, r *http.Request) error {
	transfer:=new(TransferRequest)
	if err:=json.NewDecoder(r.Body).Decode(transfer); err != nil {
		return err
	}
	defer r.Body.Close()


	return WriteJSON(w, http.StatusOK, transfer)
}


