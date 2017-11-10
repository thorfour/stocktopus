var DEBUG = true
var cp = require("child_process");
var request = require('request')
var headers = {
    'User-Agent': 'Super Agent/0.0.1',
    'Content-Type': 'application/json',
}

exports.handler = function(req, res) {

    var queryStr = JSON.stringify(req.body);

    if (DEBUG) {
        console.log(queryStr);
    }

    // Spawn the go routine to lookup stock quote
    var proc = cp.spawnSync("./serverless", [queryStr], {stdio: 'pipe', encoding: "utf8"});
    var quote = proc.stdout;

    var respType = "in_channel";
    // Check for no response, means there was an error
    if (quote === "") {
        quote = proc.stderr
        respType = "ephemeral";

        if (DEBUG) {
            console.log(proc.stderr);
        }
    }

    // Parse quote into json for slack
    var resp = '{ "response_type" : "' + respType + '", "text" : "' + quote + '" }';

    var options = {
        url: req.body.response_url,
        method: 'POST',
        headers: headers,
        form: resp
    }

    if (DEBUG) {
        console.log(req.body.response_url);
    }

    // Return json
    request(options);

    res.status(200).end();
};
