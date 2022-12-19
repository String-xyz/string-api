package model

type FPVisitor struct {
	VisitorID  string
	Country    string
	State      string
	IPAddress  string
	Timestamp  int64
	Confidence float64
	Type       string
	UserAgent  string
}
