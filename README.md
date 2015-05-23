# penv
`penv` permanently sets environment variables. It supports the following:

* bash - entries are added to `~/.bashrc`
* fish - entries are added to `~/.config/fish/config.fish`
* windows - entries are added to the registry for the current user

## Installation
`penv` is both a library and a command. To use the library in your own code see
(godoc.org/github.com/badgerodon/penv). For the command:

    go get github.com/badgerodon/penv/...

This creates a `penv` command. Here's its usage:

    usage: penv [<flags>] <command> [<args> ...]

    Permanently set/unset environment variables

    Flags:
      --help     Show help.
      --version  Show application version.

    Commands:
      help [<command>...]
        Show help.

      set <name> <value>
        Permanently NAME to VALUE in the environment

      unset <name>
        Permanently unset NAME in the environment

      append <name> <value>
        Permanently append VALUE to NAME in the environment

## License
MIT
