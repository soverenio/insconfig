package utils

import (
	"strconv"
)

type Path string

func (p Path) AppendStructKey(key string) Path {
	if p != "." {
		p += "."
	}
	return p + Path(key)
}

func (p Path) AppendMapKey(key string) Path {
	return p + Path("['"+key+"']")
}

func (p Path) AppendArrayIdx(idx int) Path {
	return p + Path("["+strconv.Itoa(idx)+"]")
}
