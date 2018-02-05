package sender

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

// Metric class.
type Metric struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int64  `json:"clock"`
}

// Metric class constructor.
func NewMetric(host, key, value string, clock ...int64) *Metric {
	m := &Metric{Host: host, Key: key, Value: value}
	// use current time, if `clock` is not specified
	if m.Clock = time.Now().Unix(); len(clock) > 0 {
		m.Clock = int64(clock[0])
	}
	return m
}

// Packet class.
type Packet struct {
	Request string    `json:"request"`
	Data    []*Metric `json:"data"`
	Clock   int64     `json:"clock"`
}

// Packet class cunstructor.
func NewPacket(data []*Metric, clock ...int64) *Packet {
	p := &Packet{Request: `sender data`, Data: data}
	// use current time, if `clock` is not specified
	if p.Clock = time.Now().Unix(); len(clock) > 0 {
		p.Clock = int64(clock[0])
	}
	return p
}

type DiscoveryPayload struct {
	Data []map[string]string `json:"data"`
}

func NewDiscoveryMetric(host, key string, value []map[string]string, clock ...int64) *Metric {
	payload := &DiscoveryPayload{Data: value}
	jsonString, err := json.Marshal(&payload)

	if err != nil {
		fmt.Errorf("Error marshaling: %s", err.Error)
	}

	m := &Metric{Host: host, Key: key, Value: string(jsonString)}

	// use current time, if `clock` is not specified
	if m.Clock = time.Now().Unix(); len(clock) > 0 {
		m.Clock = int64(clock[0])
	}
	return m
}

// DataLen Packet class method, return 8 bytes with packet length in little endian order.
func (p *Packet) DataLen() []byte {
	dataLen := make([]byte, 8)
	JSONData, _ := json.Marshal(p)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(JSONData)))
	return dataLen
}

// Response struct
type Response struct {
	Response string `json:"response"`
	Info     string `json:"info"`
}

// Response constructor
func NewResponse(data []byte) (response *Response, err error) {
	if err = json.Unmarshal(data[13:], &response); err != nil {
		return
	}
	return
}
