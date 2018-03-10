package main

import (
        "net"
        "strconv"
        "socks/log"
        "socks/config"
        "socks/context"    
        "socks/session"
)

type Server struct {
    config	*config.Config
}

func New(config *config.Config) (*Server) {
    return &Server{config : config}
}

func (server *Server) Start () (bool) {
    
    listener, err := net.Listen(server.config.Server.Protocol, net.JoinHostPort(server.config.Server.Address , strconv.Itoa(server.config.Server.Listen)))
    
    if (err != nil) {
        log.Errorf("Error : %s", err)
        return false
    }
   
    // Start to accept incoming connections
    for {
        
        // Accept incoming connections
        connection, err := listener.Accept()
        
        if (err != nil) {
            
            log.Errorf("Error in accepting incoming connection: %s\n", err.Error())
            continue
        }
        
        // Is it a valid connection?
        if (connection == nil) {
            log.Errorf("Error in accepting incoming connection: connection is nil\n")            
            continue
        }
 
        // Handle the incoming connections.
        go server.handleIncoming(connection)
    }
}

func (server *Server) handleIncoming (conn net.Conn) {
    
    log.Infof("Incomming: %s, Remote Addr: %s\n", conn.LocalAddr().Network(), conn.RemoteAddr().String())
    
    // create the context 
    contxt, err := context.New(conn, server.config)
    
    // Check the context is valid
    if (contxt == nil) {
        log.Errorf("error during create Context: %s\n", err.Error())
        conn.Close()
        return
    }
    
    log.Infof("Context created, version: %d\n", contxt.Version())
        
    // create a session
    var sessionCtxt *session.SessionContext = session.New(contxt)
    
    // start the session
    err = sessionCtxt.Start()
    
    // If there is any error, log it
    if (err != nil) {
        log.Errorf("Session create failed: %s\n", err.Error())
    }

    // Done
    conn.Close()
}