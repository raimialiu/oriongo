package constants

type DbConnectionStatus string

const (
	CONNECTED              DbConnectionStatus = "connected"
	CONNECTING             DbConnectionStatus = "connecting"
	CONNECTION_NOT_STARTED DbConnectionStatus = "connection_not_started"
	CONNECTION_DROPPED     DbConnectionStatus = "connection_dropped"
	CONNECTION_ERROR       DbConnectionStatus = "connection_error"
)
