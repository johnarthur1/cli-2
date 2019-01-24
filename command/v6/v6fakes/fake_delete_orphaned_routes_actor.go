// Code generated by counterfeiter. DO NOT EDIT.
package v6fakes

import (
	sync "sync"

	v2action "code.cloudfoundry.org/cli/actor/v2action"
	v6 "code.cloudfoundry.org/cli/command/v6"
)

type FakeDeleteOrphanedRoutesActor struct {
	DeleteUnmappedRoutesStub        func(string) (v2action.Warnings, error)
	deleteUnmappedRoutesMutex       sync.RWMutex
	deleteUnmappedRoutesArgsForCall []struct {
		arg1 string
	}
	deleteUnmappedRoutesReturns struct {
		result1 v2action.Warnings
		result2 error
	}
	deleteUnmappedRoutesReturnsOnCall map[int]struct {
		result1 v2action.Warnings
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutes(arg1 string) (v2action.Warnings, error) {
	fake.deleteUnmappedRoutesMutex.Lock()
	ret, specificReturn := fake.deleteUnmappedRoutesReturnsOnCall[len(fake.deleteUnmappedRoutesArgsForCall)]
	fake.deleteUnmappedRoutesArgsForCall = append(fake.deleteUnmappedRoutesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteUnmappedRoutes", []interface{}{arg1})
	fake.deleteUnmappedRoutesMutex.Unlock()
	if fake.DeleteUnmappedRoutesStub != nil {
		return fake.DeleteUnmappedRoutesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteUnmappedRoutesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutesCallCount() int {
	fake.deleteUnmappedRoutesMutex.RLock()
	defer fake.deleteUnmappedRoutesMutex.RUnlock()
	return len(fake.deleteUnmappedRoutesArgsForCall)
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutesCalls(stub func(string) (v2action.Warnings, error)) {
	fake.deleteUnmappedRoutesMutex.Lock()
	defer fake.deleteUnmappedRoutesMutex.Unlock()
	fake.DeleteUnmappedRoutesStub = stub
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutesArgsForCall(i int) string {
	fake.deleteUnmappedRoutesMutex.RLock()
	defer fake.deleteUnmappedRoutesMutex.RUnlock()
	argsForCall := fake.deleteUnmappedRoutesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutesReturns(result1 v2action.Warnings, result2 error) {
	fake.deleteUnmappedRoutesMutex.Lock()
	defer fake.deleteUnmappedRoutesMutex.Unlock()
	fake.DeleteUnmappedRoutesStub = nil
	fake.deleteUnmappedRoutesReturns = struct {
		result1 v2action.Warnings
		result2 error
	}{result1, result2}
}

func (fake *FakeDeleteOrphanedRoutesActor) DeleteUnmappedRoutesReturnsOnCall(i int, result1 v2action.Warnings, result2 error) {
	fake.deleteUnmappedRoutesMutex.Lock()
	defer fake.deleteUnmappedRoutesMutex.Unlock()
	fake.DeleteUnmappedRoutesStub = nil
	if fake.deleteUnmappedRoutesReturnsOnCall == nil {
		fake.deleteUnmappedRoutesReturnsOnCall = make(map[int]struct {
			result1 v2action.Warnings
			result2 error
		})
	}
	fake.deleteUnmappedRoutesReturnsOnCall[i] = struct {
		result1 v2action.Warnings
		result2 error
	}{result1, result2}
}

func (fake *FakeDeleteOrphanedRoutesActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteUnmappedRoutesMutex.RLock()
	defer fake.deleteUnmappedRoutesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDeleteOrphanedRoutesActor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ v6.DeleteOrphanedRoutesActor = new(FakeDeleteOrphanedRoutesActor)
