//+build windows

package penv

import (
	"log"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// WindowsDAO is the data access object for windows
type WindowsDAO struct {
}

func init() {
	RegisterDAO(1000, func() bool {
		return runtime.GOOS == "windows"
	}, &WindowsDAO{})
}

// Load loads the environment
func (dao *WindowsDAO) Load() (*Environment, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	ki, err := key.Stat()
	if err != nil {
		return nil, err
	}

	names, err := key.ReadValueNames(int(ki.ValueCount))
	if err != nil {
		return nil, err
	}

	env := &Environment{
		Setters:   make([]NameValue, 0),
		Appenders: make([]NameValue, 0),
		Unsetters: make([]NameValue, 0),
	}

	for _, name := range names {
		value, _, err := key.GetStringValue(name)
		if err != nil {
			return nil, err
		}
		env.Setters = append(env.Setters, NameValue{strings.ToUpper(name), value})
	}

	return env, nil
}

// Save saves the environment
func (dao *WindowsDAO) Save(env *Environment) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	ki, err := key.Stat()
	if err != nil {
		return err
	}

	names, err := key.ReadValueNames(int(ki.ValueCount))
	if err != nil {
		return err
	}

	log.Println(names)
	log.Println("SETTERS", env.Setters)
	log.Println("APPENDERS", env.Appenders)

	// set
	for _, nv := range env.Setters {
		err = key.SetStringValue(nv.Name, nv.Value)
		if err != nil {
			return err
		}
	}

	// append
	for _, nv := range env.Appenders {
		values := []string{}
		for _, name := range names {
			if strings.ToUpper(name) == strings.ToUpper(nv.Name) {
				value, _, err := key.GetStringValue(name)
				if err != nil {
					return err
				}
				values = append(values, strings.Split(value, ";")...)
			}
		}
		values = uniquei(append(values, nv.Value))
		err = key.SetStringValue(nv.Name, strings.Join(values, ";"))
		if err != nil {
			return err
		}
	}

	// unset
	for _, name := range names {
		for _, nv := range env.Unsetters {
			if nv.Name == name {
				err = key.DeleteValue(name)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
