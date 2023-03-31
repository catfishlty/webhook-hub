package exp

import (
	"github.com/alexflint/go-arg"
	"github.com/catfishlty/webhooks-hub/internal/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

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
		panic(types.CommonError{
			Code: http.StatusBadRequest,
			Msg:  "bind json failed",
		})
	}
}

func HandleRequestInvalid(err error) {
	if err != nil {
		log.Debugf("request validate failed: %v", err)
		panic(types.CommonError{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
	}
}

func HandleDB(err error, msg string) {
	if err != nil {
		log.Errorf("%s : %v", msg, err)
		panic(types.CommonError{
			Code: http.StatusInternalServerError,
			Msg:  msg,
		})
	}
}
