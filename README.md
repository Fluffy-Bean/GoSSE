shrimple go sse example

| path              | what dat do                                                                      |
|-------------------|----------------------------------------------------------------------------------|
| /cmd/web          | entrypoint                                                                       |
| /internal/handler | web handler, passing around reference to sse, could add database connection here |
| /internal/routes  | web routes                                                                       |
| /internal/sse     | handles connections, heartbeats and broadcasting messages                        |
| /templates        | contains the js to start an `EventSource`                                        |

todo:
- authenticate that messages are coming from a specific user by id or whatnot
- save messages to repopulate chat history on refresh