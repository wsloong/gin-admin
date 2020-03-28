package util

import (
	"fmt"
	"os"
	"time"

	"github.com/LyricTian/structs"
)

var (
	pid = os.Getpid()
)

// NewTraceID 创建追踪ID
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s",
		pid,
		time.Now().Format("2006.01.02.15.04.05.999999"))
}

// NewRecordID 创建记录ID
func NewRecordID() string {
	return NewObjectID().Hex()
}

// StructMapToStruct 结构体映射
func StructMapToStruct(s, ts interface{}) error {
	if !structs.IsStruct(s) || !structs.IsStruct(ts) {
		return nil
	}

	ss, tss := structs.New(s), structs.New(ts)
	for _, field := range tss.Fields() {
		if !field.IsExported() {
			continue
		}

		if sf, ok := ss.FieldOk(field.Name()); ok {
			err := field.Set2(sf.Value())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
