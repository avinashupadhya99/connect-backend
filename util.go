package main

func userInSlice(user User, users []User) bool {
	for _, u := range users {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}
