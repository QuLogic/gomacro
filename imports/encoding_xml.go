// this file was generated by gomacro command: import _b "encoding/xml"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	"encoding/xml"
)

// reflection: allow interpreted code to import "encoding/xml"
func init() {
	Packages["encoding/xml"] = Package{
	Binds: map[string]Value{
		"CopyToken":	ValueOf(xml.CopyToken),
		"Escape":	ValueOf(xml.Escape),
		"EscapeText":	ValueOf(xml.EscapeText),
		"HTMLAutoClose":	ValueOf(&xml.HTMLAutoClose).Elem(),
		"HTMLEntity":	ValueOf(&xml.HTMLEntity).Elem(),
		"Header":	ValueOf(xml.Header),
		"Marshal":	ValueOf(xml.Marshal),
		"MarshalIndent":	ValueOf(xml.MarshalIndent),
		"NewDecoder":	ValueOf(xml.NewDecoder),
		"NewEncoder":	ValueOf(xml.NewEncoder),
		"NewTokenDecoder":	ValueOf(xml.NewTokenDecoder),
		"Unmarshal":	ValueOf(xml.Unmarshal),
	}, Types: map[string]Type{
		"Attr":	TypeOf((*xml.Attr)(nil)).Elem(),
		"CharData":	TypeOf((*xml.CharData)(nil)).Elem(),
		"Comment":	TypeOf((*xml.Comment)(nil)).Elem(),
		"Decoder":	TypeOf((*xml.Decoder)(nil)).Elem(),
		"Directive":	TypeOf((*xml.Directive)(nil)).Elem(),
		"Encoder":	TypeOf((*xml.Encoder)(nil)).Elem(),
		"EndElement":	TypeOf((*xml.EndElement)(nil)).Elem(),
		"Marshaler":	TypeOf((*xml.Marshaler)(nil)).Elem(),
		"MarshalerAttr":	TypeOf((*xml.MarshalerAttr)(nil)).Elem(),
		"Name":	TypeOf((*xml.Name)(nil)).Elem(),
		"ProcInst":	TypeOf((*xml.ProcInst)(nil)).Elem(),
		"StartElement":	TypeOf((*xml.StartElement)(nil)).Elem(),
		"SyntaxError":	TypeOf((*xml.SyntaxError)(nil)).Elem(),
		"TagPathError":	TypeOf((*xml.TagPathError)(nil)).Elem(),
		"Token":	TypeOf((*xml.Token)(nil)).Elem(),
		"TokenReader":	TypeOf((*xml.TokenReader)(nil)).Elem(),
		"UnmarshalError":	TypeOf((*xml.UnmarshalError)(nil)).Elem(),
		"Unmarshaler":	TypeOf((*xml.Unmarshaler)(nil)).Elem(),
		"UnmarshalerAttr":	TypeOf((*xml.UnmarshalerAttr)(nil)).Elem(),
		"UnsupportedTypeError":	TypeOf((*xml.UnsupportedTypeError)(nil)).Elem(),
	}, Proxies: map[string]Type{
		"Marshaler":	TypeOf((*Marshaler_encoding_xml)(nil)).Elem(),
		"MarshalerAttr":	TypeOf((*MarshalerAttr_encoding_xml)(nil)).Elem(),
		"TokenReader":	TypeOf((*TokenReader_encoding_xml)(nil)).Elem(),
		"Unmarshaler":	TypeOf((*Unmarshaler_encoding_xml)(nil)).Elem(),
		"UnmarshalerAttr":	TypeOf((*UnmarshalerAttr_encoding_xml)(nil)).Elem(),
	}, Untypeds: map[string]string{
		"Header":	"string:<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n",
	}, 
	}
}

// --------------- proxy for encoding/xml.Marshaler ---------------
type Marshaler_encoding_xml struct {
	Object	interface{}
	MarshalXML_	func(_proxy_obj_ interface{}, e *xml.Encoder, start xml.StartElement) error
}
func (Proxy *Marshaler_encoding_xml) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return Proxy.MarshalXML_(Proxy.Object, e, start)
}

// --------------- proxy for encoding/xml.MarshalerAttr ---------------
type MarshalerAttr_encoding_xml struct {
	Object	interface{}
	MarshalXMLAttr_	func(_proxy_obj_ interface{}, name xml.Name) (xml.Attr, error)
}
func (Proxy *MarshalerAttr_encoding_xml) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return Proxy.MarshalXMLAttr_(Proxy.Object, name)
}

// --------------- proxy for encoding/xml.TokenReader ---------------
type TokenReader_encoding_xml struct {
	Object	interface{}
	Token_	func(interface{}) (xml.Token, error)
}
func (Proxy *TokenReader_encoding_xml) Token() (xml.Token, error) {
	return Proxy.Token_(Proxy.Object)
}

// --------------- proxy for encoding/xml.Unmarshaler ---------------
type Unmarshaler_encoding_xml struct {
	Object	interface{}
	UnmarshalXML_	func(_proxy_obj_ interface{}, d *xml.Decoder, start xml.StartElement) error
}
func (Proxy *Unmarshaler_encoding_xml) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return Proxy.UnmarshalXML_(Proxy.Object, d, start)
}

// --------------- proxy for encoding/xml.UnmarshalerAttr ---------------
type UnmarshalerAttr_encoding_xml struct {
	Object	interface{}
	UnmarshalXMLAttr_	func(_proxy_obj_ interface{}, attr xml.Attr) error
}
func (Proxy *UnmarshalerAttr_encoding_xml) UnmarshalXMLAttr(attr xml.Attr) error {
	return Proxy.UnmarshalXMLAttr_(Proxy.Object, attr)
}
