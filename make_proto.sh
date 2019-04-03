#!/bin/sh
set -x

protoc -I routeguide/ routeguide/route_guide.proto --go_out=plugins=grpc:routeguide
