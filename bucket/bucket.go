package bucket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tiredkangaroo/ajiteshcc/env"
)

type Object struct {
	Name      string `json:"name"`       // object key
	PublicURL string `json:"public_url"` // public URL of the object
}

type listAllObjectsInBucketResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Result  []struct {
		Key string `json:"key"`
	} `json:"result"`
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
	var objects []Object
	pub_url := *env.DefaultEnv.R2_PHOTOS_BUCKET_PUBLIC_URL
	for _, obj := range listResp.Result {
		pub_url.Path = "/" + obj.Key // set the path to the object key
		objects = append(objects, Object{
			Name:      obj.Key,
			PublicURL: pub_url.String(),
		})
	}
	return objects, nil
}
