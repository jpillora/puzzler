package leetcode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jpillora/puzzler/internal/pzlr/x"
)

type problemSpec struct {
	Stat       problemStat `json:"stat"`
	IsPaidOnly bool        `json:"paid_only"`
}

func (p problemSpec) ID() string {
	return fmt.Sprintf("%04d", p.Stat.ID)
}

func (p problemSpec) Slug() string {
	return p.Stat.TitleSlug
}

type problemStat struct {
	ID        int    `json:"question_id"`
	TitleSlug string `json:"question__title_slug"`
}

func getAllLeetcodeProblems() ([]problemSpec, error) {
	rc, err := x.NetCached("leetcode.com", "leetcode-problems", func() (io.ReadCloser, error) {
		fmt.Printf("Fetching leetcode problem list...\n")
		resp, err := http.Get("https://leetcode.com/api/problems/all/")
		return resp.Body, err
	})
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	var data struct {
		StatStatusPairs []problemSpec `json:"stat_status_pairs"`
	}
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return nil, fmt.Errorf("could not decode leetcode problems: %w", err)
	}
	return data.StatStatusPairs, nil
}

func getProblemSpec(id string) (problemSpec, error) {
	num := 0
	slug := ""
	if n, err := strconv.ParseInt(id, 10, 64); err == nil && n > 0 {
		num = int(n)
	} else {
		slug = id
	}
	problems, err := getAllLeetcodeProblems()
	if err != nil {
		return problemSpec{}, fmt.Errorf("could not fetch leetcode problems: %w", err)
	}
	for _, p := range problems {
		if num > 0 && p.Stat.ID == num {
			return p, nil
		}
		if slug != "" && p.Stat.TitleSlug == slug {
			return p, nil
		}
	}
	s := fmt.Sprintf("slug: %s", slug)
	if num > 0 {
		s = fmt.Sprintf("id: %d", num)
	}
	return problemSpec{}, fmt.Errorf("problem spec not found for %s", s)

}

type problemCode struct {
	Data struct {
		Question struct {
			QuestionID         string `json:"questionId"`
			QuestionFrontendID string `json:"questionFrontendId"`
			CodeSnippets       []struct {
				Lang     string `json:"lang"`
				LangSlug string `json:"langSlug"`
				Code     string `json:"code"`
			} `json:"codeSnippets"`
			EnvInfo       string `json:"envInfo"`
			EnableRunCode bool   `json:"enableRunCode"`
		} `json:"question"`
	} `json:"data"`
}

func getProblemCode(slug string) (code string, err error) {
	rc, err := x.NetCached("leetcode.com", "leetcode-problem-"+slug, func() (io.ReadCloser, error) {
		query := `query questionEditorData($s: String!) {
			question(titleSlug: $s) {
				codeSnippets {
					langSlug
					code
				}
			}
		}`
		body, _ := json.Marshal(struct {
			Query     string            `json:"query"`
			Variables map[string]string `json:"variables"`
		}{
			Query: query,
			Variables: map[string]string{
				"s": slug,
			},
		})
		fmt.Printf("Fetching problem code for %s...\n", slug)
		resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		return resp.Body, err
	})
	if err != nil {
		// attempt to improve the error message
		if rc != nil {
			body := bytes.Buffer{}
			io.Copy(&body, rc)
			if body.Len() > 0 {
				err = errors.New(body.String())
			}
		}
		return "", fmt.Errorf("could not fetch leetcode problem code: %w", err)
	}
	defer rc.Close()
	var data problemCode
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return "", fmt.Errorf("could not decode leetcode problem code: %w", err)
	}
	for _, snippet := range data.Data.Question.CodeSnippets {
		if snippet.LangSlug == "golang" {
			return snippet.Code, nil
		}
	}
	return "", fmt.Errorf("go code snippet not found for slug: %s", slug)
}

func fetchQuestionText(slug string) (string, error) {
	fmt.Printf("Fetching problem code for %s...\n", slug)
	url := fmt.Sprintf("https://leetcode.com/problems/%s/", slug)
	rc, err := x.Get(url)
	if err != nil {
		return "", errors.New("failed to fetch aoc description")
	}
	defer rc.Close()
	doc, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		return "", errors.New("invalid html")
	}
	metaDesc := doc.Find(`meta[name="description"]`)
	if metaDesc.Length() == 0 {
		return "", errors.New("no meta desc")
	}
	c, ok := metaDesc.Attr("content")
	if !ok {
		return "", errors.New("no meta desc content")
	}
	c = regexp.MustCompile(`(?m)^[\x{25A0}\x{00A0}\s]+\n`).ReplaceAllString(c, "\n")
	// wrap at 80 chars
	const w = 80
	lines := strings.Split(c, "\n")
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if len(l) < w {
			continue // dont need to wrap
		}
		j := strings.LastIndex(l[:w], " ")
		if j == -1 {
			continue // can't wrap
		}
		head := l[:j]
		tail := l[j+1:]
		lines = append(lines[0:i], append([]string{head, tail}, lines[i+1:]...)...)
	}
	c = strings.Join(lines, "\n")
	return c, nil
}
