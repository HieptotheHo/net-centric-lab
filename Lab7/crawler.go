package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
func HTMLParser(resp *http.Response) *html.Node {
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// Parse the HTML content
	content, err := html.Parse(bytes.NewReader(body))
	checkError(err)

	return content
}

func main() {
	mangaByGenre := make(map[string][]string)
	// Make a GET request to the WEBTOON homepage
	resp, err := http.Get("https://www.webtoons.com/en/")
	checkError(err)

	defer resp.Body.Close()

	doc := HTMLParser(resp)

	// Find all the links on the page
	// links := findAllLinks(doc)

	// // Print the links
	// for _, link := range links {
	// 	fmt.Println(link)
	// }
	linkToGenres := findLinkWithText(doc, "GENRES")

	if linkToGenres != "" {
		fmt.Println(linkToGenres)
		genresRsp, err := http.Get("https://www.webtoons.com/en/genres")
		checkError(err)

		defer genresRsp.Body.Close()

		genresDoc := HTMLParser(genresRsp)

		genres, genreNames := findAllGenreLinks(genresDoc)

		numberOfGenres := 10
		numberOfTitles := 10
		for i := 0; i < numberOfGenres; i++ {

			genreRsp, err := http.Get(genres[i])
			checkError(err)

			defer genreRsp.Body.Close()

			genreDoc := HTMLParser(genreRsp)

			titles := findAllTitlesOfAGenre(genreDoc)

			fmt.Println("############################################")
			fmt.Println(genreNames[i])

			mangaByGenre[genreNames[i]] = titles[0:numberOfTitles]

			for j := 0; j < numberOfTitles; j++ {
				fmt.Println(j+1, ". ", titles[j])

			}
		}
	}
	// Encode the data to JSON
	jsonData, err := json.MarshalIndent(mangaByGenre, "", "  ")
	checkError(err)

	// Create the JSON file
	file, err := os.Create("manga_by_genre.json")
	checkError(err)
	defer file.Close()

	// Write the JSON data to the file
	if _, err := file.Write(jsonData); err != nil {
		checkError(err)
	}

	fmt.Println("Data successfully written to manga_by_genre.json")
}

func findLinkWithText(n *html.Node, text string) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode && c.Data == text {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						return attr.Val
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findLinkWithText(c, text); result != "" {
			return result
		}
	}
	return ""
}

func findAllGenreLinks(n *html.Node) ([]string, []string) {
	var genres []string
	var genreNames []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "data-genre" {
					// Found an li with data-genre attribute
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "a" {
							var link string
							var name string
							// Extract the href attribute and text content of the <a> tag
							for _, aAttr := range c.Attr {
								if aAttr.Key == "href" {
									link = aAttr.Val
								}
							}
							for aChild := c.FirstChild; aChild != nil; aChild = aChild.NextSibling {
								if aChild.Type == html.TextNode {
									// name = aChild.Data
									name = strings.ReplaceAll(aChild.Data, " ", "")
									name = strings.ReplaceAll(name, "\t", "")
									name = strings.ReplaceAll(name, "\n", "")
								}
							}
							if link != "" && name != "" {
								genres = append(genres, link)
								genreNames = append(genreNames, name)
							}
						}
					}
				}
			}
		}
		// Recursively traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return genres, genreNames
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractText(c)
	}
	return result
}

func findAllTitlesOfAGenre(n *html.Node) []string {

	// Initialize a slice to hold the contents
	var titles []string

	// Function to traverse the nodes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "subj" {
					// Found a p with class subj, extract its text content
					titles = append(titles, extractText(n))
					break
				}
			}
		}
		// Recursively traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return titles
}

func findAllLinks(n *html.Node) []string {
	var links []string

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Println(attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return links
}
