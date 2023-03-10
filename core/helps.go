package core

const cmd = `
  start    Start process.
	-d     Daemon process.
  stop     Stop daemon process.
`

const NoCmdHelp = `Welcome To Go-Gourd Command Console!
Usage: %s COMMAND
Commands:` + cmd

const UndefinedCmdHelp = `The command '%s' is not exists!
Welcome To Go-Gourd Command Console!
Usage: %s COMMAND
Commands:` + cmd
