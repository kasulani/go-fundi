# Contributing Guidelines for Fundi

Thank you for considering contributing to Fundi! Please follow these guidelines to make the process as smooth as possible for everyone involved.

## Code of Conduct
We have adopted a Code of Conduct that we expect all contributors to adhere to. Please read it before making any contributions.

## Reporting Bugs and Issues
If you find a bug or issue with Fundi, please open an issue in the project's issue tracker. Please provide as much detail as possible, including steps to reproduce the issue, expected behavior, and actual behavior.

## Contributing Code
If you would like to contribute code to Fundi, please follow these guidelines:

1. Fork the Project
2. Create your Feature [Branch](#branch-naming-convention-and-commit-message-format) (`git checkout -b major/AmazingFeature`)
3. [Commit](#commit-message-guidelines) your Changes (`git commit -m 'major: Add some AmazingFeature'`)
4. Push to the Branch (`git push origin minor/AmazingFeature`)
5. Open a Pull Request

Before submitting a pull request, please ensure that your code adheres to the following guidelines:

Follow the existing code style and formatting. Write clear and concise [commit messages](#commit-message-guidelines). Include tests for any new
functionality or bug fixes. Ensure that your changes do not break any existing functionality. By contributing to Fundi,
you agree to license your contributions under the terms of the MIT License.

If you have any questions or issues, please open an issue in this repository.

### Code Review
All code contributions will be reviewed by a maintainer of the project. The maintainer may provide feedback or request changes to the code. Please be patient during the review process.

## Branch naming convention and commit message format

The branch naming convention and commit message format are as follows:
Branch naming convention: `type/branch-name`
Commit message format: `type: commit message`
The `type` can be one of the following:

- `minor`: Minor changes or a new feature
- `major`: Major changes or breaking change
- `patch`: A bug fix
- `test`: Adding tests
- `chore`: Maintenance tasks such as updating dependencies or configuration files or bootstrap code

### Commit Message Guidelines

To maintain consistency and clarity in our project history, all commit messages should follow the format: `type: commit message`

#### Accepted Types
- **minor**: For minor changes or new features.
- **major**: For major changes or breaking changes.
- **patch**: For bug fixes.
- **test**: For adding or modifying tests.
- **chore**: For maintenance tasks, such as updating dependencies or configuration files or bootstrap code.

#### Examples
- `minor: Add user profile page`
- `major: Remove support for deprecated API`
- `patch: Fix null pointer exception in login handler`
- `test: Add unit tests for authentication module`
- `chore: Update CI configuration`

#### Why This Matters
Using a consistent format for commit messages helps:
- Easily identify the purpose and impact of each commit.
- Streamline the release process by automatically generating changelogs.
- Improve collaboration and understanding among team members.

Make sure to follow these guidelines for every commit to keep our project history clean and meaningful!

## License
By contributing to Fundi, you agree to license your contributions under the terms of the MIT License.

If you have any questions or issues, please open an issue in this repository.