package on_off

type OnOff struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID      string   `json:"id"`
	Actions []Action `json:"actions"`
}

type Action struct {
	Type  string `json:"type"`
	State State  `json:"state"`
}

type State struct {
	Instance string `json:"instance"`
	Value    bool   `json:"value"`
}

func New(id string, value bool) OnOff {
	return OnOff{
		Devices: []Device{
			{
				ID: id,
				Actions: []Action{
					{
						Type: "devices.capabilities.on_off",
						State: State{
							Instance: "on",
							Value:    value,
						},
					},
				},
			},
		},
	}
}
