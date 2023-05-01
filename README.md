
# bigJSON

  

bigJSON is a Go project designed to handle large JSON files effectively. It processes JSON data with high efficiency using Go's inherent concurrency capabilities.

  

## JSON File Format

  

The application expects a JSON file in the root directory in the following format:

  


`[ { "postcode": "10224", "recipe": "Creamy Dill Chicken", "delivery": "Wednesday 1AM - 7PM" }, ... ]`

  

## Function

The program is able to:

  

Count the number of unique recipe names.

Count the number of occurrences for each unique recipe name (alphabetically ordered by recipe name).

Find the postcode with most delivered recipes.

Count the number of deliveries to postcode 10120 that lie within the delivery time between 10AM and 3PM, examples (12AM denotes midnight):

NO - 9AM - 2PM

YES - 10AM - 2PM

List the recipe names (alphabetically ordered) that contain in their name one of the provided words, like:

Potato

Veggie

Mushroom

  
  

## Building Docker Image

To build the Docker image for the application, use the following command:

  

`docker build -t bigjson .`

  
  

## Running the Program Inside Docker Container

To execute the program inside the Docker container, use the following command:

  

`docker run -it --rm -v "$(pwd):/app" bigjson ./main -file /app/bigJSON.json -keywords "Tex-Mex Tilapia" -postcode 10213 -time "8PM-7AM"`

The flags used in the above command are:

  

-file: Specifies the path to the JSON file.

-keywords: Specifies the keywords to search in the JSON file.

-postcode: Specifies the postcode to search in the JSON file.

-time: Specifies the time window.

Ensure you replace the values in the flags with those relevant to your use case.
