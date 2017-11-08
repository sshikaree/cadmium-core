package xmpp

import (
	"encoding/xml"
	"fmt"
)

// XMPP predefined statuses
type Show struct {
	XMLName xml.Name `xml:"show"`
	Val     string   `xml:",chardata"`
}

func (s *Show) String() string {
	return fmt.Sprintf("<show>%s</show>", s.Val)
}

var (
	Show = struct {
		Chat, Away, XA, DND string
	}{
		"<show>chat</show>",
		"",
		"",
		"",
	}
)
