var fs = require('fs');
const {logToConsole, logRequest} = require('./logger.js');

var mapping = {}

module.exports = (isLocal) => {
    var client
    var model = {}
    if(!isLocal) {
        client = new (require('mongodb').MongoClient)(process.env.MONGODB_URI);
        client.connect()
            .then(client => {
                model.db = client.db("simple-url-shortener");
                model["mapping"] = model.db.collection("mapping")
                logToConsole.GREEN(0,"+","Fetched MongoDB database and collections")
        });
    }

    function extract_mapping(json) {
        ret = {}

        json.forEach( (json_map) => {
            json_map['shortUrls'].forEach( (shortUrl) => {
                ret[shortUrl] = json_map['longUrl'];
            });
        });
        return ret;
    }

    return {
        mapping : mapping,
        updateMapping :
            isLocal ?
                ()=> {
                    fs.readFile('mapping.json', 'utf8', (err, data) => {
                        if (err) throw err;
                
                        json = JSON.parse(data);
                                        
                        newMapping = extract_mapping(json);
                        Object.keys(mapping).forEach((k)=>{delete mapping[k]})
                        Object.keys(newMapping).forEach((k)=>{mapping[k] = newMapping[k]})

                        logToConsole.GREEN(0,"+",`Imported mapping.json`);
                    });
                }:


                async () => {
                    newMapping = extract_mapping(await model.mapping.find().toArray())

                    Object.keys(mapping).forEach((k) => {
                        delete mapping[k]
                    })
                    Object.keys(newMapping).forEach((k) => {
                        mapping[k] = newMapping[k]
                    })

                    logToConsole.GREEN(0, "+", `Imported mapping.json`);
                }
    }



}