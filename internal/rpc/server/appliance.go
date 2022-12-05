package server

import (
	"github.com/hardcaporg/hardcap/internal/db"
	"github.com/hardcaporg/hardcap/internal/rpc/common"
)

type Appliance struct{}

func (s *Appliance) Register(args common.ApplianceRegisterArgs, reply *common.ApplianceRegisterReply) error {
    err := db.Appliance.Create(args.Appliance)
    if err != nil {
        return err
    }

    reply.Appliance = args.Appliance
    return nil
}

func (s *Appliance) List(args common.ApplianceListArgs, reply *common.ApplianceListReply) error {
	var err error
	if args.Limit == 0 {
		args.Limit = 100
	}

	reply.List, _, err = db.Appliance.FindByPage(args.Offset, args.Limit)
	if err != nil {
		return err
	}

	return nil
}
