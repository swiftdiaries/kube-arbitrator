/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package schedulercache

import (
	"fmt"
	"log"
	"math"

	"k8s.io/api/core/v1"
)

//Resource struct to define what resources are equalized
type Resource struct {
	MilliCPU  float64
	Memory    float64
	NvidiaGPU float64
}

func Decorator(fn func(r *Resource)) func(r *Resource) {
	return func(r *Resource) {
		log.Println("starting")
		fn(r)
		log.Println("completed")
	}
}

func EmptyResource() *Resource {
	return &Resource{
		MilliCPU:  0,
		Memory:    0,
		NvidiaGPU: 0,
	}
}

func (r *Resource) Clone() *Resource {
	clone := &Resource{
		MilliCPU:  r.MilliCPU,
		Memory:    r.Memory,
		NvidiaGPU: r.NvidiaGPU,
	}
	return clone
}

var minMilliCPU float64 = 10
var minMemory float64 = 10 * 1024 * 1024
var minNvidiaGPU float64 = 1

func NewResource(rl v1.ResourceList) *Resource {
	r := EmptyResource()
	for rName, rQuant := range rl {
		switch rName {
		case v1.ResourceCPU:
			r.MilliCPU += float64(rQuant.MilliValue())
		case v1.ResourceMemory:
			r.Memory += float64(rQuant.Value())
		case v1.ResourceNvidiaGPU:
			r.NvidiaGPU += float64(rQuant.Value())
		}
	}
	return r
}

func (r *Resource) IsEmpty() bool {
	return r.MilliCPU < minMilliCPU && r.Memory < minMemory && r.NvidiaGPU < minNvidiaGPU
}

func (r *Resource) Add(rr *Resource) *Resource {
	r.MilliCPU += rr.MilliCPU
	r.Memory += rr.Memory
	r.NvidiaGPU += rr.NvidiaGPU
	return r
}

//A function to Subtract two Resource objects.
func (r *Resource) Sub(rr *Resource) *Resource {
	if r.Less(rr) == false {
		r.MilliCPU -= rr.MilliCPU
		r.Memory -= rr.Memory
		r.NvidiaGPU -= rr.NvidiaGPU
		return r
	}
	panic("Resource is not sufficient to do operation: Sub()")
}

func (r *Resource) Less(rr *Resource) bool {
	return r.MilliCPU < rr.MilliCPU && r.Memory < rr.Memory && r.NvidiaGPU < rr.NvidiaGPU
}

func (r *Resource) LessEqual(rr *Resource) bool {
	return (r.MilliCPU < rr.MilliCPU || math.Abs(rr.MilliCPU-r.MilliCPU) < 0.01) &&
		(r.Memory < rr.Memory || math.Abs(rr.Memory-r.Memory) < 1) &&
		(r.NvidiaGPU < rr.NvidiaGPU || math.Abs(rr.NvidiaGPU-r.NvidiaGPU) < 1)
}

func (r *Resource) String() string {
	return fmt.Sprintf("cpu %f, memory %f, NvidiaGPU %f", r.MilliCPU, r.Memory, r.NvidiaGPU)
}
