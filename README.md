# WebMemo

WebMemo is a collaborative note-taking application that supports Markdown
editing with real-time preview. It features web publishing, tag-based
organization, interactive checkboxes, image embedding
(using [Imageer](https://github.com/isutare412/imageer)), natural language search,
and team collaboration with approval workflows.

## Deployment Architecture

![web memo deployment architecture](docs/assets/deployment-diagram.drawio.png)

## Hybrid Semantic Search

WebMemo supports hybrid search that combines:

- Semantic vector search (embedding-based retrieval)
- BM25 keyword search (sparse retrieval)

Search results are fused with Reciprocal Rank Fusion (RRF), so both semantic
similarity and exact keyword signals contribute to ranking.

## Screenshots

### Memo Edit
Create and edit memos using a Markdown editor with Write/Preview tabs. Add tags
for organization, use the formatting toolbar for headings, lists, checkboxes,
and images.

![Memo Edit](docs/assets/memo-edit.png)

### Memo View
View rendered memos with interactive checkboxes, embedded images.

![Memo View](docs/assets/memo-view.png)

### Collaboration
Publish your memos to the web and collaborate with others.

![Collaboration](docs/assets/collaboration.png)
