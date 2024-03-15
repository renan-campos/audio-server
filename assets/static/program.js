import {BackendClient} from './javascript/client.js'

// Define a UI structure (could be an object or a class)
class UI {
    constructor(backendClient) {
        this.backendClient = backendClient;

        this.MovieDisplay = document.createElement("h1");
        this.MovieDisplay.id = "DisplayMovie";
        this.MovieDisplay.textContent = "";

        // Create a button element and configure it
        this.AddMovieButton = createButton(
            'AddMovieButton',
            'Add Movie',
            () => {
                const movieName = prompt("Enter a movie name:");
                if (movieName !== null) {
                    this.backendClient.PostMovie(movieName, (movieName) => {
                        console.log("Movie posted.");
                        console.log(movieName);
                    });
                }
            });
        this.PickMovieButton = createButton(
            'PickMovieButton',
            'Pick a Movie',
            () => {
                this.backendClient.ListMovies((movieList) => {
                    const movie  = getRandomElement(movieList.Items);
                    this.MovieDisplay.textContent = movie;
                    console.log("Movie list attained.");
                    console.log(movieList);
                    console.log("Random movie chosen: " + movie);
                });
        });
    }

    Render() {
        // Append the list and button to the document
        for (const element of [
            this.MovieDisplay,
            this.AddMovieButton,
            this.PickMovieButton,
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

function getRandomElement(array) {
  // Generate a random index within the length of the array
  const randomIndex = Math.floor(Math.random() * array.length);
  // Return the element at that random index
  return array[randomIndex];
} 
