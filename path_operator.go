package natssecretsengine

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/nats-io/nkeys"
)

const (
	operatorStoragePath = "operator"
)

// natsOperator stores informations about the operator.
type natsOperator struct {
	Name       string `json:"name"`
	PublicKey  string `json:"public-key"`
	PrivateKey string `json:"private-key"`
	Seed       string `json:"seed"`
}

// pathOperator extends the Vault API with a `/operator`
// endpoint for the natsBackend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the operatoruration.
func pathOperator(b *natsBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "operator",
			Fields:  map[string]*framework.FieldSchema{},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathOperatorRead,
				},
			},
			ExistenceCheck:  b.pathOperatorExistenceCheck,
			HelpSynopsis:    pathOperatorHelpSyn,
			HelpDescription: pathOperatorHelpDesc,
		},
		{
			Pattern: "operator/generate",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the operator",
					Required:    true,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathOperatorGenerate,
			},
			HelpSynopsis:    pathOperatorGenerateHelpSyn,
			HelpDescription: pathOperatorGenerateHelpDesc,
		},
	}
}

// pathOperatorExistenceCheck verifies if the operatoruration exists.
func (b *natsBackend) pathOperatorExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// pathOperatorRead reads the operatoruration and outputs non-sensitive information.
func (b *natsBackend) pathOperatorRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	operator, err := getOperator(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"name":        operator.Name,
			"public-key":  operator.PublicKey,
			"private-key": operator.PrivateKey,
			"seed":        operator.Seed,
		},
	}, nil
}

// pathOperatorGenerate generates a new operator for the backend
func (b *natsBackend) pathOperatorGenerate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	operator, err := getOperator(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if operator == nil {
		operator = new(natsOperator)
	}

	// store provided operator name
	if name, ok := data.GetOk("name"); ok {
		operator.Name = name.(string)
	} else if !ok {
		return nil, fmt.Errorf("missing name for operator")
	}

	// create operator keypair
	keypair, err := nkeys.CreateOperator()
	if err != nil {
		return nil, err
	}

	// store public key
	operator.PublicKey, err = keypair.PublicKey()
	if err != nil {
		return nil, err
	}

	// store private key
	var pKey []byte
	pKey, err = keypair.PrivateKey()
	if err != nil {
		return nil, err
	}
	operator.PrivateKey = base64.StdEncoding.EncodeToString(pKey)

	// store seed
	var seed []byte
	seed, err = keypair.Seed()
	if err != nil {
		return nil, err
	}
	operator.Seed = base64.StdEncoding.EncodeToString(seed)

	// convert to json
	entry, err := logical.StorageEntryJSON(operatorStoragePath, operator)
	if err != nil {
		return nil, err
	}

	// store in vault storage
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// reset the client so the next invocation will pick up the new configuration
	b.reset()

	return nil, nil
}

func getOperator(ctx context.Context, s logical.Storage) (*natsOperator, error) {
	entry, err := s.Get(ctx, operatorStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	operator := new(natsOperator)
	if err := entry.DecodeJSON(&operator); err != nil {
		return nil, fmt.Errorf("error reading root operatoruration: %w", err)
	}

	// return the operator, we are done
	return operator, nil
}

// pathOperatorHelpSynopsis summarizes the help text for the operatoruration
const pathOperatorHelpSyn = `Short Description.`

// pathOperatorHelpDescription describes the help text for the operatoruration
const pathOperatorHelpDesc = `

Detailed Descrition

`

const pathOperatorGenerateHelpSyn = `Short Description.`
const pathOperatorGenerateHelpDesc = `

Detailed Descrition

`
