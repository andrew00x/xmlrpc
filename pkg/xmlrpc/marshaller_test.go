package xmlrpc

import (
	"testing"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/stretchr/testify/assert"
)

func TestMarshalWithoutParams(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("test")
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>test</methodName>
    <params>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithSingleStringParam(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("message", []interface{}{"hello"}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>message</methodName>
    <params>
        <param><value><string>hello</string></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithSingleBytesParam(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("data", []interface{}{[]byte("hello")}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>data</methodName>
    <params>
        <param><value><base64>aGVsbG8=</base64></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithSingleBooleanParam(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("yesOrNo", []interface{}{true}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>yesOrNo</methodName>
    <params>
        <param><value><boolean>1</boolean></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithSingleIntParam(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("number", []interface{}{123}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>number</methodName>
    <params>
        <param><value><i4>123</i4></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithSingleInt64Param(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("number", []interface{}{int64(123)}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>number</methodName>
    <params>
        <param><value><i8>123</i8></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithArray(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("array", []interface{}{[]interface{}{"hello", []byte("hello"), true, 123}}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
   <methodName>array</methodName>
   <params>
		<param><value><array><data>
       	<value><string>hello</string></value>
       	<value><base64>aGVsbG8=</base64></value>
       	<value><boolean>1</boolean></value>
       	<value><i4>123</i4></value>
		</data></array></value></param>
   </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithNestedArray(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("array", []interface{}{[]interface{}{"hello", []interface{}{"to be or not to be", []byte("hello"), true, 123, int64(321)}, "bye"}}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
   <methodName>array</methodName>
   <params>
      	<param><value><array><data>
			<value><string>hello</string></value>
			<value><array><data>
		        <value><string>to be or not to be</string></value>
		        <value><base64>aGVsbG8=</base64></value>
		        <value><boolean>1</boolean></value>
		        <value><i4>123</i4></value>
		        <value><i8>321</i8></value>
			</data></array></value>
			<value><string>bye</string></value>
		</data></array></value></param>
   </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithMap(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("map", []interface{}{map[string]interface{}{"msg": "hello"}}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
   <methodName>map</methodName>
	<params>
      	<param><value><struct>
			<member>
				<name>msg</name>
				<value><string>hello</string></value>
			</member>
		</struct></value></param>
	</params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

type Struct1 struct {
	Foo    string
	Nested Struct2
}

type Struct2 struct {
	Bar string
}

func TestMarshalWithStruct(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("struct", []interface{}{Struct2{"bar"}}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
  <methodName>struct</methodName>
	<params>
     	<param><value><struct>
			<member>
				<name>Bar</name>
				<value><string>bar</string></value>
			</member>
		</struct></value></param>
	</params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithNestedStruct(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("struct", []interface{}{Struct1{"foo", Struct2{"bar"}}}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
  <methodName>struct</methodName>
	<params>
     	<param><value><struct>
			<member>
				<name>Foo</name>
				<value><string>foo</string></value>
			</member>
			<member>
				<name>Nested</name>
				<value><struct>
					<member>
						<name>Bar</name>
						<value><string>bar</string></value>
					</member>
				</struct></value>
			</member>
		</struct></value></param>
	</params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func TestMarshalWithMultipleParams(t *testing.T) {
	m := marshaller{}
	xml, err := m.marshal("multiple", []interface{}{"hello", []byte("hello"), true, 123, int64(321)}...)
	assert.Nil(t, err)
	expected := formatXml(`<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>multiple</methodName>
    <params>
        <param><value><string>hello</string></value></param>
        <param><value><base64>aGVsbG8=</base64></value></param>
        <param><value><boolean>1</boolean></value></param>
        <param><value><i4>123</i4></value></param>
        <param><value><i8>321</i8></value></param>
    </params>
</methodCall>`)
	assert.Equal(t, expected, formatXml(string(xml)))
}

func formatXml(in string) string {
	return xmlfmt.FormatXML(in, "", "  ")
}
