package xmlrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMessageWithSingleStringParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><string>hello</string></value></param>
	</params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"hello"}, res)
}

func TestUnmarshalMessageWithSingleBase64Param(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><base64>aGVsbG8=</base64></value></param>
	</params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{[]byte("hello")}, res)
}

func TestUnmarshalMessageWithSingleBooleanParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><boolean>1</boolean></value></param>
	</params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{true}, res)
}

func TestUnmarshalMessageWithSingleIntParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><i4>123</i4></value></param>
	</params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{123}, res)
}

func TestUnmarshalMessageWithSingleInt64Param(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><i8>123</i8></value></param>
	</params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{int64(123)}, res)
}

func TestUnmarshalMessageWithArrayParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
	<params>
		<param><value><array><data>
			<value><string>hello</string></value>
			<value><base64>aGVsbG8=</base64></value>
			<value><boolean>1</boolean></value>
			<value><i4>123</i4></value>
			<value><i8>321</i8></value>
		</data></array></value></param>
	</params>
</methodResponse>
	`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{[]interface{}{"hello", []byte("hello"), true, 123, int64(321)}}, res)
}

func TestUnmarshalMessageWithNestedArraysParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
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
</methodResponse>
	`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{[]interface{}{"hello", []interface{}{"to be or not to be", []byte("hello"), true, 123, int64(321)}, "bye"}}, res)
}

func TestUnmarshalMessageWithStructParam(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
	<params>
      	<param><value><struct>
			<member>
				<name>msg</name>
				<value><string>hello</string></value>
			</member>
		</struct></value></param>
	</params>
</methodResponse>
	`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{map[string]interface{}{"msg": "hello"}}, res)
}

func TestUnmarshalWithMultipleParams(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param><value><string>hello</string></value></param>
	    <param><value><array><data>
			<value><string>hello</string></value>
		    <value><base64>aGVsbG8=</base64></value>
		    <value><boolean>1</boolean></value>
		    <value><i4>123</i4></value>
		    <value><i8>321</i8></value>
		</data></array></value></param>
        <param><value><base64>aGVsbG8=</base64></value></param>
        <param><value><boolean>1</boolean></value></param>
        <param><value><i4>123</i4></value></param>
      	<param><value><struct>
			<member>
				<name>msg</name>
				<value><string>hello</string></value>
			</member>
		</struct></value></param>
    </params>
</methodResponse>
`)

	u := unmarshaller{}
	res, err := u.unmarshal(xml)

	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"hello", []interface{}{"hello", []byte("hello"), true, 123, int64(321)}, []byte("hello"), true, 123, map[string]interface{}{"msg": "hello"}}, res)
}

func TestUnmarshalFaultMessage(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
   <fault>
      <value>
         <struct>
            <member>
               <name>faultCode</name>
               <value><int>3</int></value>
               </member>
            <member>
               <name>faultString</name>
               <value><string>something went wrong</string></value>
               </member>
            </struct>
         </value>
      </fault>
</methodResponse>
	`)

	u := unmarshaller{}
	_, err := u.unmarshal(xml)

	assert.Equal(t, "error response, code: 3, text: something went wrong", err.Error())
}
