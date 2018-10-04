package peer

const (
	brokerIP        = "127.0.0.1"
	brokerPort      = 3080
	delimiter       = "%;%"
	filenamePeer    = "peer.json"
	filenamePrv     = "prvkey.pem"
	filenamePub     = "pubkey.pem"
	keySize         = 1024 // 4096 significantly increases test duration
	listenMultiAddr = "/ip4/0.0.0.0/tcp/3001"
	permissionDir   = 0750
	permissionFile  = 0750
	protocolID      = "/sakaban/v0.0.0"
)
