from fastapi import FastAPI
from pydantic import BaseModel
import grpc
import retrievedoc_pb2
import retrievedoc_pb2_grpc
from transformers import pipeline

app = FastAPI()

# Initialize the QA model
oracle = pipeline(model="deepset/roberta-base-squad2")

# Define a Pydantic model for the request body
class ChatRequest(BaseModel):
    query: str

def retrieve_document(query):
    with grpc.insecure_channel('tfidf-service:50051') as channel:
        stub = retrievedoc_pb2_grpc.DocumentScorerStub(channel)
        response = stub.RetrieveDocument(retrievedoc_pb2.Query(text=query))
        print("Retrieved Document:", response.text)  # Debug print
        return response.text
    #with grpc.insecure_channel('tfidf-service:50051') as channel:
    #with grpc.insecure_channel('localhost:50051') as channel:

def generate_response(query, document):
    # Use the QA model to extract the answer from the document
    result = oracle(question=query, context=document)
    print("Generated Response:", result)  # Debug print
    return result['answer']

@app.post("/chat")
async def chat(request: ChatRequest):
    document = retrieve_document(request.query)
    response = generate_response(request.query, document)
    return {"response": response}
