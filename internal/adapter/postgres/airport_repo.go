package postgres

import (
	"context"
	"fmt"

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
	rows, err := r.db.pool.Query(ctx, `SELECT
			airport_code,
			airport_name,
			city,
			country,
			coordinates[0],
			coordinates[1],
			timezone
		FROM bookings.airports_data
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var code, timezone string
		var name, city, country localizedText
		var lon, lat float64

		if err := rows.Scan(&code, &name, &city, &country, &lon, &lat, &timezone); err != nil {
			return nil, fmt.Errorf("scan airport: %w", err)
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
			Coordinates: [2]float64{lon, lat},
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
	err := r.db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings.airports_data").Scan(&count)
	return count, err
}
