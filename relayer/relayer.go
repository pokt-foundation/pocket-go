// Package relayer is a helper for doing relays with simpler input that just using the package Provider
// Underneath uses the package Provider
package relayer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/big"

	"golang.org/x/crypto/sha3"

	"github.com/pokt-foundation/pocket-go/provider"
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
)

// Provider interface representing provider functions necessary for Relayer Package
type Provider interface {
	RelayWithCtx(ctx context.Context, rpcURL string, input *provider.RelayInput, options *provider.RelayRequestOptions) (*provider.RelayOutput, error)
}

// Signer interface representing signer functions necessary for Relayer Package
type Signer interface {
	Sign(payload []byte) (string, error)
}

// Relayer implementation of relayer interface
type Relayer struct {
	signer   Signer
	provider Provider
}

// NewRelayer returns instance of Relayer with given input.
// signer is the `client` that was used during AAT generation.
// The signer is often synonymous with the `gateway` which may or may not be the
// same as the `application` depending on how the AAT was generated.
func NewRelayer(signer Signer, provider Provider) *Relayer {
	return &Relayer{
		signer:   signer,
		provider: provider,
	}
}

func (r *Relayer) validateRelayRequest(input *Input) error {
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

	if input.Session.Header == (provider.SessionHeader{}) {
		return ErrNoSessionHeader
	}

	return nil
}

func getNode(input *Input) (*provider.Node, error) {
	if input.Node == nil {
		return GetRandomSessionNode(input.Session)
	}

	if !IsNodeInSession(input.Session, input.Node) {
		return nil, ErrNodeNotInSession
	}

	return input.Node, nil
}

// getSignedProofBytes returns the relay proof bytes signed by the signer
func (r *Relayer) getSignedProofBytes(proof *provider.RelayProof) (string, error) {
	// Prepare the relay proof bytes to be signed
	proofBytes, err := GenerateProofBytes(proof)
	if err != nil {
		return "", err
	}

	// Sign the relay proof bytes using the private key of the the signer, also
	// known as the client from the AAT generation process.
	return r.signer.Sign(proofBytes)
}

// buildRelay creates a Pocket relay using the RPC payload provided
func (r *Relayer) buildRelay(
	servicerNode *provider.Node,
	input *Input,
	options *provider.RelayRequestOptions,
) (*provider.RelayInput, error) {
	relayPayload := &provider.RelayPayload{
		Data:    input.Data,
		Method:  input.Method,
		Path:    input.Path,
		Headers: input.Headers,
	}

	relayMeta := &provider.RelayMeta{
		BlockHeight: input.Session.Header.SessionHeight,
	}

	requestHash := &RequestHash{
		Payload: relayPayload,
		Meta:    relayMeta,
	}

	hashedReq, err := HashRequest(requestHash)
	if err != nil {
		return nil, err
	}

	entropy, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	// Prepare the RelayProof object
	relayProof := provider.RelayProof{
		RequestHash:        hashedReq,
		Entropy:            entropy.Int64(),
		SessionBlockHeight: input.Session.Header.SessionHeight,
		ServicerPubKey:     servicerNode.PublicKey,
		Blockchain:         input.Blockchain,
		AAT:                input.PocketAAT,
	}

	// Sign the relay proof object
	signedProofBytes, err := r.getSignedProofBytes(&relayProof)
	if err != nil {
		return nil, err
	}

	// Update the Signature of the RelayProof
	relayProof.Signature = signedProofBytes

	return &provider.RelayInput{
		Payload: relayPayload,
		Meta:    relayMeta,
		Proof:   &relayProof,
	}, nil
}

// Relay does relay request with given input
// Will always return with an output that includes the status code from the request
func (r *Relayer) Relay(input *Input, options *provider.RelayRequestOptions) (*Output, error) {
	return r.RelayWithCtx(context.Background(), input, options)
}

// RelayWithCtx does relay request with given input
// Will always return with an output that includes the status code from the request
func (r *Relayer) RelayWithCtx(ctx context.Context, input *Input, options *provider.RelayRequestOptions) (*Output, error) {
	defaultOutput := &Output{
		RelayOutput: &provider.RelayOutput{
			StatusCode: provider.DefaultStatusCode,
		},
	}

	err := r.validateRelayRequest(input)
	if err != nil {
		return defaultOutput, err
	}

	node, err := getNode(input)
	if err != nil {
		return defaultOutput, err
	}

	relayInput, err := r.buildRelay(node, input, options)
	if err != nil {
		return defaultOutput, err
	}

	relayOutput, relayErr := r.provider.RelayWithCtx(ctx, node.ServiceURL, relayInput, options)
	if relayErr != nil {
		defaultOutput.RelayOutput = relayOutput
		return defaultOutput, relayErr
	}

	return &Output{
		RelayOutput: relayOutput,
		Proof:       relayInput.Proof,
		Node:        node,
	}, nil
}

// GetRandomSessionNode returns a random node from given session
func GetRandomSessionNode(session *provider.Session) (*provider.Node, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(session.Nodes))))
	if err != nil {
		return nil, err
	}

	node := session.Nodes[index.Int64()]
	return &node, nil
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

	if _, err = hasher.Write(marshaledProof); err != nil {
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
