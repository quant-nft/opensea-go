#!/usr/bin/env bash
git tag -d $*
git push --delete origin $*