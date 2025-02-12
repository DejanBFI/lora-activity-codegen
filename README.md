# LORA Activity Code Generator

DISLCLAIMER: This only works when the activity only calls one proxy. If the activity has multiple calls to proxies (which are possible), please resort to other ways.

## Prerequisites

- Go 1.23.4

## Preparation

Do this once after pulling the repository or every dependency update.

1. Run `make vendor`

## Code Generation

1. Edit `schemaName` in `main.go`. This refers to the schema name in LGS, without `.schema.json` suffix.
2. Edit `processAndActivityName` in `main.go`. This refers to the package name and the activity name in Temporal.
3. Edit `readSet`. This relates to the document fields needed to run the process. In this case, it's the request fields needed by the LGS schema.
4. Edit `writeSet`. This relates to the document fields where we will put the result. In the case, this will be the fields obtained from the LGS schema.
5. Run `go run main.go`.
6. Copy the output into `impl.go` file into any LORA worker repo.
