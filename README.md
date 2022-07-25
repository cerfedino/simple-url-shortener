# Description
This is a simple url shortener written in Javascript.\
Takes the URL mapping either from a local json file or from a mongodb database.

---
**TOC**
<!-- TOC -->

- [Description](#description)
- [Installing and running](#installing-and-running)
        - [NPM run scripts](#npm-run-scripts)
- [URL mapping format](#url-mapping-format)
- [Configuration](#configuration)
- [Routes](#routes)
        - [URL shortening middleware](#url-shortening-middleware)
        - [Serving static files inside of the public folder](#serving-static-files-inside-of-the-public-folder)
- [IP logging feature](#ip-logging-feature)
- [systemd service unit file](#systemd-service-unit-file)

<!-- /TOC -->
---

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

# IP logging feature
All the IPs making requests to the shortening service get logged inside the `log.csv`.\
A friend of mine told me that it is basically spyware :see_no_evil:.\
I guess until no one complains it is legal.

# systemd service unit file
I've added a systemd unit file so that I can automatically run the application as a background process when booting my machine.