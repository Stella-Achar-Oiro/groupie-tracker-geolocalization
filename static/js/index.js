function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

mapboxgl.accessToken = 'pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw';

const searchInput = document.getElementById('search-input');
const suggestionsContainer = document.getElementById('suggestions');
const creationYearSlider = document.getElementById('creation-year');
const creationYearDisplay = document.getElementById('creation-year-display');
const firstAlbumYearSlider = document.getElementById('first-album-year');
const firstAlbumYearDisplay = document.getElementById('first-album-year-display');
const memberCheckboxes = document.getElementById('member-checkboxes');
const locationCheckboxes = document.getElementById('location-checkboxes');
let allArtists = [];
let allLocations = new Set();

// Event Listeners
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

// Add event listeners for filters
creationYearSlider.addEventListener('input', updateCreationYearDisplay);
firstAlbumYearSlider.addEventListener('input', updateFirstAlbumYearDisplay);

function updateCreationYearDisplay() {
    creationYearDisplay.textContent = creationYearSlider.value;
    applyFilters();
}

function updateFirstAlbumYearDisplay() {
    firstAlbumYearDisplay.textContent = firstAlbumYearSlider.value;
    applyFilters();
}

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

function searchArtists(query = '') {
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
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        displayResults(data.artists);
        if (!allArtists.length) {
            allArtists = data.artists;
        }
        hideLoading();
    })
    .catch(error => {
        console.error('Error:', error);
        showError('An error occurred while searching for artists. Please try again later.');
        hideLoading();
    });
}

function getFilterValues() {
    return {
        creationYearMin: parseInt(document.getElementById('creation-year').value),
        creationYearMax: 2023,
        firstAlbumYearMin: parseInt(document.getElementById('first-album-year').value),
        firstAlbumYearMax: 2023,
        members: Array.from(document.querySelectorAll('#member-checkboxes input:checked')).map(cb => parseInt(cb.value)),
        locations: Array.from(document.querySelectorAll('#location-checkboxes input:checked')).map(cb => cb.value)
    };
}

function displayResults(artists) {
    const container = document.getElementById('results-container');
    container.innerHTML = '';
    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.innerHTML = `
            <img src="placeholder.jpg" data-src="${artist.image}" alt="${artist.name}" class="lazy-image">
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

function showError(message) {
    const errorElement = document.getElementById('error-message');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
    setTimeout(() => {
        errorElement.style.display = 'none';
    }, 5000);
}

// Initialize the page
window.addEventListener('load', () => {
    searchArtists('');
});