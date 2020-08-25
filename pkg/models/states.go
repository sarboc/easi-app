package models

const (
	// CreateNewIntake is
	CreateNewIntake = iota
	// EditDraftIntake is
	EditDraftIntake
	// SubmitIntake is
	SubmitIntake
	//DecideIntakeAccepted is
	DecideIntakeAccepted
	// DecideIntakeApproved is
	DecideIntakeApproved
	// DecideIntakeClosed is
	DecideIntakeClosed
	// WithdrawIntake is
	WithdrawIntake
	// BadTransition is
	BadTransition
)
