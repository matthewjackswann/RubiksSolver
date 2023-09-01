module github.com/matthewjackswann/rubiks

go 1.18

replace github.com/davidminor/uint128 v0.0.0-20141227063632-5745f1bf8041 => ./util/uint128

require (
	github.com/davidminor/uint128 v0.0.0-20141227063632-5745f1bf8041
	github.com/mattn/go-sqlite3 v1.14.15
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f
)
