package cgminer

import (
	"fmt"
	"bufio"
	"net"
	"strings"
	"encoding/json"
	"errors"
)

type CGMiner struct {
	server string
}

type status struct {
	Code        int
	Description string
	Status      string `json:"STATUS"`
	When        int64
}

type Summary struct {
	Accepted               string
	BestShare              int64   `json:"Best Share"`
	DeviceHardwarePercent  float64 `json:"Device Hardware%"`
	DeviceRejectedPercent  float64 `json:"Device Rejected%"`
	DifficultyAccepted     float64 `json:"Difficulty Accepted"`
	DifficultyRejected     float64 `json:"Difficulty Rejected"`
	DifficultyStale        float64 `json:"Difficulty Stale"`
	Discarded              int64
	Elapsed                int64
	FoundBlocks            int64 `json:"Found Blocks"`
	GetFailures            int64 `json:"Get Failures"`
	Getworks               int64
	HardwareErrors         int64   `json:"Hardware Errors"`
	LocalWork              int64   `json:"Local Work"`
	NetworkBlocks          int64   `json:"Network Blocks"`
	PoolRejectedPercentage float64 `json:"Pool Rejected%"`
	PoolStalePercentage    float64 `json:"Pool Stale%"`
	Rejected               string
	RemoteFailures         int64 `json:"Remote Failures"`
	Stale                  int64
	TotalMH                float64 `json:"Total MH"`
	Utilty                 float64
	WorkUtility            float64 `json:"Work Utility"`
	HashrateAvr		float64 `json:"GHS av"`
	
}

type summaryResponse struct {
	Status  []status  `json:"STATUS"`
	Summary []Summary `json:"SUMMARY"`
	Id      int64     `json:"id"`
}



func New(hostname string, port int64) *CGMiner {
	miner := new(CGMiner)
	server := fmt.Sprintf("%s:%d", hostname, port)
	miner.server = server

	return miner
}

func (miner *CGMiner) runCommand(command, argument string) (string, error) {
	conn, err := net.Dial("tcp", miner.server)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	type commandRequest struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	request := &commandRequest{
		Command: command,
	}

	if argument != "" {
		request.Parameter = argument
	}

	requestBody, err := json.Marshal(request)

	if err != nil {
		return "", err
	}

	fmt.Fprintf(conn, "%s", requestBody)

	result, err := bufio.NewReader(conn).ReadString('\x00')

	if err != nil {
		return "", err
	}

	return strings.TrimRight(result, "\x00"), nil
}

func (miner *CGMiner) Summary() (*Summary, error) {
	result, err := miner.runCommand("summary", "")
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>"+result)

	var summaryResponse summaryResponse

	err = json.Unmarshal([]byte(result), &summaryResponse)
	if err != nil {
		return nil, err
	}

	if len(summaryResponse.Summary) != 1 {
		return nil, errors.New("Received multiple Summary objects")
	}

	var summary = summaryResponse.Summary[0]
	return &summary, err
}