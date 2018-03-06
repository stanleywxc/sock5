package context

import (
    "fmt"
    "bufio"
    "errors"
    "net"
    "socks"
    
)

type Context struct {
    version		byte
    connection 	net.Conn
    reader		*bufio.Reader
    writer		*bufio.Writer
}

func New(conn net.Conn) (*Context, error) {
    
    var reader *bufio.Reader = bufio.NewReader(conn)
    var writer *bufio.Writer = bufio.NewWriter(conn)
    
    // Read the version
    var version 	byte
    var err		error
    version, err = reader.ReadByte()
   
    // if version is not presented or error, then return nil, and err
    if (err != nil) {
        return nil, err
    }
    
    fmt.Printf("version: %d\n", version)
    
    // Check version is supported?
    if ((version != socks.SOCKS_VERSION_V4) && (version != socks.SOCKS_VERSION_V5)) {
        return nil, errors.New("Version is not supported")
    }
    
    // Rewind the position of version
    reader.UnreadByte();
    
    // return the Context object
    return &Context {   version		: version,
                        connection 	: conn,
                        reader		: reader,
                        writer		: writer }, nil
}

func (contxt *Context)Version() byte {
    return contxt.version
}

func (context *Context) Reader() (*bufio.Reader) {
    return context.reader
}

func (context *Context) Writer() (*bufio.Writer) {
    return context.writer
}

func (context *Context) LocalAddr() (string) {
    return context.connection.LocalAddr().String()
}