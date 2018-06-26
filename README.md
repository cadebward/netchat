# netchat

A small TCP server for handling chat.

### Running

```bash
$ go run main.go
```

```bash
$ telnet 127.0.0.1 3000
```

```bash
$ telnet 127.0.0.1 3000
```

### Known Bugs

- If your client receives a message while you are typing, your message is broken up into multiple lines

### Approach


The approach I took was largely driven around the main parts of the application that I could see: 

1. The client
2. The server
3. The main initialization of the application

`main.go` is the file that gathers the configurations and starts up the entire applications. `server/client.go` is the code that drives the behavior for each individual client that connects to the server. In this file you’ll find convenient methods for sends messages and receiving messages from a single client. In `server/server.go` I wrote the glue to tie many clients together. `server.go` utilizes all the convenience methods in `client.go` to broadcast messages to all clients, send messages to singles clients, and relay each client’s messages to all the other clients. I spent a good amount of time really trying to make these boundaries as clean as I could, and I feel like it really paid off in making the code easy to follow.
