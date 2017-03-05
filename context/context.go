package context

import (
	"github.com/koolay/econfig/config"
	jww "github.com/spf13/jwalterweatherman"
)

// FlagsType flags
type FlagsType struct {
	Global *config.GlobalFlag
	Serve  *config.ServeFlag
}

// Flags collection of flag
var Flags FlagsType

var Logger *jww.Notepad
