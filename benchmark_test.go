package gomemql

import "testing"

func BenchmarkTest(b *testing.B) {

	var tabData []*tableDef

	// 预估每个索引包含的记录
	for i := 0; i < 100; i++ {
		tabData = append(tabData, &tableDef{
			Id:    int32(i + 1),
			Level: int32(i * 10),
			Name:  "kitty",
		})
	}

	tab := NewTable(new(tableDef))
	for _, r := range tabData {
		tab.AddRecord(r)
	}

	b.ResetTimer()
	// 并发查询量
	for i := 0; i < 3000; i++ {
		NewQuery(tab).Where("Id", ">", int32(50)).Where("Level", "==", int32(500)).VisitRawResult(nil)
	}
}
