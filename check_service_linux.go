﻿package snclient

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	AvailableChecks["check_service"] = CheckEntry{"check_service", new(CheckService)}
}

type CheckService struct {
	noCopy noCopy
	data   CheckData
}

var ServiceStates = map[string]string{
	"stopped":      "1",
	"dead":         "1",
	"startpending": "2",
	"stoppending":  "3",
	"running":      "4",
	"started":      "4",
}

/* check_service_linux
 * Description: Checks the state of a service on the host.
 * Thresholds: status
 * Units: stopped, dead, startpending, stoppedpending, running, started
 */
func (l *CheckService) Check(args []string) (*CheckResult, error) {
	// default state: OK
	state := int64(0)
	l.data.detailSyntax = "%(service)=%(state)"
	l.data.topSyntax = "%(crit_list), delayed (%(warn_list))"
	l.data.okSyntax = "All %(count) service(s) are ok."
	argList := ParseArgs(args, &l.data)
	var output string
	var services []string
	var checkData map[string]string

	// parse treshold args
	for _, arg := range argList {
		switch arg.key {
		case "service":
			services = append(services, arg.value)
		}
	}

	metrics := make([]*CheckMetric, 0, len(services))
	okList := make([]string, 0, len(services))
	warnList := make([]string, 0, len(services))
	critList := make([]string, 0, len(services))

	warnTreshold.value = ServiceStates[warnTreshold.value]
	critTreshold.value = ServiceStates[critTreshold.value]

	for _, service := range services {
		// collect service state
		out, err := exec.Command("systemctl", "status", fmt.Sprintf("%s.service", service)).Output()
		if err != nil {
			return &CheckResult{
				State:  CheckExitUnknown,
				Output: fmt.Sprintf("Service %s not found: %s", service, err),
			}, nil
		}

		re := regexp.MustCompile(`Active:\s*[A-Za-z]+\s*\(([A-Za-z]+)\)`)
		match := re.FindStringSubmatch(string(out))

		stateStr := ServiceStates[match[1]]
		stateFloat, _ := strconv.ParseFloat(stateStr, 64)

		metrics = append(metrics, &CheckMetric{Name: service, Value: stateFloat})

		mdata := []MetricData{{name: "state", value: stateStr}}
		sdata := map[string]string{
			"service": service,
			"state":   stateStr,
		}

		// compare ram metrics to tresholds
		if CompareMetrics(mdata, l.data.critTreshold) {
			critList = append(critList, ParseSyntax(l.data.detailSyntax, sdata))

			continue
		}

		if CompareMetrics(mdata, l.data.warnTreshold) {
			warnList = append(warnList, ParseSyntax(l.data.detailSyntax, sdata))

			continue
		}

		okList = append(okList, ParseSyntax(l.data.detailSyntax, sdata))
	}

	if len(critList) > 0 {
		state = CheckExitCritical
	} else if len(warnList) > 0 {
		state = CheckExitWarning
	}

	checkData = map[string]string{
		"status":    strconv.FormatInt(state, 10),
		"count":     strconv.FormatInt(int64(len(services)), 10),
		"ok_list":   strings.Join(okList, ", "),
		"warn_list": strings.Join(warnList, ", "),
		"crit_list": strings.Join(critList, ", "),
	}

	if state == CheckExitOK {
		output = ParseSyntax(l.data.okSyntax, checkData)
	} else {
		output = ParseSyntax(l.data.topSyntax, checkData)
	}

	return &CheckResult{
		State:   state,
		Output:  output,
		Metrics: metrics,
	}, nil
}