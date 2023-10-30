package database

import (
	"fmt"
	"strconv"
	"strings"
)

type Int64 int64

func (i Int64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

//func (i Int64) Value() (driver.Value, error) {
//	return int64(i), nil
//}
//func (i *Int64) Scan(value interface{}) error {
//	intValue, ok := value.(int64)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal int64 value:", value))
//	}
//	*i = Int64(intValue)
//	return nil
//}

// ///////////////////////////////////
func (i Int64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", i)), nil
}
func (i *Int64) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Replace(s, "\"", "", -1)
	if s == "null" || s == "" {
		*i = Int64(0)
	}
	iv, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64(iv)
	return nil
}
