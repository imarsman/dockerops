#!/bin/bash

# This script takes an argument to docker run and turns it into a complete call
# invocation for ops. If this is not done then sub-parts of the call like
# "image list" in "ops image list" will not be handled properly. Thus this
# script makes a proper invocation and then uses it to call ops.

# Get args to docker run call as a single string
args="$*"

# Remove any leading ops invocaton using greedy
args="${args#*ops }"

# Create the call, adding back ops invocation
to_run="/app/ops $args"

echo $to_run

# Run created call
eval "$to_run"