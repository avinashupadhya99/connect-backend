package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Slot struct {
	gorm.Model
	StartTime string         `json:"starttime" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Date      datatypes.Date `json:"-" gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	DateStr   string         `json:"date" gorm:"-"`
	Users     []*User        `json:"users" gorm:"many2many:slot_users;"`
	UserID    int            `json:"userid" gorm:"-"`
}

func BookSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slot Slot
	json.NewDecoder(r.Body).Decode(&slot)
	date, err := time.Parse("2006-01-02", slot.DateStr)
	fmt.Println("date", date)
	if err == nil {
		slot.Date = datatypes.Date(date)
		var user User
		DB.First(&user, slot.UserID)
		if reflect.DeepEqual(user, User{}) {
			http.Error(w, fmt.Sprintf("User with ID %d does not exist", slot.UserID), http.StatusBadRequest)
		} else {
			DB.Preload("Users").Where("start_time = ? AND date = ?", slot.StartTime, slot.DateStr).FirstOrInit(&slot)
			for index := range slot.Users {
				slot.Users[index].Password = ""
			}
			if !userInSlice(user, slot.Users) {
				slot.Users = append(slot.Users, &user)
			}
			DB.Save(&slot)
			json.NewEncoder(w).Encode(slot)
		}
		return
	} else {
		log.Println(err.Error())
		http.Error(w, "Invalid date format, should be YYYY-MM-DD", http.StatusBadRequest)
		return
	}
}

func GetSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slot Slot
	starttime := r.URL.Query().Get("starttime")
	datestr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", datestr)
	fmt.Println("date", date)

	if err == nil {
		// DB.Preload("Users").Where(&Slot{StartTime: starttime, EndTime: endtime, Date: datatypes.Date(date)}).First(&slot)
		DB.Preload("Users").Where("start_time = ? AND date = ?", starttime, datestr).First(&slot)
		slot.DateStr = datestr
		json.NewEncoder(w).Encode(slot)
	} else {
		log.Println(err.Error())
		http.Error(w, "Invalid date format, should be YYYY-MM-DD", http.StatusBadRequest)
	}
}

func DeleteSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slot Slot
	starttime := r.URL.Query().Get("starttime")
	datestr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", datestr)
	fmt.Println("date", date)

	if err == nil {
		DB.Where("start_time = ? AND date = ?", starttime, datestr).First(&slot)
		// DB.Where(&Slot{StartTime: starttime, EndTime: endtime, Date: datatypes.Date(date)}).First(&slot)
		fmt.Println(slot)
		if reflect.DeepEqual(slot, Slot{}) {
			message := fmt.Sprintf("No slot at %s  on %s", starttime, datestr)
			http.Error(w, message, http.StatusBadRequest)
		} else {
			DB.Delete(&slot)
			message := fmt.Sprintf("Slot at %s on %s deleted", starttime, datestr)
			json.NewEncoder(w).Encode(message)
		}
	} else {
		log.Println(err.Error())
		http.Error(w, "Invalid date format, should be YYYY-MM-DD", http.StatusBadRequest)
	}
}

func UnBookSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slot Slot
	starttime := r.URL.Query().Get("starttime")
	datestr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", datestr)
	userid := r.URL.Query().Get("userid")
	fmt.Println("date", date)
	var user User
	if err == nil {
		DB.Where("start_time = ? AND date = ?", starttime, datestr).First(&slot)
		fmt.Println(slot)
		if reflect.DeepEqual(slot, Slot{}) {
			message := fmt.Sprintf("No slot at %s  on %s", starttime, datestr)
			http.Error(w, message, http.StatusBadRequest)
		} else {
			DB.First(&user, userid)
			if reflect.DeepEqual(user, User{}) {
				message := fmt.Sprintf("No user with ID %s", userid)
				http.Error(w, message, http.StatusBadRequest)
			} else {
				DB.Model(&user).Association("Slots").Delete(&slot)
				message := fmt.Sprintf("Slot at %s on %s deleted for user %s", starttime, datestr, userid)
				json.NewEncoder(w).Encode(message)
			}
		}
	} else {
		log.Println(err.Error())
		http.Error(w, "Invalid date format, should be YYYY-MM-DD", http.StatusBadRequest)
	}
}
