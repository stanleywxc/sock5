package main

import (
        "os"
)


type Args struct {
    args		map[string]string    
}

func (arg *Args) Get(key string) (string) {
    return arg.args[key]
}

func parseArgs() (*Args) {
    
    var arg Args = Args {}
    
    arg.args = make(map[string]string)
    
    for i := 1; i < len(os.Args); i++ {
        
        switch os.Args[i] {
            case "-f" :
                if ((i+1) >= len(os.Args)) {
                    continue
                }
                arg.args[os.Args[i]] = os.Args[i+1]
                i++
                break
            default:
                break
        }
    }
    
    return &arg
}