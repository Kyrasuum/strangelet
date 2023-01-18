//adapted from 'micro's way of accomplishing this task https://github.com/zyedidia/micro

package config

import (
	"io/ioutil"
	"path"
	"path/filepath"

	rt "strangelet/runtime"

	"github.com/zyedidia/highlight"
)

const (
	RTColorscheme  = 0
	RTSyntax       = 1
	RTHelp         = 2
	RTPlugin       = 3
	RTSyntaxHeader = 4
)

var (
	NumTypes   = 5 // How many filetypes are there
	SyntaxDefs = []*highlight.Def{}
)

type RTFiletype int

// RuntimeFile allows the program to read runtime data like colorschemes or syntax files
type RuntimeFile interface {
	// Name returns a name of the file without paths or extensions
	Name() string
	// Data returns the content of the file.
	Data() ([]byte, error)
}

// allFiles contains all available files, mapped by filetype
var allFiles [][]RuntimeFile
var realFiles [][]RuntimeFile

func init() {
	allFiles = make([][]RuntimeFile, NumTypes)
	realFiles = make([][]RuntimeFile, NumTypes)
}

// NewRTFiletype creates a new RTFiletype
func NewRTFiletype() int {
	NumTypes++
	allFiles = append(allFiles, []RuntimeFile{})
	realFiles = append(realFiles, []RuntimeFile{})
	return NumTypes - 1
}

// some file on filesystem
type realFile string

// some asset file
type assetFile string

// some file on filesystem but with a different name
type namedFile struct {
	realFile
	name string
}

// a file with the data stored in memory
type memoryFile struct {
	name string
	data []byte
}

func (mf memoryFile) Name() string {
	return mf.name
}
func (mf memoryFile) Data() ([]byte, error) {
	return mf.data, nil
}

func (rf realFile) Name() string {
	fn := filepath.Base(string(rf))
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

func (rf realFile) Data() ([]byte, error) {
	return ioutil.ReadFile(string(rf))
}

func (af assetFile) Name() string {
	fn := path.Base(string(af))
	return fn[:len(fn)-len(path.Ext(fn))]
}

func (af assetFile) Data() ([]byte, error) {
	return rt.Asset(string(af))
}

func (nf namedFile) Name() string {
	return nf.name
}

// AddRuntimeFile registers a file for the given filetype
func AddRuntimeFile(fileType RTFiletype, file RuntimeFile) {
	allFiles[fileType] = append(allFiles[fileType], file)
}

// AddRealRuntimeFile registers a file for the given filetype
func AddRealRuntimeFile(fileType RTFiletype, file RuntimeFile) {
	allFiles[fileType] = append(allFiles[fileType], file)
	realFiles[fileType] = append(realFiles[fileType], file)
}

// AddRuntimeFilesFromDirectory registers each file from the given directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromDirectory(fileType RTFiletype, directory, pattern string) {
	files, _ := ioutil.ReadDir(directory)
	for _, f := range files {
		if ok, _ := filepath.Match(pattern, f.Name()); !f.IsDir() && ok {
			fullPath := filepath.Join(directory, f.Name())
			AddRealRuntimeFile(fileType, realFile(fullPath))
		}
	}
}

// AddRuntimeFilesFromAssets registers each file from the given asset-directory for
// the filetype which matches the file-pattern
func AddRuntimeFilesFromAssets(fileType RTFiletype, directory, pattern string) {
	files, err := rt.AssetDir(directory)
	if err != nil {
		return
	}
	for _, f := range files {
		if ok, _ := path.Match(pattern, f); ok {
			AddRuntimeFile(fileType, assetFile(path.Join(directory, f)))
		}
	}
}

// FindRuntimeFile finds a runtime file of the given filetype and name
// will return nil if no file was found
func FindRuntimeFile(fileType RTFiletype, name string) RuntimeFile {
	for _, f := range ListRuntimeFiles(fileType) {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

// ListRuntimeFiles lists all known runtime files for the given filetype
func ListRuntimeFiles(fileType RTFiletype) []RuntimeFile {
	return allFiles[fileType]
}

// ListRealRuntimeFiles lists all real runtime files (on disk) for a filetype
// these runtime files will be ones defined by the user and loaded from the config directory
func ListRealRuntimeFiles(fileType RTFiletype) []RuntimeFile {
	return realFiles[fileType]
}

// InitRuntimeFiles initializes all assets file and the config directory
func InitRuntimeFiles() {
	add := func(fileType RTFiletype, dir, pattern string) {
		AddRuntimeFilesFromDirectory(fileType, filepath.Join(ConfigDir, dir), pattern)
		AddRuntimeFilesFromAssets(fileType, path.Join("runtime", dir), pattern)
	}

	add(RTColorscheme, "colorschemes", "*.scheme")
	add(RTSyntax, "syntax", "*.yaml")
	add(RTSyntaxHeader, "syntax", "*.hdr")
	add(RTHelp, "help", "*.md")

	LoadDefinitions()
}

func LoadDefinitions() {
	rtfiles := ListRuntimeFiles(RTSyntax)
	for _, rtfile := range rtfiles {
		SyntaxData, err := rtfile.Data()
		if err != nil {
			continue
		}
		syntaxDef, err := highlight.ParseDef(SyntaxData)
		SyntaxDefs = append(SyntaxDefs, syntaxDef)
	}
}

func DetectType(fileName string, firstLine []byte) *highlight.Def {
	return highlight.DetectFiletype(SyntaxDefs, fileName, firstLine)
}
