package main

import (
	"log"
	"os"

	"github.com/golang-book/penv"

	"gopkg.in/alecthomas/kingpin.v2"
)

var help = `usage: penv command

Where command is:
  set     Set an environmental variable
  unset   Unset an environmental variable
`
var helpSet = `usage: penv set [NAME=VALUE]...
Permanently set each NAME to VALUE in the environment.
`
var helpUnset = `usage: penv unset [NAME]...
Permanently unset each NAME in the environment. Only variables added by
'set' can be 'unset'.
`

func main() {
	log.SetFlags(0)

	app := kingpin.New("penv", "Permanently set/unset environment variables")
	app.Version("1.0")

	setcmd := app.Command("set", "Permanently NAME to VALUE in the environment")
	setcmdName := setcmd.Arg("name", "environment variable name").Required().String()
	setcmdValue := setcmd.Arg("value", "environment variable value").Required().String()

	unsetcmd := app.Command("unset", "Permanently unset NAME in the environment")
	unsetcmdName := unsetcmd.Arg("name", "environment variable name").Required().String()

	appendcmd := app.Command("append", "Permanently append VALUE to NAME in the environment")
	appendcmdName := appendcmd.Arg("name", "environment variable name").Required().String()
	appendcmdValue := appendcmd.Arg("value", "environment vairable value").Required().String()

	var err error
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case setcmd.FullCommand():
		err = penv.SetEnv(*setcmdName, *setcmdValue)
	case unsetcmd.FullCommand():
		err = penv.UnsetEnv(*unsetcmdName)
	case appendcmd.FullCommand():
		err = penv.AppendEnv(*appendcmdName, *appendcmdValue)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
