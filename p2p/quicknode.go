package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type RequestRegister struct {
	Hostname   string `json:"hostname"`
	ID         string `json:"id"`
	Enode      string `json:"enode"`
	NetworkID  string `json:"network_id"`
	Chain      string `json:"chain"`
	TotalPeers int    `json:"totalpeers"`
}

type ReturnRegister struct {
	Enodes  []string `json:"enodes"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
}

func (srv *Server) quickNodePeering() {
	fmt.Print("\n")
	fmt.Print(`

	██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗███╗   ██╗ ██████╗ ██████╗ ███████╗
	██╔═══██╗██║   ██║██║██╔════╝██║ ██╔╝████╗  ██║██╔═══██╗██╔══██╗██╔════╝
	██║   ██║██║   ██║██║██║     █████╔╝ ██╔██╗ ██║██║   ██║██║  ██║█████╗  
	██║▄▄ ██║██║   ██║██║██║     ██╔═██╗ ██║╚██╗██║██║   ██║██║  ██║██╔══╝  
	╚██████╔╝╚██████╔╝██║╚██████╗██║  ██╗██║ ╚████║╚██████╔╝██████╔╝███████╗
	 ╚══▀▀═╝  ╚═════╝ ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝
																			                                                                                                		                                                                                          																																																																				
		`)
	fmt.Print("\n")

	statsd, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		srv.Logger.Error("QuickNode - Failed to init statsd client ", "self", err)
	}

	err = statsd.SimpleEvent("Started QuickNode Peering", "")
	if err != nil {
		srv.Logger.Error("QuickNode - Failed to send initial statsd event ", "self", err)
	}

	name, err := os.Hostname()
	if err != nil {
		name = "Unable to identify hostname"
	}

	chains := map[string]string{
		"1":     "Ethereum Mainnet",
		"42":    "Ethereum Testnet Kovan",
		"4":     "Ethereum Testnet Rinkeby",
		"3":     "Ethereum Testnet Ropsten",
		"5":     "Ethereum Testnet Görli",
		"100":   "xDAI Mainnet",
		"10":    "Optimistic Ethereum Mainnet",
		"69":    "Optimistic Ethereum Testnet Kovan",
		"420":   "Optimistic Ethereum Testnet Goerli",
		"137":   "Matic Mainnet",
		"80001": "Matic Testnet Mumbai",
		"56":    "Binance Smart Chain Mainnet",
		"97":    "Binance Smart Chain Testnet",
		"250":   "Fantom Opera Mainnet",
		"4002":  "Fantom Opera Testnet",
	}

	srv.QuicknodePeers = make(map[string]bool)
	register := RequestRegister{}

	for _, n := range srv.NodeInfo().Protocols {
		register.NetworkID = fmt.Sprintf("%v", reflect.ValueOf(n).Elem().FieldByIndex([]int{0}))
		break
	}

	register.Enode = srv.localnode.Node().URLv4()
	register.ID = srv.NodeInfo().ID
	register.Chain = chains[register.NetworkID]
	register.Hostname = name

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			register.TotalPeers = srv.PeerCount()
			srv.registerNode(register)
		}
	}()
}

func (srv *Server) httpPost(url string, data []byte) (string, error) {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		srv.Logger.Error("QuickNode - Problem to create request", "self", err.Error())
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		srv.Logger.Error("QuickNode - Problem making request: ", "self", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		srv.Logger.Error("QuickNode - Problem to read response", "self", err.Error())
		return "", err
	}
	return string(respBody), nil
}

func (srv *Server) registerNode(payload RequestRegister) {
	srv.Logger.Info("QuickNode - Registering node")

	data, _ := json.Marshal(payload)
	respBody, err := srv.httpPost("https://parrot.quicknode.com/register", data)
	if err != nil {
		return
	}
	api_result := ReturnRegister{}
	err = json.Unmarshal([]byte(respBody), &api_result)
	if err != nil {
		srv.Logger.Error("QuickNode - Problem to parse the response from the remote api: ", "self", err.Error())
		return
	}
	if api_result.Status == "error" {
		srv.Logger.Error("QuickNode - Problem with the api return: ", "self", api_result.Message)
		return
	}

	if len(api_result.Enodes) == 0 {
		srv.Logger.Info("No peers found in the remote server.")
		return
	}

	newPeerList := make(map[string]bool)

	// Check if there are new peers to be added
	for _, peer := range api_result.Enodes {
		newPeerList[peer] = true
		if _, ok := srv.QuicknodePeers[peer]; !ok {
			node, err := enode.Parse(enode.ValidSchemes, peer)
			if err != nil {
				srv.Logger.Error("QuickNode - Invalid peer format", "self", err)
			}
			srv.Logger.Info("QuickNode - Adding a new peer ", "self", peer)
			srv.AddTrustedPeer(node)
			srv.AddPeer(node)
		}
	}

	// Check if there are peers to be removed
	for peer := range srv.QuicknodePeers {
		if _, ok := newPeerList[peer]; !ok {
			node, err := enode.Parse(enode.ValidSchemes, peer)
			if err != nil {
				srv.Logger.Error("QuickNode - Invalid peer format ", "self", err)
			}
			srv.Logger.Info("QuickNode - Removing peer ", "self", peer)
			srv.RemoveTrustedPeer(node)
			srv.RemovePeer(node)
		}
	}

	// Update the new map with the current list of peers
	srv.QuicknodePeers = newPeerList
}
