package xid

type IIdGenerator interface {
	NewLong() uint64
}

type ISnowWorker interface {
	NextId() int64
}

type IdGeneratorOptions struct {
	Method            uint16 // 雪花计算方法,（1-漂移算法|2-传统算法），默认1
	BaseTime          int64  // 基础时间（ms单位），不能超过当前系统时间
	WorkerId          uint16 // 机器码，必须由外部设定，最大值 2^WorkerIdBitLength-1
	WorkerIdBitLength byte   // 机器码位长，默认值6，取值范围 [1, 15]（要求：序列数位长+机器码位长不超过22）
	SeqBitLength      byte   // 序列数位长，默认值6，取值范围 [3, 21]（要求：序列数位长+机器码位长不超过22）
	MaxSeqNumber      uint32 // 最大序列数（含），设置范围 [MinSeqNumber, 2^SeqBitLength-1]，默认值0，表示最大序列数取最大值（2^SeqBitLength-1]）
	MinSeqNumber      uint32 // 最小序列数（含），默认值5，取值范围 [5, MaxSeqNumber]，每毫秒的前5个序列数对应编号0-4是保留位，其中1-4是时间回拨相应预留位，0是手工新值预留位
	TopOverCostCount  uint32 // 最大漂移次数（含），默认2000，推荐范围500-10000（与计算能力有关）
}

func NewIdGeneratorOptions(workerId uint16) *IdGeneratorOptions {
	return &IdGeneratorOptions{
		Method:            1,
		WorkerId:          workerId,
		BaseTime:          1533715688000,
		WorkerIdBitLength: 6,
		SeqBitLength:      6,
		MaxSeqNumber:      0,
		MinSeqNumber:      5,
		TopOverCostCount:  2000,
	}
}

type OverCostActionArg struct {
	ActionType             int32
	TimeTick               int64
	WorkerId               uint16
	OverCostCountInOneTerm int32
	GenCountInOneTerm      int32
	TermIndex              int32
}

func (ocaa OverCostActionArg) OverCostActionArg(workerId uint16, timeTick int64, actionType int32, overCostCountInOneTerm int32, genCountWhenOverCost int32, index int32) {
	ocaa.ActionType = actionType
	ocaa.TimeTick = timeTick
	ocaa.WorkerId = workerId
	ocaa.OverCostCountInOneTerm = overCostCountInOneTerm
	ocaa.GenCountInOneTerm = genCountWhenOverCost
	ocaa.TermIndex = index
}
