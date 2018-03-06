package socks

// Socks Version define
const (
    SOCKS_VERSION_V4	 	= byte(0x04)
    SOCKS_VERSION_V5		= byte(0x05)
)

/* RFC 1928
          o  X'00' NO AUTHENTICATION REQUIRED
          o  X'01' GSSAPI
          o  X'02' USERNAME/PASSWORD
          o  X'03' to X'7F' IANA ASSIGNED
          o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
          o  X'FF' NO ACCEPTABLE METHODS

*/
const (
    SOCKS_AUTH_NOAUTHENTICATION	= byte(0x00)
    SOCKS_AUTH_GSSAPI			= byte(0x01)
    SOCKS_AUTH_USERPASSWORD		= byte(0x02)
    SOCKS_AUTH_NOACCEPTABLE		= byte(0xFF) 
)

const (
    SOCKS_V5_ATYP_IP4			= byte(0x01)
    SOCKS_V5_ATYP_FQDN			= byte(0x03)
    SOCKS_V5_ATYP_IP6			= byte(0x04)
)

const (
    SOCKS_COMMAND_CONNECT		= byte(0x01)
    SOCKS_COMMAND_BIND			= byte(0x02)
    SOCKS_COMMAND_UDP_ASSOCIATE	= byte(0x03)
)

/* RFC 1928
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
*/
const (
    SOCKS_V5_STATUS_SUCCESS				= byte(0x00)
    SOCKS_V5_STATUS_SERVER_FAILURE		= byte(0x01)
    SOCKS_V5_STATUS_NOT_ALLOWED			= byte(0x02)
    SOCKS_V5_STATUS_NETWORK_UNREACHABLE	= byte(0x03)
    SOCKS_V5_STATUS_HOST_UNREACHABLE		= byte(0x04)
    SOCKS_V5_STATUS_CONN_REFUSED			= byte(0x05)
    SOCKS_V5_STATUS_TTL_EXPIRED			= byte(0x06)
    SOCKS_V5_STATUS_COMMAND_UNSUPPORTED	= byte(0x07)
    SOCKS_V5_STATUS_ADDR_UNSUPPORTED		= byte(0x08)
    SOCKS_V5_STATUS_UNASSIGNED			= byte(0xff)
)
