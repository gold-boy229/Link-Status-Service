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
