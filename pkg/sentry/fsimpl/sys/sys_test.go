// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sys_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/context/contexttest"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs/kernfstest"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/sys"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

func newTestSystem(t *testing.T) *kernfstest.System {
	ctx := contexttest.Context(t)
	creds := auth.CredentialsFromContext(ctx)
	v := vfs.New()
	v.MustRegisterFilesystemType("sysfs", sys.FilesystemType{})
	mns, err := v.NewMountNamespace(ctx, creds, "", "sysfs", &vfs.GetFilesystemOptions{})
	if err != nil {
		t.Fatalf("Failed to create new mount namespace: %v", err)
	}
	return kernfstest.NewSystem(ctx, t, v, mns)
}

func TestReadCPUFile(t *testing.T) {
	s := newTestSystem(t)
	expected := fmt.Sprintf("0-%d", sys.MaxCPUCores-1)

	for _, fname := range []string{"online", "possible", "present"} {
		pop := s.PathOpAtRoot(fmt.Sprintf("devices/system/cpu/%s", fname))
		fd, err := s.VFS.OpenAt(s.Ctx, s.Creds, &pop, &vfs.OpenOptions{})
		if err != nil {
			t.Fatalf("OpenAt(pop:%+v) = %+v failed: %v", pop, fd, err)
		}
		defer fd.DecRef()
		content, err := s.ReadToEnd(fd)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if diff := cmp.Diff(expected, content); diff != "" {
			t.Fatalf("Read returned unexpected data:\n--- want\n+++ got\n%v", diff)
		}
	}
}

func TestSysRootContainsExpectedEntries(t *testing.T) {
	s := newTestSystem(t)
	pop := s.PathOpAtRoot("/")
	s.AssertDirectoryContains(&pop, map[string]kernfstest.DirentType{
		"block":    linux.DT_DIR,
		"bus":      linux.DT_DIR,
		"class":    linux.DT_DIR,
		"dev":      linux.DT_DIR,
		"devices":  linux.DT_DIR,
		"firmware": linux.DT_DIR,
		"fs":       linux.DT_DIR,
		"kernel":   linux.DT_DIR,
		"module":   linux.DT_DIR,
		"power":    linux.DT_DIR,
	})
}
