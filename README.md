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
       <li><a href="#generate-project-structure">Generate project structure</a></li>
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

### Generate project structure

The tool can generate a project structure based on a YAML file that specifies the directory structure. The generated
project can be customized by specifying `go` templates for files in each directory.

### Generate project files from local templates

The tool supports reading templates from a local directory on the user's machine and generate project files based on 
these templates. This allows developers to create their own templates and use them with the tool.

### Generate project files from remote templates

The tool can read templates from a remote Git repository and generate project files based on these templates. This feature 
enables developers to use templates hosted on GitHub or any other Git hosting platform.

<!-- INSTALLATION -->

## Installation

To install Fundi, you can use the go install command:

```bash
$ go install github.com/kasulani/go-fundi
```

This will install **Fundi** in your `$GOBIN` directory, which is typically located at `$GOPATH/bin`. Make sure that this
directory is added to your `$PATH` environment variable so that you can run the `fundi` command from anywhere in your
terminal.

If you encounter any errors during the installation process, try running

```bash
$ go clean -i github.com/kasulani/go-fundi
```

to clean any existing build artifacts before running the go install command again.

<!-- USAGE EXAMPLES -->

## Usage

To use **Fundi**, create a YAML file that specifies the project structure and any `go` templates that specify the contents of the files in
each directory.

**Example YAML file**

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

When you execute the generate command with the above configuration file, **Fundi** will create a project directory structure,
add files to the project directories. The file contents in those files will be generated based on the templates provided.
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

Edit the `example yaml file` and remove the templates from the configuration file.
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

## Roadmap

See the [open issues](https://github.com/kasulani/go-fundi/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->

## Contributing

We welcome contributions from the community! If you would like to contribute to **Fundi**, please follow these guidelines:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Before submitting a pull request, please ensure that your code adheres to the following guidelines:

Follow the existing code style and formatting. Write clear and concise commit messages. Include tests for any new
functionality or bug fixes. Ensure that your changes do not break any existing functionality. By contributing to Fundi,
you agree to license your contributions under the terms of the MIT License.

If you have any questions or issues, please open an issue in this repository.

Contributing Guidelines

For more detailed information on contributing to Fundi, please see the CONTRIBUTING file in this repository.

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

If you have any questions, suggestions or feedback, feel free to contact me via [email](mailto:kasulani@gmail.com). You can also find me
on [LinkedIn](https://ug.linkedin.com/in/kasulani). I'm always happy to hear from you!

<!-- CONCLUSION -->

## Conclusion

Fundi is a powerful tool for generating project directories and files in a flexible and customizable way. It supports
local and remote templates and provides a command line interface for easy use.