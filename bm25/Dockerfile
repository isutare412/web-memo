FROM python:3.13-slim AS builder

RUN pip install uv

WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen

FROM python:3.13-slim AS runner

WORKDIR /app
COPY --from=builder /app/.venv /app/.venv
COPY main.py .

CMD ["/app/.venv/bin/uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]
