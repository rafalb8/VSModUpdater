// Theme toggle
function toggleTheme() {
    const html = document.documentElement;
    const currentTheme = html.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    
    const btn = document.querySelector('.theme-toggle');
    btn.textContent = newTheme === 'dark' ? '☀️' : '🌙';
}

// Load saved theme or detect system preference
(function() {
    let theme = localStorage.getItem('theme');
    
    // If no saved preference, detect system preference
    if (!theme) {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        theme = prefersDark ? 'dark' : 'light';
    }
    
    document.documentElement.setAttribute('data-theme', theme);
    const btn = document.querySelector('.theme-toggle');
    if (btn) {
        btn.textContent = theme === 'dark' ? '☀️' : '🌙';
    }
})();

// Search/filter function
function filterTable() {
    const input = document.getElementById('search');
    const filter = input.value.toLowerCase();
    const table = document.getElementById('modTable');
    const rows = table.getElementsByTagName('tr');
    
    for (let i = 1; i < rows.length; i++) {
        const row = rows[i];
        const cells = row.getElementsByTagName('td');
        let found = false;
        
        for (let j = 0; j < cells.length; j++) {
            const cell = cells[j];
            if (cell) {
                const text = cell.textContent || cell.innerText;
                if (text.toLowerCase().indexOf(filter) > -1) {
                    found = true;
                    break;
                }
            }
        }
        
        row.style.display = found ? '' : 'none';
    }
}

// Table sorting
let sortOrder = {};

function sortTable(columnIndex) {
    const table = document.getElementById('modTable');
    const tbody = table.getElementsByTagName('tbody')[0];
    const rows = Array.from(tbody.getElementsByTagName('tr'));
    
    // Toggle sort order
    if (!sortOrder[columnIndex]) {
        sortOrder[columnIndex] = 'asc';
    } else if (sortOrder[columnIndex] === 'asc') {
        sortOrder[columnIndex] = 'desc';
    } else {
        sortOrder[columnIndex] = 'asc';
    }
    
    // Sort rows
    rows.sort((a, b) => {
        const aValue = a.getElementsByTagName('td')[columnIndex].textContent.trim();
        const bValue = b.getElementsByTagName('td')[columnIndex].textContent.trim();
        
        let comparison = 0;
        if (aValue < bValue) {
            comparison = -1;
        } else if (aValue > bValue) {
            comparison = 1;
        }
        
        return sortOrder[columnIndex] === 'asc' ? comparison : -comparison;
    });
    
    // Reorder DOM
    rows.forEach(row => tbody.appendChild(row));
    
    // Update sort indicators
    const headers = table.getElementsByTagName('th');
    for (let i = 0; i < headers.length; i++) {
        headers[i].classList.remove('sorted');
        const indicator = headers[i].querySelector('.sort-indicator');
        if (indicator) {
            if (i === columnIndex) {
                headers[i].classList.add('sorted');
                indicator.textContent = sortOrder[columnIndex] === 'asc' ? '↑' : '↓';
            } else {
                indicator.textContent = '↕';
            }
        }
    }
}

// Column toggle
function toggleColumn(columnIndex) {
    const table = document.getElementById('modTable');
    const headers = table.getElementsByTagName('th');
    const rows = table.getElementsByTagName('tr');
    
    // Toggle header
    headers[columnIndex].classList.toggle('hidden');
    
    // Toggle cells in all rows
    for (let i = 1; i < rows.length; i++) {
        const cells = rows[i].getElementsByTagName('td');
        if (cells[columnIndex]) {
            cells[columnIndex].classList.toggle('hidden');
        }
    }
}
