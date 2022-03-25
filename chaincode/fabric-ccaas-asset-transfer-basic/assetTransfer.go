/*
SPDX-License-hostnameentifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type serverConfig struct {
	CCID    string
	Address string
}

// SmartContract provhostnamees functions for managing an asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	Hostname	string `json:"hostname"`
	IP			string `json:"ip"`
	TTL			string `json:"ttl"`
	Time		string `json:"time"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Asset
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{Hostname: "khanghy.net", IP: "1.1.1.1", TTL: "300", Time: "500"},
		{Hostname: "khanghy.org", IP: "2.2.2.2", TTL: "300", Time: "500"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.Hostname, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, hostname, ip, ttl, time string) error {
	asset := Asset{
		Hostname:	hostname,
		IP:			ip,
		TTL:		ttl,
		Time:		time,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hostname, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given hostname.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, hostname string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", hostname)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provhostnameed parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, hostname, ip, ttl, time string) error {
	exists, err := s.AssetExists(ctx, hostname)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", hostname)
	}

	// overwritting original asset with new asset
	asset := Asset{
		Hostname:	hostname,
		IP:			ip,
		TTL:		ttl,
		Time:		time,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hostname, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, hostname string) error {
	exists, err := s.AssetExists(ctx, hostname)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", hostname)
	}

	return ctx.GetStub().DelState(hostname)
}

// AssetExists returns true when asset with given hostname exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, hostname string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(hostname)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the hostname field of asset with given hostname in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, hostname string, newHostname string) error {
	asset, err := s.ReadAsset(ctx, hostname)
	if err != nil {
		return err
	}

	asset.Hostname = newHostname
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hostname, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &asset}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {
	// See chaincode.env.example
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(&SmartContract{})

	if err != nil {
		log.Panicf("error create fabdns chaincode: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: getTLSProperties(),
	}

	if err := server.Start(); err != nil {
		log.Panicf("error starting fabdns chaincode: %s", err)
	}
}

func getTLSProperties() shim.TLSProperties {
	// Check if chaincode is TLS enabled
	tlsDisabledStr := getEnvOrDefault("CHAINCODE_TLS_DISABLED", "true")
	key := getEnvOrDefault("CHAINCODE_TLS_KEY", "")
	cert := getEnvOrDefault("CHAINCODE_TLS_CERT", "")
	clientCACert := getEnvOrDefault("CHAINCODE_CLIENT_CA_CERT", "")

	// convert tlsDisabledStr to boolean
	tlsDisabled := getBoolOrDefault(tlsDisabledStr, false)
	var keyBytes, certBytes, clientCACertBytes []byte
	var err error

	if !tlsDisabled {
		keyBytes, err = ioutil.ReadFile(key)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		certBytes, err = ioutil.ReadFile(cert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}
	// Dhostname not request for the peer cert verification
	if clientCACert != "" {
		clientCACertBytes, err = ioutil.ReadFile(clientCACert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	return shim.TLSProperties{
		Disabled: tlsDisabled,
		Key: keyBytes,
		Cert: certBytes,
		ClientCACerts: clientCACertBytes,
	}
}

func getEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		value = defaultVal
	}
	return value
}

// Note that the method returns default value if the string
// cannot be parsed!
func getBoolOrDefault(value string, defaultVal bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err!= nil {
		return defaultVal
	}
	return parsed
}
