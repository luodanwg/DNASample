package httpjsonrpc

import (
	"DNA/client"
	. "DNA/common"
	"DNA/common/log"
	. "DNA/core/asset"
	"DNA/core/contract"
	"DNA/core/signature"
	"DNA/core/transaction"
	"strconv"
)

const (
	ASSETPREFIX = "DNA"
)

func NewRegTx(rand string, index int, admin, issuer *client.Account) *transaction.Transaction {
	name := ASSETPREFIX + "-" + strconv.Itoa(index) + "-" + rand
	asset := &Asset{name, byte(0x00), AssetType(Share), UTXO}
	amount := Fixed64(1000)
	controller, _ := contract.CreateSignatureContract(admin.PubKey())
	tx, _ := transaction.NewRegisterAssetTransaction(asset, amount, issuer.PubKey(), controller.ProgramHash)
	return tx
}

func SignTx(admin *client.Account, tx *transaction.Transaction) {
	signdate, err := signature.SignBySigner(tx, admin)
	if err != nil {
		log.Error(err, "signdate SignBySigner failed")
	}
	transactionContract, _ := contract.CreateSignatureContract(admin.PublicKey)
	transactionContractContext := contract.NewContractContext(tx)
	transactionContractContext.AddContract(transactionContract, admin.PublicKey, signdate)
	tx.SetPrograms(transactionContractContext.GetPrograms())
}

func NewIssueTx(admin *client.Account, hash Uint256) *transaction.Transaction {
	//UTXO　Output generate.
	signatureRedeemScript, err := contract.CreateSignatureRedeemScript(admin.PublicKey)
	if err != nil {
		log.Error("EncodePoint error.")
	}
	hashx, err := ToCodeHash(signatureRedeemScript)
	if err != nil {
		log.Error("TocodeHash hash error.")
	}
	issueTxOutput := &transaction.TxOutput{
		AssetID:     hash,
		Value:       Fixed64(100),
		ProgramHash: hashx,
	}
	outputs := []*transaction.TxOutput{}
	outputs = append(outputs, issueTxOutput)

	// generate transaction
	tx, _ := transaction.NewIssueAssetTransaction(outputs)
	return tx
}

func NewTransferTx(regHash, issueHash Uint256, toUser *client.Account) *transaction.Transaction {
	//UTXO  Input Generate.
	transferUTXOInput := &transaction.UTXOTxInput{
		ReferTxID:          issueHash,
		ReferTxOutputIndex: uint16(0),
	}
	inputs := []*transaction.UTXOTxInput{}
	inputs = append(inputs, transferUTXOInput)

	//UTXO　Output generate.
	signatureRedeemScript, err := contract.CreateSignatureRedeemScript(toUser.PublicKey)
	if err != nil {
		log.Error("EncodePoint error.")
	}
	hashx, err := ToCodeHash(signatureRedeemScript)
	if err != nil {
		log.Error("TocodeHash hash error.")
	}
	transferTxOutput := &transaction.TxOutput{
		AssetID:     regHash,
		Value:       Fixed64(100),
		ProgramHash: hashx,
	}
	outputs := []*transaction.TxOutput{}
	outputs = append(outputs, transferTxOutput)

	// generate transaction
	tx, _ := transaction.NewTransferAssetTransaction(inputs, outputs)
	return tx
}

func NewRecordTx(rand string) *transaction.Transaction {
	recordType := string("txt")
	recordData := []byte("hello world " + rand)
	tx, _ := transaction.NewRecordTransaction(recordType, recordData)
	return tx
}

