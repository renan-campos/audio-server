root = window.location.href

function listAudio() {
    fetch(root + 'audio')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.text();
      })
      .then(data => {
        document.writeln(data + "<br>");
      })
      .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
      });
}

function listAudio(id) {
    fetch(root + 'audio/' + id)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.text();
      })
      .then(data => {
        document.writeln(data + "<br>");
      })
      .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
      });
}

function getAudioOgg(id) {
    fetch(root + 'audio/' + id + '/ogg')
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.text();
      })
      .then(data => {
        obj = data; document.writeln(data);
      })
      .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
      });
}

async function createAudioPromise(username, password) {
  return new Promise((resolve, reject) => {
    // Encode the username and password for Basic Authentication
    const base64Credentials = btoa(username + ':' + password);

    fetch(root + 'admin/audio', {
      method: 'POST',
      headers: {
        'Authorization': 'Basic ' + base64Credentials,
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.text(); // Change to .json() if the response is JSON
      })
      .then(data => {
        // Handle the response data here
        const audio_uuid = extractUuidFromText(data);
        resolve(audio_uuid); // Resolve the Promise with the audio_uuid
        console.log('new audio uuid: ', audio_uuid)
      })
      .catch(error => {
        console.error('There was a problem with the fetch operation:', error);
        reject(error); // Reject the Promise if there's an error
      });
  });
}

function createAudioMetadata(username, password, id) {
    // Encode the username and password for Basic Authentication
    const base64Credentials = btoa(username + ':' + password);

    var audio_uuid;

    fetch(root + 'admin/audio/' + id, {
      method: 'POST',
      headers: {
        'Authorization': 'Basic ' + base64Credentials,
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.text(); // Change to .json() if the response is JSON
  })
  .then(data => {
    // Handle the response data here
    console.log(data);
    audio_uuid = extractUuidFromText(data);
      console.log(audio_uuid);
  })
  .catch(error => {
    console.error('There was a problem with the fetch operation:', error);
  });
}

function uploadAudioOgg(username, password, id) {
    // Encode the username and password for Basic Authentication
    const base64Credentials = btoa(username + ':' + password);

    var audio_uuid;

    fetch(root + 'admin/audio/' + id + '/ogg', {
      method: 'POST',
      headers: {
        'Authorization': 'Basic ' + base64Credentials,
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.text(); // Change to .json() if the response is JSON
  })
  .then(data => {
    // Handle the response data here
    console.log(data);
    audio_uuid = extractUuidFromText(data);
      console.log(audio_uuid);
  })
  .catch(error => {
    console.error('There was a problem with the fetch operation:', error);
  });
}

function extractUuidFromText(text) {
  // Regular expression to match a UUID pattern
  const uuidRegex = /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/;

  // Find the first match in the text
  const match = text.match(uuidRegex);

  // Check if a match was found
  if (match && match.length > 0) {
    return match[0]; // The first match is the UUID
  } else {
    return null; // No UUID found in the text
  }
}
