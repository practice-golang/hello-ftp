package model

type Host struct {
	Hostname string
	Port     string
	Username string
	Password string
	SrcBase  string
	DstBase  string
	// QueSheets []QueSheet
}
