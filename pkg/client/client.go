package client

import (
	"encoding/gob"
	"errors"
	"net"
)


type Client struct {
		addresses []string
}

func NewClient(addresses []string) *Client {
		return &Client{addresses: addresses}
}

func (cli *Client) sendRequest(opearation, key, value string) (map[string]interface{}, error) {
	// Implement simple retry/failover logic
	for _, address := range cli.addresses {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}
		defer conn.Close()

		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		
		request := struct {
			Operation	string
			Key			string
			Value		string
		}{operation, key, value}


		if err := encoder.Encode(request); err != nil {
			continue
		}

		var response map[string]interface{}
		if err := decoder.Decode(&response); err != nil {
			continue
		}


		return response, nil
	}

	return nil, errors.New("failed to connect to any node")
}

func (cli *Client) Get(key string) (string, bool, error) {
		response, err := cli.sendRequest("GET", key, "")
		if err != nil {
			return "", false, err
		}
		return response["value"].(string),response["ok"].(bool), nil
}

func (cli *Client) Set(key, value string) error {
		_, err := cli.sendRequest("SET", key, value)
		return err
}

func (cli *Client) Delete(key string) error {
		_, err := cli.sendRequest("DELETE", key, "")
		return err
}

