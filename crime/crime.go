package crime

import (
	"fmt"
	"strings"
)

// Crime structs hold information about criminal activity reported by Clery
// act reports
type Crime struct {
	// DateReported records when the criminal activity was disclosed to the
	// police
	DateReported string

	// DateOccurred records when the criminal activity took place
	DateOccurred string

	// ReportID is the unique ID used to identify the criminal activity
	// report in the reporting institution's internal system
	ReportID string

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
		"Occurred: %s\n"+
		"ID: %s\n"+
		"Location: %s\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s",
		c.DateReported,
		c.DateOccurred,
		c.ReportID,
		c.Location,
		strings.Join(c.Incidents, ","),
		strings.Join(c.Descriptions, ","),
		c.Remediation)
}
