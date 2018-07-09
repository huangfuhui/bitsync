package coin

import "bitsync/object"

type Coin struct {
	object.Base
	Name              string
	NameCn            string
	FullName          string
	Icoin             string
	OfficialWebsite   string
	WhitePaper        string
	IssueAmount       int64
	FlowAmount        int64
	IcoPrice          string
	BlockchainBrowser string
	Introduction      string
}
