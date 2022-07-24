const express = require('express')
const router = express.Router()
const mapping = require('../modules/mapping.js')
const { logToConsole, logRequest } = require('../modules/logger.js')('shortener');

/**
 * 
 * Router for the URL shortening service.
 * 
 */


module.exports = router

/**
 * GET /help
 * 
 * Returns the URL mapping.
 */
router.get('/help', (req,res)=>{
    res.send(mapping)
})


/**
 * GET /*
 * 
 * Check if the request url is mapped. If so, redirects.
 */
router.get('/*', (req, res, next) => {
    res.set('Cache-Control', 'no-store')

    url = req.url.substr(1);
    logToConsole.YELLOW(0,".",`Received request url "${url}"`);

    if (url in mapping) {
        logToConsole.GREEN(1,"+",`Found match for "${url}"`);
        logToConsole.YELLOW(2,"-",`Redirecting to ${mapping[url]} ...`);
        logRequest(req);
        res.writeHead(301, {'Location' : mapping[url]}).end();
    } else {
        logToConsole.RED(1,".",`Did not find any matches for "${url}"`);
        next();
    }
});