package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/coffee22coder/bookings/internal/domain"
)

type AirportRepo struct {
	db *Pool
}

type localizedText struct {
	En string `json:"en"`
	Ru string `json:"ru"`
}

func New(db *Pool) *AirportRepo {
	return &AirportRepo{
		db: db,
	}
}

func (r *AirportRepo) List(ctx context.Context, limit int, offset int) ([]domain.Airport, error) {
	list := make([]domain.Airport, 0, limit)
	rows, err := r.db.pool.Query(ctx, "SELECT * FROM bookings.airports LIMIT=$1 OFFSET=$2", &limit, &offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var code, coordinates, timezone string
		var name, city, country localizedText

		err = rows.Scan(&code, &name, &city, &country, &coordinates, &timezone)
		coordinates = strings.Trim(coordinates, "()")
		parseCoordinate := strings.Split(coordinates, ",")
		if len(parseCoordinate) != 2 {
			return nil, fmt.Errorf("invalid coordinates: %s", coordinates)
		}
		var floatCoordinates [2]float64
		for idx, val := range parseCoordinate {
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			floatCoordinates[idx] = floatVal
		}
		airport := domain.Airport{
			Code: code,
			Name: domain.Location{
				En: name.En,
				Ru: name.Ru,
			},
			City: domain.Location{
				En: city.En,
				Ru: city.Ru,
			},
			Country: domain.Location{
				En: country.En,
				Ru: country.Ru,
			},
			Coordinates: floatCoordinates,
			Timezone:    timezone,
		}
		list = append(list, airport)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AirportRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings.airports").Scan(&count)
	return count, err
}
