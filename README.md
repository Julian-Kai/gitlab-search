# gitlab-search

It's inspired by [gitlab-search](https://github.com/phillipj/gitlab-search).

## Overview
This is a CLI tool to perform a thorough search of your own GitLab projects.
If you cannot use the advanced search, a premium feature, I hope it's helpful for you.

## Prerequisites
1. Install [Go](https://go.dev/doc/install)
2. Prepare your [GitLab Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token) with at least `read_api` scope

## Usage
```bash
$ make build
$ ./bin/gs search --url https://gitlab.com --token your_pat --keyword your_keyword

$ ./bin/gs search --help
Perform a thorough search of your GitLab projects

Usage:
  gs search [flags]

Flags:
  -h, --help             help for search
  -k, --keyword string   search keyword
  -t, --token string     personal access token
  -u, --url string       gitlab url
```

## Notice
1. Due to [GitLab rate limits](https://docs.gitlab.com/ee/user/gitlab_com/index.html#gitlabcom-specific-rate-limits), the search would be performed every 6s.
You may feel a little slow, but it can protect your self-hosted GitLab instance from too many requests.    
2. Considering that the goal is to find the target projects and not all occurrences, the results of each project will only appear at most 10.

## License
MIT