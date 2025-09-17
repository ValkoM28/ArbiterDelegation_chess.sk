# Chess Arbiter Delegation Generator - Release v1.0.0

## 🎉 First Release!

This is the initial release of the Chess Arbiter Delegation Generator, a web application designed to automate the creation of delegation forms for chess arbiters in Slovakia.

## ✨ Features

- **🔄 Automatic Data Loading**: Fetches arbiters and leagues from chess.sk API
- **📊 Excel Processing**: Downloads and processes tournament data from chess-results.com
- **📄 PDF Generation**: Creates delegation forms with arbiter and match information
- **🌐 Web Interface**: Simple, user-friendly interface in Slovak
- **📦 Batch Processing**: Handles multiple arbiters and matches simultaneously
- **💾 ZIP Downloads**: Automatically packages generated PDFs for easy download

## 🚀 Quick Start

1. **Download** the appropriate ZIP file for your operating system and architecture
2. **Extract** the files to a folder on your computer
3. **Run** the application:
   - **Windows x64**: Double-click `server-windows-amd64.exe`
   - **Windows ARM64**: Double-click `server-windows-arm64.exe`
   - **Linux x64**: Run `./server-linux-amd64` in terminal
   - **Linux ARM64**: Run `./server-linux-arm64` in terminal
   - **macOS x64**: Run `./server-macos-amd64` in terminal
   - **macOS ARM64 (Apple M1/M2/M3)**: Run `./server-macos-arm64` in terminal
4. **Open** your browser and go to `http://localhost:8080`

## 📱 Supported Platforms

- **Windows**: x64 (Intel/AMD) and ARM64 (Snapdragon)
- **Linux**: x64 (Intel/AMD) and ARM64 (Raspberry Pi, ARM servers)
- **macOS**: x64 (Intel) and ARM64 (Apple M1/M2/M3 processors)

## 📋 Usage

1. Click **"Načítaj dáta z chess.sk"** to load arbiters and leagues
2. Select an arbiter and league from the dropdowns
3. Optionally download tournament data from chess-results.com
4. Generate PDF delegation forms
5. Download the ZIP file with all generated PDFs

## 🏗️ Technical Details

- **Language**: Go 1.19+
- **Web Framework**: Gin
- **PDF Processing**: PDFCPU
- **Excel Processing**: Excelize
- **Frontend**: HTML, CSS, JavaScript (Tailwind CSS)


## 📦 What's Included

Each release package contains:
- Pre-compiled binary for your operating system
- Web interface files
- PDF templates
- README.md with usage instructions
- MIT License

## 📞 Support

For questions or support, please create an issue on GitHub.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Made for the Slovak Chess Federation** 🇸🇰 ♟️

*Release Date: September 17, 2025*
