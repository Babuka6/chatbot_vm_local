# Use official Python image
FROM python:3.11-slim

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Create working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

RUN pip install "numpy<2"
# RUN pip install numpy==1.17

# Copy the proto-generated Python files
COPY retrievedoc_pb2.py retrievedoc_pb2_grpc.py ./

# Copy the Python source code
COPY chatbot.py ./

# Expose the FastAPI port
EXPOSE 5000

# Run the application
CMD ["uvicorn", "chatbot:app", "--host", "0.0.0.0", "--port", "5000"]
