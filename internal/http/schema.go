package http

import "strings"

type VoterInfoResponse struct {
	Kind             string            `json:"kind"`
	Election         Election          `json:"election"`
	OtherElections   []Election        `json:"otherElections"`
	NormalizedInput  Address           `json:"normalizedInput"`
	PollingLocations []PollingLocation `json:"pollingLocations"`
	EarlyVoteSites   []VoteSite        `json:"earlyVoteSites"`
	DropOffLocations []VoteSite        `json:"dropOffLocations"`
	Contests         []Contest         `json:"contests"`
	State            []State           `json:"state"`
	MailOnly         bool              `json:"mailOnly"`
}

// Election Resource
type Election struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ElectionDay   string `json:"electionDay"`
	OcdDivisionId string `json:"ocdDivisionId"`
}

// Address Resource
type Address struct {
	LocationName string `json:"locationName"`
	Line1        string `json:"line1"`
	Line2        string `json:"line2"`
	Line3        string `json:"line3"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
}

func (a Address) String() string {
	var b strings.Builder
	separator := ""

	if a.LocationName != "" {
		b.WriteString(a.LocationName)
		separator = ", "
	}
	if a.Line1 != "" {
		b.WriteString(separator)
		b.WriteString(a.Line1)
		separator = ", "
	}
	if a.Line2 != "" {
		b.WriteString(separator)
		b.WriteString(a.Line2)
	}
	if a.Line3 != "" {
		b.WriteString(separator)
		b.WriteString(a.Line3)
	}
	if a.City != "" {
		b.WriteString(separator)
		b.WriteString(a.City)
		separator = ", "
	}
	if a.State != "" {
		b.WriteString(separator)
		b.WriteString(a.State)
	}
	if a.Zip != "" {
		// Add a space before the zip if there's a state present
		if a.State != "" {
			b.WriteString(" ")
		} else {
			b.WriteString(separator)
		}
		b.WriteString(a.Zip)
	}

	// If no components were written, return an empty string
	if b.Len() == 0 {
		return ""
	}

	return b.String()
}

// PollingLocation Resource
type PollingLocation struct {
	Address       Address  `json:"address"`
	Notes         string   `json:"notes"`
	PollingHours  string   `json:"pollingHours"`
	Name          string   `json:"name"`
	VoterServices string   `json:"voterServices"`
	StartDate     string   `json:"startDate"`
	EndDate       string   `json:"endDate"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	Sources       []Source `json:"sources"`
}

// VoteSite Resource (used for earlyVoteSites and dropOffLocations)
type VoteSite struct {
	Address       Address  `json:"address"`
	Notes         string   `json:"notes"`
	PollingHours  string   `json:"pollingHours"`
	Name          string   `json:"name"`
	VoterServices string   `json:"voterServices"`
	StartDate     string   `json:"startDate"`
	EndDate       string   `json:"endDate"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	Sources       []Source `json:"sources"`
}

// Contest Resource
type Contest struct {
	Type                       string      `json:"type"` // TODO: Convert to ENUM
	PrimaryParty               string      `json:"primaryParty"`
	ElectorateSpecifications   string      `json:"electorateSpecifications"`
	Special                    string      `json:"special"`
	BallotTitle                string      `json:"ballotTitle"`
	Office                     string      `json:"office"`
	Level                      []string    `json:"level"`
	Roles                      []string    `json:"roles"`
	District                   District    `json:"district"`
	NumberElected              string      `json:"numberElected"`   // Schema says long, but API returns a string
	NumberVotingFor            string      `json:"numberVotingFor"` // Schema says long, but API returns a string
	BallotPlacement            string      `json:"ballotPlacement"` // Schema says long, but API returns a string
	Candidates                 []Candidate `json:"candidates"`
	ReferendumTitle            string      `json:"referendumTitle"`
	ReferendumSubtitle         string      `json:"referendumSubtitle"`
	ReferendumUrl              string      `json:"referendumUrl"`
	ReferendumBrief            string      `json:"referendumBrief"`
	ReferendumText             string      `json:"referendumText"`
	ReferendumProStatement     string      `json:"referendumProStatement"`
	ReferendumConStatement     string      `json:"referendumConStatement"`
	ReferendumPassageThreshold string      `json:"referendumPassageThreshold"`
	ReferendumEffectOfAbstain  string      `json:"referendumEffectOfAbstain"`
	ReferendumBallotResponses  []string    `json:"referendumBallotResponses"`
	Sources                    []Source    `json:"sources"`
}

// Candidate Resource
type Candidate struct {
	Name          string    `json:"name"`
	Party         string    `json:"party"`
	CandidateUrl  string    `json:"candidateUrl"`
	Phone         string    `json:"phone"`
	PhotoUrl      string    `json:"photoUrl"`
	Email         string    `json:"email"`
	OrderOnBallot int64     `json:"orderOnBallot"`
	Channels      []Channel `json:"channels"`
}

// Channel Resource
type Channel struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// District Resource
type District struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	ID    string `json:"id"`
}

// Source Resource
type Source struct {
	Name     string `json:"name"`
	Official bool   `json:"official"`
}

// State Resource
type State struct {
	Name                       string                     `json:"name"`
	ElectionAdministrationBody ElectionAdministrationBody `json:"electionAdministrationBody"`
	LocalJurisdiction          *AdministrationRegion      `json:"local_jurisdiction"`
	Sources                    []Source                   `json:"sources"`
}

// ElectionAdministrationBody Resource
type ElectionAdministrationBody struct {
	Name                                string             `json:"name"`
	ElectionInfoUrl                     string             `json:"electionInfoUrl"`
	ElectionRegistrationUrl             string             `json:"electionRegistrationUrl"`
	ElectionRegistrationConfirmationUrl string             `json:"electionRegistrationConfirmationUrl"`
	ElectionNoticeText                  string             `json:"electionNoticeText"`
	ElectionNoticeUrl                   string             `json:"electionNoticeUrl"`
	AbsenteeVotingInfoUrl               string             `json:"absenteeVotingInfoUrl"`
	VotingLocationFinderUrl             string             `json:"votingLocationFinderUrl"`
	BallotInfoUrl                       string             `json:"ballotInfoUrl"`
	ElectionRulesUrl                    string             `json:"electionRulesUrl"`
	VoterServices                       []string           `json:"voter_services"`
	HoursOfOperation                    string             `json:"hoursOfOperation"`
	CorrespondenceAddress               Address            `json:"correspondenceAddress"`
	PhysicalAddress                     Address            `json:"physicalAddress"`
	ElectionOfficials                   []ElectionOfficial `json:"electionOfficials"`
}

// ElectionOfficial Resource
type ElectionOfficial struct {
	Name              string `json:"name"`
	Title             string `json:"title"`
	OfficePhoneNumber string `json:"officePhoneNumber"`
	FaxNumber         string `json:"faxNumber"`
	EmailAddress      string `json:"emailAddress"`
}

// AdministrationRegion Resource
type AdministrationRegion struct {
	Name                       string                     `json:"name"`
	ElectionAdministrationBody ElectionAdministrationBody `json:"electionAdministrationBody"`
}
