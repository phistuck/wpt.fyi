// Copyright 2018 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package query

import (
	"github.com/web-platform-tests/wpt.fyi/shared"
)

// Binder is a mechanism for binding a query over a slice of test runs to
// a particular query service mechanism.
type Binder interface {
	// Bind produces an query execution Plan and/or error after binding its inputs
	// to a query service mechanism. E.g., an in-memory cache may verify that the
	// given runs are in the cache and extract results data subsets that pertain
	// to the runs before producing a Plan implementation that can operate over
	// the subsets directly.
	Bind([]shared.TestRun, ConcreteQuery) (Plan, error)
}

// Plan a query execution plan that returns results.
type Plan interface {
	// Execute runs the query execution plan. The result set type depends on the
	// underlying query service mechanism that the Plan was bound with.
	Execute([]shared.TestRun) interface{}
}

// ConcreteQuery is an AbstractQuery that has been bound to specific test runs.
type ConcreteQuery interface {
	Size() int
}

// ConcreteItemQuery is an AbstractItemQuery that has been bound to a specific test run.
type ConcreteItemQuery interface {
	ConcreteQuery
}

// Exists is constrains search results to require that at least one run meets the
// requirements of its Arg tree.
type Exists struct {
	Runs []shared.TestRun
	Args []ConcreteItemQuery
}

// RunTestStatusEq constrains search results to include only test results from a
// particular run that have a particular test status value. Run IDs are those
// values automatically assigned to shared.TestRun instances by Datastore.
// Status IDs are those codified in shared.TestStatus* symbols.
type RunTestStatusEq struct {
	Run    int64
	Status int64
}

// RunTestStatusNeq constrains search results to include only test results from a
// particular run that do not have a particular test status value. Run IDs are
// those values automatically assigned to shared.TestRun instances by Datastore.
// Status IDs are those codified in shared.TestStatus* symbols.
type RunTestStatusNeq struct {
	Run    int64
	Status int64
}

// Or is a logical disjunction of ConcreteItemQuery instances.
type Or struct {
	Args concreteItemQueries
}

// And is a logical conjunction of ConcreteItemQuery instances.
type And struct {
	Args concreteItemQueries
}

// Not is a logical negation of ConcreteItemQuery instances.
type Not struct {
	Arg ConcreteItemQuery
}

// True is a true-valued ConcreteItemQuery.
type True struct{}

// False is a false-valued ConcreteItemQuery.
type False struct{}

// Size of TestNamePattern has a size of 1: servicing such a query requires a
// substring match per test.
func (TestNamePattern) Size() int { return 1 }

// Size of RunTestStatusEq is 1: servicing such a query requires a single lookup
// in a test run result mapping per test.
func (RunTestStatusEq) Size() int { return 1 }

// Size of RunTestStatusNeq is 1: servicing such a query requires a single
// lookup in a test run result mapping per test.
func (RunTestStatusNeq) Size() int { return 1 }

// Size of Or is the sum of the sizes of its constituent ConcretQuery instances.
func (o Or) Size() int { return o.Args.Size() }

// Size of And is the sum of the sizes of its constituent ConcretQuery
// instances.
func (a And) Size() int { return a.Args.Size() }

// Size of Not is one unit greater than the size of its ConcreteQuery argument.
func (n Not) Size() int { return 1 + n.Arg.Size() }

// Size of True is 0: It should be optimized out of queries in practice.
func (True) Size() int { return 0 }

// Size of False is 0: It should be optimized out of queries in practice.
func (False) Size() int { return 0 }

type concreteQueries []ConcreteQuery

func (c concreteQueries) Size() int {
	s := 0
	for _, q := range c {
		s += q.Size()
	}
	return s
}

type concreteItemQueries []ConcreteItemQuery

func (c concreteItemQueries) Size() int {
	s := 0
	for _, q := range c {
		s += q.Size()
	}
	return s
}
