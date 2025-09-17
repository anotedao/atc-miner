package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func callMine() error {
	var networkByte = byte(55)
	var nodeURL = AnoteNodeURL

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rec, err := proto.NewRecipientFromString(DappAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	call := proto.FunctionCall{
		Name:      "mine",
		Arguments: proto.Arguments{},
	}

	// payments := proto.ScriptPayments{}
	// payments.Append(proto.ScriptPayment{
	// 	Amount: abi.Balance - RewardFee,
	// })

	fa := proto.OptionalAsset{}

	// Current time in milliseconds
	ts := uint64(time.Now().Unix() * 1000)

	tr := proto.NewUnsignedInvokeScriptWithProofs(
		2,
		networkByte,
		sender,
		rec,
		call,
		nil,
		fa,
		RewardFee,
		ts)

	tr.Proofs = proto.NewProofs()

	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		return err
	}

	tr.Sign(55, sk)

	// // Send the transaction to the network
	resp, err := cl.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getPublicKey(address string) string {
	pk := ""

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a, err := proto.NewAddressFromString(address)
	if err != nil {
		log.Println(err)
	}

	transactions, resp, err := client.Transactions.Address(ctx, a, 100)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	for _, tr := range transactions {
		at := AnoteTransaction{}
		trb, err := json.Marshal(tr)
		json.Unmarshal(trb, &at)
		pk, err := crypto.NewPublicKeyFromBase58(at.SenderPublicKey)
		if err != nil {
			log.Println(err)
		}
		addr, err := proto.NewAddressFromPublicKey(55, pk)
		if err != nil {
			log.Println(err)
		}
		if addr.String() == address {
			return at.SenderPublicKey
		}
	}

	transactions = nil

	return pk
}

type AnoteTransaction struct {
	SenderPublicKey string `json:"senderPublicKey"`
}

func getHeight() uint64 {
	height := uint64(0)

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bh, _, err := cl.Blocks.Height(ctx)
	if err != nil {
		log.Println(err)
	} else {
		height = bh.Height
	}

	return height
}
