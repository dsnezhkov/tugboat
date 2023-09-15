package tech

import (
	"fmt"

	"tugboat/action"
	"tugboat/defs"
	"tugboat/util"
)


func (component *ComponentTech) invoke(msg defs.Message, handoffTo []string) bool {

	var err error
	var cmdOp []string
	var logMessage string

	if component.Options.Directive != "" {
		cmdOp = append( cmdOp, component.Options.Directive)
	}
	if component.Options.DirectiveOpts != nil {
		cmdOp = append( cmdOp, component.Options.DirectiveOpts...)
	}

	// Add data (additional options perhaps) to the command
	if util.CheckMessage(&msg) == true {
		cmdOp = append( cmdOp, msg.Data...)
	}

	err = util.CmdExec(&component.Sout, &component.Serr, cmdOp...)
	if err!=nil {
		logMessage = fmt.Sprintf("Error in execution: %v, %s", err, component.Serr.String())

		component.Tlog.Log(component.Name, "Error", logMessage )
		logMessage = fmt.Sprintf("Output collected so far: \n%s\n", component.Sout.String()  )

		logMessage = fmt.Sprintf("Handing off empty message to (%s) chain\n", handoffTo )
		component.Tlog.Log(component.Name, "Debug", logMessage )

		for _,nextHandoff := range handoffTo {
			logMessage = fmt.Sprintf("Handing off message to (%s) chain\n", nextHandoff)
			component.Tlog.Log(component.Name, "Debug", logMessage )
			action.Handoff(nextHandoff, defs.Message{})
		}
	}else{

		logMessage = fmt.Sprintf("Execution Results: \nStdErr:\n%s\nStdOut:\n%s\n",component.Serr.String(),component.Sout.String() )
		component.Tlog.Log(component.Name, "Info", logMessage )



		// ---------- Process output and decide what to forward ----- //
		// ...
		// ...
		// ...

		handoffMessage := defs.Message{Data: []string{"Z:"}}

		// ...
		// ...
		// ---------- Process output and decide what to forward ----- //



		for _,nextHandoff := range handoffTo {
			logMessage = fmt.Sprintf(
				"Handing off message to (%s) chain\n", handoffTo )
			component.Tlog.Log(component.Name, "Debug", logMessage )

			action.Handoff(nextHandoff, handoffMessage)
		}

	}

	// TODO: what should it be ?
	return true
}