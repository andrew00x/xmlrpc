package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type unmarshaller struct {
	last *xml.Token
}

type valueElement struct {
	Data string `xml:",chardata"`
}

type fault struct {
	code int
	text string
}

const any = ""

func (u *unmarshaller) unmarshal(b []byte) (params []interface{}, err error) {
	d := xml.NewDecoder(bytes.NewReader(b))
	var se *xml.StartElement
	if se, err = u.startElement(d, "methodResponse"); err != nil || se == nil {
		return
	}
	if se, err = u.startElement(d, any); err != nil {
		return
	}
	name := se.Name.Local
	switch name {
	case "params":
		if params, err = u.unmarshalParams(d); err != nil {
			return
		}
		_, err = u.mustEndElement(d, "params")
	case "fault":
		var f fault
		if f, err = u.unmarshalFault(d); err != nil {
			return
		}
		if _, err = u.mustEndElement(d, "fault"); err != nil {
			return
		}
		err = errors.New(fmt.Sprintf("error response, code: %d, text: %s", f.code, f.text))
	default:
		err = errors.New(fmt.Sprintf("invalid xml, unknown element %s", name))
		return
	}
	if err == nil {
		_, err = u.mustEndElement(d, "methodResponse")
	}
	return
}

func (u *unmarshaller) unmarshalParams(d *xml.Decoder) (params []interface{}, err error) {
	var se *xml.StartElement
	for {
		if se, err = u.startElement(d, "param"); err != nil || se == nil {
			return
		}
		var v interface{}
		v, err = u.unmarshalValue(d)
		if err != nil {
			return
		}
		params = append(params, v)
		if _, err = u.mustEndElement(d, "param"); err != nil {
			return
		}
	}
}

func (u *unmarshaller) unmarshalFault(d *xml.Decoder) (f fault, err error) {
	var v interface{}
	v, err = u.unmarshalValue(d)
	fm, ok := v.(map[string]interface{})
	if !ok {
		err = errors.New(fmt.Sprintf("unable parse fault message correctly: %v", v))
	}
	fcode, ok := fm["faultCode"]
	if !ok {
		err = errors.New(fmt.Sprintf("no code in fault message: %v", fm))
	}
	fmsg, ok := fm["faultString"]
	if !ok {
		err = errors.New(fmt.Sprintf("no code in fault message: %v", fm))
	}
	code := fcode.(int)
	msg := fmsg.(string)
	f = fault{code: code, text: msg}
	return
}

func (u *unmarshaller) unmarshalValue(d *xml.Decoder) (v interface{}, err error) {
	var se *xml.StartElement
	var vn valueElement
	if se, err = u.startElement(d, "value"); err != nil || se == nil {
		return
	}
	if se, err = u.startElement(d, any); err != nil {
		return
	}
	name := se.Name.Local
	switch name {
	case "string", "base64", "int", "i4", "i8", "boolean":
		if err = d.DecodeElement(&vn, se); err != nil {
			return
		}
		v, err = decodeValue(vn.Data, name)
		u.last = nil
	case "array":
		v, err = u.unmarshalArray(d)
		_, err = u.mustEndElement(d, "array")
	case "struct":
		v, err = u.unmarshalStruct(d)
		_, err = u.mustEndElement(d, "struct")
	default:
		err = errors.New(fmt.Sprintf("unsupported type: %s", name))
	}
	if err == nil {
		_, err = u.mustEndElement(d, "value")
	}
	return
}

func (u *unmarshaller) unmarshalArray(d *xml.Decoder) (arr []interface{}, err error) {
	var se *xml.StartElement
	if se, err = u.startElement(d, "data"); err != nil || se == nil {
		return
	}
	for {
		var v interface{}
		v, err = u.unmarshalValue(d)
		if v == nil {
			break
		}
		arr = append(arr, v)
	}
	_, err = u.mustEndElement(d, "data")
	return
}

func (u *unmarshaller) unmarshalStruct(d *xml.Decoder) (m map[string]interface{}, err error) {
	var se *xml.StartElement
	var vn valueElement
	m = make(map[string]interface{})
	for {
		if se, err = u.startElement(d, "member"); err != nil || se == nil {
			return
		}
		if se, err = u.startElement(d, "name"); err != nil {
			return
		}
		if err = d.DecodeElement(&vn, se); err != nil {
			return
		}
		n := vn.Data
		u.last = nil
		var v interface{}
		if v, err = u.unmarshalValue(d); err != nil {
			return
		}
		m[n] = v
		_, err = u.mustEndElement(d, "member")
	}
}

func (u *unmarshaller) startElement(d *xml.Decoder, name string) (se *xml.StartElement, err error) {
	var t xml.Token
	if u.last != nil {
		t = *u.last
		u.last = nil
	}
	for se == nil {
		switch e := t.(type) {
		case xml.StartElement:
			if name == any || e.Name.Local == name {
				se = &e
			} else {
				u.last = &t
			}
			return
		case xml.EndElement:
			u.last = &t
			return
		}
		t, err = d.Token()
		if t == nil && err == io.EOF {
			err = nil
			return
		}
		if err != nil {
			return
		}
	}
	return
}

func (u *unmarshaller) endElement(d *xml.Decoder, name string) (ee *xml.EndElement, err error) {
	var t xml.Token = u.last
	if u.last != nil {
		t = *u.last
		u.last = nil
	}
	for ee == nil {
		switch e := t.(type) {
		case xml.EndElement:
			if name == any || e.Name.Local == name {
				ee = &e
			} else {
				u.last = &t
			}
			return
		case xml.StartElement:
			u.last = &t
			return
		}
		t, err = d.Token()
		if t == nil && err == io.EOF {
			err = nil
			return
		}
		if err != nil {
			return
		}
	}
	return
}

func (u *unmarshaller) mustEndElement(d *xml.Decoder, name string) (ee *xml.EndElement, err error) {
	if ee, err = u.endElement(d, name); err != nil {
		return
	}
	if ee == nil {
		err = errors.New(fmt.Sprintf("missing end element: %s", name))
		return
	}
	return
}

func decodeValue(raw string, t string) (v interface{}, err error) {
	switch t {
	case "string":
		v = raw
	case "base64":
		v, err = base64.StdEncoding.DecodeString(raw)
	case "i4", "int":
		var i64 int64
		if i64, err = strconv.ParseInt(raw, 10, 32); err == nil {
			v = int(i64)
		}
	case "i8":
		v, err = strconv.ParseInt(raw, 10, 64)
	case "boolean":
		v, err = strconv.ParseBool(raw)
	}
	return
}
