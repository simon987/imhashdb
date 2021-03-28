# imhashdb

### Setup (development):

1. Install PostgreSQL 12
1. Install [pg_hamming](https://github.com/simon987/pg_hamming)
1. Create a database called `imhashdb` (`su postgres; createdb imhashdb`)
1. Install redis-server
1. Install [fastimagehash](https://github.com/simon987/fastimagehash)
1. Start an instance of `imhashdb hasher` (See `./imhashdb --help` for CLI options)
1. Start an image crawler (e.g. `reddit_feed`), with `ARC_LIST=imhash,...`
1. Start an instance of `imhashdb web`
1. Start the development web server [imhashdb-frontend](https://github.com/simon987/imhashdb-frontend)
 - Run `npm run start`
