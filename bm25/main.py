import logging

from fastapi import FastAPI
from fastembed import SparseTextEmbedding
from pydantic import BaseModel


class _HealthCheckFilter(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        return record.getMessage().find("/health") == -1


logging.getLogger("uvicorn.access").addFilter(_HealthCheckFilter())

app = FastAPI()

model = SparseTextEmbedding(model_name="Qdrant/bm25")


class EmbedSparseRequest(BaseModel):
    inputs: list[str]


class SparseVector(BaseModel):
    indices: list[int]
    values: list[float]


@app.post("/embed-sparse")
def embed_sparse(req: EmbedSparseRequest) -> list[SparseVector]:
    embeddings = list(model.embed(req.inputs))
    return [
        SparseVector(
            indices=e.indices.tolist(),
            values=e.values.tolist(),
        )
        for e in embeddings
    ]


@app.get("/health")
def health():
    return {"status": "ok"}
