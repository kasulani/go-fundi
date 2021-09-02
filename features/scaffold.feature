Feature: Scaffold a new project directory structure

  Scenario: Scaffold command exits with code 0
    Given I have "a good fundi file"
    When I execute the cli command
    """
    fundi scaffold -f {{.File}}
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

  Scenario: Scaffold command exits with code 1
    Given I have "a bad fundi file"
    When I execute the cli command
    """
    fundi scaffold -f {{.File}}
    """
    Then I must get an exit code 1
