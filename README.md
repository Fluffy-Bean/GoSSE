shrimple go sse example

| path              | what dat do                                                                             |
|-------------------|-----------------------------------------------------------------------------------------|
| /cmd/web          | entrypoint                                                                              |
| /internal/handler | web handler, passing around reference to sse, could add database connection here        |
| /internal/routes  | web routes                                                                              |
| /internal/sse     | handles connections, heartbeats and broadcasting messages                               |
| /internal/jwt     | encoding and decoding of jwt tokens, useful for this example as we dont have a database |
| /templates        | contains the js to start an `EventSource`                                               |

todo:
- save messages to repopulate chat history on refresh
