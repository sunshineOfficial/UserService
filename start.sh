#!/bin/sh
if [ "$ENV" == "utests" ]; then
    go test -v;
else
    user-service;
fi