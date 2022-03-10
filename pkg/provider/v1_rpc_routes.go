package provider

// V1RPCRoute enum listing all the routes in the V1 RPC interface
type V1RPCRoute string

const (
	// ClientChallenge represents client challenge route
	ClientChallenge V1RPCRoute = "/v1/client/challenge"
	// ClientDispatch represents client dispatch route
	ClientDispatch V1RPCRoute = "/v1/client/dispatch"
	// ClientRawTX represents cliente raw TX route
	ClientRawTX V1RPCRoute = "/v1/client/rawtx"
	// ClientRelay represtns client realy route
	ClientRelay V1RPCRoute = "/v1/client/relay"
	// QueryAccount represents query account route
	QueryAccount V1RPCRoute = "/v1/query/account"
	// QueryAccountsTXs represents query account TXs route
	QueryAccountTXs V1RPCRoute = "/v1/query/accounttxs"
	// QueryAllParams represents query all params route
	QueryAllParams V1RPCRoute = "/v1/query/allparams"
	// QueryApp represents query app route
	QueryApp V1RPCRoute = "/v1/query/app"
	// QueryAppParams represents query app params route
	QueryAppParams V1RPCRoute = "/v1/query/appparams"
	// QueryApps represents query apps route
	QueryApps V1RPCRoute = "/v1/query/apps"
	// QueryBalance represents query balance route
	QueryBalance V1RPCRoute = "/v1/query/balance"
	// QueryBlock represents query block route
	QueryBlock V1RPCRoute = "/v1/query/block"
	// QueryBlocksTXs represents query block TXs route
	QueryBlockTXs V1RPCRoute = "/v1/query/blocktxs"
	// QueryHeight represents query height route
	QueryHeight V1RPCRoute = "/v1/query/height"
	// QueryNode represents query node route
	QueryNode V1RPCRoute = "/v1/query/node"
	// QueryNodeClaim represents query node claim route
	QueryNodeClaim V1RPCRoute = "/v1/query/nodeclaim"
	// QueryNodeClaims represents query node claims route
	QueryNodeClaims V1RPCRoute = "/v1/query/nodeclaims"
	// QueryNodeParams represents query node params route
	QueryNodeParams V1RPCRoute = "/v1/query/nodeparams"
	// QueryNodeReceipt represents query node receipt route
	QueryNodeReceipt V1RPCRoute = "/v1/query/nodereceipt"
	// QueryNodeReceipts represents query node receipts route
	QueryNodeReceipts V1RPCRoute = "/v1/query/nodereceipts"
	// QueryNodes represents query nodes route
	QueryNodes V1RPCRoute = "/v1/query/nodes"
	// QueryPocketParams represents query pocket params route
	QueryPocketParams V1RPCRoute = "/v1/query/pocketparams"
	// QuerySupply represents query supply route
	QuerySupply V1RPCRoute = "/v1/query/supply"
	// QuerySupportedChains represents query supported chains route
	QuerySupportedChains V1RPCRoute = "/v1/query/supportedchains"
	// QueryTX represents query TX route
	QueryTX V1RPCRoute = "/v1/query/tx"
	// QueryUpgrade reprents query upgrade route
	QueryUpgrade V1RPCRoute = "/v1/query/upgrade"
)
