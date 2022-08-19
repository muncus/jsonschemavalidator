#!/bin/bash
# The entrypoint script ensures the args from action.yaml (inputs.documents)
# become multiple args, instead of a single quoted string. This also expands
# file globs.
/app/schemacheck $*