var cSpawn = require("child_process").spawn;

exports.handler = function(event, context) {

    var proc = cSpawn("./colinmc", [event.text.trim()]);
    var quote = "";
    proc.stdout.on("data", function(buf) {
        quote += buf;
    });

    proc.on("close", function(code) {
        if (code == 0) {
            context.succeed(quote);
        }
    });
};
