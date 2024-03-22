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
dagger call install --dir https://github.com/dagger/dagger
```

## Generate

Generate the Schema for the given Dagger version.

```bash
dagger call generate --dir https://github.com/dagger/dagger
```

## Init module

```bash
mkdir java-sdk
cd java-sdk

dagger init
dagger develop --sdk go
```
