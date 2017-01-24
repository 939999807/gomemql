package gomemql

import (
	"errors"
	"reflect"
)

type Table struct {
	fieldByName map[string]*tableField
	fields      []*tableField
}

// 添加一行数据
func (self *Table) AddRecord(record interface{}) {

	vRecord := reflect.Indirect(reflect.ValueOf(record))

	// 根据字段进行索引
	for i := 0; i < len(self.fields); i++ {

		recordField := self.fields[i]

		// 结构体中数值
		key := vRecord.Field(i).Interface()

		// 将数值添加到字段索引中, 同一个值可能有多个, 引用记录集合
		recordField.Add(key, record)

	}
}

func (self *Table) FieldByName(name string) *tableField {

	if v, ok := self.fieldByName[name]; ok {

		return v
	}

	return nil
}

func (self *Table) FieldByIndex(index int) *tableField {
	return self.fields[index]
}

func (self *Table) GenFieldIndex(name string, matchTypeStr string, begin, end int32) error {

	if self.NumFields() == 0 {
		return nil
	}

	if self.FieldByIndex(0).KeyCount() == 0 {
		return errors.New("require table data to gen index")
	}

	field := self.FieldByName(name)
	if field == nil {
		return errors.New("field not found: " + name)
	}

	matchType := getMatchTypeBySign(matchTypeStr)
	if matchType == MatchType_Unknown {
		return errors.New("unknown match type: " + matchTypeStr)
	}

	// 做出索引
	field.genIndex(matchType)

	var i int32

	var lastNotNil *RecordList
	for i = begin; i <= end; i++ {
		list := field.getByKey(i)

		if list == nil {

			if lastNotNil != nil {
				field.addIndexData(matchType, i, lastNotNil)
			}

		} else {
			field.addIndexData(matchType, i, list)
		}

		if list != nil {
			lastNotNil = list
		}

	}

	return nil
}

func (self *Table) NumFields() int {
	return len(self.fields)
}

func NewTable(userStruct interface{}) *Table {

	tStruct := reflect.TypeOf(userStruct)
	if tStruct.Kind() == reflect.Ptr {
		tStruct = tStruct.Elem()
	}

	self := &Table{
		fieldByName: make(map[string]*tableField),
		fields:      make([]*tableField, tStruct.NumField()),
	}

	for i := 0; i < tStruct.NumField(); i++ {

		fd := tStruct.Field(i)

		tf := newTableField()

		self.fields[i] = tf
		self.fieldByName[fd.Name] = tf
	}

	return self
}
