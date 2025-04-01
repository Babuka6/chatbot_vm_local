package main

import (
	"context"
	"log"
	"math"
	"net"
	"strings"
	"unicode"

	pb "tfidf-service/proto/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDocumentScorerServer
}

func (s *server) RetrieveDocument(ctx context.Context, query *pb.Query) (*pb.Document, error) {
	// Example documents
	documents := []string{
		"Document 1: Our shop offers an extensive range of high-quality products, including state-of-the-art electronics, stylish and comfortable clothing, and a variety of home goods designed to enhance your living space. Whether you are looking for the latest smartphones with cutting-edge features, premium laptops for work and gaming, or home appliances that make daily tasks easier, our store is committed to providing top-notch items from reputable brands. We also carry a wide selection of fashion apparel, including seasonal wear, activewear, and accessories, ensuring you always have access to trendy and comfortable outfits. Additionally, our home goods department features everything from kitchen essentials and furniture to decorative items that help personalize your home. We strive to offer competitive pricing, frequent promotions, and exclusive deals to make shopping with us a rewarding experience.",

		"Document 2: If you are not fully satisfied with your purchase, you can return any unused items in their original, unopened packaging within 30 days of the purchase date for a full refund or exchange. Our hassle-free return policy is designed to ensure customer satisfaction, allowing you to shop with confidence. To initiate a return, simply visit our online returns portal, where you can enter your order details and receive step-by-step instructions on how to send your item back. If you prefer, you can also visit one of our physical store locations, where our staff will be happy to assist you with processing your return. Please note that certain items, such as personal care products, perishable goods, and digital downloads, may be ineligible for return due to hygiene and licensing restrictions. Refunds are typically processed within 5-7 business days after we receive and inspect the returned item.",

		"Document 3: We are pleased to offer fast and reliable shipping to our customers, with free standard shipping on all orders over $50. Our delivery service ensures that your purchases arrive safely and on time, typically within 3-5 business days for domestic shipments. We also provide expedited shipping options for those who need their items sooner, including two-day and overnight delivery services at an additional cost. Our logistics team works with leading courier partners to provide efficient and secure shipping solutions. Upon dispatch, you will receive a tracking number via email, allowing you to monitor your order's journey in real time. If you experience any shipping delays or issues, our dedicated customer support team is available to assist you with resolving any concerns. International shipping is available for select locations, with estimated delivery times varying based on destination and customs processing.",

		"Document 4: All of our products come with a standard one-year warranty that covers manufacturing defects and performance issues. This warranty ensures that your purchase is protected against unexpected failures, giving you peace of mind. If you experience any defects within the warranty period, you may be eligible for a free repair, replacement, or refund, depending on the product category and issue reported. Our warranty program includes coverage for electronics, appliances, and select home goods, with an option to extend the warranty period for an additional fee. To make a warranty claim, simply contact our support team with your order details and a description of the issue. We will guide you through the process and, if necessary, provide a prepaid return label for you to send the item back for inspection. Please note that accidental damage, misuse, and unauthorized repairs may void the warranty.",

		"Document 5: For any issues, questions, or concerns regarding your order, our dedicated customer support team is available to assist you via multiple communication channels. You can reach us by email at support@example.com, where our representatives typically respond within 24 hours. Additionally, our toll-free customer service hotline, 1-800-123-4567, is available during business hours for immediate assistance with order inquiries, troubleshooting, and general support. We also offer a live chat feature on our website, allowing you to connect with a support agent in real time. Our help center contains a comprehensive list of frequently asked questions (FAQs) covering topics such as account management, payment methods, order tracking, returns, and product warranties. Our goal is to provide exceptional customer service and ensure that every shopping experience is smooth and hassle-free. We value customer feedback and continuously strive to improve our services based on user input.",
	}

	// Find the most relevant document using TF-IDF and cosine similarity
	topDocument := RetrieveDocument(query.Text, documents)
	return &pb.Document{Text: topDocument}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDocumentScorerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// RetrieveDocument returns the most relevant document for the query
func RetrieveDocument(query string, documents []string) string {
	// Step 1: Tokenize the query and documents
	queryTokens := Tokenize(query)
	documentTokens := make([][]string, len(documents))
	for i, doc := range documents {
		documentTokens[i] = Tokenize(doc)
	}

	// Step 2: Compute TF-IDF scores for the query and documents
	queryTfidf := ComputeQueryTFIDF(queryTokens, documentTokens)
	documentTfidfs := ComputeDocumentTFIDFs(documentTokens)

	// Step 3: Calculate cosine similarity between the query and each document
	similarities := make([]float64, len(documents))
	for i, docTfidf := range documentTfidfs {
		similarities[i] = CosineSimilarity(queryTfidf, docTfidf)
	}

	// Step 4: Find the document with the highest similarity score
	maxScore := -1.0
	topDocumentIndex := 0
	for i, score := range similarities {
		if score > maxScore {
			maxScore = score
			topDocumentIndex = i
		}
	}

	return documents[topDocumentIndex]
}

// Tokenize splits a document into words
func Tokenize(document string) []string {
	// Convert to lowercase
	document = strings.ToLower(document)

	// Remove punctuation
	document = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, document)

	// Split into words
	return strings.Fields(document)
}

// ComputeQueryTFIDF calculates the TF-IDF vector for the query
func ComputeQueryTFIDF(queryTokens []string, documentTokens [][]string) map[string]float64 {
	// Calculate TF for the query
	queryTf := TermFrequency(queryTokens)

	// Calculate IDF for the query terms
	idf := InverseDocumentFrequency(documentTokens)

	// Compute TF-IDF for the query
	queryTfidf := make(map[string]float64)
	for word := range queryTf {
		queryTfidf[word] = queryTf[word] * idf[word]
	}

	return queryTfidf
}

// ComputeDocumentTFIDFs calculates the TF-IDF vectors for all documents
func ComputeDocumentTFIDFs(documentTokens [][]string) []map[string]float64 {
	idf := InverseDocumentFrequency(documentTokens)
	documentTfidfs := make([]map[string]float64, len(documentTokens))

	for i, tokens := range documentTokens {
		tf := TermFrequency(tokens)
		tfidf := make(map[string]float64)

		for word := range tf {
			tfidf[word] = tf[word] * idf[word]
		}

		documentTfidfs[i] = tfidf
	}

	return documentTfidfs
}

// TermFrequency calculates the term frequency for a list of tokens
func TermFrequency(tokens []string) map[string]float64 {
	tf := make(map[string]float64)

	// Count word occurrences
	for _, word := range tokens {
		tf[word]++
	}

	// Normalize by the total number of words
	totalWords := float64(len(tokens))
	for word := range tf {
		tf[word] /= totalWords
	}

	return tf
}

// InverseDocumentFrequency calculates the IDF for each word in the corpus
func InverseDocumentFrequency(documentTokens [][]string) map[string]float64 {
	idf := make(map[string]float64)
	totalDocuments := float64(len(documentTokens))

	// Count the number of documents containing each word
	for _, tokens := range documentTokens {
		uniqueWords := make(map[string]bool)
		for _, word := range tokens {
			if !uniqueWords[word] {
				idf[word]++
				uniqueWords[word] = true
			}
		}
	}

	// Calculate IDF
	for word := range idf {
		idf[word] = math.Log(totalDocuments / idf[word])
	}

	return idf
}

// CosineSimilarity calculates the cosine similarity between two TF-IDF vectors
func CosineSimilarity(vectorA, vectorB map[string]float64) float64 {
	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0

	// Compute dot product and magnitudes
	for word, scoreA := range vectorA {
		scoreB := vectorB[word]
		dotProduct += scoreA * scoreB
		magnitudeA += scoreA * scoreA
	}
	for _, scoreB := range vectorB {
		magnitudeB += scoreB * scoreB
	}

	// Avoid division by zero
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
}
