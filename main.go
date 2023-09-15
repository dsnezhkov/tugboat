package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/dsnezhkov/tugboat/action"
	"github.com/dsnezhkov/tugboat/comms"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/fs"
	"github.com/dsnezhkov/tugboat/logger"
	"github.com/dsnezhkov/tugboat/payloads"
	"github.com/dsnezhkov/tugboat/workflow"
)

// Embedding configuration within binary

//go:embed config/*
var cfs embed.FS

//go:embed plugins/*
var dfs embed.FS

// /////////////////////////////////////////

func main() {

	defs.Ps = defs.NewPubsub()
	// FS Mgmt layer
	defs.FsManager = fs.CreateFSManager(&cfs, &dfs)
	// Payload Mgmt Layer
	defs.PayManager = payloads.CreatePayloadManager(defs.FsManager)
	// Comms Mgmt Layer
	defs.CommsManager = comms.CreateCommsManager(comms.CommsManagerOptions{Payman: defs.PayManager})

	defs.Tlog = logger.CreateLogManager(defs.FsManager)

	log.Println("Static resources FS ToC:")
	_ = defs.FsManager.ListEmbedFS(&dfs)
	log.Println("Conf resources FS ToC:")
	_ = defs.FsManager.ListEmbedFS(&cfs)

	defs.Tlog.Log("main", "DEBUG", "Parse Flow config file ")
	defs.FlowConf = defs.NewFlowConfig()
	defs.FlowConf.CreateFromConfig(defs.FsManager, "config.hcl")

	defs.Tlog.Log("main", "DEBUG", "Process Global config section  ")
	defs.FlowConf.PopulateGlobalConfig()

	defs.Tlog.Log("main", "DEBUG", "Process and Initialize  vectors")
	defs.FlowConf.PopulateAvailableVectors()

	defs.Tlog.Log("main", "DEBUG", "Initialize Workflows")
	workflow.Workflow = workflow.CreateWorkFlow()

	defs.Tlog.Log("main", "DEBUG", "Initializing workflows")
	workflow.Workflow.Setup()

	// Health check
	defs.Tlog.Log("main", "INFO", "Starting memory profiler")
	action.StartProfiler()

	defs.Tlog.Log("main", "INFO", "Starting the workflow chain...")
	workflowHead := workflow.Workflow.GetWHead()

	logMessage := fmt.Sprintf("%s : %s", "Starting the chain from head", workflowHead)
	defs.Tlog.Log("main", "INFO", logMessage)

	time.Sleep(50 * time.Millisecond)
	action.Handoff(workflowHead, defs.Message{})

	// Send signal to component
	defs.Tlog.Log("main", "INFO", "Sending singal to comp_tech for runaway")
	go func() {
		time.Sleep(20 * time.Second)
		defs.ChanSignal2component["comp_tech"] <- true
	}()

	defs.Tlog.Log("main", "INFO", "Waiting for workers to finish")
	defs.Wg.Wait()
	defs.Tlog.Log("main", "INFO", "Workers finished")

	// // Post-processing ///

	defs.Tlog.Log("main", "Info", "Listing Payloads")

	defs.Tlog.Log("main", "Info", "All tasks finished. Closing subscriptions")
	defs.Ps.Close()

	// defs.Tlog.Log("main", "Info", "Reading Logs")
	// logsCollected := defs.FsManager.ListMemFS(defs.FsManager.GetLogFS())
	//
	// println(logsCollected)
	// for _, l := range logsCollected {
	//
	// 	l = strings.Replace(l, ".", "", 1)
	//
	// 	log.Printf("-- %s --\n", l)
	// 	data, err := defs.FsManager.ReadFile("log", l)
	// 	if err != nil {
	// 		log.Printf("Error Reading file: %v", err)
	// 	} else {
	// 		log.Printf("Read: \n%s\n", data)
	// 	}
	//
	// }
	defs.Tlog.Log("main", "Info", "-- DONE --")

}
