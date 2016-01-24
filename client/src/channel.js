export default class Channel {
  constructor (conn, name, cb) {
    this.conn = conn
    this.name = name
    this.cb = cb
  }

  subscribe () {
    var data = { command: "subscribe", channel: this.name}
    this.conn.send(data)
  }
}
