const PROXY_CONFIG = {
  "/pp": {
    "target": {
      "host": "stratos-ui.ngrok.io",
      "protocol": "https:",
      "port": 443
    },
    "secure": false,
    "changeOrigin": true
  }
}

module.exports = PROXY_CONFIG;
