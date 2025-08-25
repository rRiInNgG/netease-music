// song_wiki_summary_service_test.go

package service

import (
	"fmt"
	"testing"
)

func TestSongWikiSummaryService_GetSongWikiSummary(t *testing.T) {
	service := NewSongWikiSummaryService(nil)

	code, resp := service.GetSongWikiSummary("2122976910")
	fmt.Println(code, string(resp))

	if code != 200 {
		t.Errorf("code error: %f", code)
	}
}
