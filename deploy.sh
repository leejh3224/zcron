#!/usr/bin/env bash

cd functions

GOOS=linux go build -o ../bin/zcron .

cd ../infra

yarn deploy
