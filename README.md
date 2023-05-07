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

Step 2: Once the project is cloned, change your directory into the src folder, and build the Go project and save the scraper binary in the project's root folder
```bash
cd go_scraper/src
go build -o ../scraper
cd ..
./scraper -domain "YOUR_WEBSITE_DOMAIN" -protocol "https | http" -maxExplorationDepth 5
```

Executable/binary flags 
---
Here is a brief list of each optional command-line argument that the scraper binary accepts:
- **-domain (string)**: This is the allowed domain that the scraper traverses and mines text data from. ex: example.com
- **-protocol (string)**: This is the protocol of the HTTP request that Colly will make on the web (it should be either http or https, based on your absolute URL)
- **-maxExplorationDepth (unsigned int)**: This is the maximum recursive depth that the Colly collector will scrape hyperlinks through, while traversing the pages for text data. Fluctuating this value may yield less/more data in the output JSON file, as potentially more page routes could be asynchronously visited. 

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
