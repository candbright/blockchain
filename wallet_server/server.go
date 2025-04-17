package main

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
	}
}

func (s *WalletServer) Port() uint16 {
	return s.port
}

func (s *WalletServer) Gateway() string {
	return s.gateway
}

func (s *WalletServer) Run() {

}
