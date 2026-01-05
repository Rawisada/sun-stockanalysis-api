package user

import "time"

type User struct {
	ID        	uint 			`gorm:"primaryKeyวautoIncrement;"`
	Email     	string
	Name      	string
	AvatarURL 	string
	Provider  	string 
	CreatedAt 	time.Time
}