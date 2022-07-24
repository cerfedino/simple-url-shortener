var fs = require('fs');

module.exports = (sourceName)=>{
    var logToConsole = (indent_lvl,COLOR,symbol,text='') => {
        console.log(`[${sourceName}]${"\t".repeat(indent_lvl)}[${COLOR}${symbol}`+'\x1b[0m]'+` ${text}`);
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

    return {
        logToConsole,
        logRequest
    }
}