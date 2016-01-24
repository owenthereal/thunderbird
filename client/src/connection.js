import Channel from "./channel"

export default class Connection {
  constructor(url, cb) {
    this.url = url
    this.cb = cb
    this.connect()
  }

  connect () {
    var self = this

    this.conn = new WebSocket(this.url);
    this.conn.onopen = function(evt) {
      console.log("connected to " + self.url)
      self.cb(self)
    }
    this.conn.onclose = function(evt) {
      console.log("disconnected")
    }
    this.conn.onerror = function(evt) {
      console.log("error")
    }
    this.conn.onmessage = function(evt) {
      console.log(evt.data)
    }
  }

  subscribe (channel, cb) {
    var data = { command: "subscribe", channel: channel}
    this.conn.send(JSON.stringify(data))
  }

  trigger (channel, evt) {
    var data = { command: "broadcast", channel: channel, body: JSON.stringify(evt) }
    this.conn.send(JSON.stringify(data))
  }
}
