package main

type Page struct {
	Number int
	Label  int
	Active bool
}

func paginator(currentPage int, total int, padding int) []Page {
	if total == 0 {
		return []Page{{Number: 0, Label: 1, Active: true}}
	}

	totalDisplayedPages := padding*2 + 1
	if totalDisplayedPages > total {
		totalDisplayedPages = total
	}
	if totalDisplayedPages == total {
		return makePages(0, total, currentPage)
	}

	paddingLeft := calculatePaddingLeft(currentPage, total, totalDisplayedPages, padding)
	paddingRight := totalDisplayedPages - paddingLeft

	first := currentPage - paddingLeft
	last := currentPage + paddingRight

	return makePages(first, last, currentPage)
}

func calculatePaddingLeft(currentPage int, total int, totalDisplayedPages int, padding int) int {
	if currentPage <= padding {
		return currentPage
	}
	if currentPage >= total-padding {
		return totalDisplayedPages - (total - currentPage)
	}
	return padding
}

func makePages(start int, end int, currentPage int) []Page {
	pages := make([]Page, end-start)
	for i := start; i <= end-1; i++ {
		pages[i-start] = Page{Number: i, Label: i + 1, Active: i == currentPage}
	}
	return pages
}
