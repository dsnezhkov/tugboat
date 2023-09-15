package util

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/dsnezhkov/tugboat/defs"
)

type MessageCheckState struct {
	Data bool
}

func CheckMessage(msg *defs.Message) bool {

	msgCheckState := MessageCheckState{
		Data: false,
	}
	if msg.Data != nil {
		msgCheckState.Data = true
	}
	return true
}

// TODO: Logic and err handling, abstract type

func FindVectorForComponent(componentName string) (component string, vector string) {

	for k, v := range defs.VectorsAvailable {
		for _, i := range v {
			if i == componentName {
				fmt.Printf("C -> %s :  V:%s\n", i, k)
				return i, k
			}
		}
	}

	return "", ""
}

func FindNextVectorLinkInFlow(config *defs.Config, vectorName string) (vectorNamesNext []string) {

	for _, i := range config.FlowLinkBlock.FlowLinks {
		if i.Name == vectorName {
			return i.VectorNames
		}
	}
	return []string{}
}

// CmdExec Execute a command
func CmdExec(sout *strings.Builder, serr *strings.Builder, args ...string) error {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Printf("CmdExec: %v", args)

	cmd := exec.Command(baseCmd, cmdArgs...)

	//var stdoutBuf, stderrBuf bytes.Buffer

	// create a pipe for the output of the script
	cmdStdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("StdoutPipe failed with %v\n", err)
		return err
	}

	// create a pipe for the error of the script
	cmdStdErrReader, err := cmd.StderrPipe()
	if err != nil {
		err = fmt.Errorf("StderrPipe failed with %v\n", err)
		return err
	}

	scannerOut := bufio.NewScanner(cmdStdOutReader)
	go func() {
		err := func() error {
			for scannerOut.Scan() {
				// TODO: recover on ErrTooLarge
				_, err = sout.Write(scannerOut.Bytes())
				if err != nil {
					return err
				}

				// https://github.com/golang/go/issues/28822
				os := runtime.GOOS
				switch os {
				case "windows":
					_, err = sout.WriteString("\n")
				case "darwin":
				case "linux":
				default:
				}

			}
			return nil
		}()
		if err != nil {
			log.Printf("Error in processing cmd stdout")
		}
	}()

	scannerErr := bufio.NewScanner(cmdStdErrReader)
	go func() {
		err := func() error {
			for scannerErr.Scan() {
				// TODO: recover on ErrTooLarge
				_, err = serr.Write(scannerErr.Bytes())
				if err != nil {

					return err
				}

				// https://github.com/golang/go/issues/28822
				os := runtime.GOOS
				switch os {
				case "windows":
					_, err = sout.WriteString("\n")
				case "darwin":
				case "linux":
				default:
				}

			}
			return nil
		}()
		if err != nil {
			log.Printf("Error in processing cmd stderr")
		}
	}()

	//cmd.Stdout = io.MultiWriter(&stdoutBuf)
	//cmd.Stderr = io.MultiWriter(&stderrBuf)

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("cmd.Start() failed with %v\n", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		err = fmt.Errorf("cmd.Wait() failed with %v\n", err)
		return err
	}

	//outStr, errStr := string(stdoutBuf.Bytes()),
	//string(stderrBuf.Bytes())

	//fmt.Printf("E: %s | O: %s", stderrBuf.String(), stdoutBuf.String())

	return nil
}

func ListModulesLoadsForMe(component defs.IComponent) {
	for _, cm := range component.GetModules() {
		fmt.Printf("module %s:\n", cm.Name)
		for i, cl := range cm.CLoads {
			fmt.Printf("%d) %s - (%s, %s)\n", i, cl.Source, cl.DataType, cl.NativeFormat)

		}
	}

}

func GetModuleByName(component defs.IComponent, cmoduleName string) (*defs.CModule, bool) {
	for _, v := range component.GetModules() {
		if v.Name == cmoduleName {
			return &v, true
		}
	}
	return nil, false
}
func GetModuleLoadByName(cmodule defs.CModule, cmloadName string) (*defs.CLoad, bool) {
	for _, l := range cmodule.CLoads {
		if l.Name == cmloadName {
			return &l, true
		}
	}
	return nil, false
}

func GetModuleLoadchain(component defs.IComponent, cmoduleName string, cmloadName string) (*defs.CLoad, bool) {
	if m, ok := GetModuleByName(component, cmoduleName); ok {
		if l, ok := GetModuleLoadByName(*m, cmloadName); ok {
			return l, true
		}
	}
	return nil, false
}
