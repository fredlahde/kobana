#!/bin/bash

for p in $(ls -d */ | grep -v test_fixtures | sed -e 's|/||'); do
  go test "./$p"
done
