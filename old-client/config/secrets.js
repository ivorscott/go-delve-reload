// Docker secrets helper for SPAs

const fs = require("fs");
const path = require("path");

exports = {};

const SECRET_DIR = "/run/secrets";

function getSecrets(secretDir) {
  const secrets = {};
  if (fs.existsSync(secretDir)) {
    const files = fs.readdirSync(secretDir);

    files.forEach(file => {
      const fullPath = path.join(secretDir, file);
      const key = file;
      const data = fs
        .readFileSync(fullPath, "utf8")
        .toString()
        .trim();

      secrets[key] = data;
    });
  }
  return secrets;
}

function getSecretFactory(secrets) {
  return name => secrets[name];
}

function isSecret(key) {
  return process.env[key].includes(SECRET_DIR);
}

const secrets = getSecrets(SECRET_DIR);

exports.secrets = secrets;
exports.isSecret = isSecret;
exports.getSecret = getSecretFactory(secrets);
exports.getSecrets = getSecrets;

module.exports = exports;
