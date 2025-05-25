package repo

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type State int

const (
	CreatingDevice State = 1

	CreatingCommandDevice State = 2
	CreatingCommandAction State = 3
	CreatingCommandText   State = 4

	DeleteCommand State = 5
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
	fmt.Printf("[REPO] CreateUser - User:%d", tgID)

	sql := `
		insert into users (tg_id) values ($1)
		on conflict (tg_id) do nothing
	`
	_, err := r.db.Exec(ctx, sql, tgID)
	return err
}

func (r *Repo) GetUser(ctx context.Context, tgID int64) (User, error) {
	fmt.Printf("[REPO] GetUser - User:%d", tgID)

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
	fmt.Printf("[REPO] ChangeState - User:%d\n", tgID)

	sql := `
		update users set state = $1 where tg_id = $2
	`
	_, err := r.db.Exec(ctx, sql, state, tgID)
	return err
}

func (r *Repo) GetCommands(ctx context.Context, tgID int64) ([]Command, error) {
	fmt.Printf("[REPO] GetCommands - User:%d", tgID)

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
	fmt.Println("[REPO] GetCommandList")

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
	fmt.Printf("[REPO] CreateDevice - User:%d", tgID)

	sql := `
		insert into devices (tg_id, device_id) values ($1, $2)
	`
	_, err := r.db.Exec(ctx, sql, tgID, device)
	return err
}

func (r *Repo) CreateCommandDevice(ctx context.Context, tgID int64, device string) error {
	fmt.Printf("[REPO] CreateCommandDevice - User:%d", tgID)

	sql := `
		insert into commands (tg_id, device_id) values ($1, $2)
	`
	_, err := r.db.Exec(ctx, sql, tgID, device)
	return err
}

func (r *Repo) CreateCommandText(ctx context.Context, tgID int64, command string) error {
	fmt.Printf("[REPO] CreateCommandText - User:%d", tgID)

	sql := `
		update commands set command = $1, done = true where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, command, tgID)
	return err
}

func (r *Repo) DeleteCommand(ctx context.Context, tgID int64, command string) error {
	fmt.Printf("[REPO] DeleteCommand - User:%d, Command:%s", tgID, command)

	sql := `
		delete from commands where command = $1 and tg_id = $2
	`
	_, err := r.db.Exec(ctx, sql, command, tgID)
	return err
}

func (r *Repo) CreateCommandAction(ctx context.Context, tgID int64, action string) error {
	fmt.Printf("[REPO] CreateCommandAction - User:%d", tgID)

	sql := `
		update commands set action = $1 where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, action, tgID)
	return err
}

func (r *Repo) CreateCommandColor(ctx context.Context, tgID int64, color string) error {
	fmt.Printf("[REPO] CreateCommandColor - User:%d", tgID)

	sql := `
		update commands set color = $1 where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID, color)
	return err
}

func (r *Repo) CreateCommandDone(ctx context.Context, tgID int64) error {
	fmt.Printf("[REPO] CreateCommandDone - User:%d", tgID)

	sql := `
		update commands set done = true where tg_id = $2 and done is false
	`
	_, err := r.db.Exec(ctx, sql, tgID)
	return err
}
