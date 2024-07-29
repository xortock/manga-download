package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xortock/mangafire-download/models"
	"github.com/xortock/mangafire-download/uris"
	"golang.org/x/net/html"
)

type IMangaClient interface {
}

type MangaClient struct {
	httpClient http.Client
}

func NewMangaClient() *MangaClient {
	return &MangaClient{
		httpClient: http.Client{},
	}
}

func (client *MangaClient) GetChapters(mangaCode string) []models.Chapter {

	var response, err = client.httpClient.Get(uris.MANGA_FIRE + "ajax/read/" + mangaCode + "/chapter/en")
	if err != nil {
		cli.Exit("Error: Failed to retrieve manga chapters", 1)
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		cli.Exit("Error: Failed to read request body", 1)
	}

	var chaptersApiResult models.ChaptersApiResult
	var errJson = json.Unmarshal(body, &chaptersApiResult)
	if errJson != nil {
		cli.Exit("Error: Failed to deserialize request body", 1)
	}

	document, err := html.Parse(strings.NewReader(chaptersApiResult.HtmlResult.HtmlContent))
	if err != nil {
		cli.Exit("Error: Failed to parse response to html", 1)
	}

	var matchingNodes = FindDescendants(document, "a")

	// check for matching attributes
	var chapters = []models.Chapter{}
	for _, node := range matchingNodes {

		var numberAsString string
		var idAsString string

		for _, attribute := range node.Attr {
			if attribute.Key == "data-number" {
				numberAsString = attribute.Val
			}
			if attribute.Key == "data-id" {
				idAsString = attribute.Val
			}
		}

		var number, _ = strconv.ParseFloat(numberAsString, 64)
		var id, _ = strconv.Atoi(idAsString)

		var chapter = models.Chapter{
			Number: number,
			Id:     id,
		}
		chapters = append(chapters, chapter)
	}

	return chapters
}

func FindDescendants(node *html.Node, tag string) []html.Node {
	var matchingNodes = []html.Node{}

	// The signature needs to be defined here otherwise it cant be use down below
	var traverse func(node *html.Node, tag string)
	traverse = func(node *html.Node, tag string) {
		if node.Type == html.ElementNode && node.Data == tag {
			// process the Product details within each <li> element
			matchingNodes = append(matchingNodes, *node)
		}
		// traverse the child nodes
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child, tag)
		}
	}

	traverse(node, tag)
	return matchingNodes
}

func (client *MangaClient) GetChapterUris(chapterId int) []string {

	var response, err = client.httpClient.Get(uris.MANGA_FIRE + "/ajax/read/chapter/" + strconv.Itoa(chapterId))
	if err != nil {
		cli.Exit("Error: Failed to retrieve manga chapters", 1)
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		cli.Exit("Error: Failed to read request body", 1)
	}

	var chapterApiResult models.ChapterApiResult
	var errJson = json.Unmarshal(body, &chapterApiResult)
	if errJson != nil {
		cli.Exit("Error: Failed to deserialize request body", 1)
	}

	var uris []string

	for _, item := range chapterApiResult.Result.Images {
		var items = item.([]interface{})
		uris = append(uris, items[0].(string))
	}

	return uris
}

func (client *MangaClient) GetChapterImages(uri string) []byte {

	var response, err = client.httpClient.Get(uri)
	if err != nil {
		cli.Exit("Error: Failed to retrieve manga chapters", 1)
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		cli.Exit("Error: Failed to read request body", 1)
	}

	return body
}
