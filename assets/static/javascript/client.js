class BackendClient {
  constructor() {
    this.root = window.location.href;
  }

  async ListAudio(handler) {
    fetchAndHandle(this.root + "audio", handler);
  }

  async GetAudioMetadata(id, handler) {
    fetchAndHandle(this.root + "audio/" + id, handler);
  }

  async GetAudioOgg(id, handler) {
    fetchAndHandle(this.root + "audio/" + id + "/ogg", handler);
  }

  async CreateAudio(auth, handler) {
    fetchAndHandleWithAuth(this.root + "admin/audio", auth, handler);
  }

  async UploadAudioMetadata(auth, metadata, handler) {
    fetchAndHandleWithAuth(
      this.root + "admin/audio" + id,
      auth,
      metadata,
      handler,
    );
  }

  async UploadAudio(auth, audiofile, handler) {
    const formData = new FormData();
    formData.append("audioFile", audioFile);
    fetchAndHandleWithAuth(
      this.root + "admin/audio" + id + "/ogg",
      auth,
      formData,
      handler,
    );
  }
}

// helpers {
function fetchAndHandle(route, handler) {
  fetch(route)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.text();
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

function fetchAndHandleWithAuth(route, auth, handler, body) {
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
