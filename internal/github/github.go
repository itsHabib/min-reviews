package github

const baseRepoAPIPath = "https://api.github.com/repos"

type PullRequest struct {
	ID                 int                 `json:"id"`
	URL                string              `json:"url"`
	State              string              `json:"state"`
	Number             int                 `json:"number"`
	RequestedTeams     []RequestedTeam     `json:"requested_teams"`
	RequestedReviewers []RequestedReviewer `json:"requested_reviewers"`
}

type RequestedReviewer struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type RequestedTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type User struct {
	Login string `json:"login"`
}
