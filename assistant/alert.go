package assistant

// All alerts should return a non-empty string if a notification is required.
type handlerFunc func() (message string)

type alert struct {
	name    string
	handler handlerFunc
}

// TODO: more alerts.
var alerts = []*alert{
	{
		name:    "Crypto",
		handler: cryptoAlertHandler,
	},
}
