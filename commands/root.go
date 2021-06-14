package commands

type RootCmd struct{}

func (c *RootCmd) Help() string {
	return "Compute the SSZ hash-tree-root of the spec object"
}

func (c *RootCmd) Cmd(route string) (cmd interface{}, err error) {
	// TODO: per type, return a cmd
	return nil, err
}
