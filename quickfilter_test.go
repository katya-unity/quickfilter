package quickfilter_test

import (
	"fmt"
	"testing"

	"github.com/jussi-kalliokoski/quickfilter"
)

func Test(t *testing.T) {
	t.Run("Add and Iterate", func(t *testing.T) {
		data := generateData(20)
		qf := quickfilter.New(len(data))

		for i := range data {
			if data[i].index%2 == 0 {
				qf = qf.Add(i)
			}
		}
		newData := make([]mockData, 0, qf.Len())
		for it := qf.Iterate(); !it.Done(); it = it.Next() {
			newData = append(newData, data[it.Value()])
		}

		validate(t, len(data), newData)
	})

	t.Run("Fill and Iterate", func(t *testing.T) {
		data := generateData(20)
		qf := quickfilter.New(len(data))
		expectedLen := len(data)

		qf = qf.Fill()
		newData := make([]mockData, 0, qf.Len())
		for it := qf.Iterate(); !it.Done(); it = it.Next() {
			newData = append(newData, data[it.Value()])
		}
		receivedLen := len(newData)

		if expectedLen != receivedLen {
			t.Errorf("expected %d, got %d", expectedLen, receivedLen)
		}
	})

	t.Run("Clear and Iterate", func(t *testing.T) {
		data := generateData(20)
		qf := quickfilter.New(len(data))
		expectedLen := 0

		for i := range data {
			if data[i].index%2 == 0 {
				qf = qf.Add(i)
			}
		}
		qf = qf.Clear()
		newData := make([]mockData, 0, qf.Len())
		for it := qf.Iterate(); !it.Done(); it = it.Next() {
			newData = append(newData, data[it.Value()])
		}
		receivedLen := len(newData)

		if expectedLen != receivedLen {
			t.Errorf("expected %d, got %d", expectedLen, receivedLen)
		}
	})

	t.Run("Fill and Copy", func(t *testing.T) {
		data := generateData(20)
		qf := quickfilter.NewFilled(len(data))
		expectedLen := len(data)

		qf2 := qf.Copy()
		newData := make([]mockData, 0, qf.Len())
		for it := qf2.Iterate(); !it.Done(); it = it.Next() {
			newData = append(newData, data[it.Value()])
		}
		receivedLen := len(newData)

		if expectedLen != receivedLen {
			t.Errorf("expected %d, got %d", expectedLen, receivedLen)
		}
	})

	t.Run("Has", func(t *testing.T) {
		t.Run("found", func(t *testing.T) {
			qf := quickfilter.New(24).Add(12)

			found := qf.Has(12)

			if !found {
				t.Error("expected Has to return true")
			}
		})

		t.Run("not found", func(t *testing.T) {
			qf := quickfilter.New(24).Add(12)

			found := qf.Has(11)

			if found {
				t.Error("expected Has to return false")
			}
		})
	})

	t.Run("Cap", func(t *testing.T) {
		expectedCap := 64
		qf := quickfilter.New(expectedCap)

		receivedCap := qf.Cap()

		if expectedCap != receivedCap {
			t.Errorf("expected %d, got %d", expectedCap, receivedCap)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should decrement len", func(t *testing.T) {
			qf := quickfilter.New(128)
			expectedLen := qf.Cap() / 4

			for i := 0; i < qf.Cap(); i += 2 {
				qf = qf.Add(i)
			}
			for i := 0; i < qf.Cap(); i += 4 {
				qf = qf.Delete(i)
			}
			receivedLen := qf.Len()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
		})

		t.Run("should make items non-iterable", func(t *testing.T) {
			qf := quickfilter.New(128)
			expectedIterated := qf.Cap() / 4

			for i := 0; i < qf.Cap(); i += 2 {
				qf = qf.Add(i)
			}
			for i := 0; i < qf.Cap(); i += 4 {
				qf = qf.Delete(i)
			}
			receivedIterated := 0
			for it := qf.Iterate(); !it.Done(); it = it.Next() {
				receivedIterated++
			}

			if expectedIterated != receivedIterated {
				t.Errorf("expected %d, got %d", expectedIterated, receivedIterated)
			}
		})

		t.Run("double delete should be a no-op", func(t *testing.T) {
			qf := quickfilter.NewFilled(128)
			expectedLen := qf.Len() - 2

			qf = qf.Delete(qf.Cap() / 4)
			qf = qf.Delete(qf.Cap() / 4)
			qf = qf.Delete(qf.Cap() / 2)
			qf = qf.Delete(qf.Cap() / 2)
			receivedLen := qf.Len()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
		})
	})

	t.Run("CopyFrom", func(t *testing.T) {
		qf1 := quickfilter.NewFilled(128)
		qf2 := quickfilter.New(qf1.Cap())
		expectedLen := qf1.Len()

		qf2 = qf2.CopyFrom(qf1)
		receivedLen := qf2.Len()

		if expectedLen != receivedLen {
			t.Errorf("expected %d, got %d", expectedLen, receivedLen)
		}
	})

	t.Run("Resize", func(t *testing.T) {
		t.Run("shrink", func(t *testing.T) {
			expectedCap := 64
			expectedLen := expectedCap
			qf := quickfilter.New(expectedCap * 2)

			qf = qf.Resize(expectedCap)
			qf = qf.Fill()
			receivedCap := qf.Cap()
			receivedLen := qf.Cap()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
			if expectedCap != receivedCap {
				t.Errorf("expected %d, got %d", expectedCap, receivedCap)
			}
		})

		t.Run("grow", func(t *testing.T) {
			expectedCap := 128
			expectedLen := expectedCap
			qf := quickfilter.New(expectedCap / 2)

			qf = qf.Resize(expectedCap)
			qf = qf.Fill()
			receivedCap := qf.Cap()
			receivedLen := qf.Cap()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
			if expectedCap != receivedCap {
				t.Errorf("expected %d, got %d", expectedCap, receivedCap)
			}
		})

		t.Run("noop", func(t *testing.T) {
			expectedCap := 128
			expectedLen := expectedCap
			qf := quickfilter.New(expectedCap)

			qf = qf.Resize(expectedCap)
			qf = qf.Fill()
			receivedCap := qf.Cap()
			receivedLen := qf.Cap()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
			if expectedCap != receivedCap {
				t.Errorf("expected %d, got %d", expectedCap, receivedCap)
			}
		})

		t.Run("shrink & grow", func(t *testing.T) {
			expectedCap := 128
			expectedLen := expectedCap
			qf := quickfilter.New(expectedCap)

			qf = qf.Resize(expectedCap / 2)
			qf = qf.Resize(expectedCap)
			qf = qf.Fill()
			receivedCap := qf.Cap()
			receivedLen := qf.Cap()

			if expectedLen != receivedLen {
				t.Errorf("expected %d, got %d", expectedLen, receivedLen)
			}
			if expectedCap != receivedCap {
				t.Errorf("expected %d, got %d", expectedCap, receivedCap)
			}
		})

		t.Run("should not break batch operations", func(t *testing.T) {
			t.Run("grow", func(t *testing.T) {
				oldCap := 60
				expectedCap := 128
				qf1 := quickfilter.New(oldCap)
				qf2 := quickfilter.New(expectedCap)

				qf1 = qf1.Resize(expectedCap)
				qf1 = qf1.IntersectionOf(qf1, qf2)
				qf2 = qf2.UnionOf(qf1, qf2)
				receivedCap := qf1.Cap()

				if expectedCap != receivedCap {
					t.Errorf("expected %d, got %d", expectedCap, receivedCap)
				}
			})

			t.Run("shrink", func(t *testing.T) {
				oldCap := 128
				expectedCap := 60
				qf1 := quickfilter.New(oldCap)
				qf2 := quickfilter.New(expectedCap)

				qf1 = qf1.Resize(expectedCap)
				qf1 = qf1.IntersectionOf(qf1, qf2)
				qf2 = qf2.UnionOf(qf1, qf2)
				receivedCap := qf1.Cap()

				if expectedCap != receivedCap {
					t.Errorf("expected %d, got %d", expectedCap, receivedCap)
				}
			})
		})
	})

	t.Run("Fill, Intersect and check length", func(t *testing.T) {
		expectedLen := 10
		qf1 := quickfilter.NewFilled(expectedLen)
		qf2 := quickfilter.NewFilled(expectedLen)
		qf2 = qf1.IntersectionOf(qf1, qf2)

		got := qf2.Len()
		if expectedLen != got {
			t.Errorf("expected %d, got %d", expectedLen, got)
		}
	})

	t.Run("Fill, Intersect, Delete and check length", func(t *testing.T) {
		sourceLen := 10
		qf1 := quickfilter.NewFilled(sourceLen)
		qf2 := quickfilter.NewFilled(sourceLen)
		qf2.Delete(8)
		qf2.Delete(1)
		qf2 = qf1.IntersectionOf(qf1, qf2)

		got := qf2.Len()
		expectedLen := 8
		if expectedLen != got {
			t.Errorf("expected %d, got %d", expectedLen, got)
		}
	})

	t.Run("Fill, Union and check length", func(t *testing.T) {
		expectedLen := 10
		qf1 := quickfilter.NewFilled(expectedLen)
		qf2 := quickfilter.NewFilled(expectedLen)
		qf2 = qf1.UnionOf(qf1, qf2)

		got := qf2.Len()
		if expectedLen != got {
			t.Errorf("expected %d, got %d", expectedLen, got)
		}
	})

	t.Run("Fill, Union, Delete and check length", func(t *testing.T) {
		expectedLen := 10
		qf1 := quickfilter.NewFilled(expectedLen)
		qf2 := quickfilter.NewFilled(expectedLen)
		qf2.Delete(8)
		qf2.Delete(1)
		qf2 = qf1.UnionOf(qf1, qf2)

		got := qf2.Len()
		if expectedLen != got {
			t.Errorf("expected %d, got %d", expectedLen, got)
		}
	})
}

func Example() {
	data := make([]int, 0, 8)
	for len(data) < cap(data) {
		data = append(data, len(data))
	}
	qf := quickfilter.New(len(data))
	for i := range data {
		if data[i]%2 == 0 {
			qf = qf.Add(i)
		}
	}
	newData := make([]int, 0, qf.Len())
	for it := qf.Iterate(); !it.Done(); it = it.Next() {
		newData = append(newData, data[it.Value()])
	}
	// Output: [0 2 4 6]
	fmt.Println(newData)
}

func Example_union() {
	data := make([]int, 0, 16)
	for len(data) < cap(data) {
		data = append(data, len(data))
	}
	qf1 := quickfilter.New(len(data))
	for i := range data {
		if data[i]%2 == 0 {
			qf1 = qf1.Add(i)
		}
	}
	qf2 := quickfilter.New(len(data))
	for i := range data {
		if data[i]%3 == 0 {
			qf2 = qf2.Add(i)
		}
	}
	qf := quickfilter.New(len(data))
	qf = qf.UnionOf(qf1, qf2)
	newData := make([]int, 0, qf.Len())
	for it := qf.Iterate(); !it.Done(); it = it.Next() {
		newData = append(newData, data[it.Value()])
	}
	// Output: [0 2 3 4 6 8 9 10 12 14 15]
	fmt.Println(newData)
}

func Example_intersection() {
	data := make([]int, 0, 16)
	for len(data) < cap(data) {
		data = append(data, len(data))
	}
	qf1 := quickfilter.New(len(data))
	for i := range data {
		if data[i]%2 == 0 {
			qf1 = qf1.Add(i)
		}
	}
	qf2 := quickfilter.New(len(data))
	for i := range data {
		if data[i]%3 == 0 {
			qf2 = qf2.Add(i)
		}
	}
	qf := quickfilter.New(len(data))
	qf = qf.IntersectionOf(qf1, qf2)
	newData := make([]int, 0, qf.Len())
	for it := qf.Iterate(); !it.Done(); it = it.Next() {
		newData = append(newData, data[it.Value()])
	}
	// Output: [0 6 12]
	fmt.Println(newData)
}

func Benchmark(b *testing.B) {
	const size = 20000

	b.Run("QuickFilter", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			data := generateData(size)
			qf := quickfilter.New(len(data))
			for i := range data {
				if data[i].index%2 == 0 {
					qf = qf.Add(i)
				}
			}
			newData := make([]mockData, 0, qf.Len())
			for it := qf.Iterate(); !it.Done(); it = it.Next() {
				newData = append(newData, data[it.Value()])
			}
			validate(b, size, newData)
		}
	})

	b.Run("dynamic allocations", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			data := generateData(size)
			newData := make([]mockData, 0)
			for i := range data {
				if data[i].index%2 == 0 {
					newData = append(newData, data[i])
				}
			}
			validate(b, size, newData)
		}
	})

	b.Run("in place", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			data := generateData(size)
			offset := 0
			for i := range data {
				if data[i].index%2 == 0 {
					data[i-offset] = data[i]
				} else {
					offset++
				}
			}
			data = data[:len(data)-offset]
			validate(b, size, data)
		}
	})

	b.Run("in place copied", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			data := generateData(size)
			newData := make([]mockData, len(data))
			copy(newData, data)
			offset := 0
			for i := range newData {
				if newData[i].index%2 == 0 {
					newData[i-offset] = newData[i]
				} else {
					offset++
				}
			}
			newData = newData[:len(newData)-offset]
			validate(b, size, newData)
		}
	})
}

type mockData struct {
	index int
	trash [1000]int
}

func generateData(dataLen int) []mockData {
	data := make([]mockData, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		data = append(data, mockData{index: i})
		_ = data[i].trash
	}
	return data
}

func validate(tb testing.TB, oldLen int, newData []mockData) {
	tb.Helper()
	expectedLen := oldLen / 2
	receivedLen := len(newData)
	if expectedLen != receivedLen {
		tb.Fatalf("expected length %d, received %d", expectedLen, receivedLen)
	}
	for i := range newData {
		if newData[i].index%2 != 0 {
			tb.Fatalf("unexpected index %d", newData[i].index)
		}
	}
}
