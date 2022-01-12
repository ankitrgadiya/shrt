package datastore

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"argc.in/shrt/internal/model"
	"github.com/pkg/errors"
)

func NewSQLiteStore(path string) (RouteStore, error) {
	dsn := fmt.Sprintf("file:%s?_journal=WAL", path)

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "opening %s", path)
	}

	i := &sqliteImpl{conn: conn}

	if err := i.Initialize(); err != nil {
		return nil, err
	}

	return i, nil
}

type sqliteImpl struct {
	conn *sql.DB
}

func (i *sqliteImpl) Initialize() error {
	ctx := context.Background()

	query := `CREATE TABLE IF NOT EXISTS links (
                slug TEXT PRIMARY KEY,
                url TEXT NOT NULL
              )`

	stmt, err := i.conn.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "preparing statement")
	}

	if _, err := stmt.Exec(); err != nil {
		return errors.Wrap(err, "creating table")
	}

	return nil
}

func (i *sqliteImpl) Close() error {
	return i.conn.Close()
}

func (i *sqliteImpl) Save(ctx context.Context, r *model.Route) error {
	query := `INSERT INTO links (slug, url)
			  VALUES(:1, :2)
              ON CONFLICT(slug)
              DO UPDATE SET url=excluded.url`

	if _, err := i.conn.ExecContext(ctx, query, r.Slug, r.URL); err != nil {
		return err
	}

	return nil
}

func (i *sqliteImpl) Query(ctx context.Context, r *model.Route) error {
	query := `SELECT url FROM links WHERE slug = ?`

	if err := i.conn.QueryRowContext(ctx, query, r.Slug).Scan(&r.URL); err != nil {
		return err
	}

	return nil
}

func (i *sqliteImpl) QueryAll(ctx context.Context) (routes []model.Route, err error) {
	query := `SELECT slug, url FROM links`

	rows, err := i.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r model.Route
		if err := rows.Scan(&r.Slug, &r.URL); err != nil {
			return nil, err
		}

		routes = append(routes, r)
	}

	return routes, nil
}

func (i *sqliteImpl) Delete(ctx context.Context, r *model.Route) error {
	query := `DELETE FROM links WHERE slug = ?`

	if _, err := i.conn.ExecContext(ctx, query, r.Slug); err != nil {
		return err
	}

	return nil
}
