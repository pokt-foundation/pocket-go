package provider

// V1RPCRoute enum listing all the routes in the V1 RPC interface
type V1RPCRoute string

const (
	// ClientChallengeRoute represents client challenge route
	ClientChallengeRoute V1RPCRoute = "/v1/client/challenge"
	// ClientDispatchRoute represents client dispatch route
	ClientDispatchRoute V1RPCRoute = "/v1/client/dispatch"
	// ClientRawTXRoute represents cliente raw TX route
	ClientRawTXRoute V1RPCRoute = "/v1/client/rawtx"
	// ClientRelayRoute represents client realy route
	ClientRelayRoute V1RPCRoute = "/v1/client/relay"
	// QueryAccountRoute represents query account route
	QueryAccountRoute V1RPCRoute = "/v1/query/account"
	// QueryAccountTXsRoute represents query account TXs route
	QueryAccountTXsRoute V1RPCRoute = "/v1/query/accounttxs"
	// QueryAllParamsRoute represents query all params route
	QueryAllParamsRoute V1RPCRoute = "/v1/query/allparams"
	// QueryAppRoute represents query app route
	QueryAppRoute V1RPCRoute = "/v1/query/app"
	// QueryAppParamsRoute represents query app params route
	QueryAppParamsRoute V1RPCRoute = "/v1/query/appparams"
	// QueryAppsRoute represents query apps route
	QueryAppsRoute V1RPCRoute = "/v1/query/apps"
	// QueryBalanceRoute represents query balance route
	QueryBalanceRoute V1RPCRoute = "/v1/query/balance"
	// QueryBlockRoute represents query block route
	QueryBlockRoute V1RPCRoute = "/v1/query/block"
	// QueryBlockTXsRoute represents query block TXs route
	QueryBlockTXsRoute V1RPCRoute = "/v1/query/blocktxs"
	// QueryHeightRoute represents query height route
	QueryHeightRoute V1RPCRoute = "/v1/query/height"
	// QueryNodeRoute represents query node route
	QueryNodeRoute V1RPCRoute = "/v1/query/node"
	// QueryNodeClaimRoute represents query node claim route
	QueryNodeClaimRoute V1RPCRoute = "/v1/query/nodeclaim"
	// QueryNodeClaimsRoute represents query node claims route
	QueryNodeClaimsRoute V1RPCRoute = "/v1/query/nodeclaims"
	// QueryNodeParamsRoute represents query node params route
	QueryNodeParamsRoute V1RPCRoute = "/v1/query/nodeparams"
	// QueryNodeReceiptRoute represents query node receipt route
	QueryNodeReceiptRoute V1RPCRoute = "/v1/query/nodereceipt"
	// QueryNodeReceiptsRoute represents query node receipts route
	QueryNodeReceiptsRoute V1RPCRoute = "/v1/query/nodereceipts"
	// QueryNodesRoute represents query nodes route
	QueryNodesRoute V1RPCRoute = "/v1/query/nodes"
	// QueryPocketParamsRoute represents query pocket params route
	QueryPocketParamsRoute V1RPCRoute = "/v1/query/pocketparams"
	// QuerySupplyRoute represents query supply route
	QuerySupplyRoute V1RPCRoute = "/v1/query/supply"
	// QuerySupportedChainsRoute represents query supported chains route
	QuerySupportedChainsRoute V1RPCRoute = "/v1/query/supportedchains"
	// QueryTXRoute represents query TX route
	QueryTXRoute V1RPCRoute = "/v1/query/tx"
	// QueryUpgradeRoute represents query upgrade route
	QueryUpgradeRoute V1RPCRoute = "/v1/query/upgrade"
)
