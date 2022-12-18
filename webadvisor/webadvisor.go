package webadvisor

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

type coursesResponse struct {
	courses []struct {
		SectionID   []string `json:"MatchingSectionIds"`
		ID          string   `json:"Id"`
		SubjectCode string   `json:"SubjectCode"`
		Number      string   `json:"Number"`
	}
}

type WebAdvisor struct {
	client  *http.Client
	cookies []*http.Cookie
}

func New() WebAdvisor {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return WebAdvisor{client: client}
}

func (w *WebAdvisor) CheckCourse() {
	token, err := w.getToken()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(token)
	//w.updateCookies()
	w.getSections(token)
}

//func (w *WebAdvisor) updateCookies() {
//	for _, cookie := range w.cookies {
//		w.client.Jar.SetCookies(cookie.Domain, w.cookies)
//	}
//
//}

func (w *WebAdvisor) getToken() (string, error) {

	const url = "https://colleague-ss.uoguelph.ca/Student/Courses"
	req, _ := http.NewRequest("GET", url, nil)

	res, err := w.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	w.cookies = res.Cookies()
	token := w.extractToken(res.Body)
	if token == "" {
		return "", fmt.Errorf("no token found")
	}
	return token, nil
}

func (w *WebAdvisor) extractToken(rb io.Reader) string {
	data, _ := io.ReadAll(rb)
	// Yes regex is hard.
	r, _ := regexp.Compile(`<input name="__RequestVerificationToken" type="hidden" value="[^"]*" />`)
	match := r.FindString(string(data))
	if match == "" {
		log.Println("No match")
		return ""
	}
	line := strings.Split(match, " ")[3]
	token := line[7 : len(line)-1]
	return token
}

func (w *WebAdvisor) getSections(token string) {
	postUrl := "https://colleague-ss.uoguelph.ca/Student/Student/Courses/SearchAsync"
	data := bytes.NewBufferString(`{"searchParameters":"{\"keyword\":\"cis2500\",\"terms\":[],\"requirement\":null,\"subrequirement\":null,\"courseIds\":null,\"sectionIds\":null,\"requirementText\":null,\"subrequirementText\":\"\",\"group\":null,\"startTime\":null,\"endTime\":null,\"openSections\":null,\"subjects\":[],\"academicLevels\":[],\"courseLevels\":[],\"synonyms\":[],\"courseTypes\":[],\"topicCodes\":[],\"days\":[],\"locations\":[],\"faculty\":[],\"onlineCategories\":null,\"keywordComponents\":[],\"startDate\":null,\"endDate\":null,\"startsAtTime\":null,\"endsByTime\":null,\"pageNumber\":1,\"sortOn\":\"None\",\"sortDirection\":\"Ascending\",\"subjectsBadge\":[],\"locationsBadge\":[],\"termFiltersBadge\":[],\"daysBadge\":[],\"facultyBadge\":[],\"academicLevelsBadge\":[],\"courseLevelsBadge\":[],\"courseTypesBadge\":[],\"topicCodesBadge\":[],\"onlineCategoriesBadge\":[],\"openSectionsBadge\":\"\",\"openAndWaitlistedSectionsBadge\":\"\",\"subRequirementText\":null,\"quantityPerPage\":30,\"openAndWaitlistedSections\":null,\"searchResultsView\":\"CatalogListing\"}"}`)
	req, err := http.NewRequest("POST", postUrl, data)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json, charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("__RequestVerificationToken", token)
	//w.client.Jar.SetCookies(req.URL, w.cookies)
	// Get cookie values from the req
	for _, cookie := range w.cookies {
		fmt.Println("Cookie", cookie.Name)
		fmt.Println("Domain", cookie.Value)
		req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}
	resp, err := w.client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.ContentLength)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString[:500])
}
