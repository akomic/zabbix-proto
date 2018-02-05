package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"
	"zabbix-proto/activechecks"
	"zabbix-proto/sender"
)

type Client struct {
	Host string `json:"host"`
	Port int
}

func NewClient(Host string, Port int) *Client {
	c := &Client{Host: Host, Port: Port}
	return c
}

// Method Client class, return zabbix header.
func (c *Client) getHeader() []byte {
	return []byte("ZBXD\x01")
}

// Method Client class, resolve uri by name:port.
func (c *Client) getTCPAddr() (iaddr *net.TCPAddr, err error) {
	// format: hostname:port
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	// Resolve hostname:port to ip:port
	iaddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		err = fmt.Errorf("Connection failed: %s", err.Error())
		return
	}

	return
}

// Method Client class, make connection to Zabbix Server.
func (c *Client) connect() (conn *net.TCPConn, err error) {

	type DialResp struct {
		Conn  *net.TCPConn
		Error error
	}

	// Open connection to zabbix host
	iaddr, err := c.getTCPAddr()
	if err != nil {
		return
	}

	// dial tcp and handle timeouts
	ch := make(chan DialResp)

	go func() {
		conn, err = net.DialTCP("tcp", nil, iaddr)
		ch <- DialResp{Conn: conn, Error: err}
	}()

	select {
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("Connection Timeout")
	case resp := <-ch:
		if resp.Error != nil {
			err = resp.Error
			break
		}

		conn = resp.Conn
	}

	return
}

// Method Client class, read data from connection.
func (c *Client) read(conn *net.TCPConn) (res []byte, err error) {
	res = make([]byte, 0, 4096)
	res, err = ioutil.ReadAll(conn)
	if err != nil {
		err = fmt.Errorf("Error whule receiving the data: %s", err.Error())
		return
	}

	return
}

// Method Client, Retrieve Zabbix Agent Active checks from server
// Variadic arguments: Host_metadata
func (c *Client) GetActiveItems(host string, params ...string) (items *activechecks.Response, err error) {
	conn, err := c.connect()
	if err != nil {
		return
	}
	defer conn.Close()

	packet := activechecks.NewPacket(host, params...)
	dataPacket, _ := json.Marshal(packet)

	// fmt.Printf("HEADER: % x (%s)\n", c.getHeader(), c.getHeader())
	// fmt.Printf("DATALEN: % x, %d byte\n", packet.DataLen(), len(packet.DataLen()))
	// fmt.Printf("BODY: %s\n", string(dataPacket))

	// Fill buffer
	buffer := append(c.getHeader(), packet.DataLen()...)
	buffer = append(buffer, dataPacket...)

	// Sent packet to zabbix
	_, err = conn.Write(buffer)
	if err != nil {
		err = fmt.Errorf("Error while sending the data: %s", err.Error())
		return
	}

	var res []byte
	res, err = c.read(conn)
	if err != nil {
		fmt.Errorf("Error while reading data: %s", err.Error())
	}
	items, err = activechecks.NewResponse(res)
	if err != nil {
		fmt.Errorf("Error decoding response from Zabbix", err.Error)
	}
	return
}

// Method Sender class, send packet to zabbix.
func (c *Client) Send(packet *sender.Packet) (response *sender.Response, err error) {
	conn, err := c.connect()
	if err != nil {
		return
	}
	defer conn.Close()

	dataPacket, _ := json.Marshal(packet)

	// fmt.Printf("HEADER: % x (%s)\n", c.getHeader(), c.getHeader())
	// fmt.Printf("DATALEN: % x, %d byte\n", packet.DataLen(), len(packet.DataLen()))
	// fmt.Printf("BODY: %s\n", string(dataPacket))

	// Fill buffer
	buffer := append(c.getHeader(), packet.DataLen()...)
	buffer = append(buffer, dataPacket...)

	_, err = conn.Write(buffer)
	if err != nil {
		err = fmt.Errorf("Error while sending the data: %s", err.Error())
		return
	}

	var res []byte
	res, err = c.read(conn)
	if err != nil {
		fmt.Errorf("Error while reading data: %s", err.Error())
	}

	response, err = sender.NewResponse(res)
	if err != nil {
		fmt.Errorf("Error decoding response from Zabbix", err.Error)
	}

	return
}
