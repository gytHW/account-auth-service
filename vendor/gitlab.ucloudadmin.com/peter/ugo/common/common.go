package common


type CommonErr struct {

	Message             string
	ReturnCode          int32
}

func (inst *CommonErr) Error() string {

	return inst.Message
}

func (inst *CommonErr) RetCode() int32 {

	return inst.ReturnCode
}











