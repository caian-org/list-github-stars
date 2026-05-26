package cli

import "fmt"

var (
	ProgramVersion   = "0.0.0-dev"
	ProgramCommitSHA = ""
	ProgramBuildTime = ""
)

func programVersion() string {
	v := ProgramVersion
	if ProgramCommitSHA != "" {
		v = fmt.Sprintf("%s (%s)", v, ProgramCommitSHA)
	}
	if ProgramBuildTime != "" {
		v = fmt.Sprintf("%s, built %s", v, ProgramBuildTime)
	}
	return v
}
