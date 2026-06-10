package entity

type Flat struct {
	ID      uint   `json:"id" db:"id"`
	HouseID uint   `json:"house_id" db:"house_id"`
	Number  int    `json:"number" db:"number"`
	Price   int    `json:"price" db:"price"`
	Rooms   int    `json:"rooms" db:"rooms"`
	Status  string `json:"status" db:"status"`
}
