package config

import (
        "os"
        "os/user"
        "encoding/json"
        "io/ioutil"
        "path/filepath"
        "socks/log"
)

const	DEFAULT_CONF_FILE = "socks5.conf"

type Config struct {
    Daemon	bool
    Server	ServerConf
    Auth		AuthConf
    Log		LogConf
}

type ServerConf	struct {
    Protocol		string
    Address		string
    Listen		int
}

type AuthConf struct {
    Username		string
    Password		string
}

type LogConf	 struct {
    Level	int
    Path		string
}

func readConf(path string) (*Config) {

    var err		error
    var config 	Config = Config{}
    
    // looking for 'socks5.conf' file
    if (len(path) == 0) {
        return nil
    }
        
    // check if file exist
    fileinfo, err := os.Stat(path)
        
    if (err != nil) {
        log.Errorf("Stat conf file: %s\n", err.Error())    
        return nil
    }
        
    if (fileinfo.Mode().IsRegular() != true) {
        log.Errorf("Conf file is not readable\n")    
        return nil        
    }
 
    bytes, err := ioutil.ReadFile(path)
    
    if (err != nil) {
        log.Errorf("Reading conf file: %s\n", err.Error())
        return nil
    }
    
    if err = json.Unmarshal(bytes, &config); err != nil {
        log.Errorf("Unmarshal conf file: %s\n", err.Error())
        return nil
    }       
    
    return &config
}

// Initialization always returns a Config
func Initialize (path string) (*Config) {
    
    var config 	*Config = nil
    
    // looking for 'socks5.conf' file
    config = readConf(path)
    if (config != nil) {
        return config
    }
    
    var home string = ""
    
    // Find the home directory of the user that socks5 is running
    usr, err := user.Current()
    if (err == nil) {
        home = filepath.Join(usr.HomeDir, ".socks5")
    }
    
    // path doesn't exist, trying to find conf by the following order:
    // 1. ~/.socks5/socks5.conf
    // 2. /usr/local/etc/socks5.conf
    // 3. /etc/socks5.conf
    var paths []string = []string{filepath.Join(home, DEFAULT_CONF_FILE), filepath.Join("/usr/local/etc", DEFAULT_CONF_FILE), filepath.Join("/etc", DEFAULT_CONF_FILE)}
    
    for _, path := range paths {
        config = readConf(path)
        if (config != nil) {
            return config
        }
    }
    
    // Do we find conf file?
    if config == nil {
        // No there is no conf file
        // set it to default value
        config = &Config { Daemon: false, Server: ServerConf{Protocol: "tcp", Address: "", Listen: 1080}, Auth: AuthConf{Username: "", Password: ""}, Log: LogConf{Level: 1, Path: "~/tmp"}}
    }
    
    // Done
    return config
}