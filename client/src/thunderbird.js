import Connection from "./connection"

export function connect(url) {
  return new Connection(url)
}
