package repo

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type State int

const (
	CreatingDevice State = 1

	CreatingCommandDevice State = 2
	CreatingCommandVoice  State = 3
	CreatingCommandAction State = 4
	CreatingCommandColor  State = 5
)

type Repo struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateUser(ctx context.Context, tgID int64) error {
	sql := `
		insert into users (tg_id) values ($1)
	`
	_, err := r.db.Exec(ctx, sql, tgID)
	return err
}

func (r *Repo) GetUser(ctx context.Context, tgID int64) (User, error) {
	var (
		user User
		sql  = `
			select * from users where tg_id = $1
		`
	)
	if err := pgxscan.Get(ctx, r.db, &user, sql, tgID); err != nil {
		if pgxscan.NotFound(err) {
			if err = r.CreateUser(ctx, tgID); err != nil {
				return User{}, err
			}
			return r.GetUser(ctx, tgID)
		}
		return User{}, err
	}
	return user, nil
}

func (r *Repo) ChangeState(ctx context.Context, tgID int64, state State) error {
	sql := `
		update users set state = $1 where tg_id = $2
	`
	_, err := r.db.Exec(ctx, sql, state, tgID)
	return err
}

func (r *Repo) GetCommands(ctx context.Context, tgID int64) ([]Command, error) {
	var (
		commands []Command
		sql      = `
			select * from commands where tg_id = $1 and done is true
		`
	)

	if err := pgxscan.Select(ctx, r.db, &commands, sql, tgID); err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *Repo) GetCommandList(ctx context.Context) ([]Command, error) {
	var (
		commands []Command
		sql      = `select * from command_list;`
	)

	if err := pgxscan.Select(ctx, r.db, &commands, sql); err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *Repo) CreateDevice(ctx context.Context, tgID int64, device string) error {
	sql := `
		insert into devices (tg_id, device_id) values ($1, $2)
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
		insert into commands (tg_id, device_id) values ($1, $2)
	`
	_, err := r.db.Exec(ctx, sql, tgID, device)
	return err
}

func (r *Repo) CreateCommandVoice(ctx context.Context, tgID int64, command string) error {
	sql := `
		update commands set command = $1 where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, command)
	return err
}

func (r *Repo) CreateCommandAction(ctx context.Context, tgID int64, action string) error {
	sql := `
		update commands set action = $1 where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, action)
	return err
}

func (r *Repo) CreateCommandColor(ctx context.Context, tgID int64, color string) error {
	sql := `
		update commands set color = $1 where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, color)
	return err
}

func (r *Repo) CreateCommandDone(ctx context.Context, tgID int64) error {
	sql := `
		update commands set done = true where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID)
	return err
}
