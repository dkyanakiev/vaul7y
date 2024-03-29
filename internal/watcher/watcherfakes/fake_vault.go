// Code generated by counterfeiter. DO NOT EDIT.
package watcherfakes

import (
	"sync"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/dkyanakiev/vaulty/internal/watcher"
	"github.com/hashicorp/vault/api"
)

type FakeVault struct {
	AddressStub        func() string
	addressMutex       sync.RWMutex
	addressArgsForCall []struct {
	}
	addressReturns struct {
		result1 string
	}
	addressReturnsOnCall map[int]struct {
		result1 string
	}
	AllMountsStub        func() (map[string]*models.MountOutput, error)
	allMountsMutex       sync.RWMutex
	allMountsArgsForCall []struct {
	}
	allMountsReturns struct {
		result1 map[string]*models.MountOutput
		result2 error
	}
	allMountsReturnsOnCall map[int]struct {
		result1 map[string]*models.MountOutput
		result2 error
	}
	AllPoliciesStub        func() ([]string, error)
	allPoliciesMutex       sync.RWMutex
	allPoliciesArgsForCall []struct {
	}
	allPoliciesReturns struct {
		result1 []string
		result2 error
	}
	allPoliciesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	GetPolicyInfoStub        func(string) (string, error)
	getPolicyInfoMutex       sync.RWMutex
	getPolicyInfoArgsForCall []struct {
		arg1 string
	}
	getPolicyInfoReturns struct {
		result1 string
		result2 error
	}
	getPolicyInfoReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetSecretInfoStub        func(string, string) (*api.Secret, error)
	getSecretInfoMutex       sync.RWMutex
	getSecretInfoArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getSecretInfoReturns struct {
		result1 *api.Secret
		result2 error
	}
	getSecretInfoReturnsOnCall map[int]struct {
		result1 *api.Secret
		result2 error
	}
	ListNamespacesStub        func() ([]string, error)
	listNamespacesMutex       sync.RWMutex
	listNamespacesArgsForCall []struct {
	}
	listNamespacesReturns struct {
		result1 []string
		result2 error
	}
	listNamespacesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	ListNestedSecretsStub        func(string, string) ([]models.SecretPath, error)
	listNestedSecretsMutex       sync.RWMutex
	listNestedSecretsArgsForCall []struct {
		arg1 string
		arg2 string
	}
	listNestedSecretsReturns struct {
		result1 []models.SecretPath
		result2 error
	}
	listNestedSecretsReturnsOnCall map[int]struct {
		result1 []models.SecretPath
		result2 error
	}
	ListSecretsStub        func(string) (*api.Secret, error)
	listSecretsMutex       sync.RWMutex
	listSecretsArgsForCall []struct {
		arg1 string
	}
	listSecretsReturns struct {
		result1 *api.Secret
		result2 error
	}
	listSecretsReturnsOnCall map[int]struct {
		result1 *api.Secret
		result2 error
	}
	SetNamespaceStub        func(string)
	setNamespaceMutex       sync.RWMutex
	setNamespaceArgsForCall []struct {
		arg1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVault) Address() string {
	fake.addressMutex.Lock()
	ret, specificReturn := fake.addressReturnsOnCall[len(fake.addressArgsForCall)]
	fake.addressArgsForCall = append(fake.addressArgsForCall, struct {
	}{})
	stub := fake.AddressStub
	fakeReturns := fake.addressReturns
	fake.recordInvocation("Address", []interface{}{})
	fake.addressMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVault) AddressCallCount() int {
	fake.addressMutex.RLock()
	defer fake.addressMutex.RUnlock()
	return len(fake.addressArgsForCall)
}

func (fake *FakeVault) AddressCalls(stub func() string) {
	fake.addressMutex.Lock()
	defer fake.addressMutex.Unlock()
	fake.AddressStub = stub
}

func (fake *FakeVault) AddressReturns(result1 string) {
	fake.addressMutex.Lock()
	defer fake.addressMutex.Unlock()
	fake.AddressStub = nil
	fake.addressReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeVault) AddressReturnsOnCall(i int, result1 string) {
	fake.addressMutex.Lock()
	defer fake.addressMutex.Unlock()
	fake.AddressStub = nil
	if fake.addressReturnsOnCall == nil {
		fake.addressReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.addressReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeVault) AllMounts() (map[string]*models.MountOutput, error) {
	fake.allMountsMutex.Lock()
	ret, specificReturn := fake.allMountsReturnsOnCall[len(fake.allMountsArgsForCall)]
	fake.allMountsArgsForCall = append(fake.allMountsArgsForCall, struct {
	}{})
	stub := fake.AllMountsStub
	fakeReturns := fake.allMountsReturns
	fake.recordInvocation("AllMounts", []interface{}{})
	fake.allMountsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) AllMountsCallCount() int {
	fake.allMountsMutex.RLock()
	defer fake.allMountsMutex.RUnlock()
	return len(fake.allMountsArgsForCall)
}

func (fake *FakeVault) AllMountsCalls(stub func() (map[string]*models.MountOutput, error)) {
	fake.allMountsMutex.Lock()
	defer fake.allMountsMutex.Unlock()
	fake.AllMountsStub = stub
}

func (fake *FakeVault) AllMountsReturns(result1 map[string]*models.MountOutput, result2 error) {
	fake.allMountsMutex.Lock()
	defer fake.allMountsMutex.Unlock()
	fake.AllMountsStub = nil
	fake.allMountsReturns = struct {
		result1 map[string]*models.MountOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) AllMountsReturnsOnCall(i int, result1 map[string]*models.MountOutput, result2 error) {
	fake.allMountsMutex.Lock()
	defer fake.allMountsMutex.Unlock()
	fake.AllMountsStub = nil
	if fake.allMountsReturnsOnCall == nil {
		fake.allMountsReturnsOnCall = make(map[int]struct {
			result1 map[string]*models.MountOutput
			result2 error
		})
	}
	fake.allMountsReturnsOnCall[i] = struct {
		result1 map[string]*models.MountOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) AllPolicies() ([]string, error) {
	fake.allPoliciesMutex.Lock()
	ret, specificReturn := fake.allPoliciesReturnsOnCall[len(fake.allPoliciesArgsForCall)]
	fake.allPoliciesArgsForCall = append(fake.allPoliciesArgsForCall, struct {
	}{})
	stub := fake.AllPoliciesStub
	fakeReturns := fake.allPoliciesReturns
	fake.recordInvocation("AllPolicies", []interface{}{})
	fake.allPoliciesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) AllPoliciesCallCount() int {
	fake.allPoliciesMutex.RLock()
	defer fake.allPoliciesMutex.RUnlock()
	return len(fake.allPoliciesArgsForCall)
}

func (fake *FakeVault) AllPoliciesCalls(stub func() ([]string, error)) {
	fake.allPoliciesMutex.Lock()
	defer fake.allPoliciesMutex.Unlock()
	fake.AllPoliciesStub = stub
}

func (fake *FakeVault) AllPoliciesReturns(result1 []string, result2 error) {
	fake.allPoliciesMutex.Lock()
	defer fake.allPoliciesMutex.Unlock()
	fake.AllPoliciesStub = nil
	fake.allPoliciesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) AllPoliciesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.allPoliciesMutex.Lock()
	defer fake.allPoliciesMutex.Unlock()
	fake.AllPoliciesStub = nil
	if fake.allPoliciesReturnsOnCall == nil {
		fake.allPoliciesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.allPoliciesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) GetPolicyInfo(arg1 string) (string, error) {
	fake.getPolicyInfoMutex.Lock()
	ret, specificReturn := fake.getPolicyInfoReturnsOnCall[len(fake.getPolicyInfoArgsForCall)]
	fake.getPolicyInfoArgsForCall = append(fake.getPolicyInfoArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetPolicyInfoStub
	fakeReturns := fake.getPolicyInfoReturns
	fake.recordInvocation("GetPolicyInfo", []interface{}{arg1})
	fake.getPolicyInfoMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) GetPolicyInfoCallCount() int {
	fake.getPolicyInfoMutex.RLock()
	defer fake.getPolicyInfoMutex.RUnlock()
	return len(fake.getPolicyInfoArgsForCall)
}

func (fake *FakeVault) GetPolicyInfoCalls(stub func(string) (string, error)) {
	fake.getPolicyInfoMutex.Lock()
	defer fake.getPolicyInfoMutex.Unlock()
	fake.GetPolicyInfoStub = stub
}

func (fake *FakeVault) GetPolicyInfoArgsForCall(i int) string {
	fake.getPolicyInfoMutex.RLock()
	defer fake.getPolicyInfoMutex.RUnlock()
	argsForCall := fake.getPolicyInfoArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeVault) GetPolicyInfoReturns(result1 string, result2 error) {
	fake.getPolicyInfoMutex.Lock()
	defer fake.getPolicyInfoMutex.Unlock()
	fake.GetPolicyInfoStub = nil
	fake.getPolicyInfoReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) GetPolicyInfoReturnsOnCall(i int, result1 string, result2 error) {
	fake.getPolicyInfoMutex.Lock()
	defer fake.getPolicyInfoMutex.Unlock()
	fake.GetPolicyInfoStub = nil
	if fake.getPolicyInfoReturnsOnCall == nil {
		fake.getPolicyInfoReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getPolicyInfoReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) GetSecretInfo(arg1 string, arg2 string) (*api.Secret, error) {
	fake.getSecretInfoMutex.Lock()
	ret, specificReturn := fake.getSecretInfoReturnsOnCall[len(fake.getSecretInfoArgsForCall)]
	fake.getSecretInfoArgsForCall = append(fake.getSecretInfoArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.GetSecretInfoStub
	fakeReturns := fake.getSecretInfoReturns
	fake.recordInvocation("GetSecretInfo", []interface{}{arg1, arg2})
	fake.getSecretInfoMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) GetSecretInfoCallCount() int {
	fake.getSecretInfoMutex.RLock()
	defer fake.getSecretInfoMutex.RUnlock()
	return len(fake.getSecretInfoArgsForCall)
}

func (fake *FakeVault) GetSecretInfoCalls(stub func(string, string) (*api.Secret, error)) {
	fake.getSecretInfoMutex.Lock()
	defer fake.getSecretInfoMutex.Unlock()
	fake.GetSecretInfoStub = stub
}

func (fake *FakeVault) GetSecretInfoArgsForCall(i int) (string, string) {
	fake.getSecretInfoMutex.RLock()
	defer fake.getSecretInfoMutex.RUnlock()
	argsForCall := fake.getSecretInfoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeVault) GetSecretInfoReturns(result1 *api.Secret, result2 error) {
	fake.getSecretInfoMutex.Lock()
	defer fake.getSecretInfoMutex.Unlock()
	fake.GetSecretInfoStub = nil
	fake.getSecretInfoReturns = struct {
		result1 *api.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) GetSecretInfoReturnsOnCall(i int, result1 *api.Secret, result2 error) {
	fake.getSecretInfoMutex.Lock()
	defer fake.getSecretInfoMutex.Unlock()
	fake.GetSecretInfoStub = nil
	if fake.getSecretInfoReturnsOnCall == nil {
		fake.getSecretInfoReturnsOnCall = make(map[int]struct {
			result1 *api.Secret
			result2 error
		})
	}
	fake.getSecretInfoReturnsOnCall[i] = struct {
		result1 *api.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListNamespaces() ([]string, error) {
	fake.listNamespacesMutex.Lock()
	ret, specificReturn := fake.listNamespacesReturnsOnCall[len(fake.listNamespacesArgsForCall)]
	fake.listNamespacesArgsForCall = append(fake.listNamespacesArgsForCall, struct {
	}{})
	stub := fake.ListNamespacesStub
	fakeReturns := fake.listNamespacesReturns
	fake.recordInvocation("ListNamespaces", []interface{}{})
	fake.listNamespacesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) ListNamespacesCallCount() int {
	fake.listNamespacesMutex.RLock()
	defer fake.listNamespacesMutex.RUnlock()
	return len(fake.listNamespacesArgsForCall)
}

func (fake *FakeVault) ListNamespacesCalls(stub func() ([]string, error)) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = stub
}

func (fake *FakeVault) ListNamespacesReturns(result1 []string, result2 error) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = nil
	fake.listNamespacesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListNamespacesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = nil
	if fake.listNamespacesReturnsOnCall == nil {
		fake.listNamespacesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.listNamespacesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListNestedSecrets(arg1 string, arg2 string) ([]models.SecretPath, error) {
	fake.listNestedSecretsMutex.Lock()
	ret, specificReturn := fake.listNestedSecretsReturnsOnCall[len(fake.listNestedSecretsArgsForCall)]
	fake.listNestedSecretsArgsForCall = append(fake.listNestedSecretsArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.ListNestedSecretsStub
	fakeReturns := fake.listNestedSecretsReturns
	fake.recordInvocation("ListNestedSecrets", []interface{}{arg1, arg2})
	fake.listNestedSecretsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) ListNestedSecretsCallCount() int {
	fake.listNestedSecretsMutex.RLock()
	defer fake.listNestedSecretsMutex.RUnlock()
	return len(fake.listNestedSecretsArgsForCall)
}

func (fake *FakeVault) ListNestedSecretsCalls(stub func(string, string) ([]models.SecretPath, error)) {
	fake.listNestedSecretsMutex.Lock()
	defer fake.listNestedSecretsMutex.Unlock()
	fake.ListNestedSecretsStub = stub
}

func (fake *FakeVault) ListNestedSecretsArgsForCall(i int) (string, string) {
	fake.listNestedSecretsMutex.RLock()
	defer fake.listNestedSecretsMutex.RUnlock()
	argsForCall := fake.listNestedSecretsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeVault) ListNestedSecretsReturns(result1 []models.SecretPath, result2 error) {
	fake.listNestedSecretsMutex.Lock()
	defer fake.listNestedSecretsMutex.Unlock()
	fake.ListNestedSecretsStub = nil
	fake.listNestedSecretsReturns = struct {
		result1 []models.SecretPath
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListNestedSecretsReturnsOnCall(i int, result1 []models.SecretPath, result2 error) {
	fake.listNestedSecretsMutex.Lock()
	defer fake.listNestedSecretsMutex.Unlock()
	fake.ListNestedSecretsStub = nil
	if fake.listNestedSecretsReturnsOnCall == nil {
		fake.listNestedSecretsReturnsOnCall = make(map[int]struct {
			result1 []models.SecretPath
			result2 error
		})
	}
	fake.listNestedSecretsReturnsOnCall[i] = struct {
		result1 []models.SecretPath
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListSecrets(arg1 string) (*api.Secret, error) {
	fake.listSecretsMutex.Lock()
	ret, specificReturn := fake.listSecretsReturnsOnCall[len(fake.listSecretsArgsForCall)]
	fake.listSecretsArgsForCall = append(fake.listSecretsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ListSecretsStub
	fakeReturns := fake.listSecretsReturns
	fake.recordInvocation("ListSecrets", []interface{}{arg1})
	fake.listSecretsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVault) ListSecretsCallCount() int {
	fake.listSecretsMutex.RLock()
	defer fake.listSecretsMutex.RUnlock()
	return len(fake.listSecretsArgsForCall)
}

func (fake *FakeVault) ListSecretsCalls(stub func(string) (*api.Secret, error)) {
	fake.listSecretsMutex.Lock()
	defer fake.listSecretsMutex.Unlock()
	fake.ListSecretsStub = stub
}

func (fake *FakeVault) ListSecretsArgsForCall(i int) string {
	fake.listSecretsMutex.RLock()
	defer fake.listSecretsMutex.RUnlock()
	argsForCall := fake.listSecretsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeVault) ListSecretsReturns(result1 *api.Secret, result2 error) {
	fake.listSecretsMutex.Lock()
	defer fake.listSecretsMutex.Unlock()
	fake.ListSecretsStub = nil
	fake.listSecretsReturns = struct {
		result1 *api.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) ListSecretsReturnsOnCall(i int, result1 *api.Secret, result2 error) {
	fake.listSecretsMutex.Lock()
	defer fake.listSecretsMutex.Unlock()
	fake.ListSecretsStub = nil
	if fake.listSecretsReturnsOnCall == nil {
		fake.listSecretsReturnsOnCall = make(map[int]struct {
			result1 *api.Secret
			result2 error
		})
	}
	fake.listSecretsReturnsOnCall[i] = struct {
		result1 *api.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeVault) SetNamespace(arg1 string) {
	fake.setNamespaceMutex.Lock()
	fake.setNamespaceArgsForCall = append(fake.setNamespaceArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.SetNamespaceStub
	fake.recordInvocation("SetNamespace", []interface{}{arg1})
	fake.setNamespaceMutex.Unlock()
	if stub != nil {
		fake.SetNamespaceStub(arg1)
	}
}

func (fake *FakeVault) SetNamespaceCallCount() int {
	fake.setNamespaceMutex.RLock()
	defer fake.setNamespaceMutex.RUnlock()
	return len(fake.setNamespaceArgsForCall)
}

func (fake *FakeVault) SetNamespaceCalls(stub func(string)) {
	fake.setNamespaceMutex.Lock()
	defer fake.setNamespaceMutex.Unlock()
	fake.SetNamespaceStub = stub
}

func (fake *FakeVault) SetNamespaceArgsForCall(i int) string {
	fake.setNamespaceMutex.RLock()
	defer fake.setNamespaceMutex.RUnlock()
	argsForCall := fake.setNamespaceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeVault) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addressMutex.RLock()
	defer fake.addressMutex.RUnlock()
	fake.allMountsMutex.RLock()
	defer fake.allMountsMutex.RUnlock()
	fake.allPoliciesMutex.RLock()
	defer fake.allPoliciesMutex.RUnlock()
	fake.getPolicyInfoMutex.RLock()
	defer fake.getPolicyInfoMutex.RUnlock()
	fake.getSecretInfoMutex.RLock()
	defer fake.getSecretInfoMutex.RUnlock()
	fake.listNamespacesMutex.RLock()
	defer fake.listNamespacesMutex.RUnlock()
	fake.listNestedSecretsMutex.RLock()
	defer fake.listNestedSecretsMutex.RUnlock()
	fake.listSecretsMutex.RLock()
	defer fake.listSecretsMutex.RUnlock()
	fake.setNamespaceMutex.RLock()
	defer fake.setNamespaceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVault) recordInvocation(key string, args []interface{}) {
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

var _ watcher.Vault = new(FakeVault)
