package domain

import "time"

type Weather struct {
	ID                 *int       `json:"id"`
	CreatedAt          *time.Time `json:"created_at"`
	CreatedBy          *string    `json:"created_by"`
	UpdatedAt          *time.Time `json:"updated_at"`
	UpdatedBy          *string    `json:"updated_by"`
	DeletedAt          *time.Time `json:"deleted_at"`
	DeletedBy          *string    `json:"deleted_by"`
	Name               *string    `json:"name"`
	Latitude           *float64   `json:"latitude"`
	Longitude          *float64   `json:"longitude"`
	Icon               *string    `json:"icon"`
	IconNum            *int       `json:"icon_num"`
	Summary            *string    `json:"summary"`
	Temperature        *float32   `json:"temperature"`
	WindSpeed          *float32   `json:"wind_speed"`
	WindAngel          *int       `json:"wind_angel"`
	WindDir            *string    `json:"wind_dir"`
	PrecipitationTotal *float32   `json:"precipitation_total"`
	PrecipitationType  *string    `json:"precipitation_type"`
	CloudCover         *int       `json:"cloud_cover"`
}
