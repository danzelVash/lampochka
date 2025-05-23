package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type State int

const (
	CreatingDevice State = 1

	CreatingCommandDevice State = 2
	CreatingCommandVoice  State = 2
	CreatingCommandAction State = 2
	CreatingCommandColor  State = 2
)

type Repo struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) ChangeState(ctx context.Context, tgID int64, state State) error {
	sql := `
		update users set state = $1 where id = $2
	`
	_, err := r.db.Exec(ctx, sql, state, tgID)
	return err
}

func (r *Repo) CreateDevice(ctx context.Context, tgID int64, device string) error {
	sql := `
		insert into devices (user, device_id) values ($1, $2)
	`
	_, err := r.db.Exec(ctx, sql, tgID, device)
	return err
}

// девайс
// гс
// че делает
// если врубает цвет, то выбор цвета

func (r *Repo) CreateCommandDevice(ctx context.Context, tgID int64, device string) error {
	sql := `
		insert into commands (user, device_id) values ($1, $2)
	`
	_, err := r.db.Exec(ctx, sql, tgID, device)
	return err
}

func (r *Repo) CreateCommandVoice(ctx context.Context, tgID int64, command string) error {
	sql := `
		update commands set command = $1 where user = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, command)
	return err
}

func (r *Repo) CreateCommandAction(ctx context.Context, tgID int64, action string) error {
	sql := `
		update commands set action = $1 where user = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, action)
	return err
}

func (r *Repo) CreateCommandColor(ctx context.Context, tgID int64, color string) error {
	sql := `
		update commands set color = $1 where user = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, color)
	return err
}

func (r *Repo) CreateCommandDone(ctx context.Context, tgID int64) error {
	sql := `
		update commands set done = true where user = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID)
	return err
}
