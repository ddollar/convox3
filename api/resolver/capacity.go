package resolver

type Capacity struct {
	cpu CapacityMetric
	mem CapacityMetric
}

type CapacityMetric struct {
	total int32
	used  int32
}

func (c *Capacity) Cpu() CapacityMetric {
	return c.cpu
}

func (c *Capacity) Mem() CapacityMetric {
	return c.mem
}

func (cm CapacityMetric) Total() int32 {
	return cm.total
}

func (cm CapacityMetric) Used() int32 {
	return cm.used
}
