# sha-versioning
Deriving the number of revisions of a file in git based on SHA value in releases

## To Do

Implement a check to retrieve the tags rather than provide them as an input. The essential URL is:
```
https://qa.door43.org/api/v1/repos/unfoldingword/en_twl/releases
```

## setup

```sh
$ go mod init shaversioning
go: creating new go.mod: module shaversioning
go: to add module requirements and sums:
        go mod tidy
$ go mod tidy
go: finding module for package golang.org/x/mod/semver
go: found golang.org/x/mod/semver in golang.org/x/mod v0.4.2
```

Now you can run the bash scripts, for example:
```
sh run_twl.sh
```
which produces `en_twl_revs.csv`

**WARNING: at the moment, this only works with "tsv" type resources!!**