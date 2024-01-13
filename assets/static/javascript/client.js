export class BackendClient {
  constructor() {
    this.root = window.location.href;
    this.version = "v0/"
  }

  async ListAudio(handler) {
    fetchAndHandle(this.root + this.version + "audio", handler);
  }

  async GetAudioMetadata(id, handler) {
    fetchAndHandle(this.root + this.version + "audio/" + id, handler);
  }

  async GetAudioOgg(id, handler) {
    fetchAndHandleBlob(this.root + this.version + "audio/" + id + "/ogg", handler);
  }

  async CreateAudio(auth, handler) {
    fetchAndHandleWithAuth(this.root + this.version + "admin/audio", auth, handler);
  }

  async UploadAudioMetadata(id, auth, metadata, handler) {
    fetchAndHandleWithAuth(
      this.root + this.version + "admin/audio/" + id,
      auth,
      metadata,
      handler,
    );
  }

  async UploadAudio(id, auth, audioFile, handler) {
    const formData = new FormData();
    formData.append("audioFile", audioFile);
  const base64Credentials = btoa(auth.username + ":" + auth.password);

  fetch(
    this.root + this.version + "admin/audio/" + id + "/ogg", {
    method: "POST",
    headers: {
      Authorization: "Basic " + base64Credentials,
    },
    body: formData,
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.text(); // Change to .json() if the response is JSON
    })
    .then((data) => {
      if (typeof handler === "function") {
        handler(data);
      } else {
        console.log(data);
      }
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
      reject(error); // Reject the Promise if there's an error
    });
  }
}

// helpers {
function fetchAndHandle(route, handler) {
  fetch(route)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    })
    .then((data) => {
      if (typeof handler === "function") {
        handler(data);
      } else {
        console.log(data);
      }
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function fetchAndHandleBlob(route, handler) {
  fetch(route)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.blob();
    })
    .then((blob) => {
      if (typeof handler === "function") {
        handler(blob);
      } else {
        console.log(blob);
      }
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function fetchAndHandleWithAuth(route, auth, body, handler) {
  const base64Credentials = btoa(auth.username + ":" + auth.password);

  fetch(route, {
    method: "POST",
    headers: {
      Authorization: "Basic " + base64Credentials,
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: body,
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.text(); // Change to .json() if the response is JSON
    })
    .then((data) => {
      if (typeof handler === "function") {
        handler(data);
      } else {
        console.log(data);
      }
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
      reject(error); // Reject the Promise if there's an error
    });
}

function extractUuidFromText(text) {
  // Regular expression to match a UUID pattern
  const uuidRegex =
    /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/;

  // Find the first match in the text
  const match = text.match(uuidRegex);

  // Check if a match was found
  if (match && match.length > 0) {
    return match[0]; // The first match is the UUID
  } else {
    return null; // No UUID found in the text
  }
}
// } helpers
