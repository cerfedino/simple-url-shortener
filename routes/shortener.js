const express = require('express')
const router = express.Router()
const mapping = require('../modules/mapping.js')
const { logToConsole, logRequest } = require('../modules/logger.js');

module.exports = router


router.get('/*', (req, res, next) => {
    res.set('Cache-Control', 'no-store')

    url = req.url.substr(1);
    logToConsole.YELLOW(0,".",`Received request url "${url}"`);

    if (url in mapping) {
        logToConsole.GREEN(1,"+",`Found match for "${url}"\n\t\t Redirecting to ${mapping[url]} ...`);
        logRequest(req);
        res.writeHead(301, {'Location' : mapping[url]});
    } else {
        logToConsole.RED(1,"-",`Did not find any matches for "${url}"`);
    }

    next();
});