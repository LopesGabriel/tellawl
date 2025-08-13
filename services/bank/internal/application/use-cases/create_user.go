package usecases

import (
	"errors"

	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/bank/internal/domain/repository"
)

type createUserUseCase struct {
	userRepository repository.UserRepository
}

type CreateUserUseCaseInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func NewCreateUserUseCase(userRepository repository.UserRepository) *createUserUseCase {
	return &createUserUseCase{
		userRepository: userRepository,
	}
}

func (usecase *createUserUseCase) Execute(input CreateUserUseCaseInput) (*models.User, error) {
	if input.FirstName == "" {
		return nil, MissingRequiredFieldsError("FirstName")
	}

	if input.LastName == "" {
		return nil, MissingRequiredFieldsError("LastName")
	}

	if input.Email == "" {
		return nil, MissingRequiredFieldsError("Email")
	}

	if input.Password == "" {
		return nil, MissingRequiredFieldsError("Password")
	}

	existingUser, err := usecase.userRepository.FindByEmail(input.Email)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.Join(errors.New("could not validate existing user"), err)
		}
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user, err := models.CreateNewUser(input.FirstName, input.LastName, input.Email, input.Password)
	if err != nil {
		return nil, errors.Join(errors.New("could not create user"), err)
	}

	if err := usecase.userRepository.Save(user); err != nil {
		return nil, errors.Join(errors.New("could not persist the user"), err)
	}

	return user, nil
}
