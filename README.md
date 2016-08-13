# findowner

findowner is used to find the "reviewer"s of a git directory.

The program uses number of commits to rank. It will pull past 1 year commits excluding "merge" commits. It ranks people by the number of commits they made and select top N (default is 3).

## Directory Crawling Policy

The program will walk from the top to each sub directories recursively. It has a global limit on depth (default is 3).
For example, "./pkg/apis/apps" is allowed, while "./pkg/apis/apps/install" isn't.

It also has rules to exclude following directories:
```
		"vendor",
		"contrib/mesos/", // we don't need to go recursively
		// exclude generated code: `find . | grep "generated"` + some guessing
		"staging",
		"cmd/libs/go2idl/client-gen",
		"federation/client/clientset_generated",
		"pkg/client/clientset_generated",
```
This is hardcoded currently.

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