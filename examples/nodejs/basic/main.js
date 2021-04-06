const process = require('process');
const { Person } = require('./pb/swapi/resources_pb');

async function main() {
  process.stdin.once('data', chunk => {
    const luke = Person.deserializeBinary(chunk);
    console.log('hi', luke.getName());
  });
}

main();
