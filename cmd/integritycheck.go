package cmd

import (
	"bitbucket.org/level27/lvl/types"
	"bitbucket.org/level27/lvl/utils"
)

func outputFormatIntegrityCheckTable(checks []types.IntegrityCheck) {
	outputFormatTableFuncs(
		checks,
		[]string{"ID", "STATUS", "DATE"},
		[]interface{}{"Id", "Status", func(s types.IntegrityCheck) string {
			return utils.FormatUnixTime(s.DtRequested)
		}})
}