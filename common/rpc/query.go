package rpc

import (
	"encoding/xml"
)

type Query struct {
	XMLName  xml.Name `xml:"query"`
	InnerXML []byte   `xml:",innerxml"`
}

type Roster struct {
	XMLName  xml.Name  `xml:"jabber:iq:roster query" json:"-"`
	Contacts []Contact `xml:"item"`
}

// type Roster []Contact

type Contact struct {
	XMLName      xml.Name `xml:"item" json:"-"`
	JID          string   `xml:"jid,attr"`
	Subscription string   `xml:"subscription,attr"`
	Name         string   `xml:",attr"`
	Group        []string `xml:"group"`
}

func ParseQuery(query []byte) (stanza interface{}, err error) {
	// name := struct {
	// 	XMLName xml.Name `xml:"query"`
	// }{}
	q := Query{}
	err = xml.Unmarshal(query, &q)
	if err != nil {
		return
	}

	switch q.XMLName.Space {
	// Is it roster?
	case "jabber:iq:roster":
		r := Roster{}
		err = xml.Unmarshal(query, &r)
		return r, err
	}

	// default
	return nil, err
}
