package main

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Slot struct {
	gorm.Model
	StartTime string `json:"starttime" gorm:"primaryKey"`
	EndTime   string `json:"endtime" gorm:"primaryKey"`
	Date      string `json:"date" gorm:"primaryKey"`
	Users     []User `json:"-" gorm:"many2many:slot_users;"`
	UserID    int    `json:"userid" gorm:"-"`
}

func BookSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slot Slot
	json.NewDecoder(r.Body).Decode(&slot)
	DB.FirstOrInit(&slot)
	var user User
	DB.First(&user, slot.UserID)
	if !userInSlice(user, slot.Users) {
		slot.Users = append(slot.Users, user)
	}
	DB.Save(&slot)
	json.NewEncoder(w).Encode(slot)
}