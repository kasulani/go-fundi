Feature: Scaffold a new project directory structure

  Scenario: Scaffold command exits with code 0
    Given I execute the cli command
    """
    fundi scaffold -f /project/pkg/behaviour/.bdd.test.fundi.yaml
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
    Given I execute the cli command
    """
    fundi scaffold -f /project/pkg/behaviour/.non-existing.fundi.yaml
    """
    Then I must get an exit code 1
