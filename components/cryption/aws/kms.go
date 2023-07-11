package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"mosn.io/layotto/components/cryption"
	"mosn.io/pkg/log"
)

type cy struct {
	client *kms.KMS
	keyID  string
}

func NewCryption() cryption.CryptionService {
	return &cy{}
}

func (k *cy) Init(ctx context.Context, conf *cryption.Config) error {
	accessKey := conf.Metadata[cryption.ClientKey]
	secret := conf.Metadata[cryption.ClientSecret]
	region := conf.Metadata[cryption.Region]
	keyID := conf.Metadata[cryption.KeyID]
	staticCredentials := credentials.NewStaticCredentials(accessKey, secret, "")

	awsConf := &aws.Config{
		Region:      aws.String(region),
		Credentials: staticCredentials,
	}
	client := kms.New(session.New(), awsConf)
	k.client = client
	k.keyID = keyID
	return nil
}

func (k *cy) Decrypt(ctx context.Context, request *cryption.DecryptRequest) (*cryption.DecryptResponse, error) {
	decryptRequest := &kms.DecryptInput{
		CiphertextBlob: request.CipherText,
	}
	decryptResp, err := k.client.Decrypt(decryptRequest)
	if err != nil {
		log.DefaultLogger.Errorf("fail decrypt data, err: %+v", err)
		return nil, fmt.Errorf("fail decrypt data with error: %+v", err)
	}
	resp := &cryption.DecryptResponse{KeyId: *decryptResp.KeyId, PlainText: decryptResp.Plaintext}
	return resp, nil
}

func (k *cy) Encrypt(ctx context.Context, request *cryption.EncryptRequest) (*cryption.EncryptResponse, error) {
	encryptRequest := &kms.EncryptInput{
		KeyId:     aws.String(k.keyID),
		Plaintext: request.PlainText,
	}

	encryptResp, err := k.client.Encrypt(encryptRequest)
	if err != nil {
		log.DefaultLogger.Errorf("fail encrypt data, err: %+v", err)
		return nil, fmt.Errorf("fail encrypt data with error: %+v", err)
	}
	resp := &cryption.EncryptResponse{KeyId: *encryptResp.KeyId, CipherText: encryptResp.CiphertextBlob}
	return resp, nil
}
