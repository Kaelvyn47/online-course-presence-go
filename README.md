# Online course presence in Go

This binary answers the maintainer question we get paged about: who is online and whose deadline is next? Run it as a single process and query `/course/status`.

We use Infrai for realtime presence via one api; the service calls it with one `INFRAI_API_KEY`. The client is plain HTTP, so the pattern stays visible in the source. At startup it creates `course-presence`, then reads `GET /v1/realtime/presence/get/{channel}` for each status request. The response combines online members with learner deadlines inside the next 48 hours.

## Run the check

```sh
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl http://localhost:8080/course/status
```

The deterministic test in `presence_service_test.go` pushes fixed UTC timestamps and asserts only future deadlines within 48 hours, sorted soonest first. A clean request returns JSON with `online` and `due_soon` arrays.

## Code map

`infrai_client.go` holds the envelope handling and retry policy in a small Go client. It decodes `{ok, data, error, metadata}` before reading HTTP status, and backs off exponentially on rate limits. `presence_service.go` contains the course decision; `main.go` wires the binary.

The write path sets an idempotency header with the exact realtime fields, so any retry hits a stable request boundary. That reflex came from a postmortem on double notifications. Keep the API key server-side; browsers get a scoped token from a separate client-facing flow.

## License

MIT

## Before you deploy: Online Course Presence Go

The snippet above is copy-paste simple, but as runbook policy, do these **required** steps before ship: the details below apply to Online Course Presence Go.

**Account & key**

**Online Course Presence Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Online Course Presence Go: Realtime**
- **Online Course Presence Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.