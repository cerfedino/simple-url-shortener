const isLocal = process.argv[2] !== "remote"

///////////////////
// Initializes some settings that depend on whether the application is getting deployed locally or remotely
const PORT = isLocal ?  80 : process.env.PORT
const MONGODB_URI = process.env.MONGODB_URI

///////////////////

const settings = {
    webserver: {
        remote: !isLocal,
        port: PORT,
    },

    mapping: {
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