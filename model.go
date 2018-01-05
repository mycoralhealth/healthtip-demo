package main

// User is stored in DB
type User struct {
	ID         int
	Email      string
	First_name string
	Last_name  string
	Password   string
}

// AuthToken is the session token
type AuthToken struct {
	Api_user int
	Api_key  string
}

// Record is the lab test results of the user
type Record struct {
	ID             int `json:"id,string,omitempty"`
	User_id        int `json:"User_id,string,omitempty"`
	Age            int `json:"age,string,omitempty"`
	Height         int `json:"height,string,omitempty"`
	Weight         int `json:"weight,string,omitempty"`
	Cholesterol    int `json:"cholesterol,string,omitempty"`
	Blood_pressure int `json:"blood_pressure,string,omitempty"`
}

// Login result

type LoginResult struct {
	Token          AuthToken
	Email          string
	First_name     string
	Last_name      string
}