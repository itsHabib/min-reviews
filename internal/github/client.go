package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Client is used to interact with the GitHub API to get PRs and team members
type Client struct {
	c     *http.Client
	token string
}

// NewClient initiates a new GitHub client to interact with the API.
// Usage example:
//
//	c, err := github.NewClient(token)
//	handleErr(err)
//	if err != nil {  }
//	// get PR
//	pr, err := c.GetPR(repo, prNumber)
//	handleErr(err)
//	// get team members
//	members, err := c.GetTeamMembers(teamURL)
//	handleErr(err)
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("unable to initialize new client due to missing token")
	}

	return &Client{
		c:     new(http.Client),
		token: token,
	}, nil
}

// GetPR takes in a repo and pr number and returns a PullRequest object.
// Returns an error if failing to action the request or receiving
// a non-200 status code.
func (c *Client) GetPR(repo string, prNumber int) (*PullRequest, error) {
	url := repoURL(repo) + "/pulls/" + strconv.Itoa(prNumber)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get pr: %w", err)
	}
	defer cleanupResponse(resp)

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non 200 status code: %d, body: %s", resp.StatusCode, string(b))
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("unable to decode response body: %w", err)
	}

	return &pr, nil
}

// GetTeamMembers returns a list of users belonging to the team URL given.
// Returns an error if failing to action the request or if receiving a
// non-200 status code.
func (c *Client) GetTeamMembers(teamURL string) ([]User, error) {
	req, err := http.NewRequest(http.MethodGet, teamURL+"/members", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to form request to get team members: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to submit request to get team memvers: %w", err)
	}

	defer cleanupResponse(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return users, nil
}

func cleanupResponse(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_, _ = io.Copy(io.Discard, resp.Body)
}

func repoURL(repo string) string {
	return baseRepoAPIPath + "/" + repo
}
