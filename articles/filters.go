package articles

import "strings"

type Matcher interface {
	MatchArticle(article Article) bool
}

type ArticleFilterFunc func(article Article) bool

func (f ArticleFilterFunc) MatchArticle(article Article) bool {
	return f(article)
}

var ArticleAboutCovid = ArticleFilterFunc(func(article Article) bool {
	// TODO: Not the smartest way to compare strings but don't want to use regexp here
	// TODO: Proper case insensitive match with regexp?
	for _, sub := range []string{"Corona", "Covid-19"} {
		if strings.Contains(article.Content.Title, sub) {
			return true
		}
		if strings.Contains(article.Content.TitleHead, sub) {
			return true
		}
		if strings.Contains(article.Content.Lead, sub) {
			return true
		}
	}

	return false
})
