#!/bin/bash

echo "running"

to_run="/app/ops $@"
echo "running $to_run"


eval $to_run