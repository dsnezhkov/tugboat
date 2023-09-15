package common

import (
	"strings"

	"tugboat/defs"
	"tugboat/logger"
)

type Component struct {
	Name    string
	Active  bool
	Options defs.OpCmd
	SignalChan chan bool
	Sout       strings.Builder
	Serr       strings.Builder
	Data       []string
	Timeout    uint
	Modules    []defs.CModule
	Tlog       *logger.LogManager
}
