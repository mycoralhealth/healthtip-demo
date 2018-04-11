package healthtip

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/mailgun/mailgun-go.v1"
)

func emailHealthTipRequest(userinfo UserInfo, record Record) error {

	log.Println(os.Getenv("MG_DOMAIN"))
	log.Println(os.Getenv("MG_API_KEY"))
	log.Println(os.Getenv("MG_PUBLIC_API_KEY"))

	mailgun.Debug = true

	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"), "")
	message := mg.NewMessage(
		"Health Tips <no-reply@mail.mycoralhealth.com>",
		"New Health Tip Request",
		"Hello, you received a new Health Tip request. Details as follows.\n\nUser: "+userinfo.Name+ //+user.FirstName+" "+user.LastName+
			"\nEmail: "+userinfo.Email+"\n"+
			"Test result details: \n"+
			"  Age: "+strconv.Itoa(record.Age)+"\n"+
			"  Height: "+strconv.Itoa(record.Height)+"\n"+
			"  Weight: "+strconv.Itoa(record.Weight)+"\n"+
			"  Heart Rate: "+strconv.Itoa(record.Weight)+"\n"+
			"  Breath Rate: "+strconv.Itoa(record.BloodPressure)+"\n\n"+
			"Please respond in 48 hours as per our demo website SLA.\nThanks!",
		os.Getenv("MAIL_TO"))

	resp, id, err := mg.Send(message)

	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}
