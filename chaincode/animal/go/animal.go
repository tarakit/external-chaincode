/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for listing animals
type SmartContract struct {
	contractapi.Contract
}

// Animal describes basic details of the animal
type Animal struct {
	Origin string `json:"origin"`
	Name   string `json:"name"`
	Colour string `json:"colour"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Animal
}

// InitLedger adds a base set of animals to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	animals := []Animal{
		Animal{Origin: "Africa", Name: "African Elefant", Colour: "grey"},
		Animal{Origin: "Europe", Name: "Cow", Colour: "brown"},
		Animal{Origin: "Asia", Name: "Asian Elefant", Colour: "grey"},
	}

	for i, animal := range animals {
		animalAsByte, _ := json.Marshal(animal)
		err := ctx.GetStub().PutState("ANIMAL"+strconv.Itoa(i), animalAsByte)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateAnimal adds a new animal to the world state with given details
func (s *SmartContract) CreateAnimal(ctx contractapi.TransactionContextInterface, animalNumber string, origin string, name string, colour string) error {
	animal := Animal{
		Origin: origin,
		Name:   name,
		Colour: colour,
	}

	animalAsByte, _ := json.Marshal(animal)

	return ctx.GetStub().PutState(animalNumber, animalAsByte)
}

// QueryAnimal returns the animal stored in the world state with given id
func (s *SmartContract) QueryAnimal(ctx contractapi.TransactionContextInterface, animalNumber string) (*Animal, error) {
	animalAsByte, err := ctx.GetStub().GetState(animalNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if animalAsByte == nil {
		return nil, fmt.Errorf("%s does not exist", animalNumber)
	}

	animal := new(Animal)
	_ = json.Unmarshal(animalAsByte, animal)

	return animal, nil
}

// QueryAllAnimals returns all animals found in world state
func (s *SmartContract) QueryAllAnimals(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "ANIMAL0"
	endKey := "ANIMAL99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		animal := new(Animal)
		_ = json.Unmarshal(queryResponse.Value, animal)

		queryResult := QueryResult{Key: queryResponse.Key, Record: animal}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create aninal chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting animal chaincode: %s", err.Error())
	}
}
