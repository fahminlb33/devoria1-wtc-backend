package users

import "time"

type UserLoginDto struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	AccessToken string `json:"accessToken"`
}

type UserProfileDto struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"lastModifiedAt"`
}
