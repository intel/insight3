package ghclient

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GHClient :
type GHClient struct {
	ClientV3 *github.Client
	ClientV4 *githubv4.Client
}

// Setup : setup github client for v3 and v4
func (cli *GHClient) Setup(ctx context.Context, authToken string) error {
	//setup v2 client for go-client apis

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli.ClientV3 = github.NewClient(tc)

	cli.ClientV4 = githubv4.NewClient(tc)

	return nil
}
