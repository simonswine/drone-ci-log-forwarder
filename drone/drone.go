package drone

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/drone/drone-go/drone"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/simonswine/drone-ci-log-forwarder/drone/model"
	"golang.org/x/oauth2"
)

type Drone struct {
	client drone.Client
	logger log.Logger
	host   string

	httpClient *http.Client
}

func New(host string, token string) *Drone {
	// create an http client with oauth authentication.
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	// create the drone client with authenticator

	return &Drone{
		httpClient: auther,
		client:     drone.NewClient(host, auther),
		logger:     log.NewNopLogger(),
		host:       host,
	}
}

func (d *Drone) WithLogger(logger log.Logger) *Drone {
	d.logger = logger
	return d
}

func (d *Drone) eventWatcher(eventCh chan *model.Event) error {
	logger := log.WithPrefix(d.logger, "component", "event-watcher")
	url := strings.TrimRight(d.host, "/") + "/api/stream"

	req, err := http.NewRequest("GET", url, nil) //ioutil.NopCloser(pr))
	if err != nil {
		return err
	}
	req.Header.Set("accept", "text/event-stream")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected status code %d (expected 200)", resp.StatusCode)
	}
	level.Debug(logger).Log("msg", "connected to event stream", "url", url, "proto", resp.Proto)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {

		text := scanner.Text()

		if text == "" {
			continue
		}

		if text == ": ping" {
			level.Debug(logger).Log("msg", "received ping")
			continue
		}

		dataPrefix := "data:"
		if len(text) > len(dataPrefix) && text[0:len(dataPrefix)] == dataPrefix {
			var event model.Event
			if err := json.Unmarshal(scanner.Bytes()[len(dataPrefix):], &event); err != nil {
				level.Warn(logger).Log("msg", "failed parsing json data", "data", text, "error", err)
				continue
			}
			eventCh <- &event
			continue
		}

		level.Debug(logger).Log("msg", "unknown event", "text", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (d *Drone) Run() error {
	// check if I am able to login
	user, err := d.client.Self()
	if err != nil {
		return fmt.Errorf("unable to authenticate with drone: %w", err)
	}
	level.Info(d.logger).Log("msg", "successfully authenticated with drone server", "server", d.host, "user", user.Login)

	var (
		eventCh = make(chan *model.Event)
		wg      sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := d.eventWatcher(eventCh); err != nil && !strings.Contains(err.Error(), "http2: server sent GOAWAY") {
				level.Error(d.logger).Log("msg", "event watcher failed", "error", err)
				break
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range eventCh {
			if e.Namespace != "grafana" {
				continue
			}
			if e.Name != "backend-enterprise" && e.Name != "gex-plugins" {
				continue
			}
			level.Info(d.logger).Log("msg", "event received", "repo", e.Name, "namespace", e.Namespace)
		}
	}()

	wg.Wait()

	return nil
}

func stepKeysFromEvent(e *model.Event) []stepKey {
	base := stepKey{
		Namespace:   e.Namespace,
		Name:        e.Name,
		BuildNumber: e.Build.Number,
	}

	var stepKeys []stepKey

	for _, stage := range e.Build.Stages {
		for _, step := range stage.Steps {
			sk := base
			sk.Stage = stage.Name
			sk.Step = step.Name
			stepKeys = append(stepKeys, sk)
		}
	}

	return stepKeys
}

type stepKey struct {
	Namespace   string
	Name        string
	BuildNumber int
	Stage       string
	Step        string
}
