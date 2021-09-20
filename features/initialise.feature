Feature: Scaffold and create files using a single command

  Scenario: Init command exits with code 0
    Given I have "a good fundi file"
    When I execute the cli command
    """
    fundi init --skip-templates --use-config {{.File}}
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

  Scenario: Init command exits with code 1
    Given I have "a bad fundi file"
    When I execute the cli command
    """
    fundi init --skip-templates --use-config {{.File}}
    """
    Then I must get an exit code 1
