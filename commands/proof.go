package commands

type ProofCmd struct{}

func (c *ProofCmd) Help() string {
	return "Produce and verify arbitrary SSZ merkle proofs over any spec object"
}

func (c *ProofCmd) Cmd(route string) (cmd interface{}, err error) {
	// TODO: per type, return a cmd
	return nil, err
}
