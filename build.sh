#!/usr/bin/env bash
project_path=$(cd `dirname $0`; pwd)
mkdir -p "${project_path}/version"
GITCOMMIT=$(git rev-parse --short HEAD)
VERSION=$(git describe --abbrev=0 --tags)
BUILDTIME=$(date -u -d "@${SOURCE_DATE_EPOCH:-$(date +%s)}" --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
cat > ${project_path}/version/version.go <<DVEOF
// +build autogen
// Package version is auto-generated at build-time
package version
// Default build-time variable for library-import.
// This file is overridden on build with build-time information.
const (
	GitCommit             string = "$GITCOMMIT"
	Version               string = "$VERSION"
	BuildTime             string = "$BUILDTIME"
)
DVEOF
go build -work ${project_path}/cmd/p2pssh/p2pssh.go