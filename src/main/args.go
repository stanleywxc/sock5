package main

import (
        "fmt"
        "os"
)

const HELP_MSG = "Usage: main --help | -f conf-file\n"

type Args struct {
    args		map[string]string    
}

func (arg *Args) Get(key string) (string) {
    return arg.args[key]
}

func parseArgs() (*Args, string) {
    
    var msg	string 	= ""
    var arg Args 	= Args {}
    
    arg.args = make(map[string]string)
    
    for i := 1; i < len(os.Args); i++ {
        
        switch os.Args[i] {
            case "-f" :
                if ((i+1) >= len(os.Args)) {
                    msg = "Param '-f' is provided, but missing '-f' value\nExample usage: 'main -f /tmp/sock5.conf'\nUsing switch '--help' for help info"
                    return &arg, msg
                }
                arg.args[os.Args[i]] = os.Args[i+1]
                i++
                break
            case "--help":
                return &arg, HELP_MSG
            default:
                msg = fmt.Sprintf("Unsupported switch: '%s'\nTo get help info: 'main --help'\n", os.Args[i])
                return &arg, msg
        }
    }
    
    return &arg, msg
}