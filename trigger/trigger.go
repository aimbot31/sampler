package trigger

import (
	"fmt"
	"github.com/sqshq/sampler/asset"
	"github.com/sqshq/sampler/config"
	"os"
	"os/exec"
)

const (
	TrueIndicator = "1"
	BellCharacter = "\a"
)

type Trigger struct {
	title     string
	condition string
	actions   Actions
	data      map[string]Data
}

type Actions struct {
	terminalBell bool
	sound        bool
	visual       bool
	script       *string
}

type Data struct {
	previousValue interface{}
	currentValue  interface{}
}

func NewTriggers(configs []config.TriggerConfig) []Trigger {

	triggers := make([]Trigger, len(configs))

	for _, c := range configs {
		triggers = append(triggers, NewTrigger(c))
	}

	return triggers
}

func NewTrigger(config config.TriggerConfig) Trigger {
	return Trigger{
		title:     config.Title,
		condition: config.Condition,
		data:      make(map[string]Data),
		actions: Actions{
			terminalBell: *config.Actions.TerminalBell,
			sound:        *config.Actions.Sound,
			visual:       *config.Actions.Visual,
			script:       config.Actions.Script,
		},
	}
}

func (t Trigger) Execute(label string, value interface{}) {

	if data, ok := t.data[label]; ok {
		data.previousValue = data.currentValue
		data.currentValue = value
	} else {
		t.data[label] = Data{previousValue: nil, currentValue: value}
	}

	t.evaluate(label, t.data[label])
}

func (t Trigger) evaluate(label string, data Data) {

	output, err := runScript(t.condition, label, data)

	if err != nil {
		//println(err) // TODO visual notification
	}

	if string(output) != TrueIndicator {
		return
	}

	if t.actions.terminalBell {
		fmt.Print(BellCharacter)
	}

	if t.actions.sound {
		_ = asset.Beep()
	}

	if t.actions.visual {
		// TODO visual notification
	}

	if t.actions.script != nil {
		_, _ = runScript(*t.actions.script, label, data)
	}
}

func runScript(script, label string, data Data) ([]byte, error) {
	cmd := exec.Command("sh", "-c", script)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("prev=%v", data.previousValue),
		fmt.Sprintf("cur=%v", data.currentValue),
		fmt.Sprintf("label=%v", label))

	return cmd.Output()
}