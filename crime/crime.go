package crime

import (
	"fmt"
	"strings"
	"time"
)

// Crime structs hold information about criminal activity reported by Clery
// act reports
type Crime struct {
	// DateReported records when the criminal activity was disclosed to the
	// police
	DateReported time.Time

	// DateOccurredStart records when the criminal activity started taking
	// place
	DateOccurredStart time.Time

	// DateOccurredEnd records when the criminal activity stopped taking
	// place
	DateOccurredEnd time.Time

	// ReportSuperID is the first portion of police report ID associated
	// with the reported crime.
	ReportSuperID uint

	// ReportID is the second portion of the police report ID associated
	// with the reported crime
	ReportID uint

	// Location is the place where the criminal activity occurred
	Location string

	// Incidents holds the official classifications of the criminal
	// activity which took place
	Incidents []string

	// Descriptions holds any details about the specific incidents which
	// took place
	Descriptions []string

	// Remediation is the action taken by the institution who reported the
	// crime to deal with the criminal activity
	Remediation string
}

func (c Crime) String() string {
	return fmt.Sprintf("Reported: %s\n"+
		"Occurred Start: %s\n"+
		"Occurred End: %s\n"+
		"ID: %d-%d\n"+
		"Location: %s\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s",
		c.DateReported,
		c.DateOccurredStart,
		c.DateOccurredEnd,
		c.ReportSuperID,
		c.ReportID,
		c.Location,
		strings.Join(c.Incidents, ","),
		strings.Join(c.Descriptions, ","),
		c.Remediation)
}
