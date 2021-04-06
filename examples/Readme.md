some other fun examples:

curl -s https://swapi.dev/api/people/1/ | pout -f swapi/resources.proto Person -- | node ./nodejs/basic/main.js 