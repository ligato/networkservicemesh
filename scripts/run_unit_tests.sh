#!/bin/bash
CUR_DIR=$(dirname "$0")
TEST_LIST=$("$CUR_DIR"/for-each-module.sh "go list ./..." | grep -v -e "sample" | grep -e "^github")
# shellcheck disable=SC2086
gotestsum --junitfile ~/junit/unit-tests.xml -- -short $TEST_LIST
