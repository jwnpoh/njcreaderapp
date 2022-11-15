package broker

import (
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

func (b *broker) GetTitle(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Url string `json:"url"`
	}

	s := serializer.NewSerializer(false, "", nil)
	s.Decode(w, r, &userInput)

	title := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"[^\"]*\"\s?\/?>`)

	resp, err := http.Get(userInput.Url)
	if err != nil {
		s := serializer.NewSerializer(true, "url seems to be invalid", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}
	defer resp.Body.Close()

	b2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s := serializer.NewSerializer(true, "unable to parse article content", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	t := title.Find(b2)

	regexHead := regexp.MustCompile(`\s?<meta[^p]*property=\"og:title\"\s?content=\"`)
	regexTail := regexp.MustCompile(`\"\s?\/?>`)

	head := regexHead.Find(t)
	tail := regexTail.Find(t)

	output := bytes.TrimPrefix(t, head)
	output = bytes.TrimSuffix(output, tail)
	output = bytes.TrimSpace(output)

	titleText := html.UnescapeString(string(output))
	if titleText == "" {
		s := serializer.NewSerializer(true, "could not find title, please enter manually", err)
		s.ErrorJson(w, err)
		b.Logger.Error(s, r)
		return
	}

	s = serializer.NewSerializer(false, "successfully retrieved title", titleText)
	s.Encode(w, http.StatusAccepted)
	// b.Logger.Success(s, r)
}

func formatDate(date string) string {
	day := regexp.MustCompile(`^\S+\s`)
	dateString := regexp.MustCompile(`^\S+\s\S+\s\S+`)

	res := day.ReplaceAllString(date, "")
	res = dateString.FindString(res)

	return res
}
