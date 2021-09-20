Feature: Add files to an existing project directory structure

  Scenario: Generate files skip templates exits with code 0
    Given I have "a good fundi file"
    When I execute the cli command
    """
    fundi generate scaffold --use-config {{.File}}
    """
    Then I must get an exit code 0
    When I execute the cli command
    """
    fundi generate files --skip-templates --use-config {{.File}}
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
    docker-compose.yml
    docker-files
    docs
    features
    pkg
    """
    When I execute the cli command
    """
    ls funditest/pkg
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    app
    behaviour
    domain
    """
    When I execute the cli command
    """
    ls funditest/pkg/app
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    doc.go
    """

  Scenario: Generate files from templates exits with code 0
    Given I have "a good fundi file"
    When I execute the cli command
    """
    fundi generate scaffold --use-config {{.File}}
    """
    Then I must get an exit code 0
    When I execute the cli command
    """
    fundi generate files --use-config {{.File}}
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
    docker-compose.yml
    docker-files
    docs
    features
    pkg
    """
    When I execute the cli command
    """
    ls funditest/pkg
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    app
    behaviour
    domain
    """
    When I execute the cli command
    """
    ls funditest/pkg/app
    """
    Then I must get an exit code 0
    And I must get a command output
    """
    doc.go
    """
    And file "doc.go" has contents
    """
    // Package app provides some functionality.
    package app

    """

  Scenario: Generate files skip templates exits with code 1
    Given I have "a bad fundi file"
    When I execute the cli command
    """
    fundi generate scaffold --use-config {{.File}}
    """
    Then I must get an exit code 1
    When I execute the cli command
    """
    fundi generate files --skip-templates --use-config {{.File}}
    """
    Then I must get an exit code 1

  Scenario: Generate files from templates exits with code 1
    Given I have "a bad fundi file"
    When I execute the cli command
    """
    fundi generate scaffold --use-config {{.File}}
    """
    Then I must get an exit code 1
    When I execute the cli command
    """
    fundi generate files --use-config {{.File}}
    """
    Then I must get an exit code 1

  Scenario: No project structure
    Given I have "a good fundi file"
    When I execute the cli command
    """
    fundi generate files --skip-templates --use-config {{.File}}
    """
    Then I must get an exit code 1
