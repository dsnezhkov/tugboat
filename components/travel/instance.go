package travel

import (
	"fmt"
	"sync"
	"time"

	"github.com/dsnezhkov/tugboat/components/common"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/logger"
)

func init() {
	name := "comp_travel"
	defs.ComponentAvailable[name] = CreateComponent(name)
}

type ComponentTravel struct {
	common.Component
}

func CreateComponent(name string) *ComponentTravel {
	comp := &ComponentTravel{}
	comp.Active = true
	comp.Name = name
	comp.Options = defs.OpCmd{}
	comp.SignalChan = nil
	comp.Data = []string{}
	comp.Timeout = defs.MAX_TIMEOUT
	comp.Modules = []defs.CModule{}
	comp.Tlog = nil
	return comp
}

func (component *ComponentTravel) SetLogger(tlog *logger.LogManager) {
	component.Tlog = tlog
}
func (component *ComponentTravel) GetName() string {
	return component.Name
}

func (component *ComponentTravel) GetModules() []defs.CModule {
	return component.Modules
}
func (component *ComponentTravel) SetModules(modules []defs.CModule) {
	component.Modules = make([]defs.CModule, len(modules))
	copy(component.Modules, modules)
}
func (component *ComponentTravel) SetCmdOptions(op *defs.OpCmd) {
	component.Options.Directive = op.Directive
	component.Options.DirectiveOpts = op.DirectiveOpts
}
func (component *ComponentTravel) SetCmdDir(opDir string) {
	component.Options.Directive = opDir
}
func (component *ComponentTravel) SetCmdDirOpt(opDirOpts []string) {
	component.Options.DirectiveOpts = opDirOpts
}
func (component *ComponentTravel) SetActive(active bool) {
	component.Active = active
}
func (component *ComponentTravel) SetSignalChan(signal chan bool) {
	component.SignalChan = signal
}
func (component *ComponentTravel) SetData(data []string) {
	component.Data = data
}
func (component *ComponentTravel) SetTimeout(timeout uint) {
	if timeout > 0 {
		component.Timeout = timeout
	}
}
func (component *ComponentTravel) InvokeComponent(
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
		message = fmt.Sprintf(
			"Partial output so far:\nsOut:\n%s\nsErr:\n%s\n", component.Sout.String(), component.Serr.String())
		component.Tlog.Log(component.Name, "INFO", message)

	case signalC := <-component.SignalChan:
		message = fmt.Sprintf("Processing received signal: %v", signalC)
		component.Tlog.Log(component.Name, "INFO", message)
		message = fmt.Sprintf("Partial output so far: \nsOut: %s\n sErr: %s", component.Sout.String(), component.Serr.String())
		component.Tlog.Log(component.Name, "INFO", message)
	}
}
