export class BackendClient {
  constructor() {
    this.root = window.location.origin;
    this.version = "/v0/";
  }

  async ListAudio(handler) {
    fetchAndHandle(this.root + this.version + "audio", handler);
  }

  async GetAudioMetadata(id, handler) {
    fetchAndHandle(this.root + this.version + "audio/" + id, handler);
  }

  async GetAudio(id, handler) {
    fetchAndHandleBlob(
      this.root + this.version + "audio/" + id + "/webm",
      handler,
    );
  }

  async GetAudioBuffer(id, start, end, handler) {
    return fetchAndHandleBuffer(
      this.root + this.version + "audio/" + id + "/webm",
      start,
      end,
      handler,
    );
  }

  async GetAudioHeader(id, handler) {
    return fetchAndHandleHeader(
      this.root + this.version + "audio/" + id + "/webm",
      handler,
    );
  }

  async CreateAudio(auth, handler) {
    fetchAndHandleWithAuth(
      this.root + this.version + "admin/audio",
      auth,
      handler,
    );
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

    fetch(this.root + this.version + "admin/audio/" + id + "/webm", {
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

function fetchAndHandleHeader(route, handler) {
  return fetch(route, { method: "HEAD" })
    .then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.headers;
    })
    .then((header) => {
      if (typeof handler === "function") {
        return handler(header);
      } else {
        console.log(header);
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

function fetchAndHandleBuffer(route, start, end, handler) {
  return fetch(route, { headers: { Range: `bytes=${start}-${end}` } })
    .then(async (response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.arrayBuffer();
    })
    .then((buffer) => {
      if (typeof handler === "function") {
        return handler(buffer);
      } else {
        console.log(buffer);
      }
    })
    .catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

// Function to retrieve a cookie by name
function getCookie(name) {
  const cookies = document.cookie.split(";");
  for (let cookie of cookies) {
    const [cookieName, cookieValue] = cookie.trim().split("=");
    if (cookieName === name) {
      return decodeURIComponent(cookieValue);
    }
  }
  return null;
}

function fetchAndHandleWithAuth(route, auth, body, handler) {
  const jwt = getCookie("jwt");
  if (!jwt) {
    console.log(
      "JWT token not found. Redirecting to authentication endpoint...",
    );
    // Token not found, redirect to authentication endpoint
    window.location.href =
      "/auth?redirect=" + encodeURIComponent(window.location.href);
  }

  fetch(route, {
    method: "POST",
    headers: {
      Authorization: "Bearer " + jwt,
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
