#! /bin/bash
env GOOS=js GOARCH=wasm go build -o gandermerge.wasm github.com/iancanderson/gandermerge

git add gandermerge.wasm
git stash
git checkout gh-pages
git stash pop
git commit -am "New web release"
git checkout main
