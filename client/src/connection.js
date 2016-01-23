import Channel from "./channel"

export default class Connection {
  constructor(url) {
    this.url = url
    this.connect()
  }

  connect () {
    this.conn = new WebSocket(this.url);
    var self = this
    this.conn.onopen = function(evt) {
      console.log("connected to " + self.url)
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

  subscribe (channel) {
    return new Channel(channel)
  }
}
