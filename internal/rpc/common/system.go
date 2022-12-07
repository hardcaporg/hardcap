package common

import "github.com/hardcaporg/hardcap/internal/model"

type SystemListArgs struct {
    Limit, Offset int
}

type SystemListReply struct {
    List []*model.System
}
