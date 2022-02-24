Disclaimer: This is a one-day project therefore its not as refined as it should be
# Description
This is a simple url shortener written in Javascript.
# Installing and running
To install the dependencies, do:
`npm install`
then run by using
`sudo npm start`

# JSON mapping format
Example:
`yourdomain.com/very/long/url` will map into `mapping.json` as: 
`"very/long/url"`<br><br>
**Make sure to always specify *https* or *https* in the destination url**
```json
[
  {
    "shortUrls" : ["git", "github"],
    "longUrl" : "https://github.com/AlbertCerfeda"
  },
  {
    "shortUrls" : ["gl", "google"],
    "longUrl" : "https://www.google.com/"
  }
]
```
The file can be edited on runtime because every 5-minutes the program imports the mapping again
