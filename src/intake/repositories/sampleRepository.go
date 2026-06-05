package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)
import "github.com/qmessentials/qmessentials/intake/models"

type SampleRepository interface {
	Get() ([]models.Sample, error)
}

type SampleRepositoryPG struct {
	host     string
	port     string
	database string
	user     string
	password string
}

func NewSampleRepositoryPG(host string, port string, database string, user string, password string) *SampleRepositoryPG {
	return &SampleRepositoryPG{host, port, database, user, password}
}

func (r *SampleRepositoryPG) Get() ([]models.Sample, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", r.user, r.password, r.host, r.port, r.database)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = db.Close() //TODO: handle error
	}()
	results := make([]models.Sample, 0)
	rows, err := db.Query("select id, serial_number, part_number, status, created_at, updated_at from samples")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close() //TODO: handle error
	}()
	for rows.Next() {
		var sample models.Sample
		err = rows.Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, sample)
	}
	return results, err
}
