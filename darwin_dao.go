//+build darwin

package penv

import (
	"os"
	"path/filepath"
	"runtime"
)

var darwinPlist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
   <key>Label</key>
   <string>com.github.badgerodon.penv</string>
   <key>Program</key>
   <string>` + filepath.Join(os.Getenv("HOME"), ".config", "penv.sh") + `</string>
   <key>RunAtLoad</key>
   <true/>
</dict>
</plist>`

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
	panic("NOT IMPLEMENTED")
}

// Save saves the environment
func (dao *DarwinDAO) Save(env *Environment) error {
	panic("NOT IMPLEMENTED")
}
