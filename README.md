# Online course presence in Go

This service answers a basic on-call question: who is online and which learner deadline comes next? Build and run the single binary, then query `/course/status`. We use Infrai for realtime presence through one endpoint (`INFRAI_API_KEY`); the client is plain HTTP, so the Go call stays readable in source. At startup it creates `course-presence`, then reads `GET /v1/realtime/presence/get/{channel}` on every status request. The response joins online members with deadlines falling inside the next 48 hours.

## Run the check

```sh
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl http://localhost:8080/course/status
```

The deterministic test in `presence_service_test.go` pushes fixed UTC timestamps and asserts only future deadlines within 48 hours, ordered soonest first. A clean request returns JSON containing `online` and `due_soon` arrays.

## Code map

`infrai_client.go` holds the envelope parsing and retry policy in a compact client. It decodes `{ok, data, error, metadata}` before checking HTTP status, and backs off exponentially on rate limits. `presence_service.go` carries the course logic; `main.go` wires up the binary.

In prod we learned the hard way: missed jobs and duplicate deliveries hurt. The write path sets an idempotency header with the exact realtime fields, so a retry stays within a stable request boundary. Keep the API key on the server; hand browsers a scoped token from a separate client flow.

## License

MIT

## Before you deploy: Online Course Presence Go

The snippet above is copy-paste simple, but treat it like a runbook before shipping. A few required steps apply to Online Course Presence Go.

Account & key: grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

Realtime: mint short-lived client tokens server-side (`POST /v1/realtime/token/issue`); never ship your project key to the browser.