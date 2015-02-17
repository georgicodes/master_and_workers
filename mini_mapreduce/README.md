Approach
======

## Single machine, serial - singlecore.go
* reads in file
* runs map phase where it keeps count of each word found
* program exits prints listing of all words in sorted order

## Multi machine, multi-core, parallel
* reads in file
* splits up file based on gomaxprocs
* runs gomaxprocs number of goroutines for map phase where it keeps count of each word found (fan out)
* fan's in responses from each map routine. It will perform reduce phase
* program exits prints listing of all words in sorted order

## Distributed
TBD
