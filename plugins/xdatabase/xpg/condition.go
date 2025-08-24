package xpg

// Where 添加 WHERE 条件
func (i *Inst) Where(query interface{}, args ...interface{}) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.Where(query, args...)
	return inst
}

func (i *Inst) WithId(id int64) *Inst {
	return i.Where("id = ?", id)
}
func (i *Inst) WithIds(ids []int64) *Inst {
	return i.Where("id = ANY(?)", ids)
}

func (i *Inst) XWhere(condition bool, query interface{}, args ...interface{}) *Inst {
	if !condition {
		return i
	}
	return i.Where(query, args...)
}

func (i *Inst) Limit(limit int) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.Limit(uint64(limit))
	return inst
}

func (i *Inst) Group(field ...string) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.GroupBy(field...)
	return inst
}

func (i *Inst) Offset(offset int) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.Offset(uint64(offset))
	return inst
}

func (i *Inst) Order(order ...string) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.OrderBy(order...)
	return inst
}

func (i *Inst) Scopes(f ...func(inst *Inst) *Inst) *Inst {
	inst := i.Copy()
	for _, scope := range f {
		inst = scope(inst)
	}
	return inst
}
