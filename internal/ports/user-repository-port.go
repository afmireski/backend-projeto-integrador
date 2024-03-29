package ports

import (
	"time"

	"github.com/afmireski/garchop-api/internal/models"
	myTypes "github.com/afmireski/garchop-api/internal/types"
)

type UserRepositoryPort interface {
	Create(input CreateUserInput) (string, error)

	FindById(id string) (*models.UserModel, error)

	Update(id string, input myTypes.AnyMap, where myTypes.Where) (myTypes.Any, error)

	Delete(id string) error
}

type CreateUserInput struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
}