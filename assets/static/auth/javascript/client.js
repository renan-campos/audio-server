export class BackendClient {
  constructor() {
    this.root = window.location.origin;
    this.version = "/v0"
  }

    async SendOtp(otp, respHandler, bodyHandler, errorHandler) {
        Post(this.root + this.version + "/token", otp, respHandler, bodyHandler, errorHandler);
    }

}

// helpers {
function Post(route, body, respHandler, textHandler, errorHandler) {
    fetch(route, 
        {method: "POST", headers: { "Content-Type": "application/text" }, body: body})
        .then(resp => { return respHandler(resp); })
        .then(text => { textHandler(text); })
        .catch(error => { errorHandler(error); });
}
// } helpers
