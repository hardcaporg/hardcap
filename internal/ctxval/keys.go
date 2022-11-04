package ctxval

type commonKeyId int

const (
	loggerCtxKey    commonKeyId = iota
	requestIdCtxKey commonKeyId = iota
)
