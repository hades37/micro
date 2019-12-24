package conf

import (
	"github.com/go-ini/ini"
)

var  Conf *ini.File
var err error

func init(){
	Conf,err=ini.Load("./conf.ini")
}