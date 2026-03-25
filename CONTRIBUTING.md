# Contributing to Daedalus

First off, thank you for contributing to Daedalus. This is a proprietary monorepo. Access is strictly restricted to authorized core team members. 

To maintain the highest software engineering standards and ensure our Jenkins CI/CD pipelines run smoothly, please adhere to the following guidelines.

## 1. Branching Strategy (GitFlow)
We use a structured branching model:
* `main` : Production-ready code. Do not push directly to this branch.
* `develop` : Integration branch for upcoming releases.
* `feature/<ticket-id>-<short-description>` : For new features (e.g., `feature/PB-051-auth-service`).
* `fix/<ticket-id>-<short-description>` : For bug fixes.

## 2. Commit Message Convention
We strictly follow [Conventional Commits](https://www.conventionalcommits.org/):
* `feat:` A new feature.
* `fix:` A bug fix.
* `docs:` Documentation only changes.
* `test:` Adding missing tests or correcting existing tests.
* `chore:` Changes to the build process, CI/CD, or auxiliary tools.

*Example: `feat(procurement-agent): integrate LangChain sourcing logic`*

## 3. Pull Request (PR) Process
1. Push your feature branch to the repository.
2. Open a Pull Request against the `develop` branch.
3. Ensure your PR description links to the relevant Scrum Product Backlog ID.
4. **Quality Gates:** - Code must compile without warnings.
   - **Test Coverage MUST be ≥ 80%** (Unit & Integration tests). Jenkins will automatically fail the build if this threshold is not met.
5. Require at least one peer review approval before merging.

## 4. Local Development
Please refer to the `README.md` inside each specific microservice folder (`/services/*`) for language-specific setup instructions (Golang, FastAPI, Node.js).