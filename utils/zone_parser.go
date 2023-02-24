package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//
// Parser for zone files.
// References:
// * RFC 1035 (section 5): https://datatracker.ietf.org/doc/html/rfc1035
// * RFC 2308 (section 4): https://datatracker.ietf.org/doc/html/rfc2308
// * RFC 2136:             https://datatracker.ietf.org/doc/html/rfc2136
// * https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml
// * https://bind9.readthedocs.io/en/v9_18_4/chapter3.html
//
// This is mostly a parser for the structure of the zone file,
// it does not intend to derive meaning.
// e.g. outputted data does not automatically merge TTL values and such.
// e.g. it only has a list of record types,
//      since these are necessary to disambiguate the parsing of RRs.
//      Parsing of record data is left to using code.
//

// $ORIGIN example.com.
// foo\.bar IN A 1.1.1.1
// Actions.domains@example.com
// @   IN  SOA     VENERA      Action\.domains (
//
//	20     ; SERIAL
//	7200   ; REFRESH
//	600    ; RETRY
//	3600000; EXPIRE
//	60)    ; MINIMUM
type RecordTtl uint32

func (ttl RecordTtl) String() string {
	return fmt.Sprint(uint32(ttl))
}

type DnsClass uint16

const (
	DnsClassIN DnsClass = 1
	DnsClassCH DnsClass = 3
	DnsClassHS DnsClass = 4
)

var classMap = map[string]DnsClass{
	"IN": DnsClassIN,
	"CH": DnsClassCH,
	"HS": DnsClassCH,
}

var classMapReverse = reverseMap(classMap)

func (t DnsClass) String() string {
	return classMapReverse[t]
}

type RecordType uint16

// Only record types supported by Level27 DNS are listed here.

const (
	RecordTypeA     RecordType = 1
	RecordTypeNS    RecordType = 2
	RecordTypeCNAME RecordType = 5
	// SOA is not supported by Level27 but it shows up a lot so we should be able to parse it.
	RecordTypeSOA  RecordType = 6
	RecordTypeMX   RecordType = 15
	RecordTypeTXT  RecordType = 16
	RecordTypeAAAA RecordType = 28
	RecordTypeSRV  RecordType = 33
	RecordTypeDS   RecordType = 43
	RecordTypeTLSA RecordType = 52
	RecordTypeCAA  RecordType = 257
)

var typeMap = map[string]RecordType{
	"A":     RecordTypeA,
	"NS":    RecordTypeNS,
	"CNAME": RecordTypeCNAME,
	"SOA":   RecordTypeSOA,
	"MX":    RecordTypeMX,
	"TXT":   RecordTypeTXT,
	"AAAA":  RecordTypeAAAA,
	"SRV":   RecordTypeSRV,
	"DS":    RecordTypeDS,
	"TLSA":  RecordTypeTLSA,
	"CAA":   RecordTypeCAA,
}

var typeMapReverse = reverseMap(typeMap)

func (t RecordType) String() string {
	return typeMapReverse[t]
}

type ZoneParser struct {
	reader    *bufio.Reader
	lineIndex int32
}

// From https://www.reddit.com/r/golang/comments/q4a70y/how_do_experienced_go_developers_model_sum_types/
// Declaring a sum type in Go.

// Base type for all entry types parseable in a zone file.
type ZoneEntry interface {
	IsZoneEntry()
}

// $ORIGIN zone file entry.
type ZoneEntryOrigin struct {
	DomainName string
}

func (ZoneEntryOrigin) IsZoneEntry() {}

func (e ZoneEntryOrigin) String() string {
	return fmt.Sprintf("$TTL %s;", e.DomainName)
}

// $INCLUDE zone file entry.
type ZoneEntryInclude struct {
	FileName   string
	DomainName string
}

func (ZoneEntryInclude) IsZoneEntry() {}

// $TTL zone file entry.
type ZoneEntryTtl struct {
	Ttl RecordTtl
}

func (ZoneEntryTtl) IsZoneEntry() {}

func (e ZoneEntryTtl) String() string {
	return fmt.Sprintf("$TTL %d;", e.Ttl)
}

// Resource Record zone file entry.
type ZoneEntryRr struct {
	// Can be nil to indicate not given for record
	DomainName *string
	// Can be nil to indicate not given for record
	Class *DnsClass
	// Can be nil to indicate not given for record
	Ttl *RecordTtl
	// Type of the record.
	Type RecordType
	// Record data is given as an unstructed set of strings.
	Data []string
}

func (ZoneEntryRr) IsZoneEntry() {}

func (e ZoneEntryRr) String() string {
	domain := ""
	if e.DomainName != nil {
		domain = *e.DomainName
	}

	class := ""
	if e.Class != nil {
		class = e.Class.String()
	}

	ttl := ""
	if e.Ttl != nil {
		ttl = e.Ttl.String()
	}

	value := strings.Join(e.Data, ", ")

	return fmt.Sprintf("%s\t%s\t%s\t%v\t%s", domain, class, ttl, e.Type, value)
}

func NewZoneParser(reader io.Reader) ZoneParser {
	return ZoneParser{
		reader: bufio.NewReader(reader),
	}
}

// Parse an entry from the input. This moves the parser forward in the input.
func (z *ZoneParser) NextEntry() (ZoneEntry, error) {
	// The parser should always be at the start of a new entry's line when this function is called.
	// Empty entries (blank lines) are implicitly skipped.
	// If a parser errors occurs, we automatically move to the next line in the hope that
	// parsing can somewhat recover (so we can report more than one error if possible)
	// This means that if we suddenly read an unexpected newline character,
	// it should be unread from the bufreader so we can skip to it.

	// State for the current parsing operation. No parse state is carried through between entries.

	var startLine int32
	var state *zoneEntryParseState

	var entry ZoneEntry
	for {
		startLine = z.lineIndex
		state = &zoneEntryParseState{}

		// nextEntryCore always leaves the current read position right before a newline.
		// This allows us to consistently handle error and success scenarios
		// for going to the next directive.
		// It also means nextItem() will consistently "stick" at an EOL until moved up by this code.
		var err error
		entry, err = z.nextEntryCore(state)
		if err != nil {
			// We have to skip until the next line to hopefully allow the parser to recover.
			// These may produce generic IO/EOL errors. If they do it's not a big deal.
			z.skipUntilEol()
			z.skipNewlines()
			if err == io.EOF {
				return nil, io.EOF
			}

			return nil, fmt.Errorf("error on directive starting at line %d: %s", startLine+1, err)
		}

		if entry == nil {
			// Empty line. Skip newlines and let the loop try to read a new
			err = z.skipNewlines()
			if err != nil {
				return nil, err
			}

			continue
		}

		break
	}

	// Assuming a well-formed file, we should be at the end of a line (or EOF).
	// For RRs this is guaranteed since they always try to parse as much items as possible.
	// Special directives like $TTL expect a fixed count of items however,
	// so anything after the TTL (invalid data or just whitespace/comment) is not read over yet.
	// This nextItem() *should* give an end-of-directive mark.
	// It will also make sure to close parentheses if we're in those, assert that below.
	item, err := z.nextItem(state)
	if !z.parseIsDirectiveEnd(err) {
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("unexpected item found after directive starting at line %d: '%s'", startLine+1, item)
	}

	if z.isInParentheses(state) {
		return nil, fmt.Errorf("unclosed parentheses pair starting at line %d", *state.ParenthesesStartLine+1)
	}

	// Sanity assert we're right before a newline or EOL. If not, it's a bug in the parser.
	read, err := z.skipUntilEol()
	if err != nil {
		return nil, err
	}

	if read != 0 {
		panic("parser was not before EOL/EOF after parsing last item in directive")
	}

	// This skipNewlines() call (and the one above) may eat multiple lines if they're empty.
	// This doesn't really matter since we keep count of the index for error reporting either way.
	err = z.skipNewlines()
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// Core parsing code. Does not guarantee leaving read position in consistent state.
func (z *ZoneParser) nextEntryCore(state *zoneEntryParseState) (ZoneEntry, error) {
	// Starting directives or domain names MUST be at the hard start of the line.
	keyItem, err := z.parseLooseItem()
	if err != nil {
		return nil, err
	}

	// Note: keyItem will be an empty string if the start of the entry is blank.
	// Which happens for an RR that inherits the domain from the previous RR.

	if strings.HasPrefix(keyItem, "$") {
		// Indeed a special directive.
		switch keyItem {
		case "$TTL":
			return z.parseTtldirective(state)
		}
	}

	// Regular resource record.
	return z.parseRrDirective(keyItem, state)
}

func (z *ZoneParser) parseTtldirective(state *zoneEntryParseState) (ZoneEntry, error) {
	valueItem, err := z.nextItem(state)
	if err != nil {
		return nil, err
	}

	ttl, err := strconv.ParseUint(valueItem, 10, 32)
	if err != nil {
		return nil, err
	}

	return ZoneEntryTtl{
		Ttl: RecordTtl(ttl),
	}, nil
}

func (z *ZoneParser) parseRrDirective(keyItem string, state *zoneEntryParseState) (ZoneEntry, error) {
	// Note: keyItem may be empty string if RR has no specified domain name.

	// Read first two items to allow us to clearly tell the order of [<class>], [<TTL>] and <type>.
	// In the shortest form this may be part of the <RDATA> we're reading,
	// that's fine since there should be always at least one item in RDATA.
	firstItem, err := z.nextItem(state)
	if err != nil {
		if err == errZoneParseEol && keyItem == "" {
			// No item at start of line and no items afterwards
			// means there's actually just nothing on this line!
			// Just return nil up the chain and let NextEntry() loop for the next line.
			return nil, nil
		}
		return nil, fmt.Errorf("failed reading first directive item: %v", err)
	}

	secondItem, err := z.nextItem(state)
	if err != nil {
		return nil, fmt.Errorf("failed reading second directive item: %v", err)
	}

	// Pointer types as nillables.
	var ttl *RecordTtl
	var class *DnsClass
	var recType *RecordType
	var rdata []string

	err = z.checkFirstRecordItem(firstItem, &ttl, &class, &recType)
	if err != nil {
		return nil, err
	}

	if recType == nil {
		// First item wasn't record type, check second.
		err = z.checkFirstRecordItem(secondItem, &ttl, &class, &recType)
		if err != nil {
			return nil, err
		}
	} else {
		// First item was record type so second item must've been part of <RDATA>
		rdata = append(rdata, secondItem)
	}

	if recType == nil {
		// If we still don't have the record type, it MUST be the next item.
		item, err := z.nextItem(state)
		if err != nil {
			return nil, err
		}

		recTypeValue, ok := typeMap[item]
		if !ok {
			return nil, fmt.Errorf("unknown record type: '%s'", item)
		}

		recType = &recTypeValue
	}

	for {
		item, err := z.nextItem(state)
		if err != nil {
			if z.parseIsDirectiveEnd(err) {
				break
			}

			return nil, err
		}

		rdata = append(rdata, item)
	}

	entry := ZoneEntryRr{
		Class: class,
		Ttl:   ttl,
		Type:  *recType,
		Data:  rdata,
	}

	if keyItem == "" {
		entry.DomainName = nil
	} else {
		entry.DomainName = &keyItem
	}

	return entry, nil
}

func (z *ZoneParser) checkFirstRecordItem(
	item string,
	ttl **RecordTtl,
	class **DnsClass,
	recType **RecordType) error {

	// Check if it's a TTL value.
	if isAsciiDigit(item[0]) {
		if *ttl != nil {
			return fmt.Errorf("found second number when TTL value already given: '%s'", item)
		}

		ttlVal, err := strconv.ParseUint(item, 10, 16)
		if err != nil {
			return fmt.Errorf("error parsing TTL value '%s': %s", item, err.Error())
		}

		ttlValue16 := RecordTtl(ttlVal)
		*ttl = &ttlValue16
		return nil
	}

	// Check if it's a DNS class.
	if classValue, ok := classMap[item]; ok {
		if *class != nil {
			return fmt.Errorf("found second DNS class when DNS class already given: '%s'", item)
		}

		*class = &classValue
		return nil
	}

	// Check if it's a record type.
	if recTypeValue, ok := typeMap[item]; ok {
		*recType = &recTypeValue
		return nil
	}

	return fmt.Errorf("invalid domain class, TTL value or record type: '%s'", item)
}

// Does actual lexing work.
func (z *ZoneParser) nextItem(state *zoneEntryParseState) (string, error) {
	// Loop to allow repeat skipping of whitespace and/or comments.
	for {
		err := z.skipWhitespace(state)
		if err != nil {
			return "", err
		}

		chr, _, err := z.reader.ReadRune()
		if err != nil {
			return "", err
		}

		if chr == ';' {
			// Comment, skip until newline.
			_, err = z.skipUntilEol()
			if err != nil {
				return "", err
			}

			// Let the loop catch the newline.
			continue
		}

		if chr == '\n' || chr == '\r' {
			if z.isInParentheses(state) {
				// We're in parentheses, just keep chomping lines.
				z.countNewline(chr)
				continue
			}

			z.reader.UnreadRune()
			return "", errZoneParseEol
		}

		if chr == '(' {
			if z.isInParentheses(state) {
				return "", errors.New("found opening parentheses while already inside parentheses")
			}

			copy := z.lineIndex
			state.ParenthesesStartLine = &copy
			continue
		}

		if chr == ')' {
			if !z.isInParentheses(state) {
				return "", errors.New("found closing parentheses but not inside parentheses")
			}

			state.ParenthesesStartLine = nil
			continue
		}

		if chr == '"' {
			return z.parseQuotedItem()
		}

		// Non-special character, start reading it into an item.
		z.reader.UnreadRune()
		item, err := z.parseLooseItem()
		if err != nil {
			return "", err
		}

		if item == "" {
			return "", errZoneParseEol
		}

		return item, nil
	}
}

// Returns the amount of runes skipped before reaching EOL/error.
func (z *ZoneParser) skipUntilEol() (int, error) {
	read := 0
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return read, err
		}

		if chr == '\n' || chr == '\r' {
			z.reader.UnreadRune()
			break
		}

		read += 1
	}

	return read, nil
}

func (z *ZoneParser) skipNewlines() error {
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if chr != '\n' && chr != '\r' {
			z.reader.UnreadRune()
			break
		}

		z.countNewline(chr)
	}

	return nil
}

// Generally we treat CR and LF the same way, skip em both when in parentheses. This just works.
// For line index counting we have to avoid counting CR though.
func (z *ZoneParser) countNewline(chr rune) {
	if chr != '\r' {
		z.lineIndex += 1
	}
}

func (z *ZoneParser) parseLooseItem() (string, error) {
	var item bytes.Buffer
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			return "", err
		}

		if z.isWhitespace(chr) {
			break
		}

		if chr == '\n' || chr == '\r' || chr == ';' || chr == '(' || chr == ')' || chr == '"' {
			// Immediately unread the character.
			// The next call to nextItem will handle it appropriately.
			z.reader.UnreadRune()
			break
		}

		// TODO: Escapes.

		item.WriteRune(chr)
	}

	return item.String(), nil
}

func (z *ZoneParser) parseQuotedItem() (string, error) {
	var item bytes.Buffer
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return "", errors.New("early EOF while parsing quoted item")
			}
			return "", err
		}

		// TODO: Escapes.

		if chr == '"' {
			break
		}

		item.WriteRune(chr)
	}

	return item.String(), nil
}

func (z *ZoneParser) parseIsDirectiveEnd(err error) bool {
	return err == io.EOF || err == errZoneParseEol
}

var errZoneParseEol = errors.New("end of directive")

func (z *ZoneParser) skipWhitespace(state *zoneEntryParseState) error {
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			return err
		}

		if z.isWhitespace(chr) {
			continue
		}

		z.reader.UnreadRune()
		break
	}

	return nil
}

func (*ZoneParser) isWhitespace(chr rune) bool {
	return chr == '\t' || chr == ' '
}

func (*ZoneParser) isInParentheses(state *zoneEntryParseState) bool {
	return state.ParenthesesStartLine != nil
}

// Reads a rune, but ignore newlines when in parentheses.
/*
func (z *ZoneParser) readRune(state *zoneEntryParseState) (rune, error) {
	for {
		chr, _, err := z.reader.ReadRune()
		if err != nil {
			return 0, err
		}

		if state.InParentheses && (chr == '\r' || chr == '\n') {
			continue
		}

		return chr, nil
	}
}
*/

type zoneEntryParseState struct {
	ParenthesesStartLine *int32
}

func isAsciiDigit(chr byte) bool {
	return chr >= 0x30 && chr <= 0x39
}

func reverseMap[K comparable, V comparable](m map[K]V) map[V]K {
	new := map[V]K{}

	for k, v := range m {
		new[v] = k
	}

	return new
}
