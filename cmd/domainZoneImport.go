package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	domainCmd.AddCommand(domainZoneImportCmd)
	domainZoneImportCmd.Flags().BoolVarP(&domainZoneImportYes, "yes", "y", false, "Confirm import of file without prompt")
}

var domainZoneImportYes bool
var domainZoneImportCmd = &cobra.Command{
	Use:   "zoneimport <domain> <zone file>",
	Short: "Import DNS records for a domain from a zone file",
	Long: `Import DNS records for a domain from a zone file.
Existing records (same name/type) will be replaced by the new records.
Pass '-' as file name to read from stdin.`,

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := openArgFile(args[1])
		if err != nil {
			return fmt.Errorf("failed to open input: %s", err.Error())
		}

		defer file.Close()

		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		domain, err := Level27Client.Domain(domainID)
		if err != nil {
			return err
		}

		existingRecords, err := Level27Client.DomainRecords(domainID, "", l27.CommonGetParams{PageableParams: l27.PageableParams{Limit: 10000}})
		if err != nil {
			return err
		}

		origin := fmt.Sprintf("%s.", domain.Fullname)

		// Build index to find records to replace.
		existingRecordsIndex := zoneImportMakeExistingRecordsIndex(existingRecords)
		toReplace, toCreate := zoneDomainImportParse(origin, file, existingRecordsIndex)

		fmt.Printf(
			"%d existing records to delete (for replacement)\n%d records to create\n",
			len(toReplace),
			len(toCreate))

		if !domainZoneImportYes {
			if !confirmPrompt("Confirm importing records?") {
				return nil
			}
		}

		for id := range toReplace {
			err := Level27Client.DomainRecordDelete(domainID, id)
			if err != nil {
				return err
			}
		}

		for _, request := range toCreate {
			_, err := Level27Client.DomainRecordCreate(domainID, request)
			if err != nil {
				return err
			}
		}

		fmt.Printf("All records successfully imported")

		return nil
	},
}

func zoneDomainImportParse(
	origin string,
	file io.Reader,
	existingRecordsIndex map[zoneImportingExistingRecord][]l27.IntID,
) (map[l27.IntID]bool, []l27.DomainRecordRequest) {
	toReplace := map[l27.IntID]bool{}
	toCreate := []l27.DomainRecordRequest{}

	var currentClass utils.DnsClass = 0
	currentOrigin := origin
	warnedTtlDirective := false
	warnedTtlRecord := false
	warnedClass := false

	lastDomain := "@"

	parser := utils.NewZoneParser(file)
	for {
		entry, err := parser.NextEntry()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Printf("Error parsing record: %s\n", err.Error())
			continue
		}

		// fmt.Printf("%v\n", entry)
		if _, ok := entry.(utils.ZoneEntryTtl); ok {
			if !warnedTtlDirective {
				fmt.Printf("Note: TTL directives are not imported, set TTL manually after import.\n")
				warnedTtlDirective = true
			}
		} else if entryOrigin, ok := entry.(utils.ZoneEntryOrigin); ok {
			currentOrigin = strings.ToLower(entryOrigin.DomainName)
		} else if rr, ok := entry.(utils.ZoneEntryRr); ok {
			if rr.Ttl != nil && !warnedTtlRecord {
				fmt.Printf("Note: Level27 does not support per-record TTL values, TTL values will be ignored.\n")
				warnedTtlRecord = true
			}

			if rr.Class != nil {
				currentClass = *rr.Class
			}

			if rr.DomainName != nil {
				lastDomain = strings.ToLower(*rr.DomainName)
			}

			if currentClass == 0 {
				fmt.Printf("Warning: no DNS class given for record: %v\n", rr)
				continue
			}

			if currentClass != utils.DnsClassIN {
				if !warnedClass {
					fmt.Printf("Note: Level27 does not support non-IN records, ignoring.\n")
					warnedClass = true
				}

				continue
			}

			finalName := zoneDomainNormalizeOrigin(lastDomain, currentOrigin, origin)

			// Check if there is already an existing record in the API of this type/name.
			// Add them to the list of records to delete on commit.
			existingRecord := zoneImportingExistingRecord{
				Type: rr.Type.String(),
				Name: finalName,
			}
			for _, id := range existingRecordsIndex[existingRecord] {
				toReplace[id] = true
			}

			request := l27.DomainRecordRequest{
				Type: rr.Type.String(),
				Name: finalName,
			}

			if request.Name == "@" {
				request.Name = ""
			}

			switch rr.Type {
			case utils.RecordTypeA:
				request.Content = rr.Data[0]
			case utils.RecordTypeAAAA:
				request.Content = rr.Data[0]
			case utils.RecordTypeMX:
				priority, err := strconv.ParseInt(rr.Data[0], 10, 32)
				if err != nil {
					fmt.Printf("Invalid priority in MX record: '%s'\n", rr.Data[0])
					continue
				}

				request.Priority = int32(priority)
				request.Content = rr.Data[1]
			case utils.RecordTypeTXT:
				request.Content = strings.Join(rr.Data, "")
			case utils.RecordTypeCNAME:
				request.Content = rr.Data[0]
			case utils.RecordTypeNS:
				if request.Name == "" {
					fmt.Printf("Note: NS record at domain origin ignored.\n")
					continue
				}
				request.Content = rr.Data[0]
			case utils.RecordTypeSRV:
				request.Content = strings.Join(rr.Data, " ")
			case utils.RecordTypeTLSA:
				request.Content = strings.Join(rr.Data, " ")
			case utils.RecordTypeCAA:
				request.Content = strings.Join(rr.Data, " ")
			case utils.RecordTypeDS:
				request.Content = strings.Join(rr.Data, " ")
			default:
				fmt.Printf("Note: Level27 does not support importing %v records, ignoring.\n", rr.Type)
				continue
			}

			toCreate = append(toCreate, request)
		}
	}

	return toReplace, toCreate
}

func zoneDomainNormalizeOrigin(domain string, curOrigin string, destOrigin string) string {
	concat := zoneDomainConcat(domain, curOrigin)
	return zoneDomainRelative(concat, destOrigin)
}

// Make a domain absolute by appending the origin (if it's not yet absolute).
func zoneDomainConcat(domain string, origin string) string {
	if strings.HasSuffix(domain, ".") {
		return domain
	}

	if domain == "@" {
		return origin
	}

	return fmt.Sprintf("%s.%s", domain, origin)
}

// Make a domain relative again by splitting off the
// "xyz.foo.bar.baz.", "bar.baz." -> "xyz.foo"
// "bar.baz.", "bar.baz."         -> "@"
// "abc.xyz.", "bar.baz."         -> "abc.xyz."
func zoneDomainRelative(domain string, origin string) string {
	if domain == origin {
		// Same domain as origin
		return "@"
	}

	if strings.HasSuffix(domain, fmt.Sprintf(".%s", origin)) {
		// Subdomain of origin
		return domain[:len(domain)-len(origin)-1]
	}

	// Not related at all
	return domain
}

func zoneImportMakeExistingRecordsIndex(records []l27.DomainRecord) map[zoneImportingExistingRecord][]l27.IntID {
	index := map[zoneImportingExistingRecord][]l27.IntID{}

	for _, record := range records {
		existing := zoneImportingExistingRecord{Type: record.Type, Name: record.Name}
		ids := index[existing]
		ids = append(ids, record.ID)
		index[existing] = ids
	}

	return index
}

type zoneImportingExistingRecord struct {
	Type string
	Name string
}
