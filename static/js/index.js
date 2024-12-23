// Global variables
let allArtists = [];
let allLocations = new Set();
let memberCounts = new Set();
let showingFavorites = false;

// DOM Elements
const searchInput = document.getElementById('search-input');
const suggestionsContainer = document.getElementById('suggestions');
const favoritesBtn = document.getElementById('show-favorites');
const locationSearch = document.getElementById('location-search');

// Loading functions
function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function showError(message) {
    const errorElement = document.getElementById('error-message');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
    setTimeout(() => {
        errorElement.style.display = 'none';
    }, 5000);
}

// Search and Suggestions
searchInput.addEventListener('input', () => {
    const query = searchInput.value;
    if (query.length >= 1) {
        fetch(`/api/suggestions?q=${encodeURIComponent(query)}`)
            .then(response => response.json())
            .then(suggestions => displaySuggestions(suggestions))
            .catch(error => console.error('Error:', error));
    } else {
        suggestionsContainer.innerHTML = '';
    }
});

searchInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        const query = searchInput.value;
        suggestionsContainer.innerHTML = '';
        searchArtists(query);
    }
});

function displaySuggestions(suggestions) {
    suggestionsContainer.innerHTML = '';
    suggestions.forEach(suggestion => {
        const div = document.createElement('div');
        div.className = 'suggestion-item';
        div.textContent = `${suggestion.text} (${suggestion.type})`;
        div.onclick = () => {
            searchInput.value = suggestion.text;
            suggestionsContainer.innerHTML = '';
            searchArtists(suggestion.text);
        };
        suggestionsContainer.appendChild(div);
    });
}

// Filter Initialization and Updates
function initializeFilters(artists) {
    // Initialize member counts
    artists.forEach(artist => {
        memberCounts.add(artist.members.length);
    });

    // Populate member checkboxes
    const memberCheckboxes = document.getElementById('member-checkboxes');
    const sortedCounts = Array.from(memberCounts).sort((a, b) => a - b);
    
    memberCheckboxes.innerHTML = sortedCounts.map(count => `
        <label class="checkbox-label">
            <input type="checkbox" name="members" value="${count}">
            ${count} ${count === 1 ? 'member' : 'members'}
        </label>
    `).join('');

    // Add 10+ option if needed
    if (sortedCounts[sortedCounts.length - 1] < 10) {
        memberCheckboxes.innerHTML += `
            <label class="checkbox-label">
                <input type="checkbox" name="members" value="10+">
                10+ members
            </label>
        `;
    }

    // Initialize locations using cached data from the search response
    Promise.all(artists.map(artist => 
        fetch(`/api/artist/${artist.id}`)
            .then(response => response.json())
            .then(data => data.locations)
    ))
    .then(locationArrays => {
        // Flatten all location arrays and add to Set
        locationArrays.forEach(locations => {
            locations.forEach(location => {
                allLocations.add(location.address);
            });
        });
        
        updateLocationCheckboxes(Array.from(allLocations));
    })
    .catch(error => {
        console.error('Error fetching locations:', error);
        showError('Error loading location filters');
    });

    // Set year ranges
    initializeYearRanges(artists);
}

function initializeYearRanges(artists) {
    const creationYearMin = Math.min(...artists.map(a => a.creationDate));
    const creationYearMax = Math.max(...artists.map(a => a.creationDate));
    const albumYearMin = Math.min(...artists.map(a => parseInt(a.firstAlbum.split('-')[2])));
    const albumYearMax = Math.max(...artists.map(a => parseInt(a.firstAlbum.split('-')[2])));

    document.getElementById('creation-year-min').value = creationYearMin;
    document.getElementById('creation-year-max').value = creationYearMax;
    document.getElementById('album-year-min').value = albumYearMin;
    document.getElementById('album-year-max').value = albumYearMax;
}

function updateLocationCheckboxes(locations) {
    const locationCheckboxes = document.getElementById('location-checkboxes');
    locationCheckboxes.innerHTML = locations.sort().map(location => `
        <label class="checkbox-label">
            <input type="checkbox" name="locations" value="${location}">
            ${location}
        </label>
    `).join('');
}

// Filter Event Listeners
document.querySelectorAll('.year-input').forEach(input => {
    input.addEventListener('change', applyFilters);
});

locationSearch.addEventListener('input', (e) => {
    const searchTerm = e.target.value.toLowerCase();
    const filteredLocations = Array.from(allLocations).filter(location =>
        location.toLowerCase().includes(searchTerm)
    );
    updateLocationCheckboxes(filteredLocations);
});

document.querySelectorAll('input[name="members"], input[name="locations"]').forEach(checkbox => {
    checkbox.addEventListener('change', applyFilters);
});

// Favorites Handling
favoritesBtn.addEventListener('click', () => {
    showingFavorites = !showingFavorites;
    favoritesBtn.classList.toggle('active');
    favoritesBtn.innerHTML = showingFavorites ? 
        '<i class="fas fa-star"></i> Show All' : 
        '<i class="fas fa-star"></i> Show Favorites';
    
    if (showingFavorites) {
        const favorites = JSON.parse(localStorage.getItem('favorites') || '[]');
        const favoriteArtists = allArtists.filter(artist => 
            favorites.includes(artist.id.toString())
        );
        displayResults(favoriteArtists);
    } else {
        searchArtists(searchInput.value);
    }
});

// Filter Functions
function getFilterValues() {
    const memberCheckboxes = Array.from(document.querySelectorAll('input[name="members"]:checked'));
    const memberValues = memberCheckboxes.map(cb => {
        return cb.value === '10+' ? 10 : parseInt(cb.value);
    });

    return {
        creationYearMin: parseInt(document.getElementById('creation-year-min').value),
        creationYearMax: parseInt(document.getElementById('creation-year-max').value),
        firstAlbumYearMin: parseInt(document.getElementById('album-year-min').value),
        firstAlbumYearMax: parseInt(document.getElementById('album-year-max').value),
        members: memberValues,
        locations: Array.from(document.querySelectorAll('input[name="locations"]:checked')).map(cb => cb.value)
    };
}

function applyFilters() {
    searchArtists(searchInput.value);
}

// Clear Filters
function clearFilters() {
    // Reset year inputs
    const defaultYears = {
        'creation-year-min': 1950,
        'creation-year-max': 2023,
        'album-year-min': 1950,
        'album-year-max': 2023
    };

    Object.entries(defaultYears).forEach(([id, value]) => {
        document.getElementById(id).value = value;
    });

    // Clear member checkboxes
    document.querySelectorAll('#member-checkboxes input[type="checkbox"]')
        .forEach(checkbox => checkbox.checked = false);

    // Clear location search and checkboxes
    document.getElementById('location-search').value = '';
    document.querySelectorAll('#location-checkboxes input[type="checkbox"]')
        .forEach(checkbox => checkbox.checked = false);

    // Clear search input and suggestions
    searchInput.value = '';
    suggestionsContainer.innerHTML = '';

    // Reset favorites if showing
    if (showingFavorites) {
        showingFavorites = false;
        favoritesBtn.classList.remove('active');
        favoritesBtn.innerHTML = '<i class="fas fa-star"></i> Show Favorites';
    }

    // Refresh results
    searchArtists('');
}

// Add clear filters event listener
document.getElementById('clear-filters').addEventListener('click', clearFilters);

// Search and Display Functions
function searchArtists(query = '') {
    const container = document.getElementById('results-container');
    showLoading();
    const filters = getFilterValues();
    
    fetch(`/api/search?q=${encodeURIComponent(query)}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(filters)
    })
    .then(response => {
        if (!response.ok) {
            throw { response };
        }
        return response.json();
    })
    .then(data => {
        displayResults(data.artists);
        if (!allArtists.length) {
            allArtists = data.artists;
            initializeFilters(data.artists);
        }
        hideLoading();
    })
    .catch(error => handleError(error, container));
}

function displayResults(artists) {
    const container = document.getElementById('results-container');
    container.innerHTML = '';
    
    if (artists.length === 0) {
        container.innerHTML = '<div class="no-results">No artists found matching your criteria</div>';
        return;
    }

    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.innerHTML = `
            <img src="/static/images/placeholder.jpg" data-src="${artist.image}" alt="${artist.name}" class="lazy-image">
            <h3>${artist.name}</h3>
            <p><i class="fas fa-calendar-alt"></i> Created: ${artist.creationDate}</p>
            <p><i class="fas fa-compact-disc"></i> First Album: ${artist.firstAlbum}</p>
        `;
        card.onclick = () => {
            window.location.href = `/artist/${artist.id}`;
        };
        container.appendChild(card);
    });
    lazyLoadImages();
}

// Image Lazy Loading
function lazyLoadImages() {
    const images = document.querySelectorAll('.lazy-image');
    const options = {
        root: null,
        rootMargin: '0px',
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.remove('lazy-image');
                observer.unobserve(img);
            }
        });
    }, options);

    images.forEach(img => observer.observe(img));
}

window.addEventListener('online', () => {
    // Refresh the page when connection is restored
    window.location.reload();
});

window.addEventListener('offline', () => {
    const container = document.getElementById('results-container'); // or 'artist-details'
    createOfflineUI(container);
});

// Initialize the page
window.addEventListener('load', () => {
    searchArtists('');
});