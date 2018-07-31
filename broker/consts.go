package broker

const (
	// HTTPDefaultAddr is the addr the broker will listen at by default
	HTTPDefaultAddr = "0.0.0.0"
	// HTTPDefaultPort is the port the broker will listen at by default
	HTTPDefaultPort = 3080
	// keySize specifies the size of RSA keys
	keySize       = 2048
	restPublicKey = "publicKey"
	// tokenChars contains all possible characters of a token
	tokenChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"
	// tokenSize is the size in characters of the auth tokens
	tokenSize = 32
)
