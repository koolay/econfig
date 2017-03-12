package config

import (
	"io/ioutil"
	"log"
	"os"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	logHandle = os.Stdout
)

func NewLogger(gf *GlobalFlag) *jww.Notepad {
	var logger *jww.Notepad

	var err error
	if gf.LogFile != "" {
		logHandle, err = os.OpenFile(gf.LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}

	if gf.Verbose || viper.GetBool("flags.verbose") {
		logger = jww.NewNotepad(jww.LevelTrace, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
	} else {
		logger = jww.NewNotepad(jww.LevelError, jww.LevelTrace, logHandle, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	return logger
}
