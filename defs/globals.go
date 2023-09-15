package defs

import (
	"sync"

	"github.com/dsnezhkov/tugboat/comms"
	"github.com/dsnezhkov/tugboat/fs"
	"github.com/dsnezhkov/tugboat/logger"
	"github.com/dsnezhkov/tugboat/payloads"
)

const MAX_TIMEOUT = 86400 // 24hrs

var Ps *Pubsub
var FlowConf *Config

var ComponentAvailable = map[string]IComponent{}
var VectorsAvailable = map[string][]string{}

var Wg sync.WaitGroup

var FsManager *fs.FSManager
var PayManager *payloads.PayloadManager
var CommsManager *comms.CommsManager
var Tlog *logger.LogManager

var Chan2component map[string]<-chan Message

// Signal plane for component
var ChanSignal2component map[string]chan bool
