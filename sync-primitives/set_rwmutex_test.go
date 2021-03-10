package main

import "testing"

func BenchmarkSetRW_10To90(b *testing.B) {
	var set = NewSetRW()

	b.Run("SetRW: 10 write/90 read", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 10; i++ {
					if i%10 == 0 {
						set.Add(i)
					} else {
						set.Has(i)
					}
				}
			}
		})
	})
}

func BenchmarkSetRW_50To50(b *testing.B) {
	var set = NewSetRW()

	b.Run("SetRW: 50 write/50 read", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 10; i++ {
					if i%2 == 0 {
						set.Add(i)
					} else {
						set.Has(i)
					}
				}
			}
		})
	})
}

func BenchmarkSetRW_90To10(b *testing.B) {
	var set = NewSetRW()

	b.Run("SetRW: 90 write/10 read", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 10; i++ {
					if i%10 == 0 {
						set.Has(i)
					} else {
						set.Add(i)
					}
				}
			}
		})
	})
}
