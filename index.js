const express = require('express');
const app = express();
const port = process.env.PORT || 80;

const shortener = require('./routes/shortener')

const { logToConsole, logRequest } = require('./modules/logger.js')('index')

{
    app.listen(port, () => {
        logToConsole.GREEN(0, "+", `Server is running on port ${port}`);
    });

    app.use(shortener)

    // Serve static files
    app.use('/', express.static('public'));


    // Catch 404
    app.use((req, res, next) => {
        res.status(404).end();
    });

}