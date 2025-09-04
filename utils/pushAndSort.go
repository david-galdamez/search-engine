package utils

import (
	"github.com/david-galdamez/search-engine/models"
)

func PushAndSort(arrayResult *[]models.SearchResults, newResult models.SearchResults) {

	*arrayResult = append(*arrayResult, newResult)
	lastIdx := len(*arrayResult) - 1
	heapifyUp(arrayResult, lastIdx)
}

func heapifyUp(arrayResult *[]models.SearchResults, idx int) {
	if idx == 0 {
		return
	}

	parentIdx := (idx - 1) / 2
	parentNode := (*arrayResult)[parentIdx]
	actualNode := (*arrayResult)[idx]
	if actualNode.Score > parentNode.Score {
		(*arrayResult)[idx] = parentNode
		(*arrayResult)[parentIdx] = actualNode
		heapifyUp(arrayResult, parentIdx)
	}
}
