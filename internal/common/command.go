package common

import (
	"errors"
	"fmt"
	"github.com/catfishlty/webhook-hub/internal/check"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"runtime"
)

type StartCommand struct {
	Port      string `arg:"-p,--port,env:WH_PORT" default:"8080" help:"http server port"`
	SecretKey string `arg:"-k,--key,env:WH_SECRET_KEY" default:"" help:"secret key for jwt"`
	Salt      string `arg:"-s,--salt,env:WH_SALT" default:"" help:"salt for password"`
	LogLevel  string `arg:"-l,--log,env:WH_LOG_LEVEL" default:"info" help:"log level(panic, fatal, error, warn, warning, info, debug, trace)"`
}

func (cmd *StartCommand) Validate() error {
	return validation.ValidateStruct(cmd,
		validation.Field(&cmd.Port, validation.Required, is.Port),
		validation.Field(&cmd.LogLevel, validation.Required, check.IsLogLevel),
		validation.Field(&cmd.SecretKey, validation.Required, validation.Length(16, 128), check.IsKey),
		validation.Field(&cmd.Salt, validation.Required, validation.Length(8, 32), check.IsKey),
	)
}

func (*StartCommand) ValidateFunc(value any) error {
	if cmd, ok := value.(*StartCommand); ok {
		return cmd.Validate()
	}
	return errors.New("invalid validate value [start]]")
}

type AdminListCommand struct {
}

type AdminResetCommand struct {
	Id       string `arg:"-i,--id" help:"user id"`
	Password string `arg:"-p,--password" help:"password"`
	Salt     string `arg:"-s,--salt" default:"" help:"salt for password"`
}

func (cmd *AdminResetCommand) Validate() error {
	return validation.ValidateStruct(cmd,
		validation.Field(&cmd.Id, validation.Required, is.UUIDv4),
		validation.Field(&cmd.Password, validation.Required, is.Alphanumeric),
		validation.Field(&cmd.Salt, validation.Required, validation.Length(8, 32), is.Alphanumeric),
	)
}

func (*AdminResetCommand) ValidateFunc(value any) error {
	if cmd, ok := value.(*AdminResetCommand); ok {
		return cmd.Validate()
	}
	return errors.New("invalid validate value [admin reset]]")
}

type AdminCommand struct {
	ListCommand  *AdminListCommand  `arg:"subcommand:list" help:"list all users"`
	ResetCommand *AdminResetCommand `arg:"subcommand:reset" help:"reset user password"`
}

func (cmd *AdminCommand) Validate() error {
	return validation.ValidateStruct(cmd,
		validation.Field(&(*cmd).ListCommand, validation.Required.When(cmd.ResetCommand == nil)),
		validation.Field(&(*cmd).ResetCommand, validation.Required.When(cmd.ListCommand == nil), validation.By(cmd.ResetCommand.ValidateFunc)),
	)
}

func (*AdminCommand) ValidateFunc(value any) error {
	if cmd, ok := value.(*AdminCommand); ok {
		return cmd.Validate()
	}
	return errors.New("invalid validate value [admin]]")

}

type Command struct {
	StartCommand *StartCommand `arg:"subcommand:start" help:"start application"`
	AdminCommand *AdminCommand `arg:"subcommand:admin" help:"manage admin user"`
	DBType       string        `arg:"-t,--db-type,env:WH_DB_TYPE" default:"sqlite" help:"database type(sqlite, mysql, postgres, sqlserver), sqlite as default"`
	DBPath       string        `arg:"-d,--db,env:WH_DB_PATH" default:"./data.db" help:"database path"`
}

func (cmd *Command) Validate() error {
	return validation.ValidateStruct(cmd,
		validation.Field(&(*cmd).StartCommand, validation.Required.When(cmd.AdminCommand == nil),
			validation.When(cmd.StartCommand != nil, validation.By(cmd.StartCommand.ValidateFunc)).
				Else(validation.Nil)),
		validation.Field(&(*cmd).AdminCommand, validation.Required.When(cmd.StartCommand == nil),
			validation.When(cmd.AdminCommand != nil, validation.By(cmd.AdminCommand.ValidateFunc)).
				Else(validation.Nil)),
		validation.Field(&cmd.DBType, validation.Required, check.IsDbType),
	)
}

func (*Command) ValidateFunc(value any) error {
	if cmd, ok := value.(*Command); ok {
		return cmd.Validate()
	}
	return errors.New("invalid validate value [sub: start]]")
}

func (*Command) Version() string {
	return fmt.Sprintf("%s %s for %s", Name, Version, runtime.GOOS)
}
func (*Command) Epilogue() string {
	return "For more information visit https://github.com/catfishlty/webhook-hub"
}
