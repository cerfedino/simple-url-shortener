const app =  require('express')();
const port = process.env.PORT || 8888;
const {updateMapping, mapping} = require('./modules/mapping.js')(process.argv[2]!="remote")

const {logToConsole, logRequest} = require('./modules/logger.js')

{
    setTimeout(updateMapping, 2000)
    
    app.listen(port, () => {
        logToConsole.GREEN(0,"+",`Server is running on port ${port}`);
        setInterval(updateMapping, 5*60*1000)
    });


    app.get('/*', (req, res) => {
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

        res.end();
    });
}