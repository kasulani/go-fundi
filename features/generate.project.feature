Feature: Generate project directories

  Scenario: using an invalid configuration file
    Given I have the following configuration
    """
    *#!%
    """
    When I execute the cli command
    """
    fundi generate-cmd -f {{.File}}
    """
    Then I must get an exit code 1
    And I must get a command output
    """
    failed to unmarshal YAML data: yaml: did not find expected alphabetic or numeric character
    """

  Scenario: generate directories only
    Given I have the following configuration
    """
    metadata:
      output: "."
      templates: "./testdata"
      values: "./testdata/.values.yml"
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
            files:
              - name: domain.go
                template: domain.go.tmpl
    """
    When I execute the cli command
    """
    fundi generate-cmd --directories-only -f {{.File}}
    """
    Then I must get an exit code 0
    When I execute the cli command
    """
    ls funditest
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    cmd
    internal
    """

  Scenario: generate empty files
    Given I have the following configuration
    """
    metadata:
      output: "."
      templates: "./testdata"
      values: "./testdata/.values.yml"
    directories:
      - name: funditest
        files:
          - name: README.md
        directories:
          - name: cmd
          - name: internal
    """
    When I execute the cli command
    """
    fundi generate-cmd -f {{.File}}
    """
    Then I must get an exit code 0
    When I execute the cli command
    """
    ls funditest
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    README.md
    cmd
    internal
    """

  Scenario: generate all
    Given I have the following configuration
    """
    metadata:
      output: "."
      templates: "./testdata"
      values: "./testdata/.values.yml"
    directories:
      - name: funditest
        files:
          - name: README.md
        directories:
          - name: cmd
            files:
              - name: main.go
                template: main.go.tmpl
          - name: internal
    """
    And a "main.go.tmpl" file with the following contents
    """
    // Package {{.package}} has the entry point into your app.
    package {{.package}}
    """
    And a ".values.yml" file with the following contents
    """
    main.go.tmpl:
      package: main
    """
    When I execute the cli command
    """
    fundi generate-cmd -f {{.File}}
    """
    Then I must get an exit code 0
    When I execute the cli command
    """
    cat funditest/cmd/main.go
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    // Package main has the entry point into your app.
    package main
    """
