// Package leak follows the use of the spider package. Leak takes a list of
// URLs and corresponding forms and makes a best-effort attempt to match source
// contact information to relevant fields.
//
// The basic strategy is to make post requests as quickly as possible and handle
// failure in the most non-disruptive way possible.
package leak

type contact struct {
	First  string
	Last   string
	Street string
	City   string
	State  string
	Email  string
	Zip    string // so we don't have to think about conversion
	Phone  string
}
