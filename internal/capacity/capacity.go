package capacity

import "errors"

type Resources struct {
	VCPU    int64
	RAMMiB  int64
	DiskGiB int64
}

type Node struct {
	ID       string
	Enabled  bool
	Total    Resources
	Reserved Resources
	Used     Resources
}

var ErrNoCapacity = errors.New("no eligible node has sufficient capacity")

func (r Resources) valid() bool {
	return r.VCPU >= 0 && r.RAMMiB >= 0 && r.DiskGiB >= 0
}

func remaining(total, reserved, used int64) int64 {
	if reserved > total || used > total-reserved {
		return 0
	}
	return total - reserved - used
}

func (n Node) Available() Resources {
	if !n.Enabled || n.ID == "" || !n.Total.valid() || !n.Reserved.valid() || !n.Used.valid() {
		return Resources{}
	}
	return Resources{
		VCPU:    remaining(n.Total.VCPU, n.Reserved.VCPU, n.Used.VCPU),
		RAMMiB:  remaining(n.Total.RAMMiB, n.Reserved.RAMMiB, n.Used.RAMMiB),
		DiskGiB: remaining(n.Total.DiskGiB, n.Reserved.DiskGiB, n.Used.DiskGiB),
	}
}

func (n Node) Fits(req Resources) bool {
	a := n.Available()
	return req.VCPU > 0 && req.RAMMiB > 0 && req.DiskGiB > 0 &&
		req.VCPU <= a.VCPU && req.RAMMiB <= a.RAMMiB && req.DiskGiB <= a.DiskGiB
}

// Select chooses an eligible node by available RAM, then stable ID.
func Select(nodes []Node, req Resources) (Node, error) {
	var best Node
	found := false
	for _, n := range nodes {
		if !n.Fits(req) {
			continue
		}
		if !found || n.Available().RAMMiB > best.Available().RAMMiB ||
			(n.Available().RAMMiB == best.Available().RAMMiB && n.ID < best.ID) {
			best = n
			found = true
		}
	}
	if !found {
		return Node{}, ErrNoCapacity
	}
	return best, nil
}
