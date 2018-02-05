package activechecks

import (
	"encoding/binary"
	"encoding/json"
)

// Packet class.
type Packet struct {
	Request       string `json:"request"`
	Host          string `json:"host"`
	Host_metadata string `json:"host_metadata,omitempty"`
}

// Packet class cunstructor.
func NewPacket(host string, params ...string) *Packet {
	p := &Packet{Request: `active checks`, Host: host}
	if len(params) >= 1 {
		p.Host_metadata = params[0]
	}
	return p
}

func (p *Packet) DataLen() []byte {
	dataLen := make([]byte, 8)
	JSONData, _ := json.Marshal(p)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(JSONData)))
	return dataLen
}

// Response struct
type Response struct {
	Response string `json:"response"`
	Data     []struct {
		Key         string `json:"key"`
		Delay       int    `json:"delay"`
		Lastlogsize int    `json:"lastlogsize"`
		Mtime       int    `json:"mtime"`
	} `json:"data"`
}

// Response constructor
func NewResponse(data []byte) (response *Response, err error) {
	if err = json.Unmarshal(data[13:], &response); err != nil {
		return
	}
	return
}
