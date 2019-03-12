package xmlrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

type marshaller struct {
}

func (m *marshaller) marshal(method string, args ...interface{}) (xml []byte, err error) {
	xmlWr := bytes.NewBufferString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	if _, err = xmlWr.WriteString("<methodCall><methodName>"); err != nil {
		return
	}
	if _, err = xmlWr.WriteString(method); err != nil {
		return
	}
	if _, err = xmlWr.WriteString("</methodName>"); err != nil {
		return
	}
	if _, err = xmlWr.WriteString("<params>"); err != nil {
		return
	}
	for _, arg := range args {
		if _, err = xmlWr.WriteString("<param>"); err != nil {
			return
		}
		if err = marshalValue(xmlWr, arg); err != nil {
			return
		}
		if _, err = xmlWr.WriteString("</param>"); err != nil {
			return
		}
	}
	_, err = xmlWr.WriteString("</params></methodCall>")

	xml = xmlWr.Bytes()
	return
}

func marshalValue(buf *bytes.Buffer, v interface{}) (err error) {
	if _, err = buf.WriteString("<value>"); err != nil {
		return
	}
	err = marshalType(buf, v)
	if err == nil {
		_, err = buf.WriteString("</value>")
	}
	return
}

func marshalType(buf *bytes.Buffer, v interface{}) (err error) {
	t := reflect.TypeOf(v)
	b, ok := v.([]byte)
	if ok {
		_, err = buf.WriteString(fmt.Sprintf("<base64>%s</base64>", base64.StdEncoding.EncodeToString(b)))
		return
	}
	if t != nil {
		switch t.Kind() {
		case reflect.String:
			if _, err = buf.WriteString("<string>"); err != nil {
				return
			}
			if err = xml.EscapeText(buf, []byte(v.(string))); err != nil {
				return
			}
			_, err = buf.WriteString("</string>")
		case reflect.Bool:
			_, err = buf.WriteString(fmt.Sprintf("<boolean>%d</boolean>", asInt(v.(bool))))
		case reflect.Int:
			_, err = buf.WriteString(fmt.Sprintf("<i4>%d</i4>", v.(int)))
		case reflect.Int64:
			_, err = buf.WriteString(fmt.Sprintf("<i8>%d</i8>", v.(int64)))
		case reflect.Slice:
			err = marshalArray(buf, v.([]interface{}))
		case reflect.Map:
			err = marshalMap(buf, v.(map[string]interface{}))
		case reflect.Struct:
			err = marshalStruct(buf, t, v)
		case reflect.Ptr:
			err = marshalValue(buf, reflect.Indirect(reflect.ValueOf(v)).Interface())
		default:
			err = errors.New(fmt.Sprintf("unsupported type: %v", t))
		}
	}
	return
}

func marshalArray(buf *bytes.Buffer, arr []interface{}) (err error) {
	if _, err = buf.WriteString("<array><data>"); err != nil {
		return
	}
	for _, e := range arr {
		if err = marshalValue(buf, e); err != nil {
			break
		}
	}
	_, err = buf.WriteString("</data></array>")
	return
}

func marshalMap(buf *bytes.Buffer, m map[string]interface{}) (err error) {
	if _, err = buf.WriteString("<struct>"); err != nil {
		return
	}
	for k, v := range m {
		if _, err = buf.WriteString("<member>"); err != nil {
			return
		}
		if _, err = buf.WriteString("<name>"); err != nil {
			return
		}
		if err = xml.EscapeText(buf, []byte(k)); err != nil {
			return
		}
		if _, err = buf.WriteString("</name>"); err != nil {
			return
		}
		if err = marshalValue(buf, v); err != nil {
			return
		}
		if _, err = buf.WriteString("</member>"); err != nil {
			return
		}
	}
	_, err = buf.WriteString("</struct>")
	return
}

func marshalStruct(buf *bytes.Buffer, t reflect.Type, st interface{}) (err error) {
	v := reflect.ValueOf(st)
	if _, err = buf.WriteString("<struct>"); err != nil {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if _, err = buf.WriteString("<member>"); err != nil {
			return
		}
		if _, err = buf.WriteString("<name>"); err != nil {
			return
		}
		if err = xml.EscapeText(buf, []byte(t.Field(i).Name)); err != nil {
			return
		}
		if _, err = buf.WriteString("</name>"); err != nil {
			return
		}
		if err = marshalValue(buf, v.FieldByIndex([]int{i}).Interface()); err != nil {
			return
		}
		if _, err = buf.WriteString("</member>"); err != nil {
			return
		}
	}
	_, err = buf.WriteString("</struct>")
	return
}

func asInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

