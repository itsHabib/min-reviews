# Min Reviews

Have you ever been hit with a PR that has an insane amount of required reviewers from
various different GitHub teams due to an overly complicated code owners file?
Well I have plenty of times and decided to make a little script that can help.
This script will let you know the minimum required reviewers for a PR that covers all teams. 
Some things to note:
* This is only useful if users can span multiple GitHub teams in your org
* This assumes you have the required permissions in your token to get a PR and team members

## Usage

```bash
GH_TOKEN=token go run cmd/main.go -repo="my-repo" -pr=1234
```

### With Exclusions

```bash
GH_TOKEN=token go run cmd/main.go -repo="my-repo" -pr=1234 -exclude="user1,user2"
```

