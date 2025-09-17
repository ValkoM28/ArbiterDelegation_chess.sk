# Chess Arbiter Delegation Generator
Disclaimer, this documentation was purely generated, if you find any flaw please submit an issue and I will correct it. 
## Table of Contents
1. [Quick Start Guide](#quick-start-guide)
2. [How to Use the Application](#how-to-use-the-application)
3. [User Interface Guide](#user-interface-guide)
4. [Technical Documentation](#technical-documentation)
5. [Development Guide](#development-guide)
6. [Contributing](#contributing)

## Quick Start Guide

### What is this application?
The Chess Arbiter Delegation Generator is a web application designed to automate the creation of delegation forms for chess arbiters in Slovakia. It integrates with the chess.sk API to fetch arbiter and league data, processes tournament information from chess-results.com, and generates PDF delegation forms.

### Key Features
- **Automatic Data Loading**: Fetches arbiters and leagues from chess.sk API
- **Tournament Processing**: Downloads and processes tournament data from chess-results.com
- **PDF Generation**: Creates delegation forms with arbiter and match information
- **Easy-to-Use Interface**: Simple web interface for data management
- **Batch Processing**: Handles multiple arbiters and matches simultaneously

### Getting Started
1. **Start the application** (see Development Guide for setup)
2. **Open your browser** and go to `http://localhost:8080`
3. **Load data** from chess.sk by clicking "Načítaj dáta z chess.sk"
4. **Select arbiters and leagues** from the dropdown menus
5. **Generate PDFs** for delegation forms

## How to Use the Application

### Step 1: Load External Data
1. Open the web application in your browser
2. Click the "Načítaj dáta z chess.sk" button
3. Enter the season start year (e.g., "2024")
4. Wait for the data to load (arbiters and leagues will be fetched from chess.sk)

### Step 2: Select Arbiter and League
1. Choose an arbiter from the "Arbiter" dropdown
2. Choose a league from the "Liga" dropdown
3. The system will automatically populate the form with relevant information

### Step 3: Process Tournament Data (Optional)
1. If you need to process tournament rounds from chess-results.com:
   - Click "Stiahnuť Excel" to download tournament data
   - The system will extract round information and match details

### Step 4: Generate Delegation Forms
1. Click "Pripraviť PDF dáta" to prepare the PDF data
2. Click "Delegovať rozhodcov" to generate the PDF delegation forms
3. The system will create a ZIP file with all generated PDFs for download

## User Interface Guide

### Main Interface Elements

#### Data Loading Section
- **"Načítaj dáta z chess.sk"** button: Loads arbiters and leagues from the official Slovak chess API
- **Season Year Input**: Enter the chess season year (e.g., 2024 for 2024/2025 season)

#### Selection Dropdowns
- **Arbiter Dropdown**: Select from available arbiters loaded from chess.sk
- **League Dropdown**: Select from available leagues for the specified season

#### Action Buttons
- **"Stiahnuť Excel"**: Downloads and processes tournament data from chess-results.com
- **"Pripraviť PDF dáta"**: Prepares the data for PDF generation
- **"Delegovať rozhodcov"**: Generates and downloads PDF delegation forms

#### Status Indicators
- The interface shows whether arbiters and leagues data has been loaded
- Success/error messages appear for each operation

### Workflow Tips
- Always load data from chess.sk first before making selections
- Ensure you have the correct season year for league data
- Check that the selected league has a valid chess-results.com link for Excel processing
- Generated PDFs are automatically packaged in a ZIP file for easy download

## Technical Documentation

### Architecture

The application follows a clean, modular architecture with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │   HTTP Server   │    │   External APIs │
│   (HTML/JS)     │◄──►│   (Gin Router)  │◄──►│   (chess.sk)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  Application    │
                    │  Layer (App)    │
                    └─────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
            ┌─────────────┐    ┌─────────────┐
            │ Data Layer  │    │ PDF Layer   │
            │ (Storage)   │    │ (Generator) │
            └─────────────┘    └─────────────┘
                    │
                    ▼
            ┌─────────────┐
            │ Excel Layer │
            │ (Processor) │
            └─────────────┘
```

## Package Structure

### `/cmd/server`
**Purpose**: Application entry point and server configuration

**Files**:
- `main.go`: Server initialization, route registration, and startup

**Key Functions**:
- `main()`: Entry point that sets up Gin router, serves static files, and starts HTTP server

### `/internal/app`
**Purpose**: Application layer coordinating between packages and handling HTTP requests

**Files**:
- `app.go`: Core application structure and dependency management
- `handlers.go`: HTTP request handlers and API endpoints

**Key Types**:
- `App`: Main application struct with storage dependency

**Key Functions**:
- `New()`: Creates new application instance
- `RegisterRoutes()`: Registers all HTTP endpoints
- `LoadArbiters()`: Loads arbiters from chess.sk API
- `LoadLeagues()`: Loads leagues from chess.sk API

### `/internal/data`
**Purpose**: Data models, storage, and data processing

**Files**:
- `models.go`: Data structures and validation logic
- `storage.go`: In-memory session storage and data processing

**Key Types**:
- `SessionData`: Thread-safe in-memory storage
- `PDFData`: Structured data for PDF generation
- `Arbiter`: Chess arbiter information from API
- `League`: Chess league information from API
- `Round`: Tournament round with matches
- `MatchInfo`: Individual match details

**Key Functions**:
- `NewSessionData()`: Creates new storage instance
- `LoadData()`: Loads data from external APIs
- `GetAllArbiters()`: Retrieves all loaded arbiters
- `GetAllLeagues()`: Retrieves all loaded leagues
- `ProcessData[T]()`: Generic data processing function

### `/internal/excel`
**Purpose**: Excel file processing and tournament data extraction

**Files**:
- `processor.go`: Excel download and data extraction

**Key Functions**:
- `DownloadChessResultsExcel()`: Downloads Excel files from chess-results.com
- `ExtractTournamentIDFromLeague()`: Extracts tournament ID from league data
- `CleanupTempFile()`: Removes temporary Excel files

### `/internal/pdf`
**Purpose**: PDF generation, form filling, and file management

**Files**:
- `generator.go`: PDF form filling and generation
- `helpers.go`: Data conversion utilities
- `mapper.go`: Field mapping for PDF forms
- `validator.go`: PDF data validation
- `zipping.go`: ZIP file creation for batch downloads

**Key Functions**:
- `FillForm()`: Fills PDF forms with data
- `PreparePDFDataFromArbiterAndLeague()`: Converts API data to PDF format
- `fromArbiter()`: Converts arbiter data for PDF
- `fromLeague()`: Converts league data for PDF

## API Endpoints

### Data Loading
- `POST /load-external-data`: Load arbiters and leagues from chess.sk API
- `GET /external-data/:type`: Get raw external data (arbiters/leagues)

### Data Retrieval
- `GET /arbiters`: Get all loaded arbiters
- `GET /arbiters/:id`: Get specific arbiter by ID
- `GET /leagues`: Get all loaded leagues
- `GET /leagues/:id`: Get specific league by ID

### PDF Generation
- `POST /prepare-pdf-data`: Prepare PDF data for specific arbiter/league
- `POST /delegate-arbiters`: Generate PDFs for multiple arbiters

### Excel Processing
- `POST /download-excel`: Download and process Excel from chess-results.com
- `POST /get-rounds`: Extract round information from Excel files

## Data Models

### Core Data Structures

#### Arbiter
```go
type Arbiter struct {
    ArbiterId    string `json:"ArbiterId"`    // Unique identifier
    PlayerId     string `json:"PlayerId"`     // Player ID in chess system
    FideId       string `json:"FideId"`       // FIDE ID
    LastName     string `json:"LastName"`     // Surname
    FirstName    string `json:"FirstName"`    // First name
    ValidTo      string `json:"ValidTo"`      // License validity end date
    Licencia     string `json:"Licencia"`     // License number
    KlubId       string `json:"KlubId"`       // Club identifier
    KlubName     string `json:"KlubName"`     // Club name
    IsActive     bool   `json:"IsActive"`     // Active status
    ArbiterLevel string `json:"ArbiterLevel"` // Certification level
}
```

#### League
```go
type League struct {
    LeagueId          string `json:"leagueId"`          // Unique identifier
    SaisonName        string `json:"saisonName"`        // Season name
    LeagueName        string `json:"leagueName"`        // Display name
    ChessResultsLink  string `json:"chessResultsLink"`  // Tournament URL
    DirectorId        string `json:"directorId"`        // Director ID
    DirectorSurname   string `json:"directorSurname"`   // Director surname
    DirectorFirstName string `json:"directorFirstName"` // Director first name
    DirectorEmail     string `json:"directorEmail"`     // Director email
}
```

#### PDFData
```go
type PDFData struct {
    Arbiter       ArbiterData  // Arbiter information
    League        LeagueData   // League information
    Match         MatchData    // Match details
    Director      DirectorData // Director information
    ContactPerson string       // Contact person
}
```

### Data Flow

1. **Data Loading**: External APIs → SessionData storage
2. **Data Processing**: Raw API data → Structured models
3. **PDF Generation**: Structured data → PDF forms
4. **File Management**: Generated PDFs → ZIP archives

## Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)

### Dependencies
- **Gin**: HTTP web framework
- **Excelize**: Excel file processing
- **PDFCPU**: PDF form filling
- **UUID**: Unique identifier generation

### External Services
- **chess.sk API**: Arbiters and leagues data
- **chess-results.com**: Tournament Excel files

## Development Guide

### Prerequisites
- Go 1.19 or later
- Git

### Setup
```bash
# Clone repository
git clone <repository-url>
cd ssz_delegovanie_rozhodcov

# Install dependencies
go mod tidy

# Run development server
go run cmd/server/main.go
```

### Project Structure Guidelines
- `/cmd`: Application entry points
- `/internal`: Private application code
- `/web`: Static web assets
- `/templates`: PDF templates

### Code Organization
- **Packages**: Single responsibility principle
- **Functions**: Clear naming and documentation
- **Error Handling**: Consistent error propagation
- **Testing**: Unit tests for core functionality

### Adding New Features
1. Define data models in `/internal/data`
2. Implement business logic in `/internal/app`
3. Add API endpoints in handlers
4. Update frontend in `/web`
5. Add tests and documentation

## Deployment

### Build
```bash
# Build binary
go build -o server cmd/server/main.go

# Run binary
./server
```



## Contributing

We welcome contributions from members of the Slovak Chess Federation and the broader chess community! This project is designed to serve the Slovak chess community, and we appreciate any help in making it better.

### Who Can Contribute
- **Slovak Chess Federation members**: Your input is especially valuable as you understand the specific needs of Slovak chess
- **Chess arbiters and organizers**: Your practical experience helps improve the application
- **Developers in the chess community**: Technical contributions are always welcome
- **Anyone interested in chess technology**: Fresh perspectives help us grow

### How to Contribute

#### For Non-Technical Users
- **Report bugs**: If something doesn't work as expected, let us know
- **Suggest improvements**: Share ideas for new features or better workflows
- **Test new features**: Help us ensure everything works correctly
- **Provide feedback**: Your user experience helps us improve the interface

#### For Developers
- **Code contributions**: Fix bugs, add features, improve performance
- **Documentation**: Help make the code more understandable
- **Testing**: Add tests to ensure reliability
- **Code review**: Help maintain code quality

### Getting Started
1. **Fork the repository** or contact the maintainers
2. **Create a feature branch** for your changes
3. **Make your changes** following Go conventions
4. **Test your changes** thoroughly
5. **Submit a pull request** with a clear description

### Code Style
- Follow Go conventions and best practices
- Use meaningful variable names
- Add comprehensive documentation
- Write unit tests for new functionality
- Keep functions focused and readable

### Contact
For questions about contributing or if you need help getting started, please contact the Slovak Chess Federation technical team.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

*Last updated: $(date)*
*Version: 1.0.0*
