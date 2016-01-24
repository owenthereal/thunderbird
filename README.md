# Thunderbird - Elegant WebSockets in Go

**I'm competing at [Gopher Gala 2016](http://gophergala.com/). Please give
this project a star and [follow me](https://github.com/jingweno) if you like what you see :beers:.**

Thunderbird (a.k.a. :zap::bird:) seamlessly integrates WebSockets with your Go web application. It allows for real-time features to be written in idiomatic Go. It's a full-stack offering that provides both a client-side JavaScript framework and a server-side Go framework. Thunderbird is heavily inspired by Elixir's [Phoenix](http://www.phoenixframework.org/) and Rails' [Action Cable](https://github.com/rails/rails/tree/master/actioncable).

## How it works

Thunderbird is built around the concept of connections and channels. It
has one connection per WebSocket connection. A single user may have
multiple WebSockets open to your application if they use multiple
browser tabs or devices. The client of a WebSocket connection is called
a consumer and it can subscribe to multiple channels. Each channel responds
to a name and encapsulates a logical unit of work similar to a controller
in a regular MVC setup. Just like how pub/sub works, a broadcast can be made
to consumers that subscribed to a channel.

Here is an example of multiple channels on multiple connections:

```
          ┌──────────┐
          │          │           Connection1
          │ Client1  │◀──────[channel1, channel2]────────┐
          │          │                                   │
          └──────────┘                                   ▼
                                                ┌────────────────┐
                          Connection2           │                │
                ┌─────[channel2, channel3]─────▶│     Server     │
                │                               │                │
                ▼                               └────────────────┘
          ┌──────────┐                                   ▲
          │          │           Connection1             │
          │ Client2  │◀───────────[channel3]─────────────┘
          │          │
          └──────────┘
```

## Tutorial: a chat app

Let's walk through building a chat app with `Thunderbird` and understand
how elegant the solution is. The code is available in the
[example](https://github.com/gophergala2016/thunderbird/tree/master/example) folder.
The example is running live [here](https://thunderbird-chat.herokuapp.com/).

### Server handling WebSocket connections

The first thing is to initialize `Thunderbird` and create a route to
handle WebSocket connections:

```go
func main() {
  tb := thunderbird.New()

  mux := http.NewServeMux()
  mux.Handle("/ws", tb.HTTPHandler())
}
```

`Thunderbird.HTTPHandler()` returns a `http.Handler` that deals with
WebSocket connections.

### Server handling a channel

We then create a channel handler that implements the
`ChannelHandler` interface. `ChannelHandler`s are registered with `Thunderbird.HandleChannel` which takes a name and a `ChannelHandler`:

```go
type RoomChannel struct {
}

func (rc *RoomChannel) Received(event thunderbird.Event) {
    // handle event
}

func main() {
    ch := &RoomChannel{}
    tb.HandleChannel("room", ch)
}
```

### Server broadcasting messages to consumers

For a chat app, we want `RoomChannel` to broadcast received
messages to all consumers of the channel. We do this by calling
`Thunderbird.Broadcast` when there is a message event:

```go
type RoomChannel struct {
	  tb *thunderbird.Thunderbird
}

func (rc *RoomChannel) Received(event thunderbird.Event) {
	  switch event.Type {
	  case "message":
		  rc.tb.Broadcast(event.Channel, event.Body)
	  }
}
```

`Event.Type` represents the type of the event received. Currently the value can be `subscribed`, `unsubscribed` or `message` which mean a subscription was made to the channel, a unsubscription was made to channel or a message was received on the channel correspondingly. More event types will be support in the future.

So far the Go side of the chat app is done. That was easy, right :-)? Let's take a
look at the JavaScript part of the app.

### Client connecting

On the client, we call `Thunderbird.connect` to connect to the Go server.
The method takes a URL and a callback that accepts a connection object.
In thie case `url` should be the WebSocket URL that we set up in
[previous step](https://github.com/gophergala2016/thunderbird#server-handling-websocket-connections): `ws://localhost:8080/ws`.

```js
Thunderbird.connect(url, function (conn) {
  // do something with the connection `conn`
})
```

### Client subscribing to a channel

With the `conn`, we can subscribe to a channel (`room` in our example) and do something with the received message in the callback:

```js
conn.subscribe("room", function (msg) {
  // handle message
})
```

### Client sending messages

For a chat app, we would like to send messages to other consumers and we
do that with the `perform` method of the connection:

```js
// when enter key is pressed
conn.perform("room", msg)
```
