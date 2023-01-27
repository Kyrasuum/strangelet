package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"strangelet/internal/events"

	//tea "github.com/charmbracelet/bubbletea"
	json5 "github.com/zyedidia/json5"
)

var (
	Bindings map[string]map[string]string = make(map[string]map[string]string)
)

const ()

func InitBindings() error {
	events.InitActions()

	Bindings["Global"] = map[string]string{}
	Bindings["LogWindow"] = map[string]string{}
	Bindings["Filebrowser"] = map[string]string{}
	Bindings["Split"] = map[string]string{}
	Bindings["Terminal"] = map[string]string{}
	for scope, _ := range Bindings {
		Bindings[scope] = DefaultBindings(scope)
	}

	return LoadConfigBindings()
}

func LoadConfigBindings() error {
	filename := filepath.Join(ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)

	var parsed map[string]map[string]interface{}
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return err
		}
	}

	for scope, binding := range parsed {
		for k, v := range binding {
			switch val := v.(type) {
			case string:
				bindKey(scope, k, val)
			default:
				return fmt.Errorf("Error reading bindings.json: non-string and non-map entry: %+v", k)
			}
		}
	}

	return nil
}

func createBindingsIfNotExist(filename string) {
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		ioutil.WriteFile(filename, []byte("{}"), 0644)
	}
}

func bindKey(scope string, k string, v string) {
	if _, ok := Bindings[scope]; !ok {
		return
	}
	Bindings[scope][k] = v
}

// TryBindKey tries to bind a key by writing to config.ConfigDir/bindings.json
// Returns true if the keybinding already existed and a possible error
func TryBindKey(scope string, k string, v string, overwrite bool) (bool, error) {
	var e error
	var parsed map[string]map[string]string

	filename := filepath.Join(ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)

	if _, ok := Bindings[scope]; !ok {
		return false, errors.New("Scope does not exist for desired keybind")
	}

	if _, e = os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return false, errors.New("Error reading bindings.json file: " + err.Error())
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return false, errors.New("Error reading bindings.json: " + err.Error())
		}

		found := false
		for key := range parsed[scope] {
			if key == k {
				if overwrite {
					parsed[scope][key] = v
				}
				found = true
				break
			}
		}

		if found && !overwrite {
			return true, nil
		} else if !found {
			parsed[scope][k] = v
		}

		bindKey(scope, k, v)

		txt, _ := json.MarshalIndent(parsed, "", "    ")
		return found, ioutil.WriteFile(filename, append(txt, '\n'), 0644)
	}
	return false, e
}

// UnbindKey removes the binding for a key from the bindings.json file
func UnbindKey(scope string, k string) error {
	var e error
	var parsed map[string]map[string]string

	filename := filepath.Join(ConfigDir, "bindings.json")
	createBindingsIfNotExist(filename)
	if _, e = os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New("Error reading bindings.json file: " + err.Error())
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			return errors.New("Error reading bindings.json: " + err.Error())
		}

		for key := range parsed[scope] {
			if key == k {
				delete(parsed[scope], key)
				break
			}
		}

		defaults := DefaultBindings(scope)
		if a, ok := defaults[k]; ok {
			bindKey(scope, k, a)
		} else if _, ok := Bindings[scope][k]; ok {
			delete(Bindings[scope], k)
		}

		txt, _ := json.MarshalIndent(parsed, "", "    ")
		return ioutil.WriteFile(filename, append(txt, '\n'), 0644)
	}
	return e
}

func DefaultBindings(scope string) map[string]string {
	switch scope {
	case "Global":
		return map[string]string{
			"ctrl+q": "Quit",
			"alt+k":  "FocusFileBrowser",
			"ctrl+k": "ToggleFileBrowser",
			"ctrl+l": "ToggleLogWindow",
			"ctrl+e": "FocusCommand",
		}
	case "LogWindow":
		return map[string]string{}
	case "Filebrowser":
		return map[string]string{}
	case "Split":
		return map[string]string{
			"ctrl+n": "NewSplit",
			// "":           "CloseSplit",

			"ctrl+t": "NewTab",
			"ctrl+s": "SaveTab",
			// "":           "CloseTab",
		}
	case "Terminal":
		return map[string]string{}
	default:
		return map[string]string{}
	}
}
