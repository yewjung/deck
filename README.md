# Deck API

This is a public API built in Go that allows you to create a new deck and draw cards from it. The API provides two endpoints: `/newdeck` and `/card`.

## Installation

To run this project on your local computer, follow these steps:

1. Clone this repository using the following command:
   ```
   git clone https://github.com/yewjung/deck.git
   ```

2. Navigate to the project directory:
   ```
   cd deck
   ```

3. Start the Docker containers using `docker-compose`:
   ```
   docker-compose up -d
   ```

   This command will build and run the necessary containers defined in the `docker-compose.yml` file.

4. Once the containers are up and running, the API will be accessible at `http://localhost:8080`.

## API Endpoints

### 1. Create a New Deck

**Endpoint:** `/newdeck`\
**Method:** `GET`\
**Description:** This endpoint creates a new deck on the server and returns a unique deck ID.

**Example Request:**
```
GET /newdeck
```

**Example Response:**
```
abcdef123456
```

### 2. Draw a Card from a Deck

**Endpoint:** `/card?deckid={deckId}`\
**Method:** `GET`\
**Description:** This endpoint allows you to draw a single card from a specific deck identified by its `deckId`.

**Example Request:**
```
GET /card?deckid=abcdef123456
```

**Example Response:**
```
{
    "rank":"6",
    "suit":"S"
}
```

## Dependencies

This project relies on the following dependencies:

- Go (1.16 or higher)
- Docker
- Docker Compose

Make sure you have these dependencies installed on your system before running the project.

