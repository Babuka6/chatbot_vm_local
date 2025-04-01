# chatbot_vm_local
Performance Evaluation of a Kubernetes-Deployed Chatbot Microservice

This project evaluates the performance of a custom chatbot and TF-IDF service built with Go and Python, deployed using gRPC and REST APIs. The experiment includes setting up a Python-based chatbot using Hugging Face models, FastAPI, and a TF-IDF service to retrieve relevant documents based on user queries, feeding them into the chatbot’s question-answering system. The system is deployed on a virtual machine with an 8-core CPU, containerized with Docker, and orchestrated using Kubernetes.

The primary goal is to identify bottleneck services and assess how vertical and horizontal scaling impact the performance of computationally intensive AI applications. By emulating real-world user interactions with the wrk benchmarking tool, the study highlights how CPU-bound limitations can restrict scalability and emphasizes the importance of GPU integration for AI workloads.

Refer to the README.pdf file to see project summary and full analysis! 
