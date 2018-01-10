package main

// User is stored in DB
type User struct {
	ID         int		`json:"id,string,omitempty"`
	Email      string	`json:"email,omitempty"`
	First_name string	`json:"firstName,omitempty"`
	Last_name  string	`json:"lastName,omitempty"`
	Password   string	`json:"password,omitempty"`
	Last_tip   int64	`json:"lastTipEpoch,omitempty"`
}

// AuthToken is the session token
type AuthToken struct {
	Api_user int		`json:"apiUser,string,omitempty"`
	Api_key  string		`json:"apiKey,omitempty"`
}

// Record is the lab test results of the user
type Record struct {
	ID             int `json:"id,string,omitempty"`
	User_id        int `json:"userId,string,omitempty"`
	Age            int `json:"age,string,omitempty"`
	Height         int `json:"height,string,omitempty"`
	Weight         int `json:"weight,string,omitempty"`
	Cholesterol    int `json:"cholesterol,string,omitempty"`
	Blood_pressure int `json:"bloodPressure,string,omitempty"`
	Tip_sent       int `json:"tipSent,string,omitempty"`
}

// Login result

type LoginResult struct {
	Token          AuthToken 	`json:"token,omitempty"`
	Email          string		`json:"email,omitempty"`
	First_name     string		`json:"firstName,omitempty"`
	Last_name      string		`json:"lastName,omitempty"`
}