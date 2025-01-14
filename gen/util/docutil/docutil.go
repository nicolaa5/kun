package docutil

import (
	"strings"

	"github.com/RussellLuo/kun/gen/util/annotation"
)

type Transport int

const (
	TransportHTTP  Transport = 0b0001
	TransportGRPC  Transport = 0b0010
	TransportEvent Transport = 0b0100
	TransportAll   Transport = 0b0111
)

type Doc []string

func (d Doc) Transport() (t Transport) {
	for _, comment := range d {
		switch dir := annotation.Directive(comment); dir.Dialect() {
		case annotation.DialectHTTP:
			t = t | TransportHTTP
		case annotation.DialectGRPC:
			t = t | TransportGRPC
		case annotation.DialectEvent:
			t = t | TransportEvent
		}
	}
	return t
}

// JoinComments joins backslash-continued comments.
func (d Doc) JoinComments() (joined Doc) {
	incompleteComment := ""

	for _, comment := range d {
		if incompleteComment == "" {
			if HasContinuationLine(comment) {
				incompleteComment = strings.TrimSuffix(comment, `\`)
			} else {
				joined = append(joined, comment)
			}
			continue
		}

		noPrefix := strings.TrimPrefix(comment, "//")
		c := incompleteComment + strings.TrimSpace(noPrefix)

		if HasContinuationLine(c) {
			incompleteComment = strings.TrimSuffix(c, `\`)
		} else {
			joined = append(joined, c)
			incompleteComment = ""
		}
	}

	return
}

func HasContinuationLine(comment string) bool {
	return strings.HasSuffix(comment, `\`)
}
