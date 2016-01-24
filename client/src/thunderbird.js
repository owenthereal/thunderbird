import Connection from "./connection"

export function connect(url, cb) {
  return new Connection(url, cb)
}
