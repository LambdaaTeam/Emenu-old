package auth_test

import (
	"testing"

	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEncodeRestaurantToken(t *testing.T) {
	restaurantID := primitive.NewObjectID()
	token, err := auth.GenerateRestaurantToken(restaurantID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	decodedRestaurantID, err := auth.DecodeToken(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if decodedRestaurantID != restaurantID.Hex() {
		t.Errorf("Expected %v, got %v", restaurantID.Hex(), decodedRestaurantID)
	}
}

func TestEncodeClientToken(t *testing.T) {
	clientCPF := "12345678900"
	token, err := auth.GenerateClientToken(clientCPF)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	decodedClientCPF, err := auth.DecodeToken(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if decodedClientCPF != clientCPF {
		t.Errorf("Expected %v, got %v", clientCPF, decodedClientCPF)
	}
}
