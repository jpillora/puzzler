package kernel

import (
	"errors"
	"fmt"
	"io"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/jpillora/puzzler/internal/pzlr/x"
)

func fetchQuestion(year, day int, session string) (string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)
	rc, err := x.GetWith(url, map[string]string{
		"Cookie": "session=" + session,
	})
	if err != nil {
		return "", errors.New("failed to fetch aoc question")
	}
	defer rc.Close()
	doc, err := goquery.NewDocumentFromReader(rc) // io.TeeReader(rc, os.Stdout))
	if err != nil {
		return "", errors.New("invalid html")
	}
	article := doc.Find("body > main > article")
	if article.Length() == 0 {
		return "", errors.New("aoc description HTML node not found")
	}
	converter := md.NewConverter("", true, nil)
	md := converter.Convert(article)
	return md, nil
}

func fetchUserInput(year, day int, session string) (string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	rc, err := x.GetWith(url, map[string]string{
		"Cookie": "session=" + session,
	})
	if err != nil {
		return "", errors.New("failed to fetch aoc input")
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", errors.New("failed to read aoc input")
	}
	return strings.TrimSpace(string(b)), nil
}
