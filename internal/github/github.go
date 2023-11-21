package github

const baseRepoAPIPath = "https://api.github.com/repos"

// PullRequest represents the data returned back from the GitHub API. It is
// intentionally truncated to only expose the data we need.
type PullRequest struct {
	ID                 int                 `json:"id"`
	URL                string              `json:"url"`
	State              string              `json:"state"`
	Number             int                 `json:"number"`
	RequestedTeams     []RequestedTeam     `json:"requested_teams"`
	RequestedReviewers []RequestedReviewer `json:"requested_reviewers"`
}

// RequestedReviewer represents the user that was requested to review a PR
type RequestedReviewer struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// RequestedTeam represents the team that was requested/assigned to review a PR
type RequestedTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// User represents a truncated version of a user object from the GitHub API
type User struct {
	Login string `json:"login"`
}
