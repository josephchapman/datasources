FROM python:latest
WORKDIR /usr/src/app
RUN pip install --no-cache-dir \
    influxdb-client \
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
EXPOSE 8501
ENTRYPOINT ["streamlit", "run", "🏠_Home.py", "--server.port=8501", "--server.address=0.0.0.0"]