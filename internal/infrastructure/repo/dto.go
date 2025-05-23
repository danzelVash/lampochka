package repo

type Command struct {
	User    int64  `db:"tg_id" json:"tg_id"`
	Device  string `db:"device_id" json:"device_id"`
	Command string `db:"command" json:"command"`
	Action  string `db:"action" json:"action"`
	Color   string `db:"color" json:"color"`
}

type User struct {
	TgID  int64 `db:"tg_id" json:"user"`
	State State `db:"state" json:"state"`
}
