package health

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"C"
	"github.com/dsnezhkov/tugboat/components/common"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/logger"
)

func init() {
	name := "comp_health"
	defs.ComponentAvailable[name] = CreateComponent(name)
}

type ComponentHealth struct {
	common.Component
}

func CreateComponent(name string) *ComponentHealth {
	comp := &ComponentHealth{}
	comp.Active = true
	comp.Name = name
	comp.Options = defs.OpCmd{}
	comp.SignalChan = nil
	comp.Data = []string{}
	comp.Timeout = defs.MAX_TIMEOUT
	comp.Modules = []defs.CModule{}
	comp.Tlog = nil
	comp.Sout = strings.Builder{}
	comp.Serr = strings.Builder{}
	return comp
}

func (component *ComponentHealth) SetLogger(tlog *logger.LogManager) {
	component.Tlog = tlog
}
func (component *ComponentHealth) GetName() string {
	return component.Name
}
func (component *ComponentHealth) GetModules() []defs.CModule {
	return component.Modules
}
func (component *ComponentHealth) SetModules(modules []defs.CModule) {
	component.Modules = make([]defs.CModule, len(modules))
	copy(component.Modules, modules)
}
func (component *ComponentHealth) SetCmdOptions(op *defs.OpCmd) {
	component.Options.Directive = op.Directive
	component.Options.DirectiveOpts = op.DirectiveOpts
}
func (component *ComponentHealth) SetCmdDir(opDir string) {
	component.Options.Directive = opDir
}
func (component *ComponentHealth) SetCmdDirOpt(opDirOpts []string) {
	component.Options.DirectiveOpts = opDirOpts
}
func (component *ComponentHealth) SetActive(active bool) {
	component.Active = active
}
func (component *ComponentHealth) SetData(data []string) {
	component.Data = data
}
func (component *ComponentHealth) SetTimeout(timeout uint) {
	if timeout > 0 {
		component.Timeout = timeout
	}
}
func (component *ComponentHealth) SetSignalChan(signal chan bool) {
	component.SignalChan = signal
}

func (component *ComponentHealth) InvokeComponent(
	wg *sync.WaitGroup, msg defs.Message, handoffTo []string) {
	defer wg.Done()

	resCh := make(chan bool, 1)
	go func() {
		res := component.invoke(msg, handoffTo)
		resCh <- res
	}()

	var message string
	select {
	case res := <-resCh:
		message = fmt.Sprintf("Processing returned on its own: %t", res)
		component.Tlog.Log(component.Name, "INFO", message)
	case <-time.After(time.Duration(component.Timeout * uint(time.Second))):
		message = fmt.Sprintf("Processing out of time internally")
		component.Tlog.Log(component.Name, "INFO", message)
	case signalC := <-component.SignalChan:
		message = fmt.Sprintf("Processing received signal: %v", signalC)
		component.Tlog.Log(component.Name, "INFO", message)
		message = fmt.Sprintf("Partial output so far: \nsOut: %s\n sErr: %s", component.Sout, component.Serr)
		component.Tlog.Log(component.Name, "INFO", message)
	}

}
