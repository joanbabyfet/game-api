package operator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"game-api/pkg"
)

func (c *Client) post(ctx context.Context, path string, req any, resp any) error {
	
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		path,
		bytes.NewReader(body),
	)
	if err != nil {
		fmt.Println("xxxxxx")
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	body, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Code != pkg.SUCCESS {
		return fmt.Errorf("operator error: %s", result.Msg)
	}

	if resp != nil {
		if err := json.Unmarshal(result.Data, resp); err != nil {
			return err
		}
	}

	return nil
}