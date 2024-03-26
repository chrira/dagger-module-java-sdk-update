# Update Dagger JavaSDK

Experimental [Dagger.io](https://dagger.io) module for the Java SDK.

It helps generate generate artifacts for new Dagger engine version.

Artifacts:

* Schemas for code generation
* Library Jar

## Update Pom

```bash
dagger call -o . update --dir . --version "0.10.2"
```

## Install Dagger CLI

```bash
dagger call install-dagger
```

## Run Java Install

Install command returns a container.

```bash
dagger call install --dir https://github.com/dagger/dagger
```

Show install output:

```bash
dagger call install --dir https://github.com/dagger/dagger stdout
```

## Generate

Generate the Schema for the given Dagger version.

```bash
dagger call generate --dir https://github.com/dagger/dagger
```

Show generate output:

```bash
dagger call generate --dir https://github.com/dagger/dagger stdout
```

Export the generated schema:

```bash
dagger call generate --dir https://github.com/dagger/dagger export --path ./my-dagger-build
```

Export the generated schema for a specific version:

```bash
dagger call generate --version 0.10.1 --dir https://github.com/dagger/dagger directory --path ./target/generated-schema export --path ./my-schema
```

## Init module

```bash
mkdir java-sdk
cd java-sdk

dagger init
dagger develop --sdk go
```
