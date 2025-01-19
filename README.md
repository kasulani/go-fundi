<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
     <a href="#features">Features</a>
      <ul>
       <li><a href="#generate-project-structure-from-yaml">Generate project structure</a></li>
       <li><a href="#generate-project-files-from-local-templates">Generate project files from local templates</a></li>
       <li><a href="#generate-project-files-from-remote-templates">Generate project files from remote templates</a></li>
      </ul>
    </li>
    <li>
      <a href="#installation">Installation</a>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#conclusion">Conclusion</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

**Fundi** is a scaffolding and code generation tool built using `go`. The tool is designed to help developers create a
directory structure for their projects, and generate code within those directories using `go` templates. The tool reads
a YAML file that specifies the project structure and creates the directories and files according to that structure.

**Fundi** is swahili word that means "expert" or "specialist" in English. It's a commonly used word in East Africa to
refer to a skilled trades-person or artisan.

<!-- FEATURES -->

## Features

### Generate Project Structure from YAML

Fundi generates a project structure using a YAML configuration file that specifies the desired directory structure.
Files within the directories can be customized using `go` templates.

### Generate Project Files from Local Templates

Fundi can use templates stored on your local machine to create project files. This enables developers to define and
reuse custom templates tailored to specific project requirements.

### Generate Project Files from Remote Templates (Coming Soon)

Fundi can fetch templates from remote Git repositories (e.g., GitHub or other Git hosting platforms) and generate
project files. This simplifies collaboration by allowing teams to share and maintain templates in a centralized
repository.
<!-- INSTALLATION -->

## Installation
Follow these steps to install **Fundi** on your system:

---
### Step 1: Download the Release
1. Visit the [Releases Page](https://github.com/kasulani/go-fundi/releases).
2. Download the appropriate archive for your operating system and architecture:
    - **Linux**: `fundi_<version>_Linux_x86_64.tar.gz`
    - **Windows**: `fundi_<version>_Windows_x86_64.zip`
    - **macOS**: `fundi_<version>_Darwin_x86_64.tar.gz`

---

### Step 2: Extract the Archive
- **Linux/macOS**:
```bash
  tar -xvf fundi_<version>_<os>_<arch>.tar.gz
```

- **Windows**:
  - Right-click on the downloaded ZIP file and select "Extract All".
  - Choose a destination folder and click "Extract".

---

### Step 3: Move the Binary to a Directory in Your PATH
- **Linux/macOS**:
  - Move the `fundi` binary to a directory in your PATH (e.g., `/usr/local/bin`):
  ```bash
    sudo mv fundi /usr/local/bin
  ```
  - Verify the installation by running:
  ```bash
    fundi --version
  ```
- **Windows**:
  - Move the `fundi.exe` binary to a directory in your PATH (e.g., `C:\Windows\System32`).
  - Open a new Command Prompt window and run:
    ```bash
    fundi --version
    ```

---

### Install from source using go install
To install Fundi, you can use the go install command:

```bash
go install github.com/kasulani/go-fundi/cmd/fundi@latest
```

This will install **Fundi** in your `$GOBIN` directory, which is typically located at `$GOPATH/bin`. Make sure that this
directory is added to your `$PATH` environment variable so that you can run the `fundi` command from anywhere in your
terminal.

If you encounter any errors during the installation process, try running

```bash
#!/bin/bash
# Clean Go build and module cache for fundi
go clean -cache -modcache -i -r github.com/kasulani/go-fundi

# Remove binaries from GOPATH or GOBIN
rm -f $(go env GOBIN)/fundi
rm -f $(go env GOPATH)/bin/fundi

# Remove local artifact directories
rm -rf ./bin ./dist
```

to clean any existing build artifacts before running the go install command again.

<!-- USAGE EXAMPLES -->

## Usage

To use **Fundi**, create a YAML file that specifies the project structure and any `go` templates that specify the
contents of the files in
each directory.

**Example YAML configuration file**

```yaml
metadata:
  output: "."
  templates: "./templates"
  values: "./values.yml"
directories:
  - name: funditest
    files:
      - name: README.md
        template: readme.md.tmpl
    directories:
      - name: cmd
        files:
          - name: main.go
            template: main.go.tmpl
      - name: internal
        skip: true
        files:
          - name: domain.go
            template: domain.go.tmpl
```

**Generate the project directories and files using templates:**

When you execute the generate command with the above configuration file, **Fundi** will create a project directory
structure,
add files to the project directories. The file contents in those files will be generated based on the templates
provided.

```bash
$ fundi generate -f /path/to/yaml/file.yaml
```

**Generate only the project directories:**

Edit the `example yaml file` and remove the files from the configuration file.

```yaml
metadata:
  output: "."
  templates: "./testdata"
  values: "./testdata/.values.yml"
directories:
  - name: funditest
    directories:
      - name: cmd
      - name: internal
```

When you execute the command below, you will only have directories created without any files.

```bash
$ fundi generate -f /path/to/yaml/file.yaml
```

**Generate only the project files:**

Edit the `example yaml configuration file` and remove the templates from the configuration file.

```yaml
metadata:
  output: "."
  templates: "./templates"
  values: "./values.yml"
directories:
  - name: funditest
    files:
      - name: README.md
    directories:
      - name: cmd
        files:
          - name: main.go
      - name: internal
        skip: true
        files:
          - name: domain.go
```

When you execute the command below, you will only have directories and empty files created.

```bash
$ fundi generate -f /path/to/yaml/file.yaml
```

<!-- CONTRIBUTING -->

## Contributing

For more detailed information on contributing to Fundi, please see the [CONTRIBUTING.md](https://github.com/kasulani/go-fundi/blob/master/CONTRIBUTING.md) file in this repository.

<!-- VERSIONING -->
## Versioning

This project follows **Semantic Versioning 2.0.0** guidelines to manage releases and version numbers. Semantic
versioning ensures that version numbers convey meaning about the underlying changes, helping developers and users know
what to expect.

### Version Format

`vMAJOR.MINOR.PATCH`

- **`MAJOR`**: Increased when there are incompatible changes or breaking changes to the API or functionality.
    - Example: Removing or changing a feature that existing users depend on.
    - Example Version: `v2.0.0` (Breaking change introduced from `v1.x.x`).
- **`MINOR`**: Increased when new features are added in a backward-compatible manner.
    - Example: Adding a new feature that doesn’t interfere with existing functionality.
    - Example Version: `v1.2.0` (New features added to `v1.1.x`).
- **`PATCH`**: Increased when backward-compatible bug fixes are introduced.
    - Example: Fixing a typo, correcting a bug, or updating documentation.
    - Example Version: `v1.1.1` (Bug fix for `v1.1.0`).

### How Versioning Works in This Project

1. **Branch Prefixes**:
    - Branch naming conventions indicate the type of version bump:
        - `minor/*` → Minor version bump.
        - `patch/*` → Patch version bump.
        - `major/*` → Major version bump.

2. **Automated Tagging and Releases**:
    - Tags are created automatically based on the type of changes introduced.
    - The release workflow ensures that every new tag is packaged and published
      using [GoReleaser](https://goreleaser.com).

3. **Commit History and Changelog**:
    - Version numbers are mapped to the changelog to reflect new features, bug fixes, or breaking changes in every
      release.

### Examples

- `v1.0.0` → Initial release with the first stable version.
- `v1.1.0` → New features added.
- `v1.1.1` → Minor bug fix.

### Learn More

For a detailed explanation of semantic versioning, see
the [Semantic Versioning Specification (SemVer 2.0.0)](https://semver.org/).

## Roadmap

See the [open issues](https://github.com/kasulani/go-fundi/issues) for a list of proposed features (and known issues).

<!-- LICENSE -->

## License

**Fundi** is licensed under the MIT License. You can find a copy of the license at the
following [link](https://opensource.org/licenses/MIT).

In summary, the MIT License grants you permission to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of Fundi, subject to certain conditions. These conditions include including the license notice and
disclaimer in all copies or substantial portions of the software.

We believe that open-source software is important for the advancement of technology and welcome contributions from the
community. If you would like to contribute to Fundi, please read our contributing guidelines in the CONTRIBUTING file of
this repository.

<!-- CONTACT -->

## Contact

If you have any questions, suggestions or feedback, feel free to contact me via [email](mailto:kasulani@gmail.com). You
can also find me
on [LinkedIn](https://ug.linkedin.com/in/kasulani). I'm always happy to hear from you!

<!-- CONCLUSION -->

## Conclusion

Fundi is a powerful tool for generating project directories and files in a flexible and customizable way. It supports
local and remote templates and provides a command line interface for easy use.