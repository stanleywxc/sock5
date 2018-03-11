//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package request

import (
        "errors"
        "net"
        "socks"
        "socks/address"
        "socks/command"
        "socks/context"
)

type Request interface {
        Start	() (bool, error)
        Command	() (*command.Command)
}

type RequestV4 struct {
        version			byte
        commandIndex		byte
        reserved			byte
        atyp				byte
        address			address.Address
        command			command.Command
        context			*context.Context
}

type RequestV5 struct {
        version			byte
        commandIndex		byte
        reserved			byte
        atyp				byte
        address			address.Address
        command			command.Command
        context			*context.Context    
}

/*----------------------------------------------------------
   Create a request
-----------------------------------------------------------*/
func New(context *context.Context) (Request) {
    
    var request Request
    
    // create Request based on Socks V4 or V5
    switch context.Version() {
        case socks.SOCKS_VERSION_V4:
            request = newRequestV4(context)
            break
        case socks.SOCKS_VERSION_V5:
            request = newRequestV5(context)
            break
    }
    
    return request
}

/*----------------------------------------------------------
    Handle Socks V5 requests
-----------------------------------------------------------*/
/* RFC 1928
4.  Requests

   Once the method-dependent subnegotiation has completed, the client
   sends the request details.  If the negotiated method includes
   encapsulation for purposes of integrity checking and/or
   confidentiality, these requests MUST be encapsulated in the method-
   dependent encapsulation.

   The SOCKS request is formed as follows:

        +----+-----+-------+------+----------+----------+
        |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
        +----+-----+-------+------+----------+----------+
        | 1  |  1  | X'00' |  1   | Variable |    2     |
        +----+-----+-------+------+----------+----------+

     Where:

          o  VER    protocol version: X'05'
          o  CMD
             o  CONNECT X'01'
             o  BIND X'02'
             o  UDP ASSOCIATE X'03'
          o  RSV    RESERVED
          o  ATYP   address type of following address
             o  IP V4 address: X'01'
             o  DOMAINNAME: X'03'
             o  IP V6 address: X'04'
          o  DST.ADDR       desired destination address
          o  DST.PORT desired destination port in network octet
             order

   The SOCKS server will typically evaluate the request based on source
   and destination addresses, and return one or more reply messages, as
   appropriate for the request type.





Leech, et al                Standards Track                     [Page 4]

RFC 1928                SOCKS Protocol Version 5              March 1996


5.  Addressing

   In an address field (DST.ADDR, BND.ADDR), the ATYP field specifies
   the type of address contained within the field:

          o  X'01'

   the address is a version-4 IP address, with a length of 4 octets

          o  X'03'

   the address field contains a fully-qualified domain name.  The first
   octet of the address field contains the number of octets of name that
   follow, there is no terminating NUL octet.

          o  X'04'

   the address is a version-6 IP address, with a length of 16 octets.

6.  Replies

   The SOCKS request information is sent by the client as soon as it has
   established a connection to the SOCKS server, and completed the
   authentication negotiations.  The server evaluates the request, and
   returns a reply formed as follows:

        +----+-----+-------+------+----------+----------+
        |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
        +----+-----+-------+------+----------+----------+
        | 1  |  1  | X'00' |  1   | Variable |    2     |
        +----+-----+-------+------+----------+----------+

     Where:

          o  VER    protocol version: X'05'
          o  REP    Reply field:
             o  X'00' succeeded
             o  X'01' general SOCKS server failure
             o  X'02' connection not allowed by ruleset
             o  X'03' Network unreachable
             o  X'04' Host unreachable
             o  X'05' Connection refused
             o  X'06' TTL expired
             o  X'07' Command not supported
             o  X'08' Address type not supported
             o  X'09' to X'FF' unassigned
          o  RSV    RESERVED
          o  ATYP   address type of following address



Leech, et al                Standards Track                     [Page 5]

RFC 1928                SOCKS Protocol Version 5              March 1996


             o  IP V4 address: X'01'
             o  DOMAINNAME: X'03'
             o  IP V6 address: X'04'
          o  BND.ADDR       server bound address
          o  BND.PORT       server bound port in network octet order

   Fields marked RESERVED (RSV) must be set to X'00'.

   If the chosen method includes encapsulation for purposes of
   authentication, integrity and/or confidentiality, the replies are
   encapsulated in the method-dependent encapsulation.
*/
func (request *RequestV5) Start () (bool, error) {
    
    var err error
    
    // Read version, command, rsv, atyp
    request.version 		, err = request.context.Reader().ReadByte()
    request.commandIndex	, err = request.context.Reader().ReadByte()
    request.reserved		, err = request.context.Reader().ReadByte()
    request.atyp			, err = request.context.Reader().ReadByte()
    
    if (err != nil) {
       return false, err 
    }
    
    request.address, err = request.getAddress()
    
    if (err != nil) {
        return false, err
    }
    
    _, err = request.getCommand()
    
    return err == nil, err
}

func (request *RequestV5) Command() (*command.Command) {
    return &request.command
}

func newRequestV5(context *context.Context) (*RequestV5) {
    return &RequestV5 { context : context}
}

/*----------------------------------------------------------
    Handle Socks V4 requests
-----------------------------------------------------------*/
func (request *RequestV4) Start () (bool, error) {
    
    return true, nil
}

func (request *RequestV4) Command() (*command.Command) {
    return &request.command
}

func newRequestV4(context *context.Context) (*RequestV4) {
    return &RequestV4 { context : context}
}

/*----------------------------------------------------------
    private method
-----------------------------------------------------------*/
func (request *RequestV5) getAddress() (address.Address, error) {
    
    var ipaddress	string
    var port			int
    var err			error
    
    switch request.atyp {
        case socks.SOCKS_V5_ATYP_IP4:
            var ipbytes []byte = make([]byte, 4)
               _, err = request.context.Reader().Read(ipbytes)
               
            // If there is no error, continue
            if (err == nil) {
                var bite byte
                bite, err = request.context.Reader().ReadByte()
                port |= (int(bite) << 8)
                bite, err = request.context.Reader().ReadByte()
                port |= int(bite)                
                ipaddress = net.IP(ipbytes).String()
            }
            break
        case socks.SOCKS_V5_ATYP_FQDN:
            var count byte
            count, err = request.context.Reader().ReadByte()
            
            // If there is no error, continue
            if (err == nil) {                
                var ipFQDN []byte = make([]byte, count)
                _, err = request.context.Reader().Read(ipFQDN)
                var bite byte
                bite, err = request.context.Reader().ReadByte()
                port |= (int(bite) << 8)
                bite, err = request.context.Reader().ReadByte()
                port |= int(bite)                
                ipaddress = string(ipFQDN[:])
            }
            break
        case socks.SOCKS_V5_ATYP_IP6:
            var ipbytes []byte = make([]byte, 16)
               _, err = request.context.Reader().Read(ipbytes)
               
            // If there is no error, continue
            if (err == nil) {
                var bite byte
                bite, err = request.context.Reader().ReadByte()
                port |= (int(bite) << 8)
                bite, err = request.context.Reader().ReadByte()
                port |= int(bite)                
                ipaddress = net.IP(ipbytes).String()
            }
            
            break
        default :
            err = errors.New("No supported address")
            break
    }
    
    if (err != nil) {
        return nil, err
    }
  
    if (request.atyp != socks.SOCKS_V5_ATYP_FQDN) {
    
        // The atyp in reply has to be SOCKS_V5_ATYP_FQDN, otherwise the certificate returned by 
        // target server will not be able to be verified, and thus cause https handshake failure.
        var hosts, err = net.LookupAddr(ipaddress)
        
        // any error? If there is any error, don't convert
        if (err == nil) {
                 
            // Check if the returned hosts is empty
            // if it is empty don't convert to FQDN
            if (len(hosts) != 0) {
                
                //it is not empty, convert the atyp to SOCKS_V5_ATYP_FQDN in reply
                ipaddress 	= hosts[0]
                request.atyp	= socks.SOCKS_V5_ATYP_FQDN
            }
        }
    }
    
    // create the destination address object
    request.address = address.New(request.atyp, ipaddress, int(port))
    
    return request.address, nil
}

func (request *RequestV5) getCommand() (*command.Command, error) {
    
    var err error
    
    switch (request.commandIndex) {
        case socks.SOCKS_COMMAND_CONNECT:
            request.command = command.NewCommandConnect(&request.address, request.context)
            break
        case socks.SOCKS_COMMAND_BIND:
            request.command = command.NewCommandBind(&request.address, request.context)
            break
        case socks.SOCKS_COMMAND_UDP_ASSOCIATE:
            request.command = command.NewCommandUDPAssociation(&request.address, request.context)
            break
        default:
            err = errors.New("Not supported command")
    }
    
    return &request.command, err
}
