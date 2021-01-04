const express = require('express')
var fs = require('fs');
const app = express()
const port = 80
var mapping = {}

const updateMapping = function(){
    fs.readFile('mapping.json', 'utf8', function (err, data) {
        if (err) throw err;
        mapping = JSON.parse(data);
    });
}

updateMapping();
app.listen(port, () => {
    console.log(`Server is running on port ${port}`)
    setInterval(function() {
        updateMapping();
    }, 5 * 60 * 1000)
})

app.get('/*', (req, res) => {
    res.set('Cache-Control', 'no-store')
    url = req.url.substr(1);
    console.log(`[.] Received request url "${url}"`);
    if (url in mapping){
        console.log(`\t[+] Found match for "${url}"\n\t Redirecting to ${mapping[url]} ...`);
        res.writeHead(301, {'Location' : mapping[url]});
    }else{
        console.log(`\t[-] Did not find any matches for "${url}"`);
    }

    res.end();
})



