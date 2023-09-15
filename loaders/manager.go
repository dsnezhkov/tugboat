package loaders

import (
	"fmt"
	"path"

	"github.com/Binject/universal"
	"tugboat/defs"
)

type TLoader interface {
	Load(string)
}

type UniversalLoader struct {
	loader *universal.Loader
}

// Singleton loader
var uloader *UniversalLoader = nil

func GetUniversalLoader() (*UniversalLoader, error) {

	var err error

	if uloader != nil {
		return uloader, nil
	}else{
		uloader = &UniversalLoader{}
		uloader.loader, err = universal.NewLoader()
		if err != nil {
			defs.Tlog.Log("Loader", "ERROR", fmt.Sprintf("UniversalLoader error: %v", err) )
			return nil, err
		}
		return uloader, nil
	}
}

func (uloader *UniversalLoader) Load(modulePath string ) (*universal.Library, error) {

	var err error

	content, err := uloader.getPayLoadContent(modulePath)
	if err != nil {
		return nil, err
	}

	library, err := uloader.loader.LoadLibrary(path.Base(modulePath), &content)
	if err != nil {
		return nil, err
	}
	return library, nil
}

func (uloader *UniversalLoader) getPayLoadContent(modulePath string ) ([]byte, error)  {
	contentBytes, err := defs.PayManager.GetPayload(modulePath)
	if err != nil {
		defs.Tlog.Log("Loader", "ERROR", fmt.Sprintf("UniversalLoader error: %v", err) )
		return nil, err
	}
	return contentBytes, nil
}


func (payman *UniversalLoader) RunExport(library *universal.Library, funcName string, args ...uintptr ) (uintptr, error) {

	val, err := library.Call(funcName, args...)
	if err != nil {
		err = fmt.Errorf("loader library call error: %e", err)
		return 0, err
	}
	return val, nil
}


func (payman *UniversalLoader) StatExport(library *universal.Library, funcName string) {

	ptr, found := library.FindProc("Runme")
	if found {
		fmt.Printf("Name: %s Base: %d, FPtr: %v\n", library.Name, library.BaseAddress, ptr)
	}else {
		fmt.Printf("Export %s not found\n", funcName)
	}

}
func (payman *UniversalLoader) ListExports(library *universal.Library) {

	for k, v := range library.Exports {
		fmt.Printf("%s -> %d\n", k, v)
	}

}