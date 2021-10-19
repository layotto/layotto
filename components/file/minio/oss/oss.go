/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oss

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/minio/minio-go/v6"

	"mosn.io/layotto/components/file"
)

const (
	listPrefix  = "listPrefix"
	endpointKey = "endpoint"
	bucketKey   = "bucket"
	fileSize    = "fileSize"
)

var (
	ErrMissingBucket      error = errors.New("missing bucket info in metadata")
	ErrMissingEndPoint    error = errors.New("missing endpoint info in metadata")
	ErrClientNotExist     error = errors.New("specific client not exist")
	ErrInvalidConfig      error = errors.New("invalid minio oss config")
	ErrNotSpecifyEndPoint error = errors.New("not specify endpoint in metadata")
)

type MinioOss struct {
	client map[string]*minio.Client
	meta   map[string]*MinioMetaData
}

type MinioMetaData struct {
	Region          string `json:"region"`
	EndPoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	SSL             bool   `json:"SSL"`
}

func NewMinioOss() file.File {
	return &MinioOss{
		client: make(map[string]*minio.Client),
		meta:   make(map[string]*MinioMetaData),
	}
}

func (m *MinioOss) Init(config *file.FileConfig) error {
	md := make([]*MinioMetaData, 0)
	err := json.Unmarshal(config.Metadata, &md)
	if err != nil {
		return ErrInvalidConfig
	}
	for _, data := range md {
		if !data.isMinioMetaValid() {
			return ErrInvalidConfig
		}
		client, err := m.createOssClient(data)
		if err != nil {
			continue
		}
		m.client[data.EndPoint] = client
		m.meta[data.EndPoint] = data
	}
	return nil
}

func (m *MinioOss) Put(st *file.PutFileStu) error {
	var (
		bucket string
		ok     bool
		size   int64 = -1
	)
	if bucket, ok = st.Metadata[bucketKey]; !ok {
		return ErrMissingBucket
	}
	client, err := m.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	// specify file size from metadata, default unknown size is -1
	if info, ok := st.Metadata[fileSize]; ok {
		size, err = strconv.ParseInt(info, 10, 64)
		if err != nil {
			return err
		}
	}
	_, err = client.PutObject(bucket, st.FileName, st.DataStream, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioOss) Get(st *file.GetFileStu) (io.ReadCloser, error) {
	var (
		bucket string
		ok     bool
	)
	if bucket, ok = st.Metadata[bucketKey]; !ok {
		return nil, ErrMissingBucket
	}
	client, err := m.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	obj, err := client.GetObject(bucket, st.FileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *MinioOss) List(st *file.ListRequest) (*file.ListResp, error) {
	var (
		bucket string
		prefix string
		ok     bool
		resp   = &file.ListResp{}
	)
	if bucket, ok = st.Metadata[bucketKey]; !ok {
		return nil, ErrMissingBucket
	}
	if p, ok1 := st.Metadata[listPrefix]; ok1 {
		prefix = p
	}
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)
	isRecursive := true
	client, err := m.selectClient(st.Metadata)
	if err != nil {
		return nil, err
	}
	objectCh := client.ListObjects(bucket, prefix, isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			continue
		}
		resp.FilesName = append(resp.FilesName, object.Key)
	}
	return resp, nil
}

func (m *MinioOss) Del(st *file.DelRequest) error {
	var (
		bucket string
		ok     bool
	)
	if bucket, ok = st.Metadata[bucketKey]; !ok {
		return ErrMissingBucket
	}
	client, err := m.selectClient(st.Metadata)
	if err != nil {
		return err
	}
	return client.RemoveObject(bucket, st.FileName)
}

func (m *MinioOss) createOssClient(meta *MinioMetaData) (*minio.Client, error) {
	if meta.Region == "" {
		return minio.New(meta.EndPoint, meta.AccessKeyID, meta.AccessKeySecret, meta.SSL)
	}
	return minio.NewWithRegion(meta.EndPoint, meta.AccessKeyID, meta.AccessKeySecret, meta.SSL, meta.Region)
}

func (m *MinioOss) selectClient(meta map[string]string) (client *minio.Client, err error) {
	var (
		endpoint string
		ok       bool
	)
	if endpoint, ok = meta[endpointKey]; !ok {
		if len(m.client) == 1 {
			for _, client := range m.client {
				return client, nil
			}
		}
		return nil, ErrNotSpecifyEndPoint
	}
	if client, ok = m.client[endpoint]; !ok {
		err = ErrClientNotExist
		return
	}
	return
}

// isMinioMetaValid check if the metadata is valid
func (mm *MinioMetaData) isMinioMetaValid() bool {
	if mm.AccessKeySecret == "" || mm.EndPoint == "" || mm.AccessKeyID == "" {
		return false
	}
	return true
}
