package webadvisor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Raeein/gmc"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

type WebAdvisor struct {
	client *http.Client
}
type CourseSearchResponse struct {
	Courses []struct {
		MatchingSectionIds []string
		Id                 string
		SubjectCode        string
		Number             string
	}
}
type SectionListResponse struct {
	TermsAndSections []struct {
		Sections []WebAdvisorSection
	}
	Course struct {
		Id string
	}
}

type WebAdvisorSection struct {
	Section struct {
		Capacity  uint
		Available uint
		CourseId  string
		Id        string
		Number    string
		TermId    string
	}
}

func New() WebAdvisor {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return WebAdvisor{client: client}
}

func (w *WebAdvisor) getToken() (string, error) {

	const url = "https://colleague-ss.uoguelph.ca/Student/Courses"
	req, _ := http.NewRequest("GET", url, nil)

	res, err := w.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

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

func (w *WebAdvisor) Exists(section gmc.Section) error {

	token, err := w.getToken()
	if err != nil {
		return err
	}

	courseID, sectionIDs, err := w.searchCourses(token, section)
	if err != nil {
		return fmt.Errorf("failed to search for course: %w", err)
	}

	webAdvisorSections, err := w.listSections(token, courseID, sectionIDs)
	if err != nil {
		return fmt.Errorf("failed to list sections: %w", err)
	}

	for _, webAdvisorSection := range webAdvisorSections {
		if webAdvisorSection.Section.Number == section.Code && webAdvisorSection.Section.TermId == section.Term {
			return nil
		}
	}

	return nil
}

func (w *WebAdvisor) searchCourses(token string, section gmc.Section) (string, []string, error) {

	postUrl := "https://colleague-ss.uoguelph.ca/Student/Courses/SearchAsync"
	data := bytes.NewBufferString(fmt.Sprintf(`{"searchParameters":"{\"keyword\":null,\"terms\":[],\"requirement\":null,\"subrequirement\":null,\"courseIds\":null,\"sectionIds\":null,\"requirementText\":null,\"subrequirementText\":\"\",\"group\":null,\"startTime\":null,\"endTime\":null,\"openSections\":null,\"subjects\":[\"%s\"],\"academicLevels\":[],\"courseLevels\":[],\"synonyms\":[],\"courseTypes\":[],\"topicCodes\":[],\"days\":[],\"locations\":[],\"faculty\":[],\"onlineCategories\":null,\"keywordComponents\":[],\"startDate\":null,\"endDate\":null,\"startsAtTime\":null,\"endsByTime\":null,\"pageNumber\":1,\"sortOn\":\"None\",\"sortDirection\":\"Ascending\",\"subjectsBadge\":[],\"locationsBadge\":[],\"termFiltersBadge\":[],\"daysBadge\":[],\"facultyBadge\":[],\"academicLevelsBadge\":[],\"courseLevelsBadge\":[],\"courseTypesBadge\":[],\"topicCodesBadge\":[],\"onlineCategoriesBadge\":[],\"openSectionsBadge\":\"\",\"openAndWaitlistedSectionsBadge\":\"\",\"subRequirementText\":null,\"quantityPerPage\":500,\"openAndWaitlistedSections\":null,\"searchResultsView\":\"CatalogListing\"}"}`, section.Course.Department))
	req, err := http.NewRequest("POST", postUrl, data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/json, charset=utf-8")
	req.Header.Set("__RequestVerificationToken", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")

	res, err := w.client.Do(req)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	var courseList CourseSearchResponse
	err = json.NewDecoder(res.Body).Decode(&courseList)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode json: %w", err)
	}
	for _, course := range courseList.Courses {
		if course.SubjectCode == section.Course.Department && course.Number == fmt.Sprintf("%d", section.Course.Code) {
			return course.Id, course.MatchingSectionIds, nil
		}
	}

	return "", nil, fmt.Errorf("%s*%d*%s*%s not found", section.Course.Department, section.Course.Code, section.Term, section.Code)
}

func (w *WebAdvisor) listSections(token, courseId string, sectionIds []string) ([]WebAdvisorSection, error) {
	data := bytes.NewBufferString(fmt.Sprintf(`{"courseId":"%s","sectionIds":%s}`+"\n", courseId, "[\""+strings.Join(sectionIds, "\",\"")+"\"]"))
	req, err := http.NewRequest("POST", "https://colleague-ss.uoguelph.ca/Student/Courses/SectionsAsync", data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json, charset=utf-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("__RequestVerificationToken", token)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	res, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	var sectionList SectionListResponse
	err = json.NewDecoder(res.Body).Decode(&sectionList)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	var results []WebAdvisorSection
	for _, termAndSection := range sectionList.TermsAndSections {
		results = append(results, termAndSection.Sections...)
	}

	return results, nil
}
