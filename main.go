package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/schollz/progressbar/v3"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func main() {
	var (
		RedditAppID     = os.Getenv("REDDIT_APP_ID")
		RedditAppSecret = os.Getenv("REDDIT_APP_SECRET")
		RedditUsername  = os.Getenv("REDDIT_USERNAME")
		RedditPassword  = os.Getenv("REDDIT_PASSWORD")
	)

	credentials := reddit.Credentials{
		ID:       RedditAppID,
		Secret:   RedditAppSecret,
		Username: RedditUsername,
		Password: RedditPassword,
	}
	client, _ := reddit.NewClient(credentials)

	ctx := context.Background()

	for {
		var posts []*reddit.Post
		var comments []*reddit.Comment
		var hasMore bool
		err := doRequest(func() (res *reddit.Response, err error) {
			posts, comments, res, err = client.User.Saved(ctx, &reddit.ListUserOverviewOptions{
				ListOptions: reddit.ListOptions{
					Limit: 100,
				},
			})
			if res != nil && res.After != "" {
				hasMore = true
			}
			return
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Got %v posts\n", len(posts))
		fmt.Printf("Got %v comments\n", len(comments))

		// Loop through posts and unsave each one
		if len(posts) > 0 {
			fmt.Println("Unsaving received posts")
			bar := progressbar.Default(int64(len(posts)))
			for _, post := range posts {
				if err := doRequest(func() (*reddit.Response, error) {
					return client.Post.Unsave(ctx, post.FullID)
				}); err != nil {
					panic(err)
				}
				_ = bar.Add(1)
			}
		}

		// Loop through comments and unsave each one
		if len(comments) > 0 {
			fmt.Println("Unsaving received comments")
			bar := progressbar.Default(int64(len(comments)))
			for _, comment := range comments {
				if err := doRequest(func() (*reddit.Response, error) {
					return client.Comment.Unsave(ctx, comment.FullID)
				}); err != nil {
					panic(err)
				}
				_ = bar.Add(1)
			}
		}

		if !hasMore {
			break
		}
	}

	fmt.Println("Everything unsaved")
}

func doRequest(req func() (*reddit.Response, error)) error {
	retries := 1
	for {
		res, err := req()
		if err != nil {
			if rerr := (&reddit.RateLimitError{}); errors.As(err, &rerr) {
				handleRateLimit(rerr.Rate)
				retries++
			} else {
				return fmt.Errorf("doing request: %w", err)
			}
		} else {
			handleRateLimit(res.Rate)
			return nil
		}

		if retries > 5 {
			return fmt.Errorf("too many retries")
		}
	}
}

func handleRateLimit(rate reddit.Rate) {
	if !rate.Reset.IsZero() && rate.Remaining == 0 && time.Now().Before(rate.Reset) {
		fmt.Println()
		fmt.Printf("Rate limit exceeded, waiting until %v\n", rate.Reset)
		time.Sleep(time.Until(rate.Reset))
	}
}
