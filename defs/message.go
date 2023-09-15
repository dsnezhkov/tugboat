package defs

import (
	"sync"

	"github.com/dsnezhkov/tugboat/logger"
)

type Message struct {
	Data []string
}
type OpCmd struct {
	Directive     string
	DirectiveOpts []string
}
type OpCmdFunc func(wg *sync.WaitGroup, message Message, handoffTo string)

type IComponent interface {
	InvokeComponent(
		wg *sync.WaitGroup, msg Message, handoffTo []string)
	SetLogger(logger *logger.LogManager)
	SetSignalChan(ch chan bool)
	SetCmdOptions(op *OpCmd)
	SetCmdDir(opDir string)
	SetCmdDirOpt(opDirOpts []string)
	SetData(data []string)
	SetTimeout(timeout uint)
	SetActive(active bool)
	GetName() string
	GetModules() []CModule
	SetModules([]CModule)
}

type CModule struct {
	Name   string
	CLoads []CLoad
}
type CLoad struct {
	Name           string
	Identifier     string
	DataType       string
	NativeFormat   string
	CEncapsulation *CLoadEncapsulation
	Source         string
}

type CLoadEncapsulation struct {
	CStoredEncrypted  *CLoadEncapsulationEncrypted
	CStoredCompressed *CLoadEncapsulationCompressed
	CStoredEncoded    *CLoadEncapsulationEncoded
	Order             string
}

type CLoadEncapsulationEncrypted struct {
	Algorithm string
	Key       string
}
type CLoadEncapsulationCompressed struct {
	Variant string
}
type CLoadEncapsulationEncoded struct {
	Variant string
}
