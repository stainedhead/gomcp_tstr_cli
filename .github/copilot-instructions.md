# copilot-instructions.md
# This file contains instructions for Copilot to generate code for the project.
- the code should follow best practices for Go development, including idiomatic Go code, error handling, and performance considerations.
- a .gitignore file should be used to exclude unnecessary files from the repository.
- the executable should be listed in the .gitignore file to avoid committing built binaries.
- a makefile should be used to automate common tasks such as building, testing, and running the project.
- a directory named "internal" should be created to hold internal packages that are not meant to be used by external consumers.
- a directory named "tests" should be created to hold test files and test data.
- temporary files should be stored in a "tmp" directory, which should be excluded from the repository using .gitignore.
- temporary scripts or tests should be stored in tmp and cleaned up when no longer needed.
- string or numeric constants should be stored in a "constants" package to avoid duplication and ensure consistency.


# instructions on documentation
- documentation should be written in markdown format and stored in the "documentation" directory.
- product documentation should be used to understand context and features of the project. 
- ./documentation/product-overview.md should be referenced to understand the product and its features.
- ./readme.md should be used to understand the project structure and how to run the project.
- ./readme.md, ./documentation/product-overview.md, and ./documentation/technical-details.md should be updated as required to keep code and documentation in sync.
- a directory named "prompts" should be created to hold prompt files used for generating code or documentation.


# instructions on SDLC and Defintion of Done
- when changes are made to the product documentation, the code should be updated to reflect those changes.
- the code should be written in a way that is easy to understand and maintain.
- the code should be modular and reusable.
- testing code should be consolidated into a tests directory, and be used to validate the functionality of the code.
- unit testing code should be written to ensure that individual components work as expected.
- the code should be well-documented and include comments where necessary.
- snyk is available as a CLI, and should be used during code change testing to ensure security vulnerabilities are identified and addressed.
- logging, configuration, flags, and other aspects should be implemented consistently throughout the codebase.
- unit tests, code linting, and other checks should be run before code is committed to the repository.
- a final build of the code should be done before changes are merged into the main branch.


