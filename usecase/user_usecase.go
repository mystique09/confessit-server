package usecase

import (
	"cnfs/domain"
)

type signUpUsecase struct {
	userRepository domain.IUserRepository
}

func NewSignUpUsecase(userRepository domain.IUserRepository) domain.ISignupUserUseCase {
	return &signUpUsecase{userRepository}
}

func (signup_usecase signUpUsecase) CheckUsernameAvailability(username domain.IUsername) bool {
	_, err := signup_usecase.userRepository.FindByUsername(username)
	return err != nil
}

func (signup_usecase signUpUsecase) Signup(payload domain.CreateUserDTO) domain.UserResponse {
	new_user := payload.ToUser()
	save, err := signup_usecase.userRepository.Create(new_user)

	if err != nil {
		return domain.UserResponse{
			Message: domain.USER_ALREADY_EXISTS,
		}
	}

	return save.IntoResponse()
}
