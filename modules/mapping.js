var fs = require('fs');
const { logToConsole } = require('./logger.js')('mapping');
const config = require('../config');

/**
 * Imports the URL map, used by the shortening controller.
 */
var mapping = {}

// If the database is remote, connects to the MongoDB database with the provided connection string.
const mongodbClient = config.webserver.remote ? new (require('mongodb').MongoClient)(config.database.mongodb_uri).connect() : undefined;


function extract_mapping(json) {
    ret = {}

    json.forEach((json_map) => {
        json_map['shortUrls'].forEach((shortUrl) => {
            ret[shortUrl] = json_map['longUrl'];
        });
    });
    return ret;
}

const updateMapping = !config.webserver.remote ?
    () => {
        // If the db is local, reads the mapping.json file
        fs.readFile('mapping.json', 'utf8', (err, data) => {
            if (err) throw err;
            
            json = JSON.parse(data);

            newMapping = extract_mapping(json);
            Object.keys(mapping).forEach((k) => { delete mapping[k] })
            Object.keys(newMapping).forEach((k) => { mapping[k] = newMapping[k] })

            logToConsole.GREEN(0, "+", `Imported mapping.json`);
        });
        // If the database is remote, pulls the 'mapping' collection from the mongodb database.
    } : async () => {
            model = {}
            await mongodbClient.then(client => {
                model.db = client.db("simple-url-shortener");
                model["mapping"] = model.db.collection(config.database.collection)
            });

            newMapping = extract_mapping(await (model.mapping.find().toArray()))

            Object.keys(mapping).forEach((k) => {
                delete mapping[k]
            })
            Object.keys(newMapping).forEach((k) => {
                mapping[k] = newMapping[k]
            })

            logToConsole.GREEN(0, "+", "Fetched MongoDB database and collections")
    }


updateMapping();
setInterval(updateMapping, config.mapping.mappingImportDelay);

module.exports = mapping