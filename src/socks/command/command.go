package command

import (
        "fmt"
        "bufio"
        "io"
        "net"
        "strconv"
        "sync"
        "socks"
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
        fmt.Printf("Connect to target: %s:%d failed\n", (*command.address).DstAddr(), (*command.address).DstPort())
        return 
    }

    fmt.Printf("Target host connected: %s via %s\n", connection.RemoteAddr().String(), connection.RemoteAddr().Network())
    
    // Split the IP address and port
    targetHost, targetPort, err := net.SplitHostPort(command.context.LocalAddr())
    
    if (err != nil) {
        command.response(socks.SOCKS_V5_STATUS_ADDR_UNSUPPORTED, nil)
        fmt.Printf("Target host IP is incorrect: %s\n", targetHost)
        return         
    }
    
    // parse the IP address into net.IP
    ipPort, _ := strconv.Atoi(targetPort)
    ipBytes := net.ParseIP(targetHost)
    
    // Check if there is any error
    if (ipBytes == nil) {
        command.response(socks.SOCKS_V5_STATUS_ADDR_UNSUPPORTED, nil)
        fmt.Printf("Target host IP is incorrect: %s\n", targetHost)
        return 
    }
    
    fmt.Printf("ipBytes length: %d\n", len(ipBytes))
    
    ipBytes = append(ipBytes, byte(ipPort >> 8))
    ipBytes = append(ipBytes, byte(ipPort & 0xFF))
    
    fmt.Printf("ipBytes length: %d\n", len(ipBytes))
    
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
        
    fmt.Printf("Waiting for proxying to finish\n")
    
    // wait for proxying being done
    command.waiter.Wait()
 
    connection.Close()
    
    fmt.Printf("Proxying finished\n")

    // Done
    return 
}

func (command *CommandConnect) listenUpstream() {
    
    fmt.Printf("Entering listenUpstream\n")

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
                fmt.Printf("-------sent data to downstream from: %s -- %d-------\n%s\n----------------------------\n", command.connection.RemoteAddr().String(), len(data), string(data[:]))   
            case status = <-command.upstop :
                fmt.Printf("Received upstop\n")
                break
        }
        if (status == true) {
            break
        }
    }
    
    fmt.Printf("Leaving listenUpstream\n")
    // Done
    return 
}

func (command *CommandConnect) listenDownstream() {

    fmt.Printf("Entering listenDownstream\n")
    
    //command.waiter.Add(1)
    defer command.waiter.Done()
    
    var writer *bufio.Writer = bufio.NewWriter(command.connection)
    
    // Wait for data coming
    var status bool = false
    for {
        
        select {
            case data := <-command.downstream :
                writer.Write(data)
                writer.Flush()
                fmt.Printf("-------sent data to upstream: %s -- %d -------\n%s\n----------------------------\n", command.connection.RemoteAddr().String(), len(data), string(data[:]))   
            case status = <-command.downstop :
                fmt.Printf("received downstop\n")
                break
        }
        
        if (status == true) {
            break
        }
    }

    fmt.Printf("Leaving listenDownstream\n")
    
    // Done
    return     
}

func (command *CommandConnect) upstreamProxy() {

    fmt.Printf("Entering upstreamProxy\n")
    
    //command.waiter.Add(1)
    defer command.waiter.Done()
    
    var reader *bufio.Reader = bufio.NewReader(command.connection)
    
    // Read data from upstream
    for {
        // clear the buffer
        var buffer []byte = make([]byte, reader.Size())
        
        // read the data from stream
        count, err := reader.Read(buffer)

        
        var rbuffer []byte = make([]byte, count)
        copy(rbuffer, buffer)
        
        fmt.Printf("------- read data from upstream: %d -------\n%s\n----------------------------\n", count, string(rbuffer[:]))
        
        // any data read?
        if (count != 0) {
            command.upstream <- rbuffer
        }
        
        // Reach end of stream?
        if (err == io.EOF) {
            fmt.Printf("Sending upstop\n")
            command.upstop <- true
            break
        }
    }
    
    fmt.Printf("Leaving upstreamProxy\n")

    // Done
    return
}

func (command *CommandConnect) downstreamProxy() {

    fmt.Printf("Entering downstreamProxy\n")

    //command.waiter.Add(1)
    defer command.waiter.Done()
    
    var reader *bufio.Reader	= command.context.Reader()
    
    // Read data from upstream
    for {
        // clear the buffer
        var buffer []byte = make([]byte, reader.Size())
        
        // read the data from stream
        count, err := reader.Read(buffer)
        
        var wbuffer []byte = make([]byte, count)
        copy(wbuffer, buffer)
 
        fmt.Printf("------- read data from downstream: %d -------\n%s\n----------------------------\n", count, string(wbuffer[:]))
       
        // any data read?
        if (count != 0) {
            command.downstream <- wbuffer
        }
        
        // Reach end of stream?
        if (err == io.EOF) {
            fmt.Printf("Sending downstop\n")
            command.downstop <- true
            break
        }
    }
    
    fmt.Printf("Leaving downstreamProxy\n")

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
