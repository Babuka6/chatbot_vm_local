#!/bin/bash

# Define the output JSON file
JSON_FILE="random_requests.json"

# Check if we have a word source: either 'shuf' with a system dictionary or a ~/words.txt fallback
if ! command -v shuf &>/dev/null && [ ! -f ~/words.txt ]; then
  echo "Error: neither 'shuf' command nor ~/words.txt file is available for random word generation."
  exit 1
fi

# Start the JSON array
echo "[" > "$JSON_FILE"

# Generate 100 random JSON entries, each containing 3 words
for i in $(seq 1 100); do
    # Get 3 random lines and merge them into a single line separated by spaces
    # 'paste -sd " "' merges the 3 lines from shuf into one space-delimited line
    RANDOM_WORDS="$(
      shuf -n3 /usr/share/dict/words 2>/dev/null || \
      shuf -n3 ~/words.txt
    )"
    RANDOM_WORDS="$(echo "$RANDOM_WORDS" | paste -sd " ")"

    # Remove any possible carriage returns and extra spaces
    RANDOM_WORDS="$(echo "$RANDOM_WORDS" | tr -d '\r' | sed 's/  */ /g; s/^ //; s/ $//')"

    # Escape any double quotes to ensure valid JSON
    ESCAPED_WORDS="$(echo "$RANDOM_WORDS" | sed 's/"/\\"/g')"

    # Build JSON entry, avoiding a trailing comma on the last item
    if [ "$i" -eq 100 ]; then
        echo "  {\"text\": \"$ESCAPED_WORDS\"}" >> "$JSON_FILE"
    else
        echo "  {\"text\": \"$ESCAPED_WORDS\"}," >> "$JSON_FILE"
    fi
done

# Close the JSON array
echo "]" >> "$JSON_FILE"

echo "✅ Generated valid JSON requests in $JSON_FILE"
