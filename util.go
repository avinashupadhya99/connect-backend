package main

import (
	"math"
	"math/rand"
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

func ShuffleArray(array []User) []User {
	var randomIndex int
	var temp User
	for index := range array {
		randomIndex = int(math.Floor(rand.Float64() * float64(index+1)))
		temp = array[index]
		array[index] = array[randomIndex]
		array[randomIndex] = temp
	}
	return array
}
