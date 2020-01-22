#!/bin/sh

cd mock_oriservice && CGO_ENABLED=0 go test srv_mock_test.go

cd ..

cd srv && CGO_ENABLED=0 go test main_test.go
