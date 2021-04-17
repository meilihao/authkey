#! /bin/bash

Version=$(git describe --tags --dirty --always)
GitBranch=$(git rev-parse --abbrev-ref HEAD)
GitHash=$(git rev-parse HEAD)
BuildTS=$(date -u --rfc-3339=seconds)

LDFLAGS="-X authkey/pkg/util.version=${Version}
         -X authkey/pkg/util.gitBranch=${GitBranch}
         -X authkey/pkg/util.gitHash=${GitHash}
         -X 'authkey/pkg/util.buildTimestamp=${BuildTS}'" # 如果要赋值的变量包含空格，需要用引号将 -X 后面的变量和值都扩起来, 否则会编译报错

go build -tags=jsoniter -ldflags "$LDFLAGS" -v -o apiserver ./main.go