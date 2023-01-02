package commands

import (
	"context"
	"fmt"
	"github.com/protolambda/zrnt/eth2"
)

type VersionCmd struct {
}

const Version = "v0.5.0"

func (c *VersionCmd) Help() string {
	return "Print ZCLI and ZRNT version"
}

func (c *VersionCmd) Run(ctx context.Context, args ...string) error {
	fmt.Printf(`
ZCLI version: %s
ZRNT version: %s
`, Version, eth2.VERSION)
	return nil
}
