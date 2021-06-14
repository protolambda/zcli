package commands

type TreeCmd struct{}

func (c *TreeCmd) Help() string {
	return "Check ssz payload"
}

func (c *TreeCmd) Cmd(route string) (cmd interface{}, err error) {
	// TODO: per type, return a cmd
	return nil, err
}
