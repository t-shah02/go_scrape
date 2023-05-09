# Go Web Scraper
A web scraper written in Go, that fetches textual content from a single website domain concurrently using the Colly library. 

Libraries used
---
- https://github.com/gocolly/colly

Running the project locally 
---

Step 1: You will need to clone this repository somewhere on your machine, so you can run this command to do so:
```bash
git clone https://github.com/t-shah02/go_scrape.git
```

Step 2: Now that repository is cloned, you can enter the folder and run the make command to automatically build, and compile the source code.
```bash
make
```

Step 3: In order to run the generated binary, you need to go into the generated binary folder and run the scraper.
```bash
cd bin
./scraper -domain fireship.io -maxExplorationDepth 5 
```

Executable/binary flags 
---
Here is a brief list of each optional command-line argument that the scraper binary accepts:
- **-domain (string)**: This is the allowed domain that the scraper traverses and mines text data from. ex: example.com
- **-protocol (string)**: This is the protocol of the HTTP request that Colly will make on the web (it should be either http or https, based on your absolute URL)
- **-maxExplorationDepth (unsigned int)**: This is the maximum recursive depth that the Colly collector will scrape hyperlinks through, while traversing the pages for text data. Fluctuating this value may yield less/more data in the output JSON file, as potentially more page routes could be asynchronously visited. 
- **-outputFolderPath (string)**: This is the relative path to the folder that you want to save scraped data JSON file in. It will automatically create this folder at that path, if the directory doesn't currently exist
- **-tags (string)**: A comma seperated list of HTML tag/elements, which indicate the DOM nodes that the scraper will target, when making requests to each page associated with the given domain. 

Output Files
---
The program will output a JSON file, which maps page routes on the scraper's specified domain, to various DOM elements that correspond to potentially useful text. These files are in the generated output folder at the root level, are typically labelled in this convention: ```domain_name```_```unix_time_stamp```.json

Here is an example of the JSON fields that you will find in an output file:
```json
{
  "https://example.com/": [
      {
        "elementTag": "h1",
        "innerText": "Example.com is such an amazing site, with useful content"
      },
      {
        "elementTag": "p",
        "innerText": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Fermentum posuere urna nec tincidunt praesent semper feugiat nibh sed. Magna sit amet purus gravida quis blandit turpis cursus in."
      }
  ],
  "https://example.com/vistors/": [
    {
      "elementTag": "h1",
      "innerText": "Visitors of example.com so far"
    },
    {
      "elementTag": "span"
      "innerText": "Joe"
    },
    {
      "elementTag": "span",
      "innerText": "Bob"
    }
  ]
}
```

Issues and Discussion
---
Current project issues and tickets can be tracked here: https://github.com/t-shah02/go_scrape/issues

Feel free to contact me regarding any questions you have about how the program works, or any code snippets that need clarification. You can also fork this repo, and make a pull request if you have any elegant solutions/ideas you want to see in this project.

Future Goals
---
1. Make this web scraper a CLI command, so that it is more widely accessibly to users that don't have the GO compiler installed on their machines.
2. Allow this scraper to run on an endpoint, in a RESTful API, containerize it with Docker and finally deploy it someplace in the cloud :)

