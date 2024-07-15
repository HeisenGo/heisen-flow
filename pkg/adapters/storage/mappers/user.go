package mappers

import (
	"server/internal/user"
	"server/pkg/adapters/storage/entities"
	"server/pkg/fp"
)

func UserEntityToDomain(entity *entities.User) *user.User {
	return &user.User{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
		Password:  entity.Password,
		Role:      user.Role(entity.Role),
	}
}
func userEntityToDomain(entity entities.User) user.User {
	return user.User{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
		Password:  entity.Password,
		Role:      user.Role(entity.Role),
	}
}

func BatchUserEntityToDomain(entities []entities.User) []user.User {
	return fp.Map(entities, userEntityToDomain)
}

func UserDomainToEntity(domainUser *user.User) *entities.User {
	return &entities.User{
		FirstName: domainUser.FirstName,
		LastName:  domainUser.LastName,
		Email:     domainUser.Email,
		Password:  domainUser.Password,
	}
}
