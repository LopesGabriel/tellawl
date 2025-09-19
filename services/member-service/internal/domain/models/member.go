package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/events"
	"golang.org/x/crypto/bcrypt"
)

type Member struct {
	Id string

	FirstName      string
	LastName       string
	Email          string
	HashedPassword string

	CreatedAt time.Time
	UpdatedAt *time.Time

	events []events.DomainEvent
}

func CreateNewUser(firstName, lastName, email, password string) (*Member, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	currentTime := time.Now()

	member := &Member{
		Id:             id,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		HashedPassword: hashedPassword,

		CreatedAt: currentTime,
		UpdatedAt: nil,
	}

	member.AddEvent(events.MemberCreatedEvent{
		UserId:    member.Id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     member.Email,
		Timestamp: currentTime,
	})

	return member, nil
}

func (u *Member) AddEvent(event events.DomainEvent) {
	u.events = append(u.events, event)
}

func (u *Member) Events() []events.DomainEvent {
	return u.events
}

func (u *Member) ClearEvents() {
	u.events = nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *Member) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
