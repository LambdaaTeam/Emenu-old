package services

import (
	"context"

	"github.com/LambdaaTeam/Emenu/cmd/ws/shared"
)

func HandlePacket(ctx context.Context, packet *shared.Packet) shared.Packet {
	// only return the packet with the same data to notifiy ws clients
	return shared.Packet{
		Type:         packet.Type,
		RestaurantID: packet.RestaurantID,
		OrderID:      packet.OrderID,
		ItemID:       packet.ItemID,
		TableID:      packet.TableID,
		Data:         packet.Data,
	}
}
