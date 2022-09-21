package models

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Andosius/valorant-exporter/cfg"
	"github.com/Andosius/valorant-exporter/helpers"
)

type Client struct {
	Lockfile        Lockfile
	LocalAuthHeader string
	HttpClient      http.Client
	Entitlements    Entitlements
	ConfigManager   ConfigManager
}

type Entitlements struct {
	AccessToken string `json:"accessToken"`
	Subject     string `json:"subject"`
	Token       string `json:"token"`
}

type Lockfile struct {
	Name     string
	PID      string
	Port     string
	Password string
	Protocol string
}

func NewClient() *Client {
	c := Client{}
	/*
		TODO:
			- Check if riot games client is running if ReadLockFile produces any
			  type of error (rewrite to produce / return error)
			- Move Content to seperate function like c.Setup()
	*/

	c.ReadLockfile()
	c.SetupHttpClient()
	c.SetupLocalHeader()
	c.SetupEntitlements()
	return &c
}

func (c *Client) ReadLockfile() {
	// Get %LocalAppData% Path and store it in path
	path, err := os.UserCacheDir()

	// Check for errors; true? -> print debug string
	helpers.Fatal("Client.ReadLockfile:1", err)

	// Add static information to path
	path = path + `\Riot Games\Riot Client\Config\lockfile`

	// Read config file
	cfg, err := os.ReadFile(path)

	// Check for errors; true? -> print debug string
	if err != nil {
		err = errors.New("riot games client is not running or you are not logged in")
	}
	helpers.Fatal("Client.ReadLockfile:2", err)

	// Split string by ":"-char and provide all data to Client.Lockfile-Element
	data := strings.Split(string(cfg), ":")
	c.Lockfile = Lockfile{
		Name:     data[0],
		PID:      data[1],
		Port:     data[2],
		Password: data[3],
		Protocol: data[4],
	}
}

func (c *Client) SetupHttpClient() {
	// Get all elements in current x509-certpool; otherwise create new one
	certs, _ := x509.SystemCertPool()

	if certs == nil {
		certs = x509.NewCertPool()
	}

	// Read in official Riot Games cert to make sure API-Calls don't get
	// rejected for such a dumb reason
	cert, err := os.ReadFile(`certs\riotgames.pem`)

	helpers.Fatal("Client.SetupHttpClient:1", err)

	// Try to add the certificate into the pool
	if ok := certs.AppendCertsFromPEM(cert); !ok {
		helpers.Fatal("Client.SetupHttpClient:2", errors.New("unable to append riotgames.pem to certpool"))
	}

	config := &tls.Config{
		// Skip unnecessary checks
		InsecureSkipVerify: false,
		// Add our certs to the tls conf
		RootCAs: certs,
	}

	// Define Transport so our requests don't end up delaying the whole application
	transport := &http.Transport{
		TLSClientConfig:     config,
		TLSHandshakeTimeout: 10 * time.Second,
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
	}

	// Create the http client and provide these settings
	client := http.Client{
		Transport: transport,
	}
	c.HttpClient = client
}

func (c *Client) SetupLocalHeader() {
	// Get a byte-slice of "riot:password" as login credentials
	raw := []byte("riot:" + c.Lockfile.Password)

	// Base64-Encode it and set the header as Basic Auth
	token := base64.StdEncoding.EncodeToString(raw)

	c.LocalAuthHeader = "Basic " + token
}

func (c *Client) SetupEntitlements() {
	url := c.Lockfile.Protocol + "://127.0.0.1:" + c.Lockfile.Port + "/entitlements/v1/token"
	req := helpers.CreateAPIRequest("GET", url, c.GetLocalHeader(), "")

	resp, err := c.HttpClient.Do(req)
	helpers.Fatal("c.SetupEntitlements:1", err)
	defer resp.Body.Close()

	// Read entitlements data
	body, err := io.ReadAll(resp.Body)
	helpers.Fatal("c.SetupEntitlements:2", err)

	// Unmarshal data :=)
	var ent Entitlements
	err = json.Unmarshal(body, &ent)
	helpers.Fatal("c.SetupEntitlements:3", err)

	c.Entitlements = ent
}

func (c Client) GetLocalHeader() map[string]string {
	return map[string]string{
		"Authorization": c.LocalAuthHeader,
	}
}

func (c Client) GetGlobalHeader() map[string]string {

	// Declare valorant-api.com data
	type VersionData struct {
		RiotClientVersion string `json:"riotClientVersion"`
	}

	type APIResult struct {
		Data VersionData `json:"data"`
	}

	// Prepare request and execute it
	req := helpers.CreateAPIRequest("GET", "https://valorant-api.com/v1/version", map[string]string{}, "")

	resp, err := c.HttpClient.Do(req)
	helpers.Fatal("c.GetGlobalHeader:1", err)

	defer resp.Body.Close()

	// Read received data and parse into struct
	body, err := io.ReadAll(resp.Body)
	helpers.Fatal("c.GetGlobalHeader:2", err)

	var api APIResult
	err = json.Unmarshal(body, &api)
	helpers.Fatal("c.GetGlobalHeader:3", err)

	return map[string]string{
		"Authorization":           "Bearer " + c.Entitlements.AccessToken,
		"X-Riot-Entitlements-JWT": c.Entitlements.Token,
		"X-Riot-ClientPlatform":   cfg.CLIENT_PLATFORM,
		"X-Riot-ClientVersion":    api.Data.RiotClientVersion,
	}
}

func (c Client) PushConfigToServer(idx int) {
	// Create structures the server accepts and format data
	type SendConfig struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}

	var sc SendConfig = SendConfig{
		Type: "Ares.PlayerSettings",
		Data: c.ConfigManager.Configs[idx].Data,
	}

	// Marshal struct to JSON
	body, err := json.Marshal(sc)
	helpers.Fatal("c.PushConfigToServer:1", err)

	// Create request and send it to server
	req := helpers.CreateAPIRequest("PUT", cfg.PUT_PLAYER_PREFERENCES, c.GetGlobalHeader(), string(body))

	resp, err := c.HttpClient.Do(req)
	helpers.Fatal("c.PushConfigToServer:1", err)

	// Send success messages
	fmt.Println("Status:", resp.Status)
	resp.Body.Close()
}

func (c *Client) GetCurrentAccountSettings() Config {
	req := helpers.CreateAPIRequest("GET", cfg.GET_PLAYER_PREFERENCES, c.GetGlobalHeader(), "")

	resp, err := c.HttpClient.Do(req)
	helpers.Fatal("c.GetCurrentAccountSettings:1", err)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	helpers.Fatal("c.GetCurrentAccountSettings:2", err)

	var cfg Config
	err = json.Unmarshal(body, &cfg)

	helpers.Fatal("c.GetCurrentAccountSettings:3", err)

	return cfg
}
