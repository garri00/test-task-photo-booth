package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"test-task-photo-booth/pkg/clients"
	"test-task-photo-booth/pkg/clients/postgresql"
	"test-task-photo-booth/src/entities/customErrors"
	"test-task-photo-booth/src/entities/dtos"
)

type photoPgStorage struct {
	client postgresql.Client
	logger *zerolog.Logger
}

func NewPhotoStoragePG(client postgresql.Client, logger *zerolog.Logger) clients.PhotoStorage {
	return &photoPgStorage{
		client: client,
		logger: logger,
	}
}

type PhotoPG struct {
	ID         null.String `json:"id"`
	DataOrigin null.String `json:"dataOrigin"` // Stored in b64
	Data75     null.String `json:"data75"`     // Stored in b64
	Data50     null.String `json:"data50"`     // Stored in b64
	Data25     null.String `json:"data25"`     // Stored in b64
	IsDeleted  null.Bool   `json:"isDeleted"`
}

func (p photoPgStorage) Create(ctx context.Context, photo *dtos.PhotoDB) error {
	query := `
		INSERT INTO service.photos
		    (
		     data_origin, 
		     data_75,
		     data_50,
		     data_25,
		     is_deleted
		     )
		VALUES
		       ($1, $2, $3, $4, $5)
		RETURNING id
	`

	if err := p.client.QueryRow(ctx, query,
		photo.DataOrigin,
		photo.Data75,
		photo.Data50,
		photo.Data25,
		photo.IsDeleted,
	).Scan(&photo.ID); err != nil {
		return fmt.Errorf("client.QueryRow() failed: %w", err)
	}

	p.logger.Info().Msgf("photo created with id: %v", photo.ID)

	return nil
}

var ErrNoPhotoFound = errors.New("didn't find photo")

func (p photoPgStorage) FindAll(ctx context.Context) ([]dtos.PhotoDB, error) {
	query := `
		SELECT id, 
		       data_origin, 
		       data_75,
		       data_50,
		       data_25,
		       is_deleted
		FROM service.photos;
	`

	photosList := make([]dtos.PhotoDB, 0)

	rows, err := p.client.Query(ctx, query)
	if err != nil {
		return photosList, fmt.Errorf("client.Query() failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photoPG PhotoPG

		err = rows.Scan(
			&photoPG.ID,
			&photoPG.DataOrigin,
			&photoPG.Data75,
			&photoPG.Data50,
			&photoPG.Data25,
			&photoPG.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("client.Query() failed: %w", err)
		}

		photoDB := dtos.PhotoDB{
			ID:         photoPG.ID.String,
			DataOrigin: photoPG.DataOrigin.String,
			Data75:     photoPG.Data75.String,
			Data50:     photoPG.Data50.String,
			Data25:     photoPG.Data25.String,
			IsDeleted:  false,
		}

		photosList = append(photosList, photoDB)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("client.Query() failed: %w", err)
	}

	return photosList, nil
}

func (p photoPgStorage) FindOne(ctx context.Context, id string) (dtos.PhotoDB, error) {
	query := `
		SELECT id, 
		       data_origin, 
		       data_75,
		       data_50,
		       data_25,
		       is_deleted
		FROM service.photos
		WHERE id = $1;
	`

	var photoPG PhotoPG
	err := p.client.QueryRow(ctx, query, id).Scan(
		&photoPG.ID,
		&photoPG.DataOrigin,
		&photoPG.Data75,
		&photoPG.Data50,
		&photoPG.Data25,
		&photoPG.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dtos.PhotoDB{}, ErrNoPhotoFound
		}

		return dtos.PhotoDB{}, fmt.Errorf("client.QueryRow() failed: %w", err)
	}

	photoDB := dtos.PhotoDB{
		ID:         photoPG.ID.String,
		DataOrigin: photoPG.DataOrigin.String,
		Data75:     photoPG.Data75.String,
		Data50:     photoPG.Data50.String,
		Data25:     photoPG.Data25.String,
		IsDeleted:  false,
	}

	return photoDB, nil
}

func (p photoPgStorage) Update(ctx context.Context, photo dtos.PhotoDB) error {
	query := `
		   UPDATE service.photos 
		   SET  
		       data_origin = $1, 
		       data_75 = $2,
		       data_50 = $3,
		       data_25 = $4,
		       is_deleted = $5
           WHERE id = $6;
`

	_, err := p.client.Exec(ctx, query,
		photo.DataOrigin,
		photo.Data75,
		photo.Data50,
		photo.Data25,
		photo.IsDeleted,
		photo.ID,
	)
	if err != nil {
		return fmt.Errorf("client.Query() failed: %w", err)
	}

	return nil
}

func (p photoPgStorage) Delete(ctx context.Context, id string) error {
	query := `
		 UPDATE service.photos 
		   SET is_deleted = TRUE
           WHERE id = $1;
	`

	commandTag, err := p.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("client.Exec() failed: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return customErrors.ErrNoRowsFindToDelete
	}

	p.logger.Debug().Msgf("photo with id = %s sucsefuly DELETED", id)

	return nil
}
