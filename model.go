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
	User_id        int
	Age            int
	Height         int
	Weight         int
	Cholesterol    int
	Blood_pressure int
}

// Login result

type LoginResult struct {
	Token AuthToken
	Email string
	First_name string
	Last_name string
}