import { BackendClient } from "./javascript/client.js";

// Define a UI structure (could be an object or a class)
class UI {
  constructor(backendClient) {
    this.backendClient = backendClient;

    this.Player = document.getElementById("audioPlayer");

    // List of guitar sessions
    // - Has play button, enabled only if a sound file is present
    // - Has trash button
    // - Has editable name
    // - Has button for uploading ogg file
    this.SessionList = document.createElement("ul");
    this.SessionList.id = "SessionList";

    // Create a button element and configure it
    this.CreateSessionButton = createButton(
      "CreateSessionButton",
      "Create Session",
      () => {
        this.backendClient.CreateAudio((data) => {
          console.log("Did it!");
          console.log(data);
        });
        const newItem = document.createElement("li");
        this.refreshSessionList();
      },
    );
    this.RefreshSessionListButton = createButton(
      "RefreshSessionListButton",
      "Refresh List",
      () => {
        this.clearSessionList();
        this.populateSessionList();
      },
    );
    // Create a button element and configure it
    this.LogoutButton = createButton("LogoutButton", "Logout", () => {
      deleteCookie("jwt"); // Delete the 'jwt' cookie
    });

    this.clearSessionList();
    this.populateSessionList();
  }

  Render() {
    // Append the list and button to the document
    for (const element of [
      this.Player,
      this.CreateSessionButton,
      this.RefreshSessionListButton,
      this.LogoutButton,
      this.SessionList,
    ]) {
      document.body.appendChild(element);
    }
  }

  clearSessionList() {
    var sessionList = document.getElementById("SessionList");
    // Remove all child elements (list items)
    while (this.SessionList.firstChild) {
      this.SessionList.removeChild(this.SessionList.firstChild);
    }
  }

  populateSessionList() {
              const backendClient = this.backendClient;
    this.backendClient.ListAudio((data) => {
      data.Items.map((sessionId) => {
        this.backendClient.GetAudioMetadata(sessionId, (metadata) => {
          // Play button {
          const playButton = document.createElement("button");
          playButton.textContent = "▶️";
          playButton.className = "edit-button";
          playButton.addEventListener("click", async () => {
            console.log("Play button clicked for " + sessionId);

            const audioData = await this.backendClient.GetAudioHeader(
              sessionId,
              (header) => {
                return {
                  chunkSize: parseInt(header.get("X-Chunk-Size")),
                  contentType: header.get("Content-Type"),
                  fileSize: parseInt(header.get("Content-Length")),
                };
              },
            );

            var mediaSource = new MediaSource();
            var audio = document.querySelector("audio");
            audio.src = URL.createObjectURL(mediaSource);
            mediaSource.addEventListener(
              "sourceopen",
              async function () {
                // Step 4: Add a SourceBuffer and append media segments
                var sourceBuffer = mediaSource.addSourceBuffer(
                  audioData.contentType,
                );

                // Function to fetch and append the next segment
                var segmentStart = 0;
                async function appendNextSegment() {
                  if (segmentStart > audioData.fileSize) {
                    console.log("All segments appended");
                    mediaSource.endOfStream();
                    return;
                  }
                  const segmentEnd = Math.min(
                    segmentStart + audioData.chunkSize,
                    audioData.fileSize,
                  );

                  const newBuffer = backendClient.GetAudioBuffer(
                    sessionId,
                    segmentStart,
                    segmentEnd,
                    (buffer) => {
                      sourceBuffer.appendBuffer(buffer);
                      segmentStart += audioData.chunkSize;
                      return sourceBuffer;
                    },
                  );
                  return newBuffer;
                }

                const firstUpdateHandler = function () {
                  sourceBuffer.removeEventListener(
                    "updateend",
                    firstUpdateHandler,
                  );

                  var downloadTime = 0;
                  var cleanupTime = sourceBuffer.buffered.end(0);

                  const audioTimeHandler = function () {
                    audio.removeEventListener("timeupdate", audioTimeHandler);

                    if (audio.currentTime > downloadTime) {
                      downloadTime = sourceBuffer.buffered.end(0);
                      appendNextSegment();
                    }
                    if (audio.currentTime > cleanupTime) {
                      sourceBuffer.remove(
                        sourceBuffer.buffered.start(0),
                        cleanupTime,
                      );
                      cleanupTime = sourceBuffer.buffered.end(0);
                    }
                    setTimeout(
                      function () {
                        audio.addEventListener("timeupdate", audioTimeHandler);
                      },
                      (downloadTime - audio.currentTime) * 1000,
                    );
                  };
                  audio.addEventListener("timeupdate", audioTimeHandler);
                };
                sourceBuffer.addEventListener("updateend", firstUpdateHandler);

                // Append the first segment
                appendNextSegment();

                // Start playing the audio after the first segment is appended
                audio.play();
              },
              false,
            );
          });
          // } Play Button
          // Edit button {
          const editButton = document.createElement("button");
          editButton.textContent = "✏️";
          editButton.className = "edit-button";
          editButton.addEventListener("click", () => {
            console.log("Edit button clicked for " + sessionId);

            const sessionName = prompt("Enter session name:");
            console.log(sessionName);
            const formData = new URLSearchParams();
            formData.append("name", sessionName);
            console.log(formData.name);
            if (sessionName !== null && sessionName !== "") {
              this.backendClient.UploadAudioMetadata(
                sessionId,
                formData,
                (data) => {
                  console.log(data);
                },
              );
            }
            this.refreshSessionList();
          });
          // } Edit Button
          // Delete button {
          const deleteButton = document.createElement("button");
          deleteButton.textContent = "🗑️";
          deleteButton.className = "delete-button";
          deleteButton.addEventListener("click", () => {
            console.log(
              "[Unimplemented] Delete button clicked for " + sessionId,
            );
          });
          // } Delete Button
          // Cloud button {
          const cloudButton = document.createElement("button");
          cloudButton.textContent = "☁️";
          cloudButton.className = "cloud-button";
          cloudButton.addEventListener("click", () => {
            console.log(
              "[Unimplemented] Cloud button clicked for " + sessionId,
            );
          });
          // } Delete Button
          // Upload button {
          const uploadButton = document.createElement("button");
          uploadButton.textContent = "📂";
          uploadButton.className = "upload-button";
          uploadButton.addEventListener("click", () => {
            console.log("Upload button clicked for " + sessionId);

            // Create the file input element
            const fileInput = document.createElement("input");
            fileInput.type = "file";
            fileInput.id = "fileInput";
            fileInput.style.display = "none";

            // Attach an event listener to the file input element
            fileInput.addEventListener("change", (event) => {
              const selectedFile = event.target.files[0];
              if (selectedFile) {
                // Handle the selected file here
                console.log("Selected file:", selectedFile.name);
                this.backendClient.UploadAudio(sessionId, selectedFile);
              }
            });
            fileInput.click();
          });
          // } Upload Button
          // Session Div {
          const sessionDiv = document.createElement("div");
          sessionDiv.classList.add("session-div");
          sessionDiv.textContent = metadata.Name;
          // } Session Div
          // Admin Button Div {
          const adminDiv = document.createElement("div");
          adminDiv.classList.add("admin-div");
          adminDiv.appendChild(editButton);
          adminDiv.appendChild(uploadButton);
          // To add at a later point
          //adminDiv.appendChild(cloudButton);
          //adminDiv.appendChild(deleteButton);
          // } Admin Button Div
          // User Button Div {
          const userDiv = document.createElement("div");
          userDiv.classList.add("user-div");
          userDiv.appendChild(playButton);
          // } User Button Div
          // Append the edit button to the list item
          const listItem = document.createElement("li");
          listItem.appendChild(sessionDiv);
          listItem.appendChild(adminDiv);
          listItem.appendChild(userDiv);
          this.SessionList.appendChild(listItem);
        });
      });
    });
  }

  refreshSessionList() {
    this.clearSessionList();
    this.populateSessionList();
  }
}

// Create an instance of the UI structure
const ui = new UI(new BackendClient());
ui.Render();

function createButton(id, textContent, clickHandler) {
  var button = document.createElement("button");
  button.id = id;
  button.textContent = textContent;
  button.addEventListener("click", clickHandler);
  return button;
}

function deleteCookie(name) {
  document.cookie = name + "=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}
