package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func userInSlice(user User, users []*User) bool {
	for _, u := range users {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

func InterestsToStringArray(interests []Interest) []string {
	var interestStringArray []string
	for _, interest := range interests {
		interestStringArray = append(interestStringArray, interest.Name)
	}
	return interestStringArray
}

func StringArrayToInterests(interestStringArray []string) []Interest {
	var interests []Interest
	for _, interestString := range interestStringArray {
		var interest Interest
		interest.Name = interestString
		DB.Where("name = ?", interestString).FirstOrInit(&interest)
		interests = append(interests, interest)
	}
	return interests
}

func ShuffleArray(array []*User) []*User {
	var randomIndex int
	var temp User
	for index := range array {
		randomIndex = int(math.Floor(rand.Float64() * float64(index+1)))
		temp = *array[index]
		*array[index] = *array[randomIndex]
		*array[randomIndex] = temp
	}
	return array
}

func sendEmail(name string, email string, matchName string, matchEmail string, date string, time string) {
	emailContent, fileErr := os.ReadFile("scheduled.html")
	if fileErr == nil {
		from := mail.NewEmail("Team Connect", "avinash@defhacks.co")
		subject := "Meeting Scheduled on Connect"
		to := mail.NewEmail(name, email)
		htmlContent := string(emailContent)
		htmlContent = strings.Replace(htmlContent, "{user}", name, -1)
		htmlContent = strings.Replace(htmlContent, "{user2}", matchName, -1)
		htmlContent = strings.Replace(htmlContent, "{user2Email}", matchEmail, -1)
		htmlContent = strings.Replace(htmlContent, "{date}", date, -1)
		htmlContent = strings.Replace(htmlContent, "{time}", time, -1)
		message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println(err.Error())
		} else {
			fmt.Println(response.StatusCode)
		}
	} else {
		log.Println(fileErr.Error())
	}

}
