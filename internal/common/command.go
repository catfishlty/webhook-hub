package common

import (
	"fmt"
	"runtime"
)

type StartCommand struct {
	Port      int    `arg:"-p,--port" default:"8080" help:"http server port"`
	SecretKey string `arg:"-k,--key" default:"" help:"secret key for jwt"`
	LogLevel  string `arg:"-l,--log" default:"info" help:"log level(debug, info, warn, error, fatal, panic)"`
}

type AdminListCommand struct {
}

type AdminResetCommand struct {
	Id       string `arg:"-i,--id" help:"user id"`
	Password string `arg:"-p,--password" help:"password"`
}

type AdminCommand struct {
	ListCommand  *AdminListCommand  `arg:"subcommand:list" help:"list all users"`
	ResetCommand *AdminResetCommand `arg:"subcommand:reset" help:"reset user password"`
}

type Command struct {
	StartCommand *StartCommand `arg:"subcommand:start" help:"start application"`
	AdminCommand *AdminCommand `arg:"subcommand:admin" help:"manage admin user"`
	DBType       string        `arg:"-t,--db-type" default:"sqlite" help:"database type(sqlite, mysql)"`
	DBPath       string        `arg:"-d,--db" default:"./data.db" help:"database path"`
}

func (*Command) Version() string {
	return fmt.Sprintf("%s %s for %s", Name, Version, runtime.GOOS)
}
func (*Command) Epilogue() string {
	return "For more information visit https://github.com/catfishlty/webhook-hub"
}
