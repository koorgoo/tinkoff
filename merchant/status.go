package merchant

const (
	StatusNew            = "NEW"
	StatusCanceled       = "CANCELED"
	StatusPreauthorizing = "PREAUTHORIZING"
	StatusFormshowed     = "FORMSHOWED"
	// StatusAuthorizing     = "AUTHORIZING"
	// Status3DSChecking     = "3DS_CHECKING"
	// Status3DSChecked      = "3DS_CHECKED"
	// StatusAuthorized      = "AUTHORIZED"
	// StatusReversing       = "REVERSING"
	// StatusReversed        = "REVERSED"
	StatusConfirming      = "CONFIRMING"
	StatusConfirmed       = "CONFIRMED"
	StatusRefunding       = "REFUNDING"
	StatusRefunded        = "REFUNDED"
	StatusPartialRefunded = "PARTIAL_REFUNDED"
	StatusRejected        = "REJECTED"
	StatusUnknown         = "UNKNOWN"
)
