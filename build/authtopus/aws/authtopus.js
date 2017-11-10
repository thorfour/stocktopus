var cp = require("child_process");
var DEBUG = true

exports.handler = function(event, context) {

    // Parse our the request from the body
    var queryStr = event.query.code

    if (DEBUG) {
        console.log(event)
        console.log(queryStr)
    }

    // Spawn the go routine to lookup stock quote
    var proc = cp.spawnSync("./oauth", [queryStr], {stdio: 'pipe', encoding: "utf8"});
    var resp = proc.stdout;

    // Check for no response, means there was an error
    if (resp === "") {
        resp = proc.stderr

        if (DEBUG) {
            console.log(proc.stderr)
        }
    }

    context.succeed(resp);
};
