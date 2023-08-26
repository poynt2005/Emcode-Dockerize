var fs = require('fs');
var path = require('path')
var os = require('os');
var cp = require("child_process");


var configYaml = [
    "bind-addr: 0.0.0.0:8080",
    "auth: password",
    "password: " + process.env.CODE_PASSWORD,
    "cert: false"
].join("\n");

var configStoreDirectory = path.join(os.homedir(), '.config', 'code-server');
if (!fs.existsSync(configStoreDirectory)) {
    fs.mkdirSync(configStoreDirectory, { recursive: true });
}

fs.writeFileSync(path.join(configStoreDirectory, "config.yaml"), configYaml, "utf-8");

var child = cp.exec('code-server');
child.stdout.pipe(process.stdout);
child.on('exit', function () {
    process.exit();
});

