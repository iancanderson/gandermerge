#! /bin/bash
env GOOS=js GOARCH=wasm go build -o spookypaths.wasm github.com/iancanderson/spookypaths

git add spookypaths.wasm
git stash
git checkout gh-pages
git stash pop
git checkout --theirs spookypaths.wasm
git commit -am "New web release"
git push
git checkout main
