// An experimental module to generate Dagger Java SDK artifacts.
//
// It helps generate generate artifacts for new [Dagger.io](https://dagger.io) engine version.
// Artifacts:
// * Schemas for code generation
// * Library Jar

package main

import (
	"context"
	"fmt"
)

type JavaSdk struct {
	Ctr *Container
}

func (m *JavaSdk) GetJDK() *JavaSdk {
	m.Ctr = dag.Container().
		From("registry.access.redhat.com/ubi9/openjdk-17:1.18-1").
			WithUser("root").
			WithExec([]string{"microdnf", "-y", "install", "gzip"}).
			WithUser("185")
	return m
}

// install Dagger
func (m *JavaSdk) InstallDagger(
	// +optional
	// +default="0.10.2"
	daggerVersion string,
) *JavaSdk {
	m.Ctr = m.Ctr.
		WithEnvVariable("DAGGER_VERSION", daggerVersion).
		WithUser("root").
		WithWorkdir("/usr/local").
		WithExec([]string{"curl", "-L", "-o", "install.sh", "https://dl.dagger.io/dagger/install.sh"}).
		WithExec([]string{"chmod", "+x", "./install.sh"}).
		WithExec([]string{"./install.sh"}).
		WithUser("185")
	return m
}

func (m *JavaSdk) DaggerVersion(ctx context.Context, container *Container) (string, error){
	return container.
				WithExec([]string{"dagger", "version"}).
				Stdout(ctx)
}

func (m *JavaSdk) CI(
	ctx context.Context,
	// +optional
	// +default="0.10.2"
	daggerVersion string,
) (string, error){
	return m.GetJDK().InstallDagger(daggerVersion).Ctr.
				WithExec([]string{"/usr/local/bin/dagger", "version"}).
				Stdout(ctx)
}

func (m *JavaSdk) Update(
	ctx context.Context,
	dir *Directory,
	version string,
	// +optional
	path string,
) (*Directory, error) {
	file :=fmt.Sprintf("%spom.xml", path)
	sedCommand := fmt.Sprintf("sed -i \"s#<daggerengine.version>devel</daggerengine.version>#<daggerengine.version>%s</daggerengine.version>#g\" %s", version, file)
	catCommand := fmt.Sprintf("cat %s | grep 'daggerengine.version'", file)
	return dag.Container().
		From("alpine@sha256:6457d53fb065d6f250e1504b9bc42d5b6c65941d57532c072d929dd0628977d0").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"sh", "-c", sedCommand}).
		WithExec([]string{"sh", "-c", catCommand}).
		Directory("/mnt/").
		Sync(ctx)
}

func (m *JavaSdk) Updates(ctx context.Context, dir *Directory, version string) (string, error) {
	sedCommand := fmt.Sprintf("sed -i \"s#<daggerengine.version>devel</daggerengine.version>#<daggerengine.version>%s</daggerengine.version>#g\" pom.xml", version)
	return dag.Container().
		From("alpine@sha256:6457d53fb065d6f250e1504b9bc42d5b6c65941d57532c072d929dd0628977d0").
		WithMountedDirectory("/mnt", dir).
		WithWorkdir("/mnt").
		WithExec([]string{"sh", "-c", sedCommand}).
		WithExec([]string{"sh", "-c", "cat pom.xml | grep 'daggerengine.version'"}).
		Stdout(ctx)
}

func (m *JavaSdk) Install(
	ctx context.Context,
	dir *Directory,
	// +optional
	// +default="0.10.2"
	daggerVersion string,
) (string, error){
	homeDir := "/home/default"
	srcDir := "/mnt/src"
	m2CacheDir := "/home/default/.m2"
	workDir := fmt.Sprintf("%s/sdk/java", srcDir)
	return m.GetJDK().
				InstallDagger(daggerVersion).
				Ctr.
				WithEnvVariable("HOME", homeDir).
				WithEnvVariable("M2_HOME", homeDir).
				WithEnvVariable("MAVEN_OPTS", "-Xdebug").
				WithEnvVariable("MVNW_VERBOSE", "true").
				WithMountedCache(m2CacheDir, dag.CacheVolume("maven-cache"), ContainerWithMountedCacheOpts{Owner: "185"}).
				WithMountedDirectory(srcDir, dir, ContainerWithMountedDirectoryOpts{Owner: "185"}).
				WithWorkdir(workDir).
				WithExec([]string{"java", "--version"}).
				WithExec([]string{"./mvnw", "--debug", "--version"}).
				WithExec([]string{"./mvnw", "install", "-pl", "dagger-codegen-maven-plugin"}).
				WithExec(
					[]string{"./mvnw", "-X", "-N", "dagger-codegen:generateSchema", "-Ddagger.bin=/usr/local/bin/dagger"},
					ContainerWithExecOpts{ExperimentalPrivilegedNesting: true},
				).
				Stdout(ctx)
}

// Generate the Schema for the given Dagger engine version.
//
// Example usage: `dagger call generate --dir https://github.com/dagger/dagger`
func (m *JavaSdk) Generate(
	ctx context.Context,
	dir *Directory,
	// +optional
	// +default="0.10.2"
	version string,
) (string, error) {
	updDir, err := m.Update(ctx, dir, version, "sdk/java/")
	//updDir, err := m.Update(ctx, dir.Directory("sdk/java"), version)
	if err != nil {
		return "", err
	}
	return m.Install(ctx, updDir, version)
}
