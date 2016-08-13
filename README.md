# findowner

findowner is used to find the "reviewer"s of a git directory:
- It crawls all directories of given "root" directory.
  - It has a depth limit (default is 3). For example, given depth limit 3,
    "./pkg/apis/apps" is allowed, while "./pkg/apis/apps/install" isn't.
  - It has an excluded list. See "./main.go" `excludedDirList` variable.
- It selects top committers for each git directory.
  - It ranks people (identified by github handle) by the number of commits they made.
  - The commits are pulled from past 1 year history excluding "merge" commits.
  - It selects top N (default is 3). If there are more than N people having the same number of commits, it outputs them all.

## exowner - existing owners

exowner is used to print out the owners in the OWNER files. Usage:

```
$ git clone https://github.com/coreos/findowner
$ cd findowner
$ go run exowner/main.go  -gitrepo ~/src/k8s.io/kubernetes/ -top-dir pkg
...
path: pkg/volume/host_path, owners: [saad-ali thockin]
...
```