package models_test

import (
	"testing"
	"time"

	pkg "github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMenus(t *testing.T) {
	t.Run("TestMenuToPublic", func(t *testing.T) {
		// TestMenuToPublic tests the conversion of a Menu to a PublicMenu
		menu := models.Menu{
			ID:           primitive.NewObjectID(),
			RestaurantID: primitive.NewObjectID(),
			Highlights:   []models.Item{},
			Categories:   []models.Category{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		publicMenu := menu.ToPublic()

		if publicMenu.ID != menu.ID {
			t.Errorf("Expected ID to be %v, got %v", menu.ID, publicMenu.ID)
		}

		if publicMenu.RestaurantID != menu.RestaurantID {
			t.Errorf("Expected RestaurantID to be %v, got %v", menu.RestaurantID, publicMenu.RestaurantID)
		}

		if len(publicMenu.Highlights) != len(menu.Highlights) {
			t.Errorf("Expected Highlights to be %v, got %v", menu.Highlights, publicMenu.Highlights)
		}

		if len(publicMenu.Categories) != len(menu.Categories) {
			t.Errorf("Expected Categories to be %v, got %v", menu.Categories, publicMenu.Categories)
		}

		if publicMenu.CreatedAt != menu.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", menu.CreatedAt, publicMenu.CreatedAt)
		}

		if publicMenu.UpdatedAt != menu.UpdatedAt {
			t.Errorf("Expected UpdatedAt to be %v, got %v", menu.UpdatedAt, publicMenu.UpdatedAt)
		}
	})
}

func TestOrders(t *testing.T) {
	t.Run("TestOrderToPublic", func(t *testing.T) {
		// TestOrderToPublic tests the conversion of an Order to a PublicOrder
		order := models.Order{
			ID:           primitive.NewObjectID(),
			RestaurantID: primitive.NewObjectID(),
			TableID:      "1",
			Status:       models.OrderStatusOpen,
			Value:        0,
			Client: models.Client{
				Name: "Felipe Kamada",
				CPF:  "12345678900",
			},
			Items:     []models.OrderItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		publicOrder := order.ToPublic()

		if publicOrder.ID != order.ID {
			t.Errorf("Expected ID to be %v, got %v", order.ID, publicOrder.ID)
		}

		if publicOrder.RestaurantID != order.RestaurantID {
			t.Errorf("Expected RestaurantID to be %v, got %v", order.RestaurantID, publicOrder.RestaurantID)
		}

		if publicOrder.TableID != order.TableID {
			t.Errorf("Expected TableID to be %v, got %v", order.TableID, publicOrder.TableID)
		}

		if publicOrder.Status != order.Status {
			t.Errorf("Expected Status to be %v, got %v", order.Status, publicOrder.Status)
		}

		if publicOrder.Value != order.Value {
			t.Errorf("Expected Value to be %v, got %v", order.Value, publicOrder.Value)
		}

		if publicOrder.Client.Name != order.Client.Name {
			t.Errorf("Expected Client.Name to be %v, got %v", order.Client.Name, publicOrder.Client.Name)
		}

		if publicOrder.Client.CPF != order.Client.CPF {
			t.Errorf("Expected Client.CPF to be %v, got %v", order.Client.CPF, publicOrder.Client.CPF)
		}

		if len(publicOrder.Items) != len(order.Items) {
			t.Errorf("Expected Items to be %v, got %v", order.Items, publicOrder.Items)
		}

		if publicOrder.CreatedAt != order.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", order.CreatedAt, publicOrder.CreatedAt)
		}

		if publicOrder.UpdatedAt != order.UpdatedAt {
			t.Errorf("Expected UpdatedAt to be %v, got %v", order.UpdatedAt, publicOrder.UpdatedAt)
		}
	})
}

func TestRestaurants(t *testing.T) {
	t.Run("TestRestaurantToPublic", func(t *testing.T) {
		// TestRestaurantToPublic tests the conversion of a Restaurant to a PublicRestaurant
		restaurant := models.Restaurant{
			ID: primitive.NewObjectID(),
			Address: struct {
				City     string `json:"city"`
				Country  string `json:"country"`
				PostCode string `json:"post_code"`
				Number   int    `json:"number"`
				Street   string `json:"street"`
				Other    string `json:"other"`
			}{
				City:     "São Paulo",
				Country:  "Brasil",
				PostCode: "12345678",
				Number:   123,
				Street:   "Rua Teste",
				Other:    "Apt 123",
			},
			Tables:    []models.Table{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		publicRestaurant := restaurant.ToPublic()

		if publicRestaurant.ID != restaurant.ID {
			t.Errorf("Expected ID to be %v, got %v", restaurant.ID, publicRestaurant.ID)
		}

		if publicRestaurant.Name != restaurant.Name {
			t.Errorf("Expected Name to be %v, got %v", restaurant.Name, publicRestaurant.Name)
		}

		if publicRestaurant.Address.City != restaurant.Address.City {
			t.Errorf("Expected Address.City to be %v, got %v", restaurant.Address.City, publicRestaurant.Address.City)
		}

		if publicRestaurant.Address.Country != restaurant.Address.Country {
			t.Errorf("Expected Address.Country to be %v, got %v", restaurant.Address.Country, publicRestaurant.Address.Country)
		}

		if publicRestaurant.Address.PostCode != restaurant.Address.PostCode {
			t.Errorf("Expected Address.PostCode to be %v, got %v", restaurant.Address.PostCode, publicRestaurant.Address.PostCode)
		}

		if publicRestaurant.Address.Number != restaurant.Address.Number {
			t.Errorf("Expected Address.Number to be %v, got %v", restaurant.Address.Number, publicRestaurant.Address.Number)
		}

		if publicRestaurant.Address.Street != restaurant.Address.Street {
			t.Errorf("Expected Address.Street to be %v, got %v", restaurant.Address.Street, publicRestaurant.Address.Street)
		}

		if publicRestaurant.Address.Other != restaurant.Address.Other {
			t.Errorf("Expected Address.Other to be %v, got %v", restaurant.Address.Other, publicRestaurant.Address.Other)
		}

		if len(publicRestaurant.Tables) != len(restaurant.Tables) {
			t.Errorf("Expected Tables to be %v, got %v", restaurant.Tables, publicRestaurant.Tables)
		}

		if publicRestaurant.CreatedAt != restaurant.CreatedAt {
			t.Errorf("Expected CreatedAt to be %v, got %v", restaurant.CreatedAt, publicRestaurant.CreatedAt)
		}

		if publicRestaurant.UpdatedAt != restaurant.UpdatedAt {
			t.Errorf("Expected UpdatedAt to be %v, got %v", restaurant.UpdatedAt, publicRestaurant.UpdatedAt)
		}
	})

	t.Run("TestRestaurantRegiterToRestaurant", func(t *testing.T) {
		// TestRestaurantRegiterToRestaurant tests the conversion of a RestaurantRegiter to a Restaurant
		restaurantRegister := models.RestaurantRegister{
			Name:     "Test Restaurant",
			Email:    "test@email.com",
			Password: "123456",
			Address: struct {
				City     string `json:"city"`
				Country  string `json:"country"`
				PostCode string `json:"post_code"`
				Number   int    `json:"number"`
				Street   string `json:"street"`
				Other    string `json:"other"`
			}{
				City:     "São Paulo",
				Country:  "Brasil",
				PostCode: "12345678",
				Number:   123,
				Street:   "Rua Test",
				Other:    "Apt 123",
			},
		}

		restaurant := restaurantRegister.ToRestaurant()

		if restaurant.Name != restaurantRegister.Name {
			t.Errorf("Expected Name to be %v, got %v", restaurantRegister.Name, restaurant.Name)
		}

		if restaurant.Email != restaurantRegister.Email {
			t.Errorf("Expected Email to be %v, got %v", restaurantRegister.Email, restaurant.Email)
		}

		if !pkg.IsPasswordValid(restaurant.Password, restaurantRegister.Password) {
			t.Errorf("Expected Password to be hashed, got %v", restaurant.Password)
		}

		if restaurant.Address.City != restaurantRegister.Address.City {
			t.Errorf("Expected Address.City to be %v, got %v", restaurantRegister.Address.City, restaurant.Address.City)
		}

		if restaurant.Address.Country != restaurantRegister.Address.Country {
			t.Errorf("Expected Address.Country to be %v, got %v", restaurantRegister.Address.Country, restaurant.Address.Country)
		}

		if restaurant.Address.PostCode != restaurantRegister.Address.PostCode {
			t.Errorf("Expected Address.PostCode to be %v, got %v", restaurantRegister.Address.PostCode, restaurant.Address.PostCode)
		}

		if restaurant.Address.Number != restaurantRegister.Address.Number {
			t.Errorf("Expected Address.Number to be %v, got %v", restaurantRegister.Address.Number, restaurant.Address.Number)
		}

		if restaurant.Address.Street != restaurantRegister.Address.Street {
			t.Errorf("Expected Address.Street to be %v, got %v", restaurantRegister.Address.Street, restaurant.Address.Street)
		}

		if restaurant.Address.Other != restaurantRegister.Address.Other {
			t.Errorf("Expected Address.Other to be %v, got %v", restaurantRegister.Address.Other, restaurant.Address.Other)
		}

		if len(restaurant.Tables) != 0 {
			t.Errorf("Expected Tables to be empty, got %v", restaurant.Tables)
		}

		if restaurant.SchemaVersion != 1 {
			t.Errorf("Expected SchemaVersion to be 1, got %v", restaurant.SchemaVersion)
		}

		if restaurant.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set, got %v", restaurant.CreatedAt)
		}

		if restaurant.UpdatedAt.IsZero() {
			t.Errorf("Expected UpdatedAt to be set, got %v", restaurant.UpdatedAt)
		}
	})
}
