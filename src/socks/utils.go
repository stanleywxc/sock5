package socks

import (
        "fmt"
        "log"
        "sync"
)

var locker sync.Mutex

func init() { 
}

func FormatLog(data []byte) {
   
    var line 	string = ""
    var linex	string = ""
    
    for i, v := range data {
        
        line += fmt.Sprintf("%02x ", v)
        
        if ((v > 32) && (v < 127)){
            linex += fmt.Sprintf("%c", v)
        } else {
            linex += fmt.Sprintf(".")
        }            
        
        if (((i +1) % 8) == 0) {
            line += fmt.Sprintf(" ")
            linex += fmt.Sprintf(" ")
        }
        
        if (((i + 1) %16) == 0) {
            log.Printf("%s%s\n", line, linex)
            line  	= ""
            linex 	= ""
        }
    }
    
    log.Printf("\n")
}

