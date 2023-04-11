package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/catfishlty/webhooks-hub/exp"
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/data"
	"github.com/catfishlty/webhooks-hub/internal/hub"
	"github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"os"
)

var (
	g errgroup.Group
)

func main() {
	var args common.Command
	p := arg.MustParse(&args)
	exp.HandleCmd(p, args.Validate())
	dbInstance, err := data.GetDatabase(args.DBType, args.DBPath)
	exp.HandleCmdWithMsg(p, err, "failed to create database instance")
	orm, err := gorm.Open(dbInstance, &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	exp.HandleCmdWithMsg(p, err, "failed to connect database")
	switch {
	case args.StartCommand != nil:
		log.Debugf("command: Start")
		webhookHub := hub.NewHub(orm, args.StartCommand.SecretKey, args.StartCommand.Salt)
		g.Go(func() error {
			return webhookHub.Server(args.StartCommand.Port).ListenAndServe()
		})
		if err := g.Wait(); err != nil {
			log.Errorf("server start failed: %v", err)
			p.Fail("server start failed")
		}
		return
	case args.AdminCommand != nil:
		switch {
		case args.AdminCommand.ListCommand != nil:
			log.Debugf("command: Admin List")
			db := data.NewDB(orm, "")
			users, total, err := db.GetUserList(1, common.PageSizeMax)
			exp.HandleCmdWithMsg(p, err, "failed to get user list")
			fmt.Printf("total users: %d\n", total)
			for idx, user := range users {
				fmt.Printf("%3d - user: id='%s', name='%s'", idx, user.UID, user.Username)
			}
		case args.AdminCommand.ResetCommand != nil:
			log.Debugf("command: Admin Reset")
			exp.HandleCmdCondition(p, args.AdminCommand.ResetCommand.Id == "", "user id is empty")
			exp.HandleCmdWithMsg(p, err, "failed to get user list")
			exp.HandleCmdCondition(p, args.AdminCommand.ResetCommand.Password == "", "password is empty")
			db := data.NewDB(orm, args.AdminCommand.ResetCommand.Salt)
			err = db.UpdatePasswordAdmin(args.AdminCommand.ResetCommand.Id, args.AdminCommand.ResetCommand.Password)
			if err != nil {
				p.Fail(fmt.Sprintf("failed to reset user '%s'", args.AdminCommand.ResetCommand.Id))
				os.Exit(1)
			}
			fmt.Printf("reset user '%s' password success\n", args.AdminCommand.ResetCommand.Id)
		}
	default:
		p.Fail("command not found")
	}
}
