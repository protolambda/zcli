package commands

import (
	"context"
	"fmt"
)

type VersionCmd struct {
}

func (c *VersionCmd) Run(ctx context.Context) error {
	fmt.Println("Version") // TODO
	return nil

}
