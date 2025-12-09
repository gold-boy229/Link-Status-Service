package entity

type LinkGetStatusParams struct {
	Links []string
}

type LinkGetStatusResult struct {
	LinkStates []LinkState
	LinkNum    int
}

type LinkState struct {
	Link        string
	IsAvailable bool
}

type LinksStateResult struct {
	LinkStates []LinkState
	Error      error
}

type LinkBuildPDSParams struct {
	LinkNums []int
}

type LinkBuildPDSResult struct {
	LinkStates []LinkState
}

type LinkStatus struct {
	Address string
	Status  string
}
