package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/traces"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := traces.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := traces.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	interval, err := client.GetSyncedInterval(ctx)
	if err != nil {
		fmt.Println("failed to GetSyncedInterval: %w", err)
		return
	}

	fmt.Println("GetSyncedInterval:", interval)
	time.Sleep(time.Second)

	blockTrace, err := client.GetBlockTraceByNumber(ctx, traces.GetBlockTraceByNumberParam(17378176))
	if err != nil {
		fmt.Println("failed to GetBlockTraceByNumber: %w", err)
		return
	}

	fmt.Println("GetBlockTraceByNumber:", blockTrace)
	time.Sleep(time.Second)

	txTrace, err := client.GetTxTraceByNumberAndHash(ctx, traces.GetTxTraceByNumberAndHashParams{
		BlockNumber:     17378177,
		TransactionHash: "0x16897e492b2e023d8f07be9e925f2c15a91000ef11a01fc71e70f75050f1e03c",
	})
	if err != nil {
		fmt.Println("failed to GetTxTraceByNumberAndHash: %w", err)
		return
	}

	fmt.Println("GetTxTraceByNumberAndHash:", txTrace)
	time.Sleep(time.Second)

	txTraceOffset, err := client.GetTxTraceByNumberAndOffset(ctx, traces.GetTxTraceByNumberAndOffsetParams{
		BlockNumber: 17378177,
		Offset:      1,
	})
	if err != nil {
		fmt.Println("failed to GetTxTraceByNumberAndOffset: %w", err)
		return
	}

	fmt.Println("GetTxTraceByNumberAndOffset:", txTraceOffset)
	time.Sleep(time.Second)
}
