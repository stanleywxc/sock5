//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package command

import (
        "bufio"
        "io"
        "net"
        "strconv"
        "sync"
        "socks"
        "socks/log"
        "socks/address"
        "socks/context"
        
)

type Command interface {
        Execute 			()
}

type CommandConnect struct {
        connection		net.Conn
        waiter			sync.WaitGroup
        upstream			chan []byte
        downstream		chan	 []byte
        upstop			chan bool
        downstop			chan bool
        address			*address.Address
        context			*context.Context
}

type CommandBind struct {
        connection	net.Conn
        waiter		sync.WaitGroup    
        address		*address.Address
        context		*context.Context
}

type CommandUDPAssociation struct {
        connection	net.Conn
        waiter		sync.WaitGroup    
        address		*address.Address
        context		*context.Context
}

/*----------------------------------------------------------
    Connect Command
-----------------------------------------------------------*/
func NewCommandConnect (address *address.Address, context *context.Context) (*CommandConnect){
    
    return &CommandConnect { address : address, context : context}
}

func (command *CommandConnect) Execute () {

    // Reply the response with success code
    // try to connect to upstream/target host first
    connection, err := net.Dial((*command.address).GetNetwork(), net.JoinHostPort((*command.address).DstAddr(), strconv.Itoa((*command.address).DstPort())))    
    
    // Is there any error?
    if ( err != nil) {
        
        // Error happened, send error code back
        command.response(socks.SOCKS_V5_STATUS_HOST_UNREACHABLE, nil)
        log.Errorf("Connect to target: %s:%d failed\n", (*command.address).DstAddr(), (*command.address).DstPort())
        return 
    }

    log.Infof("Target host connected: %s via %s\n", connection.RemoteAddr().String(), connection.RemoteAddr().Network())
    
    var ipBytes []byte
    
    // check the atyp
    switch ((*command.address).Atyp()) {
        case socks.SOCKS_V5_ATYP_IP4:
            ipBytes = net.ParseIP((*command.address).DstAddr())
            break
        case socks.SOCKS_V5_ATYP_IP6:
            ipBytes = net.ParseIP((*command.address).DstAddr())
            break
        case socks.SOCKS_V5_ATYP_FQDN:
            ipBytes = append([]byte{byte(len((*command.address).DstAddr()))}, (*command.address).DstAddr()...)
            break
    }
    
    // Check if there is any error
    if (ipBytes == nil) {
        command.response(socks.SOCKS_V5_STATUS_ADDR_UNSUPPORTED, nil)
        connection.Close()
        
        log.Errorf("Target host IP is incorrect: %s\n", (*command.address).DstAddr())
        return 
    }
        
    ipBytes = append(ipBytes, byte((*command.address).DstPort() >> 8))
    ipBytes = append(ipBytes, byte((*command.address).DstPort() & 0xFF))
        
    // Send response
    command.response(socks.SOCKS_V5_STATUS_SUCCESS, ipBytes)
          
    // start to proxy
    command.connection 	= connection
    command.upstream 	= make(chan []byte)
    command.downstream	= make(chan []byte)
    command.upstop		= make(chan bool)
    command.downstop		= make(chan bool)
 
    // Start the proxy
    command.waiter.Add(4)
    
    go command.listenUpstream()
    go command.listenDownstream()

    // start the upstream proxy
    go command.upstreamProxy()
    
    //start the downstream proxy
    go command.downstreamProxy()
        
    log.Infof("Waiting for proxying to finish\n")
    
    // wait for proxying being done
    command.waiter.Wait()
 
    connection.Close()
    
    log.Infof("Proxying finished\n")

    // Done
    return 
}

func (command *CommandConnect) listenUpstream() {
    
    log.Infof("Entering listenUpstream\n")

    //command.waiter.Add(1)
    defer command.waiter.Done()
    
    var writer *bufio.Writer = command.context.Writer()
    
    // Wait for data coming
    var status bool = false
    for {
        select {
            case data := <-command.upstream :
                writer.Write(data)
                writer.Flush()
                log.Infof("Sending data to downstream(%s) - bytes: %d\n", (*command.context.Connection()).RemoteAddr().String(), len(data))   
            case status = <-command.upstop :
                log.Infof("Received upstream stop signal\n")
                break
        }
        if (status == true) {
            break
        }
    }
    
    log.Infof("Leaving listenUpstream\n")
    // Done
    return 
}

func (command *CommandConnect) listenDownstream() {

    log.Infof("Entering listenDownstream\n")
    
    defer command.waiter.Done()
    
    var writer *bufio.Writer = bufio.NewWriter(command.connection)
    
    // Wait for data coming
    var status bool = false
    for {
        
        select {
            case data := <-command.downstream :
                writer.Write(data)
                writer.Flush()
                log.Infof("Sending data to upstream(%s) - bytes: %d\n", command.connection.RemoteAddr().String(), len(data))  
            case status = <-command.downstop :
                log.Infof("Received downstream stop signal\n")
                break
        }
        
        if (status == true) {
            break
        }
    }

    log.Infof("Leaving listenDownstream\n")
    
    // Done
    return     
}

func (command *CommandConnect) upstreamProxy() {

    log.Infof("Entering upstreamProxy\n")
    
    defer command.waiter.Done()
    
    //io.Copy((*command.context.Connection()), command.connection)
    
    var reader *bufio.Reader = bufio.NewReader(command.connection)

    // Read data from upstream
    for {
        // clear the buffer
        var buffer []byte = make([]byte, reader.Size())
        
        // read the data from stream
        count, err := reader.Read(buffer)
        
        log.Infof("Reading data from upstream(%s): - bytes: %d, err: %#v\n", command.connection.RemoteAddr().String(), count, err)
        log.DebugBinary(buffer[:count])
       
        // any data read?
        if (count != 0) {
            command.upstream <- buffer[:count]
            
            log.Infof("Sending data to downstream(%s): - bytes: %d, err: %#v\n", (*command.context.Connection()).RemoteAddr().String(), count, err)
        }
        
        // Reach end of stream?
        if (err == io.EOF) {
            log.Infof("Sending upstream stop signal\n")
            command.upstop <- true
            break
        }
    }
    
    log.Infof("Leaving upstreamProxy\n")

    // Done
    return
}

func (command *CommandConnect) downstreamProxy() {

    log.Infof("Entering downstreamProxy\n")

    defer command.waiter.Done()
    
    //io.Copy(command.connection, (*command.context.Connection()))
    
    // Using the io.Copy to accomplish the data transfering is totally fine too.
    var reader *bufio.Reader	= command.context.Reader()
    
    // Read data from upstream
    for {       
        var buffer []byte = make([]byte, reader.Size())
        
        // read the data from stream
        count, err := reader.Read(buffer)
        
        log.Infof("Reading data from downstream(%s): - bytes: %d, err: %#v\n", (*command.context.Connection()).RemoteAddr().String(), count, err)
        log.DebugBinary(buffer[:count])
        
        // any data read?
        if (count != 0) {
            command.downstream <- buffer[:count]
            
            log.Infof("Sending data to upstream(%s): - bytes: %d, err: %#v\n", command.connection.RemoteAddr().String(), count, err)
        }
        
        // Reach at the end of stream?
        if (err == io.EOF) {
            log.Infof("Sending downstream stop signal\n")
            command.downstop <- true
            break
        }
    }
    
    log.Infof("Leaving downstreamProxy\n")

    // Done
    return
    
}

func (command *CommandConnect) response(statuscode byte, rest []byte) {
    
    // Send response back
    command.context.Writer().WriteByte(command.context.Version())
    command.context.Writer().WriteByte(statuscode)
    command.context.Writer().WriteByte(0x00)
    command.context.Writer().WriteByte((*command.address).Atyp())
    
    // In case there is an error, don't care the rest
    if ((rest != nil) && (len(rest) != 0)) {
        command.context.Writer().Write(rest)
    }
    
    // Flush the response
    command.context.Writer().Flush()
    
    // Done
    return 
}


/*----------------------------------------------------------
    Bind Command
-----------------------------------------------------------*/
func NewCommandBind (address *address.Address, context *context.Context) (*CommandBind) {
        
    return &CommandBind {  address : address, context : context }
}

func (command *CommandBind) Execute () {
    
    return 
}

/*----------------------------------------------------------
    Bind Command
-----------------------------------------------------------*/
func NewCommandUDPAssociation (address *address.Address, context *context.Context) (*CommandUDPAssociation) {
    return &CommandUDPAssociation {  address : address, context : context }
}

func (command *CommandUDPAssociation) Execute () {
    
    return 
}
