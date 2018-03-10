package log

import (
        "fmt"
        "log"
        "os"
)

type Level uint32

const (
        PanicLevel	Level = iota
        FatalLevel
        ErrorLevel	
        WarningLevel
        InfoLevel		
        DebugLevel	
)

var logLevel Level = InfoLevel

func SetOutput(logFile string) {
      
   // check where the log should be written
   // if log file is empty, using 'stderr' as log output
   if (len(logFile) == 0) {
       Infof("Use 'stderr' as log output\n")
       return
   }
   
   // Open the log file
    // check if file exist
    fileinfo, err := os.Stat(logFile)
    
    if ((err == nil) && (fileinfo.Mode().IsRegular() != true)) {
        Errorf("Log file: '%s' exists, but it is not regular file, use stderr for log\n", logFile)
        return
    }
    
    // Whatever the error it is, don't care, using logFile as log output    
    file, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0744)
    if (err != nil) {
        Errorf("Log file: '%s' can't be opened, use stderr for log\n", logFile)
        return    
    }
    
    log.SetOutput(file)
}

func SetLevel(level Level) {
    logLevel = level
}

func Printf(format string, v ...interface {}) {
    log.Printf(format, v...)
}

func Errorf(format string, v ...interface {}) {
    if (logLevel >= ErrorLevel) {
        log.Printf("ERROR: " + format, v...)
    }
}

func Warnf(format string, v ...interface {}) {
    if (logLevel >= WarningLevel) {
        log.Printf("WARN: " + format, v...)
    }
}

func Infof(format string, v ...interface {}) {
    if (logLevel >= InfoLevel) {
        log.Printf("INFO: " + format, v...)
    }
}

func Debugf(format string, v ...interface {}){
    if (logLevel >= DebugLevel) {
        log.Printf("DEBUG: " + format, v...)
    }
}

func ErrorBinary (data []byte) {
    if (logLevel >= ErrorLevel) {
        logBinary("ERROR: ", data)
    }    
}

func WarnBinary (data []byte) {
    if (logLevel >= WarningLevel) {
        logBinary("WARN: ", data)
    }    
}

func InfoBinary (data []byte) {
    if (logLevel >= InfoLevel) {
        logBinary("INFO: ", data)
    }    
}

func DebugBinary (data []byte) {
        
    if (logLevel >= DebugLevel) {
        logBinary("DEBUG: ", data)
    }    
}

func logBinary(level string, data []byte) {
   
    var line 	string = ""
    var linex	string = ""
    
    for i, v := range data {
        
        line += fmt.Sprintf("%02x ", v)
        
        if ((v > 32) && (v < 127)){
            linex += fmt.Sprintf("%c", v)
        } else {
            linex += "."
        }            
        
        if (((i +1) % 8) == 0) {
            line  += " "
            linex += " "
        }
        
        if (((i + 1) %16) == 0) {
            log.Printf(level + "%s%s\n", line, linex)
            line  	= ""
            linex 	= ""
        }
    }
    
    log.Printf("\n")
}