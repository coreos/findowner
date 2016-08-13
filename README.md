# findowner

findowner is used to find the "owner"s of a git directory.

It uses number of commits to rank. The program will pull past 1 year commits excluding "merge" commits. It ranks people by the number of commits they made and select top N (default is 3).

## Excluded Directories

It will walk from the top to each sub directories recursively excluding following directories:
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