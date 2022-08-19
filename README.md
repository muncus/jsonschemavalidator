# JSON Schema Validator

This tool validates JSON or YAML documents against [json
schemas](http://json-schema.org). It can be run either as a standalone command,
or through github actions.

## CLI Usage

From a clone of the repo, the tool can be built with `go build ./cmd/schemacheck`

General usage: `schemacheck --output text -s $SCHEMA $DOC...`

* `$SCHEMA` may be a local file in json schema format, or a url (e.g. `https://json.schemastore.org/github-workflow.json`)
* `$DOC` is one or more json or yaml documents that should be validated with the provided schema.
* The `--output` flag chooses between text and github-formatted output.

## Github Action Usage

### Inputs

* `schema`: **Required** A path or URL for a JSON Schema definition.
* `documents`: **Required** One or more space-separated file paths or globs
     containing the documents to verify with the given schema.

### Outputs

None.

### Example Action Usage

```
- name: "Validate Github Workflows"
  uses: muncus/jsonschemavalidator@v1
  with:
    schema: https://json.schemastore.org/github-workflow.json
    documents: "${GITHUB_WORKSPACE}/.github/workflows/*.yaml"
```


## Future Work

* [ ] tests
* [ ] Release tags for GHA usage