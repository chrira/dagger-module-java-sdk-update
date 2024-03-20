// A generated module for JavaSdk functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
)

type JavaSdk struct {
	Ctr *Container
}

// Returns a container that echoes whatever string argument is provided
func (m *JavaSdk) ContainerEcho(ctx context.Context, stringArg string) (string, error){
	return dag.Container().
				From("registry.access.redhat.com/ubi9/openjdk-17:1.18-1").
				WithExec([]string{"gzip", stringArg}).
				Stdout(ctx)
}


// Returns lines that match a pattern in the files of the provided Directory
func (m *JavaSdk) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

func (m *JavaSdk) GetJDK() *JavaSdk {
	m.Ctr = dag.Container().
		From("registry.access.redhat.com/ubi9/openjdk-17:1.18-1").
			WithUser("root").
			WithExec([]string{"microdnf", "-y", "install", "gzip"}).
			WithUser("185")
	return m
}

func (m *JavaSdk) InstallDagger() *JavaSdk {
	m.Ctr = m.Ctr.
		WithUser("root").
		WithWorkdir("/usr/local").
		WithExec([]string{"curl", "-L", "-o", "install.sh", "https://dl.dagger.io/dagger/install.sh"}).
		WithExec([]string{"chmod", "+x", "./install.sh"}).
		WithExec([]string{"./install.sh"}).
		WithUser("185")
	//	WithExec([]string{"DAGGER_VERSION=0.9.7", "./install.sh"})
	return m


// cd /usr/local
// curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION=0.9.7 sh


}

func (m *JavaSdk) DaggerVersion(ctx context.Context, container *Container) (string, error){
	return container.
				WithExec([]string{"dagger", "version"}).
				Stdout(ctx)
}

func (m *JavaSdk) CI(ctx context.Context) (string, error){
	return m.GetJDK().InstallDagger().Ctr.
				WithExec([]string{"/usr/local/bin/dagger", "version"}).
				Stdout(ctx)
}
