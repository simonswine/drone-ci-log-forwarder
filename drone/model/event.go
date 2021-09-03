package model

type Event struct {
	ID                     int    `json:"id"`
	UID                    string `json:"uid"`
	UserID                 int    `json:"user_id"`
	Namespace              string `json:"namespace"`
	Name                   string `json:"name"`
	Slug                   string `json:"slug"`
	Scm                    string `json:"scm"`
	GitHTTPURL             string `json:"git_http_url"`
	GitSSHURL              string `json:"git_ssh_url"`
	Link                   string `json:"link"`
	DefaultBranch          string `json:"default_branch"`
	Private                bool   `json:"private"`
	Visibility             string `json:"visibility"`
	Active                 bool   `json:"active"`
	ConfigPath             string `json:"config_path"`
	Trusted                bool   `json:"trusted"`
	Protected              bool   `json:"protected"`
	IgnoreForks            bool   `json:"ignore_forks"`
	IgnorePullRequests     bool   `json:"ignore_pull_requests"`
	AutoCancelPullRequests bool   `json:"auto_cancel_pull_requests"`
	AutoCancelPushes       bool   `json:"auto_cancel_pushes"`
	Timeout                int    `json:"timeout"`
	Counter                int    `json:"counter"`
	Synced                 int    `json:"synced"`
	Created                int    `json:"created"`
	Updated                int    `json:"updated"`
	Version                int    `json:"version"`
	Build                  Build  `json:"build"`
}
type Step struct {
	ID        int      `json:"id"`
	StepID    int      `json:"step_id"`
	Number    int      `json:"number"`
	Name      string   `json:"name"`
	Status    string   `json:"status"`
	ExitCode  int      `json:"exit_code"`
	Started   int      `json:"started,omitempty"`
	Stopped   int      `json:"stopped,omitempty"`
	Version   int      `json:"version"`
	Image     string   `json:"image"`
	DependsOn []string `json:"depends_on,omitempty"`
}
type Stage struct {
	ID        int    `json:"id"`
	RepoID    int    `json:"repo_id"`
	BuildID   int    `json:"build_id"`
	Number    int    `json:"number"`
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Errignore bool   `json:"errignore"`
	ExitCode  int    `json:"exit_code"`
	Machine   string `json:"machine"`
	Os        string `json:"os"`
	Arch      string `json:"arch"`
	Started   int    `json:"started"`
	Stopped   int    `json:"stopped"`
	Created   int    `json:"created"`
	Updated   int    `json:"updated"`
	Version   int    `json:"version"`
	OnSuccess bool   `json:"on_success"`
	OnFailure bool   `json:"on_failure"`
	Steps     []Step `json:"steps"`
}
type Build struct {
	ID           int     `json:"id"`
	RepoID       int     `json:"repo_id"`
	Trigger      string  `json:"trigger"`
	Number       int     `json:"number"`
	Status       string  `json:"status"`
	Event        string  `json:"event"`
	Action       string  `json:"action"`
	Link         string  `json:"link"`
	Timestamp    int     `json:"timestamp"`
	Message      string  `json:"message"`
	Before       string  `json:"before"`
	After        string  `json:"after"`
	Ref          string  `json:"ref"`
	SourceRepo   string  `json:"source_repo"`
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	AuthorLogin  string  `json:"author_login"`
	AuthorName   string  `json:"author_name"`
	AuthorEmail  string  `json:"author_email"`
	AuthorAvatar string  `json:"author_avatar"`
	Sender       string  `json:"sender"`
	Cron         string  `json:"cron"`
	Started      int     `json:"started"`
	Finished     int     `json:"finished"`
	Created      int     `json:"created"`
	Updated      int     `json:"updated"`
	Version      int     `json:"version"`
	Stages       []Stage `json:"stages"`
}
