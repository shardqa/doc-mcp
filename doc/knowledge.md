# Knowledge Base: MCP Server for LLM Integration

## Purpose
This project implements an MCP (Model Context Protocol) server for integration with LLMs (Large Language Models). Its main goal is to manage a markdown-based knowledge base located in the `/doc` directory, exposing tools and resources to LLMs via stdio and JSON-RPC.

## Key Requirements
- **Transport & Protocol:** Communicate over stdio using JSON-RPC 2.0.
- **Tool/Resource Registration:** Expose tools (e.g., create/edit/validate markdown files, refactor folders, aggregate chat context) and resources (list/read files) to LLMs.
- **Validation Rules:**
  - Every markdown file must have at least two internal links to other markdown files.
  - No markdown file may exceed 100 lines.
  - After every file edit or creation, `markdownlint --fix` is run.
  - No folder may contain more than 10 items (files or subfolders); if exceeded, the server should auto-refactor by creating subfolders and moving content.
  - All documentation is kept in the root-level `/doc` folder.
  - Validation is warn-only (does not block actions).
- **Aggregation:** Chat context and preferences are aggregated and appended/merged as markdown into the documentation.

## Project Workflow
- **TDD:** All features are developed using Test-Driven Development (TDD): write a failing test, implement the feature, make the test pass, then refactor.
- **Task Tracking:** Tasks are tracked in `TODO.md` as a markdown task list. Completed tasks are moved to `COMPLETED.md`.
- **Knowledge Maintenance:** This file (`doc/knowledge.md`) is incrementally updated to keep project knowledge relevant and up to date. Refactoring and documentation are ongoing priorities.

## Implementation Stack
- **Language:** Go
- **Testing:** Standard Go `testing` package with `testify` for expressive assertions and helpers.
- **Tool Registration:** Tools are statically registered for now, but may become dynamic/configurable in the future.

## Design & Process Notes
- The MCP server is not a CLI tool; it is a backend service for LLM integration, exposing tools/resources via a protocol.
- All validation and organization rules are implemented as part of the tool logic, with warnings returned in responses.
- The server is robust to closed stdin and designed for long-lived operation.
- The project is structured for extensibility, with clear separation of concerns and a focus on maintainability.

## Next Steps
- Continue implementing and testing each tool/resource endpoint using TDD.
- Evolve tool/resource schemas and metadata as requirements become clearer.
- Keep this knowledge base updated as the project grows. 