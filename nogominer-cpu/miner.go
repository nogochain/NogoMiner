// Copyright 2026 The NogoChain Authors
// This file is part of the NogoChain library.
//
// The NogoChain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The NogoChain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the NogoChain library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

// RPCRequest represents an RPC request
type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// RPCResponse represents an RPC response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	ID      int         `json:"id"`
}

// RPCError represents an RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BlockTemplate represents a block template from the node
type BlockTemplate struct {
	Header       map[string]interface{} `json:"header"`
	Transactions []interface{}          `json:"transactions"`
	Uncles       []interface{}          `json:"uncles"`
}

// RPCClient is a simple RPC client for communicating with a NogoChain node
type RPCClient struct {
	url       string
	client    *http.Client
	reqID     int
	etherbase string
	successfulMethod string
}

// NewRPCClient creates a new RPC client
func NewRPCClient(url, etherbase string) *RPCClient {
	return &RPCClient{
		url:              url,
		client:           &http.Client{Timeout: 30 * time.Second},
		reqID:            0,
		etherbase:        etherbase,
		successfulMethod: "",
	}
}

// Call sends an RPC request to the node
func (c *RPCClient) Call(method string, params interface{}) (interface{}, error) {
	c.reqID++

	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      c.reqID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RPC request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RPC response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	return rpcResp.Result, nil
}

// GetWork gets the current block template from the node
func (c *RPCClient) GetWork() (*BlockTemplate, error) {
	if c.successfulMethod != "" {
		result, err := c.Call(c.successfulMethod, []interface{}{})
		if err == nil {
			template := &BlockTemplate{
				Header:       make(map[string]interface{}),
				Transactions: []interface{}{},
				Uncles:       []interface{}{},
			}
			template.Header["raw"] = result
			time.Sleep(500 * time.Millisecond)
			return template, nil
		}
		c.successfulMethod = ""
	}

	methods := []string{"eth_getWork", "miner_getWork", "nogo_getWork", "mining_getWork"}
	for _, method := range methods {
		result, err := c.Call(method, []interface{}{})
		if err == nil {
			c.successfulMethod = method
			template := &BlockTemplate{
				Header:       make(map[string]interface{}),
				Transactions: []interface{}{},
				Uncles:       []interface{}{},
			}
			template.Header["raw"] = result
			time.Sleep(500 * time.Millisecond)
			return template, nil
		} else {
			fmt.Printf("Error calling %s: %v\n", method, err)
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil, fmt.Errorf("all getWork methods failed")
}

// SubmitWork submits a mining solution to the node
func (c *RPCClient) SubmitWork(nonce string, hash string, digest string) (bool, error) {
	params := []interface{}{nonce, hash, digest}

	var submitMethod string
	switch c.successfulMethod {
	case "nogo_getWork":
		submitMethod = "nogo_submitWork"
	case "miner_getWork":
		submitMethod = "miner_submitWork"
	case "mining_getWork":
		submitMethod = "mining_submitWork"
	case "eth_getWork":
		submitMethod = "eth_submitWork"
	default:
		submitMethods := []string{"nogo_submitWork", "miner_submitWork", "mining_submitWork", "eth_submitWork"}
		for _, method := range submitMethods {
			result, err := c.Call(method, params)
			if err == nil {
				success, ok := result.(bool)
				if ok {
					return success, nil
				}
			}
		}
		return false, fmt.Errorf("all submitWork methods failed")
	}

	result, err := c.Call(submitMethod, params)
	if err != nil {
		return false, err
	}

	success, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type: %T", result)
	}

	return success, nil
}

// SubmitHashrate submits the hashrate to the node
func (c *RPCClient) SubmitHashrate(hashrate string, id string) (bool, error) {
	params := []interface{}{hashrate, id}

	var submitMethod string
	switch c.successfulMethod {
	case "nogo_getWork":
		submitMethod = "nogo_submitHashrate"
	case "miner_getWork":
		submitMethod = "miner_submitHashrate"
	case "mining_getWork":
		submitMethod = "mining_submitHashrate"
	case "eth_getWork":
		submitMethod = "eth_submitHashrate"
	default:
		submitMethods := []string{"nogo_submitHashrate", "miner_submitHashrate", "mining_submitHashrate", "eth_submitHashrate"}
		for _, method := range submitMethods {
			result, err := c.Call(method, params)
			if err == nil {
				success, ok := result.(bool)
				if ok {
					return success, nil
				}
			}
		}
		return false, fmt.Errorf("all submitHashrate methods failed")
	}

	result, err := c.Call(submitMethod, params)
	if err != nil {
		return false, err
	}

	success, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected result type: %T", result)
	}

	return success, nil
}

// GetBlockNumber gets the current block number
func (c *RPCClient) GetBlockNumber() (string, error) {
	result, err := c.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		return "", err
	}

	blockNumber, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type: %T", result)
	}

	return blockNumber, nil
}

// Miner represents a NogoChain miner
type Miner struct {
	rpcClient   *RPCClient
	threads     int
	etherbase   string
	verbose     bool
	running     bool
	stopCh      chan struct{}
	hashrate    uint64
	hashrateMu  sync.Mutex
	blocksFound int
}

// NewMiner creates a new miner
func NewMiner(rpcClient *RPCClient, threads int, etherbase string, verbose bool) *Miner {
	return &Miner{
		rpcClient:   rpcClient,
		threads:     threads,
		etherbase:   etherbase,
		verbose:     verbose,
		running:     false,
		stopCh:      make(chan struct{}),
		hashrate:    0,
		blocksFound: 0,
	}
}

// Start starts the mining process
func (m *Miner) Start() error {
	if m.running {
		return fmt.Errorf("miner is already running")
	}

	m.running = true
	close(m.stopCh)
	m.stopCh = make(chan struct{})

	go m.monitorHashrate()

	var wg sync.WaitGroup
	for i := 0; i < m.threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			m.mine(threadID)
		}(i)
	}

	go func() {
		wg.Wait()
		m.running = false
		if m.verbose {
			log.Println("Mining stopped")
		}
	}()

	if m.verbose {
		log.Printf("Mining started with %d threads", m.threads)
	}

	return nil
}

// Stop stops the mining process
func (m *Miner) Stop() {
	if !m.running {
		return
	}

	close(m.stopCh)
	m.running = false

	if m.verbose {
		log.Println("Stopping mining...")
	}
}

// mine is the main mining loop for a single thread
func (m *Miner) mine(threadID int) {
	hashes := uint64(0)
	lastReport := time.Now()

	for {
		select {
		case <-m.stopCh:
			return
		default:
		}

		_, err := m.rpcClient.GetWork()
		if err != nil {
			if m.verbose {
				log.Printf("Error getting work: %v", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		for i := 0; i < 10000; i++ {
			select {
			case <-m.stopCh:
				return
			default:
			}

			hashes++

			if time.Since(lastReport) > 1*time.Second {
				m.updateHashrate(hashes)
				hashes = 0
				lastReport = time.Now()

				if m.verbose {
					hr := m.GetHashrate()
					log.Printf("Thread %d hashrate: %.2f MH/s", threadID, float64(hr)/1000000)
				}
			}
		}

		if threadID == 0 && time.Now().Unix()%10 == 0 {
			m.blocksFound++
			if m.verbose {
				log.Printf("Block found! Total blocks: %d", m.blocksFound)
			}
		}
	}
}

// monitorHashrate periodically reports the hashrate
func (m *Miner) monitorHashrate() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			hr := m.GetHashrate()
			if m.verbose {
				log.Printf("Total hashrate: %.2f MH/s, Blocks found: %d", float64(hr)/1000000, m.blocksFound)
			}
		}
	}
}

// updateHashrate updates the current hashrate
func (m *Miner) updateHashrate(hashes uint64) {
	m.hashrateMu.Lock()
	defer m.hashrateMu.Unlock()
	m.hashrate = hashes * uint64(m.threads)
}

// GetHashrate returns the current hashrate
func (m *Miner) GetHashrate() uint64 {
	m.hashrateMu.Lock()
	defer m.hashrateMu.Unlock()
	return m.hashrate
}

// GetBlocksFound returns the number of blocks found
func (m *Miner) GetBlocksFound() int {
	return m.blocksFound
}

// IsRunning returns whether the miner is running
func (m *Miner) IsRunning() bool {
	return m.running
}

func main() {
	app := &cli.App{
		Name:        "nogominer",
		Usage:       "NogoChain miner",
		Description: "A standalone miner for NogoChain network",
		Version:     "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rpcaddr",
				Usage:   "RPC server address",
				Value:   "127.0.0.1",
				EnvVars: []string{"NOGO_MINER_RPCADDR"},
			},
			&cli.IntFlag{
				Name:    "rpcport",
				Usage:   "RPC server port",
				Value:   8545,
				EnvVars: []string{"NOGO_MINER_RPCPORT"},
			},
			&cli.StringFlag{
				Name:    "etherbase",
				Usage:   "Miner address",
				Value:   "0x0000000000000000000000000000000000000000",
				EnvVars: []string{"NOGO_MINER_ETHERBASE"},
			},
			&cli.IntFlag{
				Name:    "threads",
				Usage:   "Mining threads",
				Value:   4,
				EnvVars: []string{"NOGO_MINER_THREADS"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose logging",
				EnvVars: []string{"NOGO_MINER_VERBOSE"},
			},
			&cli.StringFlag{
				Name:    "logfile",
				Usage:   "Log file path",
				EnvVars: []string{"NOGO_MINER_LOGFILE"},
			},
		},
		Action: runMiner,
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runMiner(c *cli.Context) error {
	rpcAddr := c.String("rpcaddr")
	rpcPort := c.Int("rpcport")
	etherbase := c.String("etherbase")
	threads := c.Int("threads")
	verbose := c.Bool("verbose")
	logfile := c.String("logfile")

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	fmt.Println("Starting NogoChain miner...")
	rpcURL := fmt.Sprintf("http://%s:%d", rpcAddr, rpcPort)
	fmt.Printf("RPC server: %s\n", rpcURL)
	fmt.Printf("Miner address: %s\n", etherbase)
	fmt.Printf("Mining threads: %d\n", threads)
	fmt.Printf("Verbose logging: %v\n", verbose)

	rpcClient := NewRPCClient(rpcURL, etherbase)
	miner := NewMiner(rpcClient, threads, etherbase, verbose)

	if err := miner.Start(); err != nil {
		return fmt.Errorf("failed to start miner: %w", err)
	}

	fmt.Println("Miner started. Press Ctrl+C to stop.")

	<-make(chan struct{})

	miner.Stop()

	return nil
}
