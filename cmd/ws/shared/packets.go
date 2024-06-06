package shared

const (
	Auth              = "auth"
	UpdateOrderStatus = "update_order_status"
	UpdateItemStatus  = "update_item_status"
	UpdateTableStatus = "update_table_status"
	Heartbeat         = "heartbeat"
	Error             = "error"
)

type Packet struct {
	Type         string `json:"type"`
	RestaurantID string `json:"restaurant_id"`
	OrderID      string `json:"order_id"`
	ItemID       string `json:"item_id"`
	TableID      string `json:"table_id"`
	Data         string `json:"data"`
}

func NewErrorPacket(restaurantID, data string) Packet {
	return Packet{
		Type:         Error,
		RestaurantID: restaurantID,
		Data:         data,
	}
}
