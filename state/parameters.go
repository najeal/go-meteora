package state

type StaticParameters struct {
	BaseFactor               uint16
	FilterPeriod             uint16
	DecayPeriod              uint16
	ReductionFactor          uint16
	VariableFeeControl       uint32
	MaxVolatilityAccumulator uint32
	MinBinID                 int32
	MaxBinID                 int32
	ProtocolShare            uint16
	Padding                  [6]uint8
}

type VariableParameters struct {
	VolatilityAccumulator uint32
	VolatilityReference   uint32
	IndexReference        int32
	Padding               [4]uint8
	LastUpdateTimestamp   int64
	Padding1              [8]uint8
}
