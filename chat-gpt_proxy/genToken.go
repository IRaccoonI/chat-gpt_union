package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Account struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}
type Proxy struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (p Proxy) Socks5URL() string {
	// Returns proxy URL (socks5)
	if p.User == "" && p.Pass == "" {
		return fmt.Sprintf("socks5://%s:%s", p.IP, p.Port)
	}
	return fmt.Sprintf("socks5://%s:%s@%s:%s", p.User, p.Pass, p.IP, p.Port)
}

// Read accounts.txt and create a list of accounts
func readAccounts() []Account {
	return []Account{
		{
			Email:    os.Getenv("PROXY_GPT_EMAIL"),
			Password: os.Getenv("PROXY_GPT_PASSWORD"),
		},
	}
}

// Read proxies from proxies.txt and create a list of proxies
func readProxies() []Proxy {
	return []Proxy{}
}

func GenToken() string {
	// Read accounts and proxies
	accounts := readAccounts()
	proxies := readProxies()

	// Loop through each account
	for _, account := range accounts {
		if os.Getenv("CF_PROXY") != "" {
			// exec warp-cli disconnect and connect
			exec.Command("warp-cli", "disconnect").Run()
			exec.Command("warp-cli", "connect").Run()
			time.Sleep(5 * time.Second)
		}
		var proxy_url string
		if len(proxies) == 0 {
			if os.Getenv("http_proxy") != "" {
				proxy_url = os.Getenv("http_proxy")
			}
		} else {
			proxy_url = proxies[0].Socks5URL()
			// Push used proxy to the back of the list
		}
		authenticator := NewAuthenticator(account.Email, account.Password, proxy_url)
		err := authenticator.Begin()
		if err != nil {
			// println("Error: " + err.Details)
			println("Location: " + err.Location)
			println("Status code: " + fmt.Sprint(err.StatusCode))
			println("Details: " + err.Details)
			println("Embedded error: " + err.Error.Error())
			return ""
		}
		access_token := authenticator.GetAccessToken()
		if access_token != "" {
			return access_token
		}
		// Write authenticated account to authenticated_accounts.txt
		// f, go_err = os.OpenFile("authenticated_accounts.txt", os.O_APPEND|os.O_WRONLY, 0600)
		// if go_err != nil {
		// 	continue
		// }
		// defer f.Close()
		// if _, go_err = f.WriteString(account.Email + ":" + account.Password + "\n"); go_err != nil {
		// 	continue
		// }
		// // Remove accounts.txt
		// os.Remove("accounts.txt")
		// // Create accounts.txt
		// f, go_err = os.Create("accounts.txt")
		// if go_err != nil {
		// 	continue
		// }
		// defer f.Close()
		// // Remove account from accounts
		// accounts = accounts[1:]
		// // Write unauthenticated accounts to accounts.txt
		// for _, acc := range accounts {
		// 	// Check if account is authenticated
		// 	if acc.Email == account.Email {
		// 		continue
		// 	}
		// 	if _, go_err = f.WriteString(acc.Email + ":" + acc.Password + "\n"); go_err != nil {
		// 		continue
		// 	}
		// }

	}
	return ""
}
