package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DataProvenanceChaincode define a estrutura do chaincode
type DataProvenanceChaincode struct {
    contractapi.Contract
}

// DataRecord representa um registro de dados armazenado na blockchain
type DataRecord struct {
    ID      string `json:"id"`
    Content string `json:"content"`
}

// InitLedger inicializa o ledger com dados básicos
func (d *DataProvenanceChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
    data := []DataRecord{
        {ID: "1", Content: "Texto de exemplo 1"},
        {ID: "2", Content: "Texto de exemplo 2"},
    }

    for _, record := range data {
        recordJSON, err := json.Marshal(record)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(record.ID, recordJSON)
        if err != nil {
            return fmt.Errorf("falha ao inserir dado no ledger: %v", err)
        }
    }

    return nil
}

// AddData adiciona um novo registro de dados na blockchain
func (d *DataProvenanceChaincode) AddData(ctx contractapi.TransactionContextInterface, id string, content string) error {
    exists, err := d.DataExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("o registro %s já existe", id)
    }

    record := DataRecord{
        ID:      id,
        Content: content,
    }
    recordJSON, err := json.Marshal(record)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, recordJSON)
}

// ReadData retorna o conteúdo de um registro de dados da blockchain
func (d *DataProvenanceChaincode) ReadData(ctx contractapi.TransactionContextInterface, id string) (*DataRecord, error) {
    recordJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("falha ao ler do ledger: %v", err)
    }
    if recordJSON == nil {
        return nil, fmt.Errorf("o registro %s não existe", id)
    }

    var record DataRecord
    err = json.Unmarshal(recordJSON, &record)
    if err != nil {
        return nil, err
    }

    return &record, nil
}

// DataExists verifica se um registro de dados existe na blockchain
func (d *DataProvenanceChaincode) DataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    recordJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("falha ao ler do ledger: %v", err)
    }

    return recordJSON != nil, nil
}

// UpdateData atualiza o conteúdo de um registro de dados existente na blockchain
func (d *DataProvenanceChaincode) UpdateData(ctx contractapi.TransactionContextInterface, id string, newContent string) error {
    exists, err := d.DataExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("o registro %s não existe", id)
    }

    record := DataRecord{
        ID:      id,
        Content: newContent,
    }
    recordJSON, err := json.Marshal(record)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, recordJSON)
}

// DeleteData remove um registro de dados da blockchain
func (d *DataProvenanceChaincode) DeleteData(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := d.DataExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("o registro %s não existe", id)
    }

    return ctx.GetStub().DelState(id)
}

func main() {
    chaincode, err := contractapi.NewChaincode(&DataProvenanceChaincode{})
    if err != nil {
        log.Panicf("Erro ao criar o chaincode: %v", err)
    }

    if err := chaincode.Start(); err != nil {
        log.Panicf("Erro ao iniciar o chaincode: %v", err)
    }
}
