package utils_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/level27/lvl/utils"
)

func TestZoneParse(t *testing.T) {
	text := `
; Zone file for bedrijfskat.be
;
$TTL 300
@	IN	SOA	dns1.cp4staging.be.	hostmaster.level27.be. (
    2100000044      ; serial
    14400     ; refresh
    3600       ; retry
    1209600      ; expire
    3600)    ; minimum
;

; NS
	NS dns1.cp4staging.be.
	NS dns2.cp4staging.be.

; MX
bedrijfskat.be.	MX	10 smtp2.cp4staging.be.
bedrijfskat.be.	MX	20 smtp1.cp4staging.be.


; TXT
amai IN TXT "asdfasdfasdf"
bedrijfskat.be. IN TXT "v=spf1 include:_mail.cp4staging.be ?all"
	TXT	foobar\036
	TXT	foobar\"
	TXT	"foobar\""
	TXT	"foobar\036"

; SRV
_autodiscover._tcp	SRV	0 5 443 autodiscover.cp4staging.be.

; TLSA

; DS

; others
urlfwtest-1	A	172.29.17.26
urlfwtest-2	A	195.225.166.11
urlfwtest-2	AAAA	2a02:5b40:4:210::dddd:10
solr-vm-size-2	A	172.29.17.29
solr-8-vm-size	A	172.29.17.29
blabla	A	195.225.166.14
letest	A	172.29.17.31
cms-install-1-wp	A	195.225.166.15
cms-install-2-drupal	A	195.225.166.15
cms-install-3-magento	A	195.225.166.15
scanner	A	172.20.17.4
cms-install-4-shopware	A	195.225.166.15
autoconfig	CNAME	autodiscover.cp4staging.be.
`
	parser := utils.NewZoneParser(bytes.NewReader([]byte(text)))
	assertTtl(t, &parser, 300)

	// SOA
	assertRr(t, &parser, dom("@"), dnsClass(utils.DnsClassIN), nil, utils.RecordTypeSOA, []string{"dns1.cp4staging.be.", "hostmaster.level27.be.", "2100000044", "14400", "3600", "1209600", "3600"})

	// NS
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeNS, []string{"dns1.cp4staging.be."})
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeNS, []string{"dns2.cp4staging.be."})

	// MX
	assertRr(t, &parser, dom("bedrijfskat.be."), nil, nil, utils.RecordTypeMX, []string{"10", "smtp2.cp4staging.be."})
	assertRr(t, &parser, dom("bedrijfskat.be."), nil, nil, utils.RecordTypeMX, []string{"20", "smtp1.cp4staging.be."})

	// TXT
	assertRr(t, &parser, dom("amai"), dnsClass(utils.DnsClassIN), nil, utils.RecordTypeTXT, []string{"asdfasdfasdf"})
	assertRr(t, &parser, dom("bedrijfskat.be."), dnsClass(utils.DnsClassIN), nil, utils.RecordTypeTXT, []string{"v=spf1 include:_mail.cp4staging.be ?all"})
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeTXT, []string{"foobar$"})
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeTXT, []string{"foobar\""})
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeTXT, []string{"foobar\""})
	assertRr(t, &parser, nil, nil, nil, utils.RecordTypeTXT, []string{"foobar$"})

	// SRV
	assertRr(t, &parser, dom("_autodiscover._tcp"), nil, nil, utils.RecordTypeSRV, []string{"0", "5", "443", "autodiscover.cp4staging.be."})

	// A/AAAA/CNAME
	assertRr(t, &parser, dom("urlfwtest-1"), nil, nil, utils.RecordTypeA, []string{"172.29.17.26"})
	assertRr(t, &parser, dom("urlfwtest-2"), nil, nil, utils.RecordTypeA, []string{"195.225.166.11"})
	assertRr(t, &parser, dom("urlfwtest-2"), nil, nil, utils.RecordTypeAAAA, []string{"2a02:5b40:4:210::dddd:10"})
	assertRr(t, &parser, dom("solr-vm-size-2"), nil, nil, utils.RecordTypeA, []string{"172.29.17.29"})
	assertRr(t, &parser, dom("solr-8-vm-size"), nil, nil, utils.RecordTypeA, []string{"172.29.17.29"})
	assertRr(t, &parser, dom("blabla"), nil, nil, utils.RecordTypeA, []string{"195.225.166.14"})
	assertRr(t, &parser, dom("letest"), nil, nil, utils.RecordTypeA, []string{"172.29.17.31"})
	assertRr(t, &parser, dom("cms-install-1-wp"), nil, nil, utils.RecordTypeA, []string{"195.225.166.15"})
	assertRr(t, &parser, dom("cms-install-2-drupal"), nil, nil, utils.RecordTypeA, []string{"195.225.166.15"})
	assertRr(t, &parser, dom("cms-install-3-magento"), nil, nil, utils.RecordTypeA, []string{"195.225.166.15"})
	assertRr(t, &parser, dom("scanner"), nil, nil, utils.RecordTypeA, []string{"172.20.17.4"})
	assertRr(t, &parser, dom("cms-install-4-shopware"), nil, nil, utils.RecordTypeA, []string{"195.225.166.15"})
	assertRr(t, &parser, dom("autoconfig"), nil, nil, utils.RecordTypeCNAME, []string{"autodiscover.cp4staging.be."})

	assertEof(t, &parser)
}

func assertTtl(t *testing.T, parser *utils.ZoneParser, ttl utils.RecordTtl) {
	entry, err := parser.NextEntry()
	if err != nil {
		t.Fatal(err)
		return
	}

	entryTtl, ok := entry.(utils.ZoneEntryTtl)
	if !ok {
		t.Fatal("Expected TTL entry, got:", entry)
	}

	if entryTtl.Ttl != ttl {
		t.Fatal("Unexpected TTL. Expected ", ttl, "got", entryTtl.Ttl)
	}
}

func assertRr(
	t *testing.T,
	parser *utils.ZoneParser,
	domain *string,
	class *utils.DnsClass,
	ttl *utils.RecordTtl,
	recordType utils.RecordType,
	params []string) {
	entry, err := parser.NextEntry()
	if err != nil {
		t.Fatal(err)
	}

	entryRr, ok := entry.(utils.ZoneEntryRr)
	if !ok {
		t.Fatal("Expected RR entry, got:", entry)
	}

	if !ptrEq(entryRr.DomainName, domain) {
		t.Fatal("Unexpected domain. Expected", domain, "got", printPtrStr(entryRr.DomainName))
	}

	if !ptrEq(entryRr.Class, class) {
		t.Fatal("Unexpected class. Expected", printPtr(class), "got", printPtr(entryRr.Class))
	}

	if !ptrEq(entryRr.Ttl, ttl) {
		t.Fatal("Unexpected TTL. Expected", printPtr(ttl), "got", printPtr(entryRr.Ttl))
	}

	if entryRr.Type != recordType {
		t.Fatal("Unexpected record type. Expected", recordType, "got", entryRr.Type)
	}

	if !reflect.DeepEqual(entryRr.Data, params) {
		t.Fatal("Mismatching record data.")
	}
}

func assertEof(t *testing.T, parser *utils.ZoneParser) {
	_, err := parser.NextEntry()
	if err != io.EOF {
		t.Fatal("Expected EOF, got", err)
	}
}

func dnsClass(class utils.DnsClass) *utils.DnsClass {
	return &class
}

func ttl(ttl utils.RecordTtl) *utils.RecordTtl {
	return &ttl
}

func dom(domain string) *string {
	return &domain
}

func printPtr[T fmt.Stringer](value *T) string {
	if value == nil {
		return "<nil>"
	}

	return (*value).String()
}

func printPtrStr(value *string) string {
	if value == nil {
		return "<nil>"
	}

	return (*value)
}

func ptrEq[T comparable](a *T, b *T) bool {
	// != is an XOR here.
	if (a == nil) != (b == nil) {
		return false
	}

	if a == nil {
		return true
	}

	return *a == *b
}
