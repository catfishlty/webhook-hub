package exp

import (
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type CommonError struct {
	Code          int
	Message       string
	Err           error
	IsSystemError bool
}

func HandleCmd(p *arg.Parser, err error) {
	if err != nil {
		p.Fail(err.Error())
		os.Exit(1)
	}
}

func HandleCmdWithMsg(p *arg.Parser, err error, msg string) {
	if err != nil {
		p.Fail(msg)
		log.Errorf("%s : %v", msg, err)
		os.Exit(1)
	}
}

func HandleCmdCondition(p *arg.Parser, cond bool, msg string) {
	if cond {
		p.Fail(msg)
		os.Exit(1)
	}
}

func HandleBindJSON(err error) {
	if err != nil {
		log.Warnf("bind json failed: %v", err)
		panic(&CommonError{
			Code:    http.StatusBadRequest,
			Message: "bind json failed",
		})
	}
}

func HandleRequestInvalid(err error) {
	if err != nil {
		log.Debugf("request validate failed: %v", err)
		panic(&CommonError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
}

func HandleDB(err error, msg string) {
	if err != nil {
		log.Errorf("%s : %v", msg, err)
		panic(&CommonError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		})
	}
}
