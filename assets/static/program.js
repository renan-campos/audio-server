function renderFullScreenImage(imageSrc) {
  // Create an image element
  var img = new Image();

  // Set the source of the image
  img.src = imageSrc;

  // Set the width and height of the image to the dimensions of the screen
  img.width = window.innerWidth;
  img.height = window.innerHeight;

  // Append the image to the body or any other container element
  document.body.appendChild(img);
}

// Example usage:
// Replace 'your_image_path.jpg' with the actual path to your image
renderFullScreenImage('success.png');

