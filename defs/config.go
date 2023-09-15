package defs

import (
	"embed"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"tugboat/fs"
)


func NewFlowConfig() *Config {
	c := &Config{}
	return c
}

func (flowconf *Config) CreateFromConfig(fsm *fs.FSManager, cfgPath string ) {
	status, err := FlowConf.parseConfig(*fsm.GetConfigFS(),
		strings.Join([]string{"config", cfgPath}, "/"))

	if err != nil && status == false {
		log.Fatalf("Error: Cannot process flow config file: %v", err)

	}
	//fmt.Printf("%#v", FlowConf)
}


func (flowconf *Config) parseConfig(cfs embed.FS, filePath string) (bool, error) {

	// Recover from Json panics
	defer func() bool {
		if err := recover(); err != nil {
			log.Printf("DAG: Execution encountered fault, recovered from: %v", err)
			return false
		}
		return true
	}()

	cfgData, err := cfs.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	err = hclsimple.Decode(filePath, cfgData, nil, flowconf)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	_, err = json.MarshalIndent(flowconf, "", "\t")
	if err!= nil {
		log.Printf("json.MarshalIndent: %v", err)
		return false, err
	}
	//fmt.Print(string(s))

	return true, nil
}

func (flowconf *Config) PopulateAvailableVectors () {

	// Map subscriptions os component names to vectors
	log.Println("Main: Creating component subscriptions to vectors - HCL ")
	for _, vc := range FlowConf.VectorBlock.Vectors {
		for _, cl := range vc.ComponentLinks {
			VectorsAvailable[vc.Name] =
				append(VectorsAvailable[vc.Name], cl)
			log.Printf("Main:\tSubscribing component %s to vector %s\n",
				vc.Name, ComponentAvailable[cl])

		}
	}
}
func (flowconf *Config) PopulateGlobalConfig () {
	log.Println("Main: Setting global options ")
	Tlog.SetLevel(FlowConf.GlobalOpts.LoggerOpts.Level)
	Tlog.SetLogLocation(FlowConf.GlobalOpts.LoggerOpts.Location)
	for _, l := range FlowConf.GlobalOpts.PayloadOpts{
		err := PayManager.SetPayloadsLocation(l.Variant, l.Location)
		if err != nil {
			// TODO: get dynamic name of function for errors
			log.Printf("PopulateGlobalConfig: %v", err)
		}
	}

}
