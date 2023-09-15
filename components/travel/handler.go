package travel

import (
	"log"

	"github.com/dsnezhkov/tugboat/action"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/util"
)

func (component *ComponentTravel) invoke(msg defs.Message, handoffTo []string) bool {

	var err error
	var cmdOp []string

	if component.Options.Directive != "" {
		cmdOp = append(cmdOp, component.Options.Directive)
	}
	if component.Options.DirectiveOpts != nil {
		cmdOp = append(cmdOp, component.Options.DirectiveOpts...)
	}

	// Add data (additional options perhaps) to the command
	log.Printf("Component::Travel: Adding data for execution")

	if util.CheckMessage(&msg) == true {
		cmdOp = append(cmdOp, msg.Data...)
	}

	err = util.CmdExec(&component.Sout, &component.Serr, cmdOp...)
	if err != nil {
		log.Printf("Component::Travel: Error in exec: %v, %s", err, component.Serr.String())
		log.Printf("Component::Travel: Output so far: \n %s\n", component.Sout.String())
		log.Printf("Component::Travel: Handing off empty message to chain\n")
		log.Printf("Component::Travel: Handing off to (%s) chain\n", handoffTo)

		for _, nextHandoff := range handoffTo {
			log.Printf("Component::Travel: Handing off to (%s) chain\n", nextHandoff)
			action.Handoff(nextHandoff,
				defs.Message{})
		}
	} else {
		log.Printf("Component::Travel After Execution Results:"+
			"\nStdErr: %s StdOut: %s \n", component.Serr.String(), component.Sout.String())

		log.Printf("Component::Travel: Handing off to (%s) chain\n", handoffTo)

		for _, nextHandoff := range handoffTo {
			log.Printf("Component::Travel: Handing off to (%s) chain\n", nextHandoff)
			action.Handoff(nextHandoff,
				defs.Message{Data: []string{"Z:"}})
		}
	}

	return true
}
