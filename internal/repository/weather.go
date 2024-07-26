package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/wazadio/realtime-weather/internal/config/postgres"
	"github.com/wazadio/realtime-weather/internal/domain"
	"github.com/wazadio/realtime-weather/pkg/logger"
	"github.com/wazadio/realtime-weather/pkg/types"
)

const (
	ERROR_AFFECTED_ROWS    = "rows affected not equal"
	ERROR_RECORD_NOT_FOUND = "record not found"
	SUCCESS_INSERT         = "success insert"
	SUCCESS_UPDATE         = "success update"
)

type weatherRepository struct {
	WriteDB *sql.DB
	ReadDB  *sql.DB
}

type WeatherRepository interface {
	SaveNewCoordinate(ctx context.Context, data domain.Weather) (err error)
	GetWeatherByLatLon(ctx context.Context, latitude, longitude float64) (weathers []domain.Weather, err error)
	GetWeatherByName(ctx context.Context, name string) (weathers []domain.Weather, err error)
}

func NewWeatherRepository(postgresDB *postgres.DB) WeatherRepository {
	return &weatherRepository{
		WriteDB: postgresDB.Write,
		ReadDB:  postgresDB.Read,
	}
}

func (r *weatherRepository) SaveNewCoordinate(ctx context.Context, data domain.Weather) (err error) {
	query := `
		INSERT INTO weather(
			created_at, created_by, updated_at, updated_by,
			name, latitude, longitude,
			icon, icon_num, summary,
			wind_speed, wind_angel, wind_dir,
			precipitation_total, precipitation_type, cloud_cover,
			temperature
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.WriteDB.Begin()
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())
		return
	}
	result, err := tx.ExecContext(
		queryContext, query, time.Now().Local(), types.SYSTEM, time.Now().Local(), types.SYSTEM,
		data.Name, data.Latitude, data.Longitude,
		data.Icon, data.IconNum, data.Summary, data.WindSpeed, data.WindAngel,
		data.WindDir, data.PrecipitationTotal, data.PrecipitationType, data.CloudCover,
		data.Temperature,
	)
	if err != nil {
		tx.Rollback()
		logger.Print(ctx, logger.ERROR, err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		logger.Print(ctx, logger.ERROR, err.Error())
		return
	}

	if rowsAffected != 1 {
		tx.Rollback()
		err = fmt.Errorf("%s, %d", ERROR_AFFECTED_ROWS, rowsAffected)
		logger.Print(ctx, logger.ERROR, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Print(ctx, logger.ERROR, err.Error())
	}

	logger.Print(ctx, logger.INFO, fmt.Errorf("%s, %d", SUCCESS_INSERT, rowsAffected))

	return
}

func (r *weatherRepository) GetWeatherByName(ctx context.Context, name string) (weathers []domain.Weather, err error) {
	query := `
		SELECT *
		FROM weather
		WHERE name = $1
			AND deleted_at is null
	`

	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.ReadDB.QueryContext(queryContext, query, name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		weather := domain.Weather{}

		s := reflect.ValueOf(&weather).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err = rows.Scan(columns...)
		if err != nil {
			return
		}

		weathers = append(weathers, weather)
	}

	return
}

func (r *weatherRepository) GetWeatherByLatLon(ctx context.Context, latitude, longitude float64) (weathers []domain.Weather, err error) {
	query := `
		SELECT *
		FROM weather
		WHERE latitude = $1
			AND longitude = $2
	`

	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.ReadDB.QueryContext(queryContext, query, latitude, longitude)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		weather := domain.Weather{}

		s := reflect.ValueOf(&weather).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err = rows.Scan(columns...)
		if err != nil {
			return
		}

		weathers = append(weathers, weather)
	}

	return
}
