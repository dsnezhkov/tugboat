package fs

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
)

type FSManager struct {
	confFS *embed.FS
	sPayFS *embed.FS
	dPayFS *memfs.MemFS
	logFS  *memfs.MemFS
}

func CreateFSManager(confFS *embed.FS, staticFS *embed.FS) *FSManager {
	fsm := &FSManager{
		confFS: confFS,
		logFS:  memfs.Create(),
		dPayFS: memfs.Create(),
		sPayFS: staticFS,
	}
	return fsm
}

func (fsm *FSManager) GetLogFS() *memfs.MemFS {
	return fsm.logFS
}

func (fsm *FSManager) GetConfigFS() *embed.FS {
	return fsm.confFS
}

func (fsm *FSManager) GetPayFS() (*memfs.MemFS) {
	return fsm.dPayFS
}
func (fsm *FSManager) GetPlugFS() (*embed.FS) {
	return fsm.sPayFS
}

// TODO: implement mkdir -p?

func (fsm *FSManager) CreateDir(fso *memfs.MemFS, path string) error {

	err := fso.Mkdir(path, 0777)
	if err != nil {
		return err
	}
	return nil
}
func (fsm *FSManager) WriteFile(kind string, path string, pay []byte) (int, error) {

	var fso *memfs.MemFS

	switch kind {
	case "log":
		fso = fsm.logFS
	case "pay_dynamic":
		fso = fsm.dPayFS
	default:
		err := fmt.Errorf("error: %s", "Invalid memfs type")
		return 0, err
	}

	var f vfs.File
	f, err := fso.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		err = fmt.Errorf("Cannot open file: %v\n", err)
		return 0, err
	} else {
		defer f.Close()
	}

	count, err := f.Write(pay)
	if err != nil {
		err = fmt.Errorf("file not written: %v", err)
		return 0, err
	}
	return count, err
}

func (fsm *FSManager) ReadFile(kind string, path string) ([]byte, error) {

	var contentBytes []byte
	var err error

	switch kind {
	case "log":

		fso := fsm.logFS

		file, err := fso.OpenFile(path, os.O_RDONLY, 0)

		if err != nil {
			return nil, err
		}
		defer file.Close()

		var f fs.FileInfo
		f, err = fso.Stat(file.Name())

		if err != nil {
			return nil, err
		}

		var size int64 = f.Size()
		contentBytes = make([]byte, size)

		bufr := bufio.NewReader(file)
		_, err = bufr.Read(contentBytes)


	case "pay_static":
		fso := fsm.sPayFS
		contentBytes, err = fso.ReadFile(path)
		if err != nil {
			return nil, err
		}

	case "pay_dynamic":

		fso := fsm.dPayFS

		file, err := fso.OpenFile(path, os.O_RDONLY, 0)

		if err != nil {
			return nil, err
		}
		defer file.Close()

		var f fs.FileInfo
		f, err = fso.Stat(file.Name())

		if err != nil {
			return nil, err
		}

		var size int64 = f.Size()
		contentBytes = make([]byte, size)

		bufr := bufio.NewReader(file)
		_, err = bufr.Read(contentBytes)

	default:
		err := fmt.Errorf("error: %s", "Invalid memfs type")
		return nil, err
	}

	return contentBytes, err
}
func (fsm *FSManager) ListEmbedFS(fso *embed.FS) error {

	err := fs.WalkDir(fso, ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			dir := ""
			if d.IsDir() == true {
				dir = "/"
			}
			fmt.Printf("%s%s\n", path, dir)
			return nil
		},
	)

	if err != nil {
		return err
	}
	return nil
}
func (fsm *FSManager) ListPlugins() error {

	err := fs.WalkDir(fsm.sPayFS, ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			dir := ""
			if d.IsDir() == true {
				dir = "/"
			}
			fmt.Printf("%s%s\n", path, dir)
			return nil
		},
	)

	if err != nil {
		return err
	}
	return nil
}
func (fsm *FSManager) CheckIfMemDirExists(fso *memfs.MemFS, path string) (bool, error) {
	var stat os.FileInfo
	stat, err := fso.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}
func (fsm *FSManager) CheckIfEmbDirExists(fso *embed.FS, path string) (int, error) {

	var dirent []fs.DirEntry
	dirent, err := fso.ReadDir(path)
	if err != nil {
		return 0, err
	}
	return len(dirent), nil
}
func (fsm *FSManager) ListMemFS(mfs *memfs.MemFS) []map[string]int64 {

	var dirs []string
	var path = "."
	var filesFound []map[string]int64
	filesFound = listDirContents(mfs, path, dirs, &filesFound)
	return filesFound
}
func (fsm *FSManager) ListPayloads() []map[string]int64 {

	var dirs []string
	var path = "."
	var filesFound []map[string]int64
	filesFound = listDirContents(fsm.dPayFS, path, dirs, &filesFound)
	return filesFound
}

func processed(fileName string, processedDirectories []string) bool {
	for i := 0; i < len(processedDirectories); i++ {
		if processedDirectories[i] != fileName {
			continue
		}
		return true
	}
	return false
}
func listDirContents(mfs *memfs.MemFS, path string, dirs []string, found *[]map[string]int64) []map[string]int64 {
	files, _ := mfs.ReadDir(path)

	for _, f := range files {
		var newPath string
		if path != "/" {
			newPath = fmt.Sprintf("%s/%s", path, f.Name())
		} else {
			newPath = fmt.Sprintf("%s%s", path, f.Name())
		}

		if f.IsDir() {
			if !processed(newPath, dirs) {
				dirs = append(dirs, newPath)
				listDirContents(mfs, newPath, dirs, found)
			}
		} else {
			fileInfo, err := mfs.Stat(newPath)
			var sz int64
			if err!=nil {
				sz = 0
			}
			sz = fileInfo.Size()
			entry := make(map[string]int64,0)
			entry[newPath] = sz
			*found = append(*found, entry)
		}
	}

	return *found
}
