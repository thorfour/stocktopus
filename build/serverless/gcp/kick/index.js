var nextFuncURL = "STOCKTOPUS ENDPOINT" 
var request = require('request')
var headers = {
    'User-Agent': 'Super Agent/0.0.1',
    'Content-Type': 'application/json',
}

exports.handler = function(req, res) {

    // Immediately send OK response
    res.status(200).end();

    var options = {
        url: nextFuncURL,
        method: 'POST',
        headers: headers,
        body: JSON.stringify(req.body)
    }

    // Pass request to next cloud function
    request(options);
};
