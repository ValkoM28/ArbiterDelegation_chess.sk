// Rounds Management JavaScript
let currentRounds = [];
let currentLeague = null;
let directorInfo = '';
let contactPerson = '';

// Load rounds data for a specific league
async function loadRoundsData(leagueId) {
    try {
        const response = await fetch('/get-rounds', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ leagueId: leagueId })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        currentRounds = data.rounds;
        currentLeague = data.league;
        
        // Extract director info and contact person from league
        directorInfo = `${currentLeague.directorFirstName} ${currentLeague.directorSurname}`;
        if (currentLeague.directorEmail) {
            directorInfo += ` (${currentLeague.directorEmail})`;
        }
        contactPerson = '';
        displayRoundsEditor();
        
        return data;
    } catch (error) {
        console.error('Error loading rounds data:', error);
        showStatus('Error loading rounds data: ' + error.message, 'error');
        throw error;
    }
}

// Display the rounds editor interface
function displayRoundsEditor() {
    const roundsContainer = document.getElementById('roundsEditor');
    if (!roundsContainer) {
        console.error('Rounds editor container not found');
        return;
    }

    roundsContainer.classList.remove('hidden');

    // JavaScript injected html, probably enough for the usecase
    let html = `
        <div class="mx-auto bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-6">Uprav Kolá</h2>
            
            <!-- Global Fields -->
            <div class="mb-8 p-4 bg-gray-50 rounded-lg">
                <h3 class="text-lg font-medium text-gray-700 mb-4">Základné Info</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Riaditeľ súťaže (kontakt)</label>
                        <input 
                            type="text" 
                            id="globalDirectorInfo" 
                            value="${directorInfo}"
                            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Poverený člen KR (kontakt)</label>
                        <input 
                            type="text" 
                            id="globalContactPerson" 
                            value="${contactPerson}"
                            placeholder="Enter contact person"
                            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                </div>
            </div>

            <!-- Rounds List -->
            <div class="space-y-6">
    `;

    currentRounds.forEach((round, roundIndex) => {
        html += `
            <div class="border border-gray-200 rounded-lg p-4">
                <div class="flex items-center justify-between mb-4">
                    <h4 class="text-lg font-medium text-gray-700">Kolo č. ${round.number}</h4>
                </div>
                
                <div class="space-y-3">
        `;

        // Add each match in the round
        round.matches.forEach((match, matchIndex) => {
            html += `
                <div class="p-3 bg-gray-50 rounded border">
                    <div class="grid grid-cols-1 md:grid-cols-4 gap-3 mb-3">
                        <div>
                            <label class="block text-xs font-medium text-gray-600 mb-1">Domáci</label>
                            <input 
                                type="text" 
                                id="round_${roundIndex}_match_${matchIndex}_home" 
                                value="${match.homeTeam}"
                                class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label class="block text-xs font-medium text-gray-600 mb-1">Hostia</label>
                            <input 
                                type="text" 
                                id="round_${roundIndex}_match_${matchIndex}_guest" 
                                value="${match.guestTeam}"
                                class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label class="block text-xs font-medium text-gray-600 mb-1">Dátum a Čas (RRRR/MM/DD HH:MM)</label>
                            <input 
                                type="text" 
                                id="round_${roundIndex}_match_${matchIndex}_datetime" 
                                value="${match.dateTime}"
                                class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label class="block text-xs font-medium text-gray-600 mb-1">Adresa hracej miestnosti</label>
                            <input 
                                type="text" 
                                id="round_${roundIndex}_match_${matchIndex}_address" 
                                value="${match.address}"
                                placeholder="Enter address"
                                class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                            />
                        </div>
                    </div>
                    
                    <!-- Match Arbiter Selection -->
                    <div class="mt-3 p-2 bg-blue-50 rounded">
                        <label class="block text-xs font-medium text-gray-600 mb-1">Delegovaný rozhodca:</label>
                        <div class="relative">
                            <input 
                                type="text" 
                                id="round_${roundIndex}_match_${matchIndex}_arbiter_search" 
                                placeholder="Hľadaj podľa priezviska..."
                                class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                                oninput="filterArbiters(${roundIndex}, ${matchIndex})"
                                onfocus="showArbiterDropdown(${roundIndex}, ${matchIndex})"
                                onblur="hideArbiterDropdown(${roundIndex}, ${matchIndex})"
                            />
                            <div 
                                id="round_${roundIndex}_match_${matchIndex}_arbiter_dropdown" 
                                class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg hidden max-h-60 overflow-y-auto"
                            >
                                <!-- Arbiter options will be populated here -->
                            </div>
                        </div>
                        <div id="round_${roundIndex}_match_${matchIndex}_arbiter_details" class="mt-1 text-xs text-gray-600 hidden">
                            <!-- Arbiter details will be shown here -->
                        </div>
                    </div>
                </div>
            `;
        });
        html += `

                </div>
            </div>
        `;
    });

    html += `
        <div class="flex space-x-4 justify-end">
            <button 
                id="prepareDelegationBtn"
                onclick="prepareDelegationData()"
                class="bg-orange-500 hover:bg-orange-600 text-white font-bold py-3 px-6 text-lg rounded-lg transition duration-200"
                disabled
            >
                Vytvoriť delegačné listy
            </button>
        </div>

    `; 
    roundsContainer.innerHTML = html;
    
    // Enable the prepare delegation button
    const prepareDelegationBtn = document.getElementById('prepareDelegationBtn');
    if (prepareDelegationBtn) {
        prepareDelegationBtn.disabled = false;
    }
    
    // Populate arbiter dropdowns for all matches
    populateMatchArbiterDropdowns();
}

// Save rounds data
async function saveRoundsData() {
    try {
        // Collect all the data from the form
        const updatedRounds = [];
        
        currentRounds.forEach((round, roundIndex) => {
            const updatedRound = {
                number: round.number,
                dateTime: document.getElementById(`round_${roundIndex}_datetime`).value,
                matches: []
            };

            round.matches.forEach((match, matchIndex) => {
                const updatedMatch = {
                    homeTeam: document.getElementById(`round_${roundIndex}_match_${matchIndex}_home`).value,
                    guestTeam: document.getElementById(`round_${roundIndex}_match_${matchIndex}_guest`).value,
                    dateTime: document.getElementById(`round_${roundIndex}_match_${matchIndex}_datetime`).value,
                    address: document.getElementById(`round_${roundIndex}_match_${matchIndex}_address`).value
                };
                updatedRound.matches.push(updatedMatch);
            });

            updatedRounds.push(updatedRound);
        });

        // Update global info
        directorInfo = document.getElementById('globalDirectorInfo').value;
        contactPerson = document.getElementById('globalContactPerson').value;

        // Send to backend
        const response = await fetch('/save-rounds', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ rounds: updatedRounds })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Rounds saved:', data);
        
        // Update current rounds
        currentRounds = updatedRounds;
        
        showStatus('Rounds data saved successfully!', 'success');
        
    } catch (error) {
        console.error('Error saving rounds data:', error);
        showStatus('Error saving rounds data: ' + error.message, 'error');
    }
}

// Hide rounds editor
function hideRoundsEditor() {
    const roundsContainer = document.getElementById('roundsEditor');
    if (roundsContainer) {
        roundsContainer.classList.add('hidden');
    }
}

// Show status message
function showStatus(message, type = 'info') {
    const statusElement = document.getElementById('roundsStatus');
    if (statusElement) {
        statusElement.textContent = message;
        statusElement.className = `text-sm ${type === 'error' ? 'text-red-600' : type === 'success' ? 'text-green-600' : 'text-blue-600'}`;
    }
}



// Populate arbiter dropdowns for all matches
async function populateMatchArbiterDropdowns() {
    try {
        // Get arbiters data
        const response = await fetch('/arbiters');
        const data = await response.json();
        
        if (data.arbiters && data.arbiters.length > 0) {
            // Store arbiters globally for filtering
            window.allArbiters = data.arbiters;
            
            // Populate all match arbiter dropdowns
            currentRounds.forEach((round, roundIndex) => {
                round.matches.forEach((match, matchIndex) => {
                    const dropdownElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
                    if (dropdownElement) {
                        // Populate dropdown with all arbiters initially
                        populateArbiterDropdown(roundIndex, matchIndex, data.arbiters);
                    }
                });
            });
        }
    } catch (error) {
        console.error('Error loading arbiters for match dropdowns:', error);
    }
}

// Populate arbiter dropdown with given arbiters list
function populateArbiterDropdown(roundIndex, matchIndex, arbiters) {
    const dropdownElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
    if (!dropdownElement) return;
    
    // Clear existing options
    dropdownElement.innerHTML = '';
    
    // Sort arbiters by last name first, then first name
    const sortedArbiters = arbiters.sort((a, b) => {
        const aName = `${a.LastName} ${a.FirstName}`.toLowerCase();
        const bName = `${b.LastName} ${b.FirstName}`.toLowerCase();
        return aName.localeCompare(bName);
    });
    
    // Add arbiter options
    sortedArbiters.forEach(arbiter => {
        const option = document.createElement('div');
        option.className = 'px-3 py-2 hover:bg-gray-100 cursor-pointer text-sm';
        option.textContent = `${arbiter.LastName} ${arbiter.FirstName} (${arbiter.ArbiterLevel})${arbiter.KlubName ? ` - ${arbiter.KlubName}` : ''}`;
        option.onclick = () => selectArbiter(roundIndex, matchIndex, arbiter);
        dropdownElement.appendChild(option);
    });
}

// Show arbiter dropdown
function showArbiterDropdown(roundIndex, matchIndex) {
    const dropdown = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
    if (dropdown) {
        dropdown.classList.remove('hidden');
    }
}

// Hide arbiter dropdown
function hideArbiterDropdown(roundIndex, matchIndex) {
    // Add a small delay to allow clicking on options
    setTimeout(() => {
        const dropdown = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
        if (dropdown) {
            dropdown.classList.add('hidden');
        }
    }, 200);
}

// Normalize text by removing diacritics and converting to lowercase
function normalizeText(text) {
    return text
        .toLowerCase()
        .normalize('NFD') // Decompose characters with diacritics
        .replace(/[\u0300-\u036f]/g, '') // Remove diacritic marks
        .replace(/[^\w\s]/g, '') // Remove special characters except letters, numbers, and spaces
        .trim();
}

// Filter arbiters based on search input
function filterArbiters(roundIndex, matchIndex) {
    const searchInput = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_search`);
    const searchTerm = normalizeText(searchInput.value);
    
    if (!window.allArbiters) return;
    
    // Filter arbiters by last name first, then first name
    const filteredArbiters = window.allArbiters.filter(arbiter => {
        const fullName = normalizeText(`${arbiter.LastName} ${arbiter.FirstName}`);
        const clubName = arbiter.KlubName ? normalizeText(arbiter.KlubName) : '';
        return fullName.includes(searchTerm) || clubName.includes(searchTerm);
    });
    
    populateArbiterDropdown(roundIndex, matchIndex, filteredArbiters);
    
    // Show dropdown if there are results
    const dropdown = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
    if (dropdown) {
        dropdown.classList.remove('hidden');
    }
}

// Select an arbiter
function selectArbiter(roundIndex, matchIndex, arbiter) {
    const searchInput = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_search`);
    const dropdown = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_dropdown`);
    
    // Update search input with selected arbiter
    searchInput.value = `${arbiter.LastName} ${arbiter.FirstName} (${arbiter.ArbiterLevel})${arbiter.KlubName ? ` - ${arbiter.KlubName}` : ''}`;
    
    // Hide dropdown
    dropdown.classList.add('hidden');
    
    // Store selected arbiter ID for later use
    searchInput.setAttribute('data-arbiter-id', arbiter.ArbiterId);
    
    // Show arbiter details
    showArbiterDetails(roundIndex, matchIndex, arbiter);
}

// Show arbiter details
async function showArbiterDetails(roundIndex, matchIndex, arbiter) {
    const detailsElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_details`);
    if (detailsElement) {
        const clubInfo = arbiter.KlubName ? ` - ${arbiter.KlubName}` : '';
        detailsElement.innerHTML = `
            <div class="text-xs text-gray-600">
                <strong>${arbiter.LastName} ${arbiter.FirstName}</strong> (${arbiter.ArbiterLevel})${clubInfo}
            </div>
        `;
        detailsElement.classList.remove('hidden');
    }
}

// Handle arbiter selection for a specific match (legacy function - keeping for compatibility)
async function onMatchArbiterSelected(roundIndex, matchIndex) {
    const selectElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter`);
    const detailsElement = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_details`);
    
    if (selectElement.value) {
        try {
            const response = await fetch(`/arbiters/${selectElement.value}`);
            const data = await response.json();
            
            if (data.arbiter) {
                const clubInfo = data.arbiter.KlubName ? ` - ${data.arbiter.KlubName}` : '';
                detailsElement.innerHTML = `
                    <strong>Selected:</strong> ${data.arbiter.FirstName} ${data.arbiter.LastName} (ID: ${data.arbiter.PlayerId})${clubInfo}
                `;
                detailsElement.classList.remove('hidden');
            }
        } catch (error) {
            console.error('Error fetching arbiter details:', error);
            detailsElement.innerHTML = '<span class="text-red-600">Error loading arbiter details</span>';
            detailsElement.classList.remove('hidden');
        }
    } else {
        detailsElement.classList.add('hidden');
    }
}

// Prepare PDFData array from current rounds data
function preparePDFDataArray() {
    const leagueSelect = document.getElementById('leagueSelect');
    
    const globalDirectorInfo = document.getElementById('globalDirectorInfo')?.value;
    const globalContactPerson = document.getElementById('globalContactPerson')?.value || '';
    
    // Get league name from the selected option
    const selectedLeagueOption = leagueSelect.options[leagueSelect.selectedIndex];
    const leagueName = selectedLeagueOption ? selectedLeagueOption.textContent.split(' (')[0] : '';
    const leagueYear = selectedLeagueOption ? selectedLeagueOption.textContent.match(/\((.+?)\)/)?.[1] || '' : '';
    
    const pdfDataArray = [];
    
    // Create PDFData for each match
    currentRounds.forEach((round, roundIndex) => {
        round.matches.forEach((match, matchIndex) => {
            // Get current form data (including any user edits)
            const homeTeam = document.getElementById(`round_${roundIndex}_match_${matchIndex}_home`)?.value || match.homeTeam;
            const guestTeam = document.getElementById(`round_${roundIndex}_match_${matchIndex}_guest`)?.value || match.guestTeam;
            const dateTime = document.getElementById(`round_${roundIndex}_match_${matchIndex}_datetime`)?.value || match.dateTime;
            const address = document.getElementById(`round_${roundIndex}_match_${matchIndex}_address`)?.value || match.address;
            
            // Get arbiter info from the searchable input
            const arbiterSearchInput = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_search`);
            const selectedArbiterId = arbiterSearchInput ? arbiterSearchInput.getAttribute('data-arbiter-id') : '';
            
            // Get arbiter details from the details element
            const arbiterDetails = document.getElementById(`round_${roundIndex}_match_${matchIndex}_arbiter_details`);
            let arbiterName = '';
            let arbiterId = '';
            let arbiterClubName = '';
            
            if (arbiterDetails && arbiterDetails.textContent) {
                // Extract name from the details text (format: "LastName FirstName (Level) - Club")
                const detailsText = arbiterDetails.textContent;
                const nameMatch = detailsText.match(/<strong>(.+?)<\/strong>/);
                if (nameMatch) {
                    arbiterName = nameMatch[1];
                    arbiterId = selectedArbiterId;
                    
                    // Extract club name if present
                    const clubMatch = detailsText.match(/- (.+?)<\/div>/);
                    if (clubMatch) {
                        arbiterClubName = clubMatch[1];
                    }
                }
            }
            
            // If no arbiter details found, try to get from search input
            if (!arbiterName && arbiterSearchInput && arbiterSearchInput.value) {
                const inputValue = arbiterSearchInput.value;
                // Parse the format: "LastName FirstName (Level) - Club"
                const arbiterMatch = inputValue.match(/^(.+?) \((.+?)\)(?: - (.+))?$/);
                if (arbiterMatch) {
                    arbiterName = arbiterMatch[1];
                    arbiterId = selectedArbiterId;
                    arbiterClubName = arbiterMatch[3] || '';
                }
            }
            
            const pdfData = {
                league: {
                    name: leagueName,
                    year: leagueYear
                },
                director: {
                    contact: globalDirectorInfo
                },
                arbiter: {
                    firstName: arbiterName.split(' ')[0] || '',
                    lastName: arbiterName.split(' ').slice(1).join(' ') || '',
                    playerId: arbiterId,
                    clubName: '' // just because of the updated ArbiterInfo in backend it wont run without this line :D 
                },
                match: {
                    homeTeam: homeTeam,
                    guestTeam: guestTeam,
                    dateTime: dateTime,
                    address: address
                },
                contactPerson: globalContactPerson
            };
            
            pdfDataArray.push(pdfData);
        });
    });
    
    return pdfDataArray;
}

// Prepare delegation data and send to backend
async function prepareDelegationData() {
    const leagueSelect = document.getElementById('leagueSelect');
    const roundsStatus = document.getElementById('roundsStatus');
    
    if (!leagueSelect.value) {
        roundsStatus.innerHTML = '<span class="text-red-600">✗ Please select a league first</span>';
        return;
    }
    
    if (!currentRounds || currentRounds.length === 0) {
        roundsStatus.innerHTML = '<span class="text-red-600">✗ Please load rounds data first</span>';
        return;
    }
    
    try {
        // Prepare PDFData array from current rounds data
        const pdfDataArray = preparePDFDataArray();
        
        // Validate that we have data
        if (pdfDataArray.length === 0) {
            roundsStatus.innerHTML = '<span class="text-red-600">✗ No match data found</span>';
            return;
        }
        
        // Check if all matches have arbiters assigned (for testing, we'll allow missing arbiters)
        const missingArbiters = pdfDataArray.filter(data => !data.arbiter.firstName || !data.arbiter.lastName);
        if (missingArbiters.length > 0) {
            console.log(`Warning: ${missingArbiters.length} matches missing arbiter assignment - proceeding anyway for testing`);
        }
        
        // Show loading state
        roundsStatus.innerHTML = '<span class="text-blue-600">⏳ Generating PDFs and creating zip file...</span>';
        
        // Send to backend
        const response = await fetch('/delegate-arbiters', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(pdfDataArray)
        });
        
        if (!response.ok) {
            // Handle error responses
            let errorMessage = `Server error: ${response.status} ${response.statusText}`;
            try {
                const errorData = await response.json();
                errorMessage = `Server error: ${errorData.error || 'Unknown error'}`;
            } catch (jsonError) {
                // If JSON parsing fails, we'll use the default error message
                console.warn('Could not parse error response as JSON:', jsonError);
            }
            throw new Error(errorMessage);
        }
        
        // Check if response is a file download (zip)
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/zip')) {
            // Handle zip file download
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            
            // Get filename from Content-Disposition header or use default
            const contentDisposition = response.headers.get('content-disposition');
            let filename = 'delegacne_listy.zip';
            if (contentDisposition) {
                const filenameMatch = contentDisposition.match(/filename="(.+)"/);
                if (filenameMatch) {
                    filename = filenameMatch[1];
                }
            }
            
            // Create download link and trigger download
            const a = document.createElement('a');
            a.href = url;
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
            
            roundsStatus.innerHTML = `
                <span class="text-green-600">✓ PDFs generated and zip file downloaded successfully!</span><br>
                <span class="text-sm text-gray-600">Count: ${pdfDataArray.length} items</span><br>
                <span class="text-sm text-gray-600">File: ${filename}</span>
            `;
        } else {
            // Fallback for JSON response (shouldn't happen with current backend)
            const result = await response.json();
            console.log('Delegation data sent:', result);
            
            roundsStatus.innerHTML = `
                <span class="text-green-600">✓ ${result.message}</span><br>
                <span class="text-sm text-gray-600">Count: ${result.count} items</span><br>
                <span class="text-sm text-gray-600">Check server console for detailed output</span>
            `;
        }
        
    } catch (error) {
        console.error('Error preparing delegation data:', error);
        roundsStatus.innerHTML = `<span class="text-red-600">✗ Error: ${error.message}</span>`;
    }
}
