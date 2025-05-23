package repo

type Command struct {
	User    int64  `db:"tg_id" json:"user"`
	Device  string `db:"device_id" json:"device"`
	Command string `db:"command" json:"command"`
	Action  string `db:"action" json:"action"`
	Color   string `db:"color" json:"color"`
}
