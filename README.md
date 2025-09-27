# Generátor delegácie rozhodcov (for English scroll down)

Webová aplikácia na automatizáciu vytvárania delegačných listov pre šachových rozhodcov na Slovensku. Systém sa integruje s API chess.sk a chess-results.com. 

Oficiálne používaná na delegáciu rozhodcov majstrovských súťaží SŠZ (Extraliga, 1. Liga, 2. Liga)

## Spustenie programu

### Možnosť 1: Stiahnutie exe súboru 

#### Predpoklady
- Internetové pripojenie (pre prístup k API)

#### Stiahnutie a spustenie
1. **Stiahnite si .exe súbor** zo stránky vydaní (Releases). Ktorý zip si stiahnete záleží od Vášho PC. Najbežnejšie však bude potrebná verzia amd64. Ak máte procesor, ktorý nie je AMD ani Intel, bude potrebná verzia arm64. 
2. **Rozbaľte zip súbor**
3. **Spustite aplikáciu**: server.exe
4. **Otvorte prehliadač** a prejdite na: `http://localhost:8080`

### Možnosť 2: Kompilácia zo zdrojového kódu (Pre vývojárov)

#### Predpoklady
- Go 1.19 alebo novší
- Internetové pripojenie (pre prístup k API)

#### Inštalácia a spustenie
```bash
git clone <repository-url>
cd ssz_delegovanie_rozhodcov

go mod tidy

go run cmd/server/main.go
```

### Prístup k aplikácii
Otvorte si prehliadač a prejdite na: `http://localhost:8080`

## 📋 Postup používania

### 1. Načítanie dát z chess.sk
- Kliknite na **"Načítaj dáta z chess.sk"**
- Zadajte rok sezóny (napr. "2024" pre sezónu 2024/2025)
- Počkajte na načítanie rozhodcov a líg

### 2. Výber rozhodcu a ligy
- Vyberte rozhodcu z rozbaľovacieho menu **"Arbiter"**
- Vyberte ligu z rozbaľovacieho menu **"Liga"**
- Formulár sa automaticky vyplní relevantnými informáciami

### 3. Spracovanie turnajových dát (Voliteľné)
- Kliknite na **"Stiahnuť Excel"** na stiahnutie turnajových dát z chess-results.com
- Systém extrahuje informácie o kolách a zápasoch

### 4. Generovanie delegácie
- Kliknite na **"Delegovať rozhodcov"** na generovanie PDF formulárov
- Stiahnite si ZIP súbor so všetkými vygenerovanými PDF


Pozrite si [DOCS.md](DOCS.md) pre podrobnú technickú dokumentáciu a pokyny pre príspevky.

## 📄 Licencia

Tento projekt je licencovaný pod MIT licenciou - pozrite si súbor [LICENSE](LICENSE) pre podrobnosti.

## 📞 Podpora

Pre otázky alebo podporu, prosím vytvorte issue.

---

**Vytvorené pre Slovenský šachový zväz** 🇸🇰 ♟️

---

# Chess Arbiter Delegation Generator

A web application for automating the creation of delegation forms for chess arbiters in Slovakia. The system integrates with chess.sk API and chess-results.com to streamline the delegation process.

## 🚀 Quick Start

### Option 1: Download Pre-compiled Binary (Recommended for Users)

#### Prerequisites
- Internet connection (for API access)

#### Download & Run
1. **Download the binary** from the releases page
2. **Extract the files** to a folder on your computer
3. **Run the application**:
   ```bash
   # On Windows
   server.exe
   
   # On Linux/macOS
   ./server
   ```
4. **Open your browser** and go to: `http://localhost:8080`

#### What's Included
- Pre-compiled binary for your operating system
- All necessary web assets and templates
- Ready to run without any setup

### Option 2: Build from Source (For Developers)

#### Prerequisites
- Go 1.19 or later
- Internet connection (for API access)

#### Installation & Run
```bash
# Clone the repository
git clone <repository-url>
cd ssz_delegovanie_rozhodcov

# Install dependencies
go mod tidy

# Start the application
go run cmd/server/main.go
```

### Access the Application
Open your browser and go to: `http://localhost:8080`

## 📋 How to Use

### 1. Load Data from chess.sk
- Click **"Načítaj dáta z chess.sk"**
- Enter the season year (e.g., "2024" for 2024/2025 season)
- Wait for arbiters and leagues to load

### 2. Select Arbiter and League
- Choose an arbiter from the **"Arbiter"** dropdown
- Choose a league from the **"Liga"** dropdown
- The form will populate with relevant information

### 3. Process Tournament Data (Optional)
- Click **"Stiahnuť Excel"** to download tournament data from chess-results.com
- The system extracts round information and match details

### 4. Generate Delegation Forms
- Click **"Pripraviť PDF dáta"** to prepare the data
- Click **"Delegovať rozhodcov"** to generate PDF forms
- Download the ZIP file with all generated PDFs

## ✨ Features

- **🔄 Automatic Data Loading**: Fetches arbiters and leagues from chess.sk API
- **📊 Excel Processing**: Downloads and processes tournament data from chess-results.com
- **📄 PDF Generation**: Creates delegation forms with arbiter and match information
- **🌐 Web Interface**: Simple, user-friendly interface in Slovak
- **📦 Batch Processing**: Handles multiple arbiters and matches simultaneously
- **💾 ZIP Downloads**: Automatically packages generated PDFs for easy download


## 🤝 Contributing

We welcome contributions from the Slovak Chess Federation community!

### For Users
- Report bugs and suggest improvements
- Test new features and provide feedback
- Share your experience with the application

### For Developers
- Fix bugs and add new features
- Improve documentation and code quality
- Add tests and enhance performance

See [DOCS.md](DOCS.md) for detailed technical documentation and contribution guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For questions or support, please submit an issue. 

---

**Made for the Slovak Chess Federation**
