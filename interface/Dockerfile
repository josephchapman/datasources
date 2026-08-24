FROM python:latest
WORKDIR /usr/src/app
RUN curl -LsSf https://astral.sh/uv/install.sh | env UV_INSTALL_DIR="/usr/local/bin" sh && \
    uv pip install --system --no-cache-dir \
    bash_kernel \
    influxdb-client \
    ipykernel \
    jupyterlab\
    langchain \
    langchain-community \
    langchain_redis \
    langchain-openai \
    langgraph \
    matplotlib \
    numpy \
    scipy \
    "starlette<0.40.0" \
    streamlit
RUN python -m ipykernel install --sys-prefix --name python3 --display-name "Python 3 (ipykernel)" && \
    python -m bash_kernel.install
EXPOSE 8501
EXPOSE 8888