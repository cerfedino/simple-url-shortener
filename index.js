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

