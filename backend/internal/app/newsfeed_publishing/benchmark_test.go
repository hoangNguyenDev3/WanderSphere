package newsfeed_publishing_svc

import (
	"testing"
)

// BenchmarkCalculateEngagementScore benchmarks the engagement score calculation
func BenchmarkCalculateEngagementScore(b *testing.B) {
	createdAtUnix := int64(1700000000)
	likes := 42
	comments := 15

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		calculateEngagementScore(createdAtUnix, likes, comments)
	}
}

// BenchmarkCalculateEngagementScore_HighEngagement tests with viral post metrics
func BenchmarkCalculateEngagementScore_HighEngagement(b *testing.B) {
	createdAtUnix := int64(1700000000)
	likes := 100000
	comments := 50000

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		calculateEngagementScore(createdAtUnix, likes, comments)
	}
}

// BenchmarkCalculateEngagementScore_ZeroEngagement tests with a brand-new post (no interactions)
func BenchmarkCalculateEngagementScore_ZeroEngagement(b *testing.B) {
	createdAtUnix := int64(1700000000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		calculateEngagementScore(createdAtUnix, 0, 0)
	}
}

// BenchmarkCalculateEngagementScore_Parallel tests concurrent score calculation
func BenchmarkCalculateEngagementScore_Parallel(b *testing.B) {
	createdAtUnix := int64(1700000000)
	likes := 42
	comments := 15

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			calculateEngagementScore(createdAtUnix, likes, comments)
		}
	})
}
