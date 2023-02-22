package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically

type Asset struct {
	Status                   string          `json:"Status"`
	ExternalRefNum           string          `json:"ExternalRefNum"`
	TicketID                 string          `json:"TicketID"`
	ISPID                    string          `json:"ISPID"`
	ActionTaken              string          `json:"ActionTaken"`
	AuthCode                 string          `json:"AuthCode"`
	PromisedDate             string          `json:"PromisedDate"`
	VisitStartDate           string          `json:"VisitStartDate"`
	VisitEndDate             string          `json:"VisitEndDate"`
	TicketDescription        string          `json:"ClaimDescription"`
	FaultKey                 FaultKey        `json:"FaultKey"`
	FaultDescription         string          `json:"FaultDescription"`
	FaultReported            string          `json:"FaultReported"`
	CoverageCode             string          `json:"CoverageCode"`
	ResolutionCode           string          `json:"ResolutionCode"`
	ProductDetail            ProductDetail   `json:"ProductDetail"`
	Consumer                 Consumer        `json:"Consumer"`
	Site                     ServiceLocation `json:"Site"`
	PartsConsumed            []Part          `json:"PartsConsumed"`
	ErrorCode                string          `json:"ErrorCode"`
	Warranty                 Warranty        `json:"Warranty"`
	ProofOfPurchaseIndicator string          `json:"ProofOfPurchaseIndicator"`
}
type FaultKey struct {
	FaultCode     string `json:"FaultCode"`
	ComponentCode string `json:"ComponentCode"`
	DefectCode    string `json:"DefectCode"`
}

type ProductDetail struct {
	ProductID    string `json:"ProductID"`
	Brand        string `json:"Brand"`
	ProductType  string `json:"ProductType"`
	ModelName    string `json:"ModelName"`
	SerialNumber string `json:"SerialNumber"`
	MLCode       string `json:"MLCode"`
	SkillKey     string `json:"SkillKey"`
	PurchaseDate string `json:"PurchaseDate"`
	Retailer     string `json:"Retailer"`
	WarrantyCode string `json:"WarrantyCode"`
}
type Consumer struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	Language   string `json:"Language"`
	Email      string `json:"Email"`
	Mobile     string `json:"Mobile"`
	Telephone  string `json:"Telephone"`
	Street     string `json:"Street"`
	City       string `json:"City"`
	Country    string `json:"Country"`
	PostalCode string `json:"PostalCode"`
	LegalID    string `json:"LegalID"`
}
type ServiceLocation struct {
	Street     string `json:"Street"`
	City       string `json:"City"`
	Country    string `json:"Country"`
	PostalCode string `json:"PostalCode"`
}

type Part struct {
	ProductCode string `json:"ProductCode"`
	Quantity    string `json:"Quantity"`
	OrderNumber string `json:"OrderNumber"`
	ClaimAmount string `json:"ClaimAmount"`
}

type Warranty struct {
	WarrantyNumber string `json:"WarrantyNumber"`
	ProviderID     string `json:"ProviderID"`
	ProviderName   string `json:"ProviderName"`
	StartDate      string `json:"StartDate"`
	EndDate        string `json:"EndDate"`
	TicketNumber   string `json:"TicketNumber"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{{
		TicketID:                 "655035",
		Status:                   "Dispatched",
		ExternalRefNum:           "1070935",
		ISPID:                    "ITXISP1412",
		PromisedDate:             "20230110",
		VisitStartDate:           "20230110",
		VisitEndDate:             "20230110",
		TicketDescription:        "AVATRICE ELECTROLUX TUTTO SPENTO",
		CoverageCode:             "Z2",
		ProofOfPurchaseIndicator: "X",
		Consumer: Consumer{
			FirstName:  "Amalia",
			LastName:   "Gennarelli",
			Language:   "IT",
			Email:      "Amalia.Gennarelli@gmail.com",
			Mobile:     "069863536",
			Telephone:  "069863536",
			Street:     "Via Calipso 00042",
			City:       "Anzio",
			Country:    "IT",
			PostalCode: "00042",
		},

		FaultKey: FaultKey{
			FaultCode:     "G10",
			ComponentCode: "121",
			DefectCode:    "22",
		},
		Warranty: Warranty{
			WarrantyNumber: "1239040",
			ProviderID:     "11",
			ProviderName:   "Nexure",
			StartDate:      "20201007",
			EndDate:        "20221007",
		},
		ProductDetail: ProductDetail{
			ProductID:    "914916405",
			Brand:        "ELX",
			ProductType:  "Washer,Front Loaded",
			ModelName:    "EW6F4922EB",
			SerialNumber: "03700396",
			MLCode:       "00",
			PurchaseDate: "20201007",
			Retailer:     "Grando",
			WarrantyCode: "2 LW",
		},
		PartsConsumed: []Part{
			Part{
				ProductCode: "1802394933",
				Quantity:    "1",
				OrderNumber: "1230002",
				ClaimAmount: "100",
			},
			Part{
				ProductCode: "1802394067",
				Quantity:    "1",
				OrderNumber: "320039920",
				ClaimAmount: "120",
			},
		},
	},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.TicketID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, TicketID string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", TicketID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, TicketID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(TicketID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) UpdateTicketStatus(ctx contractapi.TransactionContextInterface, ticketId string, newStatus string) (string, error) {
	asset, err := s.ReadAsset(ctx, ticketId)
	if err != nil {
		return "", err
	}

	oldStatus := asset.Status
	asset.Status = newStatus

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(ticketId, assetJSON)
	if err != nil {
		return "", err
	}

	return oldStatus, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) AddParts(ctx contractapi.TransactionContextInterface, ticketId string, args string) error {
	asset, err := s.ReadAsset(ctx, ticketId)
	if err != nil {
		return err
	}
	parts := make([]Part, 0)
	json.Unmarshal([]byte(args), &parts)

	for _, part := range parts {
		asset.PartsConsumed = append(asset.PartsConsumed, part)
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(ticketId, assetJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) SubmitClaim(ctx contractapi.TransactionContextInterface,
	ticketId string,
	externalRefNum string,
	actionTaken string,
	visitStartDate string,
	visitEndDate string,
	faultCode string,
	componentCode string,
	defectCode string,
	serialNumber string,
	mLCode string,
	purchaseDate string) (string, error) {

	asset, err := s.ReadAsset(ctx, ticketId)
	if err != nil {
		return "", err
	}

	if externalRefNum != "" {
		asset.ExternalRefNum = externalRefNum
	}
	if actionTaken != "" {
		asset.ActionTaken = actionTaken
	}
	if visitStartDate != "" {
		asset.VisitStartDate = visitStartDate
	}
	if visitEndDate != "" {
		asset.VisitEndDate = visitEndDate
	}
	if faultCode != "" {
		asset.FaultKey.FaultCode = faultCode
	}
	if componentCode != "" {
		asset.FaultKey.ComponentCode = componentCode
	}
	if defectCode != "" {
		asset.FaultKey.DefectCode = defectCode
	}
	if serialNumber != "" {
		asset.ProductDetail.SerialNumber = serialNumber
	}
	if mLCode != "" {
		asset.ProductDetail.MLCode = mLCode
	}
	if purchaseDate != "" {
		asset.ProductDetail.PurchaseDate = purchaseDate
	}

	asset.Status = "ClaimSubmitted"

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(ticketId, assetJSON)
	if err != nil {
		return "", err
	}

	return asset.Status, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface,
	ticketID, ispId, ticketDescription, promisedDate, coverageCode,
	firstName, lastName, language, email, mobile, street, city, country, postalCode,
	warrantyNumber, providerID, providerName, warrantyStartDate, warrantyEndDate,
	productID, purchaseDate, retailer, warrantyCode, mLCode, serialNumber string) error {
	exists, err := s.AssetExists(ctx, ticketID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ticketID)
	}
	asset := Asset{
		TicketID:          ticketID,
		ISPID:             ispId,
		TicketDescription: ticketDescription,
		PromisedDate:      promisedDate,
		CoverageCode:      coverageCode,
		Consumer: Consumer{
			FirstName:  firstName,
			LastName:   lastName,
			Language:   language,
			Email:      email,
			Mobile:     mobile,
			Telephone:  mobile,
			Street:     street,
			City:       city,
			Country:    country,
			PostalCode: postalCode,
		},
		Warranty: Warranty{
			WarrantyNumber: warrantyNumber,
			ProviderID:     providerID,
			ProviderName:   providerName,
			StartDate:      warrantyStartDate,
			EndDate:        warrantyEndDate,
		},
		ProductDetail: ProductDetail{
			ProductID:    productID,
			SerialNumber: serialNumber,
			MLCode:       mLCode,
			PurchaseDate: purchaseDate,
			Retailer:     retailer,
			WarrantyCode: warrantyCode,
		},
		PartsConsumed: []Part{
			Part{
				ProductCode: "30016",
				Quantity:    "1",
			},
		},
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ticketID, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
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

		assets = append(assets, &asset)
	}

	return assets, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
