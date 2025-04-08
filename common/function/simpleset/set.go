/*
 * @Author: majian
 * @Date: 2024-12-05 17:16:04
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-05 17:25:29
 */
package simpleset

type EmptySt struct {
}

type Set struct {
	SetMap map[interface{}]EmptySt
}

func NewSet() *Set {
	set := &Set{}
	set.SetMap = make(map[interface{}]EmptySt)
	return set
}

// 返回值表示是否新增
func (m *Set) Add(val interface{}) bool {
	if _, ok := m.SetMap[val]; ok {
		return false
	}
	m.SetMap[val] = EmptySt{}
	return true
}

func (m *Set) Rem(val interface{}) bool {
	if _, ok := m.SetMap[val]; ok {
		delete(m.SetMap, val)
		return true
	}
	return false
}

func (m *Set) Has(val interface{}) bool {
	if _, ok := m.SetMap[val]; ok {
		return true
	}
	return false
}

func (m *Set) Size() int {
	return len(m.SetMap)
}
