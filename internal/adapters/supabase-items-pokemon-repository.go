package adapters

import (
	"encoding/json"

	myTypes "github.com/afmireski/garchop-api/internal/types"
	"github.com/afmireski/garchop-api/internal/models"
	supabase "github.com/nedpals/supabase-go"
)

type SupabaseItemsRepository struct {
	client *supabase.Client
}

func NewSupabaseItemsRepository(client *supabase.Client) *SupabaseItemsRepository {
	return &SupabaseItemsRepository{
		client: client,
	}
}

func (r *SupabaseItemsRepository) serializeToModel(supabaseData myTypes.AnyMap) (*models.ItemModel, error) {
	jsonData, err := json.Marshal(supabaseData)
	if err != nil {
		return nil, err
	}

	var modelData models.ItemModel
	err = json.Unmarshal(jsonData, &modelData)
	if err != nil {
		return nil, err
	}

	return &modelData, nil
}

func (r *SupabaseItemsRepository) serializeToModels(supabaseData []myTypes.AnyMap) ([]models.ItemModel, error) {
	var modelsData []models.ItemModel = make([]models.ItemModel, 0)
	for _, item := range supabaseData {
		model, err := r.serializeToModel(item)
		if err != nil {
			return nil, err
		}
		modelsData = append(modelsData, *model)
	}
	return modelsData, nil
}

func (r *SupabaseItemsRepository) FindById(id string, where myTypes.Where) (*models.ItemModel, error) {
	panic("implement me")
}

func (r *SupabaseItemsRepository) FindAll(where myTypes.Where) ([]models.ItemModel, error) {
	var supabaseData []myTypes.AnyMap

	query := r.client.DB.From("items").Select("*, pokemons (*, pokemon_types (*, types (*)), tiers (*))").Is("deleted_at", "null")

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


	return r.serializeToModels(supabaseData)
}

func (r *SupabaseItemsRepository) Create(input myTypes.CreateItemInput) (*models.ItemModel, error) {
	var supabaseData []myTypes.AnyMap

	err := r.client.DB.From("items").Insert(input).Execute(&supabaseData); if err != nil {
		return nil, err
	}

	return r.serializeToModel(supabaseData[0])
}

func (r *SupabaseItemsRepository) UpdateMany(input myTypes.AnyMap, where myTypes.Where) ([]models.ItemModel, error) {
	var supabaseData []myTypes.AnyMap

	if len(where) == 0 {
		return []models.ItemModel{}, nil
	}

	query := r.client.DB.From("items").Update(input);
	for column, filter := range where {
		for operator, criteria := range filter {
			query = query.Filter(column, operator, criteria)
		}
	}

	err := query.Execute(&supabaseData); if err != nil {
		return nil, err
	}

	return r.serializeToModels(supabaseData)
}

func (r *SupabaseItemsRepository) Delete(id string) error {
	panic("implement me")
}

