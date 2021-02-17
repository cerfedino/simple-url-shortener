const express = require('express')
var fs = require('fs');
const app = express()
const port = process.env.PORT || 80
var mapping = {}

const updateMapping = function(){
    fs.readFile('mapping.json', 'utf8', function (err, data) {
        if (err) throw err;

        json = JSON.parse(data);

        newMapping = {};
        json.forEach(function(json_map){
            json_map['shortUrls'].forEach(function (shortUrl){
                newMapping[shortUrl] = json_map['longUrl'];
            });
        });
        mapping = newMapping;

        consoleStatus("",GREEN,"+",`Imported mapping.json`);
    });
}
updateMapping();

const logRequest = function(req){
    var ip = req.header('x-forwarded-for') || req.connection.remoteAddress;

    let date_ob = new Date(Date.now());
    let date = date_ob.getDate();
    let month = date_ob.getMonth() + 1;
    let year = date_ob.getFullYear();

    let hours = date_ob.getHours();
    let minutes = date_ob.getMinutes();
    let seconds = date_ob.getSeconds();

    fs.appendFile('log.csv', `${year}-${month}-${date};${hours}:${minutes}:${seconds};${ip};${req.url}\n`, function (err) {
        if (err) throw err;
        consoleStatus("\t",GREEN,"+",`Logged into CSV file!`);
      });
}


app.listen(port, () => {
    consoleStatus("",GREEN,"+",`Server is running on port ${port}`)
    setInterval(function() {
        updateMapping();
    }, 5 * 60 * 1000)
})

app.get('/*', (req, res) => {
    res.set('Cache-Control', 'no-store')
    url = req.url.substr(1);
    consoleStatus("",YELLOW,".",`Received request url "${url}"`);
    if (url in mapping){
        consoleStatus("\t",GREEN,"+",`Found match for "${url}"\n\t\t Redirecting to ${mapping[url]} ...`);
        logRequest(req);
        res.writeHead(301, {'Location' : mapping[url]});
    }else{
        consoleStatus("\t",RED,"-",`Did not find any matches for "${url}"`);
    }

    res.end();
})

const RED = "\x1b[31m";
const GREEN = "\x1b[32m"
const YELLOW = "\x1b[33m";
const consoleStatus = function(prepend,COLOR,symbol,text){
    console.log(`${prepend}[${COLOR}${symbol}`+'\x1b[0m]'+` ${text}`);
}


