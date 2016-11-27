package merchant

type Response struct {
	TerminalKey string
	OrderId     string
	PaymentId   string // XXX: Docs lie about Number(20) type.
	Success     bool
	Status      string
	ErrorCode   string
	Message     string
	Details     string
}

type InitResponse struct {
	Response
	Amount     int64
	PaymentURL string
}

type CancelResponse struct {
	Response
	OriginalAmount int64
	NewAmount      int64
}

type GetStateResponse struct {
	Response
}
