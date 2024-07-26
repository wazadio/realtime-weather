package request

type SaveNewCoordinate struct {
	Name      *string  `json:"name" binding:"required"`
	Latitude  *float64 `json:"latitude" binding:"required,latitude"`
	Longitude *float64 `json:"longitude" binding:"required,longitude"`
}
