package grest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

//Data is a JSON object
type Data map[string]interface{}

//Datas is a collection of Data objects
type Datas []Data

//NewData instanciates an empty Data object on the heap
func NewData() *Data {
	return &Data{}
}

//Clone this Data object
func (d Data) Clone() Data {
	newData := Data{}
	for k, v := range d {
		newData[k] = v
	}
	return newData
}

//Union merges this Data object with another creating a union. The same keys are overriden by the other.
func (d Data) Union(other Data) Data {
	newData := Data{}
	for k, v := range d {
		newData[k] = v
	}
	for k, v := range other {
		newData[k] = v
	}
	return newData
}

//Intersection takes only keys into the result which are on both sides
func (d Data) Intersection(other Data) Data {
	newData := Data{}
	for k, v := range d {
		if _, ok := other[k]; ok {
			newData[k] = v
		}
	}
	for k, v := range other {
		if _, ok := d[k]; ok {
			newData[k] = v
		}
	}
	return newData
}

//Has returns true if Data contains the key
func (d Data) Has(key string) bool {
	_, contains := d[key]
	return contains
}

//String value for key
func (d Data) String(key string) *string {
	v, ok := d[key]
	if ok {
		s := fmt.Sprint(v)
		return &s
	}
	return nil
}

//Uint64 value for key
func (d Data) Uint64(key string) *uint64 {
	v, ok := d[key]
	if ok {
		s, ok := v.(uint64)
		if !ok {
			asStr := *d.String(key)
			u, err := strconv.ParseUint(asStr, 10, 64)
			if err != nil {
				dd := d.Decimal(key)
				if dd != nil {
					asStr = dd.String()
					s, err = strconv.ParseUint(asStr, 10, 64)
					if err == nil {
						return &s
					}
				}
				return nil
			}
			s = u
		}
		return &s
	}
	return nil
}

//Int64 value for key
func (d Data) Int64(key string) *int64 {
	v, ok := d[key]
	if ok {
		s, ok := v.(int64)
		if !ok {
			asStr := *d.String(key)
			i, err := strconv.ParseInt(asStr, 10, 64)
			if err != nil {
				dd := d.Decimal(key)
				if dd != nil {
					asStr = dd.String()
					s, err = strconv.ParseInt(asStr, 10, 64)
					if err == nil {
						return &s
					}
				}
				return nil
			}
			s = i
		}
		return &s
	}
	return nil
}

//Decimal value for key
func (d Data) Decimal(key string) *decimal.Decimal {
	v, ok := d[key]
	if ok {
		s, ok := v.(decimal.Decimal)
		if !ok {
			dd, err := decimal.NewFromString(*d.String(key))
			if err != nil {
				return nil
			}
			s = dd
		}
		return &s
	}
	return nil
}

//Time value for key
func (d Data) Time(key string, format string) *time.Time {
	v, ok := d[key]
	if ok {
		s, ok := v.(time.Time)
		if !ok {
			date, err := time.Parse(format, *d.String(key))
			if err != nil {
				return nil
			}
			s = date
		}
		return &s
	}
	return nil
}

//UnixTime value for key assuming the source value is an int64
func (d Data) UnixTime(key string) *time.Time {
	v, ok := d[key]
	if !ok {
		return nil
	}

	pt, ok := v.(*time.Time)
	if ok {
		return pt
	}

	t, ok := v.(time.Time)
	if ok {
		return &t
	}

	i := d.Int64(key)
	if i != nil {
		t := time.Unix(*i, 0)
		return &t
	}

	u := d.Uint64(key)
	if u != nil {
		t := time.Unix(int64(*u), 0)
		return &t
	}

	return nil
}

//Data value for key (nested data map)
func (d Data) Data(key string) *Data {
	v, ok := d[key]
	if ok {
		innerData, ok := v.(Data)
		if !ok {
			innerMap, ok := v.(map[string]interface{})
			if !ok {
				return nil
			}
			dd := Data(innerMap)
			return &dd
		}
		return &innerData
	}
	return nil
}

//=== Manipulation ================================================================================

//Add a key:value to this Data object (if it already contains this key it is overriden)
func (d *Data) Add(key string, value interface{}) {
	(*d)[key] = value
}

//Rename renames a key in this Data object
func (d *Data) Rename(key string, newNameForKey string) {
	(*d)[newNameForKey] = (*d)[key]
	delete((*d), key)
}

//Drop deletes keys and their values in this Data object
func (d *Data) Drop(keys ...string) {
	for _, key := range keys {
		delete((*d), key)
	}
}

//Select drops all keys not mentioned
func (d *Data) Select(keys ...string) {
	contains := func(s string) bool {
		for _, x := range keys {
			if x == s {
				return true
			}
		}
		return false
	}

	for key := range *d {
		if !contains(key) {
			delete((*d), key)
		}
	}
}
