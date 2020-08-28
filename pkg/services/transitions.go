package services

// StateMachineTransitions holds all the transition functions
type StateMachineTransitions struct {
	CreateIntake CreateIntakeTransition
	UpdateDraft  UpdateIntakeTransition
	Submit       UpdateIntakeTransition
	Approve      UpdateIntakeTransition
	Accept       UpdateIntakeTransition
	Close        UpdateIntakeTransition
	Archive      DeleteIntakeTransition
}
