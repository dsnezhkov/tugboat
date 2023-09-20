package workflow

import (
	"fmt"
	"github.com/dsnezhkov/tugboat-components/travel2"
	"github.com/dsnezhkov/tugboat/components/health"
	"github.com/dsnezhkov/tugboat/components/tech"
	"github.com/dsnezhkov/tugboat/components/travel"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/util"
	"log"
	"sync"
)

type WorkFlow struct{}

var Workflow *WorkFlow

func CreateWorkFlow() *WorkFlow {
	w := &WorkFlow{}
	return w
}

// Define action broker for each channel

func (workflow *WorkFlow) Broker(
	name string,
	ch <-chan defs.Message,
	opFunc func(wg *sync.WaitGroup, message defs.Message, handoffTo []string),
	handoffTo []string) {

	for msg := range ch {
		log.Printf("Main: Broker [%s], handling message: %s\n", name, msg)
		defs.Wg.Add(1)

		msg := msg
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Main: Execution encountered fault and recovered from: %s", err)
				}
			}()
			opFunc(&defs.Wg, msg, handoffTo)
		}()
	}

	log.Printf("Main: Broker [%s] is done\n", name)
}

type ComponentHelper struct{}

func (ch *ComponentHelper) ProcessModules(compConf *defs.Component, compObj *defs.IComponent) {

	var modules []defs.CModule
	fmt.Printf("Modules: ............... %+v\n", modules)
	if compConf.Modules != nil {
		for _, cm := range compConf.Modules {
			log.Printf("Workflow::ComponentHelper: Module %s\n", cm.Name)

			// Create CModule instance
			var module defs.CModule
			// Assign module name
			module.Name = cm.Name

			if cm.Loads != nil {
				module.CLoads = make([]defs.CLoad, len(cm.Loads))
				for i, cl := range cm.Loads {

					log.Printf("Workflow::ComponentHelper: Module:Loads:Name %s\n", cl.Name)

					module.CLoads[i].Name = cl.Name
					log.Printf("Workflow::ComponentHelper: Module:Loads:DataType %s\n", cl.DataType)
					module.CLoads[i].DataType = cl.DataType

					if cl.Identifier != "" {
						log.Printf("Workflow::ComponentHelper: Module:Loads:Identifier %s\n", cl.Identifier)
						module.CLoads[i].Identifier = cl.Identifier
					}

					log.Printf(
						"Workflow::ComponentHelper: Module:Loads:Native Format %s\n", cl.NativeFormat)
					module.CLoads[i].NativeFormat = cl.NativeFormat
					log.Printf(
						"Workflow::ComponentHelper: Module:Loads:Source %s\n", cl.Source)
					module.CLoads[i].Source = cl.Source

					if cl.Encapsulation != nil {

						log.Printf(
							"Workflow::ComponentHelper: Module:Loads:Encapsulation:Order %s\n", cl.Encapsulation.Order)
						module.CLoads[i].CEncapsulation = new(defs.CLoadEncapsulation)
						module.CLoads[i].CEncapsulation.Order = cl.Encapsulation.Order
						if cl.Encapsulation.StoredEncoded != nil {
							module.CLoads[i].CEncapsulation.CStoredEncoded = new(defs.CLoadEncapsulationEncoded)

							log.Printf(
								"Workflow::ComponentHelper: Module:Loads:Encapsulation:StoredEncoded:Variant %s\n", cl.Encapsulation.StoredEncoded.Variant)
							module.CLoads[i].CEncapsulation.CStoredEncoded.Variant = cl.Encapsulation.StoredEncoded.Variant

						}
						if cl.Encapsulation.StoredEncrypted != nil {

							module.CLoads[i].CEncapsulation.CStoredEncrypted = new(defs.CLoadEncapsulationEncrypted)
							log.Printf(
								"Workflow::ComponentHelper: Module:Loads:Encapsulation:StoredEncrypted:{Algorithm: %s <-> Key: %s}\n", cl.Encapsulation.StoredEncrypted.Algorithm, cl.Encapsulation.StoredEncrypted.Key)

							module.CLoads[i].CEncapsulation.CStoredEncrypted.Algorithm = cl.Encapsulation.StoredEncrypted.Algorithm
							module.CLoads[i].CEncapsulation.CStoredEncrypted.Key = cl.Encapsulation.StoredEncrypted.Key
						}
					}

				}
			}

			modules = append(modules, module)
		}

		(*compObj).SetModules(modules)
	}
}

func (workflow *WorkFlow) unwrapComponent(o *defs.IComponent, c *defs.Component) (bool, string) {
	var ok bool
	var message string

	switch {
	case c.Name == "comp_tech":
		*o = defs.ComponentAvailable[(*c).Name].(*tech.ComponentTech)
		ok = true
	case c.Name == "comp_travel":
		*o = defs.ComponentAvailable[(*c).Name].(*travel.ComponentTravel)
		ok = true
	case c.Name == "comp_travel2":
		*o = defs.ComponentAvailable[(*c).Name].(*travel2.ComponentTravel2)
		ok = true
	case c.Name == "comp_health":
		*o = defs.ComponentAvailable[(*c).Name].(*health.ComponentHealth)
		ok = true
	default:
		ok = false
		message = fmt.Sprintf("Unknown type of object: %+v\n", *o)
	}

	return ok, message
}
func (workflow *WorkFlow) Setup() {

	// Message plane for component
	defs.Chan2component = map[string]<-chan defs.Message{}
	// Signal plane for component
	defs.ChanSignal2component = make(map[string]chan bool)

	// Source in components
	var comp, vector string
	var nextVectors []string

	// Process components
	// c -> Config options structure
	// o -> Specialized object

	for _, c := range defs.FlowConf.ComponentBlock.Components {

		chelper := new(ComponentHelper)
		message := fmt.Sprintf("Processing Component: %s: %+v\n", c.Name, c)

		defs.Tlog.Log("workflow", "DEBUG", message)

		var o defs.IComponent
		//ok, msg := workflow.unwrapComponent(&o, &c)
		//if !ok {
		//	defs.Tlog.Log("workflow", "ERROR", msg)
		//	continue
		//}

		o = defs.ComponentAvailable[c.Name]

		// Process component
		// Command directives and options

		defs.Tlog.Log("workflow", "DEBUG", "Setting SetCmdOptions ")
		o.SetCmdOptions(
			&defs.OpCmd{
				Directive:     c.Directive,
				DirectiveOpts: c.DirectiveOpts,
			})

		// Modules

		defs.Tlog.Log("workflow", "DEBUG", "Processing Modules")
		chelper.ProcessModules(&c, &o)

		// What vectors a component needs to hand off its processing to
		comp, vector = util.FindVectorForComponent(c.Name)
		nextVectors = util.FindNextVectorLinkInFlow(defs.FlowConf, vector)

		// Add default data bag
		defs.ComponentAvailable[c.Name].SetData(c.Data)

		// TODO: check conversions
		// Add timeout
		// This allows timing out runaway components (e.g. search)
		if c.Timeout != 0 {
			o.SetTimeout(uint(c.Timeout))
		}
		// log.Printf("dump: %+v\n", defs.ComponentAvailable[c.name])

		// TODO: figure out if setlog shoud be used here or set in instance fields
		o.SetLogger(defs.Tlog)
		// Set active flag
		o.SetActive(c.Active)

		if !c.Active {
			continue
		}
		message = fmt.Sprintf("Subscribing component %s to vector %+v\n", comp, vector)

		defs.Tlog.Log("workflow", "DEBUG", message)

		// Register subscriptions for components to vector
		// This allows a component to subscribe to a vector handoff
		defs.Chan2component[comp] = defs.Ps.Subscribe(vector)

		// Create signal plane and register it for component
		// This allows a component to receive an out-of-band signal
		ch := make(chan bool, 1)
		o.SetSignalChan(ch)
		defs.ChanSignal2component[c.Name] = ch

		// Register component in the chain
		// This creates a chain of events in the workflow
		message = fmt.Sprintf("Registering component: %+v channel: %+v entry: %+v, vectors: %+v ", c.Name, defs.Chan2component[c.Name],
			o.InvokeComponent, vector)

		defs.Tlog.Log("workflow", "DEBUG", message)
		// Setup broker for each component channel
		go workflow.Broker(
			o.GetName(),
			defs.Chan2component[c.Name],
			o.InvokeComponent,
			nextVectors)
	}
}

func (workflow *WorkFlow) GetWHead() string {
	return defs.FlowConf.FlowLinkBlock.Head
}
