package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	filename := flag.String("in", "", ".ics file to modify")
	outFile := flag.String("out", "", "Output filename.")
	prefix := flag.String("pre", "[MANDATORY] ", "Prefix for title of each mandatory event.")
	remind := flag.Bool("remind", false, "If supplied, set a reminder MIN minutes before event.")
	min := flag.Int("min", 30, "Reminder N minutes before event.")
	flag.Parse()

	if *filename == "" {
		log.Fatalln("Error: please supply input filename")
	}

	if *outFile == "" {
		log.Fatalln("Error: please supply output filename")
	}

	if *remind {
		if *min < 0 {
			log.Fatalln("Error: min argument must be non-negative")
		}
	}

	sep := "\n"

	b, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatalln("Error reading file", filename)
	}
	s := string(b)
	ical := strings.Split(s, sep)

	mandatory := getMandatorySummaryIndices(ical)
	fmt.Printf("Adding '%s' to %d events ...\n", *prefix, len(mandatory))
	prependPrefixToSummaries(ical, *prefix, mandatory)

	if *remind {
		fmt.Printf("Adding %d-minute alarm for %d events ...\n", *min, len(mandatory))
		ical = addAlarm(ical, *min, mandatory)
	}

	newIcal := strings.Join(ical, sep)
	err = ioutil.WriteFile(*outFile, []byte(newIcal), 0644)
	if err != nil {
		log.Fatalln("Error writing file", outFile)
	}
}

func getMandatorySummaryIndices(iCalendar []string) []int {
	lastSummaryIdx := 0
	thisEventMandatory := false
	idsToModify := make([]int, 0)

	for j, line := range iCalendar {
		line = strings.ToLower(line)
		if strings.Contains(line, "begin:vevent") {
			thisEventMandatory = false
		}
		if strings.Contains(line, "summary") {
			lastSummaryIdx = j
		}
		if strings.Contains(line, "description") && strings.Contains(line, "mandatory") {
			thisEventMandatory = true
		}
		if strings.Contains(line, "end:vevent") && thisEventMandatory {
			idsToModify = append(idsToModify, lastSummaryIdx)
		}
	}
	return idsToModify
}

func prependPrefixToSummaries(iCalendar []string, prefix string, idsToChange []int) {
	summary := "SUMMARY:"
	pre := fmt.Sprintf("%s%s", summary, prefix)
	for _, j := range idsToChange {
		n := strings.Replace(iCalendar[j], summary, pre, 1)
		iCalendar[j] = n
	}
}

func addAlarm(iCalendar []string, minBeforeEvent int, idsToChange []int) []string {
	trigger := fmt.Sprintf("TRIGGER:-PT%dM", minBeforeEvent)
	alarm := []string{"BEGIN:VALARM", trigger, "ACTION:DISPLAY", "DESCRIPTION:Reminder", "END:VALARM"}

	// Append alarm just before the summary section of the vevent.
	// https://github.com/golang/go/wiki/SliceTricks#insertvector
	alarmLen := len(alarm)
	for offset, j := range idsToChange {
		// Account for change in indices whenever we append to the array.
		j += offset * alarmLen
		iCalendar = append(iCalendar[:j], append(alarm, iCalendar[j:]...)...)
	}
	return iCalendar
}
