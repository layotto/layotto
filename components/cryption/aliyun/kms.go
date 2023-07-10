// Code generated by github.com/seeflood/protoc-gen-p6 .

// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyun

import (
	"context"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"mosn.io/layotto/components/cryption"
	"mosn.io/pkg/log"
)

type cy struct {
	client *kms20160120.Client
	keyId  string
}

/*
refer: https://help.aliyun.com/document_detail/611325.html
*/
func NewCryption() cryption.CryptionService {
	return &cy{}
}

func (k *cy) Init(ctx context.Context, conf *cryption.Config) error {
	accessKey := conf.Metadata[ClientKey]
	secret := conf.Metadata[ClientSecret]
	endpoint := conf.Metadata[EndPoint]
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(accessKey),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(secret),
		// Endpoint 请参考 https://api.aliyun.com/product/Kms
		Endpoint: tea.String(endpoint),
	}

	client, err := kms20160120.NewClient(config)
	if err != nil {
		return err
	}
	k.client = client
	k.keyId = conf.Metadata[KeyID]
	return nil
}

func (k *cy) Decrypt(ctx context.Context, request *cryption.DecryptRequest) (*cryption.DecryptResponse, error) {
	decryptRequest := &kms20160120.DecryptRequest{
		CiphertextBlob: tea.String(string(request.CipherText)),
	}
	decryptResp, err := k.client.Decrypt(decryptRequest)
	if err != nil {
		log.DefaultLogger.Errorf("failed decrypt, err: %+v", err)
		return nil, fmt.Errorf("fail decrypt with error: %+v", err)
	}
	resp := &cryption.DecryptResponse{KeyId: *decryptResp.Body.KeyId, KeyVersionId: *decryptResp.Body.KeyVersionId, PlainText: []byte(*decryptResp.Body.Plaintext)}
	return resp, nil
}

func (k *cy) Encrypt(ctx context.Context, request *cryption.EncryptRequest) (*cryption.EncryptResponse, error) {
	encryptRequest := &kms20160120.EncryptRequest{
		KeyId:     tea.String(k.keyId),
		Plaintext: tea.String(string(request.PlainText)),
	}

	encryptResp, err := k.client.Encrypt(encryptRequest)

	if err != nil {
		log.DefaultLogger.Errorf("fail encrypt, err: %+v", err)
		return nil, fmt.Errorf("fail encrypt with error: %+v", err)
	}
	resp := &cryption.EncryptResponse{CipherText: []byte(*encryptResp.Body.CiphertextBlob)}
	return resp, nil
}
