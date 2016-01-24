import Events from "events"
import Channel from "./channel"

export default class Connection extends Events.EventEmitter {
  constructor(url, cb) {
    super()

    this.url = url
    this.cb = cb
    this.connect()

    this.channels = {}
  }

  connect () {
    var self = this

    this.ws = new WebSocket(this.url);
    this.ws.onopen = function(evt) {
      console.log("connected to " + self.url)
      self.cb(self)
    }
    this.ws.onclose = function(evt) {
      console.log("disconnected")
    }
    this.ws.onerror = function(evt) {
      console.log("error")
    }
    this.ws.onmessage = function(evt) {
      console.log("received message " + evt.data)

      var data = JSON.parse(evt.data)
      if (data.command === "message") {
        self.emit(data.channel, JSON.parse(data.body))
      }
    }
  }

  send (data) {
    this.ws.send(JSON.stringify(data))
  }

  subscribe (channel, cb) {
    var ch = new Channel(this, channel, cb)
    ch.subscribe()
    this.on(channel, cb)

    return ch
  }

  perform (channel, msg) {
    var data = { command: "message", channel: channel, body: JSON.stringify(msg) }
    console.log(data)
    this.ws.send(JSON.stringify(data))
  }
}
