export class BackendClient {
    constructor() {
        this.root = window.location.href;
        this.version = "v0/"
    }

    async ListMovies(handler) {
        getAndHandle(this.root + this.version + "movies", handler);
    }

    async PostMovie(movieName, handler) {
        const formData = new URLSearchParams();
        formData.append('name', movieName);
        postAndHandle(this.root + this.version + "movie", formData, handler);
    }
}

// helpers {
function getAndHandle(route, handler) {
  fetch(route).then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    }).then((data) => {
      if (typeof handler === "function") {
        handler(data);
      } else {
        console.log(data);
      }
    }).catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}

function postAndHandle(route, body, handler) {
  fetch(route, {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: body,
    }).then((response) => {
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.text(); // Change to .json() if the response is JSON
    }).then((data) => {
      if (typeof handler === "function") {
        handler(data);
      } else {
        console.log(data);
      }
    }).catch((error) => {
      console.error("There was a problem with the fetch operation:", error);
    });
}
// } helpers
