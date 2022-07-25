# Description
This is a simple url shortener written in Javascript.\
Takes the URL mapping either from a local json file or from a mongodb database.


# Installing and running
Install the dependencies with
`npm install`
###  NPM run scripts
- ```bash
  npm run start
  ```
  Gets the URL mapping by connecting to the mongodb database. See `config.js` for the mongodb connection string.
- ```bash
  npm run local
  ```
  Gets the URL mapping from the local `mapping.json` file


# URL mapping format
Example:
`yourdomain.com/very/long/url` will map as: 
`"very/long/url"`<br><br>

```json
[
  {
    "shortUrls" : ["git", "github"],
    "longUrl" : "https://github.com/AlbertCerfeda"
  },
  {
    "shortUrls" : ["gl", "google"],
    "longUrl" : "https://www.google.com/"
  },
  {
    "shortUrls" : ["alias"],
    "longUrl" : "/git"
  }
]
```
Mapping can be changed on runtime as every 5-minutes it gets imported again.

# Configuration
```js
//////////////////
// Reads environment variables, or grabs default values.
const PORT = process.env.PORT || 80
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017'

///////////////////

const settings = {
    webserver: {
        remote: !isLocal,
        port: PORT,
    },

    mapping: {
        // How frequently the mapping gets imported. 
        mappingImportDelay: 1000*60*5
    },

    database: {
        mongodb_uri: MONGODB_URI,
        db_name: "simple-url-shortener",
        collection: "mapping"
    }
}
```

# Routes
### URL shortening middleware
This controller checks the request URL for matching entries in the URL mapping.
- `GET /help`\
  Dumps the mapping in the answer
  
### Serving static files inside of the `public` folder
It is possible to serve static files placed inside the public folder.

The URL shortening middleware has the highest priority over the static file serving.
