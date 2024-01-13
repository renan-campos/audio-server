import {BackendClient} from './javascript/client.js'

// Define a UI structure (could be an object or a class)
class UI {
    constructor(backendClient) {
        // TODO move this sensitive, secretive data
        this.auth = {username: "rcampos", password: "relax"}
        this.backendClient = backendClient;

        this.Player = document.getElementById('audioPlayer');

        // List of guitar sessions
        // - Has play button, enabled only if a sound file is present
        // - Has trash button
        // - Has editable name
        // - Has button for uploading ogg file
        this.SessionList = document.createElement('ul');
        this.SessionList.id = 'SessionList';

        // Create a button element and configure it
        this.CreateSessionButton = createButton(
            'CreateSessionButton',
            'Create Session',
            () => {
            this.backendClient.CreateAudio(this.auth, (data) => {
                console.log("Did it!");
                console.log(data);
            });
            const newItem = document.createElement('li');
            this.refreshSessionList();
        });
        this.RefreshSessionListButton = createButton(
            'RefreshSessionListButton',
            'Refresh List',
            () => {
            this.clearSessionList();
            this.populateSessionList();
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
            this.SessionList, 
        ]) {
            document.body.appendChild(element)
        }
    }

    clearSessionList() {
        var sessionList = document.getElementById('SessionList');
        // Remove all child elements (list items)
        while (this.SessionList.firstChild) {
            this.SessionList.removeChild(this.SessionList.firstChild);
        }
    }

    populateSessionList() {
        this.backendClient.ListAudio((data) => {
            data.Items.map((sessionId) => {
                    this.backendClient.GetAudioMetadata(sessionId, (metadata) => {
                    // Play button {
                    const playButton = document.createElement('button');
                    playButton.textContent = '▶️';
                    playButton.className = 'edit-button';
                    playButton.addEventListener('click', () => {
                        console.log('Play button clicked for ' + sessionId);
                        this.backendClient.GetAudioOgg(sessionId, (blob) => {
                            //Create a blob URL and set it as the audio source
                            const blobURL = URL.createObjectURL(blob);
                            this.Player.type = "audio/ogg";
                            this.Player.src = blobURL;

                            // Play the audio
                            this.Player.play();
                        })

                    });
                    // } Play Button
                    // Edit button {
                    const editButton = document.createElement('button');
                    editButton.textContent = '✏️';
                    editButton.className = 'edit-button';
                    editButton.addEventListener('click', () => {
                        console.log('Edit button clicked for ' + sessionId);

                        const sessionName = prompt("Enter session name:");
                        console.log(sessionName);
                        const formData = new URLSearchParams();
                        formData.append('name', sessionName);
                        console.log(formData.name)
                        if (sessionName !== null && sessionName !== "") {
                            this.backendClient.UploadAudioMetadata(
                                sessionId,
                                this.auth, 
                                formData,
                                (data) => { console.log(data); });
                        }
                        this.refreshSessionList();
                    });
                    // } Edit Button
                    // Delete button {
                    const deleteButton = document.createElement('button');
                    deleteButton.textContent = '🗑️';
                    deleteButton.className = 'delete-button';
                    deleteButton.addEventListener('click', () => {
                        console.log('[Unimplemented] Delete button clicked for ' + sessionId);

                    });
                    // } Delete Button
                    // Cloud button {
                    const cloudButton = document.createElement('button');
                    cloudButton.textContent = '☁️';
                    cloudButton.className = 'cloud-button';
                    cloudButton.addEventListener('click', () => {
                        console.log('[Unimplemented] Cloud button clicked for ' + sessionId);

                    });
                    // } Delete Button
                    // Upload button {
                    const uploadButton = document.createElement('button');
                    uploadButton.textContent = '📂';
                    uploadButton.className = 'upload-button';
                    uploadButton.addEventListener('click', () => {
                        console.log('Upload button clicked for ' + sessionId);

                        // Create the file input element
                        const fileInput = document.createElement('input');
                        fileInput.type = 'file';
                        fileInput.id = 'fileInput';
                        fileInput.style.display = 'none';

                        // Attach an event listener to the file input element
                        fileInput.addEventListener('change', (event) =>{
                            const selectedFile = event.target.files[0];
                            if (selectedFile) {
                                // Handle the selected file here
                                console.log('Selected file:', selectedFile.name);
                                this.backendClient.UploadAudio(
                                    sessionId, 
                                    this.auth,
                                    selectedFile);
                            }
                        });
                        fileInput.click();
                    });
                    // } Upload Button
                    // Session Div {
                    const sessionDiv = document.createElement('div');
                    sessionDiv.classList.add('session-div');
                    sessionDiv.textContent = metadata.Name;
                    // } Session Div
                    // Admin Button Div {
                    const adminDiv = document.createElement('div');
                    adminDiv.classList.add('admin-div');
                    adminDiv.appendChild(editButton);
                    adminDiv.appendChild(uploadButton);
                    // To add at a later point
                    //adminDiv.appendChild(cloudButton);
                    //adminDiv.appendChild(deleteButton);
                    // } Admin Button Div 
                    // User Button Div {
                    const userDiv = document.createElement('div');
                    userDiv.classList.add('user-div');
                    userDiv.appendChild(playButton);
                    // } User Button Div
                    // Append the edit button to the list item
                    const listItem = document.createElement('li');
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
    var button = document.createElement('button');
    button.id = id
    button.textContent = textContent;
    button.addEventListener('click', clickHandler)
    return button
}
