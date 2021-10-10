package main

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model
	SlotID int     `gorm:"default:0;type:int;not null;comment:'Slot ID'"`
	Slot   Slot    `json:"slot" gorm:"foreignKey:SlotID;references:ID"`
	Users  []*User `json:"users" gorm:"many2many:meeting_users;"`
}

func createMeetings() {
	var slots []Slot
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	if err == nil {
		// Fetch all slots for the next day
		DB.Preload("Users").Where("date >= ? AND date < ?", datatypes.Date(today), datatypes.Date(tomorrow)).Find(&slots)
	} else {
		fmt.Println(err.Error())
	}
	for index, slot := range slots {
		// If the slot has more than one user
		if len(slot.Users) > 1 {
			slots[index].Users = ShuffleArray(slot.Users)
			for userIndex, user := range slot.Users {
				// Pair users
				if userIndex%2 == 1 {
					user2 := slot.Users[userIndex-1]
					// Fetch existing meetings of user to check if already booked
					DB.Preload("Meetings").First(&user)
					DB.Preload("Meetings").First(&user2)
					user1Meetings := user.Meetings
					user1HasMeeting := false
					for _, meeting := range user1Meetings {
						DB.Preload("Slot").First(&meeting)
						if meeting.Slot.ID == slot.ID {
							user1HasMeeting = true
						}
					}
					user2Meetings := user2.Meetings
					user2HasMeeting := false
					for _, meeting := range user2Meetings {
						DB.Preload("Slot").First(&meeting)
						if meeting.Slot.ID == slot.ID {
							user2HasMeeting = true
						}
					}
					if !user1HasMeeting && !user2HasMeeting {
						var meeting Meeting
						meeting.Slot = slot
						meeting.Users = append(meeting.Users, user)
						meeting.Users = append(meeting.Users, user2)
						DB.Create(&meeting)
						fmt.Printf("Meeting with ID %d created \n", meeting.ID)
					}
				}
			}
			// Check if an extra user is present who is not matched in meeting
			if len(slot.Users)%2 == 1 {
				// TODO: Add user to existing meeting or cancel user's slot
			}
		} else {
			// TODO: Let user know, no one signed up for the slot
		}
	}
}
