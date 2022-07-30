package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mini-gambler/money"
	"net/http"
	"strings"
)

var (
	initialAmount        = money.Money(5000)
	users         []User = []User{
		{ID: 1, Username: "Lewisf95", Balance: &initialAmount},
		{ID: 2, Username: "Kimbom94", Balance: &initialAmount},
	}
)

func main() {
	handleRequests()
	fmt.Println("Starting server on port localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequests() {
	http.Handle("/", http.HandlerFunc(showRoutesHandler))
	http.Handle("/user", http.HandlerFunc(userHandler))
	http.Handle("/bet", http.HandlerFunc(betHandler))
	http.Handle("/credit", http.HandlerFunc(creditHandler))
}

func showRoutesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the showRoutesHandler"))
}

type userResponse struct {
	XMLName  xml.Name `json:"-" xml:"user"`
	Username string   `json:"username" xml:"username,attr"`
	Balance  string   `json:"balance" xml:"balance,attr"`
}
type userRequest struct {
	UserID           int64  `json:"userid"`
	PreferredContent string `json:"preferredcontent,omitempty"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var userReq userRequest
	err = json.Unmarshal(bodyBytes, &userReq)
	if err != nil {
		http.Error(w, "failed to unmarshal request", http.StatusInternalServerError)
		return
	}

	if userReq.UserID == 0 {
		http.Error(w, "invalid request, missing userid", http.StatusBadRequest)
		return
	}

	var userResp userResponse
	var found bool = false
	for _, u := range users {
		if u.ID == userReq.UserID {
			found = true
			userResp = userResponse{
				Username: u.Username,
				Balance:  u.Balance.String(),
			}
			break
		}
	}

	if !found {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var userResponseBytes []byte
	if strings.EqualFold(userReq.PreferredContent, "xml") {
		w.Header().Set("Content-Type", "application/xml")
		userResponseBytes, err = xml.Marshal(userResp)
		if err != nil {
			http.Error(w, "Could not marshal user information", http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		userResponseBytes, err = json.Marshal(userResp)
		if err != nil {
			http.Error(w, "Could not marshal user information", http.StatusInternalServerError)
			return
		}
	}

	w.Write(userResponseBytes)
}

type transactionResponse struct {
	XMLName  xml.Name `json:"-" xml:"transaction"`
	Username string   `json:"username" xml:"username,attr"`
	Balance  string   `json:"balance" xml:"balance,attr"`
	Message  string   `json:"msg" xml:"msg"`
}

type transactionRequest struct {
	UserID           int64  `json:"userid"`
	Amount           int64  `json:"amount"`
	PreferredContent string `json:"preferredcontent,omitempty"`
}

func betHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var betReq transactionRequest
	err = json.Unmarshal(bodyBytes, &betReq)
	if err != nil {
		http.Error(w, "failed to unmarshal request", http.StatusInternalServerError)
		return
	}

	var user User
	for _, u := range users {
		if u.ID == betReq.UserID {
			user = u
			break
		}
	}

	moneyAmount := money.Money(betReq.Amount)
	if moneyAmount > *user.Balance {
		http.Error(w, "user does not have sufficient funds for this transaction", http.StatusBadRequest)
		return
	}

	user.Debit(moneyAmount)

	betResponse := transactionResponse{
		Username: user.Username,
		Balance:  user.Balance.String(),
		Message:  fmt.Sprintf("%s debited %s", user.Username, moneyAmount.String()),
	}

	var respBytes []byte
	if strings.EqualFold(betReq.PreferredContent, "xml") {
		w.Header().Set("Content-Type", "application/xml")
		respBytes, err = xml.Marshal(betResponse)
		if err != nil {
			http.Error(w, "failed to marshal bet response", http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		respBytes, err = json.Marshal(betResponse)
		if err != nil {
			http.Error(w, "failed to marshal bet response", http.StatusInternalServerError)
			return
		}
	}

	w.Write(respBytes)
}

func creditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var creditReq transactionRequest
	err = json.Unmarshal(bodyBytes, &creditReq)
	if err != nil {
		http.Error(w, "failed to unmarshal request", http.StatusInternalServerError)
		return
	}

	var user User
	for _, u := range users {
		if u.ID == creditReq.UserID {
			user = u
			break
		}
	}

	moneyAmount := money.Money(creditReq.Amount)
	user.Credit(moneyAmount)

	creditResponse := transactionResponse{
		Username: user.Username,
		Balance:  user.Balance.String(),
		Message:  fmt.Sprintf("%s credited %s", user.Username, moneyAmount.String()),
	}

	var respBytes []byte
	if strings.EqualFold(creditReq.PreferredContent, "xml") {
		w.Header().Set("Content-Type", "application/xml")
		respBytes, err = xml.Marshal(creditResponse)
		if err != nil {
			http.Error(w, "failed to marshal bet response", http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		respBytes, err = json.Marshal(creditResponse)
		if err != nil {
			http.Error(w, "failed to marshal bet response", http.StatusInternalServerError)
			return
		}
	}

	w.Write(respBytes)
}
