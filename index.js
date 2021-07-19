var fs = require('fs');
const app = require('express')();
const port = process.env.PORT || 80;

var mapping = {};


{
    updateMapping();
    
    app.listen(port, () => {
        logToConsole.GREEN(0,"+",`Server is running on port ${port}`);
        setInterval(updateMapping, 5*60*1000)
    });


    app.get('/*', (req, res) => {
        res.set('Cache-Control', 'no-store')

        url = req.url.substr(1);
        logToConsole.YELLOW(0,".",`Received request url "${url}"`);

        if (url in mapping){
            logToConsole.GREEN(1,"+",`Found match for "${url}"\n\t\t Redirecting to ${mapping[url]} ...`);
            logRequest(req);
            res.writeHead(301, {'Location' : mapping[url]});
        }else{
            logToConsole.RED(1,"-",`Did not find any matches for "${url}"`);
        }

        res.end();
    });
}



function updateMapping() {
    fs.readFile('mapping.json', 'utf8', (err, data) => {
        if (err) throw err;

        json = JSON.parse(data);

        newMapping = {};
        json.forEach( (json_map) => {
            json_map['shortUrls'].forEach( (shortUrl) => {
                newMapping[shortUrl] = json_map['longUrl'];
            });
        });
        mapping = newMapping;

        logToConsole.GREEN(0,"+",`Imported mapping.json`);
    });
};


function logRequest(req) {
    var ip = req.header('x-forwarded-for') || req.connection.remoteAddress;

    let date_ob = new Date(Date.now());
    let date = date_ob.getDate();
    let month = date_ob.getMonth() + 1;
    let year = date_ob.getFullYear();

    let hours = date_ob.getHours();
    let minutes = date_ob.getMinutes();
    let seconds = date_ob.getSeconds();

    fs.appendFile('log.csv', `${year}-${month}-${date};${hours}:${minutes}:${seconds};${ip};${req.url}\n`, (err) => {
        if (err) throw err;
        logToConsole.GREEN(1,"+",`Logged into CSV file!`);
    });
};


var logToConsole = (indent_lvl,COLOR,symbol,text) => {
    console.log(`${"\t".repeat(indent_lvl)}[${COLOR}${symbol}`+'\x1b[0m]'+` ${text}`);
};

logToConsole.RED = (indent_lvl,symbol,text) => {
    logToConsole(indent_lvl, "\x1b[31m",symbol,text);
};

logToConsole.GREEN = (indent_lvl,symbol,text) => {
    logToConsole(indent_lvl, "\x1b[32m",symbol,text);
};

logToConsole.YELLOW = (indent_lvl,symbol,text) => {
    logToConsole(indent_lvl, "\x1b[33m",symbol,text);
};