# findowner

findowner is used to find the "owner"s of a git directory.

It uses number of commits to rank. The program will pull commits from last 6 months excluding "merge" commits. It ranks people by the number of commits they made and select top N (default is 3).

## exowner - existing owners

exowner is used to print out the owners in the OWNER files.
