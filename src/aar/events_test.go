package aar

import (
	"net/http/httptest"
	"testing"
	"time"
)

type fixtureResults struct {
	count    int
	returned int
}

func (fr *fixtureResults) Close() {
}

func (ft *fixtureResults) Next() bool {
	ft.returned++
	return ft.returned >= ft.count
}

func (ft *fixtureResults) Scan(dest ...interface{}) error {
	switch s := dest[0].(type) {
	case *int32:
		*s = 1
	}

	//TODO: add event data

	switch s := dest[2].(type) {
	case *time.Time:
		*s = time.Now()
	}

	return nil
}

type fixtureEventsFetcher struct {
	count int
}

func (fef *fixtureEventsFetcher) GetEvents(missionID string) (resultsWrapper, error) {
	return &fixtureResults{count: fef.count}, nil
}

func BenchmarkStreamingEvents(b *testing.B) {
	b.ReportAllocs()

	mock := fixtureEventsFetcher{1000}
	w := httptest.NewRecorder()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := outputEvents(&mock, "foo", w)
		if err != nil {
			panic(err)
		}
	}
}
