package payloads

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/dsnezhkov/tugboat/fs"
)

type PayloadManager struct {
	fsMan           *fs.FSManager
	staticLocation  string
	dynamicLocation string
}

var payloadManager *PayloadManager

func CreatePayloadManager(fsm *fs.FSManager) *PayloadManager {
	if payloadManager == nil {

		payloadManager = &PayloadManager{
			fsMan:           fsm,
			staticLocation:  "/pays",
			dynamicLocation: "/pays",
		}
	}
	return payloadManager
}

func (payman *PayloadManager) SetPayloadsLocation(variant string, location string) error {

	var err error

	switch variant {
	case "static":
		// emfs does not display first `/` so we compensate for this
		re, err := regexp.Compile(`embfs://(.*)`)
		if err != nil {
			return err
		}

		matches := re.FindStringSubmatch(location)

		if len(matches) != 2 {
			err = fmt.Errorf("invalid location format: %s", variant)
			return err
		}
		payman.staticLocation = matches[1]

	case "dynamic":
		re, err := regexp.Compile(`memfs:/(/.*)`)

		if err != nil {
			return err
		}
		matches := re.FindStringSubmatch(location)

		if len(matches) != 2 {
			err = fmt.Errorf("invalid location format: %s", variant)
			return err
		}
		payman.dynamicLocation = matches[1]

	default:
		err = fmt.Errorf("invalid payload variant: %s", variant)
		return err
	}

	return nil
}

// GetPayload Figure out if payload static or dynamic and fetch it.
func (payman *PayloadManager) GetPayload(location string) ([]byte, error) {
	var (
		err error
	)

	// Static
	embFSRe := regexp.MustCompile(`embfs://(.*)`)
	memFSRe := regexp.MustCompile(`memfs:/(/.*)`)

	EmbFSReMatches := embFSRe.FindStringSubmatch(location)

	if len(EmbFSReMatches) == 2 {
		content, err := payman.GetStaticPayload(EmbFSReMatches[1])
		if err != nil {
			return nil, err
		}
		return content, nil
	}
	memFSReMatches := memFSRe.FindStringSubmatch(location)

	if len(memFSReMatches) == 2 {
		content, err := payman.GetDynamicPayload(memFSReMatches[1])
		if err != nil {
			return nil, err
		}
		return content, nil
	}

	err = fmt.Errorf("Unable to match locaton %s to available FS type\n", location)
	return nil, err
}

func (payman *PayloadManager) GetPayloadLocation(variant string) (string, error) {

	var err error

	switch variant {
	case "static":
		return payman.staticLocation, nil

	case "dynamic":
		return payman.dynamicLocation, nil

	default:
		err = fmt.Errorf("invalid payload variant: %s", variant)
		return "", err
	}

}

func (payman *PayloadManager) GetStaticPayload(path string) ([]byte, error) {

	contentBytes, err := payman.fsMan.ReadFile("pay_static", path)
	if err != nil {
		return nil, err
	}
	return contentBytes, nil

}

func (payman *PayloadManager) GetDynamicPayload(path string) ([]byte, error) {

	contentBytes, err := payman.fsMan.ReadFile("pay_dynamic", path)
	if err != nil {
		return nil, err
	}
	return contentBytes, nil

}
func (payman *PayloadManager) SetDynamicPayload(path string, content []byte) (int, error) {
	payfs := payman.fsMan.GetPayFS()

	_, err := payman.fsMan.CheckIfMemDirExists(payfs, payman.dynamicLocation)
	if err != nil {
		log.Printf("Payload Manager: unable to verify payload directory: %v\n", err)
		log.Printf("Payload Manager: attempting to create payload directory\n")

		err = payman.fsMan.CreateDir(payfs, payman.dynamicLocation)
		if err != nil {
			log.Printf("Payload Manager: unable to create top level payload directory %v\n", err)
			return 0, err
		}
	}

	payloadFileLocation := strings.Join([]string{payman.dynamicLocation, path}, "/")
	written, err := payman.fsMan.WriteFile("pay_dynamic", payloadFileLocation, content)
	if err != nil {
		return 0, err
	}
	return written, nil
}
func (payman *PayloadManager) ListDynamicPayloads() []map[string]int64 {
	return payman.fsMan.ListPayloads()
}
