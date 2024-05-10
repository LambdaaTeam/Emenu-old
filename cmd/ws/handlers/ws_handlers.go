package handlers

import (
	"context"
	"encoding/json"

	"github.com/LambdaaTeam/Emenu/cmd/ws/shared"
	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn            *websocket.Conn
	IsAuthenticated bool
	RestaurantID    string
}

func HandleWebSocketConnection(ctx context.Context, conn *websocket.Conn) {
	client := Client{
		Conn:            conn,
		IsAuthenticated: false,
		RestaurantID:    "",
	}

	// objID, err := primitive.ObjectIDFromHex("662b03d839f330bc5518c21d")

	// if err != nil {
	// 	conn.Close()
	// }

	// token, err := auth.GenerateRestaurantToken(objID)

	// if err != nil {
	// 	conn.Close()
	// }

	// fmt.Println("test_token", token)

	clients[&client] = true

	for {
		var packet shared.Packet

		_, raw_message, err := conn.ReadMessage()

		if err != nil {
			break
		}

		err = json.Unmarshal(raw_message, &packet)

		if err != nil {
			break
		}

		if packet.Type == "" && !client.IsAuthenticated {
			continue
		}

		if client.IsAuthenticated && packet.Type == shared.Auth {
			continue
		}

		if packet.Type == shared.Auth {
			restaurantID, err := auth.DecodeToken(packet.Data)

			if err != nil {
				conn.Close()
			}

			client.IsAuthenticated = true
			client.RestaurantID = restaurantID

			conn.WriteJSON(shared.Packet{
				Type:         shared.Auth,
				RestaurantID: restaurantID,
				Data:         "authenticated",
			})

			continue
		}

		BroadcastMessage(packet)
	}
}

func BroadcastMessage(message shared.Packet) {
	for client := range clients {
		if client.RestaurantID == message.RestaurantID && client.IsAuthenticated && message.Type != shared.Auth {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				client.Conn.Close()
				delete(clients, client)
			}
		}
	}
}
