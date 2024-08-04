package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/xortock/mangafire-download/constants"
	"github.com/xortock/mangafire-download/extensions"
	"github.com/xortock/mangafire-download/models"
	"golang.org/x/net/html"
)

type IMangaClient interface {
	GetDivisionsUris(division string, divisionId int) ([]string, error)
	GetDivisions(division string, mangaCode string) ([]models.Division, error)
	GetDivisionImages(uri string) ([]byte, error)
}

type MangaClient struct {
	httpClient  *http.Client
	baseAddress string
}

func NewMangaClient() *MangaClient {
	return &MangaClient{
		httpClient:  &http.Client{},
		baseAddress: constants.URI_MANGA_FIRE,
	}
}

func (client *MangaClient) GetDivisions(division string, mangaCode string) ([]models.Division, error) {
	var divisions []models.Division

	switch division {
	case constants.DIVISION_CHAPTER:
		var tempChapters, err = client.GetDivisionsInternal(mangaCode, constants.DIVISION_CHAPTER)
		if err != nil {
			return nil, err
		}
		divisions = tempChapters
	case constants.DVISION_VOLUME:
		var tempVolumes, err = client.GetDivisionsInternal(mangaCode, constants.DVISION_VOLUME)
		if err != nil {
			return nil, err
		}
		divisions = tempVolumes
	default:
		return nil, errors.New(division + " division type not supported")
	}

	slices.Reverse(divisions)
	return divisions, nil
}

func (client *MangaClient) GetDivisionsInternal(mangaCode string, division string) ([]models.Division, error) {
	var url, errUrl = url.JoinPath(client.baseAddress, "ajax/read", mangaCode, division, "en")
	if errUrl != nil {
		return nil, errUrl
	}

	var request, errRequest = http.NewRequest("GET", url, nil)
	if errRequest != nil {
		return []models.Division{}, errors.New("error: failed to create http request")
	}

	request.Header = http.Header{"User-Agent": {constants.MOZILLA_USER_AGENT}}

	var response, err = client.httpClient.Do(request)
	var isSuccessStatusCode = extensions.IsSuccessStatusCode(response)
	if err != nil || !isSuccessStatusCode {
		return []models.Division{}, errors.New(fmt.Sprint("error: failed to retrieve manga chapters status code: ", strconv.Itoa(response.StatusCode), " ", http.StatusText(response.StatusCode)))
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		return []models.Division{}, errors.New("error: failed to read request body")
	}

	var chaptersApiResult models.ChaptersApiResult
	var errJson = json.Unmarshal(body, &chaptersApiResult)
	if errJson != nil {
		return []models.Division{}, errors.New("error: failed to deserialize request body")
	}

	document, err := html.Parse(strings.NewReader(chaptersApiResult.HtmlResult.HtmlContent))
	if err != nil {
		return []models.Division{}, errors.New("error: failed to parse response to html")
	}

	var matchingNodes = FindDescendants(document, "a")

	// check for matching attributes
	var chapters = []models.Division{}
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

		var chapter = models.Division{
			Number: number,
			Id:     id,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
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

func (client *MangaClient) GetDivisionsUris(division string, divisionId int) ([]string, error) {
	var divisionsUris []string

	switch division {
	case constants.DIVISION_CHAPTER:
		var chapterUris, err = client.GetDivisionUrisInternal(divisionId, constants.DIVISION_CHAPTER)
		if err != nil {
			return nil, err
		}
		divisionsUris = chapterUris
	case constants.DVISION_VOLUME:
		var chapterUris, err = client.GetDivisionUrisInternal(divisionId, constants.DVISION_VOLUME)
		if err != nil {
			return nil, err
		}
		divisionsUris = chapterUris
	default:
		return nil, errors.New(division + " division type not supported")
	}

	return divisionsUris, nil
}

func (client *MangaClient) GetDivisionUrisInternal(chapterId int, division string) ([]string, error) {
	var url, errUrl = url.JoinPath(client.baseAddress, "ajax/read", division, strconv.Itoa(chapterId))
	if errUrl != nil {
		return nil, errUrl
	}

	var request, errRequest = http.NewRequest("GET", url, nil)
	if errRequest != nil {
		return []string{}, errors.New("error: failed to create http request")
	}

	request.Header = http.Header{"User-Agent": {constants.MOZILLA_USER_AGENT}}

	var response, err = client.httpClient.Do(request)
	var isSuccessStatusCode = extensions.IsSuccessStatusCode(response)
	if err != nil || !isSuccessStatusCode {
		return []string{}, errors.New(fmt.Sprint("error: failed to retrieve manga chapters uri status code: ", strconv.Itoa(response.StatusCode), " ", http.StatusText(response.StatusCode)))
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		return []string{}, errors.New("error: failed to read request body")
	}

	var chapterApiResult models.ChapterApiResult
	var errJson = json.Unmarshal(body, &chapterApiResult)
	if errJson != nil {
		return []string{}, errors.New("error: failed to deserialize request body")
	}

	var uris []string

	for _, item := range chapterApiResult.Result.Images {
		var items = item.([]interface{})
		uris = append(uris, items[0].(string))
	}

	return uris, nil
}

func (client *MangaClient) GetDivisionImages(uri string) ([]byte, error) {
	var url, errUrl = url.Parse(uri)
	if errUrl != nil {
		return nil, fmt.Errorf("invalid uri passed %w", errUrl)
	}

	var request, errRequest = http.NewRequest("GET", url.String(), nil)
	if errRequest != nil {
		return []byte{}, errors.New("error: failed to create http request")
	}

	request.Header = http.Header{"User-Agent": {constants.MOZILLA_USER_AGENT}}

	var response, err = client.httpClient.Do(request)
	var isSuccessStatusCode = extensions.IsSuccessStatusCode(response)
	if err != nil || !isSuccessStatusCode {
		return []byte{}, errors.New(fmt.Sprint("error: failed to retrieve manga image status code: ", strconv.Itoa(response.StatusCode), " ", http.StatusText(response.StatusCode)))
	}

	defer response.Body.Close()
	var body, errBody = io.ReadAll(response.Body)
	if errBody != nil {
		return []byte{}, errors.New("error: failed to read request body")
	}

	return body, nil
}
