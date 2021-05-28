// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package columnar

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkCollection/add-8         	 5004789	       234.1 ns/op	      82 B/op	       0 allocs/op
// BenchmarkCollection/fetch-to-8    	92310531	        12.28 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCollection/count-8       	 1653796	       725.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCollection/count-indexed-8         	23074526	        51.51 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCollection/find-8                  	 1207858	       996.8 ns/op	     336 B/op	       2 allocs/op
// BenchmarkCollection/find-indexed-8          	 3986691	       303.9 ns/op	     336 B/op	       2 allocs/op
func BenchmarkCollection(b *testing.B) {
	players := loadPlayers()
	obj := Object{
		"name":   "Roman",
		"age":    35,
		"wallet": 50.99,
		"health": 100,
		"mana":   200,
	}

	b.Run("add", func(b *testing.B) {
		col := New()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			col.Add(obj)
			if col.Count(nil) >= 1000 {
				col = New()
			}
		}
	})

	b.Run("fetch-to", func(b *testing.B) {
		dst := make(Object, 8)
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.FetchTo(20, &dst)
		}
	})

	b.Run("count", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.Count(oldHumanMages)
		}
	})

	b.Run("count-indexed", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.Count(oldHumanMagesIndexed)
		}
	})

	b.Run("find", func(b *testing.B) {
		count := 0
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.Find(oldHumanMages, func(o Object) bool {
				count++
				return true
			}, "name")
		}
	})

	b.Run("find-indexed", func(b *testing.B) {
		count := 0
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			players.Find(oldHumanMagesIndexed, func(o Object) bool {
				count++
				return true
			}, "name")
		}
	})

}

func TestCollection(t *testing.T) {
	obj := Object{
		"name":   "Roman",
		"age":    35,
		"wallet": 50.99,
		"health": 100,
		"mana":   200,
	}

	col := New()
	idx := col.Add(obj)

	{ // Find the object by its index
		obj, ok := col.Fetch(idx)
		assert.True(t, ok)
		assert.Equal(t, "Roman", obj["name"])
	}

	{ // Remove the object
		col.Remove(idx)
		obj, ok := col.Fetch(idx)
		assert.False(t, ok)
		assert.Nil(t, obj)
	}

	{ // Add a new one, should replace
		idx := col.Add(obj)
		obj, ok := col.Fetch(idx)
		assert.True(t, ok)
		assert.Equal(t, "Roman", obj["name"])
	}
}

// loadPlayers loads a list of players from the fixture
func loadPlayers() *Collection {
	out := New()

	// index on humans
	out.Index("human", "race", func(v interface{}) bool {
		return v == "human"
	})

	// index for mages
	out.Index("mage", "class", func(v interface{}) bool {
		return v == "mage"
	})

	// index for old
	out.Index("old", "age", func(v interface{}) bool {
		return v.(float64) >= 30
	})

	// Load the items into the collection
	for _, p := range loadFixture("players.json") {
		out.Add(p)
	}
	return out
}

// loadFixture loads a fixture by its name
func loadFixture(name string) []Object {
	b, err := os.ReadFile("fixtures/" + name)
	if err != nil {
		panic(err)
	}

	var data []Object
	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	/*out := make([]Object, 0, 1000*len(data))
	for i := 0; i < 1000; i++ {
		out = append(out, data...)
	}*/
	return data
}
