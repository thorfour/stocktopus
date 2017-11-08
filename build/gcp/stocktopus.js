var cp = require("child_process");
var DEBUG = false

exports.handler = function(event, context) {

    // Parse our the request from the body
    var queryStr = unescape(event.body)

    if (DEBUG) {
        console.log(queryStr)
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
            console.log(proc.stderr)
        }
    }

    // Parse quote into json for slack
    var resp = '{ "response_type" : "' + respType + '", "text" : "' + quote + '" }';

    // Return json
    context.succeed(resp);
};
