# Knowledge Base: MCP Server for LLM Integration

## Purpose
This project implements an MCP (Model Context Protocol) server for integration with LLMs (Large Language Models). Its main goal is to manage a markdown-based knowledge base located in the `/doc` directory, exposing tools and resources to LLMs via stdio and JSON-RPC.

## Key Requirements
- **Transport & Protocol:** Communicate over stdio using JSON-RPC 2.0. Only JSON-RPC responses are sent to stdout; all logs go to stderr.
- **Tool/Resource Registration:** Expose tools (e.g., create/edit/validate markdown files, refactor folders, aggregate chat context) and resources (list/read files) to LLMs. Tool names must use only underscores (no dashes) for maximum compatibility with Cursor.
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

## MCP Server Design & Cursor Integration
- The MCP server is a backend service for LLM integration, exposing tools/resources via JSON-RPC 2.0 over stdio.
- The server must respond to `initialize` with a `serverInfo`, `protocolVersion`, and `capabilities` object in the result.
- The server must respond to `listOfferings` (and `list_tools`) with a top-level `offerings` array, each tool including `name`, `description`, and a `parameters` JSON schema.
- All logs must go to stderr; only JSON-RPC responses are sent to stdout.
- The process must stay alive after stdin closes (block with `select {}`) to match Cursor's expectations.
- Tool names must use underscores only; dashes can cause Cursor to ignore the tool.
- Place the MCP server at the top of the `mcp.json` config and, if troubleshooting, disable other servers to avoid Cursor's internal tool/server limit.
- If Cursor still fails to recognize the server, restart Cursor, try again, and use the MCP Inspector to verify server output.
- These steps are based on current community and official knowledge; Cursor's MCP support is evolving and may require further tweaks as new versions are released.

## Next Steps
- Continue implementing and testing each tool/resource endpoint using TDD.
- Evolve tool/resource schemas and metadata as requirements become clearer.
- Keep this knowledge base updated as the project grows, especially with integration lessons and platform-specific quirks. 