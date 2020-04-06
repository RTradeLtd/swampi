package swampi

import "fmt"

// APICall represents a type string indicating various APIs and their paths
type APICall string

// String implements the stringer interface
// but also returns the http path for this api
func (a APICall) String() string {
	return string(a)
}

// Method returns the http method associated with the call
func (a APICall) Method() string {
	switch a {
	case SingleFileUpload:
		return "POST"
	default:
		return ""
	}
}

// ContentType returns the content type associated with this request
func (a APICall) ContentType() string {
	switch a {
	case SingleFileUpload:
		return "text/plain"
	default:
		return ""
	}
}

// ParseArgs is used to format request arguments
func (a APICall) ParseArgs(args ...interface{}) string {
	switch a {
	case SingleFileUpload:
		return a.String()
	case SingleFileDownload:
		return fmt.Sprintf(a.String(), args...)
	default:
		return ""
	}
}

const (
	// SingleFileUpload is an api to upload a singular file
	SingleFileUpload = APICall("/bzz:/")
	// SingleFileDownload is used to download a singular file from swarm
	SingleFileDownload = APICall("/bzz:/%s/")
)

var (
	// APICalls is just a helper slice containing all known API calls
	APICalls = []APICall{SingleFileUpload}
)