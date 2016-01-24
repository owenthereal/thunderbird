var path = require("path")

module.exports = {
  context: __dirname,
  entry: "./src/thunderbird.js",
  output: {
    path: path.join(__dirname, "lib"),
    filename: "thunderbird.js",
    library: "Thunderbird",
    libraryTarget: "umd"
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'babel',
        query: {
          presets: ['es2015']
        }
      }
    ]
  }
}
