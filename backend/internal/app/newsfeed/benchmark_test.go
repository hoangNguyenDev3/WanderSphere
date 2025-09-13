package newsfeed

import (
	"testing"
)

func BenchmarkEncodeCursor(b *testing.B) {
	postID := int64(123456789)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		encodeCursor(postID)
	}
}

func BenchmarkDecodeCursor(b *testing.B) {
	cursor := encodeCursor(123456789)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decodeCursor(cursor)
	}
}

func BenchmarkEncodeCursor_Parallel(b *testing.B) {
	postID := int64(123456789)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encodeCursor(postID)
		}
	})
}

func BenchmarkDecodeCursor_Parallel(b *testing.B) {
	cursor := encodeCursor(123456789)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			decodeCursor(cursor)
		}
	})
}

// BenchmarkCursorRoundTrip tests the full encode → decode cycle
func BenchmarkCursorRoundTrip(b *testing.B) {
	postID := int64(987654321)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cursor := encodeCursor(postID)
		decodeCursor(cursor)
	}
}
