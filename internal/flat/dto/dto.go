package dto

type CreateFlatInput struct {
	HouseID uint `json:"house_id"`
	Number  int  `json:"number"`
	Price   int  `json:"price"`
	Rooms   int  `json:"rooms"`
}

type UpdateStatusuInput struct {
	FlatID    uint   `json:"id"`
	NewStatus string `json:"status"`
}