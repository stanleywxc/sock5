package address

import (
        "socks"
)

type Address interface {
        Atyp	()				(byte)
        DstAddr()			(string)
        DstPort()			(int)
        SetNetwork(string)
        GetNetwork()			(string)
}

type AddressIP struct {
        atyp			byte
        dstAddr		string
        dstPort		int
        network		string
}

/*----------------------------------------------------------
    Create an Address
-----------------------------------------------------------*/
func New(atyp byte, dstAddr string, dstPort int) (Address) {
    
    var address Address
    
    switch atyp {
        case socks.SOCKS_V5_ATYP_IP4:
            address = NewAddress("tcp4", atyp, dstAddr, dstPort)
            break
        case socks.SOCKS_V5_ATYP_FQDN:
            address = NewAddress("tcp4", atyp, dstAddr, dstPort)
            break
        case socks.SOCKS_V5_ATYP_IP6:
            address = NewAddress("tcp6", atyp, dstAddr, dstPort)
            break
    }
    
    return address
}
/*----------------------------------------------------------
    AddressIP4
-----------------------------------------------------------*/
func NewAddress(network string, atyp byte, dstAddr string, dstPort int) (*AddressIP) {
    return &AddressIP {atyp : atyp, dstAddr: dstAddr, dstPort : dstPort, network : network }
}

func (address *AddressIP) Atyp() (byte) {
    return address.atyp
}

func (address *AddressIP) DstAddr() (string) {
    return address.dstAddr
}

func (address *AddressIP) DstPort() (int) {
    return address.dstPort
}

func (address *AddressIP) SetNetwork(network string) {
    address.network = network
}

func (address *AddressIP) GetNetwork() (string) {
    return address.network
}