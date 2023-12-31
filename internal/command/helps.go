package command

const CmdDetail = `
Commands:
  start    Start app
   -d  Daemon process
  restart  Restart daemon process
  stop     Stop daemon process
`

const NoCmdHelp = `Welcome To Go-Gourd Command Console!
Usage: %s COMMAND` + CmdDetail

const UndefinedCmdHelp = `The command '%s' is not exists!
Welcome To Go-Gourd Command Console!
Usage: %s COMMAND` + CmdDetail
