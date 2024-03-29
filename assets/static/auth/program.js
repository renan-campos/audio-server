import {BackendClient} from './javascript/client.js'

// Define a UI structure (could be an object or a class)
class UI {
    constructor(backendClient) {
        this.backendClient = backendClient;

        this.KeyInput = createTextBox("KeyInput", "Enter key");
        var keyInput = this.KeyInput
        var listener = function() {
            keyInput.value = "";
            keyInput.removeEventListener("focus", listener)
        };
        this.KeyInput.addEventListener("focus", listener);

        this.SendButton = createButton(
            'SendButton',
            'Send',
            () => {
               this.backendClient.SendOtp(this.KeyInput.value, (resp => {
                   if (!resp.ok) {
                       this.KeyInput.value = "Invalid OTP, please try again";
                       throw new Error('Network response was not ok');
                   } else {
                       this.KeyInput.value = "OTP Successful!";
                       return resp.text();
                   }
               }),
                (text => {
                    console.log(text);
                    setCookie("jwt", text, 30);
                    // Function to parse URL parameters

                // Check if the redirect parameter exists in the URL
                const redirectUrl = getParameterByName('redirect');
                if (redirectUrl) {
                  // Redirect the user to the specified URL
                  window.location.href = redirectUrl;
                }

                }),
                (error => {
                    console.error('Checking OTP:', error);
                }));
               this.KeyInput.addEventListener("focus", listener);
            });
    }

    Render() {
        for (const element of [
            this.KeyInput,
            this.SendButton
        ]) {
            document.body.appendChild(element)
        }
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

function createTextBox(textBoxId, defaultText) {
    var textBox = document.createElement("input");
    textBox.setAttribute("type", "text");
    textBox.setAttribute("id", textBoxId);
    textBox.setAttribute("value", defaultText);

    return textBox
}

// Function to set a cookie
function setCookie(name, value, expirationDays) {
    const date = new Date();
    date.setTime(date.getTime() + (expirationDays * 24 * 60 * 60 * 1000)); // Convert expiration days to milliseconds
    const expires = "expires=" + date.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

function getParameterByName(name, url) {
  if (!url) url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}
