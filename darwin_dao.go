//+build darwin

package penv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	darwinPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
   <key>Label</key>
   <string>com.github.badgerodon.penv</string>
   <key>ProgramArguments</key>
   <array>
     <string>bash</string>
     <string>` + filepath.Join(os.Getenv("HOME"), ".config", "penv.sh") + `</string>
   </array>
   <key>RunAtLoad</key>
   <true/>
</dict>
</plist>`

	darwinShell = &shell{
		configFileName: filepath.Join(os.Getenv("HOME"), ".config", "penv.sh"),
		commentSigil:   " #",
		quote: func(value string) string {
			r := strings.NewReplacer(
				"\\", "\\\\",
				"'", "\\'",
				"\n", `'"\n"'`,
				"\r", `'"\r"'`,
			)
			return "'" + r.Replace(value) + "'"
		},
		mkSet: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"launchctl setenv %s %s",
				nv.Name, sh.quote(nv.Value),
			)
		},
		mkAppend: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"launchctl setenv %s ${%s}${%s:+:}%s",
				nv.Name, nv.Name, nv.Name, sh.quote(nv.Value),
			)
		},
		mkUnset: func(sh *shell, nv NameValue) string {
			return fmt.Sprintf(
				"launchctl unsetenv %s",
				nv.Name,
			)
		},
	}
)

// DarwinDAO is the data access object for OSX
type DarwinDAO struct {
}

func init() {
	RegisterDAO(500, func() bool {
		return runtime.GOOS == "darwin"
	}, &DarwinDAO{})
}

// Load loads the environment
func (dao *DarwinDAO) Load() (*Environment, error) {
	return darwinShell.Load()
}

// Save saves the environment
func (dao *DarwinDAO) Save(env *Environment) error {
	err := darwinShell.Save(env)
	if err != nil {
		return err
	}

	pListFolder := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
	_, err = os.Stat(pListFolder)
	switch {
	case os.IsNotExist(err):
		if err = os.MkdirAll(pListFolder, 0777); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	pListName := filepath.Join(pListFolder, "penv.plist")

	err = ioutil.WriteFile(pListName, []byte(darwinPlist), 0777)
	if err != nil {
		return err
	}

	if err := exec.Command("launchctl", "unload", pListName).Run(); err != nil {
		return err
	}

	return exec.Command("launchctl", "load", pListName).Run()
}
