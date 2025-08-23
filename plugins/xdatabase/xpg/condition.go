package xpg

// Where 添加 WHERE 条件
func (i *Inst) Where(query interface{}, args ...interface{}) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.Where(query, args...)
	return inst
}

func (i *Inst) XWhere(condition bool, query interface{}, args ...interface{}) *Inst {
	if !condition {
		return i
	}
	i.cond = i.cond.Where(query, args...)
	return i
}

func (i *Inst) Limit(limit int) *Inst {
	inst := i.Copy()
	inst.cond = inst.cond.Limit(uint64(limit))
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
