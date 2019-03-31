workflow "Lint" {
  resolves = ["izumin5210/actions-reviewdog"]
  on = "pull_request"
}

action "izumin5210/actions-reviewdog" {
  uses = "izumin5210/actions-reviewdog/golang@master"
  secrets = ["GITHUB_TOKEN"]
}
