const config = require('./config')
const express = require('express');
const app = express();

const port = config.webserver.port;

const shortener = require('./routes/shortener')

const {logToConsole} = require('./modules/logger.js')('index')

{
    app.listen(port, () => {
        logToConsole.GREEN(0, "+", `Server is running on port ${port}`);
    });

    // Run requests first through the shortener routes
    app.use(shortener)

    // Serve static files
    app.use('/', express.static('public'));


    // Catch 404
    app.use((req, res, next) => {
        res.status(404).end();
    });

}