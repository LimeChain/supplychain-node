package handler

import (
	"github.com/Limechain/pwc-bat-node/app/business/messages"
	"github.com/Limechain/pwc-bat-node/app/interfaces/common"
)

type BusinessMessageParser interface {
	Parse(msg *common.Message) (*messages.BusinessMessage, error)
}

type BusinessLogicHandler interface {
	Handle(msg *common.Message) error
}
