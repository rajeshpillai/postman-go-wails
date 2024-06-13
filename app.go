package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	loggedIn  bool
	username  string
	loggedInM sync.Mutex
}

type FetchRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type FetchResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) Login(username, password string) error {
	correctUsername := "admin"
	correctPassword := "admin"

	fmt.Println("Attempting login with username:", username)

	if username == correctUsername && password == correctPassword {
		a.loggedInM.Lock()
		a.loggedIn = true
		a.username = username
		a.loggedInM.Unlock()
		fmt.Println("Login successful")
		return nil
	} else {
		fmt.Println("Login failed")
		return fmt.Errorf("incorrect username or password")
	}
}

func (a *App) Logout() error {
	a.loggedInM.Lock()
	a.loggedIn = false
	a.username = ""
	a.loggedInM.Unlock()
	fmt.Println("Logged out")
	return nil
}

func (a *App) CheckLogin() (string, error) {
	a.loggedInM.Lock()
	defer a.loggedInM.Unlock()
	if a.loggedIn {
		fmt.Println("User is logged in as:", a.username)
		return a.username, nil
	}
	fmt.Println("No user is logged in")
	return "", nil
}

func (a *App) PerformFetch(request FetchRequest) (FetchResponse, error) {
	fmt.Println("Received fetch request:")
	fmt.Println("URL:", request.URL)
	fmt.Println("Method:", request.Method)
	fmt.Println("Headers:", request.Headers)
	fmt.Println("Body:", request.Body)

	client := &http.Client{}
	req, err := http.NewRequest(request.Method, request.URL, strings.NewReader(request.Body))
	if err != nil {
		return FetchResponse{}, err
	}

	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return FetchResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FetchResponse{}, err
	}

	headers := make(map[string]string)
	for key, values := range resp.Header {
		headers[key] = strings.Join(values, ", ")
	}

	fmt.Println("Response status:", resp.StatusCode)
	fmt.Println("Response headers:", headers)
	fmt.Println("Response body:", string(body))

	return FetchResponse{
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(body),
	}, nil
}

func (a *App) startup(ctx context.Context) {
	runtime.EventsOn(ctx, "login", func(data ...interface{}) {
		if len(data) == 2 {
			username, _ := data[0].(string)
			password, _ := data[1].(string)
			if err := a.Login(username, password); err != nil {
				runtime.LogError(ctx, fmt.Sprintf("Login failed: %s", err))
			}
		}
	})

	runtime.EventsOn(ctx, "logout", func(data ...interface{}) {
		if err := a.Logout(); err != nil {
			runtime.LogError(ctx, fmt.Sprintf("Logout failed: %s", err))
		}
	})

	runtime.EventsOn(ctx, "check_login", func(data ...interface{}) {
		username, err := a.CheckLogin()
		if err != nil {
			runtime.LogError(ctx, fmt.Sprintf("Check login failed: %s", err))
		}
		runtime.EventsEmit(ctx, "login_status", username)
	})

	runtime.EventsOn(ctx, "perform_fetch", func(data ...interface{}) {
		if len(data) == 1 {
			request := data[0].(map[string]interface{})
			fetchRequest := FetchRequest{
				URL:    request["url"].(string),
				Method: request["method"].(string),
				Headers: func() map[string]string {
					headers := map[string]string{}
					for key, value := range request["headers"].(map[string]interface{}) {
						headers[key] = value.(string)
					}
					return headers
				}(),
				Body: request["body"].(string),
			}
			response, err := a.PerformFetch(fetchRequest)
			if err != nil {
				runtime.LogError(ctx, fmt.Sprintf("Fetch failed: %s", err))
			}
			runtime.EventsEmit(ctx, "fetch_response", response)
		}
	})
}
