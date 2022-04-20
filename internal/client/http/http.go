package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/matvoy/cparser/internal/models"
)

type httpClient struct {
	client *http.Client
	url    string
}

func NewClient(url string) (*httpClient, error) {
	if len(url) == 0 {
		return nil, nil
	}
	return &httpClient{&http.Client{}, url}, nil
}

func (c *httpClient) GetBlockData(height uint64) (*models.Block, error) {
	if height == 0 {
		return nil, nil
	}
	var result *BlockResponse
	var txResult *TxResponse
	var wg sync.WaitGroup
	var errorResult error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		result, err = c.getBlockByHeight(height)
		if err != nil {
			errorResult = multierror.Append(errorResult, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		// get txs data from api
		// or use mapToDomainWithTxDecoding
		txResult, err = c.getTxsByHeight(height)
		if err != nil {
			errorResult = multierror.Append(errorResult, err)
		}
	}()

	wg.Wait()
	return mapToDomainWithTxResponse(result, txResult), errorResult
}

func (c *httpClient) getBlockByHeight(height uint64) (*BlockResponse, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/blocks/%v", c.url, height))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result BlockResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *httpClient) getTxsByHeight(height uint64) (*TxResponse, error) {
	// tried to get with message.action param but it doesn't work
	// I can continue research if it's necessary
	resp, err := c.client.Get(fmt.Sprintf("%s/txs?tx.height=%v", c.url, height))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result TxResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
