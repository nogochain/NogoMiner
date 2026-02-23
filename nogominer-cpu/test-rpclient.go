package main

func TestRPCClient() {
	client := NewRPCClient("http://localhost:8545", "0x0000000000000000000000000000000000000000")
	_ = client
}
