const isLocal = process.argv[2] !== "remote"

///////////////////
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

// Deep freezes the settings object.
const deepFreeze = obj => {
    Object.keys(obj).forEach(prop => {
        if (typeof obj[prop] === 'object' && !Object.isFrozen(obj[prop])) deepFreeze(obj[prop]);
    });
    return Object.freeze(obj);
};


module.exports = deepFreeze(settings)