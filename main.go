package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	base_url   = "https://api.github.com"
	owner      = "code-flu"
	repository = "website-content"
)

var headers = map[string]string{
	"Authorization": "Bearer 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",
	"Accept":        "application/vnd.github.v3+json",
}

var (
	BASE_PATH_FORMAT  = "%s/repos/%s/%s"
	BRANCHES_ENDPOINT = "%s/branches"
	CONTENTS_ENDPOINT = "%s/contents"
	COMMITS_ENDPOINT  = "%s/commits"
)

// URLSet corresponds to the <urlset> element
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL corresponds to the <url> element
type URL struct {
	Loc string `xml:"loc"`
}

var urls = []string{
	"https://www.codeflu.com",
	"https://www.codeflu.com/JsonToSql",
	"https://www.codeflu.com/mavenTree",
}

func main() {
	BASE_PATH_FORMAT = fmt.Sprintf(BASE_PATH_FORMAT, base_url, owner, repository)
	BRANCHES_ENDPOINT = fmt.Sprintf(BRANCHES_ENDPOINT, BASE_PATH_FORMAT)
	CONTENTS_ENDPOINT = fmt.Sprintf(CONTENTS_ENDPOINT, BASE_PATH_FORMAT)
	COMMITS_ENDPOINT = fmt.Sprintf(COMMITS_ENDPOINT, BASE_PATH_FORMAT)

	content := make(map[string]interface{})
	content, _ = getBranches(content)
	for branchName := range content {
		content, _ = getDirectories(branchName, content)
	}
	for branchName := range content {
		dirs := content[branchName].(map[string]interface{})
		for directoryName := range dirs {
			content, _ = getDirectoryFiles(branchName, directoryName, content)
		}
	}

	for branchName := range content {
		dirs := content[branchName].(map[string]interface{})
		for directoryName := range dirs {
			files := dirs[directoryName].(map[string]interface{})
			for fileName := range files {
				content, _ = getDirectoryFileInfo(branchName, directoryName, fileName, content)
			}
		}
	}

	jsonStr, _ := json.Marshal(content)
	os.WriteFile("content.json", jsonStr, 0644)
	writeSitemap(content)
}

func getDirectoryFileInfo(branchName string, directoryName string, fileName string, content map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?path=%s/%s&sha=%s", COMMITS_ENDPOINT, strings.Replace(directoryName, " ", "%20", -1), strings.Replace(fileName, " ", "%20", -1), branchName)
	nodes, err := getContent(url, headers)
	if err != nil {
		return nil, err
	}
	branchContent := content[branchName].(map[string]interface{})
	directoryContent := branchContent[directoryName].(map[string]interface{})
	fileContent := directoryContent[fileName].(map[string]interface{})
	for _, node := range nodes {
		commitNode := node["commit"].(map[string]interface{})
		committerNode := commitNode["committer"].(map[string]interface{})
		message := commitNode["message"].(string)
		date := committerNode["date"].(string)
		author := committerNode["name"].(string)
		fileContent["date"] = getTimeAgo(date)
		fileContent["description"] = message
		fileContent["author"] = author
		fileContent["tag"] = directoryName
		break
	}
	return content, nil
}

func getDirectoryFiles(branchName string, directoryName string, content map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s?ref=%s", CONTENTS_ENDPOINT, directoryName, branchName)
	nodes, err := getContent(url, headers)
	if err != nil {
		return nil, err
	}
	branchContent := content[branchName].(map[string]interface{})
	directoryContent := branchContent[directoryName].(map[string]interface{})
	for _, node := range nodes {
		fileName := node["name"].(string)
		filePath := strings.Replace(node["path"].(string), " ", "%20", -1)
		directoryContent[fileName] = map[string]interface{}{"path": strings.TrimSuffix(filePath, ".md")}
	}
	content[branchName] = branchContent
	return content, nil
}

func getDirectories(branchName string, content map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?ref=%s", CONTENTS_ENDPOINT, branchName)
	nodes, err := getContent(url, headers)
	if err != nil {
		return nil, err
	}
	branchContent := content[branchName].(map[string]interface{})
	for _, node := range nodes {
		dirName := node["name"].(string)
		branchContent[dirName] = make(map[string]interface{})
	}
	return content, nil
}

func getBranches(content map[string]interface{}) (map[string]interface{}, error) {
	nodes, err := getContent(BRANCHES_ENDPOINT, headers)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		name := node["name"].(string)
		if name == "main" {
			continue
		}
		content[name] = map[string]interface{}{}
	}
	return content, nil
}

func getContent(url string, headers map[string]string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Resource %s Not Found", url)
		return nil, fmt.Errorf("resource not found")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var arrayNode []map[string]interface{}
	if err := json.Unmarshal(body, &arrayNode); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}
	return arrayNode, nil
}

func getTimeAgo(inputDate string) string {
	parsedTime, _ := time.Parse(time.RFC3339, inputDate)
	timeAgo := parsedTime.Format("02-Jan-2006")
	return timeAgo
}

func writeSitemap(content map[string]interface{}) {
	generateSitemapURLs("https://www.codeflu.com/post/", content, &urls)
	// Prepare the URLSet data structure
	urlSet := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, len(urls)),
	}
	// Populate the URLSet with URLs
	for i, url := range urls {
		urlSet.URLs[i] = URL{Loc: url}
	}
	// Marshal the URLSet into XML
	output, _ := xml.MarshalIndent(urlSet, "", "    ")
	os.WriteFile("sitemap.xml", []byte(xml.Header+string(output)), 0644)
}

func generateSitemapURLs(basePath string, data map[string]interface{}, urls *[]string) {
	for key, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			nextPath := fmt.Sprintf("%s/%s", strings.Trim(basePath, "/"), strings.Replace(strings.TrimSuffix(key, ".md"), " ", "%20", -1))
			generateSitemapURLs(nextPath, nestedMap, urls)
		} else {
			*urls = append(*urls, basePath)
			return
		}
	}
}

