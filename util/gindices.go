package util

import (
	"fmt"
	"github.com/protolambda/ztyp/tree"
	"strconv"
	"strings"
)

type GindicesFlag []tree.Gindex64

func (p *GindicesFlag) String() string {
	if p == nil {
		return ""
	}
	var buf strings.Builder
	for _, v := range *p {
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	return buf.String()
}

func (p *GindicesFlag) Set(v string) error {
	if p == nil {
		return fmt.Errorf("cannot decode gindices list into nil pointer")
	}
	parts := strings.Split(v, ",")
	*p = make([]tree.Gindex64, 0, len(parts))
	for i, v := range parts {
		s := strings.TrimSpace(v)
		if s == "" {
			continue
		}
		g, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gindex list, item %d, got: %q", i, s)
		}
		*p = append(*p, tree.Gindex64(g))
	}
	return nil
}

func (p *GindicesFlag) Type() string {
	return "comma separated list of generalized indices (in any base)"
}
