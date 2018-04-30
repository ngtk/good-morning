# good-morning
`good-morning` is the tiny script to create a GitHub discussion for team morning meeting.


## Install
brew install ngtk/tools/good-morning

## Usage
You have to specify configs as env.

* `GOOD_MORNING_GITHUB_ACCESS_TOKEN`: GitHub personal access token
* `GOOD_MORNING_GITHUB_TEAM_SLUG`: GitHub team Slug which we can get from the URL like `/teams/:slug`
* `GOOD_MORNING_GITHUB_ORG_NAME`: GitHub org name which we can get from the URL like `/org/:name`

```
$ good-morning
```
