var cp = require("child_process");
var DEBUG = false 

exports.handler = function(req, res) {

    // Parse our the request from the body
    var queryStr = req.query.code

    if (DEBUG) {
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

    // Redirect to home page
    res.status(302).send(resp);
};
