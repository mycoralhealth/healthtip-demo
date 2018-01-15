package main

import (  
	"gopkg.in/mailgun/mailgun-go.v1"
	"os"
	"fmt"
	"strconv"
	"log"
)

func emailHealthTipRequest(user User, record Record) error {

	log.Println(os.Getenv("MG_DOMAIN"))
	log.Println(os.Getenv("MG_API_KEY"))
	log.Println(os.Getenv("MG_PUBLIC_API_KEY"))

	mailgun.Debug = true

	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"), "")
	message := mg.NewMessage(
		"Health Tips <no-reply@mail.mycoralhealth.com>",
		"New Health Tip Request",
		"Hello, we received a new Health Tip request. Details as follows.\n\nUser: " + user.First_name + " " + user.Last_name +
		"\nEmail: " + user.Email + "\n" +
		"Test result details: \n" + 
		"  Age: " + strconv.Itoa(record.Age) + "\n" + 
		"  Height: " + strconv.Itoa(record.Height) + "\n" + 
		"  Weight: " + strconv.Itoa(record.Weight) + "\n" + 
		"  Cholesterol: " + strconv.Itoa(record.Weight) + "\n" + 
		"  Blood pressure: " + strconv.Itoa(record.Blood_pressure) + "\n\n" + 
		"Please respond in 24 hours as per our demo website SLA.\nThanks!",
		os.Getenv("MAIL_TO"))
	
	resp, id, err := mg.Send(message)
	
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}

func emailPasswordReset(user User, url string) error {

	log.Println(os.Getenv("MG_DOMAIN"))
	log.Println(os.Getenv("MG_API_KEY"))
	log.Println(os.Getenv("MG_PUBLIC_API_KEY"))

	mailgun.Debug = true

	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"), "")
	message := mg.NewMessage(
		"Health Tips <no-reply@mail.mycoralhealth.com>",
		"Password Reset Request",
		"Hello " + user.First_name + " " + user.Last_name + ",\n" +
		"\nWe received a password reset request for your account at My Coral Health - Health Tips\n" +
		"To reset your password please click here: " + url + " \n" + 
		"\nIf this wasn't you or you clicked on the reset password link in error, please disregard this message.\n\nThe My Coral Health Team",
		user.Email)
	
	resp, id, err := mg.Send(message)
	
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}