package entity

type LinkGetStatus_Params struct {
	Links []string
}

type LinkGetStatus_Result struct {
	LinkStates []LinkState
	LinkNum    int
}

type LinkState struct {
	Link        string
	IsAvailable bool
}

type LinksState_Result struct {
	LinkStates []LinkState
	Error      error
}

type LinkBuildPDS_Params struct {
	LinkNums []int
}

type LinkBuildPDS_Result struct {
	LinkStates []LinkState
}

type LinkStatus struct {
	Address string
	Status  string
}
