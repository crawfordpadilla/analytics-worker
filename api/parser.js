const { Transform } = require('stream');
const { Console } = console;

class Parser extends Transform {
  constructor(options) {
    super(options);
    this.options = options;
  }

  _transform(chunk, encoding, callback) {
    try {
      const data = JSON.parse(chunk.toString());
      this.push(JSON.stringify({ ...data, parsed: true }));
      callback(null, chunk);
    } catch (e) {
      this.emit('error', e);
      callback();
    }
  }
}

module.exports = Parser;