package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string     `json:"firstname"`
	LastName  string     `json:"lastname"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Slots     []*Slot    `json:"slots" gorm:"many2many:slot_users;"`
	About     string     `json:"about"`
	Interests []Interest `json:"-" gorm:"many2many:user_interests;"`
	Interest  []string   `json:"interests" gorm:"-"`
	Meetings  []*Meeting `json:"meetings" gorm:"many2many:meeting_users;"`
}

type Interest struct {
	gorm.Model
	Name string `json:"name"`
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	DB.Preload("Interests").Preload("Slots").Find(&users)
	for index := range users {
		users[index].Password = ""
	}
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.Preload("Interests").Preload("Slots").First(&user, params["id"])
	for index := range user.Slots {
		date, err := user.Slots[index].Date.Value()
		datestr := fmt.Sprintf("%s", date)
		if err == nil && len(datestr) >= 10 {
			user.Slots[index].DateStr = datestr[0:10]
		}
	}
	user.Password = ""
	user.Interest = InterestsToStringArray(user.Interests)
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	user.Password = GetMD5Hash(user.Password)
	user.Interests = StringArrayToInterests(user.Interest)
	DB.Create(&user)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	passwordHash := GetMD5Hash(user.Password)
	DB.Preload("Interests").Preload("Slots").First(&user)
	if user.Password == passwordHash {
		user.Password = ""
		user.Interest = InterestsToStringArray(user.Interests)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
	}
}

func GetUserMeetings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user User
	DB.Preload("Meetings").First(&user, params["id"])
	for index, meeting := range user.Meetings {
		DB.Preload("Interests").Preload("Slot").Preload("Users").First(&meeting)
		date, err := meeting.Slot.Date.Value()
		datestr := fmt.Sprintf("%s", date)
		if err == nil && len(datestr) >= 10 {
			meeting.Slot.DateStr = datestr[0:10]
		}
		user.Meetings[index] = meeting
	}
	json.NewEncoder(w).Encode(user.Meetings)
}
