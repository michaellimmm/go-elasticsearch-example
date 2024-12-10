# Dockerfile
FROM docker.elastic.co/elasticsearch/elasticsearch:8.12.0

# Install the analysis-kuromoji and analysis-icu plugins
RUN elasticsearch-plugin install --batch analysis-kuromoji
RUN elasticsearch-plugin install --batch analysis-icu