package xmlrpc

import (
	"bytes"
	"github.com/mpl/scgiclient"
)

type Client interface {
	Send(method string, args ...interface{}) (params []interface{}, err error)
}

type SCGIXmlRpc struct {
	Addr string
	marshaller
	unmarshaller
}

func CreateSCGIClient(addr string) Client {
	return &SCGIXmlRpc{Addr: addr}
}

func (s *SCGIXmlRpc) Send(method string, args ...interface{}) (params []interface{}, err error) {
	body, err := s.marshal(method, args...)
	if err != nil {
		return nil, err
	}
	resp, err := scgiclient.Send(s.Addr, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return s.unmarshal(resp.Body)
}
