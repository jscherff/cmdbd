#!/bin/sh

git rm -f log
git rm -f {i686,x86_64}/log
git rm -f {i686,x86_64}/*.json

rm -rf log
rm -rf {i686,x86_64}/log
rm -rf {i686,x86_64}/*.json
