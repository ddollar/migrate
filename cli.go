package migrate

import (
	"fmt"

	"github.com/ddollar/stdcli"
)

type CLI struct {
	cli    *stdcli.Engine
	engine *Engine
}

func (cc *CLI) Register() {
	cc.cli.Command("up", "apply migrations", cc.up, stdcli.CommandOptions{
		Usage:    "[version]",
		Validate: stdcli.ArgsMax(1),
	})
}

func (cc *CLI) Run(args []string) error {
	if err := cc.engine.Initialize(); err != nil {
		return err
	}

	if code := cc.cli.Execute(args); code != 0 {
		return fmt.Errorf("exit %d", code)
	}

	return nil
}

func (cc *CLI) up(_ *stdcli.Context) error {
	pending, err := cc.engine.Pending()
	if err != nil {
		return err
	}

	for _, v := range pending {
		fmt.Printf("%s: ", v)

		if err := cc.engine.MigrationUp(v); err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Println("OK")
		}
	}

	return nil
}
