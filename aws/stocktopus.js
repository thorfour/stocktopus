var cp = require("child_process");

exports.handler = function(event, context) {

    var proc = cp.spawnSync("./colinmc", [event.text], {encoding : "utf8"});
    context.succeed(proc.stdout);
};
