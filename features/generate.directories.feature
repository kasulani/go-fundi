Feature: Generate project directories

  Scenario: using a valid configuration file
    Given I have the following configuration
    """
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
            files:
              - name: domain.go
                template: domain.go.tmpl
    """
    When I execute the cli command
    """
    fundi generate-directories -f {{.File}}
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

  Scenario: using invalid configuration file
    Given I have the following configuration
    """
    *#!%
    """
    When I execute the cli command
    """
    fundi generate-directories -f {{.File}}
    """
    Then I must get an exit code 1
    And I must get a command output
    """
    failed to unmarshal YAML data: yaml: did not find expected alphabetic or numeric character
    """
