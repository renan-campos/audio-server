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

function createAudio() {
}
function uploadAudioMetadata(id) {
}
function uploadAudioOgg(id) {
}
