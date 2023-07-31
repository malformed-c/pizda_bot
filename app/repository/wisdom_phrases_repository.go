package repository

import (
	"context"
	"database/sql"
	"time"
)

type MysqlUserWisdomRepository struct {
	db  *sql.DB
	ctx context.Context
}

func UserWisdomRepository(ctx context.Context, db *sql.DB) *MysqlUserWisdomRepository {
	return &MysqlUserWisdomRepository{
		db: db, ctx: ctx,
	}
}

func (r MysqlUserWisdomRepository) Get(authorId string) (*WisdomPhrase, error) {
	query := "SELECT id, text, author_id FROM wisdom_phrases where author_id = ? ORDER BY RAND() LIMIT 1;"

	ctx, cancel := context.WithTimeout(r.ctx, 15*time.Second)
	defer cancel()

	var phrase WisdomPhrase
	err := r.fetchRow(r.db.QueryRowContext(ctx, query, authorId), &phrase)
	if err != nil {
		return nil, err
	}
	return &phrase, nil
}

func (r MysqlUserWisdomRepository) Create(phrase *WisdomPhrase) (*WisdomPhrase, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 15*time.Second)
	defer cancel()

	query := `INSERT INTO wisdom_phrases (text, author)
		VALUES (?, ?)
		RETURNING id, text, author_id;`
	row := r.db.QueryRowContext(
		ctx,
		query,
		phrase.Text,
		phrase.AuthorId,
	)
	err := r.fetchRow(row, phrase)
	if err != nil {
		return nil, err
	}
	return phrase, nil
}

func (r *MysqlUserWisdomRepository) fetchRow(row *sql.Row, phrase *WisdomPhrase) error {
	return row.Scan(
		&phrase.ID,
		&phrase.Text,
		&phrase.AuthorId,
	)
}
