package common

import "github.com/hardcaporg/hardcap/internal/model"

type RegistrationListArgs struct {
    Limit, Offset int
}

type RegistrationListReply struct {
    List []*model.Registration
}