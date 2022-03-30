package relayer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/rand"

	"github.com/pokt-foundation/pocket-go/pkg/provider"
	"github.com/pokt-foundation/pocket-go/pkg/signer"
	"golang.org/x/crypto/sha3"
)

var (
	// ErrNoSigner error when no signer is provided
	ErrNoSigner = errors.New("no signer provided")
	// ErrNoSession error when no session is provided
	ErrNoSession = errors.New("no session provided")
	// ErrNoSessionHeader error when no session header is provided
	ErrNoSessionHeader = errors.New("no session header provided")
	// ErrNoProvider error when no provider is provided
	ErrNoProvider = errors.New("no provider provided")
	// ErrNoPocketAAT error when no Pocket AAT is provided
	ErrNoPocketAAT = errors.New("no Pocket AAT provided")
	// ErrSessionHasNoNodes error when provided session has no nodes
	ErrSessionHasNoNodes = errors.New("session has no nodes")
	// ErrNodeNotInSession error when given node is not in session
	ErrNodeNotInSession = errors.New("node not in session")
	// ErrUnexpectedErrorResponse error when relay error response is empty
	ErrUnexpectedErrorResponse = errors.New("unexpected error response")
)

// PocketRelayer implementation of relayer interface
type PocketRelayer struct {
	signer   signer.Signer
	provider provider.Provider
}

// NewPocketRelayer returns instance of PocketRelayer with given input
func NewPocketRelayer(signer signer.Signer, provider provider.Provider) *PocketRelayer {
	return &PocketRelayer{
		signer:   signer,
		provider: provider,
	}
}

func (r *PocketRelayer) validateRelayRequest(input *RelayInput) error {
	if r.signer == nil {
		return ErrNoSigner
	}

	if r.provider == nil {
		return ErrNoProvider
	}

	if input.Session == nil {
		return ErrNoSession
	}

	if input.PocketAAT == nil {
		return ErrNoPocketAAT
	}

	if len(input.Session.Nodes) == 0 {
		return ErrSessionHasNoNodes
	}

	if input.Session.Header == nil {
		return ErrNoSessionHeader
	}

	return nil
}

func getNode(input *RelayInput) (*provider.Node, error) {
	if input.Node == nil {
		return GetRandomSessionNode(input.Session), nil
	}

	if !IsNodeInSession(input.Session, input.Node) {
		return nil, ErrNodeNotInSession
	}

	return input.Node, nil
}

func (r *PocketRelayer) getSignedProofBytes(proof *provider.RelayProof) (string, error) {
	proofBytes, err := GenerateProofBytes(proof)
	if err != nil {
		return "", err
	}

	return r.signer.GetKeyManager().Sign(proofBytes)
}

// Relay does relay request with given input
func (r *PocketRelayer) Relay(input *RelayInput, options *provider.RelayRequestOptions) (*RelayerResponse, error) {
	err := r.validateRelayRequest(input)
	if err != nil {
		return nil, err
	}

	node, err := getNode(input)
	if err != nil {
		return nil, err
	}

	relayPayload := &provider.RelayPayload{
		Data:    input.Data,
		Method:  input.Method,
		Path:    input.Path,
		Headers: input.Headers,
	}

	relayMeta := &provider.RelayMeta{
		BlockHeight: input.Session.Header.SessionHeight,
	}

	hashedReq, err := HashRequest(&RequestHash{
		Payload: relayPayload,
		Meta:    relayMeta,
	})
	if err != nil {
		return nil, err
	}

	entropy := rand.Intn(math.MaxInt)

	signedProofBytes, err := r.getSignedProofBytes(&provider.RelayProof{
		RequestHash:        hashedReq,
		Entropy:            entropy,
		SessionBlockHeight: input.Session.Header.SessionHeight,
		ServicerPubKey:     node.PublicKey,
		Blockchain:         input.Blockchain,
		AAT:                input.PocketAAT,
	})
	if err != nil {
		return nil, err
	}

	relayProof := &provider.RelayProof{
		RequestHash:        hashedReq,
		Entropy:            entropy,
		SessionBlockHeight: input.Session.Header.SessionHeight,
		ServicerPubKey:     node.PublicKey,
		Blockchain:         input.Blockchain,
		AAT:                input.PocketAAT,
		Signature:          signedProofBytes,
	}

	relay := &provider.Relay{
		Payload: relayPayload,
		Meta:    relayMeta,
		Proof:   relayProof,
	}

	relayOutput, err := r.provider.Relay(node.ServiceURL, relay, options)
	if err != nil {
		return nil, err
	}

	return &RelayerResponse{
		Response: relayOutput,
		Proof:    relayProof,
		Node:     node,
	}, nil
}

// GetRandomSessionNode returns a random node from given session
func GetRandomSessionNode(session *provider.Session) *provider.Node {
	return session.Nodes[rand.Intn(len(session.Nodes))]
}

// IsNodeInSession verifies if given node is in given session
func IsNodeInSession(session *provider.Session, node *provider.Node) bool {
	for _, sessionNode := range session.Nodes {
		if sessionNode.PublicKey == node.PublicKey {
			return true
		}
	}

	return false
}

// GenerateProofBytes returns relay proof as encoded bytes
func GenerateProofBytes(proof *provider.RelayProof) ([]byte, error) {
	token, err := HashAAT(proof.AAT)
	if err != nil {
		return nil, err
	}

	proofMap := &relayProofForSignature{
		RequestHash:        proof.RequestHash,
		Entropy:            proof.Entropy,
		SessionBlockHeight: proof.SessionBlockHeight,
		ServicerPubKey:     proof.ServicerPubKey,
		Blockchain:         proof.Blockchain,
		Token:              token,
		Signature:          "",
	}

	marshaledProof, err := json.Marshal(proofMap)
	if err != nil {
		return nil, err
	}

	hasher := sha3.New256()

	_, err = hasher.Write(marshaledProof)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

// HashAAT returns Pocket AAT as hashed string
func HashAAT(aat *provider.PocketAAT) (string, error) {
	tokenToSend := *aat
	tokenToSend.Signature = ""

	marshaledAAT, err := json.Marshal(tokenToSend)
	if err != nil {
		return "", err
	}

	hasher := sha3.New256()

	_, err = hasher.Write(marshaledAAT)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// HashRequest creates the request hash from its structure
func HashRequest(reqHash *RequestHash) (string, error) {
	marshaledReqHash, err := json.Marshal(reqHash)
	if err != nil {
		return "", err
	}

	hasher := sha3.New256()

	_, err = hasher.Write(marshaledReqHash)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
