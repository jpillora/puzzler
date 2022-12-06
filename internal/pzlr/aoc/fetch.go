package aoc

import (
	"errors"
	"fmt"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/jpillora/pzlr/internal/pzlr/x"
)

func fetchQuestion(year, day int) (string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)
	rc, err := x.Get(url)
	if err != nil {
		return "", errors.New("failed to fetch aoc description")
	}
	defer rc.Close()
	doc, err := goquery.NewDocumentFromReader(rc)
	if err != nil {
		return "", errors.New("invalid html")
	}
	article := doc.Find("body > main > article")
	if article.Length() == 0 {
		return "", errors.New("aoc description not found")
	}
	converter := md.NewConverter("", true, nil)
	md := converter.Convert(article)
	return md, nil
}
