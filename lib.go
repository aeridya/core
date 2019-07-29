package core

import "fmt"

const (
	//NAME is the name of the librar
	NAME = "Aeridya"
	//MAJORVER is the Major Version
	MAJORVER = "1"
	//MINORVER is the Minor Version
	MINORVER = "0"
	//RELEASEVER is the Release Version
	RELEASEVER = "2"
	//DESC is a description of the library
	DESC = "Server and CMS"
)

// Version returns a formatted string of the name/version number
func Version() string {
	return fmt.Sprintf("%s v%s.%s.%s", NAME, MAJORVER, MINORVER, RELEASEVER)
}

// Info returns a formatted string of Version and the Description
func Info() string {
	return fmt.Sprintf("%s\n\t%s", Version(), DESC)
}
