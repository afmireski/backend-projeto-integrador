package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/afmireski/garchop-api/internal/models"
	"github.com/afmireski/garchop-api/internal/ports"
	supabase "github.com/nedpals/supabase-go"

	myTypes "github.com/afmireski/garchop-api/internal/types"
)

type SupabaseUsersRepository struct {
	client *supabase.Client
}

func NewSupabaseUsersRepository(client *supabase.Client) *SupabaseUsersRepository {
	return &SupabaseUsersRepository{
		client: client,
	}
}

func serializeMany(data []map[string]string) ([]models.UserModel, error) {
	timeLayout := "2006-01-02T15:04:05.999999-07:00"

	for _, d := range data {
		for key, value := range d {
			if strings.Contains(key, "birth_date") {
				t, err := time.Parse("2006-01-02", value)
				if err != nil {
					return nil, err
				}
				d[key] = t.Format(timeLayout)
			}
			if strings.Contains(key, "deleted_at") && len(value) == 0 {
				tmp := time.Time{}
				d[key] = tmp.Format(timeLayout)
			}
		}
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	var result []models.UserModel
	json.Unmarshal(jsonData, &result)

	return result, nil
}

func mapToUserModel(data map[string]interface{}) (*models.UserModel, error) {
	birthDate, _ := time.Parse("2006-01-02", data["birth_date"].(string))

	createdAt, _ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", data["created_at"].(string))

	updatedAt, _ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", data["updated_at"].(string))

	var deletedAt time.Time
	if deletedAtString, ok := data["deleted_at"].(string); ok {
		deletedAt, _ = time.Parse("2006-01-02T15:04:05.999999999Z07:00", deletedAtString)
	}

	var role models.UserModelRoleEnum
	if data["role"] == "client" {
		role = models.Client
	} else {
		role = models.Admin
	}

	return models.NewUserModel(
		data["id"].(string),
		data["name"].(string),
		data["email"].(string),
		data["phone"].(string),
		birthDate,
		role,
		createdAt,
		updatedAt,
		deletedAt), nil
}

type CreateInput struct {
	Name      string                   `json:"name"`
	Email     string                   `json:"email"`
	Phone     string                   `json:"phone"`
	Password  string                   `json:"password"`
	BirthDate *time.Time                `json:"birth_date"`
	Role      models.UserModelRoleEnum `json:"role"`
}

func (r *SupabaseUsersRepository) Create(input ports.CreateUserInput) (string, error) {
	var supabaseData []map[string]string

	data := CreateInput{
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
		BirthDate: input.BirthDate,
		Role:      input.Role,
	}

	err := r.client.DB.From("users").Insert(data).Execute(&supabaseData)
	if err != nil {
		return "", err
	}

	// SignUp the user into supabase auth table
	_, signUpErr := r.client.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    input.Email,
		Password: input.PlainPassword,
	})
	if signUpErr != nil {
		return "", err
	}

	return supabaseData[0]["id"], nil
}

func (r *SupabaseUsersRepository) FindById(id string) (*models.UserModel, error) {
	var supabaseData map[string]interface{}

	err := r.client.DB.From("users").Select("*").Single().Eq("id", id).Execute(&supabaseData)

	if err != nil {

		if strings.Contains(err.Error(), "PGRST116") { // resource not found
			return nil, nil
		}

		return nil, err
	}

	return mapToUserModel(supabaseData)
}

func (r *SupabaseUsersRepository) Update(id string, input myTypes.AnyMap, where myTypes.Where) (*models.UserModel, error) {
	var supabaseData []map[string]string
	query := r.client.DB.From("users").Update(input).Eq("id", id)
	if len(where) > 0 {
		for column, filter := range where {
			for operator, criteria := range filter {
				query = query.Filter(column, operator, criteria)
			}
		}
	}

	err := query.Execute(&supabaseData)
	if err != nil {
		return nil, err
	}

	if len(supabaseData) == 0 {
		return nil, nil
	}

	result, err := serializeMany(supabaseData)
	if err != nil {
		return nil, err
	}

	return &result[0], nil
}

func (r *SupabaseUsersRepository) Delete(id string) error {
	return errors.New("not implemented")
}
