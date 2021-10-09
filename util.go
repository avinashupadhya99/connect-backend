package main

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
		DB.FirstOrInit(&interest)
		interests = append(interests, interest)
	}
	return interests
}
