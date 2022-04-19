# Reddit Unsaver

Unsaves all saved posts and comments.

## Usage

Create a reddit script app [here](https://www.reddit.com/prefs/apps). You can
set the redirect url to anything you want.

Create a `.env` file at the root of the project and add the app id, app secret,
reddit username and reddit password as shown in the example below.

```env
REDDIT_APP_ID="your app id"
REDDIT_APP_SECRET="your app secret"
REDDIT_USERNAME="your username"
REDDIT_PASSWORD="your password"
```

Final step is to just run `go run .`