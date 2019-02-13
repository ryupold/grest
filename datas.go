package grest

//Append more data
func (ds Datas) Append(data ...Data) Datas {
	return append(ds, data...)
}

//Len length
func (ds Datas) Len() int {
	return len(ds)
}

//Get get data at index
func (ds Datas) Get(index int) Data {
	return ds[index]
}

//Distinct distinct collection (set)
func (ds Datas) Distinct(by func(Data) string) Datas {
	keys := make([]string, 0, 10)
	result := make([]Data, 0, 10)

	contains := func(k string, ks []string) bool {
		for _, key := range ks {
			if key == k {
				return true
			}
		}
		return false
	}

	for _, i := range ds {
		key := by(i)
		if !contains(key, keys) {
			keys = append(keys, key)
			result = append(result, i)
		}
	}

	return result
}

//Map transforms data
func (ds Datas) Map(transform func(int, Data) Data) Datas {
	result := make([]Data, len(ds))
	for i, d := range ds {
		result[i] = transform(i, d)
	}
	return result
}

//Filter filters data
func (ds Datas) Filter(predicate func(int, Data) bool) Datas {
	result := make([]Data, 0, len(ds))
	for i, d := range ds {
		if predicate(i, d) {
			result = append(result, d)
		}
	}
	return result
}

//ForEach apply do func to each Data object
func (ds Datas) ForEach(do func(int, Data) error) error {
	for i, d := range ds {
		if err := do(i, d); err != nil {
			return err
		}
	}
	return nil
}
