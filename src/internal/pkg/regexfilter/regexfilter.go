package regexfilter

import (
	"bufio"
	"fmt"
	"internal/pkg/configuration"
	"internal/pkg/loader"
	"regexp"
	"strings"
)

//Finding contains the matched line, the location of the match, and the confidence of the match
type Finding struct {
	Match      string
	Location   string
	Rule       string
	Confidence float64
	Severity   int
}

//IsEmpty returns true if all fields are default, false otherwise
func (finding *Finding) IsEmpty() bool {
	return finding.Match == "" && finding.Confidence == 0.0 &&
		finding.Location == "" && finding.Rule == "" && finding.Severity == 0
}

//Stringer for Finding struct
func (finding Finding) String() string {
	return fmt.Sprintf("Finding [%v] Location [%v] Rule [%v] Confidence [%v]", finding.Match, finding.Location, finding.Rule, finding.Confidence)
}

//ScanFile takes a path and a scan rule and returns a slice of findings
func ScanFile(target loader.ScanTarget, rule configuration.ScanRule, fchannel chan Finding) {
	input, _ := loader.GetBytesFromScanTarget(target)

	scanner := bufio.NewScanner(strings.NewReader(string(input)))
	index := 0
	for scanner.Scan() {
		finding := evaluateRule(scanner.Text(), rule)
		if !finding.IsEmpty() {
			finding.Location = fmt.Sprintf("%v : %v", target.Path, index)
			fchannel <- finding
		}
		index++
	}
	close(fchannel)
}

func evaluateRule(line string, rule configuration.ScanRule) Finding {
	//kb todo: these should be constructed somewhere else and referenced by pointer
	matcher, err := regexp.Compile(rule.Rule)
	if err != nil {
		panic(err)
	}
	match := matcher.FindString(line)
	if match != "" {
		return Finding{
			Match:      match,
			Location:   "NOTSET",
			Rule:       rule.Rule,
			Confidence: rule.Confidence,
			Severity:   rule.Severity}
	}
	return Finding{}
}
