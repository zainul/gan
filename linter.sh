#!/bin/bash
# make continue when error
set +e
# make break when error
# set -e

GometalinterVariable=(
  "deadcode"
  "dupl"
  "errcheck"
  "gochecknoglobals"
  "gochecknoinits"
  "goconst"
  "gocyclo"
  "gofmt"
  "goimports"
  "golint"
  "gosec"
  "gosimple"
  "gotype"
  "gotypex"
  "ineffassign"
  "interfacer"
  "lll"
  "maligned"
  "megacheck"
  "nakedret"
  "safesql"
  "staticcheck"
  "structcheck"
  "test"
  "testify"
  "unconvert"
  "unparam"
  "unused"
  "varcheck"
  "vet"
  "vetshadow"
)


Directory=(
            "migrations/seed"
            "migrations/version"
          )

arrayGometalinterVariable=${#GometalinterVariable[@]}
arrayDirectory=${#Directory[@]}


for ((k=0; k<${arrayDirectory}; k++));
do
  for ((i=0; i<${arrayGometalinterVariable}; i++));
  do
    echo "Currently linter running in ${Directory[$k]} ==> ${GometalinterVariable[$i]}"
    gometalinter -j 1 --disable-all --enable=${GometalinterVariable[$i]}  ${Directory[$k]}/  2>&1

    sleep 1
    wait

  done
done