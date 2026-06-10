package dto

type CreateHouseInput struct {
	Address   string  `json:"address"`
	Year      int     `json:"year"`
	Developer *string `json:"developer"`
}
