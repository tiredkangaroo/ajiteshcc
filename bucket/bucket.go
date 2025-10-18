package bucket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tiredkangaroo/ajiteshcc/env"
)

type Object struct {
	Key string `json:"key"`
	// Etag           string            `json:"etag"`
	// LastModified   time.Time         `json:"last_modified"`
	// Size           int64             `json:"size"`
	// HTTPMetadata   map[string]string `json:"http_metadata"`
	// CustomMetadata any               `json:"custom_metadata"` // only got a response of {}
	// StorageClass string `json:"storage_class"`
}

type listAllObjectsInBucketResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Result  []Object `json:"result"`
}

func ListAllObjectsInBucket(bucketName string) ([]Object, error) {
	// https://api.cloudflare.com/client/v4/accounts/{account_id}/r2/buckets/{bucket_name}/objects
	u, err := url.Parse(fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/r2/buckets/%s/objects", env.DefaultEnv.R2_ACCOUNT_ID, bucketName))
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{
			"Authorization": []string{"Bearer " + env.DefaultEnv.R2_API_TOKEN},
		},
	}
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listResp listAllObjectsInBucketResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, err
	}
	if !listResp.Success {
		return nil, fmt.Errorf("list objects in bucket: %v", listResp.Errors)
	}
	return listResp.Result, nil
}
