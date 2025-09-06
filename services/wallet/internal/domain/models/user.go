package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/wallet/internal/domain/events"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id string

	FirstName      string
	LastName       string
	Email          string
	HashedPassword string

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func CreateNewUser(firstName, lastName, email, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	currentTime := time.Now()

	user := &User{
		Id:             id,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		HashedPassword: hashedPassword,

		CreatedAt: currentTime,
		UpdatedAt: nil,
	}

	user.AddEvent(events.UserCreatedEvent{
		UserId:    user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Timestamp: currentTime,
	})

	return user, nil
}

func (u *User) AddEvent(event events.DomainEvent) {
	u.events = append(u.events, event)
}

func (u *User) Events() []events.DomainEvent {
	return u.events
}

func (u *User) ClearEvents() {
	u.events = nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
