package config

import (
	"io/ioutil"
	"log"
	"os"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

func NewLogger(gf *GlobalFlag) *jww.Notepad {
	var logger *jww.Notepad
	if gf.Verbose || viper.GetBool("flags.verbose") {
		logger = jww.NewNotepad(jww.LevelTrace, jww.LevelTrace, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	} else {
		logger = jww.NewNotepad(jww.LevelError, jww.LevelTrace, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}
	return logger
}
