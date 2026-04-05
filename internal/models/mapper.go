package models

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}