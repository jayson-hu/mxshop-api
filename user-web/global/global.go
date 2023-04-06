package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/jayson-hu/mxshop-api/user-web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
)
